package http

import (
	"net/http"

	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CreateThread(ctx *fasthttp.RequestCtx) {
	thread := &models.Thread{}
	if err := unmarshall(ctx, thread); err != nil {
		return
	}

	var err error
	thread.Forum = ctx.UserValue("slug").(string)
	thread.Forum, err = app.Forum.SelectForumWithCase(thread.Forum)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateThread",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}

	thread.Author, err = app.User.SelectNicknameWithCase(thread.Author)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateThread",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}

	threadInBase := models.Thread{}
	if thread.Slug != "" {
		threadInBase, err := app.Thread.SelectThreadBySlug(thread.Slug)
		if err == nil {
			logrus.WithFields(logrus.Fields{
				"pack": "http",
				"func": "CreateThread",
			}).Error(err)
			send(ctx, http.StatusConflict, threadInBase)
			return
		}
	}

	if err := app.Thread.InsertThread(thread); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateThread",
		}).Error(err)

		switch err {
		case infrastructure.ErrConflict:
			send(ctx, http.StatusConflict, threadInBase)
		default:
			send(ctx, http.StatusInternalServerError, models.Message{Message: err.Error()})
		}
	}
	send(ctx, http.StatusCreated, thread)
}

func GetThreadsByForum(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	params := getParams(ctx)

	if _, err := app.Forum.SelectForumWithCase(slug); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetThreadsByForum",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}

	threads, err := app.Thread.SelectThreadsByForum(slug, params.Limit, params.Since, params.Desc)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetThreadsByForum",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, threads)
}

func GetThreadDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	thread, err := app.Thread.SelectBySlugOrID(slug)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetThreadDetails",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, thread)
}

func UpdateThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	thread := &models.Thread{}
	if err := unmarshall(ctx, thread); err != nil {
		return
	}

	if err := app.Thread.UpdateBySlugOrID(slug, thread); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "UpdateThread",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, thread)
}

func UpdateVote(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	vote := &models.Vote{}
	if err := unmarshall(ctx, vote); err != nil {
		return
	}

	thread, err := app.Thread.Vote(*vote, slug)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "UpdateVote",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, thread)
}
