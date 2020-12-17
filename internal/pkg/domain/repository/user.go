package repository

import "forum-api/internal/pkg/domain/models"

type User interface {
	InsertInto(user *models.User) error

	GetByNickname(user *models.User) error
	GetByNicknameOrEmail(user *models.User) ([]models.User, error)
	GetStatus(status *models.Status) error

	Update(user *models.User) error

	DeleteAll() error

	Prepare() error
}
