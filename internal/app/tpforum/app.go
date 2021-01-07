package tpforum

import (
	"forum-api/internal/domain/repository"

	"forum-api/internal/app/usecases"
)

func New(r *repository.App) *usecases.App {
	return &usecases.App{
		User:   newUser(r.User),
		Forum:  newForum(r.Forum),
		Thread: newThread(r.Thread),
		Post:   newPost(r.Post),
	}
}
