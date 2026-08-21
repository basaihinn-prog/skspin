package models

import (
	"context"
	"log/slog"
	"serve/comm/db"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// 删除用户的上一局， todo 应该把betlog 相关的也删掉
func (plr *Player) RewriteLastData() {
	_, err := db.Collection("players").UpdateOne(context.TODO(), bson.M{"_id": plr.PID}, bson.M{"$set": bson.M{"lastsid": "", "ls": ""}})
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

type Player struct {
	db.DocPlayer `bson:"inline"`

	LS     string
	LastID string
}

func (plr *Player) IsEndO() (isEnd bool, params []string) {
	if plr.LastID != "" {
		params = strings.Split(plr.LastID, "_")
		if params[1] == params[2] {
			isEnd = true
		}
	} else {
		isEnd = true
	}
	return
}
