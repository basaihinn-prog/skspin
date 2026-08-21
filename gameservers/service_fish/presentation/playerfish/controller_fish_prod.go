//go:build prod
// +build prod

package playerfish

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"serve/fish_comm/flux/logger"
	"serve/fish_comm/flux/player"
	_ "serve/fish_comm/flux/refundbet-tool"
	"serve/fish_comm/flux/vip"
	"strings"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

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

	var certFile, keyFile string

	certFile = filepath.Join(cf.sslPath, cf.sslCrt)

	keyFile = filepath.Join(cf.sslPath, cf.sslKey)

	http.HandleFunc(fmt.Sprintf("/%s", cf.endPoint), h.HttpHandler)

	if err := http.ListenAndServeTLS(":"+cf.port, certFile, keyFile, nil); err != nil {
		log.Fatal(err)
	}

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

		err := cf.fastHttpUpgrader.Upgrade(ctx, func(conn *websocket.Conn) {
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
