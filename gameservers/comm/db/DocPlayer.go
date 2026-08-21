package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"serve/comm/ut"

	"github.com/samber/lo"
)

type DocPlayer struct {
	PID         int64  `bson:"_id"`
	Uid         string `bson:"Uid"`         // 外部id
	AppID       string `bson:"AppID"`       // 所属产品
	CurrencyKey string `bson:"CurrencyKey"` // 玩家货币key
	Language    string `bson:"-"`           // 语言 th,en,cn
	//CurrentRtp  int    `bson:"CurrentRtp"`  // RTP
	Status int64 `bson:"Status"` // 玩家状态

	RestrictionsStatus   int64   `bson:"RestrictionsStatus"`   //是否开启最大中奖钱数倍数
	Win                  int64   `bson:"Win"`                  //当前钱数
	Multi                float64 `bson:"Multi"`                //当前倍数
	RestrictionsMaxWin   int64   `bson:"RestrictionsMaxWin"`   //限制最大钱数
	RestrictionsMaxMulti float64 `bson:"RestrictionsMaxMulti"` //限制最大倍数

	SpinCountOfThisEnter int                `bson:"spincountofthisenter"`
	SpinCount            int                `bson:"spincount"`
	MoniterCountNode     int                `bson:"monitercountnode"` //监控数量节点
	MoniterRTPNode       int64              `bson:"moniterrtpnode"`   //监控RTP节点
	BonusGameCount       int                `bson:"bonusgamecount"`   //小游戏数量
	BetHistory           []int64            `bson:"bethistory"`
	BetAmount            int64              `bson:"betamount"`
	WinAmount            int64              `bson:"winamount"`
	BetHistLastID        primitive.ObjectID `bson:"bethistlastid"`
	LastBetId            string             `bson:"lastbetid"`

	//二期改动，直接查表，不做增加处理
	TargetRTP          int   `bson:"-"` //限制目标
	RelieveRTP         int   `bson:"-"` //释放限制
	RewardPercent      int   `bson:"-"`
	NoAwardPercent     int   `bson:"-"`
	BuyMinAwardPercent int   `bson:"-"`
	PersonWinMaxMult   int64 `bson:"-"`
	PersonWinMaxScore  int   `bson:"-"`
}

func (p *DocPlayer) GetCurrencyOrTHB() string {
	//return "" //取消关联币种  每个币种一个桶
	if p.CurrencyKey == "" {
		return "THB"
	}

	return p.CurrencyKey
}

func (p *DocPlayer) CaclAvgBet() int64 {
	n := len(p.BetHistory)
	if n == 0 {
		return 0
	}

	sum := lo.Sum(p.BetHistory)
	avg := sum / int64(n)
	return avg
}

func (p *DocPlayer) appendBetHistory(bet int64) {
	p.BetHistory = append(p.BetHistory, bet)
	p.BetHistory = ut.MaxLen(p.BetHistory, 200)
}

func (player *DocPlayer) OnSpinFinish(bet, win int64, isBuy, hasBonusGame bool, buyMul int) {
	//源代码购买游戏未累计 购买倍数 pg 75 jl? pp?
	if isBuy {
		bet *= int64(buyMul)
	}
	player.BetAmount += bet
	player.WinAmount += win

	// if bet != 0 {
	if !isBuy {
		player.appendBetHistory(bet)
	}
	player.SpinCount++
	// player.SpinCountOfThisEnter++
	if hasBonusGame {
		player.BonusGameCount++
	}
	// }
}

var (
	plrprojection = D(
		"_id", 1,
		"Uid", 1,
		"AppID", 1,
		"Status", 1,

		"RestrictionsStatus", 1,
		"Win", 1,
		"Multi", 1,
		"RestrictionsMaxWin", 1,
		"RestrictionsMaxMulti", 1,
	)
)

func getDocPlayer(pid int64) (plr *DocPlayer, err error) {

	coll := Collection2("game", "Players")

	var doc DocPlayer
	err = coll.FindOne(context.TODO(), ID(pid), options.FindOne().SetProjection(plrprojection)).Decode(&doc)
	if err != nil {
		return
	}
	plr = &doc

	type currencyItem struct {
		CurrencyKey string `bson:"CurrencyKey"`
	}
	var item *currencyItem
	coll = Collection2("GameAdmin", "AdminOperator")
	coll.FindOne(context.TODO(), bson.M{"AppID": doc.AppID}, options.FindOne().SetProjection(D("CurrencyKey", 1))).Decode(&item)
	if item != nil {
		plr.CurrencyKey = item.CurrencyKey
	}
	if doc.CurrencyKey == "" {
		doc.CurrencyKey = "THB"
	}
	return
}

func (p *DocPlayer) GetWinLose() int64 {
	coll := Collection2("game", "Players")

	project := D(
		"Bet", 1,
		"Win", 1,
	)
	var doc struct {
		Bet int64 `bson:"Bet"`
		Win int64 `bson:"Win"`
	}
	err := coll.FindOne(context.TODO(), ID(p.PID), options.FindOne().SetProjection(project)).Decode(&doc)
	if err != nil {
		return 0
	}
	return doc.Win - doc.Bet
}

// 玩家RTP控制
type PlayerRTPControlModel struct {
	ID                 int64     `bson:"_id"`
	AppID              string    `bson:"AppID"`  //运营商  关联运营商
	GameID             string    `bson:"GameID"` //游戏编号 关联游戏
	CreateAt           time.Time `bson:"CreateAt"`
	Pid                int64     `bson:"Pid"`                //用户账号
	GameName           string    `bson:"GameName"`           //游戏名称
	PlayerRTP          int64     `bson:"PlayerRTP"`          //玩家历史RTP RTP控制之前的
	ControlRTP         float64   `bson:"ControlRTP"`         //账号RTP
	RewardPercent      int64     `bson:"RewardPercent"`      //账号RTP
	NoAwardPercent     int64     `bson:"NoAwardPercent"`     //账号RTP
	AutoRemoveRTP      float64   `bson:"AutoRemoveRTP"`      //自动解除RTP
	AutoRewardPercent  int64     `bson:"AutoRewardPercent"`  //自动解除RTP
	AutoNoAwardPercent int64     `bson:"AutoNoAwardPercent"` //自动解除RTP
	Status             int       `bson:"Status"`             //状态
	BuyMinAwardPercent int64     `bson:"BuyMinAwardPercent"` //限制购买的千分比
	PersonWinMaxMult   int       `bson:"PersonWinMaxMult"`   // 玩家最大赢分倍数
	PersonWinMaxScore  int64     `bson:"PersonWinMaxScore"`  // 玩家最大赢分
}

func GetDocPlayer(pid int64) (plr *DocPlayer, err error) {

	coll := Collection2("game", "Players")

	var doc DocPlayer
	err = coll.FindOne(context.TODO(), ID(pid), options.FindOne().SetProjection(plrprojection)).Decode(&doc)
	if err != nil {
		return
	}
	plr = &doc

	type currencyItem struct {
		CurrencyKey string `bson:"CurrencyKey"`
	}
	var item *currencyItem
	coll = Collection2("GameAdmin", "AdminOperator")
	coll.FindOne(context.TODO(), bson.M{"AppID": doc.AppID}, options.FindOne().SetProjection(D("CurrencyKey", 1))).Decode(&item)
	if item != nil {
		plr.CurrencyKey = item.CurrencyKey
	}
	if doc.CurrencyKey == "" {
		doc.CurrencyKey = "THB"
	}
	return
}
