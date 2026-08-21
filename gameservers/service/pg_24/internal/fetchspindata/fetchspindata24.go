package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"

	"serve/comm/db"
	"serve/comm/plats/pg"
	"serve/service/pg_24/internal/dealData"
	"serve/service/pg_24/internal/gendata"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/samber/lo"
)

func verifyOperatorPlayerSession() {

}

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc: %v MiB, TotalAlloc: %v MiB, Sys: %v MiB, NumGC: %v, Mallocs: %v, Frees: %v, StackInuse: %v \n",
		m.Alloc/1024/1024, m.TotalAlloc/1024/1024, m.Sys/1024/1024, m.NumGC, m.Mallocs, m.Frees, m.StackInuse)
}

func main() {
	mongoaddr := "mongodb://myUser:123456@192.168.1.2:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_24")

	//rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1; i++ {
		num := i
		go normal(num)
		// go game(num)
		//go gamex2(num)
		//go gamex3(num)

		// time.Sleep(time.Second * time.Duration(rand.IntN(5)+2))
		time.Sleep(time.Second)

	}
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
	collNames, _ := db.Client().Database("pg_24").ListCollectionNames(context.TODO(), bson.M{})
	for i := 0; i < len(collNames); i++ {
		if strings.HasPrefix(collNames[i], "pgSpinData") {
			name := collNames[i]
			wg.Add(1)
			go func() {
				dealData.DealNormal(name)
				wg.Done()
			}()
		} else if strings.HasPrefix(collNames[i], "pgSpinDataGamex2") {
			name := collNames[i]
			wg.Add(1)
			go func() {
				dealData.DealGamex2(name, gendata.GameTypeGamex2)
				wg.Done()
			}()
		} else if strings.HasPrefix(collNames[i], "pgSpinDataGamex3") {
			name := collNames[i]
			wg.Add(1)
			go func() {
				dealData.DealGamex3(name, gendata.GameTypeGamex3)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func normal(i int) {
	uid := fmt.Sprintf("24tdx%d", i)

	pinfo := lo.Must(pg.GetPlayerToken(uid, "24"))

	coll := db.Collection(fmt.Sprintf("pgSpinData_%d", i))

	ps := url.Values{}
	ps.Add("cs", "0.5")
	ps.Add("ml", "1")
	ps.Add("pf", "4")
	ps.Add("id", pinfo.Sid)
	ps.Add("wk", "0_C")
	ps.Add("btt", "1")
	ps.Add("bn", "1")
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
			// fmt.Println("开始重试。。。")
			pinfo = lo.Must(pg.GetPlayerToken(uid, "24"))
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

				printMemStats()

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

func gamex2(i int) {
	uid := fmt.Sprintf("24tdxx2game%d", i)

	pinfo := lo.Must(pg.GetPlayerToken(uid, "24"))

	coll := db.Collection(fmt.Sprintf("pgSpinDataGamex%d_%d", 2, i))
	ps := url.Values{}
	ps.Add("cs", "1")
	ps.Add("ml", "1")
	ps.Add("pf", "4")
	ps.Add("id", pinfo.Sid)
	ps.Add("wk", "0_C")
	ps.Add("btt", "1")
	ps.Add("atk", pinfo.Tk)

	// ps.Add("fb", "2")
	ps.Add("bn", "2")

	var bsonDocs []interface{}
	for ; ; time.Sleep(time.Second * 2) {
		var spinret struct {
			Si json.RawMessage
			// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.14.0/bson#NewDecoder
		}

		spinurl := lo.Must(url.JoinPath(pinfo.Geu, "v2/Spin"))

		if err := pg.InvokePGService(spinurl, ps, &spinret); err != nil {
			slog.Error("spin buy game error", "error", err, "sid",
				pinfo.Sid, "token", pinfo.Tk, "uid", uid)
			time.Sleep(5 * time.Second)
			// fmt.Println("gamex2, 开始重试。。。")
			pinfo = lo.Must(pg.GetPlayerToken(uid, "24"))
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

				printMemStats()

				lo.Must(coll.InsertMany(context.TODO(), bsonDocs))
			}
			bsonDocs = []interface{}{}
		}
		bsonDocs = append(bsonDocs, bsondoc)
		//lo.Must(coll.InsertOne(context.TODO(), bsondoc))
		ps.Set("id", sid)
	}
}

func gamex3(i int) {
	uid := fmt.Sprintf("24tdxx3game%d", i)

	pinfo := lo.Must(pg.GetPlayerToken(uid, "24"))

	coll := db.Collection(fmt.Sprintf("pgSpinDataGamex%d_%d", 3, i))
	ps := url.Values{}
	ps.Add("cs", "1")
	ps.Add("ml", "1")
	ps.Add("pf", "4")
	ps.Add("id", pinfo.Sid)
	ps.Add("wk", "0_C")
	ps.Add("btt", "1")
	ps.Add("atk", pinfo.Tk)

	// ps.Add("fb", "2")
	ps.Add("bn", "3")

	var bsonDocs []interface{}
	for ; ; time.Sleep(time.Second * 2) {
		var spinret struct {
			Si json.RawMessage
			// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.14.0/bson#NewDecoder
		}

		spinurl := lo.Must(url.JoinPath(pinfo.Geu, "v2/Spin"))

		if err := pg.InvokePGService(spinurl, ps, &spinret); err != nil {
			slog.Error("spin buy game error", "error", err, "sid",
				pinfo.Sid, "token", pinfo.Tk, "uid", uid)
			time.Sleep(5 * time.Second)
			// fmt.Println("gamex3, 开始重试。。。")
			pinfo = lo.Must(pg.GetPlayerToken(uid, "24"))
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

				printMemStats()

				lo.Must(coll.InsertMany(context.TODO(), bsonDocs))
			}
			bsonDocs = []interface{}{}
		}
		bsonDocs = append(bsonDocs, bsondoc)
		//lo.Must(coll.InsertOne(context.TODO(), bsondoc))
		ps.Set("id", sid)
	}
}
