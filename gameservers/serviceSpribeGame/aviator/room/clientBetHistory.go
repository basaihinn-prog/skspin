package room

import (
	"encoding/json"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log/slog"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"strconv"
)

func ClientBetHistory(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	//req := CommBody{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientBetHistory Err", "json.Unmarshal", err)
	//}
	p, _ := data.GetSFSObject("p")
	reqBetId, _ := p.GetDouble("lastBetId")
	reqLastRoundBetIdStr := strconv.FormatFloat(reqBetId, 'f', -1, 64)
	session, err := comm.GetSession(c)
	if err != nil {
		slog.Error("ClientBetHistory Err", "GetSession err", err)
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
	historyList := &comm.GetBetLogListNewResults{}
	err = comm.GetBetLogListNew(comm.GetBetLogAviatorListNewParams{
		Pid:            plr.PID,
		AppID:          plr.AppID,
		GameId:         comm.GameID,
		PageSize:       10,
		PageIndex:      1,
		FinishType:     comm.FinishType_1,
		LastRoundBetId: reqLastRoundBetIdStr,
	}, historyList)
	if err != nil {
		slog.Error("ClientBetHistory Err", "GetBetLogListNew", err)
	}
	lastBetId := int32(0)
	bets := make([]*Bets, 0)
	for _, temp := range historyList.List {
		bet := &Bets{
			Bet:           ut.Gold2Money(temp.Bet),
			CashOutDate:   temp.CashOutDate,
			Currency:      temp.CurrencyKey,
			Payout:        temp.Payout,
			Profit:        ut.Gold2Money(temp.Profit),
			RoundBetId:    temp.RoundBetId,
			RoundId:       temp.RoundId,
			Win:           true,
			WinAmount:     ut.Gold2Money(temp.Win),
			MaxMultiplier: 0,
		}
		bets = append(bets, bet)
		tempBetId, _ := strconv.ParseInt(bet.RoundBetId, 10, 32)
		lastBetId = int32(tempBetId)
	}
	content := &Content{
		C: "betHistoryResponse",
		P: P{
			Bets:      bets,
			Code:      comm.Succ,
			LastBetId: lastBetId,
		},
	}
	rsp := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, _ := json.Marshal(rsp)
	marshal, _ := GetClientBetHistoryRsp(&rsp).ToBinary()
	c.WriteMessage(messageType, marshal)
}
func GetClientBetHistoryRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutInt("lastBetId", int32(initRsp.Content.P.LastBetId))
	bets := ut.NewSFSArray()
	for _, bet := range initRsp.Content.P.Bets {
		temp := ut.NewSFSObject()
		temp.PutDouble("maxMultiplier", bet.MaxMultiplier)
		atoi, _ := strconv.Atoi(bet.RoundBetId)
		temp.PutInt("roundBetId", int32(atoi))
		temp.PutInt("betId", int32(bet.BetID))
		temp.PutString("currency", bet.Currency)
		temp.PutInt("roundId", int32(bet.RoundId))
		temp.PutBool("win", bet.Win)
		temp.PutDouble("profit", bet.Profit)
		temp.PutDouble("winAmount", bet.WinAmount)
		temp.PutBool("isFreeBet", false)
		temp.PutDouble("payout", bet.Payout)
		temp.PutDouble("bet", bet.Bet)
		temp.PutLong("cashOutDate", bet.CashOutDate)
		bets.Add(temp, ut.SFS_OBJECT, true)
	}
	pp.PutSFSArray("bets", bets)

	p.PutSFSObject("p", pp)
	p.PutString("c", "betHistoryResponse")

	so.AddCreatePAC(p, 1, 13)
	return so
}
