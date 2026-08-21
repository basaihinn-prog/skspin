package handlers

import (
	"encoding/json"
	"fmt"
	"game/comm/db"
	"game/comm/mq"
	"testing"

	"github.com/nats-io/nats.go"
)

func TestGameList(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "game")

	gameInfo := getGameBetInfo()
	for k, v := range gameInfo {
		fmt.Println(k, v)
	}
}

func TestPokdeng(t *testing.T) {
	nc, _ := nats.Connect("nats://127.0.0.1:11002")
	mq.SetNC(nc)

	// var lu struct {
	// 	Url string
	// }
	var lu json.RawMessage

	method := "/gamecenter/player/getBalance"
	method = "/pokdeng/game/launch"
	err := mq.Invoke(method, map[string]any{
		"Pid":      100001,
		"AppID":    "fake",
		"Uid":      "123",
		"GameID":   "pokdeng",
		"Language": "en",
	}, &lu)
	if err != nil {
		return
	}

	fmt.Println(string(lu))
}
