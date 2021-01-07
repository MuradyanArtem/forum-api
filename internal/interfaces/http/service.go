package http

import (
	"forum-api/internal/domain/models"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func DeleteAll(ctx *fasthttp.RequestCtx) {
	if err := app.User.DeleteAll(); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "DeleteAll",
		}).Error(err)
		setStatus(ctx, http.StatusInternalServerError)
		return
	}
	setStatus(ctx, http.StatusOK)
}

func GetStatus(ctx *fasthttp.RequestCtx) {
	status := &models.Status{}
	if err := app.User.GetStatus(status); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetStatus",
		}).Error(err)
		setStatus(ctx, http.StatusInternalServerError)
		return
	}
	send(ctx, http.StatusOK, status)
}
