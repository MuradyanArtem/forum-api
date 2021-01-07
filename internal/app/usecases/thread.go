package usecases

import "forum-api/internal/domain/models"

type Thread interface {
	InsertThread(thread *models.Thread) error
	SelectThreadsByForum(slug string, limit int, since string, desc bool) (models.ThreadSlice, error)
	SelectBySlugOrID(slugOrID string) (*models.Thread, error)
	SelectByID(id int) (*models.Thread, error)
	SelectThreadBySlug(slug string) (*models.Thread, error)
	Update(thread *models.Thread) error
	Vote(vote models.Vote, slug string) (models.Thread, error)
	UpdateBySlugOrID(slug string, thread *models.Thread) error
	GetForumIDBySlug(slug string) (int, string, error)
}
