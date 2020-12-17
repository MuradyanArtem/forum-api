package app

import (
	"forum-api/internal/pkg/domain/models"
	"forum-api/internal/pkg/domain/repository"
	"forum-api/internal/pkg/infrastructure/utils"

	"github.com/pkg/errors"
)

type Forum struct {
	forum repository.Forum
	user  repository.User
}

func NewForum(forum repository.Forum, user repository.User) *ForumApp {
	return &ForumApp{
		forum,
		user,
	}
}

func (f *Forum) Create(forum *models.Forum) error {
	user := &models.User{}
	user.Nickname = forum.User

	if err := f.user.GetByNickname(user); err != nil {
		return errors.Wrap(err, utils.UserNotExist)
	}

	forum.User = user.Nickname

	if err = f.forum.GetBySlug(forum); err == nil {
		return errors.Wrap(err, utils.ForumExist)
	}

	if err = f.forumRepo.InsertInto(forum); err != nil {
		return err
	}

	return nil
}

func (f *ForumApp) GetForum(forum *entity.Forum) error {
	if err := f.forum.GetBySlug(forum); err != nil {
		return errors.Wrap(err, utils.ForumNotExist)
	}

	return nil
}

func (f *ForumApp) GetForumThreads(forum *models.Forum, desc, limit, since string) ([]models.Thread, error) {
	if err := f.forum.GetBySlug(forum); err != nil {
		return nil, errors.Wrap(err, utils.ForumNotExist)
	}

	threads, err := f.forum.GetThreads(forum, desc, limit, since)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (f *ForumApp) GetForumUsers(forum *models.Forum, desc, limit, since string) ([]models.User, error) {
	if err := f.forum.GetBySlug(forum); err != nil {
		return nil, errors.Wrap(err, utils.ForumNotExist)
	}

	usr, err := forumApp.forumRepo.GetUsers(f, desc, limit, since)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
