package http

import (
	"forum-api/internal/app"
	"forum-api/internal/domain/models"
	"net/http"

	json "github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

type Service struct {
	user *app.User
}

func newService(user *app.User) *Service {
	return &User{
		user,
	}
}

func (s *Service) DeleteAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	if err := s.user.DeleteAll(); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "DeleteAll",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	res, err := json.Marshal(&tools.Message{Message: "all info deleted"})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "DeleteAll",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (s *Service) GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	status := &models.Status{}

	if err := s.user.GetStatus(status); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetStatus",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	res, err := json.Marshal(&status)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetStatus",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
