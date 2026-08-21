package handlers

import (
	"context"
	"game/comm"
	"game/comm/mq"
	"game/comm/slotsmongo"
	"game/service/gamecenter/internal/operator"
	"log/slog"
	"net/url"
)

func RegSub() {
	mq.JsonNC.Subscribe("/games/betlog", onAddBetLog)
	mq.JsonNC.Subscribe("/player/status", statusChange)

	// mq.JsonNC.Subscribe("/games/lottery/openPrize", func(lotteryAlert *slotsmongo.LotteryAlert) {
	// 	alert.onOpenLottery(lotteryAlert)
	// })
}

func onAddBetLog(betlog *slotsmongo.DocBetLog) {
	app := operator.AppMgr.GetApp(betlog.AppID)
	if app == nil {
		slog.Error("onAddBetLog", "appID", betlog.AppID, "betlog", betlog, "error", "app is nil")
		return
	}

	if !app.PublishHistory {
		return
	}

	var err error
	defer func() {
		if err != nil {
			slog.Error("onAddBetLog", "appID", betlog.AppID, "error", err.Error(), "betlog", betlog)
		}
	}()

	apiUrl, err := url.JoinPath(app.Address, "/History/Spin")
	if err != nil {
		return
	}

	req := map[string]any{
		"ID":         betlog.ID,
		"Pid":        betlog.Pid,
		"UserID":     betlog.UserID,
		"GameID":     betlog.GameID,
		"Bet":        goldtype(betlog.Bet),
		"Win":        goldtype(betlog.Win),
		"InsertTime": betlog.InsertTime,
		"AppID":      betlog.AppID,
		"Balance":    goldtype(betlog.Balance),
		"WinLose":    goldtype(betlog.WinLose),
		"Grade":      betlog.Grade,
		// "Round":      cmp.Or(betlog.PGBetID, betlog.ID.Hex()),
		"RoundID": betlog.RoundID,
	}

	// ut.PrintJson(req)

	err = comm.PostJsonCode(context.TODO(), apiUrl, req, nil, nil)
}

func statusChange(player operator.PlayerStatus) {
	app := operator.AppMgr.GetApp(player.AppID)
	err := operator.AppMgr.UpdateCacheUserStatus(app, player)
	if err != nil {
		slog.Error("statusChange", "player", player.UserID, "error", err.Error(), "status", player.Status)
	}
}
