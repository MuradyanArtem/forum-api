package repository

import "forum-api/internal/domain/models"

type Thread interface {
	Update(thread *models.Thread) error
	VoteByID(vote models.Vote, id int) (models.Thread, error)
	VoteBySlug(vote models.Vote, slug string) (models.Thread, error)
	InsertThread(thread *models.Thread) error
	SelectThreadByID(id int) (*models.Thread, error)
	UpdateBySlugOrID(slug string, thread *models.Thread) error
	SelectThreadBySlug(slug string) (*models.Thread, error)
	SelectThreadsByForum(slug string, limit int, since string, desc bool) (models.ThreadSlice, error)
	GetForumIDBySlug(slug string) (int, string, error)
	SelectForumByThreadID(id int) (string, error)
}
