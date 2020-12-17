package app

import (
	"forum-api/internal/pkg/domain/models"
	"forum-api/internal/pkg/domain/repository"
	"forum-api/internal/pkg/infrastructure/utils"

	"github.com/pkg/errors"
)

type Thread struct {
	thread repository.Thread
	forum  repository.Forum
}

func NewThreadApp(thread repository.Thread, forum repository.Forum) *Thread {
	return &Thread{
		thread,
		forum,
	}
}

func (t *Thread) CreateThread(thread *models.Thread) error {
	forum := &models.Forum{}
	forum.Slug = thread.Forum

	if err := threadApp.forumRepo.GetBySlug(forum); err != nil {
		return errors.Wrap(err, utils.UserNotExist)
	}

	thread.Forum = forum.Slug

	if err = t.thread.InsertInto(thread); err != nil {
		if err = t.thread.GetBySlugOrId(thread); err != nil {
			return errors.Wrap(err, utils.UserNotExist)
		}

		return errors.Wrap(err, tools.ThreadExist)
	}

	return nil
}

func (t *Thread) GetThreadInfo(thread *models.Thread) error {
	if err := t.thread.GetBySlugOrId(thread)l err != nil {
		return errors.Wrap(err, utils.ThreadNotExist)
	}

	return nil
}

func (t *Thread) CreateVote(thread *models.Thread, vote *models.Vote) error {
	if err := t.thread.GetBySlugOrId(thread); err != nil {
		return tools.ThreadNotExist
	}

	vote.Thread = thread.ID

	if err = t.thread.InsertIntoVotes(thread, vote); err != nil {
		return errors.Wrap(err, utils.UserNotExist)
	}

	return nil
}

func (t *Thread) UpdateThread(thread *models.Thread) error {
	if err := t.thread.Update(thread); err != nil {
		return errors.Wrap(err, utils.ThreadNotExist)
	}

	return nil
}

func (t *Thread) GetThreadPosts(thread *models.Thread, desc, sort, limit, since string) ([]models.Post, error) {
	if err := thread.thread.GetBySlugOrId(thread); err != nil {
		return nil, errors.Wrap(err, utils.ThreadNotExist)
	}

	posts, err := t.thread.GetPosts(thread, desc, sort, limit, since)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
