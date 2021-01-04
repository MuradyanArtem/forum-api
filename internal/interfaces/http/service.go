package http

import (
	"forum-api/internal/app"
	"forum-api/internal/domain/models"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Service struct {
	user *app.User
}

func newService(user *app.User) *Service {
	return &Service{
		user,
	}
}

func (s *Service) DeleteAll(ctx *fasthttp.RequestCtx) {
	if err := s.user.DeleteAll(); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "DeleteAll",
		}).Error(err)
		setStatus(ctx, http.StatusInternalServerError)
		return
	}
	setStatus(ctx, http.StatusOK)
}

func (s *Service) GetStatus(ctx *fasthttp.RequestCtx) {
	status := &models.Status{}
	if err := s.user.GetStatus(status); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetStatus",
		}).Error(err)
		setStatus(ctx, http.StatusInternalServerError)
		return
	}

	res, err := marshall(ctx, status)
	if err != nil {
		setStatus(ctx, http.StatusInternalServerError)
		return
	}
	setStatus(ctx, http.StatusOK)
}
