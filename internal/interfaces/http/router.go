package http

import (
	"forum-api/internal/app/usecases"

	"github.com/fasthttp/router"
)

var app *usecases.App

func New(impl *usecases.App) *router.Router {
	app = impl

	r := router.New()
	r.GET("/api/service/status", GetStatus)
	r.POST("/api/service/clear", DeleteAll)

	r.GET("/api/user/{nickname}/profile", GetUser)
	r.POST("/api/user/{nickname}/create", CreateUser)
	r.POST("/api/user/{nickname}/profile", UpdateUser)

	r.GET("/api/forum/{slug}/details", GetForumDetails)
	r.GET("/api/forum/{slug}/users", GetUsersByForum)
	r.GET("/api/forum/{slug}/threads", GetThreadsByForum)
	r.POST("/api/forum/{slug}/create", CreateThread)
	r.POST("/api/forum/create", CreateForum)

	r.GET("/api/thread/{slug}/details", GetThreadDetails)
	r.GET("/api/thread/{slug}/posts", GetPosts)
	r.POST("/api/thread/{slug}/details", UpdateThread)
	r.POST("/api/thread/{slug}/vote", UpdateVote)
	r.POST("/api/thread/{slug}/create", CreatePost)

	r.GET("/api/post/{id}/details", GetPostDetails)
	r.POST("/api/post/{id}/details", UpdatePost)
	return r
}
