package handler

import (
	"context"

	"github.com/AdiKhoironHasan/bookservices/domain/entity"
	"github.com/AdiKhoironHasan/bookservices/proto/book"
)

// List is a function
func (c *Handler) List(_ context.Context, _ *book.BookListReq) (*book.BookListRes, error) {
	books := []entity.Book{}

	rows, err := c.repo.DB.Table("public.books").Select("id, title, description, created_at, updated_at").Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		book := entity.Book{}
		rows.Scan(&book.ID, &book.Title, &book.Description, &book.CreatedAt, &book.UpdatedAt)
		books = append(books, book)
	}

	ch := make(chan []*book.Book)
	defer close(ch)

	go func(books []entity.Book, ch chan<- []*book.Book) {
		value := []*book.Book{}
		for _, val := range books {
			value = append(value, &book.Book{
				Id:          val.ID,
				Title:       val.Title,
				Description: val.Description,
				CreatedAt:   val.CreatedAt.String(),
				UpdatedAt:   val.UpdatedAt.String(),
			})
		}

		ch <- value
	}(books, ch)

	bookResp := <-ch

	res := &book.BookListRes{
		Books: bookResp,
	}

	return res, nil
}

func (c *Handler) Store(ctx context.Context, bookReq *book.BookStoreReq) (*book.BookStoreRes, error) {
	bookEntity := entity.Book{
		Title:       bookReq.Title,
		Description: bookReq.Description,
	}

	err := c.repo.DB.Create(&bookEntity).Error
	if err != nil {
		return nil, err
	}

	return &book.BookStoreRes{}, nil
}

func (c *Handler) Detail(ctx context.Context, bookReq *book.BookDetailReq) (*book.BookDetailRes, error) {
	bookEntity := entity.Book{}

	err := c.repo.DB.Where("id = ?", bookReq.Id).First(&bookEntity).Error
	if err != nil {
		return nil, err
	}

	return &book.BookDetailRes{
		Book: &book.Book{
			Id:          bookEntity.ID,
			Title:       bookEntity.Title,
			Description: bookEntity.Description,
			CreatedAt:   bookEntity.CreatedAt.String(),
			UpdatedAt:   bookEntity.UpdatedAt.String(),
		},
	}, nil
}
