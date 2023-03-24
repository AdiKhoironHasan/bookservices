package client

import (
	protoUser "github.com/AdiKhoironHasan/bookservice-protobank/proto/user"
	"github.com/AdiKhoironHasan/bookservices/proto/book"
	"google.golang.org/grpc"
)

// GRPCClient is a struct
type GRPCClient struct {
	book book.BookServiceClient
	User protoUser.UserServiceClient
}

// NewGRPCClient is constructor
func NewGRPCClient(connBook grpc.ClientConnInterface, connUser grpc.ClientConnInterface) *GRPCClient {
	return &GRPCClient{
		book: book.NewBookServiceClient(connBook),
		User: protoUser.NewUserServiceClient(connUser),
	}
}
