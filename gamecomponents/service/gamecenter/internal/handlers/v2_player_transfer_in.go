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

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v2/player/transferIn",
		Handler:      v2_player_transfer_in,
		Desc:         "玩家充值",
		Kind:         "api/v2",
		ParamsSample: v2PlayerTransferInPs{"user_id", 123.45, "trace_id 4-36", 123450},
		Class:        "operator",
		GetArg0:      getArg0,
	})

	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v2/player/transferOut",
		Handler:      v2_player_transfer_out,
		Desc:         "玩家提现",
		Kind:         "api/v2",
		ParamsSample: v2PlayerTransferOutPs{"user_id", 123.45, "trace_id 4-36", 123450},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type v2PlayerTransferInPs struct {
	UserID string
	// Language string
	Amount     float64
	TraceId    string
	RealAmount float64
}

type v2PlayerTransferInRet struct {
	AfterBalance float64
	RealAmount   float64
}

func v2_player_transfer_in(app *operator.MemApp, ps v2PlayerTransferInPs, ret *v2PlayerTransferInRet) (err error) {
	_, err = operator.AppMgr.EnsureUserExists(app, ps.UserID)
	if err != nil {
		return err
	}

	err = transferV2(app, ps.TraceId, ps.UserID, ps.Amount, ps.RealAmount, ret, "in")
	return
}

type v2PlayerTransferOutPs struct {
	UserID     string
	Amount     float64
	TraceId    string
	RealAmount float64
}

func v2_player_transfer_out(app *operator.MemApp, ps v2PlayerTransferOutPs, ret *v2PlayerTransferInRet) (err error) {
	err = transferV2(app, ps.TraceId, ps.UserID, ps.Amount, ps.RealAmount, ret, "out")
	return
}

func transferV2(app *operator.MemApp, traceId string, userId string, amount, realAmount float64, ret *v2PlayerTransferInRet, action string /*in out*/) (err error) {
	if app.WalletMode != comm.WalletModeTransfer {
		err = define.NewErrCode("Not transfer Mode", 1009)
		return
	}

	if amount <= 0 || 1e8 < amount {
		err = define.NewErrCode("Incorrect transaction amount", 1016)
		return
	}
	changeAmount := int64(0)
	item := lazy.GetCurrencyItemNew(app.CurrencyKey)
	if app.CurrencyCtrlStatus == 0 {
		if amount*item.Multi != realAmount {
			err = define.NewErrCode("Incorrect transaction amount", 1016)
			return
		}
		changeAmount = ut.Money2Gold(amount)
	} else {
		if amount != realAmount {
			err = define.NewErrCode("Incorrect transaction amount", 1016)
			return
		}
		changeAmount = ut.Money2Gold(amount) / int64(item.Multi)
	}

	memplr, err := operator.AppMgr.GetPlr2(app, userId)
	if err != nil {
		return
	}

	// lo.Must0(app.Address == "")

	now := time.Now()
	sf := ut.NewSnowflake()
	id := strconv.Itoa(int(sf.NextID()))

	order := &TransferOrder{
		ID:         id,
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
			err = define.NewErrCode("Order already exists", 1017)
		}
		return
	}

	var (
		balance int64
	)

	var logtype int
	if action == "in" {
		logtype = 1
		balance, err = memplr.TransferIn(changeAmount)
	} else {
		logtype = 2
		balance, err = memplr.TransferOut(changeAmount)
	}
	ret.AfterBalance = ut.Gold2Money(balance)
	ret.RealAmount = realAmount

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
		betlogInfo := &slotsmongo.DocBetLog{
			ID:             order.ID,
			Pid:            memplr.Pid,
			UserID:         memplr.Uid,
			AppID:          memplr.AppID,
			LogType:        logtype,
			TransferAmount: changeAmount,
			Balance:        balance,
		}

		slotsmongo.InsertBetLog(betlogInfo)

		if logtype == 2 {
			alerterlib.OnTransferOut(&alerterlib.OnTransferOutPs{
				ID:         order.ID,
				Pid:        memplr.Pid,
				UserID:     memplr.Uid,
				AppID:      memplr.AppID,
				Amount:     ut.Gold2Money(changeAmount),
				Balance:    ret.AfterBalance,
				InsertTime: betlogInfo.InsertTime.Local().UTC().Format("2006-01-02 15:04:05"),
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
