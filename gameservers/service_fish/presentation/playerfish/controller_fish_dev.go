//go:build dev
// +build dev

package playerfish

import (
	"fmt"
	"log"
	"net/http"
	"serve/fish_comm/flux/logger"
	"serve/fish_comm/flux/player"
	"serve/fish_comm/flux/vip"

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
	logger.Service.Zap.Info(fmt.Sprintf("                           /$$  | $$  Development Http Server Listening %s   ", cf.port))
	logger.Service.Zap.Info(fmt.Sprintf("                          |  $$$$$$/          %s             ", cf.startTime))
	logger.Service.Zap.Info(fmt.Sprintf("                          \\______/           %s          ", cf.version))

	vip.Service.GameType = "FISH"

	http.HandleFunc(fmt.Sprintf("/%s", cf.endPoint), h.HttpHandler)

	if err := http.ListenAndServe(":"+cf.port, nil); err != nil {
		log.Fatal(err)
	}

	//if err := fasthttp.ListenAndServe(":"+c89.port, h.FastHttpHandler); err != nil {
	//	log.Fatal(err)
	//}
}

func (cf *controllerFish) FastHttpHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case fmt.Sprintf("/%s", cf.endPoint):
		secWebSocketKey := string(ctx.Request.Header.Peek("Sec-WebSocket-Key"))
		hostExtId := "dd12f3635bc9599ca58012601d07197c"
		remoteAddr := ctx.RemoteAddr().String()
		userAgent := string(ctx.UserAgent())

		err := cf.fastHttpUpgrader.Upgrade(ctx, func(conn *websocket.Conn) {
			logger.Service.Zap.Infow("----------------Sec-WebSocket-Protocol--------------",
				"GameUser", secWebSocketKey,
				"HostExtId", hostExtId,
				"UserAgent", userAgent,
				"RemoteAddr", remoteAddr,
			)
			conn.EnableWriteCompression(true)
			conn.SetCompressionLevel(9)
			//c89.exec(secWebSocketKey, hostExtId, remoteAddr, userAgent, conn)
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
	secWebSocketKey := req.Header.Get("Sec-WebSocket-Key")
	hostExtId := "dd12f3635bc9599ca58012601d07197c"
	remoteAddr := req.RemoteAddr
	userAgent := req.UserAgent()

	conn, err := cf.gorillaUpgrader.Upgrade(res, req, nil)

	if err != nil {
		logger.Service.Zap.Warnw("websocket connection failed",
			"GameUser", secWebSocketKey,
			"HostExtId", hostExtId,
			"UserAgent", userAgent,
			"RemoteAddr", remoteAddr,
			"Error", err,
		)

		if conn != nil {
			conn.Close()
		}
		return
	}

	logger.Service.Zap.Infow("----------------Sec-WebSocket-Protocol--------------",
		"GameUser", secWebSocketKey,
		"HostExtId", hostExtId,
		"UserAgent", userAgent,
		"RemoteAddr", remoteAddr,
	)

	conn.EnableWriteCompression(true)
	conn.SetCompressionLevel(-2)
	cf.exec(secWebSocketKey, hostExtId, remoteAddr, userAgent, conn)
}
