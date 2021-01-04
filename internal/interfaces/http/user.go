package http

import (
	"errors"
	"fmt"
	"forum-api/internal/app"
	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"
	"io/ioutil"
	"net/http"

	json "github.com/mailru/easyjson"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type User struct {
	user *app.User
}

func newUser(user *app.User) *User {
	return &User{
		user,
	}
}

func (u *User) AddUser(w http.ResponseWriter, r *http.Request) {
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

	user := &models.User{}
	user.Nickname = mux.Vars(r)["nickname"]

	if err := json.Unmarshal(data, user); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "AddUser",
		}).Error(err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if users, err := u.user.CreateUser(user); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "AddUser",
		}).Error(err)

		if users == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(res)
			return
		}

		w.WriteHeader(http.StatusConflict)
		w.Write(res)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "AddUser",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (u *User) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	user := &models.User{}
	user.Nickname = mux.Vars(r)["nickname"]

	if err := u.user.GetUser(user); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetUser",
		}).Error(err)

		res, err := json.Marshal(
			&models.Message{
				Message: fmt.Sprintf("Can't find user with id %d", user.Nickname),
			})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "GetUser",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetUser",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	user := &models.User{}
	user.Nickname = mux.Vars(r)["nickname"]

	if err := json.Unmarshal(data, user); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetUser",
		}).Error(err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := u.user.UpdateUser(user); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "UpdateUser",
		}).Error(err)

		if errors.Is(err, infrastructure.UserNotExist) {
			res, err := json.Marshal(
				&models.Message{
					Message: fmt.Sprintf("Can't find user with id %d", user.Nickname),
				})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"pack": "http",
					"func": "UpdateUser",
				}).Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNotFound)
			w.Write(res)
			return
		}

		if errors.Is(err, infrastructure.UserNotUpdated) {
			res, err := json.Marshal(
				&models.Message{
					Message: fmt.Sprintf("This email is already registered by user: %d", user.Nickname),
				})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"pack": "http",
					"func": "UpdateUser",
				}).Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusConflict)
			w.Write(res)
			return
		}
	}

	res, err := json.Marshal(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "UpdateUser",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
