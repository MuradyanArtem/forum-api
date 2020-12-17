package http

import (
	"forum-api/internal/app"
	"net/http"

	"github.com/gorilla/mux"
)

func New(app *app.Api) http.Handler {
	router := mux.NewRouter()

	uh := NewUserHandler(app.user)
	fh := NewForumHandler(app.forum)
	th := NewThreadHandler(app.thread)
	ph := NewPostHandler(app.post, app.user, app.thread, app.forum)

	router.HandleFunc("/api/service/status", uh.GetStatus)

	router.HandleFunc("/api/service/clear", uh.DeleteAll)

	router.HandleFunc("/api/user/{nickname}/profile", uh.GetUser)

	router.HandleFunc("/api/user/{nickname}/profile", uh.UpdateUser)

	router.HandleFunc("/api/user/{nickname}/create", uh.AddUser)

	router.HandleFunc("/api/forum/create", fh.CreateForum).
		Methods("POST")

	router.HandleFunc("/api/forum/{slug}/details", fh.GetForumInfo).
		Methods("GET")

	router.HandleFunc("/api/forum/{slug}/users", fh.GetForumUsers).
		Methods("GET")

	router.HandleFunc("/api/forum/{slug}/threads", fh.GetForumThreads).
		Methods("GET").
		Queries(
			"desc",
			"limit",
			"since",
		)

	router.HandleFunc("/api/forum/{forum}/create", th.CreateThread)

	router.HandleFunc("/api/thread/{slug}/details", th.GetThreadInfo)

	router.HandleFunc("/api/thread/{slug}/details", th.UpdateThread)

	router.HandleFunc("/api/thread/{slug}/vote", th.CreateVote)

	router.HandleFunc("/api/thread/{slug}/posts", th.GetThreadPosts).
		Queries(
			"desc",
			"sort",
			"limit",
			"since",
		)

	router.HandleFunc("/api/thread/{slug}/create", ph.CreatePosts).
		Methods("POST")

	router.HandleFunc("/api/post/{id}/details", ph.GetPost).
		Methods("GET").
		Queries("related")

	router.HandleFunc("/api/post/{id}/details", ph.UpdatePost).
		Methods("POST")

	return router
}
