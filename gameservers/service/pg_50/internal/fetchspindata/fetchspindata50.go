package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"
	"time"

	"serve/comm/db"
	"serve/service/pg_50/internal/dealData"

	"go.mongodb.org/mongo-driver/bson"

	"serve/comm/plats/pg"

	"github.com/samber/lo"
)

func verifyOperatorPlayerSession() {

}

func main() {
	mongoaddr := "mongodb://myUser:123456@192.168.1.2:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_50")

	for i := 0; i < 1; i++ {
		num := i
		go normal(num)
		time.Sleep(time.Second)
	}
	time.Sleep(time.Second * 20)
	num := int(0)
	for ; ; time.Sleep(time.Minute) {
		err := db.MiscGet("pullNum", &num)
		if err != nil {
			break
		}
		if num >= 350000 {
			fmt.Println("拉取数据结束")
			break
		}
	}
	wg := sync.WaitGroup{}
	collNames, _ := db.Client().Database("pg_50").ListCollectionNames(context.TODO(), bson.M{})
	for i := 0; i < len(collNames); i++ {
		if strings.HasPrefix(collNames[i], "pgSpinData_") {
			name := collNames[i]
			wg.Add(1)
			go func() {
				dealData.DealNormal(name)
				wg.Done()
			}()
		} else if strings.HasPrefix(collNames[i], "pgSpinDataGame_") {
			name := collNames[i]
			wg.Add(1)
			go func() {
				dealData.DealGame(name)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func normal(i int) {
	uid := fmt.Sprintf("ccgg50_%d", i)

	pinfo := lo.Must(pg.GetPlayerToken(uid, "50"))

	coll := db.Collection(fmt.Sprintf("pgSpinData_%d", i))

	ps := url.Values{}
	ps.Add("id", pinfo.Sid)
	ps.Add("cs", "0.01")
	ps.Add("ml", "1")
	ps.Add("wk", "0_C")
	ps.Add("btt", "1")
	ps.Add("atk", pinfo.Tk)
	ps.Add("pf", "2")
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
				db.MiscInc("pullNum", 1)
			}
			bsonDocs = []interface{}{}
		}
		bsonDocs = append(bsonDocs, bsondoc)
		//lo.Must(coll.InsertOne(context.TODO(), bsondoc))
		ps.Set("id", sid)
	}
}
