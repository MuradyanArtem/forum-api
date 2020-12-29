package app

import (
	"forum-api/internal/domain/models"
	"forum-api/internal/domain/repository"
	"forum-api/internal/infrastructure"

	"github.com/pkg/errors"
)

var (
	UserNotUpdated = errors.New("User cannot be updated")
)

type User struct {
	user repository.User
}

func newUser(user repository.User) *User {
	return &User{
		user,
	}
}

func (u *User) CreateUser(user *models.User) ([]models.User, error) {

	if err := u.user.InsertInto(user); err != nil {
		users, e := u.user.GetByNicknameOrEmail(user)
		if e != nil {
			return nil, errors.Wrap(err, e.Error())
		}
		return users, err
	}

	return nil, nil
}

func (u *User) GetUser(user *models.User) error {
	if err := u.user.GetByNickname(user); err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateUser(user *models.User) error {
	uInfo := *user
	if err := u.user.GetByNickname(&uInfo); err != nil {
		return errors.Wrap(infrastructure.UserNotExist, err.Error())
	}

	if user.Email == "" {
		user.Email = uInfo.Email
	}

	if user.About == "" {
		user.About = uInfo.About
	}

	if user.Fullname == "" {
		user.Fullname = uInfo.Fullname
	}

	if err := u.user.Update(user); err != nil {
		return errors.Wrap(infrastructure.UserNotUpdated, err.Error())
	}

	return nil
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
