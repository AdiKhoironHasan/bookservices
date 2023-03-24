package handler

import (
	"github.com/AdiKhoironHasan/bookservices-books/config"
	"github.com/AdiKhoironHasan/bookservices-books/domain/service"
	"github.com/AdiKhoironHasan/bookservices-books/grpc/client"
	"github.com/AdiKhoironHasan/bookservices-books/proto/book"
)

// Interface is an interface
type Interface interface {
	// interface of grpc handler
	book.BookServiceServer
}

// Handler is struct
type Handler struct {
	config     *config.Config
	repo       *service.Repositories
	grpcClient *client.GRPCClient

	book.UnimplementedBookServiceServer
}

// NewHandler is a constructor
func NewHandler(conf *config.Config, repo *service.Repositories, grpcClient *client.GRPCClient) *Handler {
	return &Handler{
		config:     conf,
		repo:       repo,
		grpcClient: grpcClient,
	}
}

var _ Interface = &Handler{}
