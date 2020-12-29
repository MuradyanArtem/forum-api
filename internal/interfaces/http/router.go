package http

import (
	"forum-api/internal/app"
	"net/http"

	"github.com/gorilla/mux"
)

func New(app *app.App) http.Handler {
	router := mux.NewRouter()

	uh := newUser(app.user)
	fh := newForum(app.forum, app.thread)
	th := newThread(app.thread, app.post)
	ph := newPost(app.post, app.user, app.thread, app.forum)
	sh := newService(app.user)

	router.HandleFunc("/api/service/status", sh.GetStatus).
		Methods("GET")
	router.HandleFunc("/api/service/clear", sh.DeleteAll).
		Methods("POST")

	router.HandleFunc("/api/user/{nickname}/profile", uh.GetUser).
		Methods("GET")
	router.HandleFunc("/api/user/{nickname}/profile", uh.UpdateUser).
		Methods("POST")
	router.HandleFunc("/api/user/{nickname}/create", uh.AddUser).
		Methods("POST")

	router.HandleFunc("/api/forum/{slug}/create", fh.CreateThread).
		Methods()
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

	router.HandleFunc("/api/thread/{slug}/details", th.GetThreadInfo).
	Methods("GET")
	router.HandleFunc("/api/thread/{slug}/details", th.UpdateThread).
	Methods("POST")
	router.HandleFunc("/api/thread/{slug}/vote", th.CreateVote).
	Methods("POST")
	router.HandleFunc("/api/thread/{slug}/create", th.CreatePosts).
		Methods("POST")
	router.HandleFunc("/api/thread/{slug}/posts", th.GetThreadPosts).
		Queries(
			"desc",
			"sort",
			"limit",
			"since",
		)

	router.HandleFunc("/api/post/{id}/details", ph.UpdatePost).
		Methods("POST")
	router.HandleFunc("/api/post/{id}/details", ph.GetPost).
		Methods("GET").
		Queries("related")

	return router
}
