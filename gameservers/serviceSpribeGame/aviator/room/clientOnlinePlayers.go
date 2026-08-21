package room

import (
	"github.com/lesismal/nbio/nbhttp/websocket"
)

func ClientOnlinePlayers(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
	//req := ClientOnlinePlayersReq{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientOnlinePlayers Err", "json.Unmarshal", err)
	//}
	////rsp := ClientOnlinePlayersRsp{
	////	Content: Content{
	////		Ct: 1024,
	////		Ms: 500000,
	////		Tk: req.Content.Cl,
	////	},
	////}
	////marshal, _ := json.Marshal(rsp)
	////c.WriteMessage(messageType, marshal)
	////c.WriteMessage(messageType, []byte(`CashOut我返回了哈`))
	//
	//content := &Content{C: "onlinePlayers", P: P{Code: 200, OnlinePlayers: len(r.Players)}}
	//body := CommBody{
	//	Id:               13,
	//	TargetController: 1,
	//	Content:          content,
	//}
	//marshal, err := json.Marshal(body)
	//if err != nil {
	//	slog.Error("json.Marshal", "Err", err.Error())
	//}
	//c.WriteMessage(websocket.TextMessage, marshal)
	//
	//
	//c.WriteMessage(messageType, []byte(`{"id":13,"targetController":1,"content":{"c":"onlinePlayers","p":{"code":200,"onlinePlayers":5001}}}`))

}

type ClientOnlinePlayersReq struct {
	Channel          string  `json:"channel"`
	TargetController int     `json:"targetController"`
	Content          Content `json:"content"`
}

type ClientOnlinePlayersRsp struct {
	ID               int     `json:"id"`
	TargetController int     `json:"targetController"`
	Content          Content `json:"content"`
}
