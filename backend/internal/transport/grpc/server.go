package grpc

import (
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
}

func NewServer() *Server {
	server := grpc.NewServer()

	// pb.RegisterTasksServer(server, tasksHandler)

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
