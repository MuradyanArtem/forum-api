package repository

import "forum-api/internal/domain/models"

type Forum interface {
	InsertInto(forum *models.Forum) error
	GetBySlug(forum *models.Forum) error
	GetThreads(forum *models.Forum, desc, limit, since string) ([]models.Thread, error)
	GetUsers(forum *models.Forum, desc, limit, since string) ([]models.User, error)
	Prepare() error
}
