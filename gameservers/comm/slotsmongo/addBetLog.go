package slotsmongo

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"serve/comm/ut"
	"strconv"
	"time"

	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	//新加插槽
	UserName         string //用户名
	CurrencyKey      string //币种
	ManufacturerName string //厂商
}

func addBetLog(betlog *DocBetLog) {
	now := time.Now()
	if betlog.InsertTime.IsZero() {
		betlog.InsertTime = now
	}
	if betlog.ID == "" {
		sf := ut.NewSnowflake()
		betlog.ID = strconv.Itoa(int(sf.NextID()))
	}
	betlog.WinLose = betlog.Win - betlog.Bet
	// coll := ReportsCollection("BetLog")
	// coll.InsertOne(context.TODO(), betlog)

	appColl := db.Collection2("betlog", betlog.AppID)
	appColl.InsertOne(context.TODO(), betlog)

	// mq.Invoke("/alerter/betlog/add", betlog, nil)

	mq.JsonNC.Publish("/games/betlog", betlog)
	mq.JsonNC.Publish("/games/operator_balance_update", map[string]any{"AppID": betlog.AppID, "Change": ut.Gold2Money(betlog.WinLose), "Bet": ut.Gold2Money(betlog.Bet)})

	//mq.PublishMsg("/games/betlog", betlog)
}

// 添加记录

type AddBetLogParams struct {
	ID           string `bson:"_id"` // optional
	Pid          int64
	Bet          int64  // 下注
	Win          int64  // 输赢
	Balance      int64  // 余额
	RoundID      string // 必须和 ModifyGoldPs.RoundID 一致, 并且不同对局唯一
	Completed    bool   // 最后一次spin
	TotalWinLoss int64  // 当前局结算总输赢, 只需要 Completed:true 才使用
	IsBuy        bool
	Grade        int
	PGBetID      primitive.ObjectID
	BigReward    int64 // 奖池奖励
	GameType     int
	Extra        bool   // 额外下注
	AppID        string // optional
	Uid          string // optional
	ServiceName  string // optional
	Frb          bool   // 是否免费
	BonusCode    string

	//新加插槽
	UserName    string //用户名
	CurrencyKey string //币种
}

// func AddBetLog(pid int64, bet, win, balance int64, roundId string, completed bool, totalWinLoss int64,
func AddBetLog(ps *AddBetLogParams) /* (id primitive.ObjectID)*/ {

	var (
		pid          = ps.Pid
		bet          = ps.Bet
		win          = ps.Win
		balance      = ps.Balance
		roundId      = ps.RoundID
		completed    = ps.Completed
		totalWinLoss = ps.TotalWinLoss
		isBuy        = ps.IsBuy
		pgbetid      = ps.PGBetID
		bigReward    = ps.BigReward
		gameType     = ps.GameType
		grade        = ps.Grade
		extra        = ps.Extra
		appId        = ps.AppID
		uid          = ps.Uid
		serviceName  = cmp.Or(ps.ServiceName, lazy.ServiceName)
	)

	if appId == "" || uid == "" {
		var err error
		appId, uid, err = GetPlayerInfo(pid)
		if err != nil {
			slog.Error("GetPlayerInfo Fail!", "pid", pid, "error", err)
			return
		}
	}

	// jsonbuf, _ := json.Marshal(details)

	// var pgbetid string
	// if ut.IsPGGame(lazy.ServiceName) {
	// 	pgbetid = details.(string)
	// }

	playerDoc := []any{
		"BigReward", bigReward,
		"Bet", bet,
		"Win", win,
	}
	count := 0
	if completed {
		count = 1
	}

	if isBuy {
		playerDoc = append(playerDoc,
			fmt.Sprintf("TypeInfo.%d.BuyGame", gameType), count,
			fmt.Sprintf("TypeInfo.%d.BuyGameBet", gameType), bet,
			fmt.Sprintf("TypeInfo.%d.BuyGameWin", gameType), win)
	} else {
		playerDoc = append(playerDoc,
			fmt.Sprintf("TypeInfo.%d.Bet", gameType), bet,
			fmt.Sprintf("TypeInfo.%d.SpinCnt", gameType), count,
			fmt.Sprintf("TypeInfo.%d.Win", gameType), win)
	}

	db.Collection2("game", "Players").UpdateByID(context.TODO(), pid, db.D(
		"$inc", db.D(playerDoc...),
		"$set", db.D("LastSpinTime", time.Now().UTC().Format("2006-01-02T15:04:05.000Z")),
	))
	go addReport(serviceName, appId, pid, bet, win, completed, isBuy, extra, bigReward)
	addBetLog(&DocBetLog{
		ID:               ps.ID,
		Pid:              pid,
		UserID:           uid,
		GameID:           serviceName,
		GameType:         gameType,
		RoundID:          roundId,
		Bet:              bet,
		Win:              win,
		AppID:            appId,
		Balance:          balance,
		Completed:        completed,
		TotalWinLoss:     totalWinLoss,
		Grade:            grade,
		HitBigReward:     bigReward,
		PGBetID:          pgbetid.Hex(),
		UserName:         ps.UserName,
		CurrencyKey:      ps.CurrencyKey,
		ManufacturerName: ut.ServiceNameToManufacturerName(serviceName),
	})
	return
}

type AddBetLogParamsAviator struct {
	ID                 string `bson:"_id"` // optional
	Pid                int64
	AppID              string // optional
	Uid                string // optional
	UserName           string //用户名
	CurrencyKey        string //币种
	RoomId             string
	Bet                int64 // 下注
	Win                int64 // 输赢
	Balance            int64 // 余额
	GameId             string
	CashOutDate        int64     `json:"cashOutDate,omitempty"` // 提现时间
	Payout             float64   `json:"payout"`
	Profit             int64     `json:"profit,omitempty"` // 利润
	RoundBetId         string    `json:"roundBetId,omitempty"`
	RoundId            int64     `json:"roundId,omitempty"`
	MaxMultiplier      float64   `json:"maxMultiplier,omitempty"`
	Frb                uint8     // 是否免费
	InsertTime         time.Time `bson:"InsertTime"` // 数据插入时间
	RoundMaxMultiplier float64   `json:"roundMaxMultiplier,omitempty" bson:"roundMaxMultiplier" `
	LogType            int8      `json:"logType" bson:"logType"`
	FinishType         int8      `json:"finishType" bson:"finishType"`
	ManufacturerName   string    `json:"manufacturerName" bson:"manufacturerName"`
	BetId              int       `json:"betId" bson:"betId"`
	GameType           int
	Completed          bool `bson:"Completed"`
}

// func AddBetLog(pid int64, bet, win, balance int64, roundId string, completed bool, totalWinLoss int64,
func AddBetLogAviator(ps *AddBetLogParamsAviator) /* (id primitive.ObjectID)*/ {
	var (
		pid         = ps.Pid
		bet         = ps.Bet
		win         = ps.Win
		balance     = ps.Balance
		roundId     = ps.RoundId
		appId       = ps.AppID
		uid         = ps.Uid
		serviceName = lazy.ServiceName
		gameType    = ps.GameType
		completed   = ps.Completed
	)
	if appId == "" || uid == "" {
		var err error
		appId, uid, err = GetPlayerInfo(pid)
		if err != nil {
			slog.Error("GetPlayerInfo Fail!", "pid", pid, "error", err)
			return
		}
	}
	playerDoc := []any{
		"Bet", bet,
		"Win", win,
	}
	if completed {
		count := 1
		playerDoc = append(playerDoc,
			fmt.Sprintf("TypeInfo.%d.Bet", gameType), bet,
			fmt.Sprintf("TypeInfo.%d.SpinCnt", gameType), count,
			fmt.Sprintf("TypeInfo.%d.Win", gameType), win)

		db.Collection2("game", "Players").UpdateByID(context.TODO(), pid, db.D(
			"$inc", db.D(playerDoc...),
			"$set", db.D("LastSpinTime", time.Now().UTC().Format("2006-01-02T15:04:05.000Z")),
		))
		go addReport(serviceName, appId, pid, bet, win, completed, false, false, 0)
	}

	addBetLogAviator(&DocBetLogAviator{
		ID:                 ps.ID,
		Pid:                pid,
		AppID:              appId,
		Uid:                uid,
		UserName:           ps.UserName,
		CurrencyKey:        ps.CurrencyKey,
		RoomId:             ps.RoomId,
		Bet:                bet,
		Win:                win,
		Balance:            balance,
		CashOutDate:        ps.CashOutDate,
		Payout:             ps.Payout,
		Profit:             ps.Profit,
		RoundBetId:         ps.RoundBetId,
		RoundId:            roundId,
		MaxMultiplier:      ps.MaxMultiplier,
		Frb:                ps.Frb,
		GameID:             ps.GameId,
		RoundMaxMultiplier: ps.RoundMaxMultiplier,
		LogType:            ps.LogType,
		FinishType:         ps.FinishType,
		ManufacturerName:   ps.ManufacturerName,
		BetId:              ps.BetId,
		Completed:          completed,
	})
	return
}

type DocBetLogAviator struct {
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
	LogType            int8      `json:"logType" bson:"logType"`
	FinishType         int8      `json:"finishType" bson:"finishType"`
	ManufacturerName   string    `json:"manufacturerName" bson:"manufacturerName"`
	BetId              int       `json:"betId" bson:"betId"`
	Completed          bool      `bson:"Completed"`
}

func addBetLogAviator(betlog *DocBetLogAviator) {
	now := time.Now()
	if betlog.InsertTime.IsZero() {
		betlog.InsertTime = now
	}
	if betlog.ID == "" {
		sf := ut.NewSnowflake()
		betlog.ID = strconv.Itoa(int(sf.NextID()))
	}
	// coll := ReportsCollection("BetLog")
	// coll.InsertOne(context.TODO(), betlog)

	//appColl := db.Collection2("betlog", betlog.AppID)
	//appColl.InsertOne(context.TODO(), betlog)

	// mq.Invoke("/alerter/betlog/add", betlog, nil)

	err := mq.JsonNC.Publish("/games/minibetlog", betlog)
	if err != nil {
		slog.Error("插入历史记录失败", "betlog: ", betlog)
		slog.Error("addBetLogAviator Fail!", "err", err)
	}
	if betlog.Completed {
		//mq.JsonNC.Publish("/games/betlog", betlog)// todo 逻辑需要替换
		mq.JsonNC.Publish("/games/operator_balance_update", map[string]any{"AppID": betlog.AppID, "Change": ut.Gold2Money(betlog.Profit), "Bet": ut.Gold2Money(betlog.Bet)})
	}

	//mq.PublishMsg("/games/betlog", betlog)
}

func UpdateGamePlayerMulti(pid int64, multi float64) {
	playerDoc := []any{
		"Multi", multi,
	}
	db.Collection2("game", "Players").UpdateByID(context.TODO(), pid, db.D(
		"$inc", db.D(playerDoc...),
		"$set", db.D("LastSpinTime", time.Now().UTC().Format("2006-01-02T15:04:05.000Z")),
	))
}
