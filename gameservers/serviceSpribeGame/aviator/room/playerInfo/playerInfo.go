package playerInfo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"serve/comm/db"
	"serve/comm/lazy"
	"time"
)

type SpribePlayerInfoBody struct {
	Id           primitive.ObjectID `bson:"_id"`
	PID          int64              `bson:"pid"`
	Username     string             `json:"username"`
	AppID        string             `bson:"AppID"` // 所属产品
	ProfileImage string             `json:"profileImage"`
	CreateAt     primitive.DateTime `bson:"createAt"`
	UpdateAt     primitive.DateTime `bson:"updateAt"`
}

func UpsertPlayerInfo(body SpribePlayerInfoBody) error {
	now := primitive.NewDateTimeFromTime(time.Now())
	// 定义过滤条件
	filter := bson.M{"AppID": body.AppID, "pid": body.PID}
	// 定义更新内容
	update := bson.M{
		"$set": bson.M{"username": body.Username, "profileImage": body.ProfileImage, "updateAt": now},
		"$setOnInsert": bson.M{
			"_id":      body.Id,
			"pid":      body.PID,
			"AppID":    body.AppID,
			"createAt": now,
		},
	}
	// 设置 upsert 为 true
	opts := options.Update().SetUpsert(true)
	coll2 := db.Collection2(lazy.ServiceName, "spribePlayerInfo")
	// 执行 upsert 操作
	_, err := coll2.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		slog.Error("UpsertPlayerInfo::coll2.UpdateOne Err", "err", err.Error())
		return err
	}
	return nil
}

func FindPlayerInfo(pid int64, appId string) (SpribePlayerInfoBody, error) {
	var body SpribePlayerInfoBody
	err := db.Collection2(lazy.ServiceName, "spribePlayerInfo").FindOne(context.TODO(),
		db.D("pid", pid, "AppID", appId),
		options.FindOne().SetSort(db.D("createAt", -1))).Decode(&body)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return body, nil
		} else {
			return body, err
		}
	}
	return body, nil
}
