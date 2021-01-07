package http

import (
	"strconv"

	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type params struct {
	Limit int
	Since string
	Desc  bool
	Sort  string
}

func getParams(req *fasthttp.RequestCtx) *params {
	params := &params{}
	params.Since = string(req.FormValue("since"))
	params.Sort = string(req.FormValue("sort"))
	params.Desc = false
	if string(req.FormValue("desc")) == "true" {
		params.Desc = true
	}

	var err error
	params.Limit, err = strconv.Atoi(string(req.FormValue("limit")))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "getParams",
		}).Errorln("Cannot get query parameter: limit", err)
	}
	return params
}

func marshall(req *fasthttp.RequestCtx, any easyjson.Marshaler) {
	body, err := easyjson.Marshal(any)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "marshall",
		}).Error(err)
		req.SetContentType("application/json")
		req.SetStatusCode(fasthttp.StatusInternalServerError)
	}
	req.SetBody(body)
}

func unmarshall(req *fasthttp.RequestCtx, any easyjson.Unmarshaler) error {
	err := easyjson.Unmarshal(req.Request.Body(), any)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "http",
			"func": "unmarshal",
		}).Error(err)
		req.SetContentType("application/json")
		req.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}
	return nil
}

func setStatus(req *fasthttp.RequestCtx, status int) {
	req.SetContentType("application/json")
	req.SetStatusCode(status)
}

func send(ctx *fasthttp.RequestCtx, status int, any easyjson.Marshaler) {
	setStatus(ctx, status)
	marshall(ctx, any)
}
