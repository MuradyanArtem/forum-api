package usecases

import "forum-api/internal/domain/models"

type Post interface {
	InsertPost(posts []models.Post, forum string, id int) error
	SelectPostByID(id int) (*models.Post, error)
	Update(post *models.Post) error
	GetPosts(threadID int, desc bool, since string, limit int, sort string) (models.PostSlice, error)
}
