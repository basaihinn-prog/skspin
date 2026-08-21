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
	"game/service/gamecenter/internal/gamedata"
	"game/service/gamecenter/internal/operator"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/player/transferIn",
		Handler:      v1_player_transfer_in,
		Desc:         "玩家充值",
		Kind:         "api/v1",
		ParamsSample: v1PlayerTransferInPs{"user_id", 123.45, "trace_id 4-36"},
		Class:        "operator",
		GetArg0:      getArg0,
	})

	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/player/transferOut",
		Handler:      v1_player_transfer_out,
		Desc:         "玩家提现",
		Kind:         "api/v1",
		ParamsSample: v1PlayerTransferOutPs{"user_id", 123.45, "trace_id 4-36"},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type v1PlayerTransferInPs struct {
	UserID string
	// Language string
	Amount  float64
	TraceId string
}

type v1PlayerTransferInRet struct {
	AfterBalance float64
}

type TransferOrder struct {
	ID         primitive.ObjectID `bson:"_id"`
	AppID      string
	TraceId    string
	CreateTime time.Time
	Pid        int64
	UserID     string
	Amount     float64
	Action     string // in out allout
	Error      string
	Completed  bool
}

func v1_player_transfer_in(app *operator.MemApp, ps v1PlayerTransferInPs, ret *v1PlayerTransferInRet) (err error) {
	_, err = operator.AppMgr.EnsureUserExists(app, ps.UserID)
	if err != nil {
		return err
	}

	err = transfer(app, ps.TraceId, ps.UserID, ps.Amount, ret, "in")
	return
}

type v1PlayerTransferOutPs struct {
	UserID  string
	Amount  float64
	TraceId string
}

func v1_player_transfer_out(app *operator.MemApp, ps v1PlayerTransferOutPs, ret *v1PlayerTransferInRet) (err error) {
	err = transfer(app, ps.TraceId, ps.UserID, ps.Amount, ret, "out")
	return
}

func transfer(app *operator.MemApp, traceId string, userId string, amount float64, ret *v1PlayerTransferInRet, action string /*in out*/) (err error) {
	if app.WalletMode != comm.WalletModeTransfer {
		err = errors.New("Not transfer Mode")
		return
	}

	if amount <= 0 || gamedata.GetTransferLimitMap().GetLimit(app.AppID) < amount {
		err = errors.New("Incorrect transaction amount")
		return
	}

	memplr, err := operator.AppMgr.GetPlr2(app, userId)
	if err != nil {
		return
	}

	// lo.Must0(app.Address == "")

	now := time.Now()

	order := &TransferOrder{
		ID:         primitive.NewObjectIDFromTimestamp(now),
		TraceId:    traceId,
		CreateTime: now,
		UserID:     memplr.Uid,
		AppID:      memplr.AppID,
		Amount:     amount,
		Pid:        memplr.Pid,
		Action:     action,
	}
	coll := db.Collection2("game", "orders")
	_, err = coll.InsertOne(context.TODO(), order)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			err = errors.New("Order already exists")
		}
		return
	}

	var (
		balance int64
	)

	var logtype int
	if action == "in" {
		logtype = 1
		balance, err = memplr.TransferIn(ut.Money2Gold(amount))
	} else {
		logtype = 2
		balance, err = memplr.TransferOut(ut.Money2Gold(amount))
	}
	ret.AfterBalance = ut.Gold2Money(balance)

	update := db.D(
		"Completed", true,
	)

	if err != nil {
		db.AppendE(&update, "Error", err.Error())
	}
	coll.UpdateByID(context.TODO(), order.ID, db.D("$set", update))

	if err == nil {
		// slotsmongo.AddBetLog(&slotsmongo.AddBetLogParams{
		// 	Pid: memplr.Pid,
		// 	Balance: balance,

		// })
		// slotsmongo.ReportsCollection("")

		slotsmongo.InsertBetLog(&slotsmongo.DocBetLog{
			Pid:            memplr.Pid,
			UserID:         memplr.Uid,
			AppID:          memplr.AppID,
			LogType:        logtype,
			TransferAmount: ut.Money2Gold(amount),
			Balance:        balance,
		})

		if logtype == 2 {
			alerterlib.OnTransferOut(&alerterlib.OnTransferOutPs{
				ID:      order.ID,
				Pid:     memplr.Pid,
				UserID:  memplr.Uid,
				AppID:   memplr.AppID,
				Amount:  amount,
				Balance: ret.AfterBalance,
			})
		}
	}

	operator.InsertModifyLog(&operator.DocModifyGoldLog{
		ID:      order.ID,
		Change:  lo.Ternary(action == "in", ut.Money2Gold(amount), -ut.Money2Gold(amount)),
		Balance: balance,
		Error:   ut.ErrString(err),
		Action:  action,
		ReqData: order,
	}, memplr)

	return
}
