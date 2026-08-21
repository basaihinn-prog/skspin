package main

import (
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/lazy"
	_ "game/service/fakeapp/internal"
)

func main() {
	lazy.InitWithoutGrpc("fakeapp")
	//lazy.Init("fakeapp")
	mq.ConnectServerMust()
	mux.StartHttpServer(":27000")
	lazy.SignalProc()
	// cli = lazy.NewGrpcClient("gamecenter", gamepb.NewForCustomerRpcClient)

	// lazy.ServeFn(httpServer)
}

// 这里是用来测试的。虚拟了一个App
// app 是指 来接入本平台的产品

/*
var userMap = ut2.NewSyncMap[string, int64]()

// var cli gamepb.ForCustomerRpcClient

type Svr struct {
	Url       string
	Appid     string
	Websocket string
	Debug     bool
}

var apiUrlConf = map[string]Svr{
	// "test": {Url: "http://127.0.0.1:11004/api", Appid: "fake", Websocket: "ws://h5.megakf.com:11080/login", Debug: true},
	"test": {Url: "http://192.168.1.14:11004/api", Appid: "fake", Websocket: "ws://192.168.1.14:11080/login", Debug: true},
	"play": {Url: "https://api.rpplay.com/api", Appid: "fake", Websocket: "wss://ws.rpplay.com/login", Debug: false},
}

func CrossOrigin(c echo.Context) {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS*")
}

func httpServer() {
	e := comm.NewEcho()
	e.POST("/GetBalance", GetBalance)
	e.POST("/ModifyGold", ModifyGold)
	e.GET("/GetLink", GetLink)
	e.GET("/ping", func(c echo.Context) error { return c.JSON(http.StatusOK, "hello") })
	e.Start(":27000")
}

func GetLink(c echo.Context) error {
	CrossOrigin(c)
	var req struct {
		UserID   string
		GameID   string
		Language string
		Svr      string
	}
	if err := httputil.HttpBindQuery(&req, c.Request()); err != nil {
		return err
	}

	gold, ok := userMap.Load(req.UserID)
	if !ok {
		gold = 10000 * 10000
		userMap.Store(req.UserID, gold)
	}

	if req.Svr == "" {
		req.Svr = "test"
	}

	one := apiUrlConf[req.Svr]

	full, err := url.JoinPath(one.Url, "Login")
	if err != nil {
		return err
	}

	logger.Info("xx", ut2.ToJson(req))

	resp := gamepb.LoginResp{}

	err = comm.PostJsonCode(context.TODO(), full, &gamepb.LoginReq{
		UserID:   req.UserID,
		GameID:   req.GameID,
		Language: req.Language,
	}, &resp, func(r *http.Request) {
		r.Header.Set("AppID", one.Appid)
		r.Header.Set("AppSecret", one.Appid)
	})
	if err != nil {
		return err
	}

	if !one.Debug {
		return c.Redirect(http.StatusTemporaryRedirect, resp.Url)
	}

	u, err := url.Parse(resp.Url)
	if err != nil {
		return err
	}

	var xx_password = []byte("/config/token.json")

	url := one.Websocket + "?" + u.RawQuery
	byt := xxtea.Encrypt([]byte(url), xx_password)
	str := base64.StdEncoding.EncodeToString(byt)

	return c.JSON(http.StatusOK, bson.M{
		"token": str, // 伪装为token
	})
}

func GetBalance(c echo.Context) error {
	req := gamepb.GetBalanceUidReq{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	gold, ok := userMap.Load(req.UserID)
	if !ok {
		return errors.New("用户不存在")
	}

	logger.Info("GetBalance", req.UserID, gold)

	return c.JSON(http.StatusOK, &gamepb.GetBalanceUidResp{Balance: gold})
}

func ModifyGold(c echo.Context) error {

	req := gamepb.ModifyGoldUidReq{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	gold, ok := userMap.Load(req.UserID)
	if !ok {
		return errors.New("用户不存在")
	}

	gold += req.Change
	userMap.Store(req.UserID, gold)

	logger.Info("ModifyGold", req.UserID, gold)

	return c.JSON(http.StatusOK, &gamepb.GetBalanceUidResp{Balance: gold})
}

*/
