package http

import (
	"fmt"
	"forum-api/internal/app"
	"forum-api/internal/domain/models"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	json "github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

type Post struct {
	postApp   *app.Post
	userApp   *app.User
	threadApp *app.Thread
	forumApp  *app.Forum
}

func newPost(post *app.Post, user *app.User, thread *app.Thread, forum *app.Forum) *Post {
	return &PostHandler{
		post,
		user,
		thread,
		forum,
	}
}

func (p *Post) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	related := strings.Split(r.FormValue("related"), ",")

	var err error
	post := &models.Post{}
	post.ID, err = strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetPost",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := p.post.GetPost(post); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetPost",
		}).Error(err)
		res, err := json.Marshal(&tools.Message{Message: err.Error()})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "GetPost",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	user := nil
	if sort.SearchStrings(related, "user") {
		user = &models.User{}
		user.Nickname = post.Author

		if err := p.user.GetUser(user); err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "GetPost",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	thread := nil
	if sort.SearchStrings(related, "thread") {
		thread := &models.Thread{}
		thread.Slug = strconv.FormatInt(post.Thread, 10)

		if err := p.thread.GetThreadInfo(thread); err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "GetPost",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	forum := nil
	if sort.SearchStrings(related, "forum") {
		forum := &models.Forum{}
		forum.Slug = post.Forum

		if err := p.forum.GetForum(forum); err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "GetPost",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	res, err := json.Marshal(&models.PostFull{
		Post:   post,
		Forum:  forum,
		Thread: thread,
		Author: user,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetPost",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (p *Post) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var err error
	post := &models.Post{}
	post.ID, err = strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	if err := p.post.UpdatePost(post); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "UpdatePost",
		}).Error(err)

		res, err := json.Marshal(
			&models.Message{
				Message: fmt.Sprintf("Can't find user with id %d", post.ID),
			})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "UpdatePost",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	res, err := json.Marshal(&post)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "UpdatePost",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
