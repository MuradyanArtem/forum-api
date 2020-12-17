package app

import (
	"forum-api/internal/pkg/domain/models"
	"forum-api/internal/pkg/domain/repository"
	"forum-api/internal/pkg/infrastructure/utils"

	"github.com/pkg/errors"
)

type Post struct {
	post   repository.Post
	thread repository.Thread
}

func NewPost(post repository.Post, thread repository.Thread) *Post {
	return &Post{
		post,
		thread,
	}
}

func (p *Post) CreatePosts(post []*entity.Post, thread *entity.Thread) error {
	if err := p.threadRepo.GetBySlugOrId(thread); err != nil {
		return errors.Wrap(err, utils.ThreadNotExist)
	}

	if err = postApp.postRepo.InsertInto(p, th); err != nil {
		if err.Error() == "ERROR: Parent post was created in another thread (SQLSTATE 00404)" {
			return ParentNotExist
		} else {
			return UserNotExist
		}
	}

	for _, el := range post {
		if err = postApp.threadRepo.InsertIntoForumUsers(el.Forum, el.Author); err != nil {
			return err
		}
	}

	return nil
}

func (p *Post) GetPost(post *models.Post) error {
	if err := p.postRepo.GetById(post); err != nil {
		return errors.Wrap(err, utils.PostNotExist)
	}

	return nil
}

func (p *Post) UpdatePost(post *models.Post) error {
	message := post.Message

	if err := p.post.GetById(post); err != nil {
		return errors.Wrap(err, utils.PostNotExist)
	}

	if message != "" && message != post.Message {
		post.Message = message
		if err := p.post.Update(post); err != nil {
			return errors.Wrap(err, utils.PostNotExist)
		}
	}

	return nil
}
