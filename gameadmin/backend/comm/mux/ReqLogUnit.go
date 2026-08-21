package mux

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ReqLogUnit struct {
	time    time.Time
	ReqTime string
	Elapsed string
	URI     string
	UID     string
	Err     string
	Header  http.Header
	ReqBody string
	Params  interface{}
	Result  interface{}
}

func NewReqLogUnit() *ReqLogUnit {
	return &ReqLogUnit{
		time: time.Now(),
	}
}

func (r *ReqLogUnit) Print() {
	// if lazy.CommCfg().IsDev {
	r.ReqTime = r.time.Format("2006-01-02 15:04:05.000")
	r.Elapsed = time.Since(r.time).String()

	buf, _ := json.Marshal(r)

	fmt.Println(TruncMsg(buf, 1024))
	// }
}

func TruncMsg(msg []byte, maxlen int) string {
	if len(msg) < maxlen {
		return string(msg)
	}

	return fmt.Sprintf("%s +[%d]more ", msg[:maxlen], len(msg))
}
