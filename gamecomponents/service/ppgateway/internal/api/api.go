package api

import (
	"context"
	"encoding/json"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/define"
	"game/comm/mq"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/lazy"
	"game/duck/ut2/jwtutil"
	"game/service/ppgateway/internal/gamedata"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StartApi(addr string) {
	serveMux := http.NewServeMux()

	// https://api.pg-demo.com/game-api/jungle-delight/v2/GameInfo/Get?traceId=QONQCA11
	// serveMux.HandleFunc("POST /game-api/", gameapi)
	//serveMux.HandleFunc("/game-api/", gameapi)
	serveMux.HandleFunc("/web-api/ppgame/hello", Hello)
	serveMux.HandleFunc("/gs2c/v3/gameService", gameapi)
	serveMux.HandleFunc("/gs2c/saveSettings", SaveSettingsDo)
	serveMux.HandleFunc("/gs2c/announcements/unread/", Unread)
	serveMux.HandleFunc("/gs2c/reloadBalance", ReloadBalanceDo)
	serveMux.HandleFunc("/gs2c/promo/frb/available/", PromoFrbAvailable)
	serveMux.HandleFunc("/gs2c/promo/active/", PromoActive)
	serveMux.HandleFunc("/gs2c/minilobby/games", MinilobbyGames)
	serveMux.HandleFunc("/gs2c/promo/tournament/details/", PromoTournamentDetails)
	serveMux.HandleFunc("/gs2c/promo/race/details/", PromoRaceDetails)
	serveMux.HandleFunc("/gs2c/promo/tournament/scores/", PromoTournamentScores)
	serveMux.HandleFunc("/gs2c/promo/race/prizes/", PromoRacePrizes)
	serveMux.HandleFunc("/ReplayService/api/top/winnings/list", getBetHistory)
	//
	go func() {
		//证书
		certfile, _ := lazy.RouteFile.Get("certfile")
		keyfile, _ := lazy.RouteFile.Get("keyfile")
		err := http.ListenAndServeTLS(addr, certfile, keyfile, serveMux)
		//err := http.ListenAndServe(addr, serveMux)
		log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}

func wrapWebApi[R any](fn func(*PGParams, *R) error, verify bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ret PPRetWrapper
			err error
		)

		traceId := r.URL.Query().Get("traceId")
		defer func() {
			w.Header().Set("Access-Control-Allow-Origin", "*")
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

		var pid int64

		if verify {
			tk := r.FormValue("atk")
			if r.URL.Path == "/web-api/auth/session/v2/verifyOperatorPlayerSession" {
				tk = r.FormValue("os")
			} else if r.URL.Path == "/web-api/auth/session/v2/verifySession" {
				tk = r.FormValue("tk")
			}

			pid, err = jwtutil.ParseToken(tk)
			if err != nil {
				err = define.NewErrCode("Invalid player session", 1302)
				return
			}
		}

		blocked, ip, loc := gamedata.IsBlockLoc(r)
		gi := r.FormValue("gi")
		fmt.Printf("userid=%d,gameid=%s,ip=%s,loc=%s\n", pid, "pg_"+gi, ip, loc)
		insertLoginDetail(pid, gi, ip, loc)
		if blocked {
			err = define.NewErrCode("Your country or region is restricted", 1306)
			return
		}

		ps := &PGParams{
			Path:    r.URL.Path,
			TraceId: traceId,
			Form:    r.Form,
			Pid:     pid,
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

func insertLoginDetail(pid int64, gameid, ip, loc string) {
	appId, uid, err := slotsmongo.GetPlayerInfo(pid)
	if err == nil {
		ld := comm.GameLoginDetail{
			ID:        primitive.NewObjectID(),
			Pid:       appId,
			UserID:    uid,
			GameID:    "pg_" + gameid,
			Ip:        ip,
			Loc:       ut.CountryNameByCode(loc),
			LoginTime: time.Now().Unix(),
		}
		CollGameLoginDetail := db.Collection2("GameAdmin", "GameLoginDetail")
		_, err := CollGameLoginDetail.InsertOne(context.TODO(), ld)
		if err != nil {
			log.Printf("loginDetail.InsertOne occured an error => %s", err.Error())
		}
	} else {
		log.Printf("Players.FindId occured an error => %s", err.Error())
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	ret := `{"hello:"HHHeello","err":null}`

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ret))
}

var modifyGoldResp struct {
	Balance int64
}

func ReloadBalanceDo(w http.ResponseWriter, r *http.Request) {
	var (
		ret PPRetWrapper
		err error
	)
	traceId := r.URL.Query().Get("traceId")
	defer func() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
		//jsondata, err := json.Marshal(ret)
		//if err != nil {
		//	fmt.Println(err)
		//}
		if err != nil {
			marshal, _ := json.Marshal(*ret.Err)
			w.Write(marshal)
		} else {
			w.Write(*ret.Dt)
		}

	}()
	tk := r.FormValue("mgckey")
	if len(tk) == 0 {
		err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	pid, err := jwtutil.ParseToken(tk)
	if err != nil {
		err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	err = mq.Invoke("/gamecenter/player/getBalance", map[string]any{
		"Pid": pid,
	}, &modifyGoldResp)
	if err != nil {
		return
	}
	balance := ut.Gold2Money(modifyGoldResp.Balance)
	data := json.RawMessage(fmt.Sprintf("balance_bonus=0.00&balance=%.1f&balance_cash=%.1f&stime=%v&", balance, balance, time.Now().UnixMilli()))
	ret.Dt = &data
}

func SaveSettingsDo(w http.ResponseWriter, r *http.Request) { //返回游戏设置，存在哪里？
	var (
		ret *json.RawMessage
		err error
	)
	traceId := r.URL.Query().Get("traceId")
	defer func() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err != nil {
			json.Unmarshal([]byte(err.Error()), &ret)
		}
		//var jsondata []byte
		//if ret == nil {
		//jsondata, _ := json.Marshal(ret)

		//} else {
		//	jsondata, _ = json.Marshal(settings)
		//}
		w.Write(*ret)
	}()

	//err = r.ParseForm()
	//if err != nil {
	//	return
	//}
	tk := r.FormValue("mgckey")
	pid, err := jwtutil.ParseToken(tk)
	if err != nil {
		err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	gameId := r.FormValue("id")
	r.URL.Path += "/" + gameId
	ps := &PGParams{
		Path:    r.URL.Path,
		TraceId: traceId,
		Form:    r.Form,
		Pid:     pid,
	}

	var resp json.RawMessage
	//获取settings
	err = mq.Invoke(r.URL.Path, ps, &resp)
	if err != nil {
		return
	}

	if len(resp) != 0 {
		ret = &resp
	}
}

func ReplayService(w http.ResponseWriter, r *http.Request) { //todo 待实现

}
