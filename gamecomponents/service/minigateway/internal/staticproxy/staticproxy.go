package staticproxy

import (
	"encoding/json"
	"errors"
	"game/comm/define"
	"game/duck/ut2/jwtutil"
	"log"
	"net/http"
	"strconv"
)

func StartProxy(addr string) {
	go func() {
		httpmux := http.NewServeMux()
		httpmux.HandleFunc("GET /", staticAssets)
		httpmux.HandleFunc("GET /aviator/huidustg", huidustg)
		httpmux.HandleFunc("POST /aviator/playerSession", wrapWebApi(aviatorPlayerSession, true))
		////证书
		//certfile, _ := lazy.RouteFile.Get("certfile")
		//keyfile, _ := lazy.RouteFile.Get("keyfile")
		//err := http.ListenAndServeTLS(addr, certfile, keyfile, httpmux)
		err := http.ListenAndServe(addr, httpmux)
		log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}

func wrapWebApi[R any](fn func(*PGParams, *R) error, verify bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ret MiniRetWrapper
			err error
		)

		traceId := r.URL.Query().Get("traceId")
		defer func() {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			if err != nil {
				ret.Err = &PGError{}
				if errors.As(err, &ret.Err) && err.(*PGError).Cd == "5000" {
					errors.As(err, &ret.Err)
				}
				if ec, ok := err.(define.IErrcode); ok {
					code := ec.Code()
					if code != 0 {
						ret.Err.Cd = strconv.Itoa(code)
						ret.Err.Msg = err.Error()
					}
				}
			}
			jsondata, _ := json.Marshal(ret)
			w.Write(jsondata)
		}()

		err = r.ParseForm()
		if err != nil {
			return
		}

		var pid int64

		tk := r.FormValue("otk")
		//if r.URL.Path == "/web-api/auth/session/v2/aviatorPlayerSession" {
		//	tk =
		//} else if r.URL.Path == "/web-api/auth/session/v2/verifySession" {
		//	tk = r.FormValue("tk")
		//}

		pid, err = jwtutil.ParseToken(tk)
		if err != nil {
			err = define.NewErrCode("Invalid player session", 1302)
			return
		}

		gi := r.FormValue("gi")

		ps := &PGParams{
			Path:    r.URL.Path,
			TraceId: traceId,
			Form:    r.Form,
			Pid:     pid,
			GameId:  "pg_" + gi,
		}

		var ans R
		err = fn(ps, &ans)
		if err != nil {
			return
		}

		resp, _ := json.Marshal(ans)
		if len(resp) != 0 {
			raw := json.RawMessage(resp)
			ret.Dt = &raw
		}
	}
}

// https://api.pg-demo.com/web-api/auth/session/v2/verifyOperatorPlayerSession?traceId=VXTNFQ12
func aviatorPlayerSession(ps *PGParams, ret *M) (err error) {
	os := ps.Form.Get("otk")
	tk := os
	s := M{
		"tk": tk,
	}
	*ret = s
	return
}

type PGParams = define.PGParams

//	{
//	    "dt": null,
//	    "err": {
//	        "cd": "1200",
//	        "msg": "OERR: Operator return an error. Failed to verify operator player session",
//	        "tid": "HXTZJF11"
//	    }
//	}

type PGError struct {
	Cd  string `json:"cd"`
	Msg string `json:"msg"`
	Tid string `json:"tid"`
}

func (e *PGError) Error() string {
	return e.Msg
}

type MiniRetWrapper struct {
	Dt  *json.RawMessage `json:"dt"`
	Err *PGError         `json:"err"`
}

type M = map[string]any
type D = []any
