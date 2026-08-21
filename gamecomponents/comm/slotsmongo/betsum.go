package slotsmongo

import (
	"context"
	"game/comm/db"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AddBetDoc struct {
	// ID  primitive.ObjectID `bson:"_id"`
	Pid int64
	// Seq  int
	Bets []int64
	// Sum  int64
	// Time time.Time
}

func AddBet(uid int64, bet int64) (err error) {
	coll := SlotsCollection("BetSum")

	_, err = coll.UpdateByID(context.TODO(), uid, db.D(
		"$push", db.D(
			"bets", db.D(
				"$each", []int64{bet},
				"$slice", -200,
			),
		),
	), options.Update().SetUpsert(true))

	return
}

func AvgBet(uid int64) (avg int64, err error) {
	coll := SlotsCollection("BetSum")

	var doc AddBetDoc
	err = coll.FindOne(context.TODO(), db.ID(uid)).Decode(&doc)
	if err == mongo.ErrNoDocuments || len(doc.Bets) == 0 {
		err = nil
		return
	}

	sum := lo.Sum(doc.Bets)

	avg = sum / int64(len(doc.Bets))

	/*
		cursor, err := coll.Find(context.TODO(), db.D("pid", uid), options.Find().SetLimit(count))
		if err != nil {
			return
		}
		var docs []AddBetDoc
		err = cursor.All(context.TODO(), &docs)
		if err != nil {
			return
		}

		if len(docs) == 0 {
			avg = 0
			return
		}
	*/

	return
}
