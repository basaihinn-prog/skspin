package slotsmongo

import (
	"context"
	"game/comm"
	"game/comm/db"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

// 是否被禁用
func IsBanned(game string, pid int64) bool {
	var gameData struct {
		Status int `bson:"Status"`
	}
	err := db.Collection2("game", "Games").FindOne(context.TODO(), db.ID(game), options.FindOne().SetProjection(db.D("Status", 1))).Decode(&gameData)
	if err != nil {
		return true
	}
	// 关闭和隐藏后不再处理请求
	if gameData.Status == comm.GameStatus_Stop || gameData.Status == comm.GameStatus_Hide {
		return true
	}
	var playerData struct {
		Status int    `bson:"Status"`
		AppID  string `bson:"AppID"`
	}
	err = db.Collection2("game", "Players").FindOne(context.TODO(), db.ID(pid), options.FindOne().SetProjection(db.D("AppID", 1, "Status", 1))).Decode(&playerData)
	if err != nil {
		return true
	}
	// 0是正常状态，非正常状态不处理请求
	if playerData.Status != 0 {
		return true
	}
	// 维护状态，非fake运营商请求不处理
	if gameData.Status == comm.GameStatus_Maintenance && !strings.HasPrefix(playerData.AppID, "fake") {
		return true
	}

	// 运营商判断
	var operatorData struct {
		Status int `bson:"Status"`
	}
	err = db.Collection2("GameAdmin", "AdminOperator").FindOne(context.TODO(), db.D("AppID", playerData.AppID), options.FindOne().SetProjection(db.D("Status", 1))).Decode(&operatorData)
	if err != nil {
		return true
	}
	// 运营商
	if operatorData.Status != 0 {
		return true
	}

	return false
}
