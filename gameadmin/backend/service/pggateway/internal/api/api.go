package api

import (
	"context"
	"encoding/json"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/define"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/ut2/jwtutil"
	"game/service/pggateway/internal/gamedata"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StartApi(addr string) {
	serveMux := http.NewServeMux()

	// https://api.pg-demo.com/game-api/jungle-delight/v2/GameInfo/Get?traceId=QONQCA11
	serveMux.HandleFunc("/game-api/", gameapi)
	serveMux.HandleFunc("/web-api/auth/session/v2/verifyOperatorPlayerSession", wrapWebApi(verifyOperatorPlayerSession, true))
	serveMux.HandleFunc("/web-api/auth/session/v2/verifySession", wrapWebApi(verifySession, true))
	// serveMux.HandleFunc("/web-api/auth/session/v2/verifySession", verifySession2)
	serveMux.HandleFunc("/web-api/game-proxy/v2/BetSummary/Get", wrapWebApi(getBetSummary, true))
	serveMux.HandleFunc("/web-api/game-proxy/v2/BetHistory/Get", wrapWebApi(getBetHistory, true))
	serveMux.HandleFunc("/web-api/game-proxy/v2/GameWallet/Get", wrapWebApi(getGameWallet, true))
	serveMux.HandleFunc("/back-office-proxy/Report/GetBetHistory", wrapWebApi(boGetBetHistory, false))
	serveMux.HandleFunc("/web-api/game-proxy/v2/GameName/Get", gameName)
	serveMux.HandleFunc("/web-api/game-proxy/v2/GameRule/Get", gameRule)
	serveMux.HandleFunc("/web-api/game-proxy/v2/Resources/GetByReferenceIdsResourceTypeIds", getByReferenceIdsResourceTypeIds)
	serveMux.HandleFunc("/web-api/game-proxy/v2/Resources/GetByResourcesTypeIds", getByResourcesTypeIds)
	serveMux.HandleFunc("/game-api/piggy-gold/v2/GameInfo/Get", pg39GameInfo)
	// serveMux.HandleFunc("/game-api/39/v2/GameInfo/Get", pg39GameInfo)

	go func() {
		err := http.ListenAndServe(addr, serveMux)
		log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}

func wrapWebApi[R any](fn func(*PGParams, *R) error, verify bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
