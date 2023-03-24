package client

import (
	"flag"
	"fmt"

	"github.com/AdiKhoironHasan/bookservices/config"
	"github.com/AdiKhoironHasan/bookservices/grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ServerHost string
	ServerPort int
	DSN        string
	Address    *string
)

// NewGRPCConn is a constructor
func NewGRPCConn_Book(cfg *config.Config) (*grpc.ClientConn, error) {
	ServerHost = "localhost"
	ServerPort = cfg.GRPCPort
	DSN = fmt.Sprintf("%s:%d", ServerHost, ServerPort)
	Address = flag.String("localAddress", DSN, "The address to connect books service")

	flag.Parse()

	conn, err := grpc.Dial(*Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(interceptor.UnaryAuthClientInterceptor()),
		grpc.WithStreamInterceptor(interceptor.StreamAuthClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewGRPCConn_User(cfg *config.Config) (*grpc.ClientConn, error) {
	ServerHost = cfg.Dependency.UserServices.AppHostGRPC
	ServerPort = cfg.Dependency.UserServices.AppPortGRPC
	DSN = fmt.Sprintf("%s:%d", ServerHost, ServerPort)
	Address = flag.String("userAddress", DSN, "The address to connect user service")

	flag.Parse()

	conn, err := grpc.Dial(*Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(interceptor.UnaryAuthClientInterceptor()),
		grpc.WithStreamInterceptor(interceptor.StreamAuthClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
