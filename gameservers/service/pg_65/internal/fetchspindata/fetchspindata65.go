package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"serve/comm/db"
	"serve/comm/plats/pg"

	"github.com/samber/lo"
)

func verifyOperatorPlayerSession() {

}

func main() {
	//mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	//db.DialToMongo(mongoaddr, "pg_65")

	for i := 0; i < 10; i++ {
		num := i
		normal(num)
		time.Sleep(time.Second)
	}
	select {}
}

func normal(i int) {
	uid := fmt.Sprintf("test65_%d", i)

	pinfo := lo.Must(pg.GetPlayerToken(uid, "65"))

	coll := db.Collection(fmt.Sprintf("pgSpinData_%d", i))

	ps := url.Values{}
	ps.Add("cs", "0.05")
	ps.Add("ml", "1")
	ps.Add("sn", "1")
	ps.Add("pf", "3")
	ps.Add("id", pinfo.Sid)
	ps.Add("wk", "0_C")
	ps.Add("btt", "1")
	ps.Add("atk", pinfo.Tk)
	var bsonDocs []interface{}
	for ; ; time.Sleep(time.Second) {
		var spinret struct {
			Si json.RawMessage
			// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.14.0/bson#NewDecoder
		}

		spinurl := lo.Must(url.JoinPath(pinfo.Geu, "v2/Spin"))

		if err := pg.InvokePGService(spinurl, ps, &spinret); err != nil {
			slog.Error("spin game error", "error", err, "sid",
				pinfo.Sid, "token", pinfo.Tk, "uid", uid)
			time.Sleep(5 * time.Second)
			pinfo = lo.Must(pg.GetPlayerToken(uid, "65"))
			bsonDocs = []interface{}{}
			ps.Set("id", pinfo.Sid)
			ps.Set("atk", pinfo.Tk)
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

func game(i int) {
	uid := fmt.Sprintf("testuseryygame_65_%d", i)

	pinfo := lo.Must(pg.GetPlayerToken(uid, "65"))

	coll := db.Collection("pgSpinDataGame")
	ps := url.Values{}
	ps.Add("cs", "0.05")
	ps.Add("ml", "1")
	ps.Add("pf", "1")
	ps.Add("id", pinfo.Sid)
	ps.Add("wk", "0_C")
	ps.Add("btt", "1")
	ps.Add("atk", pinfo.Tk)
	ps.Add("fb", "2")
	var bsonDocs []interface{}
	for ; ; time.Sleep(time.Second) {
		var spinret struct {
			Si json.RawMessage
			// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.14.0/bson#NewDecoder
		}

		spinurl := lo.Must(url.JoinPath(pinfo.Geu, "v2/Spin"))

		if err := pg.InvokePGService(spinurl, ps, &spinret); err != nil {
			slog.Error("spin buy game error", "error", err, "sid",
				pinfo.Sid, "token", pinfo.Tk)
			return
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
