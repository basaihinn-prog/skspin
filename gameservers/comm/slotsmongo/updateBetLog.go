package slotsmongo

import (
	"log/slog"
	"serve/comm/mq"
	"time"
)

func UpdateBetLogAviator(ps *AddBetLogParamsAviator) /* (id primitive.ObjectID)*/ {
	var (
		pid     = ps.Pid
		bet     = ps.Bet
		win     = ps.Win
		balance = ps.Balance
		roundId = ps.RoundId
		appId   = ps.AppID
		uid     = ps.Uid
	)
	//if appId == "" || uid == "" {
	//	var err error
	//	appId, uid, err = GetPlayerInfo(pid)
	//	if err != nil {
	//		slog.Error("GetPlayerInfo Fail!", "pid", pid, "error", err)
	//		return
	//	}
	//}
	//go addReport(serviceName, appId, pid, bet, win, completed, isBuy, extra, bigReward)
	updateBetLogAviator(&DocBetLogAviator{ //TODO
		ID:                 ps.ID,
		Pid:                pid,
		AppID:              appId,
		Uid:                uid,
		UserName:           ps.UserName,
		CurrencyKey:        ps.CurrencyKey,
		RoomId:             ps.RoomId,
		Bet:                bet,
		Win:                win,
		Balance:            balance,
		CashOutDate:        ps.CashOutDate,
		Payout:             ps.Payout,
		Profit:             ps.Profit,
		RoundBetId:         ps.RoundBetId,
		RoundId:            roundId,
		MaxMultiplier:      ps.MaxMultiplier,
		Frb:                ps.Frb,
		GameID:             ps.GameId,
		RoundMaxMultiplier: ps.RoundMaxMultiplier,
		ManufacturerName:   ps.ManufacturerName,
		FinishType:         ps.FinishType,
	})
	return
}

func updateBetLogAviator(betlog *DocBetLogAviator) {
	now := time.Now()
	if betlog.InsertTime.IsZero() {
		betlog.InsertTime = now
	}
	if betlog.ID == "" {
		slog.Error("betlog param has not Id")
		return
	}
	// coll := ReportsCollection("BetLog")
	// coll.InsertOne(context.TODO(), betlog)

	//appColl := db.Collection2("betlog", betlog.AppID)
	//appColl.InsertOne(context.TODO(), betlog)

	// mq.Invoke("/alerter/betlog/add", betlog, nil)

	err := mq.JsonNC.Publish("/games/minibetlogUpdate", betlog)
	if err != nil {
		slog.Error("更新历史记录失败", "betlog: ", betlog)
		slog.Error("updateBetLogAviator Fail!", "err", err)
	}

	//mq.PublishMsg("/games/betlog", betlog)
}
