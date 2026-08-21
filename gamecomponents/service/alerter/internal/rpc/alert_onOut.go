package rpc

import (
	"context"
	"fmt"
	"game/comm/db"
	"game/comm/slotsmongo"
	"game/duck/lazy"
	"game/service/alerter/alerterlib"
	"game/service/alerter/internal/gamedata"
	"time"
)

func onTransferOut(ps *alerterlib.OnTransferOutPs) {
	//currency是货币
	currency := getCurrency(ps.AppID)
	item := lazy.GetCurrencyItem(currency)
	//
	//multi := item.Multi
	//if multi <= 0 {
	multi := 1.0
	//}

	amount := ps.Amount / multi

	set := gamedata.Settings
	transferOutMonry := set.GetTransferOutMoney(ps.AppID)

	if amount < transferOutMonry {
		return
	}
	lang := gamedata.GetLanguage(ps.AppID)
	msg := fmt.Sprintf(gamedata.MapLanguage[lang].LargeAmountTransfer+"\n"+
		gamedata.MapLanguage[lang].PlayerID+"\n"+
		gamedata.MapLanguage[lang].MerchantPlayerID+"\n"+
		gamedata.MapLanguage[lang].Merchant+"\n"+
		gamedata.MapLanguage[lang].OrderNumber+"\n"+
		gamedata.MapLanguage[lang].TransFerAmount+"\n"+
		gamedata.MapLanguage[lang].LatestBalance,
		ps.Pid, ps.UserID, ps.AppID, ps.ID,
		fmtMoney(item.Symbol, ps.Amount),
		fmtMoney(item.Symbol, ps.Balance),
	)
	go slotsmongo.SendTelegramChat(msg)
	go slotsmongo.SendOperatorTelegramChat(ps.AppID, msg)
	coll := db.Collection2("reports", "alert")
	coll.InsertOne(context.TODO(), db.D(
		"pid", ps.Pid,
		"userid", ps.UserID,
		"appId", ps.AppID,
		"orderId", ps.ID,
		"amount", fmtMoney(item.Symbol, amount),
		"balance", fmtMoney(item.Symbol, ps.Balance),
		"createTime", time.Now().UTC(),
	))
}
