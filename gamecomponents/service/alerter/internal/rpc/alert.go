package rpc

import (
	"context"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/lazy"
	"game/service/alerter/internal/gamedata"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/leekchan/accounting"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Alert struct {
	WinLossSum      float64
	GamesWinLossSum map[string]float64

	//币种的输赢变化
	CurrencyWinLossSum map[string]float64

	mtx sync.Mutex
}

var (
	alert = &Alert{
		GamesWinLossSum: map[string]float64{},
	}
)

func (a *Alert) getPlayerWinLoss(betlog *slotsmongo.DocBetLog) float64 {
	return ut.Gold2Money(betlog.TotalWinLoss)
}

func (a *Alert) getPlayerMiniWinLoss(betlog *slotsmongo.DocBetLogAviator) float64 {
	return ut.Gold2Money(betlog.Profit)
}
func udpatePlayerGameStatus(betlog *slotsmongo.DocBetLog) {
	coll := db.Collection2("game", "Players")
	update := db.D("$set", db.D("CompletedGames."+betlog.GameID, betlog.Completed))
	_, err := coll.UpdateByID(context.TODO(), betlog.Pid, update)
	if err != nil {
		// mongo.WriteError
		filter := db.D("_id", betlog.Pid, "CompletedGames", nil)
		update := db.D(
			"$set", db.D("CompletedGames", db.D(
				betlog.GameID, betlog.Completed),
			),
		)
		_, err := coll.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			slog.Error("udpatePlayerGameStatus", "err", ut.ErrString(err))
		}
	}
}

func udpatePlayerMiniGameStatus(betlog *slotsmongo.DocBetLogAviator) {
	coll := db.Collection2("game", "Players")
	update := db.D("$set", db.D("CompletedGames."+betlog.GameID, betlog.Completed))
	_, err := coll.UpdateByID(context.TODO(), betlog.Pid, update)
	if err != nil {
		// mongo.WriteError
		filter := db.D("_id", betlog.Pid, "CompletedGames", nil)
		update := db.D(
			"$set", db.D("CompletedGames", db.D(
				betlog.GameID, betlog.Completed),
			),
		)
		_, err := coll.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			slog.Error("udpatePlayerGameStatus", "err", ut.ErrString(err))
		}
	}
}

func getCurrency(appId string) string {
	type currencyItem struct {
		CurrencyKey string `bson:"CurrencyKey"`
	}
	var tmp currencyItem

	coll := db.Collection2("GameAdmin", "AdminOperator")

	coll.FindOne(context.TODO(), db.D("AppID", appId), options.FindOne().SetProjection(db.D("CurrencyKey", 1))).Decode(&tmp)

	if tmp.CurrencyKey == "" {
		tmp.CurrencyKey = "THB"
	}

	return tmp.CurrencyKey
}

func (a *Alert) onAddBetLog(betlog *slotsmongo.DocBetLog) (err error) {
	if betlog.LogType != 0 {
		return
	}

	udpatePlayerGameStatus(betlog)

	a.mtx.Lock()
	defer a.mtx.Unlock()

	//currency是货币
	currency := getCurrency(betlog.AppID)
	item := lazy.GetCurrencyItem(currency)

	//multi := item.Multi
	//if multi <= 0 {
	multi := 1.0
	//}

	// a.WinLossSum += float64(betlog.WinLose / 10000)
	winMoney := ut.Gold2Money(betlog.WinLose) / multi
	a.WinLossSum -= winMoney
	a.GamesWinLossSum[betlog.GameID] -= winMoney

	//统计币种的输赢
	currencyMoney := ut.Gold2Money(betlog.WinLose)

	a.CurrencyWinLossSum[currency] -= currencyMoney

	winloss := a.getPlayerWinLoss(betlog) / multi

	set := gamedata.Settings
	cond := set.GetSingleCond(betlog.AppID)
	lang := gamedata.GetLanguage(betlog.AppID)
	if betlog.Completed && winloss > set.SingleMin {
		var doc comm.Player
		coll := db.Collection2("game", "Players")
		err = coll.FindOne(context.TODO(), db.ID(betlog.Pid), options.FindOne().SetProjection(db.D("Bet", 1, "Win", 1, "Uid", 1))).Decode(&doc)
		if err != nil {
			return
		}

		bet, win := ut.Gold2Money(doc.Bet), ut.Gold2Money(doc.Win)

		if set.MatchPlayer(winloss, bet/multi, win/multi, cond) {
			// now := time.Now()
			msg := fmt.Sprintf("【ID:%d(%s)】\n"+
				gamedata.MapLanguage[lang].Profit+"\n"+
				gamedata.MapLanguage[lang].BetID+"\n"+
				gamedata.MapLanguage[lang].TotalProfitAndLoss+"\n"+
				gamedata.MapLanguage[lang].TotalBet+"\n"+
				gamedata.MapLanguage[lang].TotalRateOfReturn+"\n"+
				gamedata.MapLanguage[lang].Currency+"\n"+
				gamedata.MapLanguage[lang].CommercialOwner,
				betlog.Pid, betlog.UserID,
				betlog.GameID, fmtMoney(item.Symbol, a.getPlayerWinLoss(betlog)),
				betlog.ID,
				fmtMoney(item.Symbol, win-bet),
				fmtMoney(item.Symbol, bet),
				100*win/bet,
				currency,
				betlog.AppID,
			)
			fmt.Printf("set.SingleMin:%.2f,msg:%s", set.SingleMin, msg)
			go slotsmongo.SendTelegramChat(msg)
			go slotsmongo.SendOperatorTelegramChat(betlog.AppID, msg)
			coll := db.Collection2("reports", "alert")
			coll.InsertOne(context.TODO(), db.D(
				"pid", betlog.Pid,
				"userid", betlog.UserID,
				"gameId", betlog.GameID,
				"winMoney", fmtMoney(item.Symbol, winloss),
				"orderId", betlog.ID,
				"totalWinLoss", fmtMoney(item.Symbol, win-bet),
				"totalBet", fmtMoney(item.Symbol, bet),
				"winRate", fmt.Sprintf("%.2f%%", 100*win/bet),
				"currency", currency,
				"appId", betlog.AppID,
				"createTime", betlog.InsertTime.Local().UTC(),
			))
		}
	}

	return
}

func (a *Alert) onAddMiniBetLog(betlog *slotsmongo.DocBetLogAviator) (err error) {
	if betlog.LogType != 0 {
		return
	}

	udpatePlayerMiniGameStatus(betlog)

	a.mtx.Lock()
	defer a.mtx.Unlock()

	//currency是货币
	currency := getCurrency(betlog.AppID)
	item := lazy.GetCurrencyItem(currency)

	//multi := item.Multi
	//if multi <= 0 {
	multi := 1.0
	//}

	// a.WinLossSum += float64(betlog.WinLose / 10000)
	winMoney := ut.Gold2Money(betlog.Profit) / multi
	a.WinLossSum -= winMoney
	a.GamesWinLossSum[betlog.GameID] -= winMoney

	//统计币种的输赢
	currencyMoney := ut.Gold2Money(betlog.Profit)

	a.CurrencyWinLossSum[currency] -= currencyMoney

	winloss := a.getPlayerMiniWinLoss(betlog) / multi

	set := gamedata.Settings
	cond := set.GetSingleCond(betlog.AppID)
	lang := gamedata.GetLanguage(betlog.AppID)
	if betlog.Completed && winloss > set.SingleMin {
		var doc comm.Player
		coll := db.Collection2("game", "Players")
		err = coll.FindOne(context.TODO(), db.ID(betlog.Pid), options.FindOne().SetProjection(db.D("Bet", 1, "Win", 1, "Uid", 1))).Decode(&doc)
		if err != nil {
			return
		}

		bet, win := ut.Gold2Money(doc.Bet), ut.Gold2Money(doc.Win)

		if set.MatchPlayer(winloss, bet/multi, win/multi, cond) {
			// now := time.Now()
			msg := fmt.Sprintf("【ID:%d(%s)】\n"+
				gamedata.MapLanguage[lang].Profit+"\n"+
				gamedata.MapLanguage[lang].BetID+"\n"+
				gamedata.MapLanguage[lang].TotalProfitAndLoss+"\n"+
				gamedata.MapLanguage[lang].TotalBet+"\n"+
				gamedata.MapLanguage[lang].TotalRateOfReturn+"\n"+
				gamedata.MapLanguage[lang].Currency+"\n"+
				gamedata.MapLanguage[lang].CommercialOwner,
				betlog.Pid, betlog.Uid,
				betlog.GameID, fmtMoney(item.Symbol, a.getPlayerMiniWinLoss(betlog)),
				betlog.ID,
				fmtMoney(item.Symbol, win-bet),
				fmtMoney(item.Symbol, bet),
				100*win/bet,
				currency,
				betlog.AppID,
			)
			fmt.Printf("set.SingleMin:%.2f,msg:%s", set.SingleMin, msg)
			go slotsmongo.SendTelegramChat(msg)
			go slotsmongo.SendOperatorTelegramChat(betlog.AppID, msg)
			coll := db.Collection2("reports", "alert")
			coll.InsertOne(context.TODO(), db.D(
				"pid", betlog.Pid,
				"userid", betlog.Uid,
				"gameId", betlog.GameID,
				"winMoney", fmtMoney(item.Symbol, winloss),
				"orderId", betlog.ID,
				"totalWinLoss", fmtMoney(item.Symbol, win-bet),
				"totalBet", fmtMoney(item.Symbol, bet),
				"winRate", fmt.Sprintf("%.2f%%", 100*win/bet),
				"currency", currency,
				"appId", betlog.AppID,
				"createTime", betlog.InsertTime.Local().UTC(),
			))
		}
	}

	return
}

func (a *Alert) systemCheck(dur time.Duration) {
	defer func() {
		if x := recover(); x != nil {
			os.Stdout.Write(debug.Stack())
			log.Println(x)
		}
	}()

	a.mtx.Lock()
	defer a.mtx.Unlock()

	now := time.Now()
	set := gamedata.Settings

	if a.WinLossSum < set.SystemMin || set.SystemMax < a.WinLossSum {
		t1 := now.Add(-dur)
		msg := fmt.Sprintf(`【平台：RP】自%s 到%s 系统收益波动为%+.2f。`, t1.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"), a.WinLossSum)
		fmt.Println("", msg)
		slotsmongo.SendTelegramChat(msg)
	}
	a.WinLossSum = 0
	// dur = time.Duration(set.SystemInternalSec) * time.Second

}

// 币种的报警
func (a *Alert) currencyCheck(dur time.Duration) {

	fmt.Println("进入币种报警...")

	//程序出现问题的报警写入日志
	defer func() {
		if x := recover(); x != nil {
			os.Stdout.Write(debug.Stack())
			log.Println(x)
		}
	}()

	a.mtx.Lock()
	defer a.mtx.Unlock()

	now := time.Now()

	//获得币种的配置
	set := gamedata.Settings

	fmt.Println("币种集合:", a.CurrencyWinLossSum)

	for currency, winloss := range a.CurrencyWinLossSum {
		//遍历配置中币种的配置

		fmt.Println("币种:", currency)

		currencySetting := set.GetSettingCurrency(currency)

		fmt.Println("currencySetting:", currencySetting)

		if currencySetting.Currency == "" {
			continue
		}
		if winloss > currencySetting.CurrencyMax || winloss < currencySetting.CurrencyMin {
			t1 := now.Add(-dur)
			msg := fmt.Sprintf(`【币种: %s】自%s 到%s 系统收益波动为%+.2f。`, currency, t1.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"), winloss)
			fmt.Println(msg)
			go slotsmongo.SendTelegramChat(msg)
		}
	}
	a.CurrencyWinLossSum = map[string]float64{}
}

func (a *Alert) runSystemCheck() {
	dur := time.Duration(gamedata.Settings.SystemInternalSec) * time.Second
	for ; ; time.Sleep(dur) {
		a.systemCheck(dur)
		dur = time.Duration(gamedata.Settings.SystemInternalSec) * time.Second
	}
}

func (a *Alert) runCurrencyCheck() { //币种的定时监控
	dur := time.Duration(gamedata.Settings.SystemInternalSec) * time.Second
	for ; ; time.Sleep(dur) {
		a.currencyCheck(dur)
		dur = time.Duration(gamedata.Settings.SystemInternalSec) * time.Second
	}
}

func (a *Alert) gameCheck(dur time.Duration) {
	defer func() {
		if x := recover(); x != nil {
			os.Stdout.Write(debug.Stack())
			log.Println(x)
		}
	}()

	a.mtx.Lock()
	defer a.mtx.Unlock()

	now := time.Now()
	set := gamedata.Settings

	games := getGameList()

	for game, winloss := range a.GamesWinLossSum {
		if winloss < set.GameMin || set.GameMax < winloss {
			t1 := now.Add(-dur)
			msg := fmt.Sprintf(`【游戏：%s, ID: %s】自%s 到%s 系统收益波动为%+.2f。`, games[game], game, t1.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"), winloss)
			go slotsmongo.SendTelegramChat(msg)
		}
		a.GamesWinLossSum[game] = 0
	}
}

func getGameList() (games map[string]string) {
	games = map[string]string{}
	coll := db.Collection2("game", "Games")
	cur, err := coll.Find(context.TODO(), db.D(), options.Find().SetProjection(db.D("Name", 1)))
	if err != nil {
		return
	}

	var list []struct {
		ID   string `json:"ID" bson:"_id"`
		Name string `json:"Name" bson:"Name"`
	}
	cur.All(context.TODO(), &list)

	for _, v := range list {
		games[v.ID] = v.Name
	}

	return
}

func (a *Alert) runGameCheck() {
	dur := time.Duration(gamedata.Settings.GameInternalSec) * time.Second
	for ; ; time.Sleep(dur) {
		a.gameCheck(dur)
		dur = time.Duration(gamedata.Settings.GameInternalSec) * time.Second
	}
}

func Start() {
	//go alert.runSystemCheck()
	// go alert.runGameCheck()
	go alert.runCurrencyCheck()
}

// var (
// 	ac = sync.OnceValue(func() *accounting.Accounting {
// 		// return &accounting.Accounting{Symbol: lazy.CurrencySymbol(), Precision: 2}
// 		return &accounting.Accounting{Symbol: /*lazy.CurrencySymbol()*/ "", Precision: 2}
// 	})
// )

func fmtMoney(symbol string, f float64) string {
	ac := accounting.Accounting{Symbol: symbol, Precision: 2}
	return ac.FormatMoneyFloat64(f)
}

func (a *Alert) onOpenLottery(lotteryAlert *slotsmongo.LotteryAlert) (err error) {
	msg := fmt.Sprintf(`%s 彩票一期共计销售%v，玩家中奖%v，公司盈利%v；
彩票二期共计销售%v，玩家中奖%v，公司盈利%v
本期彩票共计盈利%v`, lotteryAlert.OpenDay, lotteryAlert.PeriodSoldGold/10000,
		lotteryAlert.PeriodProfits/10000,
		lotteryAlert.CompanyWin_1/10000,
		lotteryAlert.PeriodSoldGuessGold/10000,
		lotteryAlert.PeriodGuessProfits/10000,
		lotteryAlert.CompanyWin_2/10000,
		lotteryAlert.CompanyWin_total/10000,
	)
	fmt.Println(lotteryAlert.OpenDay, " alert onOpenLottery:", msg)
	slotsmongo.SendTelegramChat(msg)
	return
}
