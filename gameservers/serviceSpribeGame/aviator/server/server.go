package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/jwtutil"
	"serve/comm/lazy"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"serve/serviceSpribeGame/aviator/room"
	"strconv"
	"sync"
	"time"

	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

type Server struct {
	sync.RWMutex
	engine   *nbhttp.Engine
	upgrader *websocket.Upgrader
	rooms    map[string]*room.Room
}

func NewServer() *Server {
	return &Server{rooms: make(map[string]*room.Room)}
}

func (s *Server) Start() error {
	err := s.CheckRoomAndNew()
	if err != nil {
		slog.Error("CheckRoomAndNew err:", err)
	}
	//{
	//	newRoom := room.NewRoom("room1")
	//	newRoom.Run()
	//	s.Lock()
	//	if _, ok := s.rooms["room1"]; !ok {
	//		s.rooms["room1"] = newRoom
	//	}
	//	s.Unlock()
	//}
	s.upgrader = s.newUpgrader()
	mux := &http.ServeMux{}
	mux.HandleFunc("/BlueBox/websocket", s.onWebsocket)
	mux.HandleFunc("/aviator/playerSession", wrapWebApi(aviatorPlayerSession, true))
	////证书
	//certfile, _ := lazy.RouteFile.Get("certfile")
	//keyfile, _ := lazy.RouteFile.Get("keyfile")
	//// 配置 TLS
	//cert, err := tls.LoadX509KeyPair(certfile, keyfile)
	//if err != nil {
	//	log.Fatalf("Failed to load certificate: %v", err)
	//}
	//tlsConfig := &tls.Config{
	//	Certificates:       []tls.Certificate{cert},
	//	InsecureSkipVerify: true,
	//}
	addr := lo.Must(lazy.RouteFile.Get("minigateway.http.api"))
	s.engine = nbhttp.NewEngine(nbhttp.Config{
		Network: "tcp",
		//TLSConfig: tlsConfig,
		//AddrsTLS:  []string{addr},
		Addrs:   []string{addr},
		Handler: mux,
		//ReleaseWebsocketPayload: true,
	})
	s.startTimer()
	return s.engine.Start()
}
func (s *Server) Stop() {
	s.engine.Stop()
}

func (s *Server) newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	//newBetsInfo := room.NewBetsInfo("room1")
	u.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		defer func() {
			c.SetReadDeadline(time.Now().Add(nbhttp.DefaultKeepaliveTime))
		}()
		var action string
		var p *ut.SFSObject
		nowRoom := s.rooms["room1"]
		// data解密
		binaryData, err := ut.NewFromBinaryData2(data)
		if err != nil {
			slog.Error("NewFromBinaryData err:", err)
		}
		if binaryData == nil { //没有二进制转换成功就是明文形式
			slog.Info(gjson.GetBytes(data, "channel").String())
			action = gjson.GetBytes(data, "channel").String()
		} else {
			a, _ := binaryData.GetShort("a")
			p, _ = binaryData.GetSFSObject("p")
			if a == 13 {
				pc, _ := p.GetString("c")
				action = room.ActionMap[pc]
			} else {
				action = room.ActionMap[int(a)]
			}
		}
		//登陆特殊处理
		if action == "AviatorLoginIdReq" {
			//req := room.ClientLoginReq{}
			//err := json.Unmarshal(data, &req)
			//if err != nil {
			//	slog.Error("ClientLogin Err", "json.Unmarshal", err)
			//}
			//player, err := comm.TokenGetPlr(req.Content.Cl)
			tk, _ := p.GetString("cl")
			player, err := comm.TokenGetPlr(tk)
			if err != nil {
				slog.Error("newUpgrader::comm.TokenGetPlr Err", "err", err)
			}
			if _, ok := s.rooms[player.AppID]; ok {
				nowRoom = s.rooms[player.AppID]
			} else {
				slog.Error("AviatorLoginIdReq Err: NOT Found This Room")
				rsp := room.CommBody{
					Id:               13,
					TargetController: 1,
				}
				rsp.Content = &room.Content{
					P: room.P{
						Code:         comm.ErrNoRoom,
						OperatorKey:  "release",
						ErrorMessage: "NOT Found This Room",
					},
				}
				//rspMarshal, _ := json.Marshal(rsp)
				rspMarshal, _ := room.GetClientErrRsp(&rsp).ToBinary()
				c.WriteMessage(messageType, rspMarshal)
			}
		} else {
			session, err := comm.GetSession(c)
			if err != nil {
				slog.Error("newUpgrader Err", "GetSession err", err)
				rsp := room.CommBody{
					Id:               13,
					TargetController: 1,
				}
				rsp.Content = &room.Content{
					P: room.P{
						Code:         comm.ErrNoPlr,
						OperatorKey:  "release",
						ErrorMessage: "NOT Found This Plr",
					},
				}
				rspMarshal, _ := json.Marshal(rsp)
				c.WriteMessage(messageType, rspMarshal)
			}
			appId := session.Plr.AppID
			nowRoom = s.rooms[appId]
		}
		switch action {
		case "AviatorLoginIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorLoginIdReq, Conn: c, Msg: p})

		case "AviatorGameInfoIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorGameInfoIdReq, Conn: c, Msg: p})

		case "AviatorCurrentBetsInfoIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorCurrentBetsInfoIdReq, Conn: c, Msg: p})

		case "AviatorBetIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorBetIdReq, Conn: c, Msg: p})

		case "AviatorCancelBetIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorCancelBetIdReq, Conn: c, Msg: p})

		case "AviatorCashOutIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorCashOutIdReq, Conn: c, Msg: p})

		case "AviatorBetHistoryIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorBetHistoryIdReq, Conn: c, Msg: p})

		case "AviatorGameStatePingIdReq":
			so := ut.NewSFSObject()
			p := ut.NewSFSObject()
			pp := ut.NewSFSObject()
			p.PutSFSObject("p", pp)
			p.PutString("c", "PING_RESPONSE")
			so.AddCreatePAC(p, 1, 13)
			marshal, _ := so.ToBinary()
			c.WriteMessage(messageType, marshal)

		case "AviatorAddChatMessageIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorAddChatMessageIdReq, Conn: c, Msg: p})

		case "AviatorLikeIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorLikeIdReq, Conn: c, Msg: p})

		case "ClientSearchGifs":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.ClientSearchGifs, Conn: c, Msg: p})

		case "AviatorPreviousRoundInfoIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorPreviousRoundInfoIdReq, Conn: c, Msg: p})

		case "AviatorGetHugeWinsInfoIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorGetHugeWinsInfoIdReq, Conn: c, Msg: p})

		case "AviatorGetTopWinsInfoIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorGetTopWinsInfoIdReq, Conn: c, Msg: p})

		case "AviatorGetTopRoundsInfoIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorGetTopRoundsInfoIdReq, Conn: c, Msg: p})

		case "AviatorChangeProfileImageIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorChangeProfileImageIdReq, Conn: c, Msg: p})

		case "AviatorRoundFairnessIdReq":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.AviatorRoundFairnessIdReq, Conn: c, Msg: p})

		case "ServerSeedHandler":
			nowRoom.AddMsg(&comm.MSG{Typ: comm.ServerSeedHandler, Conn: c, Msg: p})
		default:
			nowRoom.AddMsg(&comm.MSG{Typ: comm.NETMSG, Conn: c, Msg: p})
		}
	})

	u.OnOpen(func(c *websocket.Conn) { //连接升级成功
		slog.Info("OnOpen")
		//plr, err := comm.GetPlr(c)
		//if err != nil {
		//	slog.Error("newUpgrader Err", "db.GetDocPlayer err", err)
		//}
		//nowRoom := s.rooms["room1"]
		//if _, ok := s.rooms[plr.AppID]; ok {
		//	nowRoom = s.rooms[plr.AppID]
		//}
		//nowRoom.AddMsg(&comm.MSG{Typ: comm.NETOPEN, Conn: c})
	})

	u.OnClose(func(c *websocket.Conn, err error) { //已断开连接
		slog.Info("OnClose")
		session, err := comm.GetSession(c)
		if err != nil {
			slog.Error("newUpgrader Err", "db.GetDocPlayer err", err)
		}
		plr := session.Plr
		nowRoom := s.rooms["room1"]
		if _, ok := s.rooms[plr.AppID]; ok {
			nowRoom = s.rooms[plr.AppID]
		}
		nowRoom.AddMsg(&comm.MSG{Typ: comm.NETCLOSE, Conn: c})
	})

	return u
}

func (s *Server) onWebsocket(w http.ResponseWriter, r *http.Request) {
	slog.Info("成功收到wss")
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	conn.SetSession(r.URL.Query())
	log.Println("OnOpen:", conn.RemoteAddr().String())
}

func (s *Server) startTimer() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				os.Stdout.Write(debug.Stack())
			}
		}()
		time.Sleep(3 * time.Second)
		slog.Info("服务启动，延迟三秒后计时器开始")
		ticker1 := time.NewTicker(1 * time.Second)        // 每 1 秒触发一次
		ticker2 := time.NewTicker(100 * time.Millisecond) // 每 100 毫秒触发一次
		ticker3 := time.NewTicker(5 * time.Second)
		defer ticker1.Stop()
		defer ticker2.Stop()
		defer ticker3.Stop()
		for {
			select {
			case <-ticker1.C:
				s.RLock()
				for _, v := range s.rooms { //分配给每个房间时间节点-每秒
					v.AddMsg(&comm.MSG{Typ: comm.TIMER_1s})
				}
				s.RUnlock()
			case <-ticker2.C:
				s.RLock()
				for _, v := range s.rooms {
					v.AddMsg(&comm.MSG{Typ: comm.TIMER_100})
				}
				s.RUnlock()
			case <-ticker3.C: //todo 转移到单独协程定时刷新
				err := s.CheckRoomAndNew()
				if err != nil {
					slog.Error("CheckRoomAndNew err:", err)
				}
			}
		}
	}()
}

func wrapWebApi[R any](fn func(*PGParams, *R) error, verify bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ret MiniRetWrapper
			err error
		)

		traceId := r.URL.Query().Get("traceId")
		defer func() {
			if r.URL.Path == "/back-office-proxy/Report/GetBetHistory" {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			if err != nil {
				ret.Err = &PGError{}
				if errors.As(err, &ret.Err) && err.(*PGError).Cd == "5000" {
					errors.As(err, &ret.Err)
				}
				if ec, ok := err.(define.IErrcode); ok {
					code := ec.Code()
					if code != 0 {
						ret.Err.Cd = strconv.Itoa(code)
						ret.Err.Msg = err.Error()
					}
				}
			}
			jsondata, _ := json.Marshal(ret)
			w.Write(jsondata)
		}()

		err = r.ParseForm()
		if err != nil {
			return
		}

		var pid int64

		tk := r.FormValue("otk")

		pid, err = jwtutil.ParseToken(tk)
		if err != nil {
			err = define.NewErrCode("Invalid player session", 1302)
			return
		}

		gi := r.FormValue("gi")

		ps := &PGParams{
			Path:    r.URL.Path,
			TraceId: traceId,
			Form:    r.Form,
			Pid:     pid,
			GameId:  "pg_" + gi,
		}

		var ans R
		err = fn(ps, &ans)
		if err != nil {
			return
		}

		resp, _ := json.Marshal(ans)
		if len(resp) != 0 {
			raw := json.RawMessage(resp)
			ret.Dt = &raw
		}
	}
}

// https://api.pg-demo.com/web-api/auth/session/v2/verifyOperatorPlayerSession?traceId=VXTNFQ12
func aviatorPlayerSession(ps *PGParams, ret *M) (err error) {
	os := ps.Form.Get("otk")
	tk := os
	s := M{
		"tk": tk,
	}
	*ret = s
	return
}

type PGParams = define.PGParams

type PGError struct {
	Cd  string `json:"cd"`
	Msg string `json:"msg"`
	Tid string `json:"tid"`
}

func (e *PGError) Error() string {
	return e.Msg
}

type MiniRetWrapper struct {
	Dt  *json.RawMessage `json:"dt"`
	Err *PGError         `json:"err"`
}

type M = map[string]any
type D = []any

func (s *Server) CheckRoomAndNew() error {
	// 运营商判断
	type operatorData struct {
		Status      int    `bson:"Status"`
		Name        string `bson:"Name"`
		AppID       string `bson:"AppID"`
		CurrencyKey string `bson:"CurrencyKey"`
	}
	appS := make([]operatorData, 0)
	operatorList := make([]map[string]string, 0)
	hhhh := options.Find().SetProjection(bson.M{"AppID": true})
	cursor1, err := db.Collection2("GameAdmin", "GameConfig").Find(context.TODO(), bson.M{
		"GameOn": 0,
		"GameId": "spribe_01",
	}, hhhh)
	if err != nil {
		return define.NewErrCode("Operator not exist.", 1310)
	}
	err = cursor1.All(context.TODO(), &operatorList)
	if err != nil {
		return err
	}
	appList := make([]string, 0)
	for _, appName := range operatorList {
		appList = append(appList, appName["AppID"])
	}
	//appList = []string{"faketrans", "qwe456"}
	//获取开启了飞机的商户 -> 获取这些商户中，商户状态是开启的商户
	cursor, err := db.Collection2("GameAdmin", "AdminOperator").Find(context.TODO(), db.D("Status", 1, "AppID", bson.M{"$in": appList}))
	if err != nil {
		return define.NewErrCode("Operator not exist.", 1310)
	}
	err = cursor.All(context.TODO(), &appS)
	if err != nil {
		return err
	}
	s.Lock()
	for _, app := range appS {
		if _, ok := s.rooms[app.AppID]; !ok {
			newRoom := room.NewRoom(app.AppID)
			newRoom.Currency = app.CurrencyKey
			s.rooms[app.AppID] = newRoom
			newRoom.Run()
		}
	}
	s.Unlock()
	return nil
}

func GetAppIds() ([]string, error) {
	var err error
	var res []string
	type operatorData struct {
		AppID string `bson:"AppID"`
	}
	appS := make([]operatorData, 0)
	cursor, err := db.Collection2("GameAdmin", "AdminOperator").Find(context.TODO(), db.D("Status", 1, "AppID", bson.M{"$in": []string{"faketrans", "qwe456"}}))
	if err != nil {
		return nil, define.NewErrCode("Operator not exist.", 1310)
	}
	err = cursor.All(context.TODO(), &appS)
	if err != nil {
		return nil, err
	}
	for _, app := range appS {
		res = append(res, app.AppID)
	}
	return res, nil
}
