package http

import (
	"fmt"
	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"
	"net/http"

	"github.com/valyala/fasthttp"
)

func CreateForum(ctx *fasthttp.RequestCtx) {
	forum := &models.Forum{}
	if err := unmarshall(ctx, forum); err != nil {
		return
	}

	var err error
	forum.User, err = app.User.SelectNicknameWithCase(forum.User)
	if err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}

	if err := app.Forum.Create(forum); err != nil {
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
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, forum)
}

func GetUsersByForum(ctx *fasthttp.RequestCtx) {
	params := getParams(ctx)
	slug := ctx.UserValue("slug").(string)
	if _, err := app.Forum.SelectForumWithCase(slug); err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: fmt.Sprintf("Can't find forum by slug: %v", slug)})
		return
	}

	users, err := app.Forum.GetUsersByForum(slug, params.Desc, params.Since, params.Limit)
	if err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, users)
}
