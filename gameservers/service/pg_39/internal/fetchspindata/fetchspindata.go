package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"serve/comm/db"
	"serve/comm/lazy"

	"serve/comm/plats/pg"

	"github.com/samber/lo"
)

func main() {
	mongoaddr := "mongodb://myUser:123456@192.168.1.2:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_39")

	// uid := "testuserzy"

	for i := 0; i < 1; i++ {
		uid := fmt.Sprintf("testuser%s%d", "zy", i)
		go fetch(uid)
	}

	lazy.SignalProc()
}

func fetch(uid string) {
	pinfo := lo.Must(pg.GetPlayerToken(uid, "39"))

	// coll := db.Collection("pgSpinData")

	ps := url.Values{}
	ps.Add("cs", "0.5")
	ps.Add("ml", "1")
	ps.Add("pf", "1")
	ps.Add("id", pinfo.Sid)
	ps.Add("wk", "0_C")
	ps.Add("btt", "1")
	ps.Add("atk", pinfo.Tk)
	for ; ; time.Sleep(time.Second) {
		var spinret struct {
			Si json.RawMessage
			// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.14.0/bson#NewDecoder
		}

		spinurl := lo.Must(url.JoinPath(pinfo.Geu, "v2/Spin"))

		if err := pg.InvokePGService(spinurl, ps, &spinret); err != nil {
			slog.Error("spin", "error", err)
			continue
		}
		// os.Stdout.Write(spinret)

		bsondoc := lo.Must(db.Json2Bson(spinret.Si))

		if bsondoc.Lookup("aw").Double() == 0 {
			// lo.Must(coll.InsertOne(context.TODO(), bsondoc))
		}

		sid := bsondoc.Lookup("sid").StringValue()
		ps.Set("id", sid)
	}

}
