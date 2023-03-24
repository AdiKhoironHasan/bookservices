package client

import (
	"github.com/AdiKhoironHasan/bookservices/proto/book"
	"google.golang.org/grpc"
)

// GRPCClient is a struct
type GRPCClient struct {
	book book.BookServiceClient
}

// NewGRPCClient is constructor
func NewGRPCClient(conn grpc.ClientConnInterface) *GRPCClient {
	return &GRPCClient{
		book: book.NewBookServiceClient(conn),
	}
}
