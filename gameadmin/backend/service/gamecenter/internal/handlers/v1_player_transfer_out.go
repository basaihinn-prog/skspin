package handlers

import (
	"context"
	"errors"
	"game/comm"
	"game/comm/db"
	"game/comm/mux"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/service/alerter/alerterlib"
	"game/service/gamecenter/internal/operator"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/player/transferOutAll",
		Handler:      v1_player_transfer_out_all,
		Desc:         "玩家转出所有余额",
		Kind:         "api/v1",
		ParamsSample: v1PlayerTransferOutAllPs{"user_id", "trace_id 4-36"},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type v1PlayerTransferOutAllPs struct {
	UserID  string
	TraceId string
}

type v1PlayerTransferOutAllRet struct {
	Amount float64
}

func v1_player_transfer_out_all(app *operator.MemApp, ps v1PlayerTransferOutAllPs, ret *v1PlayerTransferOutAllRet) (err error) {
	if app.WalletMode != comm.WalletModeTransfer {
		err = errors.New("Not transfer Mode")
		return
	}
	memplr, err := operator.AppMgr.GetPlr2(app, ps.UserID)
	if err != nil {
		return
	}

	now := time.Now()
	order := &TransferOrder{
		ID:         primitive.NewObjectIDFromTimestamp(now),
		TraceId:    ps.TraceId,
		CreateTime: now,
		UserID:     memplr.Uid,
		AppID:      memplr.AppID,
		Pid:        memplr.Pid,
		Action:     "outall",
	}
	coll := db.Collection2("game", "orders")
	_, err = coll.InsertOne(context.TODO(), order)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			err = errors.New("Order already exists")
		}
		return
	}

	amount, err := memplr.TransferOutAll()
	// if err != nil {
	// coll.UpdateByID(context.TODO(), order.TraceId, db.D("$set", db.D("Error", err.Error())))
	// }

	update := db.D(
		"Completed", true,
	)

	if err != nil {
		db.AppendE(&update, "Error", err.Error())
	} else {
		db.AppendE(&update, "Amount", amount)
		ret.Amount = amount
	}
	coll.UpdateByID(context.TODO(), order.ID, db.D("$set", update))

	if err == nil {
		slotsmongo.InsertBetLog(&slotsmongo.DocBetLog{
			Pid:    memplr.Pid,
			UserID: memplr.Uid,
			AppID:  memplr.AppID,
			// Completed: true,
			LogType:        2,
			TransferAmount: ut.Money2Gold(amount),
			Balance:        0,
		})

		alerterlib.OnTransferOut(&alerterlib.OnTransferOutPs{
			ID:      order.ID,
			Pid:     memplr.Pid,
			UserID:  memplr.Uid,
			AppID:   memplr.AppID,
			Amount:  amount,
			Balance: 0,
		})
	}

	operator.InsertModifyLog(&operator.DocModifyGoldLog{
		ID:      order.ID,
		Change:  -ut.Money2Gold(amount),
		Balance: 0,
		Error:   ut.ErrString(err),
		Action:  "outall",
		ReqData: order,
	}, memplr)

	return
}
