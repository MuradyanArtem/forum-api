package http

import (
	"forum-api/internal/app"
	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Forum struct {
	forum *app.Forum
	user  *app.User
}

func newForum(forum *app.Forum, user *app.User) *Forum {
	return &Forum{
		forum,
		user,
	}
}

func (f *Forum) CreateForum(ctx *fasthttp.RequestCtx) {
	forum := &models.Forum{}
	if err := unmarshal(ctx, forum); err != nil {
		return
	}

	var err error
	forum.User, err = f.user.SelectNicknameWithCase(forum.User)
	if err != nil {
		send(ctx, http.StatusNotFound, models.Message{err.Error()})
		return
	}

	if err := f.forum.Create(forum); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreateForum",
		}).Error(err)

		switch err {
		case infrastructure.ErrNotExists:
			send(ctx, http.StatusNotFound, models.Message{err.Error()})
		case infrastructure.ErrConflict:
			forumInBase, _ := f.forum.SelectBySlug(forum.Slug)
			send(ctx, http.StatusConflict, forumInBase)
		default:
			send(ctx, http.StatusInternalServerError, models.Message{err.Error()})
		}

		return
	}
	send(ctx, http.StatusOK, forum)
}

func (f *Forum) Details(ctx *fasthttp.RequestCtx) {
	forum, err := f.forum.SelectBySlug(ctx.UserValue("slug").(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "Details",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{err.Error()})
	}
	send(ctx, http.StatusOK, forum)
}

func (f *Forum) GetUsersByForum(ctx *fasthttp.RequestCtx) {
	params := getParams(ctx)
	slug := ctx.UserValue("slug").(string)

	users, err := f.forum.GetUsersByForum(slug, params.Desc, params.Since, params.Limit)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetUsersByForum",
		}).Error(err)
		send(ctx, http.StatusNotFound, models.Message{err.Error()})
	}
	send(ctx, http.StatusOK, users)
}
