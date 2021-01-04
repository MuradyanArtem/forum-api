package app

import "forum-api/internal/domain/repository"

type App struct {
	User   *User
	Forum  *Forum
	Thread *Thread
	Post   *Post
}

func New(r *repository.App) *App {
	return &App{
		newUser(r.User),
		newForum(r.Forum),
		newThread(r.Thread),
		newPost(r.Post),
	}
}
