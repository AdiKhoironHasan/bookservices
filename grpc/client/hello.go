package client

import (
	"context"

	"github.com/AdiKhoironHasan/bookservices/proto/book"
)

// Ping is a method
func (r GRPCClient) List(ctx context.Context) (*book.BookResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	books, err := r.book.List(ctx, &book.BookRequest{})
	if err != nil {
		return nil, err
	}

	return books, nil
}
