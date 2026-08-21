package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/slotsmongo"
	"game/duck/lang"
	"sort"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/GetGameEarningsData", "获取游戏运行时游戏收益", "AdminInfo", getGameEarningsData, getGameEarningsDataParams{
		Game:       "",
		OperatorId: 0,
	})
}

type getGameEarningsDataParams struct {
	Game       string `json:"GameID"`
	OperatorId int
	AppID      string `json:"AppID"`
	Type       int    //专门用来区分彩票 0是1期 1是2期
}

type GameDataInfo struct {
	Id             int64 `json:"-" bson:"_id"`
	WinAmount      int64 `bson:"winamount"`      // 产出
	BetAmount      int64 `bson:"betamount"`      // 流水
	BuyWinAmount   int64 `bson:"buywinamount"`   // 产出
	BuyBetAmount   int64 `bson:"buybetamount"`   // 流水
	ExtraBetAmount int64 `bson:"extrabetamount"` // 流水
	ExtraWinAmount int64 `bson:"extrawinamount"` // 流水
}

type WeekData struct {
	Time             string
	DauValue         int64 // enterplrcount
	BetValue         int64 // bet   betcount
	WinValue         int64 // win   betcount
	SystemValue      int64 // bet - win
	SpinCount        int
	BigReward        int64
	BuyBet           int64 // buybet
	BuyWin           int64 // buywin
	BuySystemValue   int64 // buybet - buywin
	ExtraBet         int64
	ExtraWin         int64
	ExtraSystemValue int64   // extrabet - extrawin
	RateValue        float64 // win / bet
	BuyRateValue     float64 // buywin/ butbet
	ExtraRateValue   float64 // extrawin/ extrabet
}

type WeekDataList struct {
	LastWeek []WeekData
	NowWeek  []WeekData
}

type getGameEarningsDataResults struct {
	SevenData GameDataInfo
	MonthData GameDataInfo
	TotalData GameDataInfo
	WeekData  WeekDataList
}

func getGameEarningsData(ctx *Context, ps getGameEarningsDataParams, ret *getGameEarningsDataResults) (err error) {
	if ps.Game == "" {
		return errors.New(lang.GetLang(ctx.Lang, "游戏不存在"))
	}
	//var appid []string
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}

	appidlist := []string{}
	if user.GroupId != 3 {
		appidlist, err = GetOperatopAppID(ctx)
	} else {
		appidlist = append(appidlist, user.AppID)
	}
	if len(ps.AppID) > 0 && user.GroupId != 3 {
		var operator *comm.Operator
		err = CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, &operator)
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
			return errors.New("未代理该商户")
		}
		appidlist = appidlist[:0]
		appidlist = append(appidlist, operator.AppID)
	}

	if err != nil {
		return
	}
	//if ps.Game != "lottery" {
	now := time.Now()

	// 7日数据
	nowstr := now.Format("20060102")
	day7str := now.AddDate(0, 0, -6).Format("20060102")
	ret.SevenData, err = getTimeData(day7str, nowstr, ps.Game, &appidlist)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	}
	// 本月数据
	monthFirst := fmt.Sprintf("%s00", now.Format("200601"))
	monthEnd := fmt.Sprintf("%s31", now.Format("200601"))
	ret.MonthData, err = getTimeData(monthFirst, monthEnd, ps.Game, &appidlist)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	}

	// 累计数据
	ret.TotalData, err = getTotalData(ps.Game, &appidlist)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	}

	// 双周数据
	ret.WeekData, err = getWeekData(ps.Game, &appidlist)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	}

	return nil
	//}
	//else {
	//	if ps.Type == 0 {
	//		now := time.Now()
	//
	//		// 7日数据
	//		nowstr := now.Format("2006-01-02")
	//		day7str := now.AddDate(0, 0, -6).Format("2006-01-02")
	//		ret.SevenData, err = getLottery1TimeData(day7str, nowstr, appid)
	//		if err != nil {
	//			return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	//		}
	//		// 本月数据
	//		monthFirst := fmt.Sprintf("%s-00", now.Format("2006-01"))
	//		monthEnd := fmt.Sprintf("%s-31", now.Format("2006-01"))
	//		ret.MonthData, err = getLottery1TimeData(monthFirst, monthEnd, appid)
	//		if err != nil {
	//			return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	//		}
	//
	//		// 累计数据
	//		ret.TotalData, err = getLottery1TotalData(appid)
	//		if err != nil {
	//			return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	//		}
	//
	//		// 双周数据
	//		ret.WeekData, err = getLottery1WeekData(appid)
	//		if err != nil {
	//			return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	//		}
	//		return nil
	//	} else if ps.Type == 1 {
	//		now := time.Now()
	//
	//		// 7日数据
	//		nowstr := now.Format("2006-01-02")
	//		day7str := now.AddDate(0, 0, -6).Format("2006-01-02")
	//		ret.SevenData, err = getLottery2TimeData(day7str, nowstr, appid)
	//		if err != nil {
	//			return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	//		}
	//		// 本月数据
	//		monthFirst := fmt.Sprintf("%s-00", now.Format("2006-01"))
	//		monthEnd := fmt.Sprintf("%s-31", now.Format("2006-01"))
	//		ret.MonthData, err = getLottery2TimeData(monthFirst, monthEnd, appid)
	//		if err != nil {
	//			return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	//		}
	//
	//		// 累计数据
	//		ret.TotalData, err = getLottery2TotalData(appid)
	//		if err != nil {
	//			return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	//		}
	//
	//		// 双周数据
	//		ret.WeekData, err = getLottery2WeekData(appid)
	//		if err != nil {
	//			return errors.New(lang.GetLang(ctx.Lang, "查询错误"))
	//		}
	//		return nil
	//	}
	//}

}

func getTimeData(startTime, endTime string, game string, appId *[]string) (ret GameDataInfo, err error) {
	match := bson.M{
		"game": game,
		"date": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}
	//if appId != "" {
	//	match["appid"] = appId
	//}
	match["appid"] = bson.M{"$in": appId}
	cursor, err := NewOtherDB("reports").Collection("BetDailyReport").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":            1,
				"winamount":      bson.M{"$sum": "$winamount"},
				"betamount":      bson.M{"$sum": "$betamount"},
				"buybetamount":   bson.M{"$sum": "$buybetamount"},
				"buywinamount":   bson.M{"$sum": "$buywinamount"},
				"extrabetamount": bson.M{"$sum": "$extrabetamount"},
				"extrawinamount": bson.M{"$sum": "$extrawinamount"},
			},
		},
	})
	if err != nil {
		return
	}
	var datas []GameDataInfo
	err = cursor.All(context.TODO(), &datas)
	if err != nil {
		return
	}
	if len(datas) > 0 {
		ret = datas[0]
	}
	return
}

func getTotalData(game string, appId *[]string) (ret GameDataInfo, err error) {
	match := bson.M{
		"game": game,
	}
	//if appId != "" {
	//	match["appid"] = appId
	//}
	match["appid"] = bson.M{"$in": appId}
	cursor, err := NewOtherDB("reports").Collection("BetSummaryReport").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":            1,
				"winamount":      bson.M{"$sum": "$winamount"},
				"betamount":      bson.M{"$sum": "$betamount"},
				"buybetamount":   bson.M{"$sum": "$buybetamount"},
				"buywinamount":   bson.M{"$sum": "$buywinamount"},
				"extrabetamount": bson.M{"$sum": "$extrabetamount"},
				"extrawinamount": bson.M{"$sum": "$extrawinamount"},
			},
		},
	})
	if err != nil {
		return
	}
	var datas []GameDataInfo
	err = cursor.All(context.TODO(), &datas)
	if err != nil {
		return
	}
	if len(datas) > 0 {
		ret = datas[0]
	}
	return
}

func getWeekData(game string, appId *[]string) (ret WeekDataList, err error) {
	now := time.Now()
	nowstr := now.Format("20060102")
	nowWeekZero := comm.GetMondayOfWeek().Format("20060102")
	lastWeekZero := comm.GetMondayOfWeek().AddDate(0, 0, -7).Format("20060102")
	lastWeekEnd := comm.GetMondayOfWeek().Add(time.Duration(-1)).Format("20060102")

	var datas []*slotsmongo.BetDailyDoc

	aggregate := []bson.M{}

	match := bson.M{
		"$match": bson.M{
			"game":  game,
			"appid": bson.M{"$in": appId},
			"date": bson.M{
				"$gte": lastWeekZero,
				"$lte": nowstr,
			},
		}}

	group := bson.M{
		"$group": bson.M{
			"_id":            "$date",
			"date":           bson.M{"$first": "$date"},
			"betamount":      bson.M{"$sum": "$betamount"},
			"winamount":      bson.M{"$sum": "$winamount"},
			"buybetamount":   bson.M{"$sum": "$buybetamount"},
			"buywinamount":   bson.M{"$sum": "$buywinamount"},
			"extrabetamount": bson.M{"$sum": "$extrabetamount"},
			"extrawinamount": bson.M{"$sum": "$extrawinamount"},
			"bigreward":      bson.M{"$sum": "$bigreward"},
			"spincount":      bson.M{"$sum": "$spincount"},
			"spinplrcount":   bson.M{"$sum": "$spinplrcount"},
			"enterplrcount":  bson.M{"$sum": "$enterplrcount"},
		}}

	aggregate = append(aggregate, match, group)

	cursor, err := NewOtherDB("reports").Collection("BetDailyReport").Coll().Aggregate(context.TODO(), aggregate)

	if err != nil {
		return ret, err
	}

	err = cursor.All(context.TODO(), &datas)

	if err != nil {
		return ret, err
	}

	for i := range datas {
		weekData := WeekData{
			Time:             datas[i].Date,
			DauValue:         int64(datas[i].EnterPlrCount),
			BetValue:         datas[i].BetAmount,
			WinValue:         datas[i].WinAmount,
			SystemValue:      datas[i].BetAmount - datas[i].WinAmount,
			SpinCount:        datas[i].SpinCount,
			BigReward:        datas[i].BigReward,
			BuyBet:           datas[i].BuyBetAmount,
			BuyWin:           datas[i].BuyWinAmount,
			BuySystemValue:   datas[i].BuyBetAmount - datas[i].BuyWinAmount,
			ExtraBet:         datas[i].ExtraBetAmount,
			ExtraWin:         datas[i].ExtraWinAmount,
			ExtraSystemValue: datas[i].ExtraBetAmount - datas[i].ExtraWinAmount,
			RateValue:        comm.DIV(float64(datas[i].WinAmount), float64(datas[i].BetAmount)),
			BuyRateValue:     comm.DIV(float64(datas[i].BuyWinAmount), float64(datas[i].BuyBetAmount)),
			ExtraRateValue:   comm.DIV(float64(datas[i].ExtraWinAmount), float64(datas[i].ExtraBetAmount)),
		}
		if datas[i].Date >= nowWeekZero && datas[i].Date <= nowstr { // 本周
			ret.NowWeek = append(ret.NowWeek, weekData)

		} else if datas[i].Date >= lastWeekZero && datas[i].Date <= lastWeekEnd { // 上周
			ret.LastWeek = append(ret.LastWeek, weekData)
		}
	}

	nowMondayT := comm.GetMondayOfWeek()
	lastMondayT := comm.GetMondayOfWeek().AddDate(0, 0, -7)

	for i := 0; i < 7; i++ {
		date := lastMondayT.AddDate(0, 0, i).Format("20060102")
		if _, ok := lo.Find(ret.LastWeek, func(item WeekData) bool {
			return item.Time == date
		}); !ok {
			ret.LastWeek = append(ret.LastWeek, WeekData{
				Time: date,
			})
		}
	}
	for i := 0; i < 7; i++ {
		date := nowMondayT.AddDate(0, 0, i).Format("20060102")
		if date > nowstr {
			break
		}
		if _, ok := lo.Find(ret.NowWeek, func(item WeekData) bool {
			return item.Time == date
		}); !ok {
			ret.NowWeek = append(ret.NowWeek, WeekData{
				Time: date,
			})
		}
	}

	sort.Slice(ret.NowWeek, func(i, j int) bool {
		return ret.NowWeek[i].Time > ret.NowWeek[j].Time
	})

	sort.Slice(ret.LastWeek, func(i, j int) bool {
		return ret.LastWeek[i].Time > ret.LastWeek[j].Time
	})

	return ret, err
}

// ///////////////////彩票新加
// 1期买彩票
func getLottery1TimeData(startTime, endTime string, appId string) (ret GameDataInfo, err error) {
	match := bson.M{
		"buyDay": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}
	if appId != "" {
		match["appID"] = appId
	}
	cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":       1,
				"winamount": bson.M{"$sum": "$periodProfits"},
				"betamount": bson.M{"$sum": "$soldGold"},
			},
		},
	})
	if err != nil {
		return
	}
	var datas []GameDataInfo
	err = cursor.All(context.TODO(), &datas)
	if err != nil {
		return
	}
	if len(datas) > 0 {
		ret = datas[0]
	}
	return
}

// 2期猜号码
func getLottery2TimeData(startTime, endTime string, appId string) (ret GameDataInfo, err error) {

	//fmt.Println("getLottery2TimeData invoke startTime:", startTime, " endTime:", endTime)

	match := bson.M{
		"buyDay": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}
	if appId != "" {
		match["appID"] = appId
	}
	cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day_Guess").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":       1,
				"winamount": bson.M{"$sum": "$periodGuessProfits"},
				"betamount": bson.M{"$sum": "$soldGuessGold"},
			},
		},
	})
	if err != nil {
		return
	}
	var datas []GameDataInfo
	err = cursor.All(context.TODO(), &datas)
	if err != nil {
		return
	}
	if len(datas) > 0 {
		ret = datas[0]
	}
	return
}

// 1期买彩票
func getLottery1TotalData(appId string) (ret GameDataInfo, err error) {
	match := bson.M{}
	if appId != "" {
		match["appID"] = appId
	}

	cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":       1,
				"winamount": bson.M{"$sum": "$periodProfits"},
				"betamount": bson.M{"$sum": "$soldGold"},
			},
		},
	})
	if err != nil {
		return
	}
	var datas []GameDataInfo
	err = cursor.All(context.TODO(), &datas)
	if err != nil {
		return
	}
	if len(datas) > 0 {
		ret = datas[0]
	}
	return
}

// 2期猜号码
func getLottery2TotalData(appId string) (ret GameDataInfo, err error) {
	match := bson.M{}
	if appId != "" {
		match["appID"] = appId
	}

	cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day_Guess").Coll().Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":       1,
				"winamount": bson.M{"$sum": "$periodGuessProfits"},
				"betamount": bson.M{"$sum": "$soldGuessGold"},
			},
		},
	})
	if err != nil {
		return
	}
	var datas []GameDataInfo
	err = cursor.All(context.TODO(), &datas)
	if err != nil {
		return
	}
	if len(datas) > 0 {
		ret = datas[0]
	}
	return
}

// 1期买彩票
func getLottery1WeekData(appId string) (ret WeekDataList, err error) {
	now := time.Now()
	nowstr := now.Format("20060102")
	nowWeekZero := comm.GetMondayOfWeek().Format("20060102")
	lastWeekZero := comm.GetMondayOfWeek().AddDate(0, 0, -7).Format("20060102")
	lastWeekEnd := comm.GetMondayOfWeek().Add(time.Duration(-1)).Format("20060102")

	var datas []*slotsmongo.BetDailyDoc
	if appId == "" {
		// 需要聚合查询
		cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day").Coll().Aggregate(context.TODO(), []bson.M{
			{
				"$match": bson.M{
					"buyDay": bson.M{
						"$gte": comm.GetMondayOfWeek().AddDate(0, 0, -7).Format("2006-01-02"),
						"$lte": now.Format("2006-01-02"),
					},
				},
			}, {
				"$group": bson.M{
					"_id":       "$buyDay",
					"date":      bson.M{"$first": "$buyDay"},
					"betamount": bson.M{"$sum": "$soldGold"},
					"winamount": bson.M{"$sum": "$periodProfits"},
					// "buybetamount":  0,
					// "buywinamount":  0,
					// "bigreward":     0,
					// "spincount":     0,
					"spinplrcount":  bson.M{"$sum": "$buyPersonCount"},
					"enterplrcount": bson.M{"$sum": "$enterCount"},
				},
			},
		})
		if err != nil {
			return ret, err
		}
		err = cursor.All(context.TODO(), &datas)
		if err != nil {
			return ret, err
		}
	} else {
		cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day").Coll().Aggregate(context.TODO(), []bson.M{
			{
				"$match": bson.M{
					"appID": appId,
					"buyDay": bson.M{
						"$gte": comm.GetMondayOfWeek().AddDate(0, 0, -7).Format("2006-01-02"),
						"$lte": now.Format("2006-01-02"),
					},
				},
			}, {
				"$group": bson.M{
					"_id":       "$buyDay",
					"date":      bson.M{"$first": "$buyDay"},
					"betamount": bson.M{"$sum": "$soldGold"},
					"winamount": bson.M{"$sum": "$periodProfits"},
					// "buybetamount":  0,
					// "buywinamount":  0,
					// "bigreward":     0,
					// "spincount":     0,
					"spinplrcount":  bson.M{"$sum": "$buyPersonCount"},
					"enterplrcount": bson.M{"$sum": "$enterCount"},
				},
			},
		})
		if err != nil {
			return ret, err
		}
		err = cursor.All(context.TODO(), &datas)
		if err != nil {
			return ret, err
		}
	}

	//时间字符串格式不一样,查出来转换一下
	for _, item := range datas {
		item.ID = strings.ReplaceAll(item.ID, "-", "")
		item.Date = strings.ReplaceAll(item.Date, "-", "")
	}

	for i := range datas {
		if datas[i].Date >= nowWeekZero && datas[i].Date <= nowstr { // 本周
			ret.NowWeek = append(ret.NowWeek, WeekData{
				Time:           datas[i].Date,
				DauValue:       int64(datas[i].EnterPlrCount),
				BetValue:       datas[i].BetAmount,
				WinValue:       datas[i].WinAmount,
				SystemValue:    datas[i].BetAmount - datas[i].WinAmount,
				SpinCount:      datas[i].SpinCount,
				BigReward:      datas[i].BigReward,
				BuyBet:         datas[i].BuyBetAmount,
				BuyWin:         datas[i].BuyWinAmount,
				BuySystemValue: datas[i].BuyBetAmount - datas[i].BuyWinAmount,
				RateValue:      comm.DIV(float64(datas[i].WinAmount), float64(datas[i].BetAmount)),
				BuyRateValue:   comm.DIV(float64(datas[i].BuyWinAmount), float64(datas[i].BuyBetAmount)),
			})

		} else if datas[i].Date >= lastWeekZero && datas[i].Date <= lastWeekEnd { // 上周
			ret.LastWeek = append(ret.LastWeek, WeekData{
				Time:           datas[i].Date,
				DauValue:       int64(datas[i].EnterPlrCount),
				BetValue:       datas[i].BetAmount,
				WinValue:       datas[i].WinAmount,
				SystemValue:    datas[i].BetAmount - datas[i].WinAmount,
				SpinCount:      datas[i].SpinCount,
				BigReward:      datas[i].BigReward,
				BuyBet:         datas[i].BuyBetAmount,
				BuyWin:         datas[i].BuyWinAmount,
				BuySystemValue: datas[i].BuyBetAmount - datas[i].BuyWinAmount,
				RateValue:      comm.DIV(float64(datas[i].WinAmount), float64(datas[i].BetAmount)),
				BuyRateValue:   comm.DIV(float64(datas[i].BuyWinAmount), float64(datas[i].BuyBetAmount)),
			})
		}
	}

	nowMondayT := comm.GetMondayOfWeek()
	lastMondayT := comm.GetMondayOfWeek().AddDate(0, 0, -7)

	for i := 0; i < 7; i++ {
		date := lastMondayT.AddDate(0, 0, i).Format("20060102")
		if _, ok := lo.Find(ret.LastWeek, func(item WeekData) bool {
			return item.Time == date
		}); !ok {
			ret.LastWeek = append(ret.LastWeek, WeekData{
				Time: date,
			})
		}
	}
	for i := 0; i < 7; i++ {
		date := nowMondayT.AddDate(0, 0, i).Format("20060102")
		if date > nowstr {
			break
		}
		if _, ok := lo.Find(ret.NowWeek, func(item WeekData) bool {
			return item.Time == date
		}); !ok {
			ret.NowWeek = append(ret.NowWeek, WeekData{
				Time: date,
			})
		}
	}

	sort.Slice(ret.NowWeek, func(i, j int) bool {
		return ret.NowWeek[i].Time > ret.NowWeek[j].Time
	})

	sort.Slice(ret.LastWeek, func(i, j int) bool {
		return ret.LastWeek[i].Time > ret.LastWeek[j].Time
	})

	return ret, err
}

// 2期猜号码
func getLottery2WeekData(appId string) (ret WeekDataList, err error) {
	now := time.Now()
	nowstr := now.Format("20060102")
	nowWeekZero := comm.GetMondayOfWeek().Format("20060102")
	lastWeekZero := comm.GetMondayOfWeek().AddDate(0, 0, -7).Format("20060102")
	lastWeekEnd := comm.GetMondayOfWeek().Add(time.Duration(-1)).Format("20060102")

	var datas []*slotsmongo.BetDailyDoc
	if appId == "" {
		// 需要聚合查询
		cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day_Guess").Coll().Aggregate(context.TODO(), []bson.M{
			{
				"$match": bson.M{
					"buyDay": bson.M{
						"$gte": comm.GetMondayOfWeek().AddDate(0, 0, -7).Format("2006-01-02"),
						"$lte": now.Format("2006-01-02"),
					},
				},
			}, {
				"$group": bson.M{
					"_id":       "$buyDay",
					"date":      bson.M{"$first": "$buyDay"},
					"betamount": bson.M{"$sum": "$soldGuessGold"},
					"winamount": bson.M{"$sum": "$periodGuessProfits"},
					// "buybetamount":  0,
					// "buywinamount":  0,
					// "bigreward":     0,
					// "spincount":     0,
					"spinplrcount":  bson.M{"$sum": "$buyPersonCount"},
					"enterplrcount": bson.M{"$sum": "$enterCount"},
				},
			},
		})
		if err != nil {
			return ret, err
		}
		err = cursor.All(context.TODO(), &datas)
		if err != nil {
			return ret, err
		}
	} else {
		cursor, err := NewOtherDB("CaiPiao").Collection("Report_Day_Guess").Coll().Aggregate(context.TODO(), []bson.M{
			{
				"$match": bson.M{
					"appID": appId,
					"buyDay": bson.M{
						"$gte": comm.GetMondayOfWeek().AddDate(0, 0, -7).Format("2006-01-02"),
						"$lte": now.Format("2006-01-02"),
					},
				},
			}, {
				"$group": bson.M{
					"_id":       "$buyDay",
					"date":      bson.M{"$first": "$buyDay"},
					"betamount": bson.M{"$sum": "$soldGuessGold"},
					"winamount": bson.M{"$sum": "$periodGuessProfits"},
					// "buybetamount":  0,
					// "buywinamount":  0,
					// "bigreward":     0,
					// "spincount":     0,
					"spinplrcount":  bson.M{"$sum": "$buyPersonCount"},
					"enterplrcount": bson.M{"$sum": "$enterCount"},
				},
			},
		})
		if err != nil {
			return ret, err
		}
		err = cursor.All(context.TODO(), &datas)
		if err != nil {
			return ret, err
		}
	}

	//时间字符串格式不一样,查出来转换一下
	for _, item := range datas {
		item.ID = strings.ReplaceAll(item.ID, "-", "")
		item.Date = strings.ReplaceAll(item.Date, "-", "")
	}

	for i := range datas {
		if datas[i].Date >= nowWeekZero && datas[i].Date <= nowstr { // 本周
			ret.NowWeek = append(ret.NowWeek, WeekData{
				Time:           datas[i].Date,
				DauValue:       int64(datas[i].EnterPlrCount),
				BetValue:       datas[i].BetAmount,
				WinValue:       datas[i].WinAmount,
				SystemValue:    datas[i].BetAmount - datas[i].WinAmount,
				SpinCount:      datas[i].SpinCount,
				BigReward:      datas[i].BigReward,
				BuyBet:         datas[i].BuyBetAmount,
				BuyWin:         datas[i].BuyWinAmount,
				BuySystemValue: datas[i].BuyBetAmount - datas[i].BuyWinAmount,
				RateValue:      comm.DIV(float64(datas[i].WinAmount), float64(datas[i].BetAmount)),
				BuyRateValue:   comm.DIV(float64(datas[i].BuyWinAmount), float64(datas[i].BuyBetAmount)),
			})

		} else if datas[i].Date >= lastWeekZero && datas[i].Date <= lastWeekEnd { // 上周
			ret.LastWeek = append(ret.LastWeek, WeekData{
				Time:           datas[i].Date,
				DauValue:       int64(datas[i].EnterPlrCount),
				BetValue:       datas[i].BetAmount,
				WinValue:       datas[i].WinAmount,
				SystemValue:    datas[i].BetAmount - datas[i].WinAmount,
				SpinCount:      datas[i].SpinCount,
				BigReward:      datas[i].BigReward,
				BuyBet:         datas[i].BuyBetAmount,
				BuyWin:         datas[i].BuyWinAmount,
				BuySystemValue: datas[i].BuyBetAmount - datas[i].BuyWinAmount,
				RateValue:      comm.DIV(float64(datas[i].WinAmount), float64(datas[i].BetAmount)),
				BuyRateValue:   comm.DIV(float64(datas[i].BuyWinAmount), float64(datas[i].BuyBetAmount)),
			})
		}
	}

	nowMondayT := comm.GetMondayOfWeek()
	lastMondayT := comm.GetMondayOfWeek().AddDate(0, 0, -7)

	for i := 0; i < 7; i++ {
		date := lastMondayT.AddDate(0, 0, i).Format("20060102")
		if _, ok := lo.Find(ret.LastWeek, func(item WeekData) bool {
			return item.Time == date
		}); !ok {
			ret.LastWeek = append(ret.LastWeek, WeekData{
				Time: date,
			})
		}
	}
	for i := 0; i < 7; i++ {
		date := nowMondayT.AddDate(0, 0, i).Format("20060102")
		if date > nowstr {
			break
		}
		if _, ok := lo.Find(ret.NowWeek, func(item WeekData) bool {
			return item.Time == date
		}); !ok {
			ret.NowWeek = append(ret.NowWeek, WeekData{
				Time: date,
			})
		}
	}

	sort.Slice(ret.NowWeek, func(i, j int) bool {
		return ret.NowWeek[i].Time > ret.NowWeek[j].Time
	})

	sort.Slice(ret.LastWeek, func(i, j int) bool {
		return ret.LastWeek[i].Time > ret.LastWeek[j].Time
	})

	return ret, err
}
