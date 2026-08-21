package room

import (
	"github.com/lesismal/nbio/nbhttp/websocket"
	"serve/comm/ut"
)

func ClientSearchGifs(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	tag, _ := data.GetString("tag")
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()
	_ = tag
	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(200))
	gifsInfo := ut.NewSFSArray()
	//for _, gif := range topGifs {
	//	temp := ut.NewSFSObject()
	//	temp.PutString("id", gif.Id)
	//	temp.PutString("url", gif.Url)
	//	dims := ut.NewSFSArray()
	//	dims.Add(gif.Dims[0], ut.INT, true)
	//	dims.Add(gif.Dims[1], ut.INT, true)
	//	temp.PutSFSArray("dims", dims)
	//	gifsInfo.Add(temp, ut.SFS_OBJECT, true)
	//}
	pp.PutSFSArray("gifsInfo", gifsInfo)

	p.PutSFSObject("p", pp)
	p.PutString("c", "AddChatMessages")

	so.AddCreatePAC(p, 1, 13)
	marshal, _ := so.ToBinary()
	c.WriteMessage(messageType, marshal)
}

func GetClientSearchGifsRsp(initRsp *CommBody) *ut.SFSObject {
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
