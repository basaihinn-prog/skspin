package facaicomm

//
//import (
//	"context"
//	"fmt"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"serve/comm/db"
//	"serve/comm/ut"
//	"strconv"
//	"time"
//
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//)
//
//type HistoryData struct {
//	RoundId        primitive.ObjectID    `json:"roundId" bson:"_id"`
//	Pid            int64                 `json:"-" bson:"pid"`
//	BetId          string                `json:"betId" bson:"betId"`
//	DateTime       int64                 `json:"dateTime"`
//	Bet            float64               `json:"bet"`
//	Win            float64               `json:"win"`
//	Balance        float64               `json:"balance"`
//	Details        string                `json:"roundDetails,omitempty" bson:"details"`
//	RoundDetails   []*HistoryDetailsData `json:"-" `
//	Currency       string                `json:"currency"`
//	CurrencySymbol string                `json:"currencySymbol"`
//	Hash           string                `json:"hash"`
//	GameConfig     *ReplayGameConfig     `json:"gameConfig" bson:"gameConfig"`
//	BaseBet        float64               `json:"baseBet" bson:"baseBet"`
//	OId            primitive.ObjectID    `bson:"oid"`
//}
//
//type HistoryDetailsData struct {
//	RoundId        primitive.ObjectID `json:"roundId"`
//	BetId          string             `json:"betId" bson:"betId"`
//	ConfigHash     string             `json:"configHash"`
//	Request        map[string]string  `json:"request"`
//	Response       Variables          `json:"response"`
//	Currency       string             `json:"currency"`
//	CurrencySymbol string             `json:"currencySymbol"`
//}
//
//func InsertBetHistoryEvery(plr *Player, bet float64, isNew bool, isFree bool, ps, ret Variables, currency, currencySymbol string, baseBet float64, oid primitive.ObjectID, init, betId string) {
//	coll := db.Collection("BetHistory")
//	rid := plr.BetHistLastID
//	request := map[string]string{}
//	if _, ok := ps["action"]; ok {
//		request["action"] = ps.Str("action")
//	}
//	if _, ok := ps["bl"]; ok {
//		request["bl"] = ps.Str("bl")
//	}
//	if _, ok := ps["c"]; ok {
//		request["c"] = ps.Str("c")
//	}
//	if _, ok := ps["counter"]; ok {
//		request["counter"] = ps.Str("counter")
//	}
//	if _, ok := ps["index"]; ok {
//		request["index"] = ps.Str("index")
//	}
//	if _, ok := ps["l"]; ok {
//		request["l"] = ps.Str("l")
//	}
//	if _, ok := ps["ind"]; ok {
//		request["ind"] = ps.Str("ind")
//	}
//	if _, ok := ps["repeat"]; ok {
//		request["repeat"] = ps.Str("repeat")
//	}
//	if _, ok := ps["symbol"]; ok {
//		request["symbol"] = ps.Str("symbol")
//	}
//	if _, ok := ps["pur"]; ok {
//		request["pur"] = ps.Str("pur")
//	}
//	if isNew {
//		bl := ret.Currency("balance")
//		win := ret.Currency("tw")
//		c, _ := strconv.ParseFloat(request["c"], 64)
//		l, _ := strconv.ParseFloat(request["l"], 64)
//		newHistoryData := &HistoryData{
//			RoundId:  rid,
//			Pid:      plr.PID,
//			BetId:    betId,
//			DateTime: time.Now().Unix() * 1000,
//			Bet:      bet,
//			Win:      win,
//			Balance:  bl,
//			RoundDetails: []*HistoryDetailsData{
//				{
//					RoundId:        rid,
//					BetId:          betId,
//					ConfigHash:     "",
//					Request:        request,
//					Response:       ret,
//					Currency:       currency,
//					CurrencySymbol: currencySymbol,
//				},
//			},
//			Currency:       currency,
//			CurrencySymbol: currencySymbol,
//			BaseBet:        c * l,
//			OId:            oid,
//		}
//		if isFree {
//			newHistoryData.Details = "Free spin"
//		}
//		coll.InsertOne(context.TODO(), newHistoryData)
//	} else {
//		bl := ret.Currency("balance")
//		update := bson.M{
//			"balance": bl,
//		}
//		if _, ok := ret["tw"]; ok {
//			update["win"] = ret.Currency("tw")
//		}
//
//		newHistoryData := &HistoryDetailsData{
//			RoundId:        rid,
//			BetId:          betId,
//			ConfigHash:     "",
//			Request:        request,
//			Response:       ret,
//			Currency:       currency,
//			CurrencySymbol: currencySymbol,
//		}
//		result := coll.FindOneAndUpdate(context.TODO(), bson.M{"_id": rid}, bson.M{"$set": update, "$push": bson.M{"rounddetails": newHistoryData}})
//		{
//			history := HistoryData{}
//			result.Decode(&history)
//
//			if history.Win >= history.BaseBet*10 && betId == "" {
//				gameConfig := newReplayGameConfig()
//				historyBetId := history.BetId
//				gameConfig.ReplayRoundId = historyBetId
//				gameConfig.CurrencyOriginal = history.Currency
//				gameConfig.Currency = history.Currency
//				gameConfig.Lang = ps.Get("lang")
//				token := ut.GenerateRandomString2(8)
//				update := bson.M{
//					"$set": bson.M{
//						"gameConfig": gameConfig,
//						"token":      token,
//					},
//				}
//				filter := bson.M{
//					"_id": rid,
//				}
//				coll.UpdateOne(context.TODO(), filter, update)
//				//历史记录表处理完毕
//				if betId == "" {
//					coll2 := db.Collection2("game", "ppReplayTokenMap")
//					body := ppReplayTokenMapBody{
//						Token:    token,
//						Gid:      coll.Database().Name(),
//						BetId:    history.BetId,
//						Init:     init,
//						CreateAt: time.Now().UnixMilli(),
//					}
//					filter = bson.M{
//						"betid": body.BetId,
//						"gid":   body.Gid,
//					}
//					_, err := coll2.ReplaceOne(context.TODO(), filter, body, options.Replace().SetUpsert(true))
//					if err != nil {
//						fmt.Println(err)
//					}
//				}
//			}
//		}
//	}
//}
//
//type ReplayGameConfig struct {
//	AmountType              string   `json:"amountType"`
//	EnvironmentId           int      `json:"environmentId"`
//	CurrencyOriginal        string   `json:"currencyOriginal"`
//	ReplayMode              bool     `json:"replayMode"`
//	Jurisdiction            string   `json:"jurisdiction"`
//	Currency                string   `json:"currency"`
//	SessionTimeout          string   `json:"sessionTimeout"`
//	StyleName               string   `json:"styleName"`
//	Lang                    string   `json:"lang"`
//	Region                  string   `json:"region"`
//	ReplaySystemUrl         string   `json:"replaySystemUrl"`
//	BrandRequirements       string   `json:"brandRequirements"`
//	ReplaySystemContextPath string   `json:"replaySystemContextPath"`
//	ReplayRoundId           string   `json:"replayRoundId"`
//	Mgckey                  string   `json:"mgckey"`
//	Datapath                string   `json:"datapath"`
//	SessionKey              []string `json:"sessionKey"`
//	SessionKeyV2            []string `json:"sessionKeyV2"`
//}
//
//func newReplayGameConfig() ReplayGameConfig {
//	gameConfig := ReplayGameConfig{
//		AmountType:              "COIN",
//		EnvironmentId:           200,
//		CurrencyOriginal:        "", //已赋值
//		ReplayMode:              true,
//		Jurisdiction:            "99",
//		Currency:                "", //已赋值
//		SessionTimeout:          "30",
//		StyleName:               "hllgd_hollygod",
//		Lang:                    "", //已赋值
//		Region:                  "Other",
//		ReplaySystemUrl:         "", //html中赋值
//		BrandRequirements:       "",
//		ReplaySystemContextPath: "", //html中赋值
//		ReplayRoundId:           "", //已赋值
//		Mgckey:                  "", //html中赋值
//		Datapath:                "", //html中赋值
//		SessionKey:              []string{"Zdch50DPJJAMsu2dPd+nETuOp+f+qGM+9TpbBmh+P243bJAgUzhqrs5v9lu7hnPAHme2RXbE0K/JiY0M7tjBlg==", "JD11U/vU0BKBrd8Qu02WUQ==", "bgMwtGxICcPNAR0j/GjDLG8XiEljUXVvTZV5TF0PiJKHq5hAOiRroVpQLYXpW6w4MvYYycgDNnynmVKEiW8JuThSKqAOg509xHBXludwBQrKL6ZadUo3sH2WiHO499wdW34m9q5LTKCaTPv7MCcF+XKjjOmMjqoKyXgHffYqbdv07/GI5rc2om5iWGUYcJ5bsZNM3EIkZYFrqWZR37ssyQGqd+ecEii+Vm/4LDw3JZBftbsooD0rlNzZ522pAOZg9Axbt1nBamoI/6MFj4Jc2qAJ67tpGic9CCoQjNonkWz3fcHBfN1zkd+/0gYTkR09xPHfrRnDz9ZJqPxXYLPxYg=="},
//		SessionKeyV2:            []string{"cgt7gcgWwj6tMSrshMvGO11mHQ2bDcFQnDyACBOGtg1YWO4nRtNvTP16ayKTCqnKRPHdWB8zz8auY5fsRutdkw==", "L3C87Vtzte5ttKic8wzI2TG33d+bH+/yt19YppwYwF8ywAjfkNdDmn6zo9ZYk9aM3gB98ha3eWkpDHEqNjn6Gw==", "f3VCqXHBwErPhhsLUzWLq/H1IiD83z2lsBH1roL+EIEHxOEq9KjPsgQNWXl01L1Xy0FBQ98xuoUSfR61NF3J/A==", "bfeo0E1rOVTMZ0QiFqJLv4CyAGCEfOfpZGuB+uj3eHXN5co4IgbtsWS2qR6Xz9jAfb50et/tVwZHupGWUZBx2w==", "aP6KLuUnE9y1vBsUG3bJXCa7m1UIExslTjB4sb92kDW5WsnMQmGbEHpD5CaoTz7dDgJrWYTL9vvnMlTsObwGEg=="},
//	}
//	return gameConfig
//}
//
//type ppReplayTokenMapBody struct {
//	Token    string `json:"token"`
//	Gid      string `json:"gid"`
//	BetId    string `json:"betId"`
//	Init     string `json:"init"`
//	CreateAt int64  `json:"createAt"`
//}
//
//type RspHistoryData struct {
//	RoundId string `json:"roundId" bson:"betId"`
//	Pid     int64  `json:"-" bson:"pid"`
//	//BetId          string                `json:"betId" bson:"betId"`
//	DateTime       int64                 `json:"dateTime"`
//	Bet            float64               `json:"bet"`
//	Win            float64               `json:"win"`
//	Balance        float64               `json:"balance"`
//	Details        string                `json:"roundDetails,omitempty" bson:"details"`
//	RoundDetails   []*HistoryDetailsData `json:"-" `
//	Currency       string                `json:"currency"`
//	CurrencySymbol string                `json:"currencySymbol"`
//	Hash           string                `json:"hash"`
//	GameConfig     *ReplayGameConfig     `json:"gameConfig" bson:"gameConfig"`
//	BaseBet        float64               `json:"baseBet" bson:"baseBet"`
//	OId            primitive.ObjectID    `bson:"oid"`
//}
