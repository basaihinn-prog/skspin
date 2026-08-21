package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"serve/comm/db"
	"serve/comm/plats/pg"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/samber/lo"
)

func verifyOperatorPlayerSession() {

}

func main() {
	mongoaddr := "mongodb://myUser:123456@192.168.1.2:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_48")

	for i := 0; i < 1; i++ {
		num := i
		go normal(num)
		// go game(num)
		time.Sleep(time.Second)
	}
	select {}
}

func normal(i int) {
	uid := fmt.Sprintf("ccgg48_%d", i)

	pinfo := lo.Must(pg.GetPlayerToken(uid, "48"))

	coll := db.Collection(fmt.Sprintf("pgSpinData_%d", i))

	ps := url.Values{}
	ps.Add("cs", "0.01")
	ps.Add("ml", "1")
	ps.Add("pf", "2")
	ps.Add("id", pinfo.Sid)
	ps.Add("wk", "0_C")
	ps.Add("btt", "1")
	ps.Add("atk", pinfo.Tk)
	ps.Add("fb", "false")
	var bsonDocs []interface{}
	for ; ; time.Sleep(time.Second) {
		var spinret struct {
			Si json.RawMessage
			// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.14.0/bson#NewDecoder
		}

		spinurl := lo.Must(url.JoinPath(pinfo.Geu, "v2/Spin"))

		if err := pg.InvokePGService(spinurl, ps, &spinret); err != nil {
			slog.Error("spin error", "error", err)
			continue
		}
		// os.Stdout.Write(spinret)

		bsondoc := lo.Must(db.Json2Bson(spinret.Si))
		sid := bsondoc.Lookup("sid").StringValue()
		if sid == bsondoc.Lookup("psid").StringValue() && len(bsonDocs) > 0 {
			firstDoc := bsonDocs[0].(bson.Raw)
			// 判断是否是从初始盘开始拉去的
			if firstDoc.Lookup("psid").StringValue() == firstDoc.Lookup("sid").StringValue() {
				lo.Must(coll.InsertMany(context.TODO(), bsonDocs))
			}
			bsonDocs = []interface{}{}
		}
		bsonDocs = append(bsonDocs, bsondoc)
		//lo.Must(coll.InsertOne(context.TODO(), bsondoc))
		ps.Set("id", sid)
	}
}
