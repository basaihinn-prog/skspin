package handlers

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"game/comm/db"
	"game/comm/define"
	"game/comm/mux"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/mongodb"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/operator"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/record/betlist",
		Handler:      v1_record_bet_list,
		Desc:         "获取投注记录",
		Kind:         "api/v1",
		ParamsSample: betlistPs{Count: 10},
		Class:        "operator",
		GetArg0:      getArg0,
	})

	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/record/betlistcsv",
		Handler:      v1_record_bet_list_csv,
		Desc:         "获取投注记录 csv 格式",
		Kind:         "api/v1",
		ParamsSample: betlistPs{Count: 10},
		Class:        "operator",
		GetArg0:      getArg0,
	})

	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/v1/record/getbetlogs",
		Handler:      v1_record_get_bet_logs,
		Desc:         "获取投注记录",
		Kind:         "api/v1",
		ParamsSample: getBetLogsPs{UserId: "", GameId: "", StartTime: 0, EndTime: 0, Page: 1, PageSize: 10},
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type betlistPs struct {
	StartID string
	Count   int
}

type getBetLogsPs struct {
	UserId    string
	GameId    string
	StartTime int64
	EndTime   int64
	Page      int64
	PageSize  int64
}

type getOperatorListResult struct {
	List []*DocBetLog
	All  int64
}

type goldtype int64

func (g goldtype) MarshalJSON() ([]byte, error) {
	fg := ut.Gold2Money(int64(g))
	return json.Marshal(fg)
}

type DocBetLog struct {
	ID         string    `bson:"_id"`
	Pid        int64     `bson:"Pid"`        // 内部玩家ID
	UserID     string    `bson:"UserID"`     // 外部玩家ID
	GameID     string    `bson:"GameID"`     // 游戏ID
	Bet        goldtype  `bson:"Bet"`        // 下注
	Win        goldtype  `bson:"Win"`        // 输赢
	InsertTime time.Time `bson:"InsertTime"` // 数据插入时间
	AppID      string    `bson:"AppID"`      //
	Balance    goldtype  `bson:"Balance"`    // 余额
	WinLose    goldtype  `bson:"WinLose"`    // 玩家输赢金额
	Grade      int       `bson:"Grade"`
	Completed  bool      `bson:"Completed"` // 最后一次spin
	RoundID    string    `bson:"RoundID"`
}

type betlistRet struct {
	Title       []string
	List        [][]any
	NextStartID string
}

var (
	title = []string{
		"ID",
		"Pid",
		"UserID",
		"GameID",
		"Bet",
		"Win",
		"InsertTime",
		"AppID",
		"Balance",
		"WinLose",
		"Grade",
		"RoundID",
	}

	titleproj = db.D(
		"Pid", 1,
		"UserID", 1,
		"GameID", 1,
		"Bet", 1,
		"Win", 1,
		"InsertTime", 1,
		"AppID", 1,
		"Balance", 1,
		"WinLose", 1,
		"Grade", 1,
		"RoundID", 1,
	)
	titlekeys map[string]int
)

func init() {
	titlekeys = make(map[string]int, len(title))
	for i, t := range title {
		titlekeys[t] = i
	}
	// titleproj = make(primitive.D, len(title)-1)
	// for i := 1; i < len(title); i++ {
	// 	t := title[i]
	// 	titleproj[i-1].Key = t
	// 	titleproj[i-1].Value = 1
	// }
}

func v1_record_bet_list(app *operator.MemApp, ps betlistPs, ret *betlistRet) (err error) {
	if false {
		now := time.Now()
		ret.NextStartID = ps.StartID
		if now.Sub(app.LastPullHisTime) < 9*time.Second {
			err = define.NewErrCode("Too frequently", 1019)
			// ret.List = [][]any{}
			return
		}

		app.LastPullHisTime = now
	}

	ps.Count = ut.ClipInt(ps.Count, 1, 100000)

	// coll := db.Collection2("reports", "BetLog")

	coll := slotsmongo.GetBetLogColl(app.AppID)

	// {LogType: {$not: {$gt: 0}}}
	filter := db.D(
		"AppID", app.AppID,
		"_id", db.D("$gt", ps.StartID),
		"LogType", db.D("$not", db.D("$gt", 0)),
	)

	opts := options.Find().SetProjection(titleproj).SetLimit(int64(ps.Count))
	cur, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}

	var betlogs []*DocBetLog
	err = cur.All(context.TODO(), &betlogs)
	if err != nil {
		return
	}

	ret.Title = title
	ret.List = make([][]any, len(betlogs))

	for i, betlog := range betlogs {
		row := make([]any, len(title))
		row[titlekeys["ID"]] = betlog.ID
		row[titlekeys["Pid"]] = betlog.Pid
		row[titlekeys["UserID"]] = betlog.UserID
		row[titlekeys["GameID"]] = betlog.GameID
		row[titlekeys["Bet"]] = betlog.Bet
		row[titlekeys["Win"]] = betlog.Win
		row[titlekeys["InsertTime"]] = betlog.InsertTime
		row[titlekeys["AppID"]] = betlog.AppID
		row[titlekeys["Balance"]] = betlog.Balance
		row[titlekeys["WinLose"]] = betlog.WinLose
		row[titlekeys["Grade"]] = betlog.Grade
		// row[titlekeys["Round"]] = cmp.Or(betlog.PGBetID, betlog.ID.Hex())
		row[titlekeys["RoundID"]] = betlog.RoundID

		ret.List[i] = row
	}

	if n := len(betlogs); n != 0 {
		ret.NextStartID = betlogs[n-1].ID
	}

	return
}

type betlistCSVRet = string

func v1_record_bet_list_csv(app *operator.MemApp, ps betlistPs, ret *betlistCSVRet) (err error) {
	now := time.Now()
	// if now.Sub(app.LastPullHisTime) < 9*time.Second {
	// 	err = errors.New("Too frequently")
	// 	return
	// }
	app.LastPullHisTime = now

	ps.Count = ut.ClipInt(ps.Count, 1, 100000)

	// coll := db.Collection2("reports", "BetLog")
	coll := slotsmongo.GetBetLogColl(app.AppID)

	filter := db.D(
		"AppID", app.AppID,
		"_id", db.D("$gt", ps.StartID),
		"LogType", db.D("$not", db.D("$gt", 0)),
	)

	opts := options.Find().SetProjection(titleproj).SetLimit(int64(ps.Count))
	cur, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}

	var betlogs []*DocBetLog
	err = cur.All(context.TODO(), &betlogs)
	if err != nil {
		return
	}

	// w.
	// ret.Title = title
	list := make([][]string, len(betlogs)+1)
	list[0] = title

	for i, betlog := range betlogs {
		row := make([]string, len(title))
		row[titlekeys["ID"]] = betlog.ID
		row[titlekeys["Pid"]] = strconv.Itoa(int(betlog.Pid))
		row[titlekeys["UserID"]] = betlog.UserID
		row[titlekeys["GameID"]] = betlog.GameID
		row[titlekeys["Bet"]] = fmtGold(betlog.Bet)
		row[titlekeys["Win"]] = fmtGold(betlog.Win)
		row[titlekeys["InsertTime"]] = betlog.InsertTime.Format(time.RFC3339)
		row[titlekeys["AppID"]] = betlog.AppID
		row[titlekeys["Balance"]] = fmtGold(betlog.Balance)
		row[titlekeys["WinLose"]] = fmtGold(betlog.WinLose)
		row[titlekeys["Grade"]] = strconv.Itoa(betlog.Grade)
		// row[titlekeys["Round"]] = cmp.Or(betlog.PGBetID, betlog.ID.Hex())
		row[titlekeys["RoundID"]] = betlog.RoundID

		list[i+1] = row
	}

	var w bytes.Buffer
	csv.NewWriter(&w).WriteAll(list)

	*ret = w.String()

	return
}

func fmtGold(gold goldtype) string {
	money := ut.Gold2Money(int64(gold))
	return fmt.Sprintf("%.4f", money)
}
func v1_record_get_bet_logs(app *operator.MemApp, ps getBetLogsPs, ret *getOperatorListResult) (err error) {
	query := bson.M{}
	if err != nil {
		return
	}
	var startT, endT time.Time
	RDate := ps.EndTime + ps.StartTime
	switch RDate {
	case 0:
	case ps.EndTime:
		endT = time.Unix(ps.EndTime, 0)
		query["InsertTime"] = bson.M{
			"$lte": endT,
		}
	case ps.StartTime:
		startT = time.Unix(ps.StartTime, 0)
		query["InsertTime"] = bson.M{
			"$gte": startT,
		}
	default:
		startT = time.Unix(ps.StartTime, 0)
		endT = time.Unix(ps.EndTime, 0)
		query["InsertTime"] = bson.M{
			"$gte": startT,
			"$lte": endT,
		}
	}

	filter := mongodb.FindPageOpt{
		Page:     ps.Page,
		PageSize: ps.PageSize,
		Sort:     db.D("InsertTime", -1),
	}

	if ps.GameId != "" {
		query["GameID"] = ps.GameId
	}
	if ps.UserId != "" {
		query["UserID"] = ps.UserId
	}
	query["LogType"] = 0

	filter.Query = query
	ret.All, err = gcdb.NewOtherDB("betlog").Collection(app.AppID).FindPage(filter, &ret.List)
	if err != nil {
		return err
	}

	return
}
