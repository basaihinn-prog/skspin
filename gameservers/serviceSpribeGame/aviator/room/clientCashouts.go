package room

import (
	"encoding/json"
	"fmt"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log/slog"
	"math"
	"serve/comm/define"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"strconv"
	"time"
)

var DateList = []string{"day", "month", "year"}

func ClientCashOut(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) { //todo
	//req := CommBody{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("CashOut Err", "json.Unmarshal", err)
	//	rsp := GetErrMsg("cashOut", comm.ServerErr, "服务内部错误")
	//	rspMarshal, _ := json.Marshal(rsp)
	//	c.WriteMessage(messageType, rspMarshal)
	//}
	p, _ := data.GetSFSObject("p")
	tempReqBetId, _ := p.GetDouble("betId")
	reqBetId := int(tempReqBetId)
	multiplying, _ := p.GetDouble("multiplying")
	cashOutRsp := CommBody{
		Id:               13,
		TargetController: 1,
	}
	session, err := comm.GetSession(c)
	if err != nil {
		slog.Error("ClientCashOut Err", "GetSession err", err)
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
	if r.State != ROOMSTATE_FLY {
		slog.Error("ClientCashOut err::Room State Not Fly")
		cashOutRsp.Content = &Content{
			C: "cashOut",
			P: P{
				Cashouts: []*Cashouts{{
					BetID:    reqBetId,
					PlayerID: playerIdStr,
				}},
				Code:         comm.ErrTimeout,
				OperatorKey:  "release",
				ErrorMessage: "Stage time out",
			},
		}
		//cashOutMarshal, _ := json.Marshal(cashOutRsp)
		cashOutMarshal, _ := GetClientErrRsp(&cashOutRsp).ToBinary()
		c.WriteMessage(messageType, cashOutMarshal)
		return
	}
	player := r.Players[playerIdStr]
	slog.Info("ClientCashOut开始", "玩家", player.PlayerId, "玩家状态为", player.PlayerState)
	//校验有无该用户这次订单
	historyBetKey := fmt.Sprintf("%s_%v", playerIdStr, reqBetId)
	if _, ok := r.Bets[historyBetKey]; !ok {
		slog.Error("ClientCashOut err::No Bet Record")
		cashOutRsp.Content = &Content{
			C: "cashOut",
			P: P{
				Cashouts: []*Cashouts{{
					BetID:    reqBetId,
					PlayerID: playerIdStr,
				}},
				Code:         comm.ErrTimeout,
				OperatorKey:  "release",
				ErrorMessage: "Stage time out",
			},
		}
		//cashOutMarshal, _ := json.Marshal(cashOutRsp)
		cashOutMarshal, _ := GetClientErrRsp(&cashOutRsp).ToBinary()
		c.WriteMessage(messageType, cashOutMarshal)
		return
	}

	var multiplier float64
	if multiplying == float64(0) {
		multiplier = r.XArr[r.XIndex]
	} else {
		multiplier = multiplying
	}
	winAmount := math.Round(player.BetMap[reqBetId]*multiplier*100) / 100
	balance, err := slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
		Pid:     plr.PID,
		Change:  ut.Money2Gold(winAmount),
		RoundID: strconv.FormatInt(r.RoundId, 10),
		Reason:  slotsmongo.ReasonWin,
		IsEnd:   true,
	})
	if err != nil {
		rsp := GetErrMsg("cashOut", comm.ServerErr, "服务内部错误")
		//rspMarshal, _ := json.Marshal(rsp)
		rspMarshal, _ := GetClientErrRsp(&rsp).ToBinary()
		c.WriteMessage(messageType, rspMarshal)
		return
	}
	//_ = balance
	sf := ut.NewSnowflake()
	id := sf.NextID()
	betId := strconv.Itoa(int(int32(id % math.MaxInt32)))
	nowTime := time.Now().UnixMilli()
	betLog := &slotsmongo.AddBetLogParamsAviator{
		ID:                 strconv.FormatInt(id, 10),
		Pid:                plr.PID,
		AppID:              plr.AppID,
		Uid:                plr.Uid,
		UserName:           strconv.FormatInt(plr.PID, 10),
		CurrencyKey:        plr.CurrencyKey,
		RoomId:             plr.AppID,
		Bet:                ut.Money2Gold(player.BetMap[reqBetId]),
		Win:                ut.Money2Gold(winAmount),
		Balance:            balance,
		CashOutDate:        nowTime,
		Payout:             multiplier,
		Profit:             ut.Money2Gold(winAmount - player.BetMap[reqBetId]),
		MaxMultiplier:      multiplier,
		RoundMaxMultiplier: r.XArr[len(r.XArr)-1],
		RoundBetId:         betId,
		RoundId:            r.RoundId,
		GameId:             comm.GameID,
		LogType:            comm.LogType_Bet,
		FinishType:         comm.FinishType_1,
		ManufacturerName:   "spribe",
		BetId:              reqBetId,
		GameType:           define.GameType_Mini,
		Completed:          true,
	}
	slotsmongo.AddBetLogAviator(betLog)
	//将主动提现过的人做历史记录编辑标记，在飞机爆炸后进行编辑，把当次最大的倍数赋值给投注记录
	r.Bets[historyBetKey].HistoryFlag = true

	bet := Bets{
		Bet:          player.BetMap[reqBetId],
		BetID:        reqBetId,
		Currency:     plr.CurrencyKey,
		PlayerID:     playerIdStr,
		ProfileImage: player.PlayerIcon,
		Username:     player.PlayerName,
		WinAmount:    winAmount,
		Multiplier:   multiplier,
	}
	UpdateBet(redisx.GetClient(), r.Name, bet)
	topWin := TopWin{
		Bet:                     player.BetMap[reqBetId],
		EndDate:                 time.Now().UTC().UnixMilli(),
		Currency:                plr.CurrencyKey,
		MaxMultiplier:           r.XArr[len(r.XArr)-1],
		Payout:                  multiplier,
		PlayerId:                playerIdStr,
		ProfileImage:            player.PlayerIcon,
		RoundBetId:              betId,
		RoundId:                 r.RoundId,
		Username:                player.PlayerName,
		WinAmountInMainCurrency: winAmount,
		WinAmount:               winAmount,
		Zone:                    "aviator_core_inst2_demo1",
	}
	if topWin.Payout >= 4.0 {
		for _, date := range DateList {
			keyTopWin := fmt.Sprintf("%s_%s_%s", r.Name, "topwin", date)
			StoreTopWins(redisx.GetClient(), keyTopWin, topWin)
			keyHugeWin := fmt.Sprintf("%s_%s_%s", r.Name, "huge", date)
			StoreHugeWins(redisx.GetClient(), keyHugeWin, topWin)
		}
	}

	//移除自动提现池子中的该次投注
	tmp := make([]*CashOut, 0)
	for _, v := range r.CashOutArr {
		if reqBetId == v.BetId && playerIdStr == v.PlayerId {
			continue
		} else {
			tmp = append(tmp, v)
		}
	}
	r.CashOutArr = tmp

	cashContent := &Content{
		C: "cashOut",
		P: P{
			Cashouts: []*Cashouts{{
				BetAmount: player.BetMap[reqBetId],
				BetID:     reqBetId,
				PlayerID:  playerIdStr,
				WinAmount: winAmount,
				Win:       true,
			}},
			Code:        comm.Succ,
			Multiplier:  multiplier,
			OperatorKey: "release",
		},
	}
	cashOutRsp.Content = cashContent
	//cashOutMarshal, _ := json.Marshal(cashOutRsp)
	cashOutMarshal, _ := GetClientCashOutRsp(&cashOutRsp).ToBinary()
	c.WriteMessage(messageType, cashOutMarshal)
	NewBalance(c, r)
	player.PlayerState = PLAYERSTATE_OB //返回余额后将当前这个提现完的玩家状态重置成ob状态
}
func GetClientCashOutRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	cashOuts := ut.NewSFSArray()
	for _, bet := range initRsp.Content.P.Cashouts {
		temp := ut.NewSFSObject()
		temp.PutDouble("betAmount", bet.BetAmount)
		temp.PutDouble("winAmount", bet.WinAmount)
		temp.PutString("player_id", bet.PlayerID)
		temp.PutInt("betId", int32(bet.BetID))
		temp.PutBool("isMaxWinAutoCashOut", false) //todo
		cashOuts.Add(temp, ut.SFS_OBJECT, true)
	}
	pp.PutSFSArray("cashouts", cashOuts)
	pp.PutDouble("multiplier", initRsp.Content.P.Multiplier)
	pp.PutString("operatorKey", initRsp.Content.P.OperatorKey)

	p.PutSFSObject("p", pp)
	p.PutString("c", "cashOut")

	so.AddCreatePAC(p, 1, 13)
	return so
}
