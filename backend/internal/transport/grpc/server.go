package grpc

import (
	"net"

	pb "github.com/mushibana/glacier/backend/proto"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
}

func NewServer(tasksHandler *TasksHandler) *Server {
	server := grpc.NewServer()

	pb.RegisterTasksServer(server, tasksHandler)

	return &Server{
		server: server,
	}
}

func (s *Server) Serve(listener net.Listener) error {
	return s.server.Serve(listener)
}

func (s *Server) GracefulStop() {
	s.server.GracefulStop()
}

func (s *Server) Stop() {
	s.server.Stop()
}
