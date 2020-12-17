package app

import (
	"forum-api/internal/pkg/domain/models"
	"forum-api/internal/pkg/domain/repository"

	"github.com/pkg/errors"
)

type User struct {
	user repository.User
}

func NewUserApp(user repository.User) *User {
	return &User{
		user,
	}
}

func (u *User) CreateUser(user *models.User) ([]models.User, error) {
	err := u.user.InsertInto(user)
	if err != nil {
		users, e := userApp.userRepo.GetByNicknameOrEmail(user)
		err = e
		utils.HandleError(err)
		return users, errors.Wrap(err, utils.UserExist)
	}

	return nil, nil
}

func (u *User) GetUser(user *models.User) error {
	if err := u.user.GetByNickname(user); err != nil {
		return errors.Wrap(err, utils.UserNotExist)
	}

	return nil
}

func (u *User) UpdateUser(user *models.User) error {
	uInfo := *user
	if err := userApp.userRepo.GetByNickname(&uInfo); err != nil {
		return utils.UserNotExist
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

	if err := userApp.userRepo.Update(user); err != nil {
		return errors.Wrap(err, utils.UserNotUpdated)
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
