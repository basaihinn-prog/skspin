package main

import (
	"fmt"
	"game/comm"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func init() {
	RegMsgProc("/AdminInfo/GetGameList", "获取游戏列表", "AdminInfo", getGameList, &GetGameRequest{})
	RegMsgProc("/AdminInfo/GetGameListExc", "获取游戏列表(运营商)", "AdminInfo", getGameListExc, &GetGameRequest{})
	RegMsgProc("/AdminInfo/clearUri", "获取游戏列表(运营商)", "AdminInfo", clearUri, &GetGameRequest{})
}

type GetGameRequest struct {
	GameName    string
	GameId      string
	Language    string
	Maintenance string
}

type GetGameListResults struct {
	List  []*comm.Game2
	Count int64
}

func getGameList(ctx *Context, ps GetGameRequest, ret *GetGameListResults) (err error) {

	filter := bson.M{}
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	if ps.GameId != "" {

		filter["_id"] = ps.GameId

	}

	if ps.Maintenance != "" {

		filter["ManufacturerName"] = ps.Maintenance

	}

	if ps.GameName != "" {

		filterStr := fmt.Sprintf("GameNameConfig.%s.GameName", ctx.Lang)

		filter[filterStr] = bson.M{"$regex": ps.GameName}
	}
	err, GameList := GetGameListFormatName(ctx.Lang, filter)
	if err != nil {
		return err
	}

	ret.List = GameList
	return nil
}

func getGameListExc(ctx *Context, ps GetGameRequest, ret *GetGameListResults) (err error) {

	filter := bson.M{}
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	if ps.GameName != "" {

		filterStr := fmt.Sprintf("GameNameConfig.%s.GameName", ctx.Lang)

		filter[filterStr] = bson.M{"$regex": ps.GameName}
	}
	if user.AppID != "admin" {
		filter["Status"] = comm.GameStatus_Open
	}

	if ps.GameId != "" && ps.GameId != "ALL" {

		filter["_id"] = ps.GameId

	}
	Maintenances := []string{}
	if ps.Maintenance != "" && ps.Maintenance != "ALL" {
		Maintenances = append(Maintenances, ps.Maintenance)
		filter["ManufacturerName"] = bson.M{"$in": Maintenances}

	} else {
		if user.GroupId == 3 {
			var result comm.Operator_V2
			err = CollAdminOperator.FindOne(bson.M{"AppID": user.AppID}, &result)
			if len(result.DefaultManufacturerOn) != 0 {
				filter["ManufacturerName"] = bson.M{"$in": result.DefaultManufacturerOn}
			}
		}
	}
	err, GameList := GetGameListFormatName(ctx.Lang, filter)
	if err != nil {
		return err
	}

	ret.List = GameList
	return nil
}

func clearUri(ctx *Context, ps GetGameRequest, ret *GetGameListResults) (err error) {
	var gameList []*comm.Game2
	coll := NewOtherDB("game").Collection("Games_copy1")
	if err = coll.FindAll(bson.M{}, &gameList); err != nil {
		return err
	}

	for _, i2 := range gameList {

		if i2.GameNameConfig["zh"] != nil {
			i2.GameNameConfig["zh"].Icon = strings.Replace(i2.GameNameConfig["zh"].Icon, "https://admin.ppply.site", "", -1)

		}
		if i2.GameNameConfig["idr"] != nil {
			i2.GameNameConfig["idr"].Icon = strings.Replace(i2.GameNameConfig["idr"].Icon, "https://admin.ppply.site", "", -1)

		}
		if i2.GameNameConfig["th"] != nil {
			i2.GameNameConfig["th"].Icon = strings.Replace(i2.GameNameConfig["th"].Icon, "https://admin.ppply.site", "", -1)

		}
		if i2.GameNameConfig["it"] != nil {
			i2.GameNameConfig["it"].Icon = strings.Replace(i2.GameNameConfig["it"].Icon, "https://admin.ppply.site", "", -1)

		}
		if i2.GameNameConfig["es"] != nil {
			i2.GameNameConfig["es"].Icon = strings.Replace(i2.GameNameConfig["es"].Icon, "https://admin.ppply.site", "", -1)

		}
		if i2.GameNameConfig["en"] != nil {
			i2.GameNameConfig["en"].Icon = strings.Replace(i2.GameNameConfig["en"].Icon, "https://admin.ppply.site", "", -1)

		}

		NewOtherDB("game").Collection("Games").UpdateOne(bson.M{"_id": i2.ID}, bson.M{"$set": i2})

	}

	return nil

}

func GetGameListFormatName(Language string, filter bson.M) (err error, gameList []*comm.Game2) {

	coll := NewOtherDB("game").Collection("Games")
	if err = coll.FindAll(filter, &gameList); err != nil {
		return err, nil
	}

	if Language == "" {
		Language = "en"
	}
	for _, game := range gameList {

		if game.GameNameConfig != nil {
			if game.GameNameConfig[Language].GameName != "" {
				game.Name = fmt.Sprintf("【%s】%s", game.GameID, game.GameNameConfig[Language].GameName)
				game.OriginName = game.GameNameConfig[Language].GameName
			}

			if setting != nil && game.GameNameConfig[Language].Icon != "" {
				game.IconUrl = setting.IconAddr + game.GameNameConfig[Language].Icon
			}
		}
	}

	return nil, gameList
}
