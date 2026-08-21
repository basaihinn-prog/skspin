package main

import (
	"context"
	"errors"
	"game/comm"
	"game/duck/lang"
	"game/duck/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	// 每日游戏详情
	RegMsgProc("/AdminInfo/GetGameData", "获取游戏数据", "AdminInfo", getGameDataV2, getGameDataParams{
		Date: "",
		//EndDate:   "",
		Operator:  1,
		PageIndex: 0,
		PageSize:  0,
		NeedLast:  false,
		//Type:      "",
	})
}

type getGameDataParams struct {
	Date         string
	EndDate      string
	Operator     int64
	PageIndex    int64
	PageSize     int64
	NeedLast     bool
	CurrencyType string `json:"CurrencyCode"`
	//Type      string
}

type bDailyDoc struct {
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

	Date         string `bson:"date"`
	Game         string `bson:"game"`
	AppID        string `bson:"appid"`
	CurrencyCode string `json:"CurrencyType" bson:"CurrencyKey"`  //货币名称
	CurrencyName string `json:"CurrencyName" bson:"CurrencyName"` //货币名称
}

type getGameDataResults struct {
	List          []*bDailyDoc
	YesterdayList []*bDailyDoc
	All           int64
}

func getGameDataV2(ctx *Context, ps getGameDataParams, ret *getGameDataResults) (err error) {

	if ps.Date == "" {
		ps.Date = time.Now().Format("20060102")
		ps.NeedLast = false
	}

	user, err := GetUser(ctx)
	if err != nil {
		return
	}

	appidlist := []string{}
	if user.GroupId != 3 {
		appidlist, err = GetOperatopAppID(ctx)
	} else {
		appidlist = append(appidlist, user.AppID)
	}

	if ps.CurrencyType != "" {
		appidlist = appidlist[:0]
		oper := []comm.Operator_V2{}
		err = CollAdminOperator.FindAll(bson.M{"CurrencyKey": ps.CurrencyType}, &oper)
		if err != nil {
			return
		}
		for _, lv := range oper {
			appidlist = append(appidlist, lv.AppID)
		}
	}

	CurrencyCode := ""
	if ps.Operator > 0 {
		var operator *comm.Operator
		err = CollAdminOperator.FindId(ps.Operator, &operator)
		if err != nil {
			return err
		}
		temp := false
		for _, appid := range appidlist {
			if operator.AppID == appid {
				temp = true
			}
		}
		if temp == false {
			return errors.New(lang.GetLang(ctx.Lang, "该币种下没有这个商户或未代理该商户"))
		}
		appidlist = appidlist[:0]
		appidlist = append(appidlist, operator.AppID)
		CurrencyCode = operator.CurrencyKey
		if ps.CurrencyType != "" {
			if ps.CurrencyType != operator.CurrencyKey {
				return errors.New(lang.GetLang(ctx.Lang, "该商户未使用这个币种"))
			}
		}
	} else {
		CurrencyCode = ps.CurrencyType
	}

	match := bson.M{
		"$match": bson.M{
			"date": ps.Date,
			"appid": bson.M{
				"$in": appidlist,
			},
		}}

	group := bson.M{
		"$group": bson.M{
			"_id":           "$game",
			"date":          bson.M{"$first": "$date"},
			"game":          bson.M{"$first": "$game"},
			"betamount":     bson.M{"$sum": "$betamount"},
			"winamount":     bson.M{"$sum": "$winamount"},
			"spincount":     bson.M{"$sum": "$spincount"},
			"spinplrcount":  bson.M{"$sum": "$spinplrcount"},
			"enterplrcount": bson.M{"$sum": "$enterplrcount"},
		}}

	sort := bson.M{
		"$sort": bson.M{
			"game": 1,
		}}

	aggregate := []bson.M{}
	aggregate = append(aggregate, match, group, sort)

	cursor, err := NewOtherDB("reports").Collection("BetDailyReport").Coll().Aggregate(context.TODO(), aggregate)
	if err != nil {
		return err
	}
	err = cursor.All(context.TODO(), &ret.List)
	if err != nil {
		return err
	}

	//币种
	currency := comm.CurrencyType{}
	if ps.Operator != 0 || ps.CurrencyType != "" {
		err = NewOtherDB("GameAdmin").Collection("CurrencyType").FindOne(bson.M{"CurrencyCode": CurrencyCode}, &currency)
		if err != nil {
			return err
		}
	}
	for i, _ := range ret.List {
		ret.List[i].CurrencyCode = currency.CurrencyCode
		ret.List[i].CurrencyName = currency.CurrencyName
		if ctx.Lang != "zh" {
			ret.List[i].CurrencyName = currency.CurrencyCode
		}
	}

	ret.All = int64(len(ret.List))
	if ps.NeedLast {
		t, err := time.Parse("20060102", ps.Date)
		if err != nil {
			return err
		}
		yesterdayT := t.AddDate(0, 0, -1)
		yesterdayDate := yesterdayT.Format("20060102")
		match = bson.M{
			"$match": bson.M{
				"date": yesterdayDate,
				"appid": bson.M{
					"$in": appidlist,
				},
			}}

		group = bson.M{
			"$group": bson.M{
				"_id":           "$game",
				"date":          bson.M{"$first": "$date"},
				"game":          bson.M{"$first": "$game"},
				"betamount":     bson.M{"$sum": "$betamount"},
				"winamount":     bson.M{"$sum": "$winamount"},
				"spincount":     bson.M{"$sum": "$spincount"},
				"spinplrcount":  bson.M{"$sum": "$spinplrcount"},
				"enterplrcount": bson.M{"$sum": "$enterplrcount"},
			}}

		sort = bson.M{
			"$sort": bson.M{
				"game": 1,
			}}

		aggregate = []bson.M{}
		aggregate = append(aggregate, match, group, sort)

		cursor, err = NewOtherDB("reports").Collection("BetDailyReport").Coll().Aggregate(context.TODO(), aggregate)
		if err != nil {
			return err
		}
		err = cursor.All(context.TODO(), &ret.YesterdayList)
		if err != nil {
			return err
		}

		for i, _ := range ret.YesterdayList {
			ret.YesterdayList[i].CurrencyCode = currency.CurrencyCode
			ret.YesterdayList[i].CurrencyName = currency.CurrencyName
			if ctx.Lang != "zh" {
				ret.YesterdayList[i].CurrencyName = currency.CurrencyCode
			}
		}
	}

	return
}

func getGameData(ctx *Context, ps getGameDataParams, ret *getGameDataResults) (err error) {
	if ps.Date == "" {
		ps.Date = time.Now().Format("20060102")
		ps.NeedLast = false
	}

	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	if user.AppID == "admin" && ps.Operator == 0 { //获取所有运营商的
		match := bson.M{
			"$match": bson.M{
				"date": ps.Date,
			}}

		group := bson.M{
			"$group": bson.M{
				"_id":           "$game",
				"date":          bson.M{"$first": "$date"},
				"game":          bson.M{"$first": "$game"},
				"betamount":     bson.M{"$sum": "$betamount"},
				"winamount":     bson.M{"$sum": "$winamount"},
				"spincount":     bson.M{"$sum": "$spincount"},
				"spinplrcount":  bson.M{"$sum": "$spinplrcount"},
				"enterplrcount": bson.M{"$sum": "$enterplrcount"},
			}}

		sort := bson.M{
			"$sort": bson.M{
				"game": 1,
			}}

		aggregate := []bson.M{}
		aggregate = append(aggregate, match, group, sort)

		cursor, err := NewOtherDB("reports").Collection("BetDailyReport").Coll().Aggregate(context.TODO(), aggregate)
		if err != nil {
			return err
		}
		err = cursor.All(context.TODO(), &ret.List)
		if err != nil {
			return err
		}
		ret.All = int64(len(ret.List))
		if ps.NeedLast {
			t, err := time.Parse("20060102", ps.Date)
			if err != nil {
				return err
			}
			yesterdayT := t.AddDate(0, 0, -1)
			yesterdayDate := yesterdayT.Format("20060102")
			match = bson.M{
				"$match": bson.M{
					"date": yesterdayDate,
				}}

			group = bson.M{
				"$group": bson.M{
					"_id":           "$game",
					"date":          bson.M{"$first": "$date"},
					"game":          bson.M{"$first": "$game"},
					"betamount":     bson.M{"$sum": "$betamount"},
					"winamount":     bson.M{"$sum": "$winamount"},
					"spincount":     bson.M{"$sum": "$spincount"},
					"spinplrcount":  bson.M{"$sum": "$spinplrcount"},
					"enterplrcount": bson.M{"$sum": "$enterplrcount"},
				}}

			sort = bson.M{
				"$sort": bson.M{
					"game": 1,
				}}

			aggregate = []bson.M{}
			aggregate = append(aggregate, match, group, sort)

			cursor, err = NewOtherDB("reports").Collection("BetDailyReport").Coll().Aggregate(context.TODO(), aggregate)
			if err != nil {
				return err
			}
			err = cursor.All(context.TODO(), &ret.YesterdayList)
			if err != nil {
				return err
			}
		}
	} else {
		filter := mongodb.FindPageOpt{
			Page:     ps.PageIndex,
			PageSize: 1000,
			Sort:     bson.M{"_id": 1},
		}
		query := bson.M{
			"date": ps.Date,
		}
		if user.AppID == "admin" {
			var operator *comm.Operator_V2
			err = CollAdminOperator.FindId(ps.Operator, &operator)
			if err != nil {
				return err
			}
			query["appid"] = operator.AppID
		} else {
			query["appid"] = user.AppID
		}
		filter.Query = query
		count, err := NewOtherDB("reports").Collection("BetDailyReport").FindPage(filter, &ret.List)
		if err != nil {
			return err
		}
		ret.All = count
		if ps.NeedLast {
			t, err := time.Parse("20060102", ps.Date)
			if err != nil {
				return err
			}
			yesterdayT := t.AddDate(0, 0, -1)
			yesterdayDate := yesterdayT.Format("20060102")
			query["date"] = yesterdayDate
			filter.Query = query
			_, err = NewOtherDB("reports").Collection("BetDailyReport").FindPage(filter, &ret.YesterdayList)
			if err != nil {
				return err
			}
		}
	}

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
//		cursor, err := NewOtherDB("reports").Collection("BetDailyReport").Coll().Aggregate(context.TODO(), []bson.M{
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
