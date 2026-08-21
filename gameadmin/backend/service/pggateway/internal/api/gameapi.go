package api

import (
	"encoding/json"
	"game/comm/define"
	"game/comm/mq"
	"game/duck/ut2/jwtutil"
	"net/http"
	"strconv"
)

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

type PGRetWrapper struct {
	Dt  *json.RawMessage `json:"dt"`
	Err *PGError         `json:"err"`
}

func gameapi(w http.ResponseWriter, r *http.Request) {
	var (
		ret PGRetWrapper
		err error
	)

	traceId := r.URL.Query().Get("traceId")
	defer func() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			ret.Err = &PGError{
				Cd:  "1200",
				Msg: err.Error(),
				Tid: traceId,
			}
			if ec, ok := err.(define.IErrcode); ok {
				code := ec.Code()
				if code != 0 {
					ret.Err.Cd = strconv.Itoa(code)
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

	tk := r.FormValue("atk")
	pid, err := jwtutil.ParseToken(tk)
	if err != nil {
		err = define.NewErrCode("Invalid player session", 1302)
		return
	}

	ps := &PGParams{
		Path:    r.URL.Path,
		TraceId: traceId,
		Form:    r.Form,
		Pid:     pid,
	}

	var resp json.RawMessage
	err = mq.Invoke(r.URL.Path, ps, &resp)
	// if err != nil {
	// 	return
	// }

	if len(resp) != 0 {
		ret.Dt = &resp
	}

	if err == nil {
		onSpin(ps, resp)
	}
}
