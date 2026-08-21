package operator

import (
	"context"
	"errors"
	"game/comm/db"
	"game/comm/define"
	"game/comm/ut"
	"game/duck/mongodb"
	"game/service/gamecenter/internal/gcdb"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocPlayer struct {
	Pid int64 `bson:"_id"`
	// Language string             //语言 th,en,cn
	Uid      string             // 外部id
	AppID    string             // 所属产品
	LoginAt  *mongodb.TimeStamp // 最后登录时间
	CreateAt *mongodb.TimeStamp // 创建时间
	// Bet            int64
	// Win            int64
	Balance        int64 // 余额
	CurrentRtp     int
	CompletedGames map[string]bool `bson:"CompletedGames,omitempty"`
}

type MemPlr struct {
	Pid   int64 `bson:"_id"`
	AppID string
	Uid   string
}

func loadPlr(filter bson.D) (memplr *MemPlr, err error) {
	coll := gcdb.CollPlayers.Coll()
	var one DocPlayer
	projection := db.D(
		"_id", 1,
		"Uid", 1,
		"AppID", 1,
	)

	err = coll.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&one)
	if err == mongo.ErrNoDocuments {
		err = define.NewErrCode("Player account does not exist", 6007)
	}

	if err == nil {
		memplr = &MemPlr{
			AppID: one.AppID,
			Pid:   one.Pid,
			Uid:   one.Uid,
		}
	}
	return

}

func loadPlrByUid(uid string, appId string) (memplr *MemPlr, err error) {
	filter := db.D("Uid", uid, "AppID", appId)
	return loadPlr(filter)
}
func loadPlrByPid(pid int64) (memplr *MemPlr, err error) {
	filter := db.ID(pid)
	return loadPlr(filter)
}

func (mp *MemPlr) Status() int {
	var status struct {
		Status int `bson:"Status"`
	}

	coll := gcdb.CollPlayers.Coll()
	coll.FindOne(context.TODO(), db.ID(mp.Pid), options.FindOne().SetProjection(db.D("Status", 1))).Decode(&status)

	return status.Status
}

func (mp *MemPlr) SetGold(gold int64) (balance int64, err error) {
	lo.Must0(gold >= 0)
	coll := gcdb.CollPlayers.Coll()
	// _, err = coll.UpdateByID(context.TODO(), mp.Pid, db.D("$inc", db.D("Balance", gold)))
	update := db.D("$set", db.D("Balance", gold))
	opts := options.FindOneAndUpdate().SetProjection(db.D("Balance", 1)).SetReturnDocument(options.Before)
	var doc DocPlayer
	err = coll.FindOneAndUpdate(context.TODO(), db.ID(mp.Pid), update, opts).Decode(&doc)
	if err != nil {
		return
	}
	balance = doc.Balance
	return
}

func (mp *MemPlr) SetRtp(rtp int) (old_rtp int, err error) {
	coll := gcdb.CollPlayers.Coll()

	update := db.D("$set", db.D("CurrentRtp", rtp))
	opts := options.FindOneAndUpdate().SetProjection(db.D("CurrentRtp", 1)).SetReturnDocument(options.Before)
	var doc DocPlayer
	err = coll.FindOneAndUpdate(context.TODO(), db.ID(mp.Pid), update, opts).Decode(&doc)
	if err != nil {
		return
	}
	old_rtp = doc.CurrentRtp
	return
}

func (mp *MemPlr) TransferIn(gold int64) (balance int64, err error) {
	lo.Must0(gold >= 0)

	if gold == 0 {
		return mp.Balance()
	}

	coll := gcdb.CollPlayers.Coll()
	// _, err = coll.UpdateByID(context.TODO(), mp.Pid, db.D("$inc", db.D("Balance", gold)))
	update := db.D("$inc", db.D("Balance", gold))
	opts := options.FindOneAndUpdate().SetProjection(db.D("Balance", 1)).SetReturnDocument(options.After)
	var doc DocPlayer
	err = coll.FindOneAndUpdate(context.TODO(), db.ID(mp.Pid), update, opts).Decode(&doc)
	if err != nil {
		return
	}
	balance = doc.Balance
	return
}

func (mp *MemPlr) TransferOut(gold int64) (balance int64, err error) {
	lo.Must0(gold >= 0)

	if gold == 0 {
		return mp.Balance()
	}

	coll := gcdb.CollPlayers.Coll()

	// _, err = coll.UpdateByID(context.TODO(), mp.Pid, db.D("$inc", db.D("Balance", gold)))
	filter := db.D(
		"_id", mp.Pid,
		"Balance", db.D("$gte", gold),
	)
	update := db.D("$inc", db.D("Balance", -gold))
	opts := options.FindOneAndUpdate().SetProjection(db.D("Balance", 1)).SetReturnDocument(options.After)
	var doc DocPlayer
	err = coll.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = errors.New("Insufficient wallet balance")
		}
		return
	}
	balance = doc.Balance
	return
}

func (mp *MemPlr) TransferOutAll() (amount float64, err error) {
	coll := gcdb.CollPlayers.Coll()
	// _, err = coll.UpdateByID(context.TODO(), mp.Pid, db.D("$inc", db.D("Balance", ut.Money2Gold(amount))))
	filter := db.ID(mp.Pid)
	update := db.D("$set", db.D("Balance", int64(0)))
	opts := options.FindOneAndUpdate().SetProjection(db.D("Balance", 1))
	var doc DocPlayer
	err = coll.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&doc)
	if err != nil {
		return
	}
	amount = ut.Gold2Money(doc.Balance)
	return
}

func (mp *MemPlr) Balance() (balance int64, err error) {
	coll := gcdb.CollPlayers.Coll()
	// _, err = coll.UpdateByID(context.TODO(), mp.Pid, db.D("$inc", db.D("Balance", ut.Money2Gold(amount))))
	filter := db.ID(mp.Pid)
	opts := options.FindOne().SetProjection(db.D("Balance", 1))
	var doc DocPlayer
	err = coll.FindOne(context.TODO(), filter, opts).Decode(&doc)
	if err != nil {
		return
	}
	// balance = ut.Gold2Money(doc.Balance)
	balance = doc.Balance
	return
}
