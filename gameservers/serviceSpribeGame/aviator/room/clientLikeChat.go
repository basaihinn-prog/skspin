package room

import (
	"encoding/json"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log/slog"
	"serve/comm/redisx"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"strconv"
)

func ClientLikeChat(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) { //todo
	//req := CommBody{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientLogin Err", "json.Unmarshal", err)
	//}
	p, _ := data.GetSFSObject("p")
	setLike, _ := p.GetBool("setLike")
	tempMessageId, _ := p.GetDouble("messageId")
	messageId := int64(tempMessageId)
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
		c.WriteMessage(messageType, rspMarshal)
	}
	//message, err := GetOneMessage(redisx.GetClient(), "room:"+r.Name+":message", messageId)
	messages, err := GetHashOne(redisx.GetClient(), "room:"+r.Name+":message", strconv.FormatInt(messageId, 10))
	if err != nil {
		slog.Error("GetOneMessage Err", "GetOneMessage err", err)
	}

	message := ChatHistory{}
	err = json.Unmarshal([]byte(messages), &message)
	if err != nil {
		slog.Error("GetOneMessage Err", "GetOneMessage err", err)
	}
	if setLike {
		if message.Likes == nil {
			message.Likes = &Likes{
				UsersLikesNumber: 0,
				UsersWhoLike:     make([]string, 0),
			}
		}
		message.Likes.UsersLikesNumber += 1
		message.Likes.UsersWhoLike = append(message.Likes.UsersWhoLike, session.PlayerId)
	} else {
		message.Likes.UsersLikesNumber -= 1
		message.Likes.UsersWhoLike = ut.Diff(message.Likes.UsersWhoLike, session.PlayerId)
	}

	bytes, _ := json.Marshal(message)
	//err = UpdateMessage(redisx.GetClient(), "room:"+r.Name+":message", float64(messageId), bytes)
	err = SetHashOne(redisx.GetClient(), "room:"+r.Name+":message", messageId, string(bytes))
	if err != nil {
		slog.Error("Message Err", "message update err", err)
	}

	//var res struct {
	//	Id               int `json:"id"`
	//	TargetController int `json:"targetController"`
	//	Content          struct {
	//		C string `json:"c"`
	//		P struct {
	//			Code             int   `json:"code"`
	//			MessageId        int64 `json:"messageId"`
	//			UsersLikesNumber int   `json:"usersLikesNumber"`
	//		} `json:"p"`
	//	} `json:"content"`
	//}
	res := &CommBody{
		Content: &Content{P: P{
			Code:      comm.Succ,
			MessageId: messageId,
		}},
	}
	res.Id = 13
	res.TargetController = 1
	res.Content.C = "like"
	res.Content.P.UsersLikesNumber = int64(message.Likes.UsersLikesNumber)

	//resbytes, _ := json.Marshal(res)
	marshal, _ := GetClientLikeChatRsp(res).ToBinary()
	c.WriteMessage(messageType, marshal)
}

func GetClientLikeChatRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutInt("messageId", int32(initRsp.Content.P.MessageId))
	pp.PutInt("usersLikesNumber", int32(initRsp.Content.P.UsersLikesNumber))

	p.PutSFSObject("p", pp)
	p.PutString("c", "AddChatMessages")

	so.AddCreatePAC(p, 1, 13)
	return so
}
