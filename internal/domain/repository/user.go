package repository

import "forum-api/internal/domain/models"

type User interface {
	Insert(user *models.User) error
	Update(user *models.User) error
	DeleteAll() error
	GetStatus(status *models.Status) error
	SelectByNickname(nickname string) (models.User, error)
	SelectNicknameWithCase(nickname string) (string, error)
	SelectByEmailOrNickname(nickname string, email string) (models.UserSlice, error)
}
