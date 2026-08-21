package handlers

import (
	"game/comm/mux"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/operator"

	"go.mongodb.org/mongo-driver/bson"
)

// var (
// 	publicMux = mux.NewMux()
// )

func init() {
	mux.DefaultRpcMux.Add(&mux.PHandler{
		Path:         "/api/GameList",
		Handler:      getGamelist,
		Desc:         "获取游戏列表",
		Kind:         "api",
		ParamsSample: nil,
		Class:        "operator",
		GetArg0:      getArg0,
	})
}

type GameType int32

const (
	GameType_Slot  GameType = 0 // 拉霸游戏
	GameType_Fish  GameType = 1 // 捕鱼游戏
	GameType_Poker GameType = 3 // 棋牌游戏
)

type GameStatus int32

const (
	GameStatus_Open        GameStatus = 0 // 正常
	GameStatus_Maintenance GameStatus = 1 // 维护中（列表中出现）
	GameStatus_Hide        GameStatus = 2 // 隐藏（列表中不出现）
	GameStatus_Closed      GameStatus = 3 // 关闭（可以执行删除）
)

type Game struct {
	ID     string `bson:"_id"`
	Name   string
	Type   GameType
	Status GameStatus
}

type GameListResp struct {
	List []*Game
}

func getGamelist(app *operator.MemApp, ps mux.EmptyParams, ret *GameListResp) (err error) {
	err = gcdb.CollGames.FindAll(bson.M{"Status": 0}, &ret.List)
	return
}
