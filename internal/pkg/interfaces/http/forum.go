package http

import (
	"encoding/json"
	"forum-api/internal/app"
	"forum-api/internal/pkg/infrastructure/utils"
	"net/http"
	"tech-db-project/application"
	"tech-db-project/domain/entity"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ForumHandler struct {
	forumApp *app.Forum
}

func NewForumHandler(forum *app.Forum) return *ForumHandler{
	return &ForumHandler{
		forum,
	}
}

func (fh *ForumHandler) CreateForum(w http.ResponseWriter, r *http.Request) {
	forum, err := models.GetForumFromBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateForum",
		}).Error(err)
		return
	}

	if err := fh.forumApp.CreateForum(forum); err != nil {
		switch err {
		case utils.UserNotExist:
			w.WriteHeader(http.StatusNotFound)
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "CreateForum",
			}).Error(err)

			if res, err := json.Marshal(&tools.Message{Message: "User not found"}); err == nil {
				w.Write(res)
			}
			return

		case utils.ForumExist:
			w.WriteHeader(http.StatusNotFound)
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "CreateForum",
			}).Error(err)
			return

		default:
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "CreateForum",
			}).Error(err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	logrus.WithFields(logrus.Fields{
		"pack": "http",
		"func": "CreateForum",
	}).Error(err)

	if res, err := json.Marshal(&f); err == nil {
		w.Write(res)
	}
}

func (fh *ForumHandler) GetForumInfo(w http.ResponseWriter, r *http.Request) {
	f := &entity.Forum{}

	vars := mux.Vars(r)
	f.Slug = vars["slug"]

	if err := fh.forumApp.GetForum(f); err != nil {
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "User not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&f)
	w.Write(res)
	tools.HandleError(err)
}

func (fh *ForumHandler) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	f := &entity.Forum{}

	vars := mux.Vars(r)
	f.Slug = vars["slug"]

	ths, err := fh.forumApp.GetForumThreads(f, r.FormValue("desc"), r.FormValue("limit"), r.FormValue("since"))
	if err != nil {
		switch err {
		case tools.ForumNotExist:
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "forum not found"})
			tools.HandleError(err)
			w.Write(res)
			return
		default:
			tools.HandleError(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&ths)
	w.Write(res)
	tools.HandleError(err)
}

func (fh *ForumHandler) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	f := &entity.Forum{}

	vars := mux.Vars(r)
	f.Slug = vars["slug"]

	users, err := fh.forumApp.GetForumUsers(f, r.FormValue("desc"), r.FormValue("limit"), r.FormValue("since"))
	if err != nil {
		switch err {
		case tools.ForumNotExist:
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "forum not found"})
			tools.HandleError(err)
			w.Write(res)
			tools.HandleError(err)
			return
		default:
			tools.HandleError(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&users)
	w.Write(res)
	tools.HandleError(err)
}
