package slotsmongo

import (
	"context"
	"sync"

	"serve/comm/db"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocPlayer struct {
	Pid   int64  `bson:"_id"`
	Uid   string `bson:"Uid"`   // 外部id
	AppID string `bson:"AppID"` // 所属产品
	//Status int64  `bson:"Status"`
}

type Players struct {
	m   map[int64]*DocPlayer
	mtx sync.Mutex
}

type Player struct {
	db.DocPlayer `bson:"inline"`

	LS string

	// SpinCountOfThisEnter int
	// SpinCount            int
	// BonusGameCount       int
	// BetHistory           []int64
	// BetAmount            int64
	// WinAmount            int64

	Cs        float64
	Ml        float64
	LastSid   string           //上一次的Sid,判断参数的正确性
	BdRecords []map[string]any //转完了后需要增加历史纪律
	IsBuy     bool             //是否是购买小游戏
	BigReward int64            //是否是大奖
}

var players = &Players{
	m: map[int64]*DocPlayer{},
}

func (ps *Players) Get(pid int64) (plr *DocPlayer, err error) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	plr = ps.m[pid]
	if plr == nil {
		coll := db.Collection2("game", "Players")

		projection := db.D(
			"_id", 1,
			"Uid", 1,
			"AppID", 1,
		)

		var doc DocPlayer
		err = coll.FindOne(context.TODO(), db.ID(pid), options.FindOne().SetProjection(projection)).Decode(&doc)
		if err != nil {
			return
		}

		plr = &doc
		ps.m[pid] = plr
	}

	return
}

func (ps *Players) GetAppID(pid int64) string {
	plr, err := ps.Get(pid)
	if err != nil {
		return "unknown"
	}

	return plr.AppID
}

func GetPlayerInfo(pid int64) (appId, uid string, err error) {
	doc, err := players.Get(pid)
	if err != nil {
		return
	}
	appId, uid = doc.AppID, doc.Uid
	return
}
