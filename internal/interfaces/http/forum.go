package http

import (
	"fmt"
	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CreateForum(ctx *fasthttp.RequestCtx) {
	forum := &models.Forum{}
	if err := unmarshall(ctx, forum); err != nil {
		return
	}

	if err := app.Forum.Create(forum); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateForum",
		}).Error(err)

		switch err {
		case infrastructure.ErrNotExists:
			send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		case infrastructure.ErrConflict:
			forumInBase, _ := app.Forum.SelectBySlug(forum.Slug)
			send(ctx, http.StatusConflict, forumInBase)
		default:
			send(ctx, http.StatusInternalServerError, models.Message{Message: err.Error()})
		}
		return
	}
	send(ctx, http.StatusCreated, forum)
}

func GetForumDetails(ctx *fasthttp.RequestCtx) {
	forum, err := app.Forum.SelectBySlug(ctx.UserValue("slug").(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetForumDetails",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, forum)
}

func GetUsersByForum(ctx *fasthttp.RequestCtx) {
	params := getParams(ctx)
	slug := ctx.UserValue("slug").(string)
	if _, err := app.Forum.SelectForumWithCase(slug); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetUsersByForum",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{Message: fmt.Sprintf("Can't find forum by slug: %v", slug)})
		return
	}

	users, err := app.Forum.GetUsersByForum(slug, params.Desc, params.Since, params.Limit)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetUsersByForum",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, users)
}
