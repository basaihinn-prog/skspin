package room

import (
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log/slog"
	"serve/comm/redisx"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
)

func ClientCurrentBetsInfo(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	content := &Content{
		C: "currentBetsInfoHandler", /*"currentBetsInfo"*/
	}
	var err error
	content.P.Bets, err = GetSortedBets(redisx.GetClient(), r.Name)
	if err != nil {
		slog.Error("GetSortedBets Err", "err", err)
	}
	content.P = P{
		Bets: content.P.Bets,
		Code: comm.Succ,
	}
	content.P.BetsCount = len(content.P.Bets)
	if content.P.BetsCount >= 50 {
		content.P.Bets = content.P.Bets[:50]
	}
	if len(content.P.Bets) == 0 {
		content.P.Bets = []*Bets{}
	}
	rsp := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, _ := json.Marshal(rsp)
	ss := GetClientCurrentBetsInfoRsp(&rsp)
	marshal, _ := ss.ToBinary()
	c.WriteMessage(websocket.BinaryMessage, marshal)

}

type ClientCurrentBetsInfoIdReq struct {
	Channel          string  `json:"channel"`
	TargetController int     `json:"targetController"`
	Content          Content `json:"content"`
}

type ClientCurrentBetsInfoIdRsp struct {
	ID               int     `json:"id"`
	TargetController int     `json:"targetController"`
	Content          Content `json:"content"`
}

func GetClientCurrentBetsInfoRsp(initRsp *CommBody) *ut.SFSObject {
	//自己做一个
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("betsCount", int32(initRsp.Content.P.BetsCount))
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	cashOuts := ut.NewSFSArray()
	for _, bet := range initRsp.Content.P.Cashouts {
		temp := ut.NewSFSObject()
		//temp.PutDouble("bet", bet.Bet)
		temp.PutString("player_id", bet.PlayerID)
		temp.PutInt("betId", int32(bet.BetID))
		temp.PutBool("isFreeBet", false)
		temp.PutString("currency", bet.Currency)
		//temp.PutString("profileImage", bet.ProfileImage)
		//temp.PutString("username", bet.Username)
		cashOuts.Add(temp, ut.SFS_OBJECT, true)
	}
	pp.PutSFSArray("cashOuts", cashOuts)
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
	p.PutString("c", "currentBetsInfo")
	so.AddCreatePAC(p, 1, 13)
	return so
}
