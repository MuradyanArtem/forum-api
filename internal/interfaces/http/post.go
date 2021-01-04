package http

import (
	"bytes"
	"forum-api/internal/app"
	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"
	"net/http"
	"strconv"

	json "github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Post struct {
	post   *app.Post
	user   *app.User
	thread *app.Thread
	forum  *app.Forum
}

func newPost(post *app.Post, user *app.User, thread *app.Thread, forum *app.Forum) *Post {
	return &Post{
		post,
		user,
		thread,
		forum,
	}
}

func (p *Post) Create(ctx *fasthttp.RequestCtx) {
	id, forum, err := p.thread.GetForumIDBySlug(ctx.UserValue("slug").(string))
	if err != nil {
		marshall(ctx, models.Message{err.Error()})
		setStatus(ctx, http.StatusNotFound)
		return
	}

	posts := &models.PostSlice{}
	if err := unmarshal(ctx, posts); err != nil {
		setStatus(ctx, http.StatusBadRequest)
		return
	}

	if err := p.post.InsertPost(posts, forum, id); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "Create",
		}).Error(err)

		switch err {
		case infrastructure.ErrConflict:
			marshall(ctx, models.Message{err.Error()})
			setStatus(ctx, http.StatusConflict)
		default:
			marshall(ctx, models.Message{err.Error()})
			setStatus(ctx, http.StatusNotFound)
		}
		return
	}

	marshall(ctx, posts)
	setStatus(ctx, http.StatusOK)
}

func (m *PostManager) Update(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, _ := strconv.Atoi(idStr)
	post := &models.Post{}
	if err := json.Unmarshal(ctx.PostBody(), post); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "` + "unmarshal not ok : " + err.Error() + `"}`))
		return
	}
	post.ID = id
	err := m.pUC.Update(post)
	switch err {
	case nil:
		resp, _ := post.MarshalJSON()
		utils.Send(200, ctx, resp)
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}
}

func (m *PostManager) GetByID(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	related := ctx.QueryArgs().Peek("related")
	idStr := ctx.UserValue("id").(string)
	id, _ := strconv.Atoi(idStr)
	details := &models.PostDetails{}
	var err error
	details.Post, err = m.pUC.SelectPostByID(id)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	user, err := m.uUC.SelectByNickname(details.Post.Author)
	details.User = &user
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	details.Thread, err = m.tUC.SelectByID(details.Post.Thread)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	details.Forum, err = m.fUC.SelectBySlug(details.Thread.Forum)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	if !bytes.Contains(related, []byte("user")) {
		details.User = nil
	}
	if !bytes.Contains(related, []byte("forum")) {
		details.Forum = nil
	}
	if !bytes.Contains(related, []byte("thread")) {
		details.Thread = nil
	}

	resp, _ := details.MarshalJSON()
	ctx.Write(resp)
	ctx.SetStatusCode(200)
}

func (m *PostManager) GetPosts(ctx *fasthttp.RequestCtx) {
	params := utils.MustGetParams(ctx)
	slug := ctx.UserValue("slugOrID").(string)
	thread, err := m.tUC.SelectBySlugOrID(slug)
	if err != nil {
		utils.Send(404, ctx, utils.MustMarshalError(err))
		return
	}
	posts, err := m.pUC.GetPosts(thread.ID, params.Desc, params.Since, params.Limit, params.Sort)
	switch err {
	case nil:
		resp, _ := json.Marshal(posts)
		utils.Send(200, ctx, resp)
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}
}
