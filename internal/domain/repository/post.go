package repository

import "forum-api/internal/domain/models"

type Post interface {
	SelectThreadByPostID(id int) (int, error)
	InsertPost(posts *models.PostSlice, forum string, id int) error
	SelectPostByID(id int) (*models.Post, error)
	Update(post *models.Post) error
	GetPosts(threadID int, desc bool, since string, limit int, sort string) (models.PostSlice, error)
}
