package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/GetPlayRTP", "获取游戏列表(运营商)", "AdminInfo", getGamListExc, &GetPlayRTPList{})
}

type GetPlayRTPList struct {
	DateType   int     //是最后登录时间还是创建时间
	AppID      string  //商户
	EndTime    string  //起始时间
	StartTime  string  //结束时间
	NowBet     float64 //今日下注
	HistoryBet float64 //历史下注
	NowWin     float64 //今日输赢
	HistoryWin float64 //历史输赢
	Page       int64   //
	PageSize   int64
}
type playRTPList struct {
	CreateTime time.Time `bson:"CreateAt"` //创建时间
	LoginTime  time.Time `bson:"LoginAt"`  // string            //最后登录时间
	AppID      string    //所属商户
	Uid        string    `bson:"Uid"`        //用户账号
	Pid        int64     `bson:"_id"`        //玩家唯一标识
	NowBet     *int64    `bson:"betamount"`  //今日下注额
	HistoryBet *int64    `bson:"HistoryBet"` //历史下注
	NowWin     *int64    `bson:"winamount"`  //今日输赢
	HistoryWin *int64    `bson:"HistoryWin"` //历史输赢
	Balance    int64     `bson:"Balance"`    //余额
}

type GetplayRTPListResults struct {
	List  []*playRTPList
	Count int64
}

func getGamListExc(ctx *Context, ps GetPlayRTPList, ret *GetplayRTPListResults) (err error) {

	//var nowPlayer []*playRTPList
	retCache := &GetplayRTPListResults{}
	if ps.AppID == "" {
		return errors.New("商户为必填项")
	}
	ps.NowBet = ps.NowBet * 10000
	ps.NowWin = ps.NowWin * 10000
	ps.HistoryBet = ps.HistoryBet * 10000
	ps.HistoryWin = ps.HistoryWin * 10000
	ok := false
	err = getPlayerByTime(&ps, &retCache.List)
	pidList := []int64{}
	for _, v := range retCache.List {
		pidList = append(pidList, v.Pid)
	}
	if ps.NowWin+ps.NowBet+ps.HistoryBet+ps.HistoryWin != 0 {
		ok = true
	}

	onematch := bson.M{}

	group := bson.M{
		"$group": bson.M{
			"_id":        "$pid",
			"HistoryBet": bson.M{"$sum": "$betamount"},
			"HistoryWin": bson.M{"$sum": "$winamount"},
		}}

	twomatch := bson.M{
		"$match": bson.M{
			"HistoryBet": bson.M{
				"$gte": ps.HistoryBet,
			},
			"HistoryWin": bson.M{
				"$gte": ps.HistoryWin,
			},
		}}

	project := bson.M{
		"$project": bson.M{
			"_id": 0,
		},
	}

	sort := bson.M{
		"$sort": bson.M{
			"Pid": 1,
		}}

	var aggregate []bson.M
	aggregate = append(aggregate, onematch, group, twomatch, project, sort)
	coll := NewOtherDB("reports").Collection("PlayerDailyReport").Coll()

	//历史下注
	now := time.Now().Format("20060102")
	for i, v := range retCache.List {
		if ps.NowBet != 0 || ps.NowWin != 0 {
			onematch = bson.M{
				"$match": bson.M{
					"pid":  v.Pid,
					"date": now,
					"betamount": bson.M{
						"$gte": ps.NowBet,
					},
					"winamount": bson.M{
						"$gte": ps.NowWin,
					},
				}}
		} else {
			onematch = bson.M{
				"$match": bson.M{
					"pid": v.Pid,
					"betamount": bson.M{
						"$gte": ps.NowBet,
					},
					"winamount": bson.M{
						"$gte": ps.NowWin,
					},
				}}
		}

		aggregate[0] = onematch
		count, err := coll.CountDocuments(context.TODO(), bson.M{"pid": v.Pid})
		if count == 0 && err == nil && ok == false {
			retCache.List[i].HistoryWin = new(int64)
			retCache.List[i].HistoryBet = new(int64)
			ret.List = append(ret.List, retCache.List[i])
			ret.Count += 1
			continue
		} else if count == 0 && err == nil {
			continue
		} else if err != nil {
			return err
		}
		cursor, err := coll.Aggregate(context.TODO(), aggregate)
		if err != nil {
			return err
		}
		var sum []struct {
			HistoryBet int64 `bson:"HistoryBet"` //历史下注
			HistoryWin int64 `bson:"HistoryWin"` //历史输赢
		}
		err = cursor.All(context.TODO(), &sum)
		if err != nil {
			return err
		}
		if len(sum) == 0 {
			continue
		}
		retCache.List[i].HistoryWin = &sum[0].HistoryWin
		retCache.List[i].HistoryBet = &sum[0].HistoryBet
		ret.List = append(ret.List, retCache.List[i])
		ret.Count += 1
	}
	ret.List, ret.Count = paginate(ret.List, ps.Page, ps.PageSize)

	//今日下注
	for i, v := range ret.List {
		filter := bson.M{
			"pid":  v.Pid,
			"date": now,
		}
		result := coll.FindOne(context.TODO(), filter, options.FindOne().SetProjection(bson.M{"_id": 0}))
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			continue
		}
		err = result.Decode(ret.List[i])
		if err != nil {
			return
		}
	}

	return
}
func paginate(items []*playRTPList, page, pageSize int64) ([]*playRTPList, int64) {
	totalItems := int64(len(items))
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalItems {
		return []*playRTPList{}, totalItems // 返回空列表和总数
	}

	if end > totalItems {
		end = totalItems
	}

	return items[start:end], totalItems
}

func getPlayerByTime(ps *GetPlayRTPList, list *[]*playRTPList) (err error) {
	query := bson.M{}
	var endTime, startTime time.Time
	if ps.EndTime != "" && ps.StartTime != "" {
		endTime, err = time.Parse("2006-01-02 15:04:05", ps.EndTime)
		startTime, err = time.Parse("2006-01-02 15:04:05", ps.StartTime)
		endTime = endTime.Add(-8 * time.Hour)
		startTime = startTime.Add(-8 * time.Hour)
		switch ps.DateType {
		case 1:
			query["LoginAt"] = bson.M{
				"$lte": endTime,
				"$gte": startTime,
			}
		case 2:
			query["CreateAt"] = bson.M{
				"$lte": endTime,
				"$gte": startTime,
			}
		}
	}
	if ps.AppID != "" {
		query["AppID"] = ps.AppID
	}
	filter := options.FindOptions{
		Sort:       bson.M{"_id": 1},
		Projection: nil,
	}
	err = NewOtherDB("game").Collection("Players").FindAllOpt(query, list, &filter)
	return
}
