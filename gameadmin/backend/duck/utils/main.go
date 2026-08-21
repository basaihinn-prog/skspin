package main

import (
	"context"
	"encoding/json"
	"fmt"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

// 运营商游戏设置
type GameConfig struct {
	//游戏编号
	GameId string `json:"GameId" bson:"GameId"`

	//运营商
	AppID string `json:"AppID" bson:"AppID" md:"运营商"`

	//配置文件
	ConfigPath string `json:"ConfigPath" bson:"ConfigPath"`

	//RTP设置
	RTP float64 `json:"RTP" bson:"RTP"`

	//止盈止损开关
	StopLoss int64 `json:"StopLoss" bson:"StopLoss"`

	//赢取最高押注倍数
	MaxMultipleOff int64 `json:"MaxMultipleOff" bson:"MaxMultipleOff"`
	MaxMultiple    int64 `json:"MaxMultiple" bson:"MaxMultiple"`

	//游戏投注
	BetBase string `json:"BetBase" bson:"BetBase"`

	//游戏类型
	GamePattern int `json:"GamePattern" bson:"GamePattern"`

	//预设面额
	Preset float32 `json:"Preset" bson:"Preset"`

	//游戏状态
	GameOn int64 `json:"GameOn" bson:"GameOn"`

	//免费游戏开关
	FreeGameOff int64 `json:"FreeGameOff" bson:"FreeGameOff"`

	//试玩链接
	GameDemo string `json:"GameDemo" bson:"GameDemo"`

	//下注的3% 个人奖池累计配置千分比 所有slots 通用,  普通, 购买不算
	RewardPercent int `json:"RewardPercent" bson:"RewardPercent"`

	//个人奖池为负数，不中奖概率千分比  20%不中奖,  优先级最高
	NoAwardPercent int `json:"NoAwardPercent" bson:"NoAwardPercent"`
}

var client *mongo.Client

type operatorGameConfig struct {
	//游戏编号	game/Games/_Id
	GameId string `json:"GameId" bson:"GameId"`

	GameIdd string `json:"GameIdd" bson:"_id"`
	//游戏分类	game/Games/Type
	Gametype int `json:"Gametype" bson:"Type"`

	//游戏名称	game/Games/Name
	GameName       string                     `json:"GameName" bson:"Name"`
	GameNameConfig map[string]*GameConfigItem `json:"GameNameConfig" bson:"GameNameConfig"`

	//配置文件
	ConfigPath string `json:"ConfigPath" bson:"ConfigPath"`

	//RTP设置
	RTP float64 `json:"RTP" bson:"RTP"`

	//止盈止损开关
	StopLoss int64 `json:"StopLoss" bson:"StopLoss"`

	//赢取最高押注倍数
	MaxMultiple int64 `json:"MaxMultiple" bson:"MaxMultiple"`

	//游戏投注
	BetBase string `json:"BetBase" bson:"BetBase"`

	//游戏类型
	GamePattern int `json:"GamePattern" bson:"GamePattern"`

	//预设面额
	Preset float32 `json:"Preset" bson:"Preset"`

	//游戏状态
	GameOn int64 `json:"GameOn" bson:"GameOn"`

	//免费游戏开关
	FreeGameOff int64 `json:"FreeGameOff" bson:"FreeGameOff"`

	//试玩链接
	GameDemo string `json:"GameDemo" bson:"GameDemo"`

	//RTP设置功能开关
	RTPOff int64 `json:"RTPOff" bson:"RTPOff"`
	//止损止赢功能开关
	StopLossOff int64 `json:"StopLossOff" bson:"StopLossOff"`
	//赢取最高押注倍数设置功能开关
	MaxMultipleOff int64 `json:"MaxMultipleOff" bson:"MaxMultipleOff"`

	//名字和时间显示
	ShowNameAndTimeOff int64 `json:"ShowNameAndTimeOff" bson:"ShowNameAndTimeOff"`
	//退出按钮显示
	ShowExitBtnOff int64 `json:"ShowExitBtnOff" bson:"ShowExitBtnOff"`
	//退出按钮链接
	ExitLink string `json:"ExitLink" bson:"ExitLink"`
}

type GameNow struct {
	GameID        string `json:"gameId"`
	RewardPercent int64  `json:"reward_percent"`
}

type GameBet struct {
	GameID string `json:"Game_id"`
	Bet    string `json:"Bet"`
}

var MapGameBet = map[string]string{}
var MapGameNow = map[string]int64{}

func init() {
	fileBytes, err := os.ReadFile("./config/output.json")

	var data []*GameNow

	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	for _, vl := range data {
		MapGameNow[vl.GameID] = vl.RewardPercent
	}

	fileBytes, err = os.ReadFile("./config/GameBet.json")

	gameBet := []GameBet{}

	err = json.Unmarshal(fileBytes, &gameBet)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	for _, vl := range gameBet {
		MapGameBet[vl.GameID] = vl.Bet
	}

}

func main() {
	var appid = "faketrans"

	client, _ = mongodb.Connect("mongodb://myAdmin:myAdminPassword1@47.237.3.219:27017/?authSource=admin")

	//GetGameNwandAw("123", 1.2)

	//err := client.Ping(context.TODO(), readpref.Primary())
	//if err != nil {
	//	log.Fatal(err)
	//}

	collgame := client.Database("game").Collection("Games")

	var list_OperatorGameConfig []*operatorGameConfig

	cursor, err := collgame.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	err = cursor.All(context.TODO(), &list_OperatorGameConfig)
	if err != nil {
		log.Fatal(err)
	}
	collgameconfig := client.Database("GameAdmin").Collection("GameConfig")

	for _, config := range list_OperatorGameConfig {
		_, err := collgameconfig.DeleteMany(context.TODO(), bson.M{"AppID": appid, "GameId": config.GameIdd})
		if err != nil {
			log.Fatal(err)
		}
		operatorConfig := bson.M{}
		operatorConfig["AppID"] = appid
		operatorConfig["GameId"] = config.GameIdd    //游戏编号
		operatorConfig["ConfigPath"] = "config/path" //配置文件
		operatorConfig["StopLoss"] = 500             //止盈止损开关
		operatorConfig["MaxMultipleOff"] = 1         //赢取最高押注倍数
		operatorConfig["MaxMultiple"] = 10000
		operatorConfig["BetBase"] = MapGameBet[config.GameIdd] //游戏投注
		operatorConfig["GamePattern"] = 2                      //游戏类型
		operatorConfig["Preset"] = 10                          //预设面额
		operatorConfig["GameOn"] = 0                           //游戏状态  0 开启  1关闭
		operatorConfig["ShowNameAndTimeOff"] = 1               //显示游戏名称和时间
		operatorConfig["ShowExitBtnOff"] = 1                   //显示退出按钮
		operatorConfig["ExitLink"] = "http://ExitLinkPage.com" //退出按钮链接

		operatorConfig["RTP"] = 96 //RTP设置
		rw, nw := GetGameNwandAw(config.GameIdd, 96)
		operatorConfig["RewardPercent"] = rw
		operatorConfig["NoAwardPercent"] = nw
		_, err = collgameconfig.InsertOne(context.TODO(), operatorConfig)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("写入", config.GameIdd)
	}

}

type Game struct {
	ID               string `json:"ID" bson:"_id"`
	GameID           string `json:"GameID" bson:"GameID"`
	Name             string `json:"Name" bson:"Name"`
	Type             int    `json:"Type" bson:"Type"`
	Status           int    `json:"Status" bson:"Status"`
	ManufacturerName string `json:"ManufacturerName" bson:"ManufacturerName"`
	IconUrl          string
	WebIconUrl       string
	ExtraData        map[string]any             `bson:"-"`
	BuyBetMulti      int                        `json:"BuyBetMulti" bson:"BuyBetMulti"`
	Bet              int64                      `json:"-" bson:"Bet"`
	GameNameConfig   map[string]*GameConfigItem `json:"GameNameConfig" bson:"GameNameConfig"`
	NoAwardPercent   int64                      `json:"NoAwardPercent" bson:"NoAwardPercent"`
}
type GameConfigItem struct {
	GameName string
	Icon     string
}

func GetGameNwandAw(gameId string, RTP float64) (rw, nw int64) {
	//var gameinfo *Game
	//single := client.Database("game").Collection("Games").FindOne(context.TO/DO(), bson.M{"_id": gameId})
	//err := single.Decode(&gameinfo)

	var k int64
	switch RTP {
	case 10:
		k = -840
		nw = 900
	case 20:
		k = -740
		nw = 800
	case 30:
		k = -640
		nw = 700
	case 40:
		k = -540
		nw = 600
	case 50:
		k = -440
		nw = 500
	case 60:
		k = -340
		nw = 400
	case 70:
		k = -240
		nw = 300
	case 80:
		k = -140
		nw = 200
	case 85:
		k = -90
		nw = 200
	case 90:
		k = -40
		nw = 200
	case 93:
		k = -10
		nw = 200
	case 96:
		k = 20
		nw = 200
	case 99:
		k = 50
		nw = 200
	case 102:
		k = 80
		nw = 200
	case 110:
		k = 160
		nw = 200
	case 120:
		k = 260
		nw = 200
	}

	rw = k + MapGameNow[gameId]

	return
}
