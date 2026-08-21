package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"serve/comm/plats/pg"

	"go.mongodb.org/mongo-driver/bson"

	"serve/comm/db"

	"github.com/samber/lo"
)

func verifyOperatorPlayerSession() {

}

func main() {
	mongoaddr := "mongodb://myUser:123456@192.168.1.2:27017/?authSource=admin"

	db.DialToMongo(mongoaddr, "pg_1451122")

	for i := 0; i < 1; i++ {
		num := i
		go normal(num)
		//go game(num)
		time.Sleep(time.Second)
	}
	select {}
}

func normal(i int) {
	uid := fmt.Sprintf("1451122tdx%d", i)

	pinfo := lo.Must(pg.GetPlayerToken(uid, "1451122"))

	coll := db.Collection(fmt.Sprintf("pgSpinData_%d", i))

	ps := url.Values{}
	ps.Add("cs", "0.04")
	ps.Add("ml", "1")
	ps.Add("pf", "4")
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
			fmt.Println("开始重试。。。")
			pinfo = lo.Must(pg.GetPlayerToken(uid, "1451122"))
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
	uid := fmt.Sprintf("1451122yygame%d", i)

	pinfo := lo.Must(pg.GetPlayerToken(uid, "1451122"))

	coll := db.Collection(fmt.Sprintf("pgSpinDataGame_%d", i))
	ps := url.Values{}
	ps.Add("cs", "0.1")
	ps.Add("ml", "1")
	ps.Add("pf", "4")
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
				pinfo.Sid, "token", pinfo.Tk, "uid", uid)
			time.Sleep(5 * time.Second)
			fmt.Println("开始重试。。。")
			pinfo = lo.Must(pg.GetPlayerToken(uid, "1451122"))
			bsonDocs = []interface{}{}
			ps.Set("id", pinfo.Sid)
			ps.Set("atk", pinfo.Tk)
			continue
		}

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
