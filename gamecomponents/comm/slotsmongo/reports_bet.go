package slotsmongo

import (
	"context"
	"game/comm/db"
	"game/comm/ut"
	"strconv"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

func addReport(game string, appId string, pid int64, bet, win int64, completed bool, isBuy bool, bigReward int64) {
	dateStr := time.Now().Format("20060102")
	incDoc := db.D(
		"betamount", bet,
		"winamount", win,
		"spincount", lo.Ternary(completed, 1, 0),
		"bigreward", bigReward,
	)

	if isBuy {
		// incDoc["buybetamount"] = bet
		// incDoc["buywinamount"] = win
		db.AppendE(&incDoc, "buybetamount", bet)
		db.AppendE(&incDoc, "buywinamount", win)
	}

	if completed && IsFirstVisitToday(ut.JoinStr(':', "spin", game, strconv.Itoa(int(pid)))) {
		// incDoc["spinplrcount"] = 1
		db.AppendE(&incDoc, "spinplrcount", 1)
	}

	//if IsFirstVisitToday(ut.JoinStr(':', "enter", game, strconv.Itoa(int(pid)))) {
	//	// incDoc["enterplrcount"] = 1
	//	db.AppendE(&incDoc, "enterplrcount", 1)
	//}

	// addBetDailyReport(dateStr, game, appId, pid, bet, win, hasBonusGame, roundCount, isBuy)
	addBetDailyReport(dateStr, game, appId, incDoc)
	addOperatorReport(dateStr, appId, incDoc)
	addPayerReport(dateStr, appId, pid, incDoc)
	addBetSummaryReport(game, appId, incDoc)
}

func UpdateEnterPlrCount(game string, appId string, pid int64) {
	// if !IsFirstVisitToday(ut.JoinStr(':', "enter", game, strconv.Itoa(int(pid)))) {
	// 	return
	// }

	// dateStr := time.Now().Format("20060102")
	// incDoc := db.D("enterplrcount", 1)

	// addBetDailyReport(dateStr, game, appId, incDoc)
	// addOperatorReport(dateStr, appId, incDoc)
	// //addPayerReport(dateStr, appId, pid, incDoc)
	// addBetSummaryReport(game, appId, incDoc)
}

//func AddReport1(dateStr string, game string, appId string, pid int64, bet, win int64, activePlayer, loginPlayer, isBuy bool) {
//	incDoc := db.D(
//		"betamount", bet,
//		"winamount", win,
//		"spincount", 1,
//	)
//
//	if isBuy {
//		// incDoc["buybetamount"] = bet
//		// incDoc["buywinamount"] = win
//		db.AppendE(&incDoc, "buybetamount", bet)
//		db.AppendE(&incDoc, "buywinamount", win)
//	}
//
//	if activePlayer {
//		// incDoc["spinplrcount"] = 1
//		db.AppendE(&incDoc, "spinplrcount", 1)
//	}
//
//	if loginPlayer {
//		// incDoc["enterplrcount"] = 1
//		db.AppendE(&incDoc, "enterplrcount", 1)
//	}
//
//	// addBetDailyReport(dateStr, game, appId, pid, bet, win, hasBonusGame, roundCount, isBuy)
//	addBetDailyReport(dateStr, game, appId, incDoc)
//	addOperatorReport(dateStr, appId, incDoc)
//	addPayerReport(dateStr, appId, pid, incDoc)
//	addBetSummaryReport(game, appId, incDoc)
//}

// /////////////////////////////////////////////////////
type BetDailyDoc struct {
	ID             string `json:"-" bson:"_id"`
	BetAmount      int64  `bson:"betamount"`      // 支出
	WinAmount      int64  `bson:"winamount"`      // 赚取
	BuyBetAmount   int64  `bson:"buybetamount"`   // buy支出
	BuyWinAmount   int64  `bson:"buywinamount"`   // buy赚取
	ExtraBetAmount int64  `bson:"extrabetamount"` // extra支出
	ExtraWinAmount int64  `bson:"extrawinamount"` // extra赚取
	BigReward      int64  `bson:"bigreward"`      // 累计爆大奖
	SpinCount      int    `bson:"spincount"`      // 转动次数
	SpinPlrCount   int    `bson:"spinplrcount"`   // 参与转动的玩家数量
	EnterPlrCount  int    `bson:"enterplrcount"`  // 进入游戏的玩家数量

	Date  string `bson:"date"`
	Game  string `bson:"game"`
	AppID string `bson:"appid"`
}

func addBetDailyReport(dateStr, game string, appId string, incDoc bson.D /*, pid int64, bet, win int64, hasBonusGame bool, roundCount int, isBuy bool*/) {
	coll := ReportsCollection("BetDailyReport")
	id := ut.JoinStr(':', dateStr, game, appId)

	update := db.D(
		"$set", db.D(
			"date", dateStr,
			"game", game,
			"appid", appId,
		),
		"$inc", incDoc,
	)
	coll.UpdateByID(context.TODO(), id, update, db.UpsertOpt())
	/*
		doc := &BetDailyDoc{
			ID:    id,
			Date:  dateStr,
			Game:  game,
			AppID: appId,
		}
		coll.InsertOne(context.TODO(), doc)

		incDoc := bson.M{
			"betamount": bet,
			"winamount": win,
			"spincount": 1,
		}

		if isBuy {
			incDoc["buybetamount"] = bet
			incDoc["buywinamount"] = win
		}

		if IsFirstVisitToday(ut.JoinStr(':', "spin", game, strconv.Itoa(int(pid)))) {
			incDoc["spinplrcount"] = 1
		}

		if IsFirstVisitToday(ut.JoinStr(':', "enter", game, strconv.Itoa(int(pid)))) {
			incDoc["enterplrcount"] = 1
		}

		_, err := coll.UpdateByID(context.TODO(), id, db.D("$inc", incDoc))
		if err != nil {
			slog.Error("UpdateByID Error", "err", err, "incDoc", incDoc)
		}
	*/
}

// /////////////////////////////////////////////////////
type BetSummaryDoc struct {
	ID           string `json:"-" bson:"_id"`
	BetAmount    int    `bson:"betamount"`    // 支出
	WinAmount    int    `bson:"winamount"`    // 赚取
	BuyBetAmount int    `bson:"buybetamount"` // buy支出
	BuyWinAmount int    `bson:"buywinamount"` // buy赚取
	SpinCount    int    `bson:"spincount"`    // 转动次数

	Game  string `bson:"game"`
	AppID string `bson:"appid"`
}

func addBetSummaryReport(game, appId string, incDoc bson.D) {
	coll := ReportsCollection("BetSummaryReport")
	id := ut.JoinStr(':', game, appId)

	update := db.D(
		"$set", db.D(
			"game", game,
			"appid", appId,
		),
		"$inc", incDoc,
	)
	coll.UpdateByID(context.TODO(), id, update, db.UpsertOpt())
}
