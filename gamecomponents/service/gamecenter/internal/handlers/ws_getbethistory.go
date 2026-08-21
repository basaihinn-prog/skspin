package handlers

import (
	"context"
	"encoding/json"
	"game/comm/db"
	"game/comm/mux"
	"game/comm/slotsmongo"
	"game/duck/lang"
	"game/service/gamecenter/internal/gamedata"
	"game/service/gamecenter/internal/operator"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	mux.RegMsgProc("/player/getBetHistoryUrl", "获取下注历史启动URL", "player", getBetHistoryUrl, getBetHistoryPs{"YingCaiShen"})
	mux.RegHttpWithSample("/api/bethistory/list", "获取下注历史列表", "/api/bethistory", httpGetBetHistory, httpGetBetHistoryPs{"TuZi", 100060, 0, 10})
	mux.RegHttpWithSample("/api/bethistory/details", "获取下注历史详情", "/api/bethistory", httpGetSpinDetails, getSpinDetailsPs{BetID: lo.Must(primitive.ObjectIDFromHex("6572e58560211e96db1625b9"))})
	mux.RegHttpWithSample("/api/getlang", "获取多语言", "/api/lang", httpGetLang, nil)
}

type getBetHistoryPs struct {
	GameID string
	// Language string
}

type getBetHistoryRet struct {
	Url string
}

func getBetHistoryUrl(ctx *mux.WSReqContext, ps getBetHistoryPs, ret *getBetHistoryRet) (err error) {
	// u := "https://bethistory.kafa010.com/"
	pid := ctx.Pid
	u, err := url.Parse(gamedata.Get().GameHistoryUrl)
	if err != nil {
		return
	}

	// if ps.Language == "" {
	// 	ps.Language = "th"
	// }

	query := u.Query()
	query.Set("pid", strconv.Itoa(int(pid)))
	query.Set("game", ps.GameID)
	query.Set("l", ctx.Language)

	u.RawQuery = query.Encode()
	ret.Url = u.String()
	return
}

type httpGetBetHistoryPs struct {
	GameID    string
	Pid       int64
	PageIndex int64
	PageSize  int64
}
type httpGetBetHistoryRet struct {
	Count        int64
	BetTotal     int64
	WinTotal     int64
	WinLoseTotal int64
	List         []*slotsmongo.DocBetLog
}

func httpGetBetHistory(req *http.Request, ps httpGetBetHistoryPs, ret *httpGetBetHistoryRet) (err error) {

	if ps.PageIndex < 0 {
		ps.PageIndex = 0
	}

	if ps.PageSize < 0 {
		ps.PageSize = 1
	}
	if ps.PageSize > 1000 {
		ps.PageSize = 1000
	}

	plr, err := operator.AppMgr.GetPlr(ps.Pid)
	if err != nil {
		return
	}

	// coll := db.Collection2("reports", "BetLog")
	coll := slotsmongo.GetBetLogColl(plr.AppID)

	earlier := time.Now().Truncate(time.Hour).Add(-7 * 24 * time.Hour)
	firstId := primitive.NewObjectIDFromTimestamp(earlier)
	filter := db.D(
		"_id", db.D("$gt", firstId),
		"GameID", ps.GameID,
		"Pid", ps.Pid,
	)

	opts := options.Find().
		SetSort(db.D("_id", -1)).
		SetLimit(ps.PageSize).
		SetSkip(ps.PageIndex * ps.PageSize).
		SetProjection(db.D("SpinDetailsJson", 0))
	cur, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}

	// ret.List = []*slotsmongo.DocBetLog{}
	err = cur.All(context.TODO(), &ret.List)
	if err != nil {
		return
	}

	ret.Count, err = coll.CountDocuments(context.TODO(), filter, options.Count().SetLimit(10000))

	if ret.Count != 0 {
		// 	{
		//   _id: "$UserID",
		//   Bet: {
		//     $sum: "$Bet",
		//   },
		//   Win: {
		//     $sum: "$Win",
		//   },
		//   WinLose: {
		//     $sum: "$WinLose",
		//   },

		// _id "testuser1"
		// Bet 1090000
		// Win 4763500
		// WinLose 3673500

		cursor, er := coll.Aggregate(context.TODO(), mongo.Pipeline{
			db.D("$match", filter),
			db.D("$group", db.D(
				"_id", "$UserID",
				"Bet", db.D("$sum", "$Bet"),
				"Win", db.D("$sum", "$Win"),
				"WinLose", db.D("$sum", "$WinLose"),
			)),
		})

		if er == nil {

			type doc struct {
				Bet     int64 `bson:"Bet"`     // 下注
				Win     int64 `bson:"Win"`     // 输赢
				WinLose int64 `bson:"WinLose"` // 玩家输赢金额
			}
			var docs []*doc

			cursor.All(context.TODO(), &docs)
			if len(docs) != 0 {
				doc := docs[0]
				ret.BetTotal = doc.Bet
				ret.WinTotal = doc.Win
				ret.WinLoseTotal = doc.WinLose
			}
		}
	}

	return
}

type getSpinDetailsPs struct {
	BetID primitive.ObjectID
	Pid   int64
}

func httpGetSpinDetails(req *http.Request, ps getSpinDetailsPs, ret *json.RawMessage) (err error) {
	coll := db.Collection2("reports", "BetLog")
	if ps.Pid != 0 {
		plr, er := operator.AppMgr.GetPlr(ps.Pid)
		if er == nil {
			coll = slotsmongo.GetBetLogColl(plr.AppID)
		}
	}

	var doc struct {
		ID              primitive.ObjectID `bson:"_id"`
		SpinDetailsJson string             `bson:"SpinDetailsJson"`
	}
	err = coll.FindOne(context.TODO(), db.ID(ps.BetID)).Decode(&doc)
	if err != nil {
		return
	}

	*ret = json.RawMessage(doc.SpinDetailsJson)
	return
}

type GetLangListParams struct {
	Permission int32
}

type GetLangListResults struct {
	List []*lang.Lang
}

func httpGetLang(req *http.Request, ps GetLangListParams, ret *GetLangListResults) (err error) {
	ret.List = lang.GetAllArr()
	return err
}
