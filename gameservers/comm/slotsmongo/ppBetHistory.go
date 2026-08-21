package slotsmongo

import (
	"context"
	"fmt"
	"time"

	//"time"

	"go.mongodb.org/mongo-driver/bson"

	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/ut"

	//"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PPBDItem = define.PPBDItem
type PPBHItem = define.PPBHItem

func PPInsertBetHistoryEvery(objectId primitive.ObjectID, pid, serverid string, bdRecord map[string]any, betLogId string, cr, sr string, currencyKey string, isEnd bool) primitive.ObjectID {
	coll := db.Collection("BetHistory")
	pPBDItem := &PPBDItem{
		CR: cr,
		SR: sr,
	}
	if isEnd {
		pPBHItem := &PPBHItem{
			Id:         objectId,
			Tid:        betLogId,
			CC:         currencyKey,
			AgentCode:  "atou",
			UserCode:   pid,
			RoundID:    fmt.Sprintf("%d", time.Now().UnixMilli()),
			GameCode:   serverid,
			Bet:        ut.GetFloat(bdRecord["c"]) * ut.GetFloat(bdRecord["l"]),
			Data:       []*PPBDItem{pPBDItem},
			CreatedAt:  time.Now(),
			SharedLink: "",
			PlayedDate: time.Now().UnixMilli(),
		}
		coll.InsertOne(context.TODO(), pPBHItem)
		objectId = pPBHItem.Id
	} else {
		gtwla := ut.GetFloat(bdRecord["tw"])
		update := bson.M{
			"win": gtwla,
			"rtp": gtwla,
		}
		coll.UpdateOne(context.TODO(), bson.M{"_id": objectId}, bson.M{"$set": update, "$push": bson.M{"data": pPBDItem}})
	}

	return objectId
}
