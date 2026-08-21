package handlers

import (
	"context"
	"game/comm"
	"game/comm/db"
	"game/comm/define"
	"game/comm/mux"
	"game/service/gamecenter/internal/operator"

	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/transaction/queryOrder",
		Handler:      v1_transaction_query_order,
		Desc:         "查询订单",
		Kind:         "api/v1",
		ParamsSample: queryOrderPs{"operator_user_abcd"},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type queryOrderPs struct {
	TraceId string
}

type queryOrderRet struct {
	Order *TransferOrder
}

func v1_transaction_query_order(app *operator.MemApp, ps queryOrderPs, ret *queryOrderRet) (err error) {
	if app.WalletMode != comm.WalletModeTransfer {
		err = define.NewErrCode("Not transfer Mode", 1009)
		return
	}
	var order TransferOrder
	coll := db.Collection2("game", "orders")
	// filter := db.D(
	// 	"_id", ps.TraceId,
	// )

	filter := db.D("TraceId", ps.TraceId, "AppID", app.AppID)
	err = coll.FindOne(context.TODO(), filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = define.NewErrCode("Order does not exist", 1018)
		}
		return
	}

	ret.Order = &order
	return
}
