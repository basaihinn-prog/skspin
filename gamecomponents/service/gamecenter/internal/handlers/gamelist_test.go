package handlers

import (
	"fmt"
	"game/comm/db"
	"testing"
)

func TestGameList(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "game")

	gameInfo := getGameBetInfo()
	for k, v := range gameInfo {
		fmt.Println(k, v)
	}
}
