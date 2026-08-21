package main

import (
	"context"
	"game/comm"
	"game/comm/db"
	"game/duck/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	RegMsgProc("/AdminInfo/manufacturerList", "获取游戏厂商", "AdminInfo", getGameManufacturerList, GameManufacturerListParams{})
}

type GameManufacturer struct {
	Id               primitive.ObjectID `bson:"_id" json:"Id"`
	ManufacturerName string             `bson:"ManufacturerName" json:"ManufacturerName"`
	ManufacturerCode string             `bson:"ManufacturerCode" json:"ManufacturerCode"`
}

type GameManufacturerListParams struct {
}
type GameManufacturerList struct {
	List []*GameManufacturer
}

func getGameManufacturerList(ctx *Context, ps GameManufacturerListParams, ret *GameManufacturerList) (err error) {
	user, _ := IsAdminUser(ctx)
	query := bson.M{}
	if user.GroupId == 3 {
		var result comm.Operator_V2
		err = CollAdminOperator.FindOne(bson.M{"AppID": user.AppID}, &result)
		if len(result.DefaultManufacturerOn) != 0 {
			query = bson.M{"ManufacturerName": bson.M{"$in": result.DefaultManufacturerOn}}
		}
	}
	manufacturersColl := db.Collection2("game", "GameManufacturers")
	coure, err := manufacturersColl.Find(context.TODO(), query)
	err = coure.All(context.TODO(), &ret.List)
	return
}

func resetGameManufacturer() {
	manufacturersColl := db.Collection2("game", "GameManufacturers")
	GamesColl := db.Collection2("game", "Games")
	cursor, err := GamesColl.Aggregate(context.TODO(), []bson.M{
		{
			"$group": bson.M{
				"_id": "$ManufacturerName",
				"games": bson.M{
					"$addToSet": "$_id",
				},
			},
		},
	})

	if err != nil {
		logger.Err("组合mongoSql获取游标失败：" + err.Error())
		return
	}

	m := make([]map[string]interface{}, 0, 1)
	err = cursor.All(context.TODO(), &m)
	if err != nil {
		logger.Err("联合查询游戏获取厂商失败：" + err.Error())
		return
	}

	var insertModel []interface{}
	for _, item := range m {
		insertModel = append(insertModel, bson.M{
			"ManufacturerName": item["_id"],
			"ManufacturerCode": item["_id"],
		})
	}

	_, err = manufacturersColl.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		logger.Err("删除厂商表数据失败：" + err.Error())
		return
	}
	_, err = manufacturersColl.InsertMany(context.TODO(), insertModel)
	if err != nil {
		logger.Err("同步厂商表数据失败：" + err.Error())
		return
	}

}
