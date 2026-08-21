package handlers

import (
	"cmp"
	"errors"
	"fmt"
	"game/comm/mq"
	"game/comm/mux"
	"game/comm/slotsmongo"
	"game/duck/ut2/jwtutil"
	"game/service/gamecenter/internal/gamedata"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/operator"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
)

// http://local-dev:27000/GetLink?GameID=XingYunXiang&UserID=abc
func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/Login",
		Handler:      login,
		Desc:         "登录",
		Kind:         "api",
		ParamsSample: v1GameLaunchPs{"abc", "XingYunXiang", "th", true, ""},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

// type LoginReq struct {
// 	UserID   string // 玩家id
// 	GameID   string // 游戏ID
// 	Language string // 游戏语言
// }

// type LoginResp struct {
// 	Url string `protobuf:"bytes,1,opt,name=Url,proto3" ` // 游戏启动链接
// }

func login(app *operator.MemApp, ps v1GameLaunchPs, ret *v1GameLaunchRet) (err error) {
	if app.Status != 0 {
		return errors.New("E408: The Marchant has be banned.")
	}

	// err = CollGames.FindAll(bson.M{}, &ret.List)
	gd := gamedata.Get()
	replaces := gamedata.GetReplaces()
	ps.GameID = replaces.GetGameId(ps.GameID)

	ps.Language = replaces.GetLang(ps.GameID, ps.Language)

	game := Game{}
	var isfakeapp = strings.HasPrefix(app.AppID, "fake")
	{
		err = gcdb.CollGames.FindId(ps.GameID, &game)
		if err == mongo.ErrNoDocuments {
			return errors.New("E404: The game is not found.")
		}

		switch game.Status {
		case GameStatus_Maintenance:
			if !isfakeapp {
				return errors.New("E405: The game is under maintenance.")
			}
		case GameStatus_Closed:
			return errors.New("E406: The game is closed.")
		case GameStatus_Hide:
			return errors.New("E407: The game is hidden.")
		}

		// if !isfakeapp && game.Status != GameStatus_Open {
		// 	return errors.New("E405: game under maintenance.")
		// }
	}

	plr, err := operator.AppMgr.EnsureUserExists(app, ps.UserID)
	if err != nil {
		return err
	}

	if plr.Status() != 0 {
		return fmt.Errorf("E408: The Player[%s] has be banned.", plr.Uid)
	}

	if isfakeapp {
		balance, _ := plr.Balance()
		if balance < 1000000 {
			plr.TransferIn(100000000)
		}
	}

	// logger.Info("请求token", plr.Uid, ps.GameID)
	token, err := jwtutil.NewTokenWithData(plr.Pid, time.Now().Add(12*time.Hour), ps.GameID)
	if err != nil {
		return err
	}

	fullurl := ""
	if isPGGame(ps.GameID) {
		// https://m-pg.kafa010.com/39/index.html?ot=c6d81c1d8bfb0e214f632e6185f11e71&btt=1&l=en&ops=000102030404e44708090b853992e98c&or=static-pg.kafa010.com

		var urltmp = gd.PGLaunchUrl
		if ps.AllowHttpProtocol && gd.PGLaunchUrlHttp != "" {
			urltmp = gd.PGLaunchUrlHttp
		}

		var urlobj *url.URL
		urlobj, err = url.Parse(urltmp)
		if err != nil {
			return
		}

		query := urlobj.Query()
		query.Set("l", ps.Language)
		query.Set("ops", token)
		if ps.ReturnUrl != "" {
			query.Set("from", ps.ReturnUrl)
		}

		urlobj.Path = path.Join(ps.GameID[3:], "index.html")
		urlobj.RawQuery = query.Encode()

		fullurl = urlobj.String()
	} else if isJILIGame(ps.GameID) {
		var urlobj *url.URL
		urlobj, err = url.Parse(gd.JILILaunchUrl)
		if err != nil {
			return
		}
		id, name := splitJILI(ps.GameID)

		lang := cmp.Or(loginLang_JILI[ps.Language], "en")
		query := urlobj.Query()
		query.Set("lang", lang)
		query.Set("ssoKey", token)
		query.Set("gameID", id)

		urlobj.Path = "/" + name + "/"
		urlobj.RawQuery = query.Encode()
		fullurl = urlobj.String()
	} else if isPPGame(ps.GameID) {
		// https://763a46cfc8.uympierc.net/gs2c/html5Game.do?jackpotid=0&extGame=1&ext=0&cb_target=exist_tab&symbol=vs20olympx&jurisdictionID=99&minilobby=false&mgckey=AUTHTOKEN@b2c4fe8ac57885ab6a85d9807a69442ec02674b3dfce39258ae80dd9dee51b00~stylename@abconline~SESSION@eb3241dc-2259-4e8a-9708-2ff018a4aeb5~SN@492c7ae1&tabName=
		var urlobj *url.URL
		urlobj, err = url.Parse(gd.PPLaunchUrl)
		if err != nil {
			return
		}

		currency := app.CurrencyKey
		if currency == "" {
			currency = "THB"
		} else if len(currency) > 3 {
			currency = currency[:3]
		}

		lang := ps.Language
		query := urlobj.Query()
		query.Set("lang", lang)
		query.Set("mgckey", token)
		query.Set("symbol", ps.GameID[len("pp_"):])
		query.Set("gname", game.Name)
		query.Set("currency", currency)

		urlobj.RawQuery = query.Encode()
		fullurl = urlobj.String()

	} else if ps.GameID == "pokdeng" {
		var lu struct {
			Url string
		}
		err = mq.Invoke("/pokdeng/game/launch", map[string]any{
			"Pid":      plr.Pid,
			"AppID":    plr.AppID,
			"Uid":      plr.Uid,
			"GameID":   ps.GameID,
			"Language": ps.Language,
		}, &lu)
		if err != nil {
			return
		}

		fullurl = lu.Url
	} else {
		args := url.Values{}
		args.Set("t", token)
		args.Set("l", ps.Language)

		game := ps.GameID
		// if game == "lottery" && ps.Platform == "desktop" {
		// 	game = "lottery-pc"
		// }

		gameLaunch, err := url.JoinPath(gd.LaunchUrl, game, "index.html")
		if err != nil {
			return err
		}

		fullurl = gameLaunch + "?" + args.Encode()
	}
	ret.Url = fullurl

	slotsmongo.UpdatePlrLoginTime(plr.Pid, ps.GameID)
	return
}

func isPGGame(gameID string) bool {
	// jili_2_csh
	return strings.HasPrefix(gameID, "pg_")
}

func splitJILI(gameID string) (id, name string) {
	arr := strings.SplitN(gameID, "_", 3)
	lo.Must0(len(arr) == 3)

	id, name = arr[1], arr[2]
	return
}

func isJILIGame(gameID string) bool {
	return strings.HasPrefix(gameID, "jili_")
}

func isPPGame(gameID string) bool {
	return strings.HasPrefix(gameID, "pp_")
}
