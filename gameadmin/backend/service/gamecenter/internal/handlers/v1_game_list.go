package handlers

import (
	"context"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/mux"
	"game/comm/ut"
	"game/service/gamecenter/internal/gamedata"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/operator"
	"log/slog"
	"sort"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:    "/api/v1/game/list",
		Handler: v1_game_list,
		Desc:    "获取游戏列表",
		Kind:    "api/v1",
		// ParamsSample: v1GameLaunchPs{"operator_user_abcd", "XingYunXiang", "th"},
		Class:   "operator",
		GetArg0: getArg0,
	})

	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:    "/api/v1/game/list2",
		Handler: v1_game_list_2,
		Desc:    "获取游戏列表-包含hide状态的游戏",
		Kind:    "api/v1",
		// ParamsSample: v1GameLaunchPs{"operator_user_abcd", "XingYunXiang", "th"},
		Class:   "operator",
		GetArg0: getArg0,
	})
}

type gamelistPs struct {
}

type v1Game struct {
	ID   string `json:"ID" bson:"_id"`
	Name string `json:"Name" bson:"Name"`
	Type int    `json:"Type" bson:"Type"`
	// Status  int    `json:"Status,omitempty" bson:"Status"`
	IconUrl    string //`json:"IconUrl" bson:"IconUrl"`
	WebIconUrl string
	Bet        int64 `json:"-" bson:"Bet"`
}

type v1GamelistRet struct {
	List []*v1Game
}

type GameBetUpdate struct {
	lock       sync.Mutex       // 互斥锁
	Games      map[string]int64 // 游戏打码量
	NextUpdate int64            // 下次更新时间戳
}

var gameUpdate = &GameBetUpdate{
	Games: map[string]int64{},
}

type BetInfo struct {
	Game string `bson:"_id"`
	Bet  int64  `bson:"betamount"`
}

type excludeGames struct {
	ID             primitive.ObjectID `bson:"_id"`
	Name           string             `json:Name,bson:Name`
	Remark         string             `json:Remark,bson:Remark`
	ExcluedGameIds []string           `json:ExcluedGameIds,bson:ExcluedGameIds`
}

func v1_game_list(app *operator.MemApp, ps mux.EmptyParams, ret *v1GamelistRet) (err error) {
	excluedGameIds := getExcluedGameIds(app.AppID)
	filter := bson.M{"Status": 0}
	if nil != excluedGameIds && len(excluedGameIds) > 0 {
		filter["_id"] = bson.M{"$nin": excluedGameIds}
	}
	coll := gcdb.CollGames.Coll()
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return
	}
	err = cur.All(context.TODO(), &ret.List)
	if err != nil {
		return
	}
	//gameBet := getGameBetInfo()

	for _, game := range ret.List {
		game.IconUrl = fmt.Sprintf("%s%s%s%s", gamedata.Settings.IconAddr, "icon/app/", game.ID, ".png")
		game.WebIconUrl = fmt.Sprintf("%s%s%s%s", gamedata.Settings.IconAddr, "icon/web/", game.ID, ".png")
	}
	sort.Slice(ret.List, func(i, j int) bool {
		return ret.List[i].Bet > ret.List[j].Bet
	})

	return
}

type v1Game2 struct {
	ID         string `json:"ID" bson:"_id"`
	Name       string `json:"Name" bson:"Name"`
	Type       int    `json:"Type" bson:"Type"`
	Status     int    `json:"Status" bson:"Status"`
	IconUrl    string //`json:"IconUrl" bson:"IconUrl"`
	WebIconUrl string
}
type v1GamelistRet2 struct {
	List []*v1Game2
}

func v1_game_list_2(app *operator.MemApp, ps mux.EmptyParams, ret *v1GamelistRet2) (err error) {
	coll := gcdb.CollGames.Coll()

	excluedGameIds := getExcluedGameIds(app.AppID)
	filter := bson.M{"Status": db.D("$lt", GameStatus_Hide)}
	if nil != excluedGameIds && len(excluedGameIds) > 0 {
		filter["_id"] = bson.M{"$nin": excluedGameIds}
	}
	cur, err := coll.Find(context.TODO(), filter, options.Find().SetSort(db.D("$natural", -1)))
	if err != nil {
		return
	}

	err = cur.All(context.TODO(), &ret.List)
	for _, game := range ret.List {
		game.IconUrl = fmt.Sprintf("%s%s%s%s", gamedata.Settings.IconAddr, "icon/app/", game.ID, ".png")
		game.WebIconUrl = fmt.Sprintf("%s%s%s%s", gamedata.Settings.IconAddr, "icon/web/", game.ID, ".png")
	}
	return
}

func getGameBetInfo() map[string]int64 {
	gameUpdate.lock.Lock()
	defer gameUpdate.lock.Unlock()
	now := time.Now()
	if gameUpdate.NextUpdate < now.Unix() {
		coll := db.Collection2("reports", "BetDailyReport")
		match := bson.M{
			"date": bson.M{
				"$gte": now.AddDate(0, 0, -6).Format("20060102"),
				"$lte": now.Format("20060102"),
			},
		}
		cursor, err := coll.Aggregate(context.TODO(), []bson.M{
			{
				"$match": match,
			}, {
				"$group": bson.M{
					"_id":       "$game",
					"betamount": bson.M{"$sum": "$betamount"},
				},
			},
		})
		if err != nil {
			return map[string]int64{}
		}
		res := make([]BetInfo, 0, 10)
		err = cursor.All(context.TODO(), &res)
		if err != nil {
			slog.Warn("db unmarshal", "error", ut.ErrString(err))
			return map[string]int64{}
		}
		if len(res) == 0 {
			return map[string]int64{}
		}
		gameUpdate.Games = map[string]int64{}
		for i := range res {
			gameUpdate.Games[res[i].Game] = res[i].Bet
		}
		gameUpdate.NextUpdate = now.Add(24 * time.Hour).Unix()
	}
	return gameUpdate.Games
}

func getExcluedGameIds(AppID string) []string {
	var ids []string
	o := &comm.Operator{}
	coll := db.Collection2("GameAdmin", "AdminOperator")
	coll.FindOne(context.TODO(), bson.M{"AppID": AppID}).Decode(&o)
	if o.ExcluedGameId.IsZero() {
		ids = o.ExcluedGameIds
	} else {
		excG := &excludeGames{}
		db.Collection2("GameAdmin", "ExcludeGame").FindOne(context.TODO(), bson.M{"_id": o.ExcluedGameId}).Decode(excG)
		ids = excG.ExcluedGameIds
	}
	return ids
}
