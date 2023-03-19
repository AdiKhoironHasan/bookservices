package repository

import "github.com/AdiKhoironHasan/bookservices/domain/entity"

type AuthTokenRepository interface {
	CreateAuthToken(*entity.AuthToken, *entity.User) (*entity.AuthToken, error)
	GetAuthToken(int) (*entity.AuthToken, error)
	GetAuthTokenByToken(string) (*entity.AuthToken, error)
}
