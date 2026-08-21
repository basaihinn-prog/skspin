package slotsmongo

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"serve/comm/db"
	"serve/comm/ut"

	"go.mongodb.org/mongo-driver/bson"
)

type PlayerDailyDoc struct {
	ID         string `json:"-" bson:"_id"`
	BetAmount  int64  `bson:"betamount"` // 支出
	WinAmount  int64  `bson:"winamount"` // 赚取
	SpinCount  int    `bson:"spincount"` // 转动次数
	RoundCount int    `bson:"roundcount"`

	Date  string `bson:"date"`
	AppID string `bson:"appid"`
	Pid   int64  `bson:"pid"`
}

func addPayerReport(dateStr, appId string, pid int64, incDoc bson.D) {
	coll := ReportsCollection("PlayerDailyReport")

	// dateStr := time.Now().Format("20060102")
	id := ut.JoinStr(':', dateStr, appId, strconv.Itoa(int(pid)))

	update := db.D(
		"$set", db.D(
			"date", dateStr,
			"appid", appId,
			"pid", pid,
		),
		"$inc", incDoc,
	)

	coll.UpdateByID(context.TODO(), id, update, db.UpsertOpt())
}

type retentionResult struct {
	Pids []int64 `bson:"pids"`
}

func addPayerDailyRetention(dateStr, game, appId string, pid int64) {
	collDaily := ReportsCollection("PlayerDailyRetention")

	// dateStr := time.Now().Format("20060102")
	id := ut.JoinStr(':', dateStr, appId, game)

	update := db.D(
		"$setOnInsert", db.D(
			"date", time.Now(),
			"appid", appId,
			"gameId", game,
		))

	result, err := collDaily.UpdateByID(context.TODO(), id, update, options.Update().SetUpsert(true))
	if err != nil {
		slog.Info("addPayerDailyRetention", "err", ut.ErrString(err))
	}
	update = db.D("$addToSet", db.D("pids", pid))
	result, err = collDaily.UpdateByID(context.TODO(), id, update)
	if err != nil {
		slog.Info("addPayerDailyRetention", "err", ut.ErrString(err))
	}
	if result.ModifiedCount > 0 {
		coll := ReportsCollection("PlayerRetentionReport")
		idGame := ut.JoinStr(':', dateStr, appId, game)
		incDoc := db.D(
			"RetentionPlayerCount", 1,
		)
		updateGame := db.D(
			"$set", db.D(
				"date", time.Now(),
				"appid", appId,
				"gameId", game,
			),
			"$inc", incDoc)
		_, err = coll.UpdateByID(context.TODO(), idGame, updateGame, options.Update().SetUpsert(true))
		if err != nil {
			slog.Info("updatePayerGameRetention", "err", ut.ErrString(err))
		}

		query := bson.M{}
		retentionResults := []*retentionResult{}
		query["_id"] = primitive.Regex{
			Pattern: "^" + ut.JoinStr(':', dateStr, appId),
		}
		cur, err := collDaily.Find(context.TODO(), query)
		if err != nil {
			slog.Info("addPayerDailyRetention", "err", ut.ErrString(err))
		}
		if err != nil {
			return
		}

		defer cur.Close(context.TODO())

		if err = cur.All(context.TODO(), &retentionResults); err != nil {
			return
		}
		pidList := []int64{}
		for _, p := range retentionResults {
			pidList = append(pidList, p.Pids...)
		}
		count := ut.Int64InArrCount(pidList, pid)

		if count == 1 {
			idApp := ut.JoinStr(':', dateStr, appId)
			updateApp := db.D(
				"$set", db.D(
					"date", time.Now(),
					"appid", appId,
				),
				"$inc", incDoc)
			_, err = coll.UpdateByID(context.TODO(), idApp, updateApp, options.Update().SetUpsert(true))
			if err != nil {
				slog.Info("updatePayerRetention", "err", ut.ErrString(err))
			}
		}
	}
}
