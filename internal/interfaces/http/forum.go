package http

import (
	"forum-api/internal/app"
	"forum-api/internal/domain/models"
	"net/http"

	json "github.com/mailru/easyjson"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Forum struct {
	forum  *app.Forum
	thread *app.Thread
}

func newForum(forum *app.Forum, thread *app.Thread) *Forum {
	return &Forum{
		forum,
		thread,
	}
}

func (f *Forum) CreateForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	forum := &models.Forum{}

	if err := json.Unmarshal(r.Body(), forum); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateForum",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := f.forum.CreateForum(forum); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateForum",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(forum)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateForum",
		}).Error(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (f *Forum) GetForumInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	forum := &models.Forum{}
	forum.Slug = mux.Vars(r)["slug"]

	if err := f.forum.GetForum(forum); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetForumUsers",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(r)
		return
	}

	res, err := json.Marshal(&forum)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetForumUsers",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (f *Forum) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	forum := &entity.Forum{}
	forum.Slug = mux.Vars(r)["slug"]

	threads, err := f.forum.GetForumThreads(forum, r.FormValue("desc"), r.FormValue("limit"), r.FormValue("since"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetForumUsers",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(threads)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetForumUsers",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (f *Forum) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	forum := &models.Forum{}
	forum.Slug = mux.Vars(r)["slug"]

	users, err := f.forum.GetForumUsers(forum, r.FormValue("desc"), r.FormValue("limit"), r.FormValue("since"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetForumUsers",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(users)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetForumUsers",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (f *Forum) CreateThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	thread, err := models.GetThreadFromBody(r.Body)
	tools.HandleError(err)
	vars := mux.Vars(r)

	t.Forum = vars["forum"]
	if err := t.threadApp.CreateThread(t); err != nil {
		if err == tools.ThreadExist {
			w.WriteHeader(http.StatusConflict)
			res, err := json.Marshal(&t)
			w.Write(res)
			return
		}

		if err == tools.UserNotExist {
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&models.Message{Message: "user not exist"})
			w.Write(res)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	res, err := json.Marshal(&th)
	tools.HandleError(err)
	w.Write(res)
}
