package http

import (
	"forum-api/internal/app"
	"net/http"
)

func New(app *app.App) http.Handler {
	uh := newUser(app.User)
	fh := newForum(app.Forum, app.Thread, app.User)
	th := newThread(app.Thread, app.Post)
	ph := newPost(app.Post, app.User, app.Thread, app.Forum)
	sh := newService(app.User)

	r := router.New()
	r.GET("/api/service/status", sh.GetStatus)
	r.POST("/api/service/clear", sh.DeleteAll)

	r.GET("/api/user/{nickname}/profile", m.GetProfile)
	r.POST("/api/user/{nickname}/create", m.CreateUser)
	r.POST("/api/user/{nickname}/profile", m.UpdateProfile)

	r.GET("/api/forum/{slug}/details", fh.GetDetails)
	r.GET("/api/forum/{slug}/users", fh.GetUsersByForum)
	r.GET("/api/forum/{slug}/threads", m.GetThreadsByForum)
	r.POST("/api/forum/{slug}/create", m.CreateThread)
	r.POST("/api/forum/create", fh.CreateForum)

	r.GET("/api/thread/{slugOrID}/details", m.Details)
	r.GET("/api/thread/{slugOrID}/posts", m.GetPosts)
	r.POST("/api/thread/{slugOrID}/details", m.Update)
	r.POST("/api/thread/{slugOrID}/vote", m.Vote)
	r.POST("/api/thread/{slugOrID}/create", m.Create)

	r.GET("/api/post/{id}/details", m.GetByID)
	r.POST("/api/post/{id}/details", m.Update)
	return r
}
