package room

import (
	"encoding/json"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log/slog"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"strconv"
)

func ChangeState(c *websocket.Conn, r *Room) {
	content := &Content{C: "changeState"}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	if r.Timer == 0 {
		content.P = P{
			Code:       comm.Succ,
			RoundId:    r.RoundId,
			TimeLeft:   Duration,
			NewStateID: 1,
		}

	} else {
		content.P = P{
			Code:       comm.Succ,
			RoundId:    r.RoundId,
			NewStateID: 2,
		}
	}
	//marshal, err := json.Marshal(body)
	//if err != nil {
	//	slog.Error("json.Marshal", "Err", err.Error())
	//}
	marshal, _ := GetChangeStateRsp(&body).ToBinary()
	c.WriteMessage(websocket.BinaryMessage, marshal)
}
func GetChangeStateRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("newStateId", int32(initRsp.Content.P.NewStateID))
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutInt("roundId", int32(initRsp.Content.P.RoundId))
	if initRsp.Content.P.NewStateID == 1 {
		pp.PutInt("timeLeft", int32(initRsp.Content.P.TimeLeft))
	}
	p.PutSFSObject("p", pp)
	p.PutString("c", "changeState")
	so.AddCreatePAC(p, 1, 13)
	return so
}

func UpdateCurrentBetsInit(c *websocket.Conn, r *Room) { // todo
	//content := `{"id":13,"targetController":1,"content":{"c":"updateCurrentBets","p":{"betsCount":1,"code":200}}}`

	//自己做一个
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", 200)
	pp.PutInt("betsCount", 1)

	p.PutSFSObject("p", pp)
	p.PutString("c", "updateCurrentBets")
	so.AddCreatePAC(p, 1, 13)
	binary, _ := so.ToBinary()

	c.WriteMessage(websocket.BinaryMessage, binary)
}

func UpdateCurrentBets(c *websocket.Conn, r *Room) {
	//content := `{"id":13,"targetController":1,"content":{"c":"updateCurrentBets","p":{"betsCount":1,"code":200}}}`
	content := &Content{
		C: "updateCurrentBets",
	}
	var err error
	content.P.Bets, err = GetSortedBets(redisx.GetClient(), r.Name)
	if err != nil {
		slog.Error("GetSortedBets Err", "err", err)
	}
	content.P = P{
		Bets:      content.P.Bets,
		BetsCount: len(content.P.Bets),
		Code:      comm.Succ,
	}
	if content.P.BetsCount >= 50 {
		content.P.Bets = content.P.Bets[:50]
	}
	rsp := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, _ := json.Marshal(rsp)
	marshal, _ := GetUpdateCurrentBetsRsp(&rsp).ToBinary()
	//c.WriteMessage(websocket.TextMessage, []byte(`UpdateCurrentBets我返回了哈`))
	c.WriteMessage(websocket.BinaryMessage, marshal)
}
func GetUpdateCurrentBetsRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutInt("betsCount", int32(initRsp.Content.P.BetsCount))
	bets := ut.NewSFSArray()
	for _, bet := range initRsp.Content.P.Bets {
		temp := ut.NewSFSObject()
		temp.PutDouble("bet", bet.Bet)
		temp.PutString("player_id", bet.PlayerID)
		temp.PutInt("betId", int32(bet.BetID))
		temp.PutBool("isFreeBet", false)
		temp.PutString("currency", bet.Currency)
		temp.PutString("profileImage", bet.ProfileImage)
		temp.PutString("username", bet.Username)
		bets.Add(temp, ut.SFS_OBJECT, true)
	}
	pp.PutSFSArray("bets", bets)
	p.PutSFSObject("p", pp)
	p.PutString("c", "updateCurrentBets")

	so.AddCreatePAC(p, 1, 13)
	return so
}

// 投注时提现时和游戏结束时会发送
func NewBalance(c *websocket.Conn, r *Room) {
	session, err := comm.GetSession(c)
	if err != nil {
		slog.Error("ClientLikeChat Err", "GetSession err", err)
		rsp := CommBody{
			Id:               13,
			TargetController: 1,
		}
		rsp.Content = &Content{
			P: P{
				Code:         comm.ErrNoPlr,
				OperatorKey:  "release",
				ErrorMessage: "NOT Found This Plr",
			},
		}
		rspMarshal, _ := json.Marshal(rsp)
		c.WriteMessage(websocket.BinaryMessage, rspMarshal)
	}
	playerIdStr := session.PlayerId
	playerId, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		slog.Error("ClientLogin Err", "strconv.ParseInt err", err)
	}
	plr, err := db.GetDocPlayer(playerId)
	if err != nil {
		slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
	}
	player := r.Players[playerIdStr]
	if r.State == ROOMSTATE_SETTLE && player.PlayerState != PLAYERSTATE_BET { //在结算阶段 且 玩家的状态==ob时就不再发送余额了
		return
	}
	balance, err := slotsmongo.GetBalance(plr.PID)
	if err != nil {
		return
	}
	content := &Content{C: "newBalance",
		P: P{
			Code:       comm.Succ,
			NewBalance: ut.Gold2Money(balance),
		},
	}
	rsp := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, err := json.Marshal(rsp)
	//if err != nil {
	//}
	//content := `{"id":13,"targetController":1,"content":{"c":"newBalance","p":{"code":200,"newBalance":80006.23}}}`
	//c.WriteMessage(websocket.TextMessage, []byte(`NewBalance我返回了哈`))
	marshal, _ := GetNewBalanceRsp(&rsp).ToBinary()
	c.WriteMessage(websocket.BinaryMessage, marshal)
}
func GetNewBalanceRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutDouble("newBalance", initRsp.Content.P.NewBalance)

	p.PutSFSObject("p", pp)
	p.PutString("c", "newBalance")
	so.AddCreatePAC(p, 1, 13)
	return so
}

func SendCoordinate(c *websocket.Conn, r *Room) {
	x := r.XArr[r.XIndex]
	content := &Content{C: "x", P: P{X: x, Code: comm.Succ}}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, err := json.Marshal(body)
	//if err != nil {
	//	slog.Error("json.Marshal", "Err", err.Error())
	//}
	//c.WriteMessage(websocket.TextMessage, []byte(`SendCoordinate我返回了哈`))
	//contentStr := `{"id":13,"targetController":1,"content":{"c":"x","p":{"x":1,"code":200}}}`
	//c.WriteMessage(websocket.TextMessage, []byte(contentStr))
	marshal, _ := GetSendCoordinateRsp(&body).ToBinary()
	c.WriteMessage(websocket.BinaryMessage, marshal)
}
func GetSendCoordinateRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutDouble("x", initRsp.Content.P.X)
	if initRsp.Content.P.CrashX > 0 {
		pp.PutDouble("crashX", initRsp.Content.P.CrashX)
	}

	p.PutSFSObject("p", pp)
	p.PutString("c", "x")
	so.AddCreatePAC(p, 1, 13)
	return so
}

func UpdateCurrentCashOuts(c *websocket.Conn, r *Room) { //其他玩家体现通知
	//content := `{"id":13,"targetController":1,"content":{"c":"updateCurrentCashOuts","p":{"cashouts":[{"betId":1,"multiplier":1.05,"player_id":"162","winAmount":105.78}],"code":200}}}`
	content := &Content{
		C: "updateCurrentCashOuts",
	}
	var err error
	content.P.Cashouts, err = GetSortedCashOuts(redisx.GetClient(), r.Name)
	if err != nil {
		slog.Error("GetSortedBets Err", "err", err)
	}
	content.P = P{
		Cashouts: content.P.Cashouts,
		Code:     comm.Succ,
	}
	if len(content.P.Cashouts) >= 50 {
		content.P.Cashouts = content.P.Cashouts[:50]
	}
	rsp := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, _ := json.Marshal(rsp)
	//c.WriteMessage(websocket.TextMessage, []byte(`UpdateCurrentBets我返回了哈`))
	marshal, _ := GetUpdateCurrentCashOutsRsp(&rsp).ToBinary()
	c.WriteMessage(websocket.BinaryMessage, marshal)
}
func GetUpdateCurrentCashOutsRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	cashouts := ut.NewSFSArray()
	for _, cashout := range initRsp.Content.P.Cashouts {
		temp := ut.NewSFSObject()
		temp.PutString("player_id", cashout.PlayerID)
		temp.PutDouble("winAmount", cashout.WinAmount)
		temp.PutDouble("multiplier", cashout.Multiplier)
		temp.PutInt("betId", int32(cashout.BetID))
		cashouts.Add(temp, ut.SFS_OBJECT, true)
	}
	pp.PutSFSArray("cashouts", cashouts)

	p.PutSFSObject("p", pp)
	p.PutString("c", "updateCurrentCashOuts")
	so.AddCreatePAC(p, 1, 13)
	return so
}

func SendCrashCoordinate(c *websocket.Conn, r *Room) {
	x := r.XArr[len(r.XArr)-1]
	content := &Content{C: "x", P: P{X: x, CrashX: x, Code: comm.Succ, Crash: true}}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, err := json.Marshal(body)
	//if err != nil {
	//	slog.Error("json.Marshal", "Err", err.Error())
	//}
	//c.WriteMessage(websocket.TextMessage, []byte(`SendCrashCoordinate我返回了哈`))
	//contentStr := `{"id":13,"targetController":1,"content":{"c":"x","p":{"x":1.23,"crashX":1.23,"code":200,"crash":true}}}`
	//c.WriteMessage(websocket.TextMessage, []byte(contentStr))
	marshal, _ := GetSendCoordinateRsp(&body).ToBinary()
	c.WriteMessage(websocket.BinaryMessage, marshal)
}

func RoundChartInfo(c *websocket.Conn, r *Room) {
	content := &Content{C: "roundChartInfo", P: P{Code: comm.Succ, MaxMultiplier: r.XArr[len(r.XArr)-1], RoundId: r.RoundId}}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, err := json.Marshal(body)
	//if err != nil {
	//	slog.Error("json.Marshal", "Err", err.Error())
	//}
	//c.WriteMessage(websocket.TextMessage, []byte(`RoundChartInfo我返回了哈`))
	//contentStr := `{"id":13,"targetController":1,"content":{"c":"roundChartInfo","p":{"code":200,"maxMultiplier":1.23,"roundId":1320758}}}`
	//c.WriteMessage(websocket.TextMessage, []byte(contentStr))
	marshal, _ := GetRoundChartInfoRsp(&body).ToBinary()
	c.WriteMessage(websocket.BinaryMessage, marshal)
}
func GetRoundChartInfoRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutDouble("maxMultiplier", initRsp.Content.P.MaxMultiplier)
	pp.PutInt("roundId", int32(initRsp.Content.P.RoundId))

	p.PutSFSObject("p", pp)
	p.PutString("c", "roundChartInfo")
	so.AddCreatePAC(p, 1, 13)
	return so
}

func RemovePlayer(c *websocket.Conn, r *Room) {
	session, err := comm.GetSession(c)
	if err != nil {
		slog.Error("RemovePlayer Err", "GetSession err", err)
		return
	}
	delete(r.Players, session.PlayerId)
}

func OnlinePlayers(c *websocket.Conn, r *Room) {
	content := &Content{C: "onlinePlayers", P: P{Code: comm.Succ, OnlinePlayers: len(r.Players) + 4923}}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, err := json.Marshal(body)
	//if err != nil {
	//	slog.Error("json.Marshal", "Err", err.Error())
	//}
	marshal, _ := GetOnlinePlayersRsp(&body).ToBinary()
	c.WriteMessage(websocket.BinaryMessage, marshal)
}
func GetOnlinePlayersRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutLong("onlinePlayers", int64(initRsp.Content.P.OnlinePlayers))
	p.PutSFSObject("p", pp)
	p.PutString("c", "onlinePlayers")

	so.AddCreatePAC(p, 1, 13)
	return so
}

//`{
//  "id": 13,
//  "targetController": 1,
//  "content": {
//    "c": "onlinePlayers",
//    "p": {
//      "code": 200,
//      "onlinePlayers": 5037
//    }
//  }
//}`
