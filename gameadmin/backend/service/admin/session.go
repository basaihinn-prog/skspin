package main

import (
	"context"
	"errors"
	"fmt"
	"game/duck/lang"
	"game/duck/logger"
	"game/duck/ut2"
	"strconv"

	"google.golang.org/grpc/metadata"
)

type Session struct {
	Pid       int64  //
	Language  string //
	LogPrefix string // log前缀，方便查询
	Metadata  metadata.MD
}

func (ps *Session) IsLogin() bool {
	return ps.Pid != 0
}

func (ps *Session) Err(err string) error {
	return errors.New(lang.Get(ps.Language, err))
}

func (ps *Session) Errf(id string, data any, plural int) error {
	return errors.New(lang.Translate(ps.Language, id, data, plural))
}

func (ps *Session) getString(k string, def string) string {
	arr := ps.Metadata.Get(k)
	if len(arr) != 0 {
		return arr[0]
	}
	return def
}

var Key = struct{}{}

func GetSession(ctx context.Context) (ss *Session) {
	val := ctx.Value(Key)
	if val != nil {
		return val.(*Session)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}

	ss = &Session{}
	ss.Metadata = md
	ss.Language = ss.getString("lang", "en")

	pidstr := ss.getString("pid", "0")
	pid, err := strconv.ParseInt(pidstr, 10, 64)
	if err == nil {
		ss.Pid = pid
	} else {
		logger.Info("pid解析出错", pidstr)
	}

	ss.LogPrefix = ut2.RandomString(6)
	if ss.IsLogin() {
		ss.LogPrefix += fmt.Sprintf("[%d]", ss.Pid)
	}

	return
}
