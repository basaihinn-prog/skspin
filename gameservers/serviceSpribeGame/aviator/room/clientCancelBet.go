package room

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"serve/comm/define"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"strconv"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

func ClientCancelBet(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) { //todo
	//req := CommBody{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientCancelBet Err", "json.Unmarshal", err)
	//}
	p, _ := data.GetSFSObject("p")
	reqBetId, _ := p.GetDouble("betId")
	session, err := comm.GetSession(c)
	if err != nil {
		slog.Error("ClientCancelBet Err", "GetSession err", err)
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
	plr := session.Plr
	playerIdStr := strconv.Itoa(int(plr.PID))
	if r.State != ROOMSTATE_WAIT_BET {
		slog.Error("ClientCancelBet err::Room State Not Bet")
		rsp := CommBody{
			Id:               13,
			TargetController: 1,
		}

		rsp.Content = &Content{
			C: "cancelBet",
			P: P{
				Code:         comm.ErrTimeout,
				PlayerID:     playerIdStr,
				ErrorMessage: "Stage time out",
				BetID:        int(reqBetId),
			},
		}
		//cashOutMarshal, _ := json.Marshal(rsp)
		cashOutMarshal, _ := GetClientErrRsp(&rsp).ToBinary()
		c.WriteMessage(messageType, cashOutMarshal)
		return
	}
	player := r.Players[session.PlayerId]
	//if player.PlayerState != PLAYERSTATE_BET { //校验用户是否为Bet状态，如果不是Bet状态则不应该取消投注
	//	slog.Error("ClientCancelBet::PlayerState Err", "player.PlayerState != ", "OB")
	//	return
	//}
	slog.Info("ClientCancelBet开始", "玩家", player.PlayerId, "玩家状态为", player.PlayerState)
	costGold := ut.Money2Gold(player.BetMap[int(reqBetId)])
	//调用游戏中心去修改玩家钱包信息，这里是还原下注金额
	balance, err := slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
		Pid:     plr.PID,
		Change:  costGold,
		RoundID: strconv.Itoa(int(r.RoundId)),
		Reason:  slotsmongo.ReasonRefund,
	})
	if err != nil {
		slog.Error("ClientCancelBet Err", "slotsmongo.ModifyGold err", err)
		return
	}
	sf := ut.NewSnowflake()
	id := sf.NextID()
	betId := strconv.Itoa(int(int32(id % math.MaxInt32)))
	nowTime := time.Now().Unix()
	betLog := &slotsmongo.AddBetLogParamsAviator{
		ID:               strconv.FormatInt(id, 10),
		Pid:              plr.PID,
		AppID:            plr.AppID,
		Uid:              plr.Uid,
		UserName:         strconv.FormatInt(plr.PID, 10),
		CurrencyKey:      plr.CurrencyKey,
		RoomId:           plr.AppID,
		Bet:              ut.Money2Gold(player.BetMap[int(reqBetId)]),
		Balance:          balance,
		CashOutDate:      nowTime,
		RoundBetId:       betId,
		RoundId:          r.RoundId,
		GameId:           comm.GameID,
		LogType:          comm.LogType_ReMove,
		FinishType:       comm.FinishType_0,
		ManufacturerName: "spribe",
		BetId:            int(reqBetId),
		GameType:         define.GameType_Mini,
		Completed:        false,
	}
	slog.Info("ClientCancelBet时,即将插入取消投注记录", "入参: belog", betLog)
	slotsmongo.AddBetLogAviator(betLog) //待改为update
	player.BetMap[int(reqBetId)] = 0
	historyBetKey := fmt.Sprintf("%s_%v", playerIdStr, int(reqBetId))
	delete(r.Bets, historyBetKey)
	player.PlayerState = PLAYERSTATE_OB
	slog.Info("ClientCancelBet中", "玩家", player.PlayerId, "玩家状态已置为", player.PlayerState)
	content := &Content{C: "cancelBet"}
	content.P = P{
		BetID:    int(reqBetId), //客户端投注窗口区分1为左2为右
		Code:     comm.Succ,
		PlayerID: session.PlayerId,
	}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, _ := json.Marshal(body)
	marshal, _ := GetClientCancelBetRsp(&body).ToBinary()
	c.WriteMessage(messageType, marshal)
	NewBalance(c, r)

	r.CashOutFuncCancel(player.PlayerId)

	//c.WriteMessage(messageType, []byte(`AviatorLogin我返回了哈`))
	//c.WriteMessage(messageType, []byte(`{"id":13,"targetController":1,"content":{"c":"bet","p":{"bet":3,"betId":1,"code":200,"player_id":"51279","profileImage":"av-4.png","username":"5953102207"}}}`))
}

func GetClientCancelBetRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutString("player_id", initRsp.Content.P.PlayerID)
	if initRsp.Content.P.ErrorMessage != "" {
		pp.PutString("errorMessage", initRsp.Content.P.ErrorMessage)
	}
	pp.PutInt("betId", int32(initRsp.Content.P.BetID))

	p.PutSFSObject("p", pp)
	p.PutString("c", "cancelBet")

	so.AddCreatePAC(p, 1, 13)
	return so
}
