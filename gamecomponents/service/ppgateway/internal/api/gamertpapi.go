package api

import (
	"encoding/json"
	"game/comm/define"
	"game/comm/mq"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

//func (e *PGError) Error() string {
//	return e.Msg
//}

//type PGRetWrapper struct {
//	Dt  *json.RawMessage `json:"dt"`
//	Err *PGError         `json:"err"`
//}

func gamertpapi(w http.ResponseWriter, r *http.Request) {
	var (
		ret PPRetWrapper
		err error
	)
	if r.Method == http.MethodOptions {
		return
	}

	//traceId := r.URL.Query().Get("traceId")
	defer func() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			ret.Err = &PGError{
				Cd:  "1200",
				Msg: err.Error(),
				//Tid: traceId,
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

	//err = r.ParseForm()
	//if err != nil {
	//	return
	//}
	srt, err := io.ReadAll(r.Body)
	if err != nil {

	}

	//r.Form, err = url.ParseQuery(string(srt))
	//if err != nil {
	//}
	var data struct {
		Pid      int    `json:"Pid"`
		Atk      string `json:"atk"`
		Count    int    `json:"Count"`
		BuyBonus bool   `json:"BuyBonus"`
		Bet      string `json:"Bet"`
		UName    string `json:"UName"`
		GameId   string `json:"GameId"`
		GameName string `json:"GameName"`
	}

	json.Unmarshal(srt, &data)

	r.Form = url.Values{}

	r.Form.Set("Pid", strconv.Itoa(data.Pid))
	r.Form.Set("atk", data.Atk)
	r.Form.Set("Count", strconv.Itoa(data.Count))
	r.Form.Set("BuyBonus", strconv.FormatBool(data.BuyBonus))
	r.Form.Set("Bet", data.Bet)
	r.Form.Set("UName", data.UName)
	r.Form.Set("GameId", data.GameId)
	r.Form.Set("GameName", data.GameName)

	//tk := w.FormValue("atk")
	//fmt.Println("tk", tk)
	//pid, err := jwtutil.ParseToken(tk)
	//if err != nil {
	//	err = define.NewErrCode("Invalid player session", 1302)
	//	return
	//}
	ps := &PGParams{
		Path: r.URL.Path,
		//TraceId: traceId,
		Form: r.Form,
		//Pid:  pid,
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
		onRtpSpin(ps, resp)
	}
}
