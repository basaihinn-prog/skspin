package room

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
	"regexp"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"strconv"
)

type MyReg struct {
	Number string
	Label  string
}

func (m *MyReg) Mate(str, substr string) bool {
	pattern := ``
	if str == "gif" {
		pattern = "^/" + str + ":(\\d+)$"
	} else {
		pattern = fmt.Sprintf(`^/%v:(-?\d+):(-?\d+)$`, str)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	matches := re.FindStringSubmatch(substr)
	if len(matches) > 1 {
		m.Label = matches[0]
		m.Number = matches[1]
		return true
	}
	return false
}

func ClientAddChatMessage(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) { //todo
	//req := CommBody{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientLogin Err", "json.Unmarshal", err)
	//}
	p, _ := data.GetSFSObject("p")
	msg, _ := p.GetString("message")
	session, err := comm.GetSession(c)
	if err != nil {
		slog.Error("ClientAddChatMessage Err", "GetSession err", err)
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
	playerIdStr := session.PlayerId
	playerId, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		slog.Error("ClientLogin Err", "strconv.ParseInt err", err)
	}
	fmt.Println("playerId", playerId)
	plrInfo := r.Players[playerIdStr]
	id, err := IncrementKey(redisx.GetClient(), "room:"+r.Name+":incr_id")
	if err != nil {
		slog.Error("Room Err", "Room Id err", err)
	}

	res := CommBody{
		Id:               13,
		TargetController: 1,
		Content: &Content{
			C: "AddChatMessages",
			P: P{
				Code: comm.Succ,
			},
		},
	}

	var m MyReg
	var chatHistory *ChatHistory
	if m.Mate("gif", msg) {
		coll := db.Collection2("spribe_01", "gif")
		var git struct {
			Id       string  `json:"id"`
			Preview  string  `json:"preview"`
			Duration float64 `json:"duration"`
			Size     int64   `json:"size"`
			Dims     []int   `json:"dims"`
			Url      string  `json:"url"`
		}

		err = coll.FindOne(context.TODO(), bson.M{"_id": m.Number}).Decode(&git)
		if err != nil {
			slog.Error("GIF Err", "GIF Get err", err)
		}
		jsonData, _ := json.Marshal(git)

		chatHistory = &ChatHistory{
			MessageID:    int(id),
			MessageType:  "gif",
			PlayerName:   plrInfo.PlayerName,
			ProfileImage: plrInfo.PlayerIcon,
			GifInfo:      string(jsonData),
		}
	} else if m.Mate("share_bet", msg) {

		collection, err := db.ClickHouseCollection("")
		if err != nil {
			slog.Error("ClickHouse Err", "ClickHouse Failure In Link", err)
		}

		var paper struct {
			Id                 int
			Bet                int64
			Win                int64
			Payout             float64
			RoundID            int64
			RoundMaxMultiplier float64
			CashOutDate        int64
			IsFreeBet          uint8
		}
		err = collection.QueryRow("select id,Bet,Win,Payout,RoundID,RoundMaxMultiplier,CashOutDate,Frb from aviatorgamelogs where AppID = ? and RoundBetId = ?", r.Name, m.Number).Scan(&paper.Id, &paper.Bet, &paper.Win, &paper.Payout, &paper.RoundID, &paper.RoundMaxMultiplier, &paper.CashOutDate, &paper.IsFreeBet)
		if err != nil {
			slog.Error("ClickHouse Err", "ClickHouse Query Failure", err)
		}

		if paper.Id == 0 {
			var reserr struct {
				Id               int `json:"id"`
				TargetController int `json:"targetController"`
				Content          struct {
					C string `json:"c"`
					P struct {
						Code         int    `json:"code"`
						ErrorMessage string `json:"errorMessage"`
					} `json:"p"`
				} `json:"content"`
			}

			reserr.Id = 13
			reserr.TargetController = 1
			reserr.Content.C = "AddChatMessages"
			reserr.Content.P.Code = 404
			reserr.Content.P.ErrorMessage = "Invalid share bet"

			resbyteserr, _ := json.Marshal(reserr)
			c.WriteMessage(messageType, resbyteserr)
			return
		}

		chatHistory = &ChatHistory{
			MessageID:    int(id),
			MessageType:  "win_info",
			PlayerName:   plrInfo.PlayerName,
			ProfileImage: plrInfo.PlayerIcon,
			WinInfo: WinInfo{
				Bet:                ut.Gold2Money(paper.Bet),
				WinAmount:          ut.Gold2Money(paper.Win),
				PlayerName:         plrInfo.PlayerName,
				Multiplier:         paper.Payout,
				Currency:           r.Currency,
				ProfileImage:       plrInfo.PlayerIcon,
				RoundId:            int(paper.RoundID),
				RoundMaxMultiplier: paper.RoundMaxMultiplier,
				CashOutDate:        paper.CashOutDate,
				IsFreeBet:          false,
			},
		}

	} else {
		chatHistory = &ChatHistory{
			Message:      msg,
			MessageID:    int(id),
			MessageType:  "message",
			PlayerName:   plrInfo.PlayerName,
			ProfileImage: plrInfo.PlayerIcon,
		}
	}
	bytes, _ := json.Marshal(chatHistory)
	err = SetHashOne(redisx.GetClient(), "room:"+r.Name+":message", id, string(bytes))
	//err = StoreMessage(redisx.GetClient(), "room:"+r.Name+":message", float64(id), bytes)
	if err != nil {
		slog.Error("Message Err", "message insert err", err)
	}
	res.Content.P.Messages = append(res.Content.P.Messages, chatHistory)
	//message, _ := json.Marshal(res)
	marshal, _ := GetClientAddChatMessageRsp(&res).ToBinary()
	r.NotifyAll2(marshal)
}

func GetClientAddChatMessageRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	messages := ut.NewSFSArray()
	messages = GetChatHistory(messages, initRsp.Content.P.Messages)
	pp.PutSFSArray("messages", messages)

	p.PutSFSObject("p", pp)
	p.PutString("c", "AddChatMessages")

	so.AddCreatePAC(p, 1, 13)
	return so
}
