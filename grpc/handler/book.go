package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/AdiKhoironHasan/bookservices/domain/entity"
	"github.com/AdiKhoironHasan/bookservices/proto/book"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	protoUser "github.com/AdiKhoironHasan/bookservice-protobank/proto/user"
)

func (c *Handler) Ping(ctx context.Context, bookReq *book.PingReq) (*book.PingRes, error) {
	var now string
	err := c.repo.DB.Raw("select now ()").Scan(&now).Error
	if err != nil {
		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	return &book.PingRes{
		Message: now,
	}, nil
}

// List is a function
func (c *Handler) List(ctx context.Context, bookReq *book.BookListReq) (*book.BookListRes, error) {
	books := []entity.Book{}

	users, err := c.grpcClient.User.List(ctx, &protoUser.UserListReq{Role: "author"})
	if err != nil {
		return nil, err
	}

	rows, err := c.repo.DB.WithContext(ctx).Model(&entity.Book{}).Where(&entity.Book{
		Title: bookReq.Title,
	}).Select("id, author_id, title, description, created_at, updated_at").Rows()
	if err != nil {
		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	defer rows.Close()
	for rows.Next() {
		book := entity.Book{}
		rows.Scan(&book.ID, &book.AuthorId, &book.Title, &book.Description, &book.CreatedAt, &book.UpdatedAt)
		books = append(books, book)
	}

	// ch := make(chan []*book.Book)
	// defer close(ch)
	// go func(books []entity.Book, ch chan<- []*book.Book) {
	// 	value := []*book.Book{}
	// 	for _, val := range books {
	// 		value = append(value, &book.Book{
	// 			Id:          val.ID,
	// 			Title:       val.Title,
	// 			Description: val.Description,
	// 			CreatedAt:   val.CreatedAt.String(),
	// 			UpdatedAt:   val.UpdatedAt.String(),
	// 		})
	// 	}

	// 	ch <- value
	// }(books, ch)
	// booksData := <-ch

	dataMap, dataUser := func(users []*protoUser.User, books []entity.Book) (map[int64][]*book.Book, map[int64]string) {
		dataUser := map[int64]string{}
		dataMap := map[int64][]*book.Book{}

		for _, val := range users {
			if _, ok := dataMap[val.Id]; !ok {
				dataMap[val.Id] = []*book.Book{}
				dataUser[val.Id] = val.Name
			}
		}

		for _, val := range books {
			if _, ok := dataMap[val.AuthorId]; ok {
				book := &book.Book{
					Id:          val.ID,
					Title:       val.Title,
					Description: val.Description,
					CreatedAt:   val.CreatedAt.String(),
					UpdatedAt:   val.UpdatedAt.String(),
				}

				dataMap[val.AuthorId] = append(dataMap[val.AuthorId], book)
			}
		}

		return dataMap, dataUser
	}(users.Users, books)

	dataRes := []*book.BookAndAuthor{}
	for key, val := range dataMap {
		if len(val) > 0 {
			dataRes = append(dataRes, &book.BookAndAuthor{
				AuthorId:   key,
				AuthorName: dataUser[key],
				Books:      val,
			})
		}
	}

	return &book.BookListRes{
		BookAndAuthor: dataRes,
	}, nil
}

func (c *Handler) Store(ctx context.Context, bookReq *book.BookStoreReq) (*book.BookStoreRes, error) {
	bookEntity := entity.Book{
		AuthorId:    bookReq.GetAuthorId(),
		Title:       bookReq.GetTitle(),
		Description: bookReq.GetDescription(),
	}

	_, err := c.grpcClient.User.Detail(ctx, &protoUser.UserDetailReq{Id: bookReq.AuthorId})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	err = c.repo.DB.WithContext(ctx).Create(&bookEntity).Error
	if err != nil {
		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	return &book.BookStoreRes{}, nil
}

func (c *Handler) Detail(ctx context.Context, bookReq *book.BookDetailReq) (*book.BookDetailRes, error) {
	bookEntity := entity.Book{}

	err := c.repo.DB.WithContext(ctx).First(&bookEntity, bookReq.Id).Error
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
	err := c.repo.DB.WithContext(ctx).First(&entity.Book{}, bookReq.Id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.New(http.StatusNotFound, "record not found").Err()
		}

		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	err = c.repo.DB.WithContext(ctx).Model(&entity.Book{ID: bookReq.Id}).Updates(&entity.Book{
		AuthorId:    bookReq.AuthorId,
		Title:       bookReq.Title,
		Description: bookReq.Description,
	}).Error
	if err != nil {
		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	return &book.BookUpdateRes{}, nil
}

func (c *Handler) Delete(ctx context.Context, bookReq *book.BookDeleteReq) (*book.BookDeleteRes, error) {
	err := c.repo.DB.WithContext(ctx).First(&entity.Book{}, bookReq.Id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.New(http.StatusNotFound, "record not found").Err()
		}

		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	err = c.repo.DB.WithContext(ctx).Delete(&entity.Book{}, bookReq.Id).Error
	if err != nil {
		return nil, status.New(http.StatusInternalServerError, err.Error()).Err()
	}

	return &book.BookDeleteRes{}, nil
}
