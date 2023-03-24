package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/AdiKhoironHasan/bookservices/domain/entity"
	"github.com/AdiKhoironHasan/bookservices/proto/book"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// List is a function
func (c *Handler) List(_ context.Context, bookReq *book.BookListReq) (*book.BookListRes, error) {
	books := []entity.Book{}

	rows, err := c.repo.DB.Model(&entity.Book{}).Where(&entity.Book{
		Title: bookReq.Title,
	}).Select("id, title, description, created_at, updated_at").Rows()
	if err != nil {
		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
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

	return &book.BookListRes{
		Books: <-ch,
	}, nil
}

func (c *Handler) Store(ctx context.Context, bookReq *book.BookStoreReq) (*book.BookStoreRes, error) {
	bookEntity := entity.Book{
		Title:       bookReq.Title,
		Description: bookReq.Description,
	}

	err := c.repo.DB.Create(&bookEntity).Error
	if err != nil {
		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	return &book.BookStoreRes{}, nil
}

func (c *Handler) Detail(ctx context.Context, bookReq *book.BookDetailReq) (*book.BookDetailRes, error) {
	bookEntity := entity.Book{}

	err := c.repo.DB.First(&bookEntity, bookReq.Id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.New(http.StatusNotFound, "record not found").Err()
		}

		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
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

func (c *Handler) Update(ctx context.Context, bookReq *book.BookUpdateReq) (*book.BookUpdateRes, error) {
	err := c.repo.DB.First(&entity.Book{}, bookReq.Id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.New(http.StatusNotFound, "record not found").Err()
		}

		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	err = c.repo.DB.Model(&entity.Book{ID: bookReq.Id}).Updates(&entity.Book{
		Title:       bookReq.Title,
		Description: bookReq.Description,
	}).Error
	if err != nil {
		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	return &book.BookUpdateRes{}, nil
}

func (c *Handler) Delete(ctx context.Context, bookReq *book.BookDeleteReq) (*book.BookDeleteRes, error) {
	err := c.repo.DB.First(&entity.Book{}, bookReq.Id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.New(http.StatusNotFound, "record not found").Err()
		}

		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	err = c.repo.DB.Delete(&entity.Book{}, bookReq.Id).Error
	if err != nil {
		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	return &book.BookDeleteRes{}, nil
}
