package slotsmongo

import (
	"context"
	"fmt"
	"game/comm/db"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetTTLLogColl(name string, ttldays int) *mongo.Collection {
	datestr := time.Now().Format("20060102")
	collname := fmt.Sprintf("%s:%d:%s", name, ttldays, datestr)
	coll := db.Collection2("ttllog", collname)

	return coll
}

func DelExpiredTTLLogColls() {
	d := db.Client().Database("ttllog")
	names, _ := d.ListCollectionNames(context.TODO(), db.D())

	now := time.Now()

	for _, name := range names {
		// fmt.Println(name)

		arr := strings.Split(name, ":")
		if len(arr) != 3 {
			continue
		}

		date := arr[2]

		ttldays, err := strconv.Atoi(arr[1])
		if err != nil {
			slog.Error("DelExpiredTTLLogColls", "err", err.Error(), "arr", arr)
			continue
		}

		ctime, err := time.ParseInLocation("20060102", date, time.Local)
		if err != nil {
			slog.Error("DelExpiredTTLLogColls", "err", err.Error(), "date", date)
			continue
		}

		if time.Duration(ttldays)*time.Hour*24 < now.Sub(ctime) {
			d.Collection(name).Drop(context.TODO())
		}
	}
}
