package main

import (
	"context"
	"game/comm/db"
	"game/duck/lazy"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
)

func main() {
	lazy.Init("admin")
	ConnectMongoByConfig()
	deal()
	//go clearPlayerRetentionReport()
	lazy.SignalProc()
}

func deal() {
	c := cron.New()
	c.AddFunc("59 * * * *", func() {
		err := clearPlayerRetentionReport()
		if err != nil {
			slog.Error("update reports PlayerRetentionReport err:  ", err)
		}
	})
	c.Start()
}

func clearPlayerRetentionReport() error {
	var err error
	coll := db.Collection2("reports", "PlayerRetentionReport")
	result, err := coll.DeleteMany(context.TODO(), bson.D{})
	slog.Info("clearPlayerRetentionReport result count:  ", result.DeletedCount)
	coll = db.Collection2("reports", "PlayerDailyRetention")
	result, err = coll.DeleteMany(context.TODO(), bson.D{})
	slog.Info("clearPlayerDailyRetention result count:  ", result.DeletedCount)
	defer func() {
		if err != nil {
			slog.Error("clearPlayerRetentionReport err:  ", err)
		}
	}()
	return err
}
