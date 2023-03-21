package client

import (
	"context"

	"github.com/AdiKhoironHasan/bookservices/proto/book"
)

func (r GRPCClient) Ping(ctx context.Context) (*book.BookResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	books, err := r.book.List(ctx, &book.BookRequest{})
	if err != nil {
		return nil, err
	}

	return books, nil
}
