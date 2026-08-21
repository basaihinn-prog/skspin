package main

import (
	"strings"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/GetTimeUtc", "获取时区", "AdminInfo", getTimeUtc, getTimeUtcParams{})
}

type getTimeUtcParams struct {
}

type getTimeUtcResult struct {
	Str string
}

func getTimeUtc(ctx *Context, ps getTimeUtcParams, ret *getTimeUtcResult) (err error) {
	now := time.Now()
	ret.Str = strings.Split(now.String(), " ")[2]
	return nil
}
