package room

import (
	"fmt"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log/slog"
	"serve/comm/redisx"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
)

func ClientGetTopRoundsInfo(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	//req := CommBody{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientGetTopRoundsInfo Err", "json.Unmarshal", err)
	//}
	//获取top信息
	//获取top信息
	p, _ := data.GetSFSObject("p")
	period, _ := p.GetString("period")
	keyTopWin := fmt.Sprintf("%s_%s_%s", r.Name, "topround", period)
	tops, err := GetSortedTopRounds(redisx.GetClient(), keyTopWin)
	if err != nil {
		slog.Error("ClientGetTopRoundsInfo Err", "GetTopRoundsInfo err", err)
	}

	content := &Content{
		C: "getTopRoundsInfo",
		P: P{
			TopRounds: tops,
			Code:      comm.Succ,
		},
	}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, _ := json.Marshal(body)
	marshal, _ := GetClientGetTopRoundsInfoRsp(&body).ToBinary()
	c.WriteMessage(messageType, marshal)
}
func GetClientGetTopRoundsInfoRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	topWins := ut.NewSFSArray()
	for _, topRound := range initRsp.Content.P.TopRounds {
		temp := ut.NewSFSObject()
		temp.PutDouble("maxMultiplier", topRound.MaxMultiplier)
		temp.PutLong("endDate", topRound.EndDate)
		temp.PutString("zone", topRound.Zone)
		temp.PutLong("roundStartDate", topRound.RoundStartDate)
		temp.PutInt("roundId", int32(topRound.RoundId))
		temp.PutString("serverSeed", topRound.ServerSeed)

		topWins.Add(temp, ut.SFS_OBJECT, true)
	}
	pp.PutSFSArray("topRounds", topWins)

	p.PutSFSObject("p", pp)
	p.PutString("c", "getTopRoundsInfo")

	so.AddCreatePAC(p, 1, 13)
	return so
}
