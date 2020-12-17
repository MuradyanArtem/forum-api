package http

import (
	"encoding/json"
	"forum-api/internal/app"
	"forum-api/internal/pkg/domain/models"
	"net/http"
	"strconv"
	"strings"
	"tech-db-project/infrasctructure/tools"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	postApp   *app.Post
	userApp   *app.User
	threadApp *app.Thread
	forumApp  *app.Forum
}

func NewPostHandler(post *app.Post, user *app.User, thread *app.Thread, forum *app.Forum) *PostHandler {
	return &PostHandler{
		post,
		user,
		thread,
		forum,
	}
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func (ph *PostHandler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	p := make([]*models.Post, 0)

	th := models.Thread{}
	vars := mux.Vars(r)
	th.Slug = vars["slug"]

	err := ph.postApp.CreatePosts(p, &th)
	if err != nil {
		if err == tools.ParentNotExist {
			w.WriteHeader(http.StatusConflict)
			res, err := json.Marshal(&tools.Message{Message: "parent conflict"})
			tools.HandleError(err)
			w.Write(res)
			return
		}
		if err == tools.ThreadNotExist || err == tools.UserNotExist {
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "user or thread not found"})
			tools.HandleError(err)
			w.Write(res)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		res, err := json.Marshal(&tools.Message{Message: err.Error()})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res, err := json.Marshal(&p)
	tools.HandleError(err)
	w.Write(res)
}

func (ph *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	str := r.FormValue("related")
	related := strings.Split(str, ",")

	p := &models.Post{}
	th := &models.Thread{}
	f := &models.Forum{}
	u := &models.User{}

	var err error
	vars := mux.Vars(r)

	p.Id, err = strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logrus.Error("Cannot parse id")
	}
	if err := ph.postApp.GetPost(p); err != nil {
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: err.Error()})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	if Find(related, "user") {
		u.Nickname = p.Author
		err := ph.userApp.GetUser(u)
		tools.HandleError(err)
	} else {
		u = nil
	}

	if Find(related, "thread") {
		th.Slug = strconv.FormatInt(p.Thread, 10)
		err := ph.threadApp.GetThreadInfo(th)
		tools.HandleError(err)
	} else {
		th = nil
	}

	if Find(related, "forum") {
		f.Slug = p.Forum
		err := ph.forumApp.GetForum(f)
		tools.HandleError(err)
	} else {
		f = nil
	}

	ans := models.Info{
		Post:   p,
		Forum:  f,
		Thread: th,
		Author: u,
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&ans)
	tools.HandleError(err)
	w.Write(res)
}

func (ph *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	p := &models.Post{}
	vars := mux.Vars(r)
	var err error

	p.Id, err = strconv.ParseInt(vars["id"], 10, 64)
	tools.HandleError(err)
	if err := ph.postApp.UpdatePost(p); err != nil {
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "post not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&p)
	tools.HandleError(err)
	w.Write(res)
}
