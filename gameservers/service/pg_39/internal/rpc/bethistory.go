package rpc

import (
	"context"
	"encoding/json"
	"time"

	"serve/comm/db"
	"serve/comm/define"
	"serve/service/pg_39/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BHItem = define.BHItem
type BDItem = define.BDItem

func insertBetHistory(plr *models.Player, gd M, gdjson json.RawMessage) primitive.ObjectID {
	coll := db.Collection("BetHistory")
	// coll.
	bditem := &BDItem{
		Tid:  gd["sid"].(string),
		Tba:  gd["tb"].(float64),
		Twla: gd["np"].(float64),
		Bl:   gd["bl"].(float64),
		Bt:   time.Now().UnixMilli(),
		Gd:   gdjson,
	}

	bhitem := &BHItem{
		ID:  primitive.NewObjectID(),
		Pid: plr.PID,
		Tid: bditem.Tid,
		Gid: 39,
		// CC:    "PGC",
		CC:    plr.GetCurrencyOrTHB(),
		Gtba:  bditem.Tba,
		Gtwla: bditem.Twla,
		Bt:    time.Now().UnixMilli(),
		// Ge:    gd["ge"].([2]int),
		Ge:   gd["ge"],
		Bd:   []*BDItem{bditem},
		Mgcc: 0,
		Fscc: 0,
	}

	coll.InsertOne(context.TODO(), bhitem)
	return bhitem.ID
}
