package repository

import "forum-api/internal/domain/models"

type Forum interface {
	Insert(forum *models.Forum) error
	SelectBySlug(slug string) (*models.Forum, error)
	GetUsersByForum(slug string, desc bool, since string, limit int) (models.UserSlice, error)
	SelectForumWithCase(slug string) (string, error)
}
