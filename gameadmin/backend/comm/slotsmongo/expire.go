package slotsmongo

import (
	"context"
	"game/comm/db"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VisitDoc struct {
	ID        string `bson:"_id"`
	VisitTime time.Time
}

func IsFirstVisitToday1(key string) (todayFirst bool) {
	todayFirst, _ = IsFirstVisitToday(key)
	return todayFirst
}

// 检查是否是今天第一次访问, 就记录
func IsFirstVisitToday(key string) (todayFirst, historyFirst bool) {
	coll := ReportsCollection("visit")
	now := time.Now()
	var doc VisitDoc
	err := coll.FindOneAndUpdate(context.TODO(), db.ID(key), db.D(
		"$set", db.D("visittime", now),
	), options.FindOneAndUpdate().SetUpsert(true)).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		err = nil
		return true, true
	}
	if err != nil {
		return false, false
	}

	// doc.VisitTime.YearDay()

	same := isSameDate(now, doc.VisitTime.Local())
	// slog.Info("IsFirstVisitToday", "key", key, "LastVisitTime", doc.VisitTime, "err", err, "same", same)
	return !same, false
}

func isSameDate(t1, t2 time.Time) bool {
	lo.Must0(t1.Location() == t1.Location())
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
