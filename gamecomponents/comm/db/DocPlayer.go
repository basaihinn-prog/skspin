package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocPlayer struct {
	PID      int64  `bson:"_id"`
	Uid      string `bson:"Uid"`   // 外部id
	AppID    string `bson:"AppID"` // 所属产品
	Language string `bson:"-"`     //语言 th,en,cn
}

//func (ps *DocPlayer) Err(err string) error {
//	return errors.New(lang.Get(ps.Language, err))
//	return errors.New(err)
//}

var (
	plrprojection = D(
		"_id", 1,
		"Uid", 1,
		"AppID", 1,
	)
)

func getDocPlayer(pid int64) (plr *DocPlayer, err error) {

	coll := Collection2("game", "Players")

	var doc DocPlayer
	err = coll.FindOne(context.TODO(), ID(pid), options.FindOne().SetProjection(plrprojection)).Decode(&doc)
	if err != nil {
		return
	}
	plr = &doc
	return
}
