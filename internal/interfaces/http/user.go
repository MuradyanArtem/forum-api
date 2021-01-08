package http

import (
	"errors"
	"forum-api/internal/domain/models"
	"forum-api/internal/infrastructure"
	"net/http"

	"github.com/valyala/fasthttp"
)

func CreateUser(ctx *fasthttp.RequestCtx) {
	user := &models.User{}
	if err := unmarshall(ctx, user); err != nil {
		return
	}

	user.Nickname = ctx.UserValue("nickname").(string)
	if user.Nickname == "" {
		setStatus(ctx, http.StatusBadRequest)
		return
	}

	if err := app.User.Insert(user); err != nil {
		switch err {
		case infrastructure.ErrConflict:
			usersAlreadyExist, err := app.User.SelectByEmailOrNickname(user.Nickname, user.Email)
			if err != nil {
				send(ctx, http.StatusBadRequest, models.Message{Message: err.Error()})
				return
			}
			send(ctx, http.StatusConflict, usersAlreadyExist)

		default:
			send(ctx, http.StatusInternalServerError, models.Message{Message: err.Error()})
		}
		return
	}
	send(ctx, http.StatusCreated, user)
}

func UpdateUser(ctx *fasthttp.RequestCtx) {
	user := &models.User{}
	if err := unmarshall(ctx, user); err != nil {
		return
	}

	user.Nickname = ctx.UserValue("nickname").(string)
	if user.Nickname == "" {
		setStatus(ctx, http.StatusBadRequest)
		return
	}

	if err := app.User.Update(user); err != nil {
		switch err {
		case infrastructure.ErrConflict:
			send(ctx, http.StatusConflict, models.Message{Message: err.Error()})

		case infrastructure.ErrNotExists:
			send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})

		default:
			send(ctx, http.StatusInternalServerError, models.Message{Message: err.Error()})
		}
		return
	}
	send(ctx, http.StatusOK, user)
}

func GetUser(ctx *fasthttp.RequestCtx) {
	user := &models.User{}
	user.Nickname = ctx.UserValue("nickname").(string)
	if user.Nickname == "" {
		setStatus(ctx, http.StatusBadRequest)
		return
	}

	userInDB, err := app.User.SelectByNickname(user.Nickname)
	if err != nil {
		if errors.Is(err, infrastructure.ErrNotExists) {
			send(ctx, http.StatusNotFound, models.Message{Message: err.Error()})
			return
		}

		send(ctx, http.StatusInternalServerError, models.Message{Message: err.Error()})
		return
	}
	send(ctx, http.StatusOK, userInDB)
}
