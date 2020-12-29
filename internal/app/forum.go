package app

import (
	"forum-api/internal/domain/models"
	"forum-api/internal/domain/repository"
	"forum-api/internal/infrastructure"

	"github.com/pkg/errors"
)

type Forum struct {
	forum repository.Forum
	user  repository.User
}

func newForum(forum repository.Forum, user repository.User) *Forum {
	return &Forum{
		forum,
		user,
	}
}

func (f *Forum) CreateForum(forum *models.Forum) error {
	user := &models.User{}
	user.Nickname = forum.User

	if err := f.user.GetByNickname(user); err != nil {
		return errors.Wrap(err, infrastructure.UserNotExist.Error())
	}

	forum.User = user.Nickname

	if err := f.forum.GetBySlug(forum); err == nil {
		return errors.Wrap(err, infrastructure.ForumExist.Error())
	}

	if err := f.forum.InsertInto(forum); err != nil {
		return err
	}

	return nil
}

func (f *Forum) GetForum(forum *models.Forum) error {
	if err := f.forum.GetBySlug(forum); err != nil {
		return err
	}

	return nil
}

func (f *Forum) GetForumThreads(forum *models.Forum, desc, limit, since string) ([]models.Thread, error) {
	if err := f.forum.GetBySlug(forum); err != nil {
		return nil, errors.Wrap(err, infrastructure.ForumNotExist.Error())
	}

	threads, err := f.forum.GetThreads(forum, desc, limit, since)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (f *Forum) GetForumUsers(forum *models.Forum, desc, limit, since string) ([]models.User, error) {
	if err := f.forum.GetBySlug(forum); err != nil {
		return nil, errors.Wrap(err, infrastructure.ForumNotExist.Error())
	}

	usr, err := f.forum.GetUsers(forum, desc, limit, since)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
