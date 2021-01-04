package app

import (
	"forum-api/internal/domain/models"
	"forum-api/internal/domain/repository"
)

type Forum struct {
	forum repository.Forum
}

func newForum(forum repository.Forum) *Forum {
	return &Forum{
		forum,
	}
}

func (f *Forum) Create(forum *models.Forum) error {
	return f.forum.Insert(forum)
}

func (f *Forum) SelectBySlug(slug string) (*models.Forum, error) {
	return f.forum.SelectBySlug(slug)
}

func (f *Forum) GetUsersByForum(slug string, desc bool, since string, limit int) (models.UserSlice, error) {
	return f.forum.GetUsersByForum(slug, desc, since, limit)
}

func (f *Forum) SelectForumWithCase(slug string) (string, error) {
	return f.forum.SelectForumWithCase(slug)
}
