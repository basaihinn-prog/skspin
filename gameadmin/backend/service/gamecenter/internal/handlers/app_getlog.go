package handlers

/*
import (
	"context"
	"game/comm/db"
	"game/comm/mux"
	"game/duck/logger"
	"game/duck/mongodb"
	"game/pb/_gen/pb/gamepb"
	"game/service/gamecenter/internal/operator"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:    "/api/GetLog",
		Handler: getlog,
		Desc:    "获取下注日志",
		Kind:    "api",
		ParamsSample: &GetLogReq{
			From: &mongodb.TimeStamp{},
		},
		Class:   "operator",
		GetArg0: getArg0,
	})
}

type GetLogReq struct {
	From *mongodb.TimeStamp `protobuf:"bytes,1,opt,name=From,proto3" `
}

func getlog(app *operator.MemApp, req GetLogReq, ret *gamepb.DocBetLogList) (err error) {
	logger.Info("GetLog", app.AppID, req.From.AsTime())
	list, err := GetLogByCursor(req.From.AsTime(), bson.M{"AppID": app.AppID})
	if err != nil {
		return err
	}

	ret.List = list
	return
}

func GetLogByCursor(cursor time.Time, query bson.M) (ret []*gamepb.DocBetLog, err error) {
	maxCount := int64(5000)

	op := options.Find()
	op.SetLimit(maxCount)
	// op.SetSort(bson.M{"InsertTime": 1}) // 升序排列

	if query == nil {
		query = bson.M{}
	}

	// id := primitive.NewObjectIDFromTimestamp(cursor.Add(1))
	// var id  primitive.ObjectID

	// query["_id"] = bson.M{"$gt": id}

	query["InsertTime"] = bson.M{
		"$gt": cursor,
	}

	coll := db.Collection2("reports", "BetLog")
	mongoCursor, err := coll.Find(context.TODO(), query, op)
	if err != nil {
		return nil, err
	}
	defer mongoCursor.Close(context.TODO())

	err = mongoCursor.All(context.TODO(), &ret)

	return
}
*/
