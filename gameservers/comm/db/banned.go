package db

import (
	"context"
	"log/slog"
	"serve/comm/define"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// 是否被禁用
func IsBanned(game string, pid int64) error {
	var gameData struct {
		Status int `bson:"Status"`
	}
	err := Collection2("game", "Games").FindOne(context.TODO(), ID(game), options.FindOne().SetProjection(D("Status", 1))).Decode(&gameData)
	if err != nil {
		return define.NewErrCode("The game does not exist.", 1402)
	}
	// 关闭和隐藏后不再处理请求
	if gameData.Status == define.GameStatus_Stop || gameData.Status == define.GameStatus_Hide {
		slog.Info("Game is close")
		return define.NewErrCode("Invalid game.", 1302)
	}

	var playerData struct {
		Status int    `bson:"Status"`
		AppID  string `bson:"AppID"`
	}
	err = Collection2("game", "Players").FindOne(context.TODO(), ID(pid), options.FindOne().SetProjection(D("AppID", 1, "Status", 1))).Decode(&playerData)
	if err != nil {
		return define.NewErrCode("Player does not exist.", 1302)
	}
	// 1是正常状态，非正常状态不处理请求
	if playerData.Status != 1 {
		slog.Info("player is close")
		return define.NewErrCode("The player has been blocked", 1309)
	}
	// 维护状态，非fake运营商请求不处理
	if gameData.Status == define.GameStatus_Maintenance && !strings.HasPrefix(playerData.AppID, "fake") {
		return define.NewErrCode("The game is being maintained", 1302)
	}

	// 运营商判断
	var operatorData struct {
		Status                int      `bson:"Status"`
		ParentAppID           string   `bson:"Name"`
		DefaultManufacturerOn []string `bson:"DefaultManufacturerOn"`
	}
	err = Collection2("GameAdmin", "AdminOperator").FindOne(context.TODO(), D("AppID", playerData.AppID), options.FindOne().SetProjection(D("Status", 1, "Name", 1, "DefaultManufacturerOn", 1))).Decode(&operatorData)
	if err != nil {
		return define.NewErrCode("Operator not exist.", 1310)
	}
	// 运营商
	if operatorData.Status != 1 {
		slog.Info("Operator is close")
		return define.NewErrCode("Operator is close.", 1310)
	}
	// 运营商游戏厂商判断
	ManufacturerOn := false
	if len(operatorData.DefaultManufacturerOn) != 0 {
		for _, Manufacturer := range operatorData.DefaultManufacturerOn {
			if strings.HasPrefix(game, strings.ToLower(Manufacturer)) {
				ManufacturerOn = true
			}
		}
		if !ManufacturerOn {
			slog.Info("ManufacturerOn Game is close")
			return define.NewErrCode("ManufacturerOn Game is close.", 1302)
		}
	}

	// 线路商判断
	if operatorData.ParentAppID != "admin" {
		var parentOperatorData struct {
			Status int `bson:"Status"`
		}
		err = Collection2("GameAdmin", "AdminOperator").FindOne(context.TODO(), D("AppID", operatorData.ParentAppID), options.FindOne().SetProjection(D("Status", 1))).Decode(&parentOperatorData)
		if err != nil {
			return define.NewErrCode("Player Operator not exist.", 1310)
		}
		if parentOperatorData.Status != 1 {
			slog.Info("Parent Operator is close")
			return define.NewErrCode("Parent Operator is close.", 1310)
		}
	}

	// 商户内游戏状态判断
	var operatorGameStatus struct {
		GameOn int `bson:"GameOn"`
	}
	err = Collection2("GameAdmin", "GameConfig").FindOne(context.TODO(), D("AppID", playerData.AppID, "GameId", game), options.FindOne().SetProjection(D("GameOn", 1))).Decode(&operatorGameStatus)
	if err != nil {
		return define.NewErrCode("Operator not exist.", 1310)
	}
	// 运营商
	if operatorGameStatus.GameOn != 0 {
		slog.Info("Operator Game is close")
		return define.NewErrCode("Operator Game is close.", 1302)
	}

	return nil
}
