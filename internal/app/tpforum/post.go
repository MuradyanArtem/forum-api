package tpforum

import (
	"forum-api/internal/domain/models"
	"forum-api/internal/domain/repository"
)

type Post struct {
	post repository.Post
}

func newPost(post repository.Post) *Post {
	return &Post{
		post,
	}
}

func (p *Post) InsertPost(posts []models.Post, forum string, id int) error {
	return p.post.InsertPost(posts, forum, id)
}

func (p *Post) Update(post *models.Post) error {
	return p.post.Update(post)
}

func (p *Post) SelectPostByID(id int) (*models.Post, error) {
	return p.post.SelectPostByID(id)
}

func (p *Post) GetPosts(threadID int, desc bool, since string, limit int, sort string) (models.PostSlice, error) {
	return p.post.GetPosts(threadID, desc, since, limit, sort)
}
