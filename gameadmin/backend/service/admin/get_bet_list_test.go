package main

import (
	"context"
	"encoding/json"
	"fmt"
	"game/comm/db"
	"game/comm/slotsmongo"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sort"
	"testing"
)

func TestGetBetLog(t *testing.T) {
	db.DialToMongo("mongodb://127.0.0.1:27017/", "reports")

	coll := db.Collection("BetLog")
	var docs []*slotsmongo.DocBetLog
	cur, _ := coll.Find(context.TODO(), db.D("GameID", "Hilo"), options.Find().SetProjection(db.D("SpinDetailsJson", 0)))
	cur.All(context.TODO(), &docs)

	type R struct {
		Pid     int64
		UserID  string
		Bet     int64
		Win     int64
		WinLose int64
	}

	m := map[string]*R{}
	for _, doc := range docs {
		r := m[doc.UserID]
		if r == nil {
			r = &R{
				Pid:    doc.Pid,
				UserID: doc.UserID,
			}
			m[doc.UserID] = r
		}
		r.Bet += doc.Bet
		r.Win += doc.Win
		r.WinLose += doc.WinLose
	}

	vs := lo.Values(m)
	sort.Slice(vs, func(i, j int) bool {
		return vs[i].UserID < vs[j].UserID
	})
	jsonstr, _ := json.MarshalIndent(vs, "", "  ")
	os.WriteFile("HiloBetLog-summary.json", jsonstr, 0644)

	// jsonstr, _ := json.MarshalIndent(docs, "", "  ")
	// os.WriteFile("HiloBetLog.json", jsonstr, 0644)
	// fmt.Println(docs)
}

func TestRomoveBetLog(t *testing.T) {
	db.DialToMongo("mongodb://127.0.0.1:27017/", "reports")

	BetDailyReport := db.Collection("BetDailyReport")
	filter := db.D(
		"date", db.D("$lt", "20231122"),
		"buybetamount", db.D("$gt", 0),
	)
	cur, err := BetDailyReport.Find(context.TODO(), filter)
	assert.Nil(t, err)

	var docs []slotsmongo.BetDailyDoc
	cur.All(context.TODO(), &docs)
	// ut.WriteJsonFile("BetDailyReport-lt20231122.json", docs)

	BetSummaryReport := db.Collection("BetSummaryReport")
	for _, doc := range docs {
		id := doc.Game + ":" + doc.AppID

		update := db.D(
			"$inc", db.D(
				"betamount", -doc.BetAmount,
				"winamount", -doc.WinAmount,
				"buybetamount", -doc.BuyBetAmount,
				"buywinamount", -doc.BuyWinAmount,
			))
		BetSummaryReport.UpdateByID(context.TODO(), id, update)
	}

	BetDailyReport.DeleteMany(context.TODO(), filter)

}

func TestGetDayReport(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "reports")
	fmt.Println(HourReports())
}
