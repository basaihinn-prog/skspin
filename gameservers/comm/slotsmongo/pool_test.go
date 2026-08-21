package slotsmongo

import (
	"context"
	"testing"

	"serve/comm/ut"
)

func TestPool(t *testing.T) {
	/*
		db.DialToMongo("mongodb://127.0.0.1:27017/YingCaiShen")

		gold, err := GetSelfSlotsPool(123456)
		slog.Info("GetSelfSlotsPool", "gold", gold, "error", err)

		gold, err = IncSelfSlotsPool(123456, 100)
		slog.Info("IncSelfSlotsPool", "error", err, "gold", gold)

		gold, err = GetSelfSlotsPool(123456)
		slog.Info("GetSelfSlotsPool", "gold", gold, "error", err)

		gold, err = IncSelfSlotsPool(123456, 100)
		slog.Info("IncSelfSlotsPool", "error", err, "gold", gold)

		gold, err = GetSelfSlotsPool(123456)
		slog.Info("GetSelfSlotsPool", "gold", gold, "error", err)

		SetSelfSlotsPool(223456, 1)
		SetSelfSlotsPool(223456, 2)
	*/
}

func TestAddBet(t *testing.T) {
	/*
		db.DialToMongo("mongodb://127.0.0.1:27017/YingCaiShen")

		/*
			AddBet(123456, 1)
			AddBet(123456, 2)
			AddBet(123456, 3)
			AddBet(123456, 4)
			AddBet(123456, 5)
			AddBet(123456, 6)
			AddBet(123456, 7)
			AddBet(123456, 8)
			AddBet(123456, 9)
			AddBet(123456, 10)
	*/
	AvgBet(12345)
}

func TestAddBetDailyReport(t *testing.T) {
	// db.DialToMongo("mongodb://127.0.0.1:27017/YingCaiShen")

	// AddBetDailyReport("YingCaiShen", )
	coll := ReportsCollection("BetDailyReport")
	// dateStr := time.Now().Format("20060102")
	dateStr := "20231012" // time.Now().Format("20060102")
	game := "XingYunXiang"
	appId := "fake"
	id := ut.JoinStr(':', dateStr, game, appId)
	doc := &BetDailyDoc{
		ID:    id,
		Date:  dateStr,
		Game:  game,
		AppID: appId,

		BetAmount:     1900000,
		WinAmount:     180000,
		SpinCount:     138,
		SpinPlrCount:  29,
		EnterPlrCount: 31,
	}
	coll.InsertOne(context.TODO(), doc)
}

func TestVisit(t *testing.T) {
	// db.DialToMongo("mongodb://127.0.0.1:27017/YingCaiShen")
	// assert.True(t, IsFirstVisitToday("abc"))
	// assert.False(t, IsFirstVisitToday("abc1"))
	// assert.True(t, IsFirstVisitToday("abc1"))
}
