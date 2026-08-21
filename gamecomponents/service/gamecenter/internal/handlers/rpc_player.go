package handlers

import (
	"errors"
	"fmt"
	"game/comm/mux"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/operator"
	"game/service/gamecenter/msgdef"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	mux.RegRpc("/gamecenter/player/getBalance", "获取玩家余额", "player", getBalance, getPs{100001})
	mux.RegRpc("/gamecenter/player/modifyGold", "修改玩家余额", "player", modifyGold, modifyGoldReq{
		Pid:    100001,
		Change: 10000,
		// Comment: "http 测试修改",
	})
	mux.RegRpc("/gamecenter/player/getGain", "获取玩家盈利", "player", getGain, getPs{100001})
	mux.RegRpc("/gamecenter/player/getPlayerInGame", "获取游戏内玩家", "player", getPlayerInGame, getPlayerInGamePs{100001, "pg_89", "faketrans"})
}

// type GameRpcPlayer struct{}

type getPs struct {
	Pid int64
}

type getPlayerInGamePs struct {
	Pid    int64
	GameId string
	AppId  string
}
type getBalanceRet struct {
	Balance int64
	Uid     string
}

type getGainRet struct {
	Gain int64
	Uid  string
}

type DocPlayer struct {
	PID         int64  `bson:"_id"`
	Uid         string `bson:"Uid"`         // 外部id
	AppID       string `bson:"AppID"`       // 所属产品
	CurrencyKey string `bson:"CurrencyKey"` // 玩家货币key
	Language    string `bson:"-"`           // 语言 th,en,cn
	//CurrentRtp  int    `bson:"CurrentRtp"`  // RTP
	Status int64 `bson:"Status"` // 玩家状态

	SpinCountOfThisEnter int                `bson:"spincountofthisenter"`
	SpinCount            int                `bson:"spincount"`
	BonusGameCount       int                `bson:"bonusgamecount"` // 小游戏数量
	BetHistory           []int64            `bson:"bethistory"`
	BetAmount            int64              `bson:"betamount"`
	WinAmount            int64              `bson:"winamount"`
	BetHistLastID        primitive.ObjectID `bson:"bethistlastid"`

	RewardPercent  int `bson:"reward_percent"`
	NoAwardPercent int `bson:"no_award_percent"`
}

// for Admin
func getPlayerInGame(req getPlayerInGamePs, doc *DocPlayer) error {
	err := gcdb.NewOtherDB(req.GameId).Collection("players").FindOne(bson.M{"_id": req.Pid, "AppID": req.AppId}, &doc)
	if err != nil {
		return err
	}
	return nil
}

func getBalance(req getPs, ret *getBalanceRet) (err error) {
	am := operator.AppMgr
	plr, err := am.GetPlr(req.Pid)
	if err != nil {
		return
	}
	app := am.GetApp(plr.AppID)
	if app == nil {
		err = errors.New("app not found.")
		return
	}

	balance, err := app.GetBalance(plr)
	if err != nil {
		return
	}

	ret.Balance = balance
	ret.Uid = plr.Uid

	return
}

type modifyGoldReq = msgdef.ModifyGoldPs

// type modifyGoldReq struct {
// 	Pid     int64 // 玩家ID
// 	Change  int64 // 金币变化：>0 加钱， <0 扣钱
// 	SpinID  string
// 	RoundID string
// 	Reason  string
// 	Comment string // 注释
// }

func modifyGold(ps *msgdef.ModifyGoldPs, ret *msgdef.ModifyGoldRet) (err error) {
	if !ps.Valid() {
		err = fmt.Errorf("error ps[%+v]", ps)
		return
	}

	am := operator.AppMgr

	plr, err := am.GetPlr(ps.Pid)
	if err != nil {
		return
	}
	app := am.GetApp(plr.AppID)
	if app == nil {
		err = errors.New("app not found.")
		return
	}

	// balance, err := app.ModifyGold(plr, ps.Change, ps.Comment)
	balance, err := app.ModifyGold(plr, ps)
	if err != nil {
		return
	}
	ret.Balance = balance
	ret.Uid = plr.Uid

	return
}

func getGain(req getPs, ret *getGainRet) (err error) {
	am := operator.AppMgr
	plr, err := am.GetPlr(req.Pid)
	if err != nil {
		return
	}
	app := am.GetApp(plr.AppID)
	if app == nil {
		err = errors.New("app not found.")
		return
	}

	gain, err := app.GetGain(plr)
	if err != nil {
		return
	}

	ret.Gain = gain
	ret.Uid = plr.Uid

	return
}
