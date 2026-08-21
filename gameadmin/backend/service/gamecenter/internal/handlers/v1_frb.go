package handlers

import (
	"context"
	"errors"
	"game/comm/db"
	"game/comm/mux"
	"game/service/gamecenter/internal/operator"
	"log/slog"
	"sync"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:    "/api/v1/freeRoundBonus/add",
		Handler: v1_frb_add,
		Desc:    "添加免费回合奖励",
		Kind:    "api/v1/freeRoundBonus",
		ParamsSample: v1FRBAddPs{
			BonusCode:      "abc101",
			GameID:         "pp_vs20olympx",
			PlayerList:     []string{"testuser1", "testuser2"},
			TotalBet:       1.0,
			Rounds:         10,
			ExpirationDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		},
		Class:   "operator",
		GetArg0: getArg0,
	})

	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:    "/api/v1/freeRoundBonus/query",
		Handler: v1_frb_query,
		Desc:    "查询免费回合奖励",
		Kind:    "api/v1/freeRoundBonus",
		ParamsSample: v1FRBQueryPs{
			BonusCode: "abc101",
			GameID:    "pp_vs20olympx",
			UserID:    "user1",
		},
		Class:   "operator",
		GetArg0: getArg0,
	})

	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:    "/api/v1/freeRoundBonus/delete",
		Handler: v1_frb_del,
		Desc:    "删除免费回合奖励",
		Kind:    "api/v1/freeRoundBonus",
		ParamsSample: v1FRBDelPs{
			BonusCode: "abc101",
		},
		Class:   "operator",
		GetArg0: getArg0,
	})
}

func frbColl() *mongo.Collection {
	coll := db.Collection2("game", "freeRoundBonus")
	return coll
}

var (
	frbIdxOnce = sync.OnceFunc(func() {
		coll := frbColl()
		str, err := coll.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
			{

				Keys: db.D("AppID", 1,
					"UserID", 1,
					"GameID", 1,
				),
				Options: options.Index().SetUnique(false),
			},
			{
				Keys: db.D("AppID", 1,
					"BonusCode", 1,
					"GameID", 1,
					"UserID", 1,
				),
				Options: options.Index().SetUnique(true),
			},
		})
		slog.Info("create index", "coll", "game.freeRoundBonus", "indexname", str, "err", err)

	})
)

type v1FRBAddPs struct {
	BonusCode      string
	GameID         string
	TotalBet       float64
	Rounds         int
	ExpirationDate int64
	PlayerList     []string
}

type v1FRBAddRet struct {
	AddCount int
}

type FRBPlayer struct {
	ID             primitive.ObjectID `bson:"_id"`
	AppID          string             `bson:"AppID"`
	UserID         string             `bson:"UserID"`
	GameID         string             `bson:"GameID"`
	BonusCode      string             `bson:"BonusCode"`
	TotalBet       float64            `bson:"TotalBet"`
	Rounds         int                `bson:"Rounds"`
	ExpirationDate int64              `bson:"ExpirationDate"`
}

func v1_frb_add(app *operator.MemApp, ps v1FRBAddPs, ret *v1FRBAddRet) (err error) {
	// coll.InsertMany(context.TODO())
	frbIdxOnce()

	if len(ps.PlayerList) == 0 {
		err = errors.New("PlayerList is empty")
		return
	}

	now := time.Now()

	if ps.ExpirationDate < now.Unix() {
		err = errors.New("ExpirationDate is invalid")
		return
	}

	if ps.Rounds < 1 {
		err = errors.New("Rounds is invalid")
		return
	}
	if ps.TotalBet <= 0 {
		err = errors.New("TotalBet is invalid")
		return
	}
	if ps.BonusCode == "" {
		err = errors.New("BonusCode is invalid")
		return
	}
	if !isPPGame(ps.GameID) {
		err = errors.New("GameID is invalid")
		return
	}

	docs := lo.Map(ps.PlayerList, func(uid string, _ int) *FRBPlayer {
		return &FRBPlayer{
			ID:             primitive.NewObjectIDFromTimestamp(now),
			AppID:          app.AppID,
			UserID:         uid,
			GameID:         ps.GameID,
			BonusCode:      ps.BonusCode,
			TotalBet:       ps.TotalBet,
			Rounds:         ps.Rounds,
			ExpirationDate: ps.ExpirationDate,
		}
	})

	coll := frbColl()
	coll.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: db.D("AppID", 1,
			"UserID", 1,
			"GameID", 1,
		),
	})
	istRet, _ := coll.InsertMany(context.TODO(), lo.ToAnySlice(docs), options.InsertMany().SetOrdered(false))

	if istRet != nil {
		ret.AddCount = len(istRet.InsertedIDs)
	}
	return
}

type v1FRBQueryPs struct {
	UserID    string
	GameID    string
	BonusCode string
}

type v1FRBQueryRet struct {
	Players []*FRBPlayer
}

func v1_frb_query(app *operator.MemApp, ps v1FRBQueryPs, ret *v1FRBQueryRet) (err error) {
	frbIdxOnce()

	filter := primitive.M{
		"AppID": app.AppID,
	}

	if ps.UserID != "" {
		filter["UserID"] = ps.UserID
	}
	if ps.GameID != "" {
		filter["GameID"] = ps.GameID
	}
	if ps.BonusCode != "" {
		filter["BonusCode"] = ps.BonusCode
	}

	coll := frbColl()

	cur, err := coll.Find(context.TODO(), filter, options.Find().SetLimit(1000))
	if err != nil {
		return
	}

	err = cur.All(context.TODO(), &ret.Players)

	return
}

type v1FRBDelPs struct {
	ID        primitive.ObjectID
	UserID    string
	GameID    string
	BonusCode string
}

type v1FRBDelRet struct {
	DeletedCount int64
}

func v1_frb_del(app *operator.MemApp, ps v1FRBDelPs, ret *v1FRBDelRet) (err error) {
	frbIdxOnce()

	filter := primitive.M{
		"AppID": app.AppID,
	}

	if ps.UserID != "" {
		filter["UserID"] = ps.UserID
	}
	if ps.GameID != "" {
		filter["GameID"] = ps.GameID
	}
	if ps.BonusCode != "" {
		filter["BonusCode"] = ps.BonusCode
	}

	if !ps.ID.IsZero() {
		filter = primitive.M{
			"_id": ps.ID,
		}
	}

	coll := frbColl()

	delret, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		return
	}

	ret.DeletedCount = delret.DeletedCount
	return
}

func DelFRBPlayers() {
	coll := frbColl()

	now := time.Now()
	filter := db.D(
		"ExpirationDate", db.D("$gt", now.Unix()),
	)

	coll.DeleteMany(context.TODO(), filter)
}
