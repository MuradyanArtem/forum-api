package repository

import "forum-api/internal/domain/models"

type Post interface {
	InsertInto(posts []*models.Post, thread *models.Thread) error
	GetById(post *models.Post) error
	Update(post *models.Post) error
	Prepare() error
}
