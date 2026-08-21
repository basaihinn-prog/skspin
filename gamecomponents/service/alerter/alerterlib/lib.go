package alerterlib

import (
	"game/comm/mq"
)

type OnTransferOutPs struct {
	ID         string `bson:"_id"`
	Pid        int64
	AppID      string
	UserID     string
	Amount     float64
	Balance    float64
	InsertTime string `bson:"InsertTime"`
}

type BalanceAlertPs struct {
	ID              string `bson:"_id"`
	AppID           string
	CooperationType int
	WalletType      int
	MerchantBalance float64
}

func OnTransferOut(ps *OnTransferOutPs) error {
	return mq.JsonNC.Publish("/alerter/onTransferOut", ps)
}
