package room

import (
	"encoding/json"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log/slog"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"serve/serviceSpribeGame/aviator/room/playerInfo"
)

func ClientChangeProfileImage(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) { //todo
	//req := CommBody{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientChangeProfileImage Err", "json.Unmarshal", err)
	//}
	p, _ := data.GetSFSObject("p")
	profileImage, _ := p.GetString("profileImage")
	session, err := comm.GetSession(c)
	if err != nil {
		slog.Error("ClientChangeProfileImage Err", "GetSession err", err)
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
	plrInfo := r.Players[session.PlayerId]
	//profileImage := req.Content.P.ProfileImage
	spribePlayerInfo := playerInfo.SpribePlayerInfoBody{
		PID:          plr.PID,
		AppID:        r.Name,
		ProfileImage: profileImage,
		Username:     plrInfo.PlayerName,
	}
	err = playerInfo.UpsertPlayerInfo(spribePlayerInfo)
	if err != nil {
		slog.Error("ClientChangeProfileImage::playerInfo.UpsertPlayerInfo Err", "err", err.Error())
		rsp := GetErrMsg("", comm.ServerErr, "服务内部错误")
		rspMarshal, _ := json.Marshal(rsp)
		c.WriteMessage(messageType, rspMarshal)
		return
	}
	content := &Content{
		C: "changeProfileImage",
		P: P{
			Code: comm.Succ,
		},
	}
	body := CommBody{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, _ := json.Marshal(body)
	marshal, _ := GetClientChangeProfileImageRsp(&body).ToBinary()
	c.WriteMessage(messageType, marshal)
}

func GetClientChangeProfileImageRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutString("profileImageName", initRsp.Content.P.ProfileImage)

	p.PutSFSObject("p", pp)
	p.PutString("c", "changeProfileImage")

	so.AddCreatePAC(p, 1, 13)
	return so
}
