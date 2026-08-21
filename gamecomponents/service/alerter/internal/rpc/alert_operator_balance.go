package rpc

import (
	"context"
	"fmt"
	"game/comm/db"
	"game/comm/slotsmongo"
	"game/service/alerter/alerterlib"
	"game/service/alerter/internal/gamedata"
	"time"
)

func operatorBalanceAlert(ps *alerterlib.BalanceAlertPs) {
	//currency是货币
	currency := getCurrency(ps.AppID)
	//item := lazy.GetCurrencyItem(currency)
	lang := gamedata.GetLanguage(ps.AppID)
	var cooperationType, walletType string
	switch ps.WalletType {
	case 1:
		walletType = gamedata.MapLanguage[lang].TransferWallet
	case 2:
		walletType = gamedata.MapLanguage[lang].SingleWallet
	}
	switch ps.CooperationType {
	case 1:
		cooperationType = gamedata.MapLanguage[lang].RevenueShare
	case 2:
		cooperationType = gamedata.MapLanguage[lang].CashFlowShare
	}
	msg := fmt.Sprintf(gamedata.MapLanguage[lang].BalanceAlert+"\n\n"+
		gamedata.MapLanguage[lang].Merchant+"\n"+
		gamedata.MapLanguage[lang].Currency+"\n"+
		gamedata.MapLanguage[lang].CooperationType+"\n"+
		gamedata.MapLanguage[lang].WalletType+"\n"+
		gamedata.MapLanguage[lang].MerchantBalance+"\n\n"+
		gamedata.MapLanguage[lang].BalanceAlertInfo,
		ps.AppID, currency, cooperationType,
		walletType, ps.MerchantBalance,
	)
	go slotsmongo.SendTelegramChat(msg)
	go slotsmongo.SendOperatorTelegramChat(ps.AppID, msg)
	go slotsmongo.SendOperatorBalanceTelegramChat(ps.AppID, msg)
	coll := db.Collection2("reports", "alert")
	coll.InsertOne(context.TODO(), db.D(
		"appId", ps.AppID,
		"currency", currency,
		"cooperationType", ps.CooperationType,
		"walletType", ps.WalletType,
		"merchantBalance", ps.MerchantBalance,
		"createTime", time.Now().UTC(),
	))
}
