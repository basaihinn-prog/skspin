package room

import (
	"encoding/json"
	"fmt"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"math/rand/v2"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"serve/serviceSpribeGame/aviator/room/playerInfo"
	"strconv"
)

func ClientLogin(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	var err error
	//req := ClientLoginReq{}
	//err = json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientLogin Err", "json.Unmarshal", err)
	//}
	tk, _ := data.GetString("cl")
	player, err := comm.TokenGetPlr(tk)
	if err != nil {
		slog.Error("ClientLogin::comm.TokenGetPlr Err", "err", err)
		rsp := GetErrMsg("", comm.ErrTokenExpired, "Token已经失效")
		rspMarshal, _ := json.Marshal(rsp)
		c.WriteMessage(messageType, rspMarshal)
		return
	}
	playerIdStr := strconv.Itoa(int(player.PID))
	c.SetSession(comm.WsSessionBody{
		PlayerId: playerIdStr,
		Plr:      player,
	})
	//获取spribe平台玩家信息
	spribePlayerInfo, err := playerInfo.FindPlayerInfo(player.PID, r.Name)
	if err != nil {
		slog.Error("ClientLogin::playerInfo.FindPlayerInfo Err", "err", err)
		rsp := GetErrMsg("", comm.ServerErr, "服务内部错误")
		rspMarshal, _ := json.Marshal(rsp)
		c.WriteMessage(messageType, rspMarshal)
		return
	}
	if spribePlayerInfo.Id.IsZero() {
		spribePlayerInfo = playerInfo.SpribePlayerInfoBody{
			Id:           primitive.NewObjectID(),
			PID:          player.PID,
			Username:     ut.GenerateRandomString2(5),
			AppID:        r.Name,
			ProfileImage: fmt.Sprintf(`av-%v.png`, rand.IntN(72)+1),
		}
		err = playerInfo.UpsertPlayerInfo(spribePlayerInfo)
		if err != nil {
			slog.Error("ClientLogin::playerInfo.UpsertPlayerInfo Err", "err", err.Error())
			rsp := GetErrMsg("", comm.ServerErr, "服务内部错误")
			rspMarshal, _ := json.Marshal(rsp)
			c.WriteMessage(messageType, rspMarshal)
			return
		}
	}
	client := &Client{
		Conn:        c,
		PlayerState: PLAYERSTATE_OB,
		PlayerId:    playerIdStr,
		PlayerName:  spribePlayerInfo.Username,
		PlayerIcon:  spribePlayerInfo.ProfileImage,
		AppId:       player.AppID,
		BetMap:      make(map[int]float64),
	}
	r.Players[playerIdStr] = client

	//rsp := ClientLoginRsp{
	//	Content: Content{
	//		Ct: 1024,   //todo 待替换
	//		Ms: 500000, //todo 待替换
	//		Tk: tk,
	//	},
	//}
	//marshal, err := json.Marshal(rsp)
	//if err != nil {
	//	slog.Error("ClientLogin::json.Marshal Err", "err", err)
	//	rsp := GetErrMsg("", comm.ServerErr, "服务内部错误")
	//	rspMarshal, _ := json.Marshal(rsp)
	//	c.WriteMessage(messageType, rspMarshal)
	//	return
	//}
	//
	//marshal, err := json.Marshal(rsp)
	//if err != nil {
	//	slog.Error("ClientLogin::json.Marshal Err", "err", err)
	//	rsp := GetErrMsg("", comm.ServerErr, "服务内部错误")
	//	rspMarshal, _ := json.Marshal(rsp)
	//	c.WriteMessage(messageType, rspMarshal)
	//	return
	//}
	rspNew := GetClientLoginRsp(tk)
	marshal, _ := rspNew.ToBinary()
	c.WriteMessage(messageType, marshal)
}

type ClientLoginReq struct {
	Channel          string  `json:"channel"`
	TargetController int     `json:"targetController"`
	Content          Content `json:"content"`
}

type ClientLoginRsp struct {
	Content Content `json:"content"`
}

func GetClientLoginRsp(tk string) *ut.SFSObject {
	//自己做一个
	so := &ut.SFSObject{}
	so.Init()
	so.PutByte("c", 0)
	so.PutShort("a", 0)
	p := &ut.SFSObject{}
	p.Init()
	p.PutInt("ct", 1024)
	p.PutInt("ms", 500000)
	p.PutString("tk", tk)
	so.PutSFSObject("p", p)
	return so
}
