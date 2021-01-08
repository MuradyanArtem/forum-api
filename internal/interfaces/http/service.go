package http

import (
	"forum-api/internal/domain/models"
	"net/http"

	"github.com/valyala/fasthttp"
)

func DeleteAll(ctx *fasthttp.RequestCtx) {
	if err := app.User.DeleteAll(); err != nil {
		setStatus(ctx, http.StatusInternalServerError)
		return
	}
	setStatus(ctx, http.StatusOK)
}

func GetStatus(ctx *fasthttp.RequestCtx) {
	status := &models.Status{}
	if err := app.User.GetStatus(status); err != nil {
		setStatus(ctx, http.StatusInternalServerError)
		return
	}
	send(ctx, http.StatusOK, status)
}
