package api

import (
	"context"
	"encoding/json"
	"fmt"
	"game/comm/db"
	"game/comm/define"
	"game/duck/ut2/jwtutil"
	"game/service/pggateway/pgcomm"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"sync"
)

// https://api.pg-demo.com/web-api/game-proxy/v2/BetHistory/Get?traceId=KKQMKE21
// gid=39&dtf=1710432000000&dtt=1711036799999&bn=1&rc=15&atk=D663F48D-040E-46B5-B2C8-3FDD51F66316&pf=1&wk=0_C&btt=1
// gid=39&dtf=1710950400000&dtt=1711036799999&bn=2&rc=15&lbt=1711005490264&atk=4EBA5DB7-222C-4898-8F7C-36522DE464EB&pf=1&wk=0_C&btt=1

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
			Keys: db.D("createdAt", 1),
			// Options: options.Index().SetUnique(false),
		},
		{
			Keys: db.D("userCode", 1),
		},
	}
	coll.Indexes().CreateMany(context.TODO(), indexmodels)
}

//type PPBHItem struct {
//	Id         primitive.ObjectID `json:"id" bson:"_id"`
//	Tid        string             `json:"tid" bson:"tid"`
//	CC         string             `json:"cc" bson:"cc"`
//	AgentCode  string             `json:"agentCode" bson:"agentCode"`
//	UserCode   string             `json:"userCode" bson:"userCode"`
//	GameCode   string             `json:"gameCode" bson:"gameCode"`
//	RoundID    string             `json:"roundID" bson:"roundID"`
//	Bet        float64            `json:"bet" bson:"bet"`
//	Win        float64            `json:"win" bson:"win"`
//	Rtp        float64            `json:"rtp" bson:"rtp"`
//	PlayedDate int64              `json:"playedDate" bson:"playedDate"`
//	Data       []*PPBDItem        `json:"data" bson:"data"`
//	SharedLink string             `json:"sharedLink" bson:"sharedLink"`
//	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
//}

type PPBHItem struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	Pid            int64              `json:"pid" bson:"pid"`
	Datetime       int64              `json:"datetime" bson:"datetime"`
	Bet            float64            `json:"bet" bson:"bet"`
	Win            float64            `json:"win" bson:"win"`
	Balance        float64            `json:"balance" bson:"balance"`
	Details        string             `json:"details" bson:"details"`
	Rounddetails   []*PPBDItem        `json:"rounddetails" bson:"rounddetails"`
	Currency       string             `json:"currency" bson:"currency"`
	Currencysymbol string             `json:"currencysymbol" bson:"currencysymbol"`
	Hash           string             `json:"hash" bson:"hash"`
}

type PPBDItem struct {
	CR string `json:"cr" bson:"cr"`
	SR string `json:"sr" bson:"sr"`
}

//type PPBHItemRsp struct {
//	Id         primitive.ObjectID `json:"id" bson:"_id"`
//	Tid        string             `json:"tid" bson:"tid"`
//	CC         string             `json:"cc" bson:"cc"`
//	AgentCode  string             `json:"agentCode" bson:"agentCode"`
//	UserCode   string             `json:"userCode" bson:"userCode"`
//	GameCode   string             `json:"gameCode" bson:"gameCode"`
//	RoundID    string             `json:"roundID" bson:"roundID"`
//	Bet        float64            `json:"bet" bson:"bet"`
//	Win        float64            `json:"win" bson:"win"`
//	Rtp        float64            `json:"rtp" bson:"rtp"`
//	PlayedDate int64              `json:"playedDate" bson:"playedDate"`
//	Data       string             `json:"data" bson:"data"`
//	SharedLink string             `json:"sharedLink" bson:"sharedLink"`
//	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
//}

type PPBHItemRsp struct {
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
	Description string        `json:"description" bson:"description"`
	Error       int64         `json:"error" bson:"error"`
	TopList     []PPBHItemRsp `json:"topList" bson:"topList"`
}

func getBetHistory(w http.ResponseWriter, r *http.Request) {
	var (
		ret HistoryRsp
		err error
	)
	defer func() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
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
	ps := &PGParams{
		Path: r.URL.Path,
		//TraceId: traceId,
		Form: r.Form,
		Pid:  pid,
	}
	coll := db.Collection2(gid, "BetHistory")

	EnsureIndex(coll)

	// dtt := min()
	filter := db.D(
		//"createdAt", db.D(
		//	//"$gt", ps.GetInt("dtf"),
		//	"$lt", time.Now().UnixMilli(),
		//),
		"userCode", fmt.Sprintf("%v", ps.Pid),
	)

	opts := options.Find().SetLimit(int64(10)).SetSort(db.D("createdAt", -1))
	cur, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}
	var bh []PPBHItem
	err = cur.All(context.TODO(), &bh)
	if err != nil {
		return
	}
	rsp := make([]PPBHItemRsp, len(bh))
	if len(bh) == 0 {
		bh = []PPBHItem{}
	} else {
		for i := range bh { //转base64
			//for i2 := range bh[i].Data {
			//	temp := bh[i].Data[i2]
			//	temp.CR = base64.StdEncoding.EncodeToString([]byte(temp.CR))
			//	temp.SR = base64.StdEncoding.EncodeToString([]byte(temp.SR))
			//}
			rsp[i] = PPBHItemRsp{}
			err := copier.Copy(&rsp[i], &bh[i])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			//marshal, _ := json.Marshal(bh[i].Data)
			//rsp[i].Data = string(marshal)
		}
	}

	ret = HistoryRsp{
		Description: "OK",
		Error:       0,
		TopList:     rsp,
	}

	return
}

// /back-office-proxy/Report/GetBetHistory
// sid=1775083500903469568&gid=1489936
func boGetBetHistory(ps *PGParams, ret *M) (err error) {
	sid := ps.Get("sid")
	id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		return
	}
	gid := ps.Get("gid")
	coll := db.Collection2("pg_"+gid, "BetHistory")

	var bh pgcomm.BHItem
	err = coll.FindOne(context.TODO(), db.ID(id)).Decode(&bh)
	if err != nil {
		return
	}

	// if len(bh.Bd) == 0 {
	// 	bh.Bd = []*pgcomm.BDItem{}
	// }
	if len(bh.CC) > 3 {
		bh.CC = fmt.Sprintf(" %s ", bh.CC[:3])
	}
	//bh.CC = bh.CC[:3]

	*ret = M{
		"bh": &bh,
	}
	return
}
