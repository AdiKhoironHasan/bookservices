package client

import (
	"google.golang.org/grpc"
)

// GRPCClient is a struct
type GRPCClient struct {
}

// NewGRPCClient is constructor
func NewGRPCClient(conn grpc.ClientConnInterface) *GRPCClient {
	return &GRPCClient{}
}
