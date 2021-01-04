package http

import (
	"errors"
	"fmt"
	"forum-api/internal/app"
	"forum-api/internal/infrastructure"
	"io/ioutil"
	"net/http"

	"forum-api/internal/domain/models"

	json "github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type Thread struct {
	thread *app.Thread
	post   *app.Post
}

func newThread(thread *app.Thread, post *app.Post) *Thread {
	return &Thread{
		thread,
		post,
	}
}

func (t *Thread) GetThreadInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	thread := &models.Thread{}
	thread.Slug = mux.Vars(r)["slug"]

	if err := t.thread.GetThreadInfo(thread); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetThreadInfo",
		}).Error(err)

		res, err := json.Marshal(
			&models.Message{
				Message: fmt.Sprintf("Can't find user with id %d", thread.ID),
			})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "GetThreadInfo",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(res)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	res, err := json.Marshal(thread)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetThreadInfo",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (t *Thread) CreateVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateVote",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()

	vote := &models.Vote{}
	if err := json.Unmarshal(data, vote); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateVote",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	thread := &models.Thread{}
	thread.Slug = mux.Vars(r)["slug"]

	if err := t.thread.CreateVote(thread, vote); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateVote",
		}).Error(err)

		res, err := json.Marshal(
			&models.Message{
				Message: fmt.Sprintf("Can't find user with id %d", thread.ID),
			})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "CreateVote",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(res)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	res, err := json.Marshal(thread)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateVote",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (t *Thread) UpdateThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateVote",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()

	thread := &models.Thread{}
	thread.Slug = mux.Vars(r)["slug"]

	if err := json.Unmarshal(data, thread); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "UpdateThread",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := t.thread.UpdateThread(thread); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "UpdateThread",
		}).Error(err)

		res, err := json.Marshal(
			&models.Message{
				Message: fmt.Sprintf("Can't find user with id %d", thread.ID),
			})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "UpdateThread",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(res)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	res, err := json.Marshal(thread)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "UpdateThread",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (t *Thread) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	thread := &models.Thread{}
	thread.Slug = mux.Vars(r)["slug"]

	posts, err := t.thread.GetThreadPosts(
		thread,
		r.FormValue("desc"),
		r.FormValue("sort"),
		r.FormValue("limit"),
		r.FormValue("since"),
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetThreadPosts",
		}).Error(err)

		res, err := json.Marshal(
			&models.Message{
				Message: fmt.Sprintf("Can't find user with id %d", thread.ID),
			})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "GetThreadPosts",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(res)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	res, err := json.Marshal(posts)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetThreadPosts",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

func (t *Thread) CreatePosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateVote",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()

	posts := &models.Posts{}
	if err := json.Unmarshal(data, posts); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreatePosts",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	thread := &models.Thread{}
	thread.Slug = mux.Vars(r)["slug"]

	if err := t.post.CreatePosts(*posts, thread); err != nil {
		if errors.Is(err, infrastructure.ParentNotExist) {
			res, err := json.Marshal(
				&models.Message{
					Message: fmt.Sprintf("Can't find user with id %d", thread.ID),
				})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"pack": "http",
					"func": "CreatePosts",
				}).Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(res)
				return
			}

			w.WriteHeader(http.StatusConflict)
			w.Write(res)
			return
		}

		if errors.Is(err, infrastructure.ThreadNotExist) || errors.Is(err, infrastructure.UserNotExist) {
			res, err := json.Marshal(
				&models.Message{
					Message: fmt.Sprintf("Can't find user with id %d", thread.ID),
				})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"pack": "http",
					"func": "CreatePosts",
				}).Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(res)
				return
			}
		}

		res, err := json.Marshal(&models.Message{Message: err.Error()})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "CreatePosts",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(res)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	res, err := json.Marshal(posts)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreatePosts",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
