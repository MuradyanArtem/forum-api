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
		newForum(r.Forum, r.User),
		newThread(r.Thread, r.Forum),
		newPost(r.Post, r.Thread),
	}
}
