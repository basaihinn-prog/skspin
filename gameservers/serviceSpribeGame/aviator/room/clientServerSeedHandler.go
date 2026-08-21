package room

import (
	"github.com/lesismal/nbio/nbhttp/websocket"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
)

func ClientServerSeedHandler(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	seed := r.ServerSeed
	marshal, _ := GetClientServerSeedHandlerRsp(seed).ToBinary()
	c.WriteMessage(messageType, marshal)
}
func GetClientServerSeedHandlerRsp(seed string) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutString("serverSeedSHA256", slotsmongo.GenerateSha256(seed))
	pp.PutInt("code", 200)
	p.PutSFSObject("p", pp)
	p.PutString("c", "serverSeedResponse")

	so.AddCreatePAC(p, 1, 13)
	return so
}
