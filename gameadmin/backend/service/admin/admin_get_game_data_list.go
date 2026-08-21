package main

import (
	"game/comm"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/GetGameDataList", "获取游戏数据列表", "AdminInfo", getGameDataList, getGameDataListParams{
		Date:      "",
		EndDate:   "",
		Operator:  1,
		PageIndex: 0,
		PageSize:  0,
	})
}

type getGameDataListParams struct {
	Date         string
	EndDate      string
	Game         string
	Operator     int64
	PageIndex    int64
	PageSize     int64
	Manufacturer string `json:"Manufacturer"` //游戏厂商
}

type betDailyDoc struct {
	ID             string `json:"-" bson:"_id"`
	BetAmount      int64  `bson:"betamount"`      // 支出
	WinAmount      int64  `bson:"winamount"`      // 赚取
	BuyBetAmount   int64  `bson:"buybetamount"`   // buy支出
	BuyWinAmount   int64  `bson:"buywinamount"`   // buy赚取
	ExtraBetAmount int64  `bson:"extrabetamount"` // extra支出
	ExtraWinAmount int64  `bson:"extrawinamount"` // extra赚取
	BigReward      int64  `bson:"bigreward"`      // 累计爆大奖
	SpinCount      int    `bson:"spincount"`      // 转动次数
	SpinPlrCount   int    `bson:"spinplrcount"`   // 参与转动的玩家数量
	EnterPlrCount  int    `bson:"enterplrcount"`  // 进入游戏的玩家数量

	Date                 string `bson:"date"`
	Game                 string `bson:"game"`
	GameIcon             string `bson:"-"`
	GameManufacturerName string `bson:"-"`
	AppID                string `bson:"appid"`
	CurrencyKey          string `json:"CurrencyKey" bson:"CurrencyKey"`
	CurrencyName         string `json:"CurrencyName" bson:"CurrencyName"`
}

type getGameDataListResults struct {
	List          []*betDailyDoc
	YesterdayList []*betDailyDoc
	All           int64
}

func getGameDataList(ctx *Context, ps getGameDataListParams, ret *getGameDataListResults) (err error) {

	//ps.PageIndex = 1
	//ps.PageSize = 2000
	if ps.Date == "" {
		now := time.Now()
		ps.EndDate = now.Format("20060102")
		ps.Date = now.AddDate(0, 0, -15).Format("20060102")
	}
	if ps.EndDate == "" {
		ps.EndDate = ps.Date
	}

	query := bson.M{
		"date": bson.M{
			"$gte": ps.Date,
			"$lte": ps.EndDate,
		},
	}
	//商户列表
	appidlist := []string{}
	//oplist := []*comm.Operator{}
	if ps.Operator == 0 {
		appidlist, err = GetOperatopAppID(ctx)
		if err != nil {
			return
		}
	} else {
		appidlist = appidlist[:0]
		var operator *comm.Operator
		err = CollAdminOperator.FindId(ps.Operator, &operator)
		if err != nil {
			return err
		}
		appidlist = append(appidlist, operator.AppID)
	}

	query["appid"] = bson.M{"$in": appidlist}

	//游戏ID列表
	if ps.Game != "ALL" && ps.Game != "" {
		query["game"] = ps.Game
	} else {
		if ps.Manufacturer == "" {
			ps.Manufacturer = "ALL"
		}
		list, err := GetManufacturerGame(ps.Manufacturer)
		if err != nil {
			return err
		}
		gIdlist := []string{}
		for _, item := range list {
			gIdlist = append(gIdlist, item.ID)
		}
		query["game"] = bson.M{"$in": gIdlist}
	}

	filter := mongodb.FindPageOpt{
		Page:     ps.PageIndex,
		PageSize: ps.PageSize,
		Sort:     bson.M{"date": -1},
		Query:    query,
	}

	count, err := NewOtherDB("reports").Collection("BetDailyReport").FindPage(filter, &ret.List)

	ret.All = count

	// 查询游戏数据
	gameNameMap := map[string]string{}
	gameMap := map[string]*comm.Game2{}
	_, games := GetGameListFormatName(ctx.Lang, nil)

	for _, game := range games {

		gameMap[game.ID] = game
		gameNameMap[game.ID] = game.Name
	}

	//币种和游戏名称
	operck, err := GetCurrencyList(ctx)
	for i, vl := range ret.List {

		vl.GameIcon = gameMap[vl.Game].IconUrl
		vl.GameManufacturerName = gameMap[vl.Game].ManufacturerName
		vl.Game = gameNameMap[vl.Game]

		if _, ok := operck[vl.AppID]; !ok {
			operck[vl.AppID] = &comm.CurrencyType{
				CurrencyCode:   "/",
				CurrencyName:   "/",
				CurrencySymbol: "/",
			}
		}

		if ctx.Lang != "zh" {
			ret.List[i].CurrencyName = operck[vl.AppID].CurrencyCode
			ret.List[i].CurrencyKey = operck[vl.AppID].CurrencyCode
		} else {
			ret.List[i].CurrencyName = operck[vl.AppID].CurrencyName
			ret.List[i].CurrencyKey = operck[vl.AppID].CurrencyCode
		}

	}

	if err != nil {
		return err
	}
	ret.All = count
	return nil
}

//func getGameDataMonth(ctx *Context, ps getGameDataParams, ret *getGameDataResults) (err error) {
//	if ps.Date == "" || ps.EndDate == "" {
//		now := time.Now()
//		ps.EndDate = now.Format("20060102")
//		ps.Date = now.AddDate(0, 0, -15).Format("20060102")
//	}
//
//	if ctx.Username == "admin" && ps.Operator == 0 { //获取所有运营商的
//		cursor, err := NewOtherDB("reports").Collection("BetDailyReport").Coll().Aggregate(context.TO/DO(), []bson.M{
//			{
//				"$match": bson.M{
//					"date": bson.M{"$gte": ps.Date, "$lte": ps.EndDate},
//				},
//			},
//			{
//				"$group": bson.M{
//					"_id":           "$game",
//					"game":          bson.M{"$first": "$game"},
//					"betamount":     bson.M{"$sum": "$betamount"},
//					"winamount":     bson.M{"$sum": "$winamount"},
//					"spincount":     bson.M{"$sum": "$spincount"},
//					"spinplrcount":  bson.M{"$sum": "$spinplrcount"},
//					"enterplrcount": bson.M{"$sum": "$enterplrcount"},
//				},
//			},
//		})
//		if err != nil {
//			return err
//		}
//		err = cursor.All(context.TODO(), &ret.List)
//		ret.All = int64(len(ret.List))
//	} else {
//		var operator *comm.Operator
//		err = CollAdminOperator.FindId(ps.Operator, &operator)
//		if err != nil {
//			return err
//		}
//		cursor, err := NewOtherDB("reports").Collection("BetDailyReport").Coll().Aggregate(context.TODO(), []bson.M{
//			{
//				"$match": bson.M{
//					"date":  bson.M{"$gte": ps.Date, "$lte": ps.EndDate},
//					"appid": operator.AppID,
//				},
//			},
//			{
//				"$group": bson.M{
//					"_id":           "$game",
//					"date":          bson.M{"$first": "$date"},
//					"game":          bson.M{"$first": "$game"},
//					"betamount":     bson.M{"$sum": "$betamount"},
//					"winamount":     bson.M{"$sum": "$winamount"},
//					"spincount":     bson.M{"$sum": "$spincount"},
//					"spinplrcount":  bson.M{"$sum": "$spinplrcount"},
//					"enterplrcount": bson.M{"$sum": "$enterplrcount"},
//				},
//			},
//		})
//		if err != nil {
//			return err
//		}
//		err = cursor.All(context.TODO(), &ret.List)
//		ret.All = int64(len(ret.List))
//	}
//	return nil
//}
