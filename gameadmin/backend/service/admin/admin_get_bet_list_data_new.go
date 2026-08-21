package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"game/duck/mongodb"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	//下注历史
	RegMsgProc("/AdminInfo/GetBetLogListNew", "获取下注列表New", "AdminInfo", getBetLogListNew, getBetLogListNewParams{
		PageSize: 10,
	})
	// RegMsgProc("/AdminInfo/DownloadBetLogListNew", "获取下注列表New", "AdminInfo", downloadBetLogListNew, getBetLogListNewParams{
	// 	PageSize: 10,
	// })
}

type getBetLogListNewParams struct {
	OperatorId     int64  //商户ID
	OrderId        string //交易订单号
	Pid            int64  //唯一标识
	UserID         string //玩家ID
	StartTime      int64
	EndTime        int64
	GameID         string //游戏ID
	Bet            *int64 //投注金额
	Win            *int64 //赢分金额
	WinLose        *int64 //玩家输赢金额
	Balance        *int64 //玩家余额
	PageSize       int64
	PageIndex      int64
	Optional       int // 0：下一页， 1：上一页
	LastId         primitive.ObjectID
	Lottery_Peroid int    //默认值为0表示1期买彩票,传1的时候为2期猜号码
	CurrencyCode   string //币种
	Manufacturer   string
	LogType        int
	NextStartTime  int64
	NextOrderId    string
	UpEndTime      int64
	UpOrderId      string
	GameType       int
}

//	type getDocBetLog struct {
//		ID         any       `ch:"id"`     //todo：等数据库清空之后改成string 旧的数据库为primitive.ObjectIDFromHex
//		Pid        int64     `ch:"Pid"`    //内部玩家ID
//		UserID     string    `ch:"UserID"` //外部玩家ID
//		GameID     string    `ch:"GameID"` //游戏ID
//		GameType   int       `ch:"GameType"`
//		RoundID    string    `ch:"RoundID"`    // 对局ID，方便查询
//		Bet        int64     `ch:"Bet"`        // 下注
//		Win        int64     `ch:"Win"`        // 单次消除输赢
//		Comment    string    `ch:"Comment"`    // 注释
//		InsertTime time.Time `ch:"InsertTime"` // 数据插入时间
//
//		Balance int64 `ch:"Balance"` // 余额
//		// RoundCount      int                `bson:"RoundCount"`
//		Completed       bool   `ch:"Completed"`    // 最后一次spin
//		TotalWinLoss    int64  `ch:"TotalWinLoss"` // 当前局结算总输赢, 只需要 Completed:true 才使用
//		WinLose         int64  `ch:"WinLose"`      // 玩家输赢金额
//		Grade           int    `ch:"Grade"`
//		HitBigReward    int64  `ch:"HitBigReward"` // 奖池奖励
//		SpinDetailsJson string `ch:"SpinDetailsJson"`
//		PGBetID         string `ch:"PGBetID"`
//		LogType         int    `ch:"LogType"`        // 0 电子投注/结算(游戏), 1 转入, 2 转出, 3 admin设置, 4 彩票投注, 5 彩票结算, 6 百人投注/结算
//		TransferAmount  int64  `ch:"TransferAmount"` // 转移金额
//
//		AppID        string `ch:"AppID"`                            //
//
// }
type getDocBetLog struct {
	ID               string    `ch:"id"`
	Pid              uint64    `ch:"Pid"`
	GameID           string    `ch:"GameID"`
	TotalWinLoss     int64     `ch:"TotalWinLoss"`
	PGBetID          string    `ch:"PGBetID"`
	UserID           string    `ch:"UserID"`
	RoundID          string    `ch:"RoundID"`
	Win              int64     `ch:"Win"`
	Balance          int64     `ch:"Balance"`
	TransferAmount   int64     `ch:"TransferAmount"`
	Comment          string    `ch:"Comment"`
	Completed        int8      `ch:"Completed"`
	WinLose          int64     `ch:"WinLose"`
	SpinDetailsJson  string    `ch:"SpinDetailsJson"`
	GameType         int32     `ch:"GameType"`
	Bet              int64     `ch:"Bet"`
	InsertTime       time.Time `ch:"InsertTime"`
	AppID            string    `ch:"AppID"`
	Grade            int32     `ch:"Grade"`
	HitBigReward     int64     `ch:"HitBigReward"`
	LogType          int32     `ch:"LogType"`
	ManufacturerName string    `ch:"ManufacturerName"`
	CurrencyKey      string    `ch:"CurrencyKey"`
	UserName         string    `ch:"UserName"`
	CurrencyCode     string    `json:"CurrencyType" bson:"CurrencyKey"`  //货币名称
	CurrencyName     string    `json:"CurrencyName" bson:"CurrencyName"` //货币名称
}

type betLogAviator struct {
	ID                 string    `json:"ID" bson:"_id"` // optional
	Pid                int64     `json:"pid" bson:"pid"`
	AppID              string    `json:"appID" bson:"appID"`             // optional
	Uid                string    `json:"uid" bson:"uid"`                 // optional
	UserName           string    `json:"userName" bson:"userName"`       //用户名
	CurrencyKey        string    `json:"currencyKey" bson:"currencyKey"` //币种
	RoomId             string    `json:"roomId" bson:"roomId"`
	Bet                int64     `json:"bet" bson:"bet"`                            // 下注
	Win                int64     `json:"win" bson:"win"`                            // 输赢
	Balance            int64     `json:"balance" bson:"balance"`                    // 余额
	CashOutDate        int64     `json:"cashOutDate,omitempty" bson:"cashOutDate" ` // 提现时间
	Payout             float64   `json:"payout" bson:"payout" `
	Profit             int64     `json:"profit,omitempty" bson:"profit" ` // 利润
	RoundBetId         string    `json:"roundBetId,omitempty" bson:"roundBetId" `
	RoundId            int64     `json:"roundId,omitempty" bson:"roundId" `
	MaxMultiplier      float64   `json:"maxMultiplier,omitempty" bson:"maxMultiplier" `
	Frb                uint8     `json:"frb" bson:"frb"`               // 是否免费
	InsertTime         time.Time `json:"insertTime" bson:"InsertTime"` // 数据插入时间
	GameID             string    `json:"gameID" bson:"GameID"`
	RoundMaxMultiplier float64   `json:"roundMaxMultiplier,omitempty" bson:"roundMaxMultiplier" `
	LogType            int32     `ch:"LogType"`
	FinishType         int32     `ch:"FinishType"`
	ManufacturerName   string    `ch:"manufacturerName"`
	BetId              int32     `ch:"BetId"`
}

type getBetLogListNewResults struct {
	List  []*getDocBetLog
	Count int64
}

func generateQueryStr(tableName string, ps getBetLogListNewParams, appidlist []string) (queryStr string, countCache int64, ArgList []any, err error) {
	query := bson.M{}
	queryStr = "select * from " + tableName + " where"
	countStr := "select count(*) from " + tableName + " where"
	isFirst := true
	argStr := ""
	argList := []any{}
	if ps.LogType != -1 {
		argStr += " LogType = ? "
		argList = append(argList, ps.LogType)
		isFirst = false
	}
	if ps.Pid > 0 {
		query["Pid"] = ps.Pid
		if isFirst {
			argStr += " Pid = ? "
			isFirst = false
		} else {
			argStr += " and Pid = ? "
		}
		argList = append(argList, ps.Pid)
		if ps.OperatorId == 0 {
			play := &comm.Player{}
			NewOtherDB("game").Collection("Players").FindOne(bson.M{"_id": ps.Pid}, play)
			if play == nil {
				return "", 0, nil, errors.New("未找到玩家")
			}
			appidlist = appidlist[:0]
			appidlist = append(appidlist, play.AppID)
		}
	}
	if ps.Manufacturer != "" {
		if isFirst {
			argStr += " ManufacturerName = ?"
			isFirst = false
		} else {
			argStr += " and ManufacturerName = ?"
		}
		argList = append(argList, ps.Manufacturer)
	}
	if ps.UserID != "" {
		if isFirst {
			argStr += " UserID = ?"
			isFirst = false
		} else {
			argStr += " and UserID = ?"
		}
		argList = append(argList, ps.UserID)
	}

	if ps.GameID != "ALL" && ps.GameID != "" {
		if isFirst {
			argStr += " GameID = ?"
			isFirst = false
		} else {
			argStr += " and GameID = ?"
		}
		argList = append(argList, ps.GameID)
	}

	if ps.Bet != nil {
		if isFirst {
			argStr += " Bet >= ?"
			isFirst = false
		} else {
			argStr += " and Bet >= ?"
		}
		argList = append(argList, *ps.Bet)
	}

	if ps.Win != nil {
		if isFirst {
			argStr += " Win >= ?"
			isFirst = false
		} else {
			argStr += " and Win >= ?"
		}
		argList = append(argList, *ps.Win)
	}

	if ps.OrderId != "" {
		if isFirst {
			argStr += " id = ?"
			isFirst = false
		} else {
			argStr += " and id = ?"
		}
		argList = append(argList, ps.OrderId)
	}

	if ps.WinLose != nil {
		if isFirst {
			argStr += " WinLose >= ?"
			isFirst = false
		} else {
			argStr += " and WinLose >= ?"
		}
		argList = append(argList, *ps.WinLose)
	}

	if ps.Balance != nil {
		if isFirst {
			argStr += " Balance >= ?"
			isFirst = false
		} else {
			argStr += " and Balance >= ?"
		}
		argList = append(argList, *ps.Balance)
	}

	sorts := bson.D{}
	sorts = append(sorts, bson.E{Key: "_id", Value: -1})

	RDate := ps.EndTime + ps.StartTime
	switch RDate {
	case 0:
	case ps.EndTime:
		if isFirst {
			argStr += " InsertTime <= ?"
			isFirst = false
		} else {
			argStr += " and InsertTime <= ?"
		}
		argList = append(argList, ps.EndTime)
	case ps.StartTime:
		if isFirst {
			argStr += " InsertTime >= ?"
			isFirst = false
		} else {
			argStr += " and InsertTime >= ?"
		}
		argList = append(argList, ps.StartTime)
	default:
		if isFirst {
			argStr += " InsertTime >= ? and InsertTime <= ?"
			isFirst = false
		} else {
			argStr += " and InsertTime >= ? and InsertTime <= ?"
		}
		argList = append(argList, ps.StartTime, ps.EndTime)
	}
	conditions := []string{}
	for _, appid := range appidlist {
		argList = append(argList, appid)
		conditions = append(conditions, "?")
	}
	if isFirst {
		argStr += " AppID in (" + strings.Join(conditions, ",") + ")"
		isFirst = false
	} else {
		argStr += " and AppID in (" + strings.Join(conditions, ",") + ")"
	}

	conn, err := db.ClickHouseCollection("")
	if err != nil {
		return
	}
	countCache = int64(0)
	err = conn.QueryRow(countStr+argStr, argList...).Scan(&countCache)
	if err != nil {
		return "", 0, nil, err
	}
	pagestart := (ps.PageIndex - 1) * ps.PageSize
	pageend := pagestart + ps.PageSize
	if pagestart > countCache {
		return
	}

	if pageend > countCache {
		pageend = countCache
	}
	if err != nil {
		return "", 0, nil, err
	}
	if ps.NextStartTime > 0 {
		if isFirst {
			argStr += " InsertTime <= ?"
			isFirst = false
		} else {
			argStr += " and InsertTime <= ?"
		}
		argList = append(argList, ps.NextStartTime)
	}
	if ps.NextOrderId != "" {
		if isFirst {
			argStr += " id < ?"
			isFirst = false
		} else {
			argStr += " and id < ?"
		}
		argList = append(argList, ps.NextOrderId)
	}
	if ps.UpEndTime > 0 {
		if isFirst {
			argStr += " InsertTime >= ?"
			isFirst = false
		} else {
			argStr += " and InsertTime >= ?"
		}
		argList = append(argList, ps.UpEndTime)
	}
	if ps.UpOrderId != "" {
		if isFirst {
			argStr += " id > ?"
			isFirst = false
		} else {
			argStr += " and id > ?"
		}
		argList = append(argList, ps.UpOrderId)
	}
	if tableName == "aviatorgamelogs" {
		if isFirst {
			argStr += " FinishType = ?"
			isFirst = false
		} else {
			argStr += " and FinishType = ?"
		}
		argList = append(argList, 1)
	}
	argList = append(argList, ps.PageSize, (ps.PageIndex-1)*ps.PageSize)
	queryStr = queryStr + argStr + " order by InsertTime desc LIMIT ? OFFSET ?"
	return queryStr, countCache, argList, err
}

func getBetLogListNew(ctx *Context, ps getBetLogListNewParams, ret *getBetLogListNewResults) (err error) {

	user, err := GetUser(ctx)
	if err != nil {
		return err
	}
	if ps.CurrencyCode != "" {
		ctx.CurrencyKey = ps.CurrencyCode
	}

	//appId := user.AppID
	appidlist := []string{}
	if user.GroupId != 3 {
		appidlist, err = GetOperatopAppID(ctx)
	} else {
		appidlist = append(appidlist, user.AppID)
	}
	if len(appidlist) == 0 {
		return lang.Error(ctx.Lang, "这个币种下没有商户")
	}

	if ps.OperatorId != 0 {
		var operator *comm.Operator
		err = CollAdminOperator.FindId(ps.OperatorId, &operator)
		if err != nil {
			return err
		}
		temp := false
		for _, appid := range appidlist {
			if operator.AppID == appid {
				temp = true
			}
		}
		if ps.CurrencyCode != "" {
			if ps.CurrencyCode != operator.CurrencyKey {
				return errors.New(lang.GetLang(ctx.Lang, "该商户未使用该币种"))
			}
		}
		if temp == false {
			return errors.New(lang.GetLang(ctx.Lang, "未代理该商户"))
		}
		appidlist = appidlist[:0]
		appidlist = append(appidlist, operator.AppID)

	}
	tableName := ""
	if ps.GameType == 0 {
		tableName = "gamelogs"
	} else {
		tableName = "aviatorgamelogs"
	}

	if len(appidlist) == 0 {
		return lang.Error(ctx.Lang, "参数错误")
	}
	list := []*getDocBetLog{}
	conn, err := db.ClickHouseCollection("")
	if err != nil {
		return
	}
	argList := []any{}
	queryStr := ""
	queryStr, ret.Count, argList, err = generateQueryStr(tableName, ps, appidlist)
	rows, err := conn.Query(queryStr, argList...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if ps.GameType == 0 {
			betlog := &getDocBetLog{}
			//if err = rows.Scan(&betlog.ID, &betlog.Pid, &betlog.UserID, &betlog.GameID, &betlog.GameType, &betlog.RoundID, &betlog.Bet, &betlog.Win, &betlog.Comment, &betlog.InsertTime, &betlog.Balance, &betlog.Completed, &betlog.TotalWinLoss, &betlog.WinLose, &betlog.SpinDetailsJson, &betlog.LogType, &betlog.PGBetID, &betlog.AppID, &betlog.CurrencyCode, &betlog.CurrencyName, &betlog.Grade, &betlog.HitBigReward, &betlog.LogType, &betlog.TransferAmount); err != nil {
			//	return err
			//}
			if err = rows.Scan(
				&betlog.ID,
				&betlog.Pid,
				&betlog.GameID,
				&betlog.TotalWinLoss,
				&betlog.PGBetID,
				&betlog.UserID,
				&betlog.RoundID,
				&betlog.Win,
				&betlog.Balance,
				&betlog.TransferAmount,
				&betlog.Comment,
				&betlog.Completed,
				&betlog.WinLose,
				&betlog.SpinDetailsJson,
				&betlog.GameType,
				&betlog.Bet,
				&betlog.InsertTime,
				&betlog.AppID,
				&betlog.Grade,
				&betlog.HitBigReward,
				&betlog.LogType,
				&betlog.ManufacturerName,
				&betlog.CurrencyKey,
				&betlog.UserName,
			); err != nil {
				return err
			}
			list = append(list, betlog)
		} else {
			betlog := &betLogAviator{}

			if err = rows.Scan(
				&betlog.ID,
				&betlog.Pid,
				&betlog.AppID,
				&betlog.Uid,
				&betlog.UserName,
				&betlog.CurrencyKey,
				&betlog.RoomId,
				&betlog.Bet,
				&betlog.Win,
				&betlog.Balance,
				&betlog.CashOutDate,
				&betlog.Payout,
				&betlog.Profit,
				&betlog.GameID,
				&betlog.RoundBetId,
				&betlog.RoundId,
				&betlog.MaxMultiplier,
				&betlog.Frb,
				&betlog.InsertTime,
				&betlog.RoundMaxMultiplier,
				&betlog.LogType,
				&betlog.FinishType,
				&betlog.ManufacturerName,
				&betlog.BetId,
			); err != nil {
				return err
			}
			betlogDoc := &getDocBetLog{
				ID:               betlog.ID,
				Pid:              uint64(betlog.Pid),
				AppID:            betlog.AppID,
				UserID:           betlog.Uid,
				GameID:           betlog.GameID,
				LogType:          betlog.LogType,
				Bet:              betlog.Bet,
				Win:              betlog.Win,
				Balance:          betlog.Balance,
				WinLose:          betlog.Profit,
				InsertTime:       betlog.InsertTime,
				RoundID:          strconv.FormatInt(betlog.RoundId, 10),
				ManufacturerName: betlog.ManufacturerName,
			}
			list = append(list, betlogDoc)
		}
	}
	//fmt.Println(err)

	if err != nil {
		return
	}

	//币种
	operck, err := GetCurrencyList(ctx)
	for i, log := range list {
		if ctx.Lang != "zh" {
			list[i].CurrencyName = operck[log.AppID].CurrencyCode
			list[i].CurrencyCode = operck[log.AppID].CurrencyCode
		} else {
			list[i].CurrencyName = operck[log.AppID].CurrencyName
			list[i].CurrencyCode = operck[log.AppID].CurrencyCode
		}
	}
	ret.List = list

	return err
}

var (
	checkedM sync.Map
)

func ensureIndex(coll *mongo.Collection) {
	key := coll.Name()
	if _, ok := checkedM.Load(key); ok {
		return
	}
	checkedM.Store(key, true)

	indexmodels := []mongo.IndexModel{
		{
			Keys: db.D("InsertTime", 1),
		}, {
			Keys: db.D("Pid", 1),
		}, {
			Keys: db.D("GameID", 1),
		}, {
			Keys: db.D("UserID", 1),
		},
	}
	coll.Indexes().CreateMany(context.TODO(), indexmodels)
}

func downloadBetLogListNew(ctx *Context, ps getBetLogListNewParams, ID string) (err error) {
	if ps.UserID == "" {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	if ps.StartTime != 0 && ps.EndTime != 0 {
		if ps.EndTime-ps.StartTime > 7*24*60*60 {
			return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
		}
	} else if ps.StartTime == 0 && ps.EndTime != 0 {
		e := time.Unix(ps.EndTime, 0)
		ps.StartTime = e.AddDate(0, 0, -7).Unix()
	} else if ps.StartTime != 0 && ps.EndTime == 0 {
		e := time.Unix(ps.StartTime, 0)
		ps.EndTime = e.AddDate(0, 0, 7).Unix()
	}
	ps.PageSize = 100 * 1000
	br := &getBetLogListNewResults{}
	err = getBetLogListNew(ctx, ps, br)
	if err != nil {
		return
	}
	gamesWrap := GetGameListResults{}
	gamesColl := NewOtherDB("game").Collection("Games")
	if err = gamesColl.FindAll(bson.M{}, &gamesWrap.List); err != nil {
		return err
	}
	gamesMap := make(map[string]string)
	for _, v := range gamesWrap.List {
		gamesMap[v.ID] = v.Name
	}
	// 创建 CSV 文件
	file, err := os.Create(fmt.Sprintf("%s/%s%s%s", BHdownloadPath, BHdownloadFilePre, ID, CSVSuf))
	if err != nil {
		return
	}
	defer file.Close()
	// 创建 CSV 写入器
	writer := csv.NewWriter(file)
	defer writer.Flush()

	title := []string{
		lang.Get(ctx.Lang, "交易订单号"),
		lang.Get(ctx.Lang, "运营商"),
		lang.Get(ctx.Lang, "唯一标识"),
		lang.Get(ctx.Lang, "玩家ID"),
		lang.Get(ctx.Lang, "记录生成时间"),
		lang.Get(ctx.Lang, "游戏"),
		lang.Get(ctx.Lang, "事件类型"),
		lang.Get(ctx.Lang, "转移金额"),
		lang.Get(ctx.Lang, "投注金额"),
		lang.Get(ctx.Lang, "赢分金额"),
		lang.Get(ctx.Lang, "玩家输赢金额"),
		lang.Get(ctx.Lang, "玩家余额")}
	err = writer.Write(title)
	logtype := []string{
		lang.Get(ctx.Lang, "电子投注/结算"),
		lang.Get(ctx.Lang, "+ 转入"),
		lang.Get(ctx.Lang, "- 转出"),
		lang.Get(ctx.Lang, "admin设置"),
		lang.Get(ctx.Lang, "彩票投注"),
		lang.Get(ctx.Lang, "彩票结算"),
		lang.Get(ctx.Lang, "百人投注/结算"),
	}
	for _, v := range br.List {
		row := make([]string, 0)
		row = append(row, v.ID, v.AppID, fmt.Sprintf("%d", v.Pid), v.UserID, v.InsertTime.Format("2006-01-02 15:04:05"))
		row = append(row, gamesMap[v.GameID], lang.Get(ctx.Lang, logtype[v.LogType]), fmt.Sprintf("%.2f", float64(v.TransferAmount)/10000))
		row = append(row, fmt.Sprintf("%.2f", float64(v.Bet)/10000), fmt.Sprintf("%.2f", float64(v.Win)/10000))
		row = append(row, fmt.Sprintf("%.2f", float64(v.WinLose)/10000), fmt.Sprintf("%.2f", float64(v.Balance)/10000))
		writer.Write(row)
	}

	// ZIP
	// var buffer bytes.Buffer
	// for _, str := range csvContent {
	// 	buffer.WriteString(str)
	// }
	// csvContentByte := buffer.Bytes()
	// var zipBuf bytes.Buffer
	// zw := zip.NewWriter(&zipBuf)
	// fh := &zip.FileHeader{
	// 	Name:   fmt.Sprintf("%s%s.csv", BHdownloadFilePre, ID),
	// 	Method: zip.Deflate, // 使用Deflate压缩
	// }

	// fw, err := zw.CreateHeader(fh)
	// if err != nil {
	// 	return
	// }
	// _, err = fw.Write(csvContentByte)
	// if err != nil {
	// 	return
	// }
	// if err = zw.Close(); err != nil {
	// 	return
	// }
	// err = os.WriteFile(fmt.Sprintf("%s/%s%s%s", BHdownloadPath, BHdownloadFilePre, ID, CSVSuf), zipBuf.Bytes(), 0666)
	// if err != nil {
	// 	return
	// }
	return
}

type DocBetLog struct {
	ID         string    `bson:"_id"`
	Pid        int64     `bson:"Pid"`    // 内部玩家ID
	UserID     string    `bson:"UserID"` // 外部玩家ID
	GameID     string    `bson:"GameID"` // 游戏ID
	GameType   int       `bson:"GameType"`
	RoundID    string    `bson:"RoundID"`    // 对局ID，方便查询
	Bet        int64     `bson:"Bet"`        // 下注
	Win        int64     `bson:"Win"`        // 输赢
	Comment    string    `bson:"Comment"`    // 注释
	InsertTime time.Time `bson:"InsertTime"` // 数据插入时间
	AppID      string    `bson:"AppID"`      //
	Balance    int64     `bson:"Balance"`    // 余额
	// RoundCount      int                `bson:"RoundCount"`
	Completed       bool   `bson:"Completed"`    // 最后一次spin
	TotalWinLoss    int64  `bson:"TotalWinLoss"` // 当前局结算总输赢, 只需要 Completed:true 才使用
	WinLose         int64  `bson:"WinLose"`      // 玩家输赢金额
	Grade           int    `bson:"Grade"`
	HitBigReward    int64  `bson:"HitBigReward"` // 奖池奖励
	SpinDetailsJson string `bson:"SpinDetailsJson"`
	PGBetID         string `bson:"PGBetID"`
	LogType         int    `bson:"LogType"`        // 0 游戏, 1 转入, 2 转出
	TransferAmount  int64  `bson:"TransferAmount"` // 转移金额
}

type Res struct {
	oper *comm.Operator //当前玩家运营商

	//如果  OperatorType 为 1 表示为线路商
	Line *comm.Operator //该运营商的代理商户

	betlog []*DocBetLog
}

func Getrrrk(plear player) (r Res, err error) {
	appid := plear.AppID

	err = CollAdminOperator.FindOne(bson.M{"AppID": appid}, &r.oper)
	if err != nil {
		return
	}
	if r.oper.Name != "admin" {
		err = CollAdminOperator.FindOne(bson.M{"AppID": r.oper.Name}, &r.Line)
		if err != nil {
			return
		}
	} else {
		r.Line = r.oper
	}

	query := bson.M{}

	query["AppID"] = plear.AppID
	query["Pid"] = plear.Id

	now := time.Now()
	fiveMinutesAgo := now.Add(-5 * time.Minute)
	date := bson.M{
		"$gte": fiveMinutesAgo, // 大于等于开始时间
		"$lte": now,            // 小于等于当前时间
	}
	query["InsertTime"] = date

	filter := mongodb.FindPageOpt{
		Page:     1,
		PageSize: 20,
		Sort:     bson.M{"InserTime": -1},
		Query:    query,
	}
	_, err = NewOtherDB("betlog").Collection(plear.AppID).FindPage(filter, &r.betlog)
	if err != nil {
		return Res{}, err
	}

	return
}
