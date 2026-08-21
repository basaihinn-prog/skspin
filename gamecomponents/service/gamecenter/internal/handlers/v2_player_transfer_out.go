package handlers

import (
	"context"
	"game/comm"
	"game/comm/db"
	"game/comm/define"
	"game/comm/mux"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/lazy"
	"game/service/alerter/alerterlib"
	"game/service/gamecenter/internal/operator"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v2/player/transferOutAll",
		Handler:      v2_player_transfer_out_all,
		Desc:         "玩家转出所有余额",
		Kind:         "api/v1",
		ParamsSample: v1PlayerTransferOutAllPs{"user_id", "trace_id 4-36"},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type v2PlayerTransferOutAllPs struct {
	UserID  string
	TraceId string
}

type v2PlayerTransferOutAllRet struct {
	Amount     float64
	RealAmount float64
}

func v2_player_transfer_out_all(app *operator.MemApp, ps v2PlayerTransferOutAllPs, ret *v2PlayerTransferOutAllRet) (err error) {
	if app.WalletMode != comm.WalletModeTransfer {
		err = define.NewErrCode("Not transfer Mode", 1009)
		return
	}
	memplr, err := operator.AppMgr.GetPlr2(app, ps.UserID)
	if err != nil {
		return
	}

	now := time.Now()
	sf := ut.NewSnowflake()
	id := strconv.Itoa(int(sf.NextID()))
	order := &TransferOrder{
		ID:         id,
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
			err = define.NewErrCode("Order already exists", 1024)
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
		if app.CurrencyCtrlStatus == 0 {
			item := lazy.GetCurrencyItemNew(app.CurrencyKey)
			ret.RealAmount = amount * item.Multi
		} else {
			ret.RealAmount = amount
		}
	}
	coll.UpdateByID(context.TODO(), order.ID, db.D("$set", update))

	if err == nil {
		betlogInfo := &slotsmongo.DocBetLog{
			ID:     order.ID,
			Pid:    memplr.Pid,
			UserID: memplr.Uid,
			AppID:  memplr.AppID,
			// Completed: true,
			LogType:        2,
			TransferAmount: ut.Money2Gold(amount),
			Balance:        0,
		}
		slotsmongo.InsertBetLog(betlogInfo)

		alerterlib.OnTransferOut(&alerterlib.OnTransferOutPs{
			ID:         order.ID,
			Pid:        memplr.Pid,
			UserID:     memplr.Uid,
			AppID:      memplr.AppID,
			Amount:     amount,
			Balance:    0,
			InsertTime: betlogInfo.InsertTime.Local().UTC().Format("2006-01-02 15:04:05"),
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
