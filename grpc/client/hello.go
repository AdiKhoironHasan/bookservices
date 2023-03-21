package client

import (
	"context"

	"github.com/AdiKhoironHasan/bookservices/proto/hello"
)

// Ping is a method
func (r GRPCClient) Ping(ctx context.Context) (*hello.PingReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ping, err := r.hello.Ping(ctx, &hello.PingRequest{})
	if err != nil {
		return nil, err
	}

	return ping, nil
}
