package slotsmongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/ut"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BDItem = define.BDItem
type BHItem = define.BHItem

func InsertBetHistory(pid, serverid int64, bdRecords []map[string]any, betIds []primitive.ObjectID, currencyKey string) primitive.ObjectID {
	coll := db.Collection("BetHistory")

	bdItems := make([]*BDItem, 0, len(bdRecords))
	for i := range bdRecords {
		bdItems = append(bdItems, &BDItem{
			Tid:  betIds[i].Hex(),
			Tba:  ut.GetFloat(bdRecords[i]["tb"]),
			Twla: ut.GetFloat(bdRecords[i]["np"]),
			Bl:   ut.GetFloat(bdRecords[i]["bl"]),
			Bt:   time.Now().UnixMilli()*100 + int64(i),
			// Bt:   bdRecords[i]["bt"].(int64),
			Gd: lo.Must(ut.GetJsonRaw(bdRecords[i])),
		})
	}

	lastPan := bdRecords[len(bdRecords)-1]
	gtwla := ut.GetFloat(lastPan["aw"]) - ut.GetFloat(bdRecords[0]["tb"])
	bhitem := &BHItem{
		ID:  primitive.NewObjectID(),
		Pid: pid,
		Tid: lastPan["psid"].(string),
		Gid: int(serverid),
		// CC:    "THB",
		CC:    currencyKey,
		Gtba:  ut.GetFloat(lastPan["tbb"]),
		Gtwla: gtwla,
		Bt:    time.Now().UnixMilli(),
		Ge:    lastPan["ge"],
		Bd:    bdItems,
		Mgcc:  0,
		Fscc:  0,
	}

	coll.InsertOne(context.TODO(), bhitem)
	return bhitem.ID
}

func InsertBetHistoryEvery(objectId primitive.ObjectID, pid, serverid int64, num int, bdRecord map[string]any, betLogId string, currencyKey string, isEnd bool) primitive.ObjectID {
	coll := db.Collection("BetHistory")
	bdItem := &BDItem{
		Tid:  betLogId,
		Tba:  ut.GetFloat(bdRecord["tb"]),
		Twla: ut.GetFloat(bdRecord["np"]),
		Bl:   ut.GetFloat(bdRecord["bl"]),
		Bt:   time.Now().UnixMilli(),
		// Bt:   bdRecords[i]["bt"].(int64),
		Gd: lo.Must(ut.GetJsonRaw(bdRecord)),
	}
	gtwla := ut.GetFloat(bdRecord["aw"]) - ut.GetFloat(bdRecord["tbb"])
	body := PsidMapBody{}
	if isEnd {
		bhitem := &BHItem{
			ID:  objectId,
			Pid: pid,
			Tid: betLogId,
			Gid: int(serverid),
			// CC:    "THB",
			CC:    currencyKey,
			Gtba:  ut.GetFloat(bdRecord["tbb"]),
			Gtwla: gtwla,
			Bt:    time.Now().UnixMilli(),
			Ge:    bdRecord["ge"],
			Bd:    []*BDItem{bdItem},
			Mgcc:  0,
			Fscc:  0,
		}
		coll.InsertOne(context.TODO(), bhitem)
		objectId = bhitem.ID
		body.Sid, body.Psid = betLogId, betLogId
	} else {
		var res Result
		update := bson.M{
			"gtba":  ut.GetFloat(bdRecord["tbb"]),
			"gtwla": gtwla,
			"mgcc":  0,
		}
		coll.FindOneAndUpdate(context.TODO(), bson.M{"_id": objectId}, bson.M{"$set": update, "$push": bson.M{"bd": bdItem},
			"$addToSet": bson.M{"ge": bson.M{"$each": bdRecord["ge"]}}}).Decode(&res)
		body.Sid = betLogId
		body.Psid = res.Tid
	}
	coll2 := db.Collection("psidMap")
	coll2.InsertOne(context.TODO(), body)

	return objectId
}

type PsidMapBody struct {
	Sid  string `json:"sid" bson:"sid"`
	Psid string `json:"psid" bson:"psid"`
}

type Result struct {
	Tid string `json:"tid" bson:"tid"`
}
