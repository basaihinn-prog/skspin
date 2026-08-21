package gcdb

import (
	"context"
	"game/comm/db"
	"game/duck/lazy"
	"game/duck/logger"
	"game/duck/mongodb"
	"log/slog"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB = mongodb.NewDB("game")

	CollPlayers = DB.Collection("Players") // 玩家表
	// CollApps          = DB.Collection("Apps")          // 接入的app
	CollStore         = DB.Collection("Store")         // kv存储
	CollModifyGoldLog = DB.Collection("ModifyGoldLog") // 金币修改日志
	CollGames         = DB.Collection("Games")         // 游戏列表

	// CollAppInfoMonth  = DB.Collection("AppInfoMonth")  // 产品统计，月
	// CollAppInfoDay    = DB.Collection("AppInfoDay")    // 产品统计，日
	// CollGameInfoMonth = DB.Collection("GameInfoMonth") // 游戏统计，月
	// CollGameInfoDay   = DB.Collection("GameInfoDay")   // 游戏统计，日
	// CollPlayerInfoDay = DB.Collection("PlayerInfoDay") // 玩家统计，日
)

func EnsureMongoIndex() {
	// CollPlayers.CreateIndexByKeys("Uid", "AppID")

	{
		indexes := []mongo.IndexModel{
			{
				Keys:    db.D("Uid", 1, "AppID", 1),
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: db.D("Uid", 1),
			},
			{
				Keys: db.D("AppID", 1),
			},
		}
		CollPlayers.Coll().Indexes().CreateMany(context.TODO(), indexes)
	}
	// pid,gameid 用于后台查询。
	// CollBetLog.CreateIndexByKeys("AppID", "InsertTime", "Pid", "GameID")

	// CollAppInfoMonth.CreateIndexByKeys("AppID")
	// CollAppInfoDay.CreateIndexByKeys("AppID")
	// CollGameInfoMonth.CreateIndexByKeys("GameID")
	// CollGameInfoDay.CreateIndexByKeys("GameID")
	// CollPlayerInfoDay.CreateIndexByKeys("Pid")

	{
		coll := db.Collection2("reports", "BetLog")
		mongodb.CreateIndexByKeys(coll, "AppID", "InsertTime", "Pid", "GameID", "UserID")
	}

	{
		coll := db.Collection2("game", "orders")
		name, err := coll.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
			Keys:    db.D("TraceId", 1, "AppID", 1),
			Options: options.Index().SetUnique(true),
		})
		slog.Info("create indexs", "coll", "orders", "idxname", name, "error", err)
	}

}

func ConnectMongoByConfig() {
	mongoAddr := lo.Must(lazy.RouteFile.Get("mongo"))

	logger.Info("连接mongodb", mongoAddr)

	// err := DB.Connect("mongodb://" + mongoAddr)
	err := DB.Connect(mongoAddr)

	if err != nil {
		logger.Info("连接mongodb失败", err)
	} else {
		logger.Info("连接mongodb成功")
		db.SetupClient(DB.Client)
	}
}

const (
	key_idCounter = 1008111
	// key_analysisCursor = "analysisCursor"
)

func NextID(targetColl *mongodb.Collection, min int64) (int64, error) {

	op := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	mp := map[string]int64{}
	err := CollStore.Coll().FindOneAndUpdate(context.TODO(),
		bson.M{
			"_id": key_idCounter,
		},
		bson.M{
			"$inc": bson.M{
				targetColl.Name: 1,
			},
		}, op).Decode(&mp)

	if err != nil {
		return 0, err
	}

	return mp[targetColl.Name] + min, err
}

func NewOtherDB(dbName string) *mongodb.DB {
	return &mongodb.DB{Name: dbName, Client: DB.Client}
}
