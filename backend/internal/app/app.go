package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"

	"github.com/mushibana/glacier/backend/internal/config"
	"github.com/mushibana/glacier/backend/internal/repository/sqlite"
	sqliteDB "github.com/mushibana/glacier/backend/internal/repository/sqlite"
	"github.com/mushibana/glacier/backend/internal/service"
	"github.com/mushibana/glacier/backend/internal/transport/grpc"
	grpcTransport "github.com/mushibana/glacier/backend/internal/transport/grpc"
)

type App struct {
	logger *slog.Logger
	db     *sql.DB
	grpc   *grpcTransport.Server
	cfg    config.Config
}

func New(cfg config.Config, logger *slog.Logger) (*App, error) {
	db, err := sqliteDB.New(cfg.SQLitePath)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	taskRepo := sqlite.NewTaskRepositoryImplementation(db)
	taskService := service.NewTaskServiceImplementation(taskRepo)
	tasksHandler := grpc.NewTasksHandler(taskService)
	grpcServer := grpcTransport.NewServer(tasksHandler)

	return &App{
		logger: logger,
		db:     db,
		grpc:   grpcServer,
		cfg:    cfg,
	}, nil
}

func (a *App) Close() error {
	a.logger.Info("closing application")

	if err := a.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.cfg.GRPCAddress)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", a.cfg.GRPCAddress, err)
	}

	defer listener.Close()

	errCh := make(chan error, 1)

	go func() {
		a.logger.Info("grpc server listening",
			"address", listener.Addr().String(),
		)

		if err := a.grpc.Serve(listener); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("shutdown signal received")

	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("grpc server: %w", err)
		}
	}

	a.shutdown()

	return nil
}

func (a *App) shutdown() {
	a.logger.Info("shutting down grpc server")

	a.grpc.GracefulStop()

	a.logger.Info("grpc server stopped")
}
