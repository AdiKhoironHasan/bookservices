package client

import (
	protoUser "github.com/AdiKhoironHasan/bookservices-protobank/proto/user"
	"google.golang.org/grpc"
)

// GRPCClient is a struct
type GRPCClient struct {
	// book book.BookServiceClient
	User protoUser.UserServiceClient
}

// NewGRPCClient is constructor
func NewGRPCClient(connUser grpc.ClientConnInterface) *GRPCClient {
	return &GRPCClient{
		// book: book.NewBookServiceClient(connBook),
		User: protoUser.NewUserServiceClient(connUser),
	}
}
