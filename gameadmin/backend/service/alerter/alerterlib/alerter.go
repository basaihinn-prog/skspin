package alerterlib

import "go.mongodb.org/mongo-driver/bson/primitive"

type OnTransferOutPs struct {
	ID      primitive.ObjectID
	Pid     int64
	UserID  string
	AppID   string
	Amount  float64
	Balance float64
}

func OnTransferOut(ps *OnTransferOutPs) {
	// Stub implementation - alerter service not available
	// This is called when a player transfers out funds
}
