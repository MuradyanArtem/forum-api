package usecases

import "forum-api/internal/domain/models"

type User interface {
	DeleteAll() error
	Insert(user *models.User) error
	Update(user *models.User) error
	GetStatus(status *models.Status) error
	SelectByNickname(nickname string) (models.User, error)
	SelectNicknameWithCase(nickname string) (string, error)
	SelectByEmailOrNickname(nickname string, email string) (models.UserSlice, error)
}
