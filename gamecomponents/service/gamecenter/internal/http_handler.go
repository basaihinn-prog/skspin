package internal

/*
import (
	"encoding/base64"
	"errors"
	"game/comm/ut/xxtea"
	"game/duck/logger"
	"game/duck/ut2/jwtutil"
	"game/pb/_gen/pb/gamepb"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GameList(c echo.Context) error {
	ret := gamepb.GameListResp{}
	err := CollGames.FindAll(bson.M{}, &ret.List)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &ret)
}

func Login(c echo.Context) error {

	req := gamepb.LoginReq{}
	if err := BindJson(c, &req); err != nil {
		return err
	}

	game := gamepb.Game{}
	err := CollGames.FindId(req.GameID, &game)
	if err == mongo.ErrNoDocuments {
		return errors.New("game not found.")
	}
	if game.Status != gamepb.GameStatus_Open {
		return errors.New("game under maintenance.")
	}

	app := GetApp(c)

	plr, err := appMgr.EnsureUserExists(app, req.UserID)
	if err != nil {
		return err
	}

	logger.Info("请求token", plr.Uid, req.GameID)
	token, err := jwtutil.NewTokenWithData(plr.Pid, time.Now().Add(12*time.Hour), req.GameID)
	if err != nil {
		return err
	}

	args := url.Values{}
	args.Set("t", token)
	args.Set("l", req.Language)

	gameLaunch, err := url.JoinPath(config.LaunchUrl, req.GameID, "index.html")
	if err != nil {
		return err
	}

	fullurl := gameLaunch + "?" + args.Encode()

	return c.JSON(http.StatusOK, &gamepb.LoginResp{Url: fullurl})
}

func GetLog(c echo.Context) error {

	req := gamepb.GetLogReq{}
	if err := BindJson(c, &req); err != nil {
		return err
	}

	app := GetApp(c)

	logger.Info("GetLog", app.AppID, req.From.AsTime())
	list, err := GetLogByCursor(req.From.AsTime(), bson.M{"AppID": app.AppID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gamepb.DocBetLogList{List: list})
}

var xx_password = []byte("/config/token.json")

// 这个函数的作用是把 websocket 连接传给客户端。
// 避免在客户端还需要将链接信息打包
func GetToken(c echo.Context) error {
	url := config.Websocket + "?t=" + c.QueryParam("t")
	byt := xxtea.Encrypt([]byte(url), xx_password)
	str := base64.StdEncoding.EncodeToString(byt)

	return c.JSON(http.StatusOK, bson.M{
		"token": str, // 伪装为token
	})
}

*/
