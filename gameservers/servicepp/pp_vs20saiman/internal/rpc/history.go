package rpc

import (
	"context"
	"encoding/json"
	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/servicepp/ppcomm"
	"time"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	ppcomm.RegRpc("/gs2c/api/history/v2/play-session/last-items", history_last_items)
	ppcomm.RegRpc("/gs2c/api/history/v3/action/children", history_children)
	ppcomm.RegRpc("/gs2c/api/history/v2/play-session/by-round", history_by_round)
}

// https://5g6kpi7kjf.uapuqhki.net/gs2c/api/history/v2/play-session/last-items?token=5993dd000ece684f7bee1a21706f74744e1487f7ac14d2fe21b533dbb625655d&symbol=vs20olympx

func history_last_items(msg *nats.Msg) (ret []byte, err error) {
	ps := ppcomm.ParseVariables(string(msg.Data))
	pid, err := jwtutil.ParseToken(ps.Str("token"))
	if err != nil {
		return
	}
	coll := db.Collection("BetHistory")
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
	zeroTime := startOfDay.Unix() * 1000
	match := bson.M{
		"datetime": bson.M{"$lte": zeroTime},
		"pid":      pid,
	}
	findOptions := options.Find()
	findOptions.SetLimit(100)
	findOptions.SetSort(bson.M{"_id": -1})
	var docs []*ppcomm.RspHistoryData
	cur, err := coll.Find(context.TODO(), match, findOptions)
	if err != nil {
		return
	}
	defer db.CloseCursor(cur)
	err = cur.All(context.TODO(), &docs)
	if err != nil {
		return
	}

	ret, _ = json.Marshal(docs)

	//ret = []byte(`[{"roundId":39295534355091,"dateTime":1726710437000,"bet":"10.00","win":"0.00","balance":"1,000,156.60","roundDetails":null,"currency":"THB","currencySymbol":"฿","hash":"540247a327ba6ea7974c885e88bd7d43"}]`)
	return
}

// https://5g6kpi7kjf.uapuqhki.net/gs2c/api/history/v3/action/children?id=39295534355091&token=e9f224bec8a561da3ec0f1e0444c4e8d3c0e5b3ed2c69d4465a46c1e68a0d9b6&symbol=vs20olympx

func history_children(msg *nats.Msg) (ret []byte, err error) {
	ps := ppcomm.ParseVariables(string(msg.Data))
	//pid, err := jwtutil.ParseToken(ps.Str("token"))
	// if err != nil {
	// 	return
	// }

	rid := ps.Str("id")
	coll := db.Collection("BetHistory")
	var doc struct {
		RoundDetails []*ppcomm.HistoryDetailsData `bson:"rounddetails"`
	}
	//objId, err := primitive.ObjectIDFromHex(rid)
	//if err != nil {
	//	return
	//}
	err = coll.FindOne(context.TODO(), bson.M{"betId": rid}, options.FindOne().SetProjection(bson.M{"rounddetails": 1})).Decode(&doc)
	if err != nil {
		return
	}

	response := map[string]any{
		"data":        doc.RoundDetails,
		"description": "OK",
		"error":       0,
	}
	ret, _ = json.Marshal(response)

	//fmt.Println(pid)
	//
	//ret = []byte(`{"error":0,"description":"OK","data":[{"roundId":39295534355091,"request":{"symbol":"vs20olympx","c":"0.5","repeat":"0","action":"doSpin","index":"2","bl":"0","counter":"3","l":"20"},"response":{"accm":"cp","tw":"0.00","c":"0.50","acci":"0","sver":"5","index":"2","balance_cash":"1,000,156.60","bl":"0","stime":"1726710437112","counter":"4","ntp":"-10.00","rid":"39295534355091","l":"20","reel_set":"1","sa":"10,9,1,8,8,4","sb":"6,10,11,11,9,5","balance_bonus":"0.00","na":"s","s":"5,9,5,7,9,3,5,7,11,10,4,3,8,11,5,3,5,9,10,8,8,11,5,9,11,8,8,3,7,5","balance":"1,000,156.60","sh":"5","w":"0.00","accv":"0"},"currency":"THB","currencySymbol":"฿","configHash":"54a2aea7ec82be563930a69a8af45fed"}]}`)
	return
}

func history_by_round(msg *nats.Msg) (ret []byte, err error) {
	// id := string(msg.Data)

	ps := ppcomm.ParseVariables(string(msg.Data))
	id := ps.Str("id")
	coll := db.Collection("BetHistory")
	var doc *ppcomm.RspHistoryData
	//if false {
	//	objId, err := primitive.ObjectIDFromHex(id)
	//	if err != nil {
	//		return
	//	}
	//
	//	err = coll.FindOne(context.TODO(), db.ID(objId), options.FindOne().SetProjection(bson.M{"rounddetails": 0})).Decode(&doc)
	//	if err != nil {
	//		return
	//	}
	//}
	err = coll.FindOne(context.TODO(), db.D("betId", id), options.FindOne().SetProjection(bson.M{"rounddetails": 0})).Decode(&doc)
	if err != nil {
		return
	}
	ret, _ = json.Marshal(doc)
	//ret = []byte(`{"roundId":39297197956091,"dateTime":1726715459000,"bet":"0.00","win":"0.00","balance":"1,000,160.65","roundDetails":"Free spin","currency":"THB","currencySymbol":"฿","hash":"1bb82d5647704d8ff5ab181659d64364"}`)
	return
}
