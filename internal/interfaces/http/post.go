package http

import (
	"bytes"
	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CreatePost(ctx *fasthttp.RequestCtx) {
	id, forum, err := app.Thread.GetForumIDBySlug(ctx.UserValue("slug").(string))
	if err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}

	posts := &models.PostSlice{}
	if err := unmarshall(ctx, posts); err != nil {
		return
	}

	if err := app.Post.InsertPost(*posts, forum, id); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "CreatePost",
		}).Error(err)

		switch err {
		case infrastructure.ErrConflict:
			send(ctx, http.StatusConflict, models.Message{Message: err.Error()})
		default:
			send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		}
		return
	}
	send(ctx, http.StatusCreated, posts)
}

func UpdatePost(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "UpdatePost",
		}).Error(err)
		setStatus(ctx, http.StatusBadRequest)
		return
	}

	post := &models.Post{}
	if err := unmarshall(ctx, post); err != nil {
		return
	}

	post.ID = id
	if err := app.Post.Update(post); err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, post)
}

func GetPostDetails(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "GetPost",
		}).Error(err)
		setStatus(ctx, http.StatusBadRequest)
		return
	}

	details := &models.PostDetails{}
	details.Post, err = app.Post.SelectPostByID(id)
	if err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}

	user, err := app.User.SelectByNickname(details.Post.Author)
	if err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	details.Author = &user

	details.Thread, err = app.Thread.SelectByID(details.Post.Thread)
	if err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}

	details.Forum, err = app.Forum.SelectBySlug(details.Thread.Forum)
	if err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}

	related := ctx.QueryArgs().Peek("related")
	if !bytes.Contains(related, []byte("user")) {
		details.Author = nil
	}
	if !bytes.Contains(related, []byte("forum")) {
		details.Forum = nil
	}
	if !bytes.Contains(related, []byte("thread")) {
		details.Thread = nil
	}
	send(ctx, http.StatusOK, details)
}

func GetPosts(ctx *fasthttp.RequestCtx) {
	thread, err := app.Thread.SelectBySlugOrID(ctx.UserValue("slug").(string))
	if err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}

	params := getParams(ctx)
	posts, err := app.Post.GetPosts(thread.ID, params.Desc, params.Since, params.Limit, params.Sort)
	if err != nil {
		send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, posts)
}
