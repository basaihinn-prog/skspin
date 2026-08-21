package main

import (
	"game/comm"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
	"time"
)

func init() {

	RegMsgProc("/AdminInfo/GetPlayerRestrictionsList", "获取玩家赢钱约束列表", "AdminInfo", getPlayerRestrictionsList, getPlayerRestrictionsParams{
		Uid:   "",
		AppID: "",
		Pid:   0,
	})
	RegMsgProc("/AdminInfo/SetPlayerRestrictions", "设置玩家约束", "AdminInfo", setPlayerRestrictions, updateRestrictions{
		Pids: nil,
	})
	RegMsgProc("/AdminInfo/CancelPlayerRestrictions", "关闭玩家约束", "AdminInfo", cancelPlayerRestrictions, updateRestrictions{

		Pid: 0,
	})
}

type getPlayerRestrictionsParams struct {
	Pid      int64  `json:"Pid" bson:"_id"`
	Uid      string `json:"Uid" bson:"Uid"`
	AppID    string `json:"AppID" bson:"AppID"`
	PageSize int64  `json:"PageSize" bson:"-"`
	Page     int64  `json:"Page" bson:"-"`
}

type getRestrictionsListResult struct {
	AllCount int64
	List     []*comm.Player
}

type setRestrictions struct {
	NoExitPId string
}

type updateRestrictions struct {
	Pids                 []int64 `json:"Pids" bson:"-"`
	Pid                  int64   `json:"Pid" bson:"_id"`
	RestrictionsStatus   int64   `json:"RestrictionsStatus" bson:"RestrictionsStatus" description:"玩家约束状态"`
	RestrictionsMaxWin   int64   `json:"RestrictionsMaxWin" bson:"RestrictionsMaxWin" description:"玩家约束最大赢取金额"`
	RestrictionsMaxMulti float64 `json:"RestrictionsMaxMulti" bson:"RestrictionsMaxMulti" description:"玩家约束最大赢取倍数"`
}

func getPlayerRestrictionsList(ctx *Context, ps getPlayerRestrictionsParams, ret *getRestrictionsListResult) (err error) {

	user, err := GetUser(ctx)
	if err != nil {
		return err
	}

	playerColl := NewOtherDB("game").Collection("Players")
	findBson := bson.M{}

	if ps.Pid > 0 {
		findBson["_id"] = ps.Pid
	}

	if ps.Uid != "" {
		findBson["Uid"] = ps.Uid
	}

	if user.GroupId > 1 {
		findBson["AppID"] = user.AppID
	}

	if ps.AppID != "" {
		findBson["AppID"] = ps.AppID
	}

	findBson["RestrictionsStatus"] = 1

	filter := mongodb.FindPageOpt{
		Page:     ps.Page,
		PageSize: ps.PageSize,
		Sort:     bson.M{"RestrictionsTime": -1},
		Query:    findBson,
	}
	Count, err := playerColl.FindPage(filter, &ret.List)
	if err != nil {
		return err
	}
	ret.AllCount = Count

	return nil
}
func setPlayerRestrictions(ctx *Context, ps updateRestrictions, ret *setRestrictions) (err error) {

	filterBson := bson.M{}
	user, b := IsAdminUser(ctx)
	if b == false {
		filterBson["AppID"] = user.AppID
	}

	var players []comm.Player
	filterBson["_id"] = bson.M{
		"$in": ps.Pids,
	}
	playerColl := NewOtherDB("game").Collection("Players")
	// 获取用户下的所有玩家
	err = playerColl.FindAll(filterBson, &players)
	if err != nil {
		return err
	}
	noExitPId, existPlayerArr := existPlayer(players, ps.Pids)

	if len(existPlayerArr) > 0 {

		err = playerColl.Update(bson.M{"_id": bson.M{"$in": existPlayerArr}}, bson.M{
			"$set": bson.M{
				"RestrictionsStatus":   1,
				"RestrictionsMaxMulti": ps.RestrictionsMaxMulti,
				"RestrictionsMaxWin":   ps.RestrictionsMaxWin,
				"RestrictionsTime":     time.Now().UTC(),
			},
		})

		if err != nil {
			return err
		}
	}

	ret.NoExitPId = noExitPId
	return err
}
func cancelPlayerRestrictions(ctx *Context, ps updateRestrictions, ret *getRestrictionsListResult) (err error) {
	playerColl := NewOtherDB("game").Collection("Players")

	err = playerColl.UpdateOne(bson.M{"_id": ps.Pid}, bson.M{
		"$set": bson.M{
			"RestrictionsStatus":   0,
			"RestrictionsMaxMulti": 0,
			"RestrictionsMaxWin":   0,
		},
	})
	if err != nil {
		return err
	}
	return err
}

func existPlayer(players []comm.Player, pids []int64) (string, []int64) {
	var exitPidArr []int64
	var pidArr []string
	if len(players) != len(pids) {
		for _, i := range players {
			pidArr = append(pidArr, strconv.FormatInt(i.Id, 10))
		}

		pidArrStr := strings.Join(pidArr, ",") + ","

		pidArr = []string{}
		for _, i := range pids {

			if strings.Index(pidArrStr, strconv.FormatInt(i, 10)+",") < 0 {

				pidArr = append(pidArr, strconv.FormatInt(i, 10))

			} else {
				exitPidArr = append(exitPidArr, i)
			}
		}

		return strings.Join(pidArr, ","), exitPidArr
	} else {

		return "", pids
	}

}
