package app

import (
	"forum-api/internal/domain/models"
	"forum-api/internal/domain/repository"
	"forum-api/internal/infrastructure"

	"github.com/pkg/errors"
)

type Post struct {
	post   repository.Post
	thread repository.Thread
}

func newPost(post repository.Post, thread repository.Thread) *Post {
	return &Post{
		post,
		thread,
	}
}

func (p *Post) CreatePosts(post []*models.Post, thread *models.Thread) error {
	if err := p.thread.GetBySlugOrId(thread); err != nil {
		return errors.Wrap(infrastructure.ThreadNotExist, err.Error())
	}

	if err := p.post.InsertInto(post, thread); err != nil {
		return errors.Wrap(infrastructure.ParentNotExist, err.Error())
	}

	for _, el := range post {
		if err := p.thread.InsertIntoForumUsers(el.Forum, el.Author); err != nil {
			return err
		}
	}

	return nil
}

func (p *Post) GetPost(post *models.Post) error {
	if err := p.post.GetById(post); err != nil {
		return errors.Wrap(err, infrastructure.PostNotExist.Error())
	}

	return nil
}

func (p *Post) UpdatePost(post *models.Post) error {
	message := post.Message

	if err := p.post.GetById(post); err != nil {
		return errors.Wrap(err, infrastructure.PostNotExist.Error())
	}

	if message != "" && message != post.Message {
		post.Message = message
		if err := p.post.Update(post); err != nil {
			return errors.Wrap(err, infrastructure.PostNotExist.Error())
		}
	}

	return nil
}
