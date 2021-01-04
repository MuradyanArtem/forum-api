package app

import (
	"forum-api/internal/domain/models"
	"forum-api/internal/domain/repository"
	"strconv"
)

type Thread struct {
	thread repository.Thread
}

func newThread(thread repository.Thread) *Thread {
	return &Thread{
		thread,
	}
}

func (t *Thread) InsertThread(thread *models.Thread) error {
	return t.thread.InsertThread(thread)
}

func (t *Thread) SelectThreadsByForum(slug string, limit int, since string, desc bool) (models.ThreadSlice, error) {
	return t.thread.SelectThreadsByForum(slug, limit, since, desc)
}

func (t *Thread) SelectBySlugOrID(slugOrID string) (*models.Thread, error) {
	value, err := strconv.Atoi(slugOrID)
	if err != nil {
		return t.thread.SelectThreadBySlug(slugOrID)
	}
	return t.thread.SelectThreadByID(value)
}
func (t *Thread) SelectByID(id int) (*models.Thread, error) {
	return t.thread.SelectThreadByID(id)
}

func (t *Thread) SelectThreadBySlug(slug string) (*models.Thread, error) {
	return t.thread.SelectThreadBySlug(slug)
}
func (t *Thread) Update(thread *models.Thread) error {
	return t.thread.Update(thread)
}

func (t *Thread) Vote(vote models.Vote, slug string) (models.Thread, error) {
	id, err := strconv.Atoi(slug)
	if err != nil {
		return t.thread.VoteBySlug(vote, slug)
	}
	return t.thread.VoteByID(vote, id)
}

func (t *Thread) UpdateBySlugOrID(slug string, thread *models.Thread) error {
	return t.thread.UpdateBySlugOrID(slug, thread)
}

func (t *Thread) GetForumIDBySlug(slug string) (int, string, error) {
	value, err := strconv.Atoi(slug)
	if err != nil {
		return t.thread.GetForumIDBySlug(slug)
	}
	forum, err := t.thread.SelectForumByThreadID(value)
	return value, forum, err
}
