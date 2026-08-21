package staticproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/define"
	"game/duck/ut2"
	"game/duck/ut2/jwtutil"
	"game/service/ppgateway/ppcomm/html"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	checkedM sync.Map
)

func EnsureIndex(coll *mongo.Collection) {
	key := coll.Database().Name()
	if _, ok := checkedM.Load(key); ok {
		return
	}
	checkedM.Store(key, true)

	indexmodels := []mongo.IndexModel{
		{
			Keys: db.D("datetime", 1),
			// Options: options.Index().SetUnique(false),
		},
		{
			Keys: db.D("pid", 1),
		},
	}
	coll.Indexes().CreateMany(context.TODO(), indexmodels)
}

type HistoryData struct {
	RoundId        primitive.ObjectID    `json:"roundId" bson:"_id"`
	Pid            int64                 `json:"-" bson:"pid"`
	BetId          string                `json:"betId" bson:"betId"`
	DateTime       int64                 `json:"dateTime"`
	Bet            float64               `json:"bet"`
	Win            float64               `json:"win"`
	Balance        float64               `json:"balance"`
	Details        string                `json:"roundDetails,omitempty" bson:"details"`
	RoundDetails   []*HistoryDetailsData `json:"-" `
	Currency       string                `json:"currency"`
	CurrencySymbol string                `json:"currencySymbol"`
	Hash           string                `json:"hash"`
	SharedLink     string                `json:"sharedLink" bson:"sharedLink"`
	Token          string                `json:"token" bson:"token"`
	GameConfig     *ReplayGameConfig     `json:"gameConfig" bson:"gameConfig"`
	BaseBet        float64               `json:"baseBet" bson:"baseBet"`
}

type HistoryDetailsData struct {
	RoundId        primitive.ObjectID `json:"roundId"`
	BetId          string             `json:"betId" bson:"betId"`
	ConfigHash     string             `json:"configHash"`
	Request        map[string]string  `json:"request"`
	Response       Variables          `json:"response"`
	Currency       string             `json:"currency"`
	CurrencySymbol string             `json:"currencySymbol"`
}

type PPBetHistoryRsp struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	RoundID    string             `json:"roundID" bson:"roundID"`
	Bet        float64            `json:"bet" bson:"bet"`
	BaseBet    float64            `json:"base_bet" bson:"baseBet"`
	Win        float64            `json:"win" bson:"win"`
	Rtp        float64            `json:"rtp" bson:"rtp"`
	PlayedDate int64              `json:"playedDate" bson:"playedDate"`
	SharedLink string             `json:"sharedLink" bson:"sharedLink"`
}

type HistoryRsp struct {
	Description string             `json:"description" bson:"description"`
	Error       int64              `json:"error" bson:"error"`
	TopList     []*PPBetHistoryRsp `json:"topList" bson:"topList"`
}

type GetLinkRsp struct {
	Error       int    `json:"error"`
	Description string `json:"description"`
	SharedLink  string `json:"sharedLink"`
}

type replayDataRsp struct {
	Error       int       `json:"error"`
	Description string    `json:"description"`
	Init        string    `json:"init"`
	Log         []logData `json:"log"`
}

type logData struct {
	CR string `json:"cr"`
	SR string `json:"sr"`
}

func getWinList(w http.ResponseWriter, r *http.Request) {
	var (
		ret HistoryRsp
		err error
	)
	defer func() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		// 设置缓存控制
		w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		if err != nil {
			ret.Error = 1
		}
		jsondata, _ := json.Marshal(ret)
		w.Write(jsondata)
	}()
	err = r.ParseForm()
	if err != nil {
		return
	}

	tk := r.FormValue("mgckey")
	pid, gid, err := jwtutil.ParseTokenData(tk)
	if err != nil {
		err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	//ps := &PGParams{
	//	Path: r.URL.Path,
	//	//TraceId: traceId,
	//	Form: r.Form,
	//	Pid:  pid,
	//}
	//coll := db.Collection2("pp_"+gid, "BetHistory")
	coll := db.Collection2(gid, "BetHistory")

	EnsureIndex(coll)

	// dtt := min()
	filter := db.D(
		"datetime", db.D(
			//"$gt", ps.GetInt("dtf"),
			"$lte", time.Now().UnixMilli(),
		),
		"pid", pid,
		"$and", bson.A{
			bson.D{bson.E{Key: "$expr", Value: db.D("$gte", []interface{}{"$win", db.D("$multiply", []interface{}{"$baseBet", 10})})}},
			bson.D{{"token", bson.D{{"$exists", true}}}},
		},
	)

	opts := options.Find().SetLimit(int64(10)).SetSort(db.D("datetime", -1))
	slog.Info("getHistory Param: ", "dbName", coll.Database().Name(), "filter: ", filter, "opts: ", opts)
	cur, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}
	var bh []HistoryData
	err = cur.All(context.TODO(), &bh)
	if err != nil {
		return
	}
	rsp := make([]*PPBetHistoryRsp, len(bh))
	for i, history := range bh {
		rsp[i] = &PPBetHistoryRsp{}
		err = copier.Copy(rsp[i], &bh[i])
		if err != nil {
			slog.Error("copier.Copy err", "err", err)
			ret = HistoryRsp{
				Description: err.Error(),
				Error:       1,
			}
			return
		}
		if _, ok := history.RoundDetails[0].Request["c"]; !ok {
			slog.Error(" history.RoundDetails param err", "err", err)
			ret = HistoryRsp{
				Description: err.Error(),
				Error:       1,
			}
			return
		}
		if _, ok := history.RoundDetails[0].Request["l"]; !ok {
			slog.Error(" history.RoundDetails param err", "err", err)
			ret = HistoryRsp{
				Description: err.Error(),
				Error:       1,
			}
			return
		}
		c, _ := strconv.ParseFloat(history.RoundDetails[0].Response["c"], 64)
		l, _ := strconv.ParseFloat(history.RoundDetails[0].Response["l"], 64)
		if l == 0 {
			l, _ = strconv.ParseFloat(history.RoundDetails[0].Request["l"], 64)
		}
		rsp[i].BaseBet = c * l
		rsp[i].Rtp = rsp[i].Win / rsp[i].BaseBet
		rsp[i].RoundID = bh[i].BetId
		rsp[i].PlayedDate = bh[i].DateTime
	}
	ret = HistoryRsp{
		Description: "OK",
		Error:       0,
		TopList:     rsp,
	}

	return
}

func getShareLink(w http.ResponseWriter, r *http.Request) {
	var (
		ret GetLinkRsp
		err error
	)
	defer func() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		// 设置缓存控制
		w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		if err != nil {
			ret.Error = 1
		}
		jsondata, _ := json.Marshal(ret)
		w.Write(jsondata)
	}()
	err = r.ParseForm()
	if err != nil {
		return
	}
	ps := r.URL.Query()
	betId := r.FormValue("roundID")
	tk := ""
	if strings.Contains(r.URL.String(), "replayGame.do") {
		tk = r.FormValue("token")
	} else {
		tk = r.FormValue("mgckey")
	}

	_, gid, err := jwtutil.ParseTokenData(tk)
	if err != nil {
		err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	//ps := &PGParams{
	//	Path: r.URL.Path,
	//	//TraceId: traceId,
	//	Form: r.Form,
	//	Pid:  pid,
	//}
	coll := db.Collection2(gid, "BetHistory")
	//referer := r.Header.Get("Referer")
	//// 解析 URL
	//parsedUrl, err := url.Parse(referer)
	//if err != nil {
	//	fmt.Println("Error parsing URL:", err)
	//	return
	//}
	//queryParams := parsedUrl.Query()
	//sip := queryParams.Get("sip")
	history := HistoryData{}
	coll.FindOne(context.TODO(), bson.M{"betId": betId}).Decode(&history)

	mainHost := ut2.Domain(r.Host)
	sharedLink := fmt.Sprintf(`https://replay.%s/%s`, mainHost, "rpt/"+history.Token)

	update := bson.M{
		"sharedLink": sharedLink,
	}
	coll.UpdateOne(context.TODO(), bson.M{"betId": betId}, bson.M{"$set": update})
	coll2 := db.Collection2("game", "ppReplayTokenMap")
	update = bson.M{
		"lang": ps.Get("lang"),
	}
	_, err = coll2.UpdateOne(context.TODO(), bson.M{"betid": betId, "lang": bson.M{"$exists": false}}, bson.M{"$set": update})
	if err != nil {
		fmt.Println(err)
	}
	ret = GetLinkRsp{
		Error:       0,
		Description: "OK",
		SharedLink:  sharedLink,
	}

	return
}

func replayData(w http.ResponseWriter, r *http.Request) {
	var (
		ret replayDataRsp
		err error
	)
	defer func() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		// 设置缓存控制
		w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		if err != nil {
			ret.Error = 1
		}
		jsondata, _ := json.Marshal(ret)
		w.Write(jsondata)
	}()
	err = r.ParseForm()
	if err != nil {
		return
	}
	betId := r.FormValue("roundID")
	tk := r.FormValue("token")
	gid := ""
	init := ""
	if len(tk) < 10 {
		rpInfo := ppReplayTokenMapBody{}
		param := bson.M{"token": tk, "gid": gid}
		coll := db.Collection2("game", "ppReplayTokenMap")
		coll.FindOne(context.TODO(), param).Decode(&rpInfo)
		gid = rpInfo.Gid
		init = rpInfo.Init
		//gid = strings.Replace(gid, "pp_", "", 1)
	} else {
		_, gid, err = jwtutil.ParseTokenData(tk)
		if err != nil {
			err = define.NewErrCode("Invalid player session", 1302)
			return
		}
		rpInfo := ppReplayTokenMapBody{}
		param := bson.M{"betid": betId, "gid": gid}
		coll := db.Collection2("game", "ppReplayTokenMap")
		coll.FindOne(context.TODO(), param).Decode(&rpInfo)
		init = rpInfo.Init
	}
	coll := db.Collection2(gid, "BetHistory")
	filter := db.D("betId", betId)
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return
	}
	var bh []HistoryData
	err = cur.All(context.TODO(), &bh)
	if err != nil {
		return
	}
	log := make([]logData, 0)
	if len(bh) == 0 {
		slog.Error("replay data is null")
		return
	}
	for _, roundDetail := range bh[0].RoundDetails {
		cr := roundDetail.Request
		sr := roundDetail.Response
		temp := logData{
			CR: mapToString(cr),
			SR: mapToString(sr),
		}
		log = append(log, temp)
	}

	ret = replayDataRsp{
		Error:       0,
		Description: "OK",
		Init:        init,
		Log:         log,
	}
	return
}

func mapToString(m map[string]string) string {
	var sb []string
	for key, value := range m {
		sb = append(sb, key+"="+value)
	}
	return strings.Join(sb, "&")
}

func replayGame(w http.ResponseWriter, r *http.Request) {
	ps := r.URL.Query()
	mgckey := ""
	var rpInfo ppReplayTokenMapBody
	param := bson.M{}
	if ps.Get("roundID") != "" {
		mgckey = ps.Get("token")
		_, game, err := jwtutil.ParseTokenData(mgckey)
		if err != nil {
			return
		}
		param = bson.M{"betid": ps.Get("roundID"), "gid": game}
	} else {
		path := r.URL.Path
		split := strings.Split(path, "/rpt/")
		if len(split) != 2 {
			slog.Error("replay path error")
			http.Error(w, "replay path error", http.StatusInternalServerError)
			return
		}
		param = bson.M{"token": split[1]}
		mgckey = split[1]
	}
	coll := db.Collection2("game", "ppReplayTokenMap")
	err := coll.FindOne(context.TODO(), param).Decode(&rpInfo)
	if err != nil {
		slog.Error("Get ppReplayTokenMap err", "err: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file := `service/ppgateway/ppcomm/html/ppReplayTemplate.html`
	var body []byte
	if _, ok := comm.OldPPGame[strings.ReplaceAll(rpInfo.Gid, "pp_", "")]; ok {
		body = html.PPReplayTemplateHtml
	} else {
		body = html.PPReplayTemplateMonkeyHtml
		file = `service/ppgateway/ppcomm/html/ppReplayTemplateMonkey.html`
	}

	body = replayTemplateDo_hook(body, w, r, mgckey, rpInfo)
	// 设置缓存控制
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
}
