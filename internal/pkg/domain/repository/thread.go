package repository

import "forum-api/internal/pkg/domain/models"

type Thread interface {
	GetBySlug(thread *models.Thread) error
	GetById(thread *models.Thread) error
	GetBySlugOrId(thread *models.Thread) error
	GetPosts(thread *models.Thread, desc, sort, limit, since string) ([]models.Post, error)

	InsertInto(thread *models.Thread) error
	InsertIntoVotes(thread *models.Thread, vote *models.Vote) error
	InsertIntoForumUsers(forum, nickname string) error

	Update(thread *models.Thread) error

	Prepare() error
}
