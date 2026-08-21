package main

import (
	"encoding/json"
	"errors"
	"game/comm/mq"
	"game/comm/ut"
	"io"
	"log/slog"
	"net/http"

	"github.com/samber/lo"
)

func adminMQHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := invokeMQ(r)

	m := map[string]any{}
	if err != nil {
		m["error"] = err.Error()
	} else {
		m["data"] = resp
	}

	buf, _ := json.Marshal(m)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func invokeMQ(r *http.Request) (resp json.RawMessage, err error) {
	token := r.Header.Get("Authorization")

	var ctx = &Context{
		Request: r,
	}
	// if strings.HasPrefix(token, "user:") {
	// 	username := strings.TrimPrefix(token, "user:")
	// 	one := adminpb.DBAdminer{}
	// 	err = CollAdminUser.FindOne(bson.M{"Username": username}, &one)
	// 	if err != nil {
	// 		return
	// 	}
	// 	ctx.PID = one.ID
	// 	ctx.Username = username
	// } else {
	cker := &checker{}
	pid, username := cker.GetPidAndUname(token)
	if pid == 0 {
		err = errors.New("用户登录已超时")
		return
	}
	ctx.PID = pid
	ctx.Username = username
	// }

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	slog.Info("mq info", "url", r.URL.Path, "Pid", pid, "UserName", username, "params", string(payload))

	if len(payload) != 0 && payload[0] == '{' {
		obj := map[string]json.RawMessage{}
		err = json.Unmarshal(payload, &obj)
		if err != nil {
			return
		}

		obj["proxy_pid"] = lo.Must(ut.GetJsonRaw(pid))
		obj["proxy_username"] = lo.Must(ut.GetJsonRaw(username))
		obj["proxy_language"] = lo.Must(ut.GetJsonRaw(ctx.Lang))

		currentUser, isAdmin := IsAdminUser(ctx)

		isAdmin = currentUser.AppID == "admin"
		obj["proxy_isadmin"] = lo.Must(ut.GetJsonRaw(isAdmin))

		err = mq.Invoke(r.URL.Path, obj, &resp)
	} else {
		err = mq.Invoke(r.URL.Path, json.RawMessage(payload), &resp)
	}

	return
}
