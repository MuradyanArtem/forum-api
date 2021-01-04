package app

import (
	"forum-api/internal/domain/models"
	"forum-api/internal/domain/repository"
)

type User struct {
	user repository.User
}

func newUser(user repository.User) *User {
	return &User{
		user,
	}
}

func (u *User) Insert(user *models.User) error {
	return u.user.Insert(user)
}

func (u *User) Update(user *models.User) error {
	return u.user.Update(user)
}

func (u *User) SelectByNickname(nickname string) (models.User, error) {
	return u.user.SelectByNickname(nickname)
}

func (u *User) SelectByEmailOrNickname(nickname string, email string) (models.UserSlice, error) {
	return u.user.SelectByEmailOrNickname(nickname, email)
}

func (u *User) SelectNicknameWithCase(nickname string) (string, error) {
	return u.user.SelectNicknameWithCase(nickname)
}

func (u *User) DeleteAll() error {
	if err := u.user.DeleteAll(); err != nil {
		return err
	}
	return nil
}

func (u *User) GetStatus(status *models.Status) error {
	if err := u.user.GetStatus(status); err != nil {
		return err
	}
	return nil
}
