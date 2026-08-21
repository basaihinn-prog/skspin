package handlers

import (
	"errors"
	"game/service/gamecenter/internal/operator"
	"net/http"
)

func getArg0(r *http.Request) (ans any, err error) {
	h := r.Header

	AppID := h.Get("AppID")
	AppSecret := h.Get("AppSecret")
	if AppID == "" {
		AppID = "fake"
		// err = errors.New("invalid AppID")
		// return
	}
	if AppSecret == "" {
		AppSecret = "fake"
		// err = errors.New("invalid AppSecret")
		// return
	}

	app := operator.AppMgr.GetApp(AppID)
	if app == nil {
		err = errors.New("invalid AppID")
		return
	}
	if app.AppSecret != AppSecret {
		err = errors.New("invalid AppSecret")
		return
	}
	ans = app
	return
}
