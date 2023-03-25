package assembler

import (
	"github.com/AdiKhoironHasan/bookservices-books/domain/entity"
	"github.com/AdiKhoironHasan/bookservices-books/proto/book"
	protoUser "github.com/AdiKhoironHasan/bookservices-protobank/proto/user"
)

func ToResponseBookList(users []*protoUser.User, books []entity.Book) []*book.Book {
	dataUser := map[int64]string{}
	dataMap := []*book.Book{}

	for _, val := range users {
		if _, ok := dataUser[val.Id]; !ok {
			dataUser[val.Id] = val.Name
		}
	}

	for _, val := range books {
		if _, ok := dataUser[val.AuthorId]; ok {
			book := &book.Book{
				Id:          val.ID,
				AuthorId:    val.AuthorId,
				AuthorName:  dataUser[val.AuthorId],
				Title:       val.Title,
				Description: val.Description,
				CreatedAt:   val.CreatedAt.String(),
				UpdatedAt:   val.UpdatedAt.String(),
			}
			dataMap = append(dataMap, book)
		}
	}

	return dataMap
}
