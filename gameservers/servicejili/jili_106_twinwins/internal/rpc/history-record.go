package rpc

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"serve/servicejili/jili_106_twinwins/internal/models"

	"serve/comm/db"
	"serve/comm/define"

	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/history/%s/get-history-record", `tw`), get_history_record)

	jiliut.RegRpc(fmt.Sprintf("/history/%s/get-single-round-log-summary/", `tw`), get_single_round_log_summary)

	jiliut.RegRpc(fmt.Sprintf("/history/%s/get-log-plate-info/", `tw`), get_log_plate_info)
}

var (
	onceCreate sync.Once
)

func createIndexes() {
	coll := db.Collection("BetHistory")
	coll.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: db.D("pid", 1),
	})
	coll.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: db.D("historyrecord.roundindex", 1),
	})

}

// https://history-api.kafa010.com/history/csh/get-history-record?EndRowIndex=10&LangId=en-US&LogIndexAsRoundIndex=false&Minutes=60&StartRowIndex=1
func get_history_record(ps *nats.Msg) (ret []byte, err error) {
	onceCreate.Do(createIndexes)

	/*
		EndRowIndex=10&LangId=en-US&LogIndexAsRoundIndex=false&Minutes=60&StartRowIndex=1

		EndRowIndex=20&LangId=en-US&LogIndexAsRoundIndex=false&Minutes=60&StartRowIndex=11

		EndRowIndex=1&LangId=en-US&LogIndexAsRoundIndex=false&RoundIndex=1718692563149146002&StartRowIndex=1

		Date=2024-06-20&EndRowIndex=10&LangId=en-US&LogIndexAsRoundIndex=false&StartRowIndex=1&TimeZoneOffsetMinutes=480

	*/

	ul, err := url.Parse(string(ps.Data))
	if err != nil {
		return
	}

	query, err := define.ParseQuery(ul.RawQuery)
	if err != nil {
		return
	}

	var (
		endRowIndex   = query.GetInt("EndRowIndex")
		startRowIndex = query.GetInt("StartRowIndex")
		minutes       = query.GetInt("Minutes")

		roundIndex = query.Get("RoundIndex")

		date                  = query.Get("Date")
		timeZoneOffsetMinutes = query.GetInt("TimeZoneOffsetMinutes")
	)

	// header := define.Header(ps.Header)

	token := strings.TrimPrefix(ps.Header.Get("Authorization"), "Bearer ")

	// pid, err := jwtutil.ParseToken(token)
	pid, err := jiliut.ParseAuthToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}

	slog.Info("get_history_record", "query", query.Encode(), "pid", pid)

	coll := db.Collection("BetHistory")

	var (
		filter any
		opts   *options.FindOptions
	)
	if roundIndex != "" {
		if startRowIndex != 1 {
			ret = jiliut.MarshalJsonReturn([]*models.HistoryRecord{})
			return
		}
		filter = db.D(
			"historyrecord.roundindex", roundIndex,
		)
		opts = options.Find().SetProjection(db.D("historyrecord", 1))
	} else if date != "" {
		// time.Local
		loc := time.FixedZone("", timeZoneOffsetMinutes*60)

		t1, er := time.ParseInLocation(time.DateOnly, date, loc)
		if er != nil {
			err = er
			return
		}

		// db.BetHistory.countDocuments({$and: [{_id: {$gt: ObjectId('6673d13f6932e6ba5a0492e2')}}, {pid: 100050}, {_id: {$lt: ObjectId('6673d1426932e6ba5a0492e5')}}]}  )

		filter = db.D("$and", bson.A{
			db.D("pid", pid),
			db.D("_id", db.D("$gt", primitive.NewObjectIDFromTimestamp(t1))),
			db.D("_id", db.D("$lt", primitive.NewObjectIDFromTimestamp(t1.Add(24*time.Hour)))),
		})

		opts = options.Find().
			SetSort(db.D("_id", -1)).
			SetLimit(int64(endRowIndex - startRowIndex + 1)).
			SetSkip(int64(startRowIndex - 1)).
			SetProjection(db.D("historyrecord", 1))

	} else {

		now := time.Now()
		start_id := primitive.NewObjectIDFromTimestamp(now.Add(-time.Duration(minutes) * time.Minute))

		filter = db.D(
			"_id", db.D("$gt", start_id),
			"pid", pid,
		)
		opts = options.Find().
			SetSort(db.D("_id", -1)).
			SetLimit(int64(endRowIndex - startRowIndex + 1)).
			SetSkip(int64(startRowIndex - 1)).
			SetProjection(db.D("historyrecord", 1))
	}

	cur, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}

	// var records []*models.HistoryRecord
	var docs []*models.HistoryDoc
	err = cur.All(context.TODO(), &docs)
	if err != nil {
		return
	}

	records := lo.Map(docs, func(doc *models.HistoryDoc, _ int) *models.HistoryRecord {
		return doc.HistoryRecord
	})

	// ret = []byte(`{"Code":0,"Message":"成功","Data":[]}`)

	ret = jiliut.MarshalJsonReturn(records)

	return
}

// https://history-api.kafa010.com/history/csh/get-single-round-log-summary/en-US/66729856d2ad078920d25965
func get_single_round_log_summary(ps *nats.Msg) (ret []byte, err error) {
	ul, err := url.Parse(string(ps.Data))
	if err != nil {
		return
	}

	opts := options.FindOne().SetProjection(db.D("singleroundlogsummaries", 1))

	coll := db.Collection("BetHistory")

	var doc models.HistoryDoc

	err = coll.FindOne(context.TODO(), bson.M{"historyrecord.roundindex": path.Base(ul.Path)}, opts).Decode(&doc)

	if err != nil {
		return
	}

	ret = jiliut.MarshalJsonReturn(doc.SingleRoundLogSummaries)
	return
}

// https://history-api.jlfafafa3.com/history/csh/get-log-plate-info/1718692563274366002/1718692563274376002

func get_log_plate_info(ps *nats.Msg) (ret []byte, err error) {
	// strings.TrimPrefix( ps.Subject, "/history/csh/get-log-plate-info/"	)
	ul, err := url.Parse(string(ps.Data))
	if err != nil {
		return
	}
	args := strings.Split(ul.Path, "/")
	t := lo.Must(lo.Nth(args, -2))
	id2 := lo.Must(lo.Nth(args, -1))
	opts := options.FindOne().SetProjection(db.D("logplateinfos."+id2, 1))
	coll := db.Collection("BetHistory")
	var doc models.HistoryDoc
	err = coll.FindOne(context.TODO(), bson.M{"historyrecord.roundindex": t}, opts).Decode(&doc)
	if err != nil {
		return
	}

	ret = jiliut.MarshalJsonReturn(doc.LogPlateInfos[id2])

	return
}
