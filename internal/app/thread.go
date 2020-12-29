package app

import (
	"forum-api/internal/domain/models"
	"forum-api/internal/domain/repository"
	"forum-api/internal/infrastructure"

	"github.com/pkg/errors"
)

type Thread struct {
	thread repository.Thread
	forum  repository.Forum
}

func newThread(thread repository.Thread, forum repository.Forum) *Thread {
	return &Thread{
		thread,
		forum,
	}
}

func (t *Thread) CreateThread(thread *models.Thread) error {
	forum := &models.Forum{}
	forum.Slug = thread.Forum

	if err := t.forum.GetBySlug(forum); err != nil {
		return errors.Wrap(err, infrastructure.UserNotExist.Error())
	}

	thread.Forum = forum.Slug

	if err := t.thread.InsertInto(thread); err != nil {
		if err := t.thread.GetBySlugOrId(thread); err != nil {
			return errors.Wrap(err, infrastructure.UserNotExist.Error())
		}

		return errors.Wrap(err, infrastructure.ThreadExist.Error())
	}

	return nil
}

func (t *Thread) GetThreadInfo(thread *models.Thread) error {
	if err := t.thread.GetBySlugOrId(thread); err != nil {
		return errors.Wrap(err, infrastructure.ThreadNotExist.Error())
	}

	return nil
}

func (t *Thread) CreateVote(thread *models.Thread, vote *models.Vote) error {
	if err := t.thread.GetBySlugOrId(thread); err != nil {
		return errors.Wrap(err, infrastructure.ThreadNotExist.Error())
	}

	vote.Thread = thread.ID

	if err := t.thread.InsertIntoVotes(thread, vote); err != nil {
		return errors.Wrap(err, infrastructure.UserNotExist.Error())
	}

	return nil
}

func (t *Thread) UpdateThread(thread *models.Thread) error {
	if err := t.thread.Update(thread); err != nil {
		return errors.Wrap(err, infrastructure.ThreadNotExist.Error())
	}

	return nil
}

func (t *Thread) GetThreadPosts(thread *models.Thread, desc, sort, limit, since string) ([]models.Post, error) {
	if err := t.thread.GetBySlugOrId(thread); err != nil {
		return nil, errors.Wrap(err, infrastructure.ThreadNotExist.Error())
	}

	posts, err := t.thread.GetPosts(thread, desc, sort, limit, since)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
