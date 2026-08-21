package internal

import (
	"context"
	"fmt"
	"game/comm"
	"game/comm/mux"
	"game/pb/_gen/pb/gamepb"
	"net/http"
	"net/url"
)

func init() {
	mux.RegHttpWithSample("/GetLink", "get game launch url", "player", getlink, getlinkPs{"testuser1", "pg_39", "en", ""})
}

type getlinkPs struct {
	UserID     string `json:"user_id"`
	GameID     string `json:"game_id"`
	Language   string `json:"language"`
	GamePreFix string `json:"game_pre_fix"`
}

type getlinkRet struct {
	Token   string `json:"token"`
	GameUrl string `json:"gameUrl"`
}

type Svr struct {
	Url   string
	Appid string
	// Websocket string
	// Debug     bool
}

// var apiUrlConf = map[string]Svr{
// 	// "test": {Url: "http://127.0.0.1:11004/api", Appid: "fake", Websocket: "ws://h5.megakf.com:11080/login", Debug: true},
// 	"test": {
// 		Url:   "http://127.0.0.1:11004/api",
// 		Appid: "fake",
// 		// Websocket: "ws://192.168.1.14:11080/login",
// 		// Debug:     true,
// 	},
// 	"play": {
// 		Url:   "https://api.rpplay.com/api",
// 		Appid: "fake",
// 		// Websocket: "wss://ws.rpplay.com/login",
// 		// Debug:     false,
// 	},
// }

func getlink(_ *http.Request, ps getlinkPs, ret *getlinkRet) (err error) {
	// gold, ok := userMap.Load(ps.UserID)
	// if !ok {
	// 	gold = 10000 * 10000
	// 	userMap.Store(ps.UserID, gold)
	// }
	ensureUserExists(ps.UserID)

	full, err := url.JoinPath(gamecenterUrl, "/api/v1/game/launch")
	if err != nil {
		return err
	}
	var resp struct {
		Url string // 游戏启动链接
	}

	err = comm.PostJsonCode(context.TODO(), full, &gamepb.LoginReq{
		UserID:   ps.UserID,
		GameID:   ps.GamePreFix + ps.GameID, //"pg_"+"89"
		Language: ps.Language,
	}, &resp, func(r *http.Request) {
		r.Header.Set("AppID", appId)
		r.Header.Set("AppSecret", appSecret)
	})
	if err != nil {
		return err
	}
	u, err := url.Parse(resp.Url)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	//query := u.Query()
	//ops := query.Get("ops")
	//if len(ops) == 0 {
	//	ops = query.Get("ssoKey")
	//}
	//if len(ops) == 0 {
	//	u2, err := url.Parse(query.Get("game_uri"))
	//	if err != nil {
	//		fmt.Println("Error parsing URL:", err)
	//	}
	//	ops = u2.Query().Get("ops")
	//}

	//var ops string
	//query := u.Query()
	//parts := strings.Split(ps.GameID, "_")
	//if parts[0] == "jili" {
	//	var GameUri *url.URL
	//	GameUri, err = url.Parse(query.Get("game_uri"))
	//	ops = GameUri.Query().Get("ssoKey")
	//} else if parts[0] == "pp" {
	//	var GameUri *url.URL
	//	GameUri, err = url.Parse(query.Get("game_uri"))
	//	ops = GameUri.Query().Get("mgckey")
	//} else {
	//	var GameUri *url.URL
	//	GameUri, err = url.Parse(query.Get("game_uri"))
	//	ops = GameUri.Query().Get("ops")
	//}

	query := u.Query()
	var GameUri *url.URL
	ops := query.Get("ops")

	if ops == "" { //pg
		GameUri, err = url.Parse(query.Get("game_uri"))
		ops = GameUri.Query().Get("ops")
	}

	if ops == "" { //jili
		GameUri, err = url.Parse(query.Get("game_uri"))
		ops = GameUri.Query().Get("ssoKey")
		if ops == "" {
			ops = query.Get("ssoKey")
		}
	}

	if ops == "" { //pp
		GameUri, err = url.Parse(query.Get("game_uri"))
		ops = GameUri.Query().Get("mgckey")
		if ops == "" {
			ops = query.Get("mgckey")
		}
	}

	if ops == "" { //jdb
		GameUri, err = url.Parse(query.Get("game_uri"))
		ops = GameUri.Query().Get("x")
		if ops == "" {
			ops = query.Get("x")
		}
	}

	ret.Token = ops
	return

	//full2, err := url.JoinPath(gameapiUrl, "/gameapi/game-api/89/v2/GameInfo/Get")
	//if err != nil {
	//	return err
	//}
	//var ps1 = url.Values{}
	//ps1.Add("btt", "1")
	//ps1.Add("vc", "2")
	//ps1.Add("pf", "1")
	//ps1.Add("l", "en")
	//ps1.Add("atk", ops)
	//
	//psid, _, _ := jwtutil.ParseTokenData(ops)
	//err = comm.PostJsonCode(context.TODO(), full2, &define.PGParams{
	//	Form: ps1,
	//	Pid:  psid,
	//}, &resp, func(r *http.Request) {
	//	r.Header.Set("AppID", appId)
	//	r.Header.Set("AppSecret", appSecret)
	//})
	//if err != nil {
	//	return err
	//}
	// if !one.Debug {
	// return c.Redirect(http.StatusTemporaryRedirect, resp.Url)
	// }
	//ret.GameUrl = resp.Url
	//uResp, err := url.Parse(resp.Url)
	//if err != nil {
	//	return err
	//}
	//
	//var token struct {
	//	Token string `json:"token"`
	//}
	//
	//// uResp.Path
	//
	//u := lo.Must(url.Parse(gamecenterUrl))
	//// u = u.JoinPath("/config/token.json")
	//u.Path = "/config/token.json"
	//u.RawQuery = uResp.RawQuery
	//
	//// ul := lo.Must(url.JoinPath(one.Url, "/config/token.json"))
	//err = mux.HttpInvoke(u.String(), mux.EmptyParams{}, &token)
	//if err != nil {
	//	return
	//}
	//
	//ret.Token = token.Token
	//
	///*
	//	var xx_password = []byte("/config/token.json")
	//
	//	url := one.Websocket + "?" + u.RawQuery
	//	byt := xxtea.Encrypt([]byte(url), xx_password)
	//	str := base64.StdEncoding.EncodeToString(byt)
	//*/
	//
	//return
}
