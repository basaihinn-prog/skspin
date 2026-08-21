package playerfish

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/fish_comm/mediator"
	"serve/fish_comm/player"
	"serve/fish_comm/vip"
	"serve/service_fish/application/accountingmanager"
	"serve/service_fish/application/botuser"
	"serve/service_fish/application/gameroom"
	"serve/service_fish/application/gamesetting"
	"serve/service_fish/application/gameuser"
	"serve/service_fish/application/lobbyroom"
	"serve/service_fish/application/lobbysetting"
	"serve/service_fish/application/lottery"
	"serve/service_fish/application/refund"
	"serve/service_fish/models"
	"strings"
	"time"

	fastHttpWebsocket "github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sasha-s/go-deadlock"
	"github.com/valyala/fasthttp"
)

var ControllerFish = &controllerFish{
	port:             os.Getenv("HTTP_PORT"),
	startTime:        time.Now().Format("02-Jan-2006 15:04:05"),
	version:          os.Getenv("VERSION"),
	sslPath:          os.Getenv("SSL_PATH"),
	sslCrt:           os.Getenv("SSL_CRT"),
	sslKey:           os.Getenv("SSL_KEY"),
	gorillaUpgrader:  player.Controller.GorillaUpgrader,
	fastHttpUpgrader: player.Controller.FastHttpUpgrader,
}

type controllerFish struct {
	port             string
	startTime        string
	version          string
	sslPath          string
	sslCrt           string
	sslKey           string
	endPoint         string
	gorillaUpgrader  *websocket.Upgrader
	fastHttpUpgrader *fastHttpWebsocket.FastHTTPUpgrader
}

func init() {
	deadlock.Opts.DeadlockTimeout = 3 * time.Minute
	ControllerFish.endPoint = "fish"
}

func (cf *controllerFish) exec(secWebSocketKey, hostExtId, remoteAddr, userAgent string, conn *websocket.Conn) {
	go lobbysetting.Service.New(secWebSocketKey, hostExtId, remoteAddr, userAgent)

	go gamesetting.Service.New(secWebSocketKey, hostExtId, remoteAddr, userAgent)

	flux.Send(gameuser.ActionGameUserCreate, player.Controller.Id, gameuser.Service.Id,
		hostExtId,
		secWebSocketKey,
		conn,
	)
}

func (cf *controllerFish) FluxActionHandler(action *flux.Action) {
	switch action.Key().Name() {
	case gameuser.ActionJoinLobbyRoomCall:
		secWebSocketKey := action.Key().From()
		gameId := action.Payload()[0].(string)
		flux.Send(lobbyroom.ActionLobbyRoomAutoJoinCall, player.Controller.Id, lobbyroom.Service.Id, secWebSocketKey, gameId)

	case lobbyroom.ActionLobbyRoomAutoJoinRecall:
		secWebSocketKey := action.Payload()[0].(string)
		lobbyRoomUuid := action.Payload()[1].(string)
		flux.Send(gameuser.ActionJoinLobbyRoomRecall, player.Controller.Id, secWebSocketKey, lobbyRoomUuid)

	case gameuser.ActionJoinGameRoomCall:
		secWebSocketKey := action.Key().From()
		configRecall := action.Payload()[0]
		hostExtId := action.Payload()[1].(string)

		j := gameroom.NewJoinGameRoom(
			gameroom.ActionRoomAutoJoinCall,
			player.Controller.Id,
			hostExtId,
			secWebSocketKey,
			gamesetting.Service.GameId(secWebSocketKey),
			gamesetting.Service.MathModuleId(secWebSocketKey),
			gamesetting.Service.BetList(secWebSocketKey),
			gamesetting.Service.RateList(secWebSocketKey),
			gamesetting.Service.Rate(secWebSocketKey),
			make([]interface{}, 1),
		)
		j.ExtraData[0] = configRecall

		flux.Send(gameroom.ActionRoomAutoJoinCall, player.Controller.Id, gameroom.Service.Id, j)

	case gameroom.ActionRoomAutoJoinRecall:
		j := action.Payload()[0].(*gameroom.JoinGameRoom)

		configRecall := j.ExtraData[0]

		flux.Send(refund.ActionRefundInit, player.Controller.Id, refund.Service.Id,
			j.SecWebSocketKey,
			accountingmanager.Service.Get(j.SecWebSocketKey),
		)

		flux.Send(gameuser.ActionJoinGameRoomRecall, player.Controller.Id, j.SecWebSocketKey,
			configRecall,
			j.GameRoomUuid,
			j.SeatId,
			j.NextScene,
		)

		gameId := gamesetting.Service.GameId(j.SecWebSocketKey)
		subGameId := gamesetting.Service.SubgameId(j.SecWebSocketKey)

		// 處理捕魚大排檔 富豪VIP廳沒有機器人的部分
		if cf.checkJoinBot(gameId, subGameId) {
			flux.Send(botuser.ActionBotUserCheck, player.Controller.Id, botuser.Service.Id,
				j.GameRoomUuid,
				j.Players,
			)
		}

	case botuser.ActionBotUserJoinGameRoomCall:
		botUserUuid := action.Key().From()
		gameRoomUuid := action.Payload()[0].(string)

		j := gameroom.NewJoinGameRoom(
			gameroom.ActionGameRoomBotJoinCall,
			player.Controller.Id,
			"BOT",
			botUserUuid,
			gamesetting.Service.GameId(botUserUuid),
			gamesetting.Service.MathModuleId(botUserUuid),
			gamesetting.Service.BetList(botUserUuid),
			gamesetting.Service.RateList(botUserUuid),
			gamesetting.Service.Rate(botUserUuid),
			nil,
		)
		j.GameRoomUuid = gameRoomUuid

		flux.Send(gameroom.ActionGameRoomBotJoinCall, player.Controller.Id, gameroom.Service.Id, j)

	case gameroom.ActionGameRoomBotJoinRecall:
		j := action.Payload()[0].(*gameroom.JoinGameRoom)

		botUserUuid := j.SecWebSocketKey

		// bot join game room failed
		if j.Players == nil {
			flux.Send(botuser.ActionBotUserDelete, player.Controller.Id, botuser.Service.Id, botUserUuid)
			return
		}

		gameId := gamesetting.Service.GameId(j.SecWebSocketKey)
		subGameId := gamesetting.Service.SubgameId(j.SecWebSocketKey)

		// 處理捕魚大排檔 富豪VIP廳沒有機器人的部分
		if cf.checkJoinBot(gameId, subGameId) {
			flux.Send(botuser.ActionBotUserCheck, player.Controller.Id, botuser.Service.Id, j.GameRoomUuid, j.Players)
		}

	case botuser.ActionBotUserReady:
		botUserUuid := action.Key().From()
		secWebSocketKey := action.Payload()[0].(string)

		flux.Send(gameuser.ActionBotStart, player.Controller.Id, secWebSocketKey, botUserUuid)

	case botuser.ActionBotUserLeaveGameRoom:
		botUserUuid := action.Key().From()
		secWebSocketKey := action.Payload()[0].(string)
		gameRoomUuid := action.Payload()[1].(string)

		flux.Send(gameuser.ActionBotStop, player.Controller.Id, secWebSocketKey, "")

		flux.Send(gameroom.ActionRoomLeaveCall, player.Controller.Id, gameroom.Service.Id, gameroom.NewLeaveGameRoom(
			player.Controller.Id, botUserUuid, gameRoomUuid, true,
		))

	case gameuser.ActionLeaveGameRoomCall:
		secWebSocketKey := action.Key().From()
		gameRoomUuid := action.Payload()[0].(string)
		isDisconnect := action.Payload()[1].(bool)

		flux.Send(gameroom.ActionRoomLeaveCall, player.Controller.Id, gameroom.Service.Id, gameroom.NewLeaveGameRoom(
			player.Controller.Id, secWebSocketKey, gameRoomUuid, isDisconnect,
		))

	case gameroom.ActionRoomLeaveRecall:
		l := action.Payload()[0].(*gameroom.LeaveGameRoom)

		mediator.Service.Delete(l.SecWebSocketKey)

		// bot user
		if len(l.SecWebSocketKey) == len(uuid.New().String()) {
			flux.Send(lottery.ActionLotteryStop, player.Controller.Id, lottery.Service.Id, l.SecWebSocketKey)
		}

		// real game user
		if len(l.SecWebSocketKey) < len(uuid.New().String()) {
			flux.Send(gameuser.ActionLeaveGameRoomRecall, player.Controller.Id, l.SecWebSocketKey, l.GameRoomUuid)
		}

		flux.Send(botuser.ActionBotUserCheck, player.Controller.Id, botuser.Service.Id, l.GameRoomUuid, l.Players)

	case refund.ActionRefundRecall:
		// game user
		secWebSocketKey := action.Payload()[0].(string)
		isDisconnect := action.Payload()[1].(bool)

		if isDisconnect {
			flux.Send(gameuser.ActionGameUserDelete, player.Controller.Id, gameuser.Service.Id, secWebSocketKey)
		} else {
			flux.Send(gameuser.ActionGameUserLeave, player.Controller.Id, gameuser.Service.Id, secWebSocketKey)
		}

	case gameuser.ActionLeaveLobbyRoomCall:
		secWebSocketKey := action.Key().From()
		lobbyRoomUuid := action.Payload()[0].(string)
		isDisconnect := action.Payload()[1].(bool)

		flux.Send(lobbyroom.ActionLobbyRoomLeaveCall, player.Controller.Id, lobbyroom.Service.Id, lobbyroom.NewLeaveLobbyRoom(
			player.Controller.Id, secWebSocketKey, lobbyRoomUuid, isDisconnect,
		))

	case lobbyroom.ActionLobbyRoomLeaveRecall:
		l := action.Payload()[0].(*lobbyroom.LeaveLobbyRoom)
		flux.Send(gameuser.ActionLeaveLobbyRoomRecall, player.Controller.Id, l.SecWebSocketKey, l.LobbyRoomUuid)

		if l.IsDisconnect {
			flux.Send(gameuser.ActionGameUserDelete, player.Controller.Id, gameuser.Service.Id, l.SecWebSocketKey)
		}
	}
}

func (cf *controllerFish) checkJoinBot(gameId string, subgameId int) bool {
	switch {
	case gameId == models.PSF_ON_00003:
		fallthrough
	case gameId == models.PSF_ON_00004 && subgameId == 2:
		fallthrough
	case gameId == models.PSF_ON_00005 && subgameId == 2:
		fallthrough
	case gameId == models.PSF_ON_00006 && subgameId == 2:
		fallthrough
	case gameId == models.PSF_ON_00007 && subgameId == 2:
		return false
	default:
		return true
	}
}
func (cf *controllerFish) ServeHandler(h player.IHandler) {
	logger.Service.Zap.Info(fmt.Sprintf("  /$$$$$$$   /$$                      /$$$$$$    /$$                           "))
	logger.Service.Zap.Info(fmt.Sprintf(" | $$__  $$ | $$                     /$$__  $$  | $$                           "))
	logger.Service.Zap.Info(fmt.Sprintf(" | $$  \\ $$| $$  /$$$$$$  /$$   /$$| $$  \\__//$$$$$$    /$$$$$$   /$$$$$$    "))
	logger.Service.Zap.Info(fmt.Sprintf(" | $$$$$$$/ | $$ |____  $$| $$  | $$|  $$$$$$ |_  $$_/   |____  $$ /$$__  $$   "))
	logger.Service.Zap.Info(fmt.Sprintf(" | $$____/  | $$  /$$$$$$$| $$  | $$ \\____  $$ | $$      /$$$$$$$| $$ \\__/   "))
	logger.Service.Zap.Info(fmt.Sprintf(" | $$       | $$ /$$__  $$| $$  | $$ /$$  \\ $$ | $$ /$$ /$$__  $$| $$         "))
	logger.Service.Zap.Info(fmt.Sprintf(" | $$       | $$|  $$$$$$$|  $$$$$$$|  $$$$$$/  |  $$$$/|  $$$$$$$| $$         "))
	logger.Service.Zap.Info(fmt.Sprintf(" |__/       |__/ \\_______/ \\____$$ \\______/   \\___/ \\_______/|__/         "))
	logger.Service.Zap.Info(fmt.Sprintf("                           /$$  | $$  Production Https Server Listening %s   ", cf.port))
	logger.Service.Zap.Info(fmt.Sprintf("                          |  $$$$$$/          %s             ", cf.startTime))
	logger.Service.Zap.Info(fmt.Sprintf("                          \\______/           %s          ", cf.version))

	vip.Service.GameType = "FISH"
	//
	//var certFile, keyFile string
	//
	//certFile = filepath.Join(cf.sslPath, cf.sslCrt)
	//
	//keyFile = filepath.Join(cf.sslPath, cf.sslKey)

	http.HandleFunc(fmt.Sprintf("/%s", cf.endPoint), h.HttpHandler)
	if err := http.ListenAndServe(":"+cf.port, nil); err != nil {
		log.Fatal(err)
	}
	//if err := http.ListenAndServeTLS(":"+cf.port, certFile, keyFile, nil); err != nil {
	//	log.Fatal(err)
	//}

	//if err := fasthttp.ListenAndServeTLS(":"+cf.port, certFile, keyFile, h.FastHttpHandler); err != nil {
	//    log.Fatal(err)
	//}
}

func (cf *controllerFish) FastHttpHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case fmt.Sprintf("/%s", cf.endPoint):
		ss := strings.Split(
			string(ctx.Request.Header.Peek("Sec-WebSocket-Protocol")),
			",",
		)

		// fish
		gameType := ss[0]

		secWebSocketKey := string(ctx.Request.Header.Peek("Sec-WebSocket-Key"))
		hostExtId := strings.TrimSpace(ss[1])
		remoteAddr := ctx.RemoteAddr().String()
		userAgent := string(ctx.UserAgent())

		responseHeader := ctx.Response.Header
		responseHeader.Set("Sec-WebSocket-Protocol", gameType)

		//handler := func(conn *websocket.Conn) {
		//	logger.Service.Zap.Infow("----------------Sec-WebSocket-Protocol--------------",
		//		"GameUser", secWebSocketKey,
		//		"GameType", gameType,
		//		"HostExtId", hostExtId,
		//		"UserAgent", userAgent,
		//		"RemoteAddr", remoteAddr,
		//	)
		//	conn.EnableWriteCompression(true)
		//	conn.SetCompressionLevel(9)
		//	//cf.exec(secWebSocketKey, hostExtId, remoteAddr, userAgent, conn)
		//}
		err := cf.fastHttpUpgrader.Upgrade(ctx, func(conn *fastHttpWebsocket.Conn) {
			logger.Service.Zap.Infow("----------------Sec-WebSocket-Protocol--------------",
				"GameUser", secWebSocketKey,
				"GameType", gameType,
				"HostExtId", hostExtId,
				"UserAgent", userAgent,
				"RemoteAddr", remoteAddr,
			)
			conn.EnableWriteCompression(true)
			conn.SetCompressionLevel(9)
			//cf.exec(secWebSocketKey, hostExtId, remoteAddr, userAgent, conn)
		})

		if err != nil {
			logger.Service.Zap.Warnw("Websocket Connection Failed",
				"GameUser", ctx.Request.Header.Peek("Sec-WebSocket-Key"),
				"HostExtId", hostExtId,
				"UserAgent", ctx.UserAgent(),
				"RemoteAddr", ctx.RemoteAddr().String(),
				"Error", err,
			)
		}

	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func (cf *controllerFish) HttpHandler(res http.ResponseWriter, req *http.Request) {
	ss := strings.Split(req.Header.Get("Sec-WebSocket-Protocol"), ",")

	// fish
	gameType := ss[0]

	secWebSocketKey := req.Header.Get("Sec-WebSocket-Key")
	hostExtId := strings.TrimSpace(ss[1])
	remoteAddr := req.RemoteAddr
	userAgent := req.UserAgent()

	responseHeader := http.Header{}
	responseHeader.Set("Sec-WebSocket-Protocol", gameType)

	conn, err := cf.gorillaUpgrader.Upgrade(res, req, responseHeader)

	if err != nil {
		logger.Service.Zap.Warnw("Websocket Connection Failed",
			"GameUser", secWebSocketKey,
			"GameType", gameType,
			"HostExtId", hostExtId,
			"UserAgent", req.UserAgent(),
			"RemoteAddr", req.RemoteAddr,
			"Error", err,
		)

		if conn != nil {
			conn.Close()
		}
		return
	}

	logger.Service.Zap.Infow("----------------Sec-WebSocket-Protocol--------------",
		"GameUser", secWebSocketKey,
		"GameType", gameType,
		"HostExtId", hostExtId,
		"UserAgent", userAgent,
		"RemoteAddr", remoteAddr,
	)

	conn.EnableWriteCompression(true)
	conn.SetCompressionLevel(-2)
	cf.exec(secWebSocketKey, hostExtId, remoteAddr, userAgent, conn)
}
