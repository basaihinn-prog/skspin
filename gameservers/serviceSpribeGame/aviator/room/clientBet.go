package room

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"serve/comm/define"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"strconv"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

func ClientBet(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	//req := CommBody{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientLogin Err", "json.Unmarshal", err)
	//}
	p, _ := data.GetSFSObject("p")
	reqBet, _ := p.GetDouble("bet")
	clientSeed, _ := p.GetString("clientSeed")
	reqBetId, _ := p.GetDouble("betId")
	freeBet, _ := p.GetBool("freeBet")
	autoCashOut, _ := p.GetDouble("autoCashOut")
	_ = freeBet
	session, err := comm.GetSession(c)
	player := r.Players[session.PlayerId]
	if err != nil {
		slog.Error("ClientBet Err", "GetSession err", err)
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
		c.WriteMessage(messageType, rspMarshal)
	}
	plr := session.Plr
	playerIdStr := strconv.Itoa(int(plr.PID))
	r.Players[playerIdStr].Seed = clientSeed
	plrInfo := r.Players[playerIdStr]
	if r.State != ROOMSTATE_WAIT_BET {
		slog.Error("ClientBet err::Room State Not Bet")
		rsp := CommBody{
			Id:               13,
			TargetController: 1,
		}
		rsp = GetErrMsg("bet", comm.ErrTimeout, "Stage time out", P{
			Code:         comm.ErrTimeout,
			BetID:        int(reqBetId),
			PlayerID:     playerIdStr,
			ProfileImage: player.PlayerIcon,
			Username:     player.PlayerName,
			Bet:          reqBet,
			ErrorMessage: "Stage time out",
		})
		//cashOutMarshal, _ := json.Marshal(rsp)
		cashOutMarshal, _ := GetClientErrRsp(&rsp).ToBinary()
		c.WriteMessage(messageType, cashOutMarshal)
		return
	}

	//if r.State != PLAYERSTATE_OB { //校验用户是否为ob状态，如果不是ob状态则不应该进入投注,如果已是投注状态不应该再次投注
	//	slog.Error("ClientBet::PlayerState Err", "player.PlayerState != ", "OB")
	//	return
	//}
	slog.Info("ClientBet开始", "玩家", player.PlayerId, "玩家状态为", player.PlayerState)
	costGold := ut.Money2Gold(reqBet)
	//调用游戏中心去修改玩家钱包信息，这里是扣除下注金额
	balance, err := slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
		Pid:     plr.PID,
		Change:  -costGold,
		RoundID: strconv.Itoa(int(r.RoundId)),
		Reason:  slotsmongo.ReasonBet,
	})
	if err != nil {
		err = define.PGNotEnoughCashErr
		rsp := GetErrMsg("bet", comm.ErrNoNotEnoughBalance, "Not enough balance", P{
			Code:         comm.ErrNoNotEnoughBalance,
			BetID:        int(reqBetId),
			PlayerID:     playerIdStr,
			ProfileImage: player.PlayerIcon,
			Username:     player.PlayerName,
			Bet:          reqBet,
			ErrorMessage: "Not enough balance",
		})
		//rspMarshal, _ := json.Marshal(rsp)
		rspMarshal, _ := GetClientErrRsp(&rsp).ToBinary()
		c.WriteMessage(messageType, rspMarshal)
		return
	}
	bet := Bets{
		Bet:          reqBet,
		BetID:        int(reqBetId),
		Currency:     plr.CurrencyKey,
		PlayerID:     playerIdStr,
		ProfileImage: player.PlayerIcon,
		Username:     player.PlayerName,
		WinAmount:    0,
	}
	err = StoreBet(redisx.GetClient(), r.Name, bet)
	if err != nil {
		slog.Error("StoreBet Err", "err", err)
		return
	}

	player.PlayerState = PLAYERSTATE_BET
	slog.Info("ClientBet中", "玩家", player.PlayerId, "玩家状态已置为", player.PlayerState, "BetId", int(reqBetId))
	player.BetMap[int(reqBetId)] = reqBet
	historyBetKey := fmt.Sprintf("%s_%v", playerIdStr, int(reqBetId))
	r.Bets[historyBetKey] = &Bets{
		BetID: int(reqBetId),
		Bet:   reqBet,
		//CashOutDate:   0,
		Currency: plr.CurrencyKey,
		//Payout:        0,
		//Profit:        0,
		//RoundBetId:    0,
		RoundId: r.RoundId,
		//Win:           0,
		//MaxMultiplier: 0,
		PlayerID:     playerIdStr,
		ProfileImage: player.PlayerIcon,
		Username:     player.PlayerName,
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
		LogType:          comm.LogType_Bet,
		FinishType:       comm.FinishType_0,
		ManufacturerName: "spribe",
		BetId:            bet.BetID,
		GameType:         define.GameType_Mini,
		Completed:        false,
	}
	slog.Info("ClientBet时,即将插入投注记录", "入参: belog", betLog)
	slotsmongo.AddBetLogAviator(betLog)

	content := &Content{C: "bet"}
	content.P = P{
		Bet:          reqBet,
		BetID:        int(reqBetId), //客户端投注窗口区分1为左2为右
		Code:         comm.Succ,
		PlayerID:     session.PlayerId,
		ProfileImage: plrInfo.PlayerIcon,
		Username:     plrInfo.PlayerName,
	}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, _ := json.Marshal(body)
	marshal, _ := GetClientBetRsp(&body).ToBinary()
	c.WriteMessage(messageType, marshal)
	NewBalance(c, r)

	//判定是否自动结算
	r.CashOutFuncAdd(int(reqBetId), reqBet, autoCashOut, player.PlayerId)
}

func GetClientBetRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutDouble("bet", initRsp.Content.P.Bet)
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutString("player_id", initRsp.Content.P.PlayerID)
	pp.PutBool("freeBet", initRsp.Content.P.FreeBet)
	pp.PutInt("betId", int32(initRsp.Content.P.BetID))
	pp.PutString("profileImage", initRsp.Content.P.ProfileImage)
	pp.PutString("username", initRsp.Content.P.Username)

	p.PutSFSObject("p", pp)
	p.PutString("c", "bet")

	so.AddCreatePAC(p, 1, 13)
	return so
}

//{
//"channel": "AviatorBetIdReq",
//"targetController": 1,
//"content": {
//"c": "bet",
//"r": -1,
//"p": {
//"bet": 10,
//"clientSeed": "5DwNcDa0qnZYdqyJ5W6i",
//"betId": 1,
//"freeBet": false,
//"autoCashOut": 1.1
//}
//}
//}
