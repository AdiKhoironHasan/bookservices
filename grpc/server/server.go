package server

import (
	"github.com/AdiKhoironHasan/bookservices-books/config"
	"github.com/AdiKhoironHasan/bookservices-books/domain/service"
	"github.com/AdiKhoironHasan/bookservices-books/grpc/client"
	"github.com/AdiKhoironHasan/bookservices-books/grpc/handler"
	"github.com/AdiKhoironHasan/bookservices-books/grpc/interceptor"
	"github.com/AdiKhoironHasan/bookservices-books/proto/book"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is struct to hold any dependencies used for server
type Server struct {
	config     *config.Config
	repo       *service.Repositories
	grpcClient *client.GRPCClient
}

type ServerGrpcOption func(*Server)

// NewGRPCServer is constructor
// func NewGRPCServer(conf *config.Config, repo *service.Repositories) *Server {
func NewGRPCServer(options ...ServerGrpcOption) *Server {
	server := &Server{}

	for _, option := range options {
		option(server)
	}

	return server
}

// Run is a method gRPC server
func (s *Server) Run(port int) error {
	//server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.AuthorizationInterceptor))
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryLoggerServerInterceptor(),
			interceptor.UnaryAuthServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.StreamLoggerServerInterceptor(),
			interceptor.StreamAuthServerInterceptor(),
		),
	)

	handlers := handler.NewHandler(s.config, s.repo, s.grpcClient)

	// register from proto
	book.RegisterBookServiceServer(server, handlers)

	// register reflection
	reflection.Register(server)

	return RunGRPCServer(server, port)
}
