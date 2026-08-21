package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/mux"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/comm/ut/xxtea"
	"game/duck/mongodb"
	"game/duck/ut2/httputil"
	"game/duck/ut2/jwtutil"
	"game/service/gamecenter/internal/gamedata"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/operator"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	// mux.RegHttpWithSample("/config/token.json", "get token", "api", getWSEndpoint, nil)
	mux.RegHttpWithSample("/GetLinkDev", "get link", "api", getWSEndpointDev, getWSEndpointDevPs{
		UserID:   "testuser1",
		GameID:   "YingCaiShen",
		Language: "en",
		AppID:    "fake",
	}).SetOnlyDev()
}

type configtokenRet struct {
	Token string `json:"token"`
	// Pid   int64  `json:"pid"`
}

var xx_password = []byte("/config/token.json")
var CollGameLoginDetail *mongodb.Collection //游戏登录日志

func HttpGetWSEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("rp-test", "xxx")
	var ret configtokenRet
	err := getWSEndpoint(r, mux.EmptyParams{}, &ret)

	slog.Info("HttpGetWSEndpoint", "path", r.URL.Path, "query", r.URL.RawQuery, "ret", ret, "error", ut.ErrString(err))
	if err != nil {
		// http.Error(w, err.Error(), http.StatusUnauthorized)
		mux.HttpReturn2(w, nil, err)
		return
	}

	mux.HttpReturn2(w, &ret, nil)
}

// 将ws链接加密伪装
func getWSEndpoint(req *http.Request, _ mux.EmptyParams, ret *configtokenRet) (err error) {
	pid, gameid, _, err := validateToken(req)
	if err != nil {
		return
	}
	// *ret = "hello"
	config := gamedata.Get()
	replaces := gamedata.GetReplaces()
	blocked, ip, loc := replaces.IsBlockLoc(req)
	fmt.Printf("userid=%d,gameid=%s,ip=%s,loc=%s\n", pid, gameid, ip, loc)
	insertLoginDetail(pid, gameid, ip, loc)
	// if locblock(req) {
	if blocked {
		err = errors.New("Your country or region is restricted")
		return
	}
	// t := req.URL.Query().Get("t")
	qs := req.URL.Query().Encode()
	// url := config.Websocket + "?t=" + t
	url := config.Websocket + "?" + qs
	byt := xxtea.Encrypt([]byte(url), xx_password)
	str := base64.StdEncoding.EncodeToString(byt)
	ret.Token = str
	return
}

func getSecretWS(u *url.URL) string {
	qs := u.RawQuery
	config := gamedata.Get()
	url := config.Websocket + "?" + qs
	byt := xxtea.Encrypt([]byte(url), xx_password)
	str := base64.StdEncoding.EncodeToString(byt)
	return str
}

type getWSEndpointDevPs struct {
	AppID    string
	UserID   string // 玩家id
	GameID   string // 游戏ID
	Language string // 游戏语言
}

type getWSEndpointDevRet struct {
	Token   string `json:"token"`
	GameUrl string `json:"gameUrl"`
}

// https://gamecenter.kafa010.com/GetLinkDev?AppID=faketrans-IDR&UserID=testuser1&GameID=Olympus&Language=en

func getWSEndpointDev(ps getWSEndpointDevPs, ret *getWSEndpointDevRet) (err error) {
	if ps.AppID == "" {
		ps.AppID = "fake"
	}

	app := operator.AppMgr.GetApp(ps.AppID)
	if app == nil {
		err = errors.New("invalid AppID")
		return
	}

	var loginResp v1GameLaunchRet
	err = login(app, v1GameLaunchPs{
		UserID:   ps.UserID,
		GameID:   ps.GameID,
		Language: ps.Language,
	}, &loginResp)
	if err != nil {
		return
	}

	ret.GameUrl = loginResp.Url

	u := lo.Must(url.Parse(ret.GameUrl))
	ret.Token = getSecretWS(u)

	return
}

func validateToken(r *http.Request) (pid int64, game, language string, err error) {
	var arg struct {
		Token string `query:"t"`
		L     string `query:"l"`
	}

	if err = httputil.HttpBindQuery(&arg, r); err != nil {
		return
	}

	pid, game, err = jwtutil.ParseTokenData(arg.Token)
	if err != nil {
		return
	}

	language = lo.Ternary(arg.L != "", arg.L, "th")
	return
}

func insertLoginDetail(pid int64, gameid, ip, loc string) {
	appId, uid, err := slotsmongo.GetPlayerInfo(pid)
	if err == nil {
		ld := comm.GameLoginDetail{
			ID:        primitive.NewObjectID(),
			Pid:       appId,
			UserID:    uid,
			GameID:    gameid,
			Ip:        ip,
			Loc:       ut.CountryNameByCode(loc),
			LoginTime: time.Now().Unix(),
		}
		if CollGameLoginDetail == nil {
			CollGameLoginDetail = gcdb.NewOtherDB("GameAdmin").Collection("GameLoginDetail")
		}
		err := CollGameLoginDetail.InsertOne(ld)
		if err != nil {
			fmt.Printf("loginDetail.InsertOne occured an error => %s", err.Error())
		}
	} else {
		fmt.Printf("Players.FindId occured an error => %s", err.Error())
	}
}
