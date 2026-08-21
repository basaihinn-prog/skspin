package main

import (
	"context"
	"errors"
	"game/comm"
	"game/comm/mux"
	"game/duck/lang"
	"game/duck/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
)

func init() {
	//商户管理 / 批量修改RTP
	RegMsgProc("/AdminInfo/UpdateBatchGameRTP", "商户批量修改RTP", "AdminInfo", updateBatchOperatorGameRTP, updateBatchOperatorGameRTPParams{
		GameList: "",
		AppID:    "",
		RTP:      0,
	})
	// 内部调用
	mux.RegRpc("/AdminInfo/Interior/UpdateBatchGameRTP", "内部调用MQ-商户批量修改RTP", "AdminInfo", interBatchGameRTP, updateBatchOperatorGameRTPParams{
		GameList: "",
		AppID:    "",
		RTP:      0,
	})

}

type updateBatchOperatorGameRTPParams struct {
	GameList     string  `json:"GameList"`
	AppID        string  `json:"AppID"`
	RTP          float64 `json:"RTP"`
	BuyRTP       int     `json:"BuyRTP"`
	GamePattern  int64   `json:"GamePattern"`
	Operator     comm.Operator_V2
	MaxWinPoints int64 `json:"MaxWinPoints"`
	MaxMultiple  int64 `json:"MaxMultiple"`
	// 接口调用不传,内部赋值
	BuyRTPOff int
}

func updateBatchOperatorGameRTP(ctx *Context, ps updateBatchOperatorGameRTPParams, ret *comm.Empty) (err error) {

	user, err := GetUser(ctx)
	if err != nil {
		return
	}
	operator := comm.Operator_V2{}
	CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, &operator)

	if operator.RTPOff == 0 && user.AppID != "admin" {
		return errors.New(lang.Get(ctx.Lang, "未开启修改RTP权限"))
	}
	if !Contains(rtpListMap[int(operator.HighRTPOff)], int(ps.RTP)) && user.AppID != "admin" {
		return errors.New(lang.GetLang(ctx.Lang, "RTP值不在允许范围内"))
	}
	ps.Operator = operator

	// 验证RTP可以设置的最大值
	err, b := comm.VerifyMaxRTP(ps.RTP)
	if b == false {
		return err
	}
	ps.BuyRTPOff = int(operator.BuyRTPOff)

	interBatchGameRTP(ps, ret)

	return
}

func interBatchGameRTP(ps updateBatchOperatorGameRTPParams, ret *comm.Empty) (err error) {
	operator := ps.Operator

	gameList := []string{}

	if ps.GameList == "ALL" {
		for _, game := range MapGameBet {
			gameList = append(gameList, game.GameID)
		}
	} else {
		gameList = strings.Split(ps.GameList, ",")
	}

	updatelist := []mongo.WriteModel{}

	for _, gameid := range gameList {
		rw, nw := GetGameNwandAw(gameid, ps.RTP)
		settingInfo := bson.M{
			"RTP":            ps.RTP,
			"NoAwardPercent": nw,
			"RewardPercent":  rw,
			"GamePattern":    ps.GamePattern,
			"MaxMultiple":    ps.MaxMultiple,
			"MaxWinPoints":   ps.MaxWinPoints,
		}
		if ps.BuyRTPOff != 0 {
			settingInfo["BuyRTP"] = ps.BuyRTP
			settingInfo["BuyMinAwardPercent"] = MapBuyGameRTP[ps.BuyRTP]
		}
		model := &mongo.UpdateOneModel{
			Filter: bson.M{"AppID": ps.AppID, "GameId": gameid},
			Update: bson.M{"$set": settingInfo},
		}
		updatelist = append(updatelist, model)
	}

	result, err := CollGameConfig.Coll().BulkWrite(context.TODO(), updatelist, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return
	}
	logger.Infof("批量更新RTP设置,本次更新了%d条", result.UpsertedCount)
	go func() {
		gamelist := []*comm.GameConfig{}
		CollGameConfig.FindAll(bson.M{"AppID": ps.AppID, "GameId": bson.M{"$in": gameList}}, &gamelist)
		for _, config := range gamelist {
			req := updataGameCofigParams{
				AppID:              ps.AppID,
				GameId:             config.GameId,
				Preset:             config.Preset,
				StopLoss:           config.StopLoss,
				GamePattern:        config.GamePattern,
				FreeGameOff:        config.FreeGameOff,
				RTP:                config.RTP,
				BuyRTP:             config.BuyRTP,
				BuyMinAwardPercent: config.BuyMinAwardPercent,
				MaxMultiple:        strconv.Itoa(int(config.MaxMultiple)),
				MaxWinPoints:       config.MaxWinPoints,
				BetBase:            config.BetBase,
				GameOn:             config.GameOn,
				ShowNameAndTimeOff: operator.ShowNameAndTimeOff,
				ShowExitBtnOff:     operator.ShowExitBtnOff,
				ExitLink:           operator.ExitLink,
				//BetMult:            config.BetMult,
				DefaultCs:       config.DefaultCs,
				DefaultBetLevel: config.DefaultBetLevel,
				ProfitMargin:    config.ProfitMargin,
				CrashRate:       config.CrashRate,
				Scale:           config.Scale,
				OnlineUpNum:     config.OnlineUpNum,
				OnlineDownNum:   config.OnlineDownNum,
			}
			sp := strings.Split(config.BetBase, ",")
			cs := make([]float64, len(sp))
			for i := 0; i < len(sp); i++ {
				cs[i], _ = strconv.ParseFloat(sp[i], 64)
			}
			rw := int64(config.RewardPercent)
			nw := int64(config.NoAwardPercent)
			PushMsgStoreSetInfo(req, rw, nw, cs)
		}
	}()
	return
}
