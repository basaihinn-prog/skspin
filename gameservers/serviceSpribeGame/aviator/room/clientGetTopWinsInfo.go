package room

import (
	"fmt"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log/slog"
	"serve/comm/redisx"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"strconv"
)

func ClientGetTopWinsInfo(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	//req := CommBody{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientGetHugeWinsInfo Err", "json.Unmarshal", err)
	//}
	p, _ := data.GetSFSObject("p")
	period, _ := p.GetString("period")
	//获取top信息
	keyTopWin := fmt.Sprintf("%s_%s_%s", r.Name, "topwin", period)
	tops, err := GetSortedTopWins(redisx.GetClient(), keyTopWin)
	if err != nil {
		slog.Error("ClientGetTopWinsInfo Err", "GetSortedTopWins err", err)
	}

	content := &Content{
		C: "getTopWinsInfo",
		P: P{
			TopWins: tops,
			Code:    comm.Succ,
		},
	}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, _ := json.Marshal(body)
	marshal, _ := GetClientGetTopWinsInfoRsp(&body).ToBinary()
	c.WriteMessage(messageType, marshal)
}
func GetClientGetTopWinsInfoRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	topWins := ut.NewSFSArray()
	for _, topWin := range initRsp.Content.P.TopWins {
		temp := ut.NewSFSObject()
		temp.PutDouble("maxMultiplier", topWin.MaxMultiplier)
		temp.PutLong("endDate", topWin.EndDate)
		temp.PutString("currency", topWin.Currency)
		temp.PutDouble("payout", topWin.Payout)
		temp.PutDouble("bet", topWin.Bet)
		temp.PutDouble("winAmount", topWin.WinAmount)
		temp.PutBool("isFreeBet", false)
		temp.PutString("profileImage", topWin.ProfileImage)
		temp.PutString("playerId", topWin.PlayerId)
		temp.PutString("username", topWin.Username)
		atoi, _ := strconv.Atoi(topWin.RoundBetId)
		temp.PutInt("roundBetId", int32(atoi))
		temp.PutDouble("winAmountInMainCurrency", topWin.WinAmountInMainCurrency)
		temp.PutString("zone", topWin.Zone)
		temp.PutInt("roundId", int32(topWin.RoundId))

		topWins.Add(temp, ut.SFS_OBJECT, true)
	}
	pp.PutSFSArray("topWins", topWins)

	p.PutSFSObject("p", pp)
	p.PutString("c", "getTopWinsInfo")

	so.AddCreatePAC(p, 1, 13)
	return so
}
