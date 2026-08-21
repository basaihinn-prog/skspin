package handlers

import (
	"game/comm/define"
	"game/comm/mq"
	"game/service/gamecenter/internal/gamedata"
	"game/service/gamecenter/internal/operator"
	"log/slog"
	"net/http"
)

func getArg0(r *http.Request) (ans any, err error) {
	h := r.Header

	AppID := h.Get("AppID")
	AppSecret := h.Get("AppSecret")
	if AppID == "" {
		AppID = "fake"
		// err = errors.New("invalid AppID")
		// return
	}
	if AppSecret == "" {
		AppSecret = "fake"
		// err = errors.New("invalid AppSecret")
		// return
	}

	app := operator.AppMgr.GetApp(AppID)
	if app == nil {
		err = define.NewErrCode("invalid AppID", 1002)
		return
	}
	if app.AppSecret != AppSecret {
		err = define.NewErrCode("invalid AppSecret", 1011)
		return
	}
	if app.OperatorType == 1 {
		err = define.NewErrCode("Line provider does not allow access", 1013)
		return
	}
	if app.ReviewStatus != 1 {
		err = define.NewErrCode("App is not review", 1022)
		return
	}
	config := gamedata.Get()
	_, ip, _ := config.IsBlockLoc(r)
	slog.Info("access network",
		"ip address", ip,
	)
	netAllow, err := checkIp(AppID, ip)
	if !netAllow {
		err = define.NewErrCode("ip network not allow access", 1014)
	}
	ans = app
	return
}

func checkIp(appID, ip string) (allow bool, err error) {
	var netAllow struct {
		Allow bool
	}
	err = mq.Invoke("/AdminInfo/Interior/operatorVerifyIP", map[string]any{
		"AppID": appID,
		"IP":    ip,
	}, &netAllow)
	if err != nil {
		return
	}
	allow = netAllow.Allow
	return
}

// getArg1 检查用户白名单
//func getArg1(r *http.Request) (ans any, err error) {
//	h := r.Header
//
//	appID := h.Get("AppID")
//	ip := httputil.GetIPFromRequest(r)
//	netAllow, err := checkIp(appID, ip)
//	if !netAllow {
//		err = errors.New("ip")
//	}
//	app, err := getArg0(r)
//	if err != nil {
//		return
//	}
//	ans = app
//	return
//}
