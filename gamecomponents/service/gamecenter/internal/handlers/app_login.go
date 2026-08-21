package handlers

import (
	"cmp"
	"game/comm"
	"game/comm/define"
	"game/comm/mq"
	"game/comm/mux"
	"game/comm/slotsmongo"
	"game/duck/ut2/jwtutil"
	"game/service/gamecenter/internal/gamedata"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/operator"
	"net/url"
	"path"
	"slices"
	"strconv"
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
		ParamsSample: v1GameLaunchPs{"abc", "XingYunXiang", "th", ""},
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
	gd := gamedata.Get()
	defer func() {
		if slices.Contains(gd.ReplaceAppid, app.AppID) {
			var args []string
			for old, new := range gd.ReplaceDoMain {
				args = append(args, old, new)
			}

			replacer := strings.NewReplacer(args...)
			ret.Url = replacer.Replace(ret.Url)
		}
	}()

	if app.Status != 1 {
		return define.NewErrCode("The Marchant has be banned.", 1001)
	}
	parentApp := operator.AppMgr.GetApp(app.ParentAppID)
	if app.ParentAppID != "admin" {
		if parentApp == nil {
			return define.NewErrCode("invalid AppID.", 1002)
		}
		if parentApp.Status != 1 {
			return define.NewErrCode("The Parent Marchant has be banned.", 1003)
		}
	}
	// err = CollGames.FindAll(bson.M{}, &ret.List)

	ps.GameID = gd.GetGameId(ps.GameID)
	if ps.Language == "" {
		ps.Language = "en"
	}
	var isfakeapp = strings.HasPrefix(app.AppID, "fake")
	game := Game{}
	{
		err = gcdb.CollGames.FindId(ps.GameID, &game)
		if err == mongo.ErrNoDocuments {
			return define.NewErrCode("The game is not found.", 1004)
		}

		switch game.Status {
		case GameStatus_Maintenance:
			if !isfakeapp {
				return define.NewErrCode("The game is under maintenance.", 1005)
			}
		case GameStatus_Closed:
			return define.NewErrCode("The game is closed.", 1006)
		case GameStatus_Hide:
			return define.NewErrCode("The game is hidden.", 1007)
		}

		// if !isfakeapp && game.Status != GameStatus_Open {
		// 	return errors.New("E405: game under maintenance.")
		// }
	}

	if ps.UserID == "" {
		return define.NewErrCode("UserID is empty.", 1008)
	}

	plr, err := operator.AppMgr.EnsureUserExists(app, ps.UserID)
	if err != nil {
		return err
	}

	if isfakeapp {
		balance, _ := plr.Balance()
		if balance < 1000000000 {
			plr.TransferIn(1000000000)
		}
	}

	// logger.Info("请求token", plr.Uid, ps.GameID)
	token, err := jwtutil.NewTokenWithData(plr.Pid, time.Now().Add(12*time.Hour), ps.GameID)
	if err != nil {
		return err
	}
	var exitInfo struct {
		ButStatus int64  // 退出状态
		ButLink   string // 游戏链接
	}
	err = mq.Invoke("/AdminInfo/Interior/GetOperatorGameMent", map[string]any{
		"AppID": app.AppID,
	}, &exitInfo)
	if err != nil {
		return
	}

	fullurl := ""
	if isPGGame(ps.GameID) {
		// https://m-pg.kafa010.com/39/index.html?ot=c6d81c1d8bfb0e214f632e6185f11e71&btt=1&l=en&ops=000102030404e44708090b853992e98c&or=static-pg.kafa010.com
		var urlobj *url.URL
		urlobj, err = url.Parse(gd.PGLaunchUrl)
		if err != nil {
			return
		}

		query := urlobj.Query()
		query.Set("l", ps.Language)
		query.Set("ops", token)
		// 获取退出按钮和退出链接
		params := comm.GetEXParams(plr.Pid, ps.GameID)
		if params.ExitBtnOff == 1 && params.ExitLink != "" {
			query.Set("f", params.ExitLink)
		}
		if params.OfficialVerify == 1 {
			query.Set("verify", "1")
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
		lang := "en"
		if _, ok := JILIIDMap[ps.GameID]; !ok {
			lang = cmp.Or(loginLang_JILI[ps.Language], "en")
		}
		query := urlobj.Query()
		query.Set("lang", lang)
		query.Set("ssoKey", token)
		query.Set("gameID", id)
		// 窗口是否关闭
		params := comm.GetEXParams(plr.Pid, ps.GameID)
		if params.OpenScreenOff == 0 {
			query.Set("itemIdx", "1")
		}
		if params.SidebarOff == 0 {
			query.Set("SidebarSwitch", "1")
		}

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

	} else if isTadaGame(ps.GameID) {
		var urlobj *url.URL
		urlobj, err = url.Parse(gd.TADALaunchUrl)
		if err != nil {
			return
		}
		id, name := splitJILI(ps.GameID)

		lang := cmp.Or(loginLang_JILI[ps.Language], "en")
		//lang := "en"
		query := urlobj.Query()
		query.Set("lang", lang)
		query.Set("ssoKey", token)
		query.Set("gameID", id)
		// 窗口是否关闭
		params := comm.GetEXParams(plr.Pid, ps.GameID)
		if params.OpenScreenOff == 0 {
			query.Set("itemIdx", "1")
		}
		if params.SidebarOff == 0 {
			query.Set("SidebarSwitch", "1")
		}

		urlobj.Path = "/" + name + "/"
		urlobj.RawQuery = query.Encode()
		fullurl = urlobj.String()
	} else if isSpribeGame(ps.GameID) {
		var urlobj *url.URL
		urlobj, err = url.Parse(gd.SPRIBELaunchUrl)
		if err != nil {
			return
		}

		lang := cmp.Or(loginLang_SPRIBE[ps.Language], "en")
		//lang := "en"
		query := urlobj.Query()
		query.Set("lang", lang)
		query.Set("token", token)
		urlobj.RawQuery = query.Encode()
		fullurl = urlobj.String()
	} else if isJDBGame(ps.GameID) {
		var urlobj *url.URL
		urlobj, err = url.Parse(gd.JDBLaunchUrl)
		if err != nil {
			return
		}

		lang := cmp.Or(loginLang_JDB[ps.Language], "en")
		gameName := cmp.Or(JDBIDGmaeMap[ps.GameID], "LuckySeven_f11946f")
		//lang := "en"
		id := splitJDB(ps.GameID)
		query := urlobj.Query()
		query.Set("lang", lang)
		query.Set("x", token)
		query.Set("gameType", id[:len(id)-3])
		query.Set("mType", id)
		query.Set("gName", gameName)
		urlobj.RawQuery = query.Encode()
		fullurl = urlobj.String()
	} else if isHACKRAWGame(ps.GameID) {
		var urlobj *url.URL
		urlobj, err = url.Parse(gd.HACKRAWLaunchUrl)
		if err != nil {
			return
		}
		id := splitJDB(ps.GameID)
		lang := cmp.Or(loginLang_HACKRAW[ps.Language], "en")
		currency := app.CurrencyKey
		if currency == "" {
			currency = "THB"
		} else if len(currency) > 3 {
			currency = currency[:3]
		}
		query := urlobj.Query()
		query.Set("language", lang)
		query.Set("token", token)
		query.Set("gameid", id)
		query.Set("currency", currency)
		// 窗口是否关闭
		params := comm.GetEXParams(plr.Pid, ps.GameID)
		if params.OpenScreenOff == 0 {
			query.Set("itemIdx", "1")
		}
		if params.SidebarOff == 0 {
			query.Set("SidebarSwitch", "1")
		}

		urlobj.Path = "/" + id + "/1.33.0/index.html"
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
		if game == "lottery" && ps.Platform == "desktop" {
			game = "lottery-pc"
		}

		gameLaunch, err := url.JoinPath(gd.LaunchUrl, game, "index.html")
		if err != nil {
			return err
		}

		fullurl = gameLaunch + "?" + args.Encode()
	}
	if exitInfo.ButStatus != 0 {
		var baseobj *url.URL
		baseobj, err = url.Parse(gd.BaseUrl)
		baseQuery := baseobj.Query()
		baseQuery.Set("game_uri", fullurl)
		baseQuery.Set("exit", strconv.FormatInt(exitInfo.ButStatus, 10))
		baseQuery.Set("link", exitInfo.ButLink)
		baseobj.RawQuery = baseQuery.Encode()
		ret.Url = baseobj.String()
	} else {
		ret.Url = fullurl
	}

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
func isTadaGame(gameID string) bool {
	return strings.HasPrefix(gameID, "tada_")
}

func isPPGame(gameID string) bool {
	return strings.HasPrefix(gameID, "pp_")
}

func isSpribeGame(gameID string) bool {
	return strings.HasPrefix(gameID, "spribe_")
}

func isJDBGame(gameID string) bool {
	return strings.HasPrefix(gameID, "jdb_")
}

func isHACKRAWGame(gameID string) bool {
	return strings.HasPrefix(gameID, "hacksaw_")
}

func splitJDB(gameID string) (id string) {
	arr := strings.SplitN(gameID, "_", 2)
	lo.Must0(len(arr) == 2)

	id = arr[1]
	return
}
