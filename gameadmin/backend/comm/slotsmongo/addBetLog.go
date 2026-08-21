package slotsmongo

import (
	"context"
	"encoding/json"
	"fmt"
	"game/comm/db"
	"game/comm/mq"
	"game/duck/lazy"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocBetLog struct {
	ID         primitive.ObjectID `bson:"_id"`
	Pid        int64              `bson:"Pid"`    // 内部玩家ID
	UserID     string             `bson:"UserID"` // 外部玩家ID
	GameID     string             `bson:"GameID"` // 游戏ID
	GameType   int                `bson:"GameType"`
	RoundID    string             `bson:"RoundID"`    // 对局ID，方便查询
	Bet        int64              `bson:"Bet"`        // 下注
	Win        int64              `bson:"Win"`        // 单次消除输赢
	Comment    string             `bson:"Comment"`    // 注释
	InsertTime time.Time          `bson:"InsertTime"` // 数据插入时间
	AppID      string             `bson:"AppID"`      //
	Balance    int64              `bson:"Balance"`    // 余额
	// RoundCount      int                `bson:"RoundCount"`
	Completed       bool   `bson:"Completed"`    // 最后一次spin
	TotalWinLoss    int64  `bson:"TotalWinLoss"` // 当前局结算总输赢, 只需要 Completed:true 才使用
	WinLose         int64  `bson:"WinLose"`      // 玩家输赢金额
	Grade           int    `bson:"Grade"`
	HitBigReward    int64  `bson:"HitBigReward"` // 奖池奖励
	SpinDetailsJson string `bson:"SpinDetailsJson"`
	PGBetID         string `bson:"PGBetID"`
	LogType         int    `bson:"LogType"`        // 0 电子投注/结算(游戏), 1 转入, 2 转出, 3 admin设置, 4 彩票投注, 5 彩票结算, 6 百人投注/结算
	TransferAmount  int64  `bson:"TransferAmount"` // 转移金额
	Frb             bool   `bson:"Frb"`            // 是否免费
	BonusCode       string `bson:"BonusCode"`      // 免费代码
}

func InsertBetLog(betlog *DocBetLog) {
	now := time.Now()
	if betlog.InsertTime.IsZero() {
		betlog.InsertTime = now
	}
	if betlog.ID.IsZero() {
		betlog.ID = primitive.NewObjectIDFromTimestamp(now)
	}
	betlog.WinLose = betlog.Win - betlog.Bet
	// coll := ReportsCollection("BetLog")
	// coll.InsertOne(context.TODO(), betlog)

	appColl := db.Collection2("betlog", betlog.AppID)
	appColl.InsertOne(context.TODO(), betlog)

	// oldspin := betlog.SpinDetailsJson
	// betlog.SpinDetailsJson = ""
	// GetTTLLogColl("betlog", 180).InsertOne(context.TODO(), betlog)
	// betlog.SpinDetailsJson = oldspin
}

func GetBetLogColl(appid string) *mongo.Collection {
	appColl := db.Collection2("betlog", appid)
	return appColl
}

func addBetLog(betlog *DocBetLog) {
	InsertBetLog(betlog)
	betlog.SpinDetailsJson = ""
	// mq.Invoke("/alerter/betlog/add", betlog, nil)
	if !strings.Contains(betlog.GameID, "lottery") {
		mq.JsonNC.Publish("/games/betlog", betlog)
	}
}

// 彩票开奖之后添加alert消息
func AddLotteryAlert(lotteryAlert *LotteryAlert) {
	mq.JsonNC.Publish("/games/lottery/openPrize", lotteryAlert)
}

type LotteryAlert struct {
	OpenDay        string
	PeriodSoldGold int64
	CompanyWin_1   int64
	PeriodProfits  int64

	PeriodSoldGuessGold int64
	CompanyWin_2        int64
	PeriodGuessProfits  int64

	CompanyWin_total int64
}

type AddBetLogParams struct {
	ID           primitive.ObjectID
	Pid          int64
	Bet          int64  // 下注
	Win          int64  // 输赢
	Balance      int64  // 余额
	RoundID      string // 对局ID，方便查询, **废弃的字段, 兼容性考虑暂时不删除**
	Completed    bool   // 最后一次spin
	TotalWinLoss int64  // 当前局结算总输赢, 只需要 Completed:true 才使用
	IsBuy        bool
	Grade        int
	Details      any
	BigReward    int64 // 奖池奖励
	// PGBetID      primitive.ObjectID
	GameType int
	// LogType  int // 0 游戏, 1 转入, 2 转出
	Comment string // 外部引用
}

// 添加记录
// func AddBetLog(pid int64, bet, win, balance int64, roundId string, completed bool, totalWinLoss int64,
//
//	isBuy bool, grade int, details any, bigReward int64, gameType int) (id primitive.ObjectID) {
func AddBetLog(ps *AddBetLogParams) (id primitive.ObjectID) {
	var (
		pid          = ps.Pid
		bet          = ps.Bet
		win          = ps.Win
		balance      = ps.Balance
		roundId      = ps.RoundID
		completed    = ps.Completed
		totalWinLoss = ps.TotalWinLoss
		isBuy        = ps.IsBuy
		// pgbetid      = ps.PGBetID
		bigReward = ps.BigReward
		gameType  = ps.GameType
		grade     = ps.Grade
		details   = ps.Details
	)

	id = ps.ID
	if id.IsZero() {
		id = primitive.NewObjectID()
	}

	appId, uid, err := GetPlayerInfo(pid)
	if err != nil {
		slog.Error("addRecord Fail!", "pid", pid, "error", err)
		return
	}

	jsonbuf, _ := json.Marshal(details)

	// var pgbetid string
	// if ut.IsPGGame(lazy.ServiceName) {
	// 	pgbetid = details.(string)
	// }

	playerDoc := []any{
		"BigReward", bigReward,
		"Bet", bet,
		"Win", win,
	}

	if isBuy {
		playerDoc = append(playerDoc,
			fmt.Sprintf("TypeInfo.%d.BuyGame", gameType), 1,
			fmt.Sprintf("TypeInfo.%d.BuyGameBet", gameType), bet,
			fmt.Sprintf("TypeInfo.%d.BuyGameWin", gameType), win)
	} else {
		playerDoc = append(playerDoc,
			fmt.Sprintf("TypeInfo.%d.Bet", gameType), bet,
			fmt.Sprintf("TypeInfo.%d.SpinCnt", gameType), 1,
			fmt.Sprintf("TypeInfo.%d.Win", gameType), win)
	}

	db.Collection2("game", "Players").UpdateByID(context.TODO(), pid, db.D(
		"$inc", db.D(playerDoc...),
	))
	addReport(lazy.ServiceName, appId, pid, bet, win, completed, isBuy, bigReward)

	// id = primitive.NewObjectID()
	addBetLog(&DocBetLog{
		ID:       id,
		Pid:      pid,
		UserID:   uid,
		GameID:   lazy.ServiceName,
		GameType: gameType,
		RoundID:  roundId,
		Bet:      bet,
		Win:      win,
		Comment:  ps.Comment,
		AppID:    appId,
		Balance:  balance,
		// RoundCount:      roundCount,
		Completed:       completed,
		TotalWinLoss:    totalWinLoss,
		Grade:           grade,
		HitBigReward:    bigReward,
		SpinDetailsJson: string(jsonbuf),
		// PGBetID:         pgbetid,
	})
	return
}

// ps *AddBetLogParams
// 添加记录
// func AddCaiPiaoBetLog(pid int64, bet, win, balance int64, roundId string, hasBonusGame bool, completed bool,
// 	isBuy bool, grade int, details any, bigReward int64, gameType int) (id primitive.ObjectID) {

// period表示是1期还是2期 period为0的时候是1期 period为1的时候是2期
func AddCaiPiaoBetLog(ps *AddBetLogParams, period int) (id primitive.ObjectID) {
	var (
		pid          = ps.Pid
		bet          = ps.Bet
		win          = ps.Win
		balance      = ps.Balance
		roundId      = ps.RoundID
		completed    = ps.Completed
		totalWinLoss = ps.TotalWinLoss
		// isBuy        = ps.IsBuy
		// pgbetid      = ps.PGBetID
		bigReward = ps.BigReward
		gameType  = ps.GameType
		grade     = ps.Grade
		details   = ps.Details
	)

	appId, uid, err := GetPlayerInfo(pid)
	if err != nil {
		slog.Error("AddCaiPiaoBetLog Fail!", "pid", pid, "error", err)
		return
	}

	jsonbuf, _ := json.Marshal(details)

	id = primitive.NewObjectID()
	gameID := "lottery"
	if period == 1 {
		gameID = fmt.Sprintf("%s_%v", gameID, period)
	}
	addBetLog(&DocBetLog{
		ID:       id,
		Pid:      pid,
		UserID:   uid,
		GameID:   gameID,
		GameType: gameType,
		RoundID:  roundId,
		Bet:      bet,
		Win:      win,
		Comment:  id.Hex(),
		AppID:    appId,
		Balance:  balance,
		// RoundCount:      roundCount,
		Grade:           grade,
		HitBigReward:    bigReward,
		SpinDetailsJson: string(jsonbuf),
		Completed:       completed,
		TotalWinLoss:    totalWinLoss,
	})
	return
}
