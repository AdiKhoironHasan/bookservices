package service

import (
	"github.com/AdiKhoironHasan/bookservices/domain/repository"
	"github.com/AdiKhoironHasan/bookservices/infrastructure/persistence"
	"gorm.io/gorm"
)

// Repositories is a struct
type Repositories struct {
	User      repository.UserRepositoryInterface
	AuthToken repository.AuthTokenRepository
	HttpLog   repository.HttpLogRepositoryInterface
	DB        *gorm.DB
}

// NewDBService is constructor
func NewDBService(db *gorm.DB) *Repositories {
	return &Repositories{
		User:    persistence.NewUserRepository(db),
		HttpLog: persistence.NewHttpLogRepository(db),
		DB:      db,
	}
}
