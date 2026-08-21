package main

import (
	"errors"
	"game/comm/mux"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

type Context struct {
	Request     *http.Request
	PID         int64
	Username    string
	Lang        string
	CurrencyKey string //币种
}

func msgProcGetArg0(r *http.Request) (ans any, err error) {
	path := NewPath(r.URL.Path, "admin")

	var ctx = &Context{
		Request: r,
	}

	cker := &checker{}
	allow, shouldLogin := cker.IsAllowAndShoudLogin(path)
	if !allow {
		err = errors.New("not allow")
		return
	}

	token := r.Header.Get("Authorization")
	lang := r.Header.Get("Lang")
	// if false && strings.HasPrefix(token, "user:") {
	// 	username := strings.TrimPrefix(token, "user:")
	// 	one := adminpb.DBAdminer{}
	// 	err = CollAdminUser.FindOne(bson.M{"Username": username, "Status": comm.User_Status_Normal}, &one)
	// 	if err != nil {
	// 		return
	// 	}
	// 	ctx.PID = one.ID
	// 	ctx.Username = username
	// } else {
	pid, username := cker.GetPidAndUname(token)
	if shouldLogin && pid == 0 {
		err = errors.New("用户登录已超时")
		return
	}
	ctx.PID = pid
	ctx.Username = username
	// }

	tokenExp := time.Now().Add(8 * time.Hour)
	err = CollAdminUser.UpdateId(pid, bson.M{
		"$set": bson.M{
			"TokenExpireAt": tokenExp,
		},
	})
	if err != nil {
		return nil, err
	}

	tokenMap.Delete(token)
	tokenMap.Store(token, &Token{
		Token:    token,
		ExpireAt: tokenExp,
		Pid:      pid,
	})
	ctx.Lang = lang

	ans = ctx
	return
}
func RegMsgProc(path string, desc string, kind string, handler interface{}, ps interface{}) *mux.PHandler {
	data := &mux.PHandler{
		Path:         path,
		Handler:      handler,
		Desc:         desc,
		Kind:         kind,
		ParamsSample: ps,
		Class:        "admin",
		GetArg0:      msgProcGetArg0,
	}

	return mux.DefaultRpcMux.Add(data)
}
