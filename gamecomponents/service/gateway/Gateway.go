package main

import (
	"encoding/json"
	"game/comm/mq"
	"game/comm/ut"
	"game/duck/logger"
	"game/duck/ut2"
	"game/duck/ut2/httputil"
	"game/duck/ut2/jwtutil"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/samber/lo"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type Gateway struct {
	players ut2.IMap[int64, *Connection]
}

func (g *Gateway) Kick(pid int64, reason string) {

	conn, ok := g.players.Load(pid)
	if !ok {
		return
	}

	g.players.Delete(pid)

	go func() {
		conn.SendMsg(pushMsg("login_err", reason))
		time.Sleep(200 * time.Millisecond) // ensure data sended
		conn.Close()
	}()
}

func validateToken(r *http.Request) (pid int64, game, language string, err error) {
	l := slog.With("path", r.URL.Path,
		"query", r.URL.RawQuery,
		"host", r.Host,
	)
	defer func() {
		if err != nil {
			l.Error("validateTokenError", "error", ut.ErrString(err))
		} else {
			l.Info("validateTokenSuccess")
		}
	}()

	var arg struct {
		Token string `query:"t"`
		L     string `query:"l"`
	}

	if err = httputil.HttpBindQuery(&arg, r); err != nil {
		return
	}

	l = l.With("arg", arg)

	pid, game, err = jwtutil.ParseTokenData(arg.Token)
	if err != nil {
		return
	}

	language = lo.Ternary(arg.L != "", arg.L, "th")

	l = l.With("pid", pid, "game", game, "language", language)
	// _, err = client.ValidateToken(context.TODO(), &gamepb.TokenReq{Token: arg.Token})
	// if err != nil {
	// 	return 0, "", err
	// }fdfdsfdkkskkkkuuuqquuuiiiqqqzzzzzzzzz

	return
}

func updatePlrLoginTime(pid int64) (err error) {
	/*
		coll := db.Collection2("game", "Players")

		var (
			plr comm.Player
			now = time.Now()
		)
		err = coll.FindOneAndUpdate(context.TODO(), db.ID(pid), db.D("$set", db.D("LoginAt", mongodb.NewTimeStamp(now))), options.FindOneAndUpdate().SetReturnDocument(options.Before).SetProjection(db.D("LoginAt", 1, "AppID", 1))).Decode(&plr)

		if err != nil {
			return
		}

		if !ut.IsSameDate(plr.LoginAt.AsTime(), now) {
			slotsmongo.IncLoginCount(plr.AppID)
		}

	*/
	return nil
}

func pushMsg(key string, data any) []byte {
	return mq.MakeMsg(key, data)
	// return lazy.MakeRouteMsg(data, key)
}

func (g *Gateway) route_login(w http.ResponseWriter, r *http.Request) {
	ip := httputil.GetIPFromRequest(r)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Err("upgrade:", err)
		return
	}
	defer c.Close()

	pid, gameid, language, err := validateToken(r)
	// logger.Info("incoming", pid, gameid, err)
	// logger.Info(r.URL.RawQuery)
	if err != nil {
		c.WriteMessage(websocket.TextMessage, pushMsg("login_err", err))
		time.Sleep(100 * time.Millisecond)
		return
	}

	// updatePlrLoginTime(pid)

	g.Kick(pid, "@multi_login")

	conn := NewConnection(c, ip, g)
	conn.pid = pid
	conn.gameid = gameid
	conn.language = language
	g.players.Store(pid, conn)

	// TODO: online msg

	c.WriteMessage(websocket.TextMessage, pushMsg("login_ok", nil))

	conn.ServeWs()
}

// recv from game
func (g *Gateway) OnRecv(topic string, pids []int64, message []byte) {

	for _, pid := range pids {
		player, ok := g.players.Load(pid)
		if ok {
			player.SendMsg(message)
		}
	}
}

type RawBase struct {
	Svr string `json:"svr"`
	// Route string          `json:"route"`
	CBID int `json:"cbId"`
	// Args  json.RawMessage `json:"args"`
}

// recv from player
func (g *Gateway) OnRecvFrom(conn *Connection, message []byte) {
	logger.Info("收到消息", conn.gameid, string(message))
	// "{\"subj\":\"Roma\",\"route\":\"Roma.Gold\",\"args\":null,\"cbId\":1}"
	subj := conn.gameid

	var top RawBase
	json.Unmarshal(message, &top)
	if top.Svr == "gamecenter" {
		subj = top.Svr
	}

	mq.Publish(subj, &mq.Forward{
		Pids: []int64{conn.pid},
		Data: message,
		L:    conn.language,
	})
}

func (g *Gateway) OnClose(conn *Connection) {
	store, ok := g.players.Load(conn.pid)
	if !ok {
		return
	}

	if store == conn {
		g.players.Delete(conn.pid)
		// TODO: offline msg
	}
}
