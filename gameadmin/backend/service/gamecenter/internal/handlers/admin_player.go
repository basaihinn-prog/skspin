package handlers

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/mux"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/lang"
	"game/service/gamecenter/internal/operator"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	mux.RegRpc("/gamecenter/player/balance", "获取玩家余额 - 包含未结算 (过期函数, 不应继续使用)", "admin rpc", rpcPlayerBalance, nil)
	mux.RegRpc("/gamecenter/player/setBalance", "设置玩家余额", "admin rpc", rpcSetBalance, nil)
}

type rpcPlayerBalancePs struct {
	Pid int64
}

// 获取玩家余额 - 包含未结算 (过期函数, 不应继续使用)
func rpcPlayerBalance(ps rpcPlayerBalancePs, ret *v1PlayerBalanceRet) (err error) {
	memplr, err := operator.AppMgr.GetPlr(ps.Pid)
	if err != nil {
		return
	}

	err = getBalanceByMemplr(memplr, ret)
	return
}

type setBalanceReq struct {
	Pid           int64   // 玩家ID
	Balance       float64 // 金币变化：>0 加钱， <0 扣钱
	BeforeBalance float64

	Proxy_pid      int64
	Proxy_username string
	Proxy_isadmin  bool
	Proxy_language string
}

type setBalanceRet struct {
	BeforeBalance float64
}

func rpcSetBalance(ps setBalanceReq, ret *setBalanceRet) (err error) {
	if !ps.Proxy_isadmin {
		err = lang.Error(ps.Proxy_language, "权限不足")
		return
	}

	// before := time.Now()
	if ps.Balance < 0 {
		err = fmt.Errorf("error balance value %f", ps.Balance)
		return
	}
	am := operator.AppMgr

	plr, err := am.GetPlr(ps.Pid)
	if err != nil {
		return
	}
	app := am.GetApp(plr.AppID)
	if app == nil {
		err = errors.New("app not found.")
		return
	}
	if app.WalletMode != comm.WalletModeTransfer {
		err = errors.New("Not transfer Mode")
		return
	}

	balance, err := plr.Balance()
	if err != nil {
		return
	}
	if ps.BeforeBalance != 0 && 1.0 < math.Abs(ut.Gold2Money(balance)-ps.BeforeBalance) {
		err = errors.New("The balance has changed")
		return
	}

	newgold := ut.Money2Gold(ps.Balance)
	oldgold, err := plr.SetGold(newgold)
	ret.BeforeBalance = ut.Gold2Money(oldgold)

	// logid := primitive.NewObjectID()
	// operator.InsertModifyLog(logid, plr, newgold-oldgold, "后台设置余额", before, newgold, err)

	id := primitive.NewObjectID()

	operator.InsertModifyLog(&operator.DocModifyGoldLog{
		ID:      id,
		Change:  newgold - oldgold,
		Balance: newgold,
		Error:   ut.ErrString(err),
		Action:  "adminSetBalance",
		ReqData: ps,
	}, plr)

	slotsmongo.InsertBetLog(&slotsmongo.DocBetLog{
		ID:     id,
		Pid:    plr.Pid,
		UserID: plr.Uid,
		AppID:  plr.AppID,
		// Completed: true,
		LogType:        3,
		TransferAmount: newgold - oldgold,
		Balance:        newgold,
	})

	coll := db.Collection2("GameAdmin", "AdminOperator")
	operator := &comm.Operator{}
	err = coll.FindOne(context.TODO(), bson.M{"AppID": plr.AppID}).Decode(&operator)
	if err != nil {
		return err
	}

	// 后台玩家信息里面使用
	coll = db.Collection2("GameAdmin", "SlotsPoolHistory")
	coll.InsertOne(context.TODO(), bson.M{
		"_id":          id,
		"OpName":       ps.Proxy_username,
		"OpPid":        ps.Proxy_pid,
		"AnimUserPid":  ps.Pid,
		"AnimUserName": plr.Uid,
		"Type":         3,
		"OldGold":      oldgold,
		"NewGold":      newgold,
		"Change":       newgold - oldgold,
		"time":         time.Now().Format("2006-01-02 15:04:05"),
		"Currency":     operator.CurrencyKey,
		"AppID":        operator.AppID,
	})

	return
}
