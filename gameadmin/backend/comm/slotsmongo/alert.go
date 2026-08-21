package slotsmongo

import (
	"context"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/ut"
	"game/duck/lazy"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func SendTelegramChat(msg string) (err error) {
	cfg := lazy.CommCfg().Alert
	// uri, _ := lazy.RouteFile.Get("telegramNotify")
	if cfg.TelegramNotify == "" {
		return
	}

	hostname, _ := os.Hostname()

	msg = fmt.Sprintf("%s\nfrom host: %s, pid: %d", msg, strings.ReplaceAll(hostname, "dou", "***"), os.Getpid())

	resp, err := http.Get(cfg.TelegramNotify + url.QueryEscape(msg))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	return
}

type BigRewardDealParams struct {
	Pid         int64
	AppID       string
	CurrencyKey string
	ServiceName string
	OriginWin   float64
	RealWin     float64
	WinLose     int64
	Desc        string
}

func BigRewardDeal(ps *BigRewardDealParams) {
	data := bson.M{
		"ID":        ps.Pid,
		"AppID":     ps.AppID,
		"Time":      time.Now(),
		"Currency":  ps.CurrencyKey,
		"GameID":    ps.ServiceName,
		"OriginWin": ut.Money2Gold(ps.OriginWin),
		"RealWin":   ut.Money2Gold(ps.RealWin),
		"WinLose":   ps.WinLose,
		"Desc":      ps.Desc,
	}
	db.Collection2("GameAdmin", "WinLoseLimit").InsertOne(context.TODO(), data)
	// 飞机预警
	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("时间：%s\n", time.Now().Format("20060102 15:04:05")))
	msg.WriteString(fmt.Sprintf("ID：%d\n", ps.Pid))
	msg.WriteString(fmt.Sprintf("运营商：%s\n", ps.AppID))
	msg.WriteString(fmt.Sprintf("游戏ID：%s\n", ps.ServiceName))
	msg.WriteString("中奖金额：" + comm.ThousandSeparator(fmt.Sprintf("%.2f\n", ps.OriginWin)))
	msg.WriteString(fmt.Sprintf("币种：%s\n", ps.CurrencyKey))
	msg.WriteString("实际金额：" + comm.ThousandSeparator(fmt.Sprintf("%.2f\n", ps.RealWin)))
	msg.WriteString("中奖时净盈亏：" + comm.ThousandSeparator(fmt.Sprintf("%.2f\n", float64(ps.WinLose)/10000)))
	msg.WriteString(fmt.Sprintf("触发限制条件：%s\n", ps.Desc))
	go SendTelegramChat(msg.String())
}

/*
import (
	"context"
	"fmt"
	"game/comm/ut"
	"game/duck/lazy"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Alert struct {
	WinLossSum     float64
	lastUpdateTime time.Time
	mtx            sync.Mutex
}

const (
	ClearDuration time.Duration = 10 * time.Minute
)

var (
	once  sync.Once
	alert *Alert
)




// 2.单个玩家在单款游戏盈利超20000会触发
// 【ID:xxxxxxx】，2023-11-20 09:03:01在XXXX，盈利XXXXX
func (a *Alert) onAddBetLog(betlog *DocBetLog) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	// a.WinLossSum += float64(betlog.WinLose / 10000)
	winMoney := ut.Gold2Money(betlog.WinLose)
	a.WinLossSum += winMoney

	now := time.Now()
	if winMoney > lazy.CommCfg().Alert.AlertSingleWinThreshold {
		// msg := fmt.Sprintf("下注:[%s]【ID:%d(%s)】，%s 在%s，盈利%+.2f",betlog.ID.Hex(), betlog.Pid, betlog.UserID, now.Format("2006-01-02 15:04:05"), lazy.ServiceName, winMoney)
		msg := fmt.Sprintf("【ID:%d(%s)】，%s 在%s，盈利%+.2f, 下注ID:[%s]", betlog.Pid, betlog.UserID, now.Format("2006-01-02 15:04:05"), lazy.ServiceName, winMoney, betlog.ID.Hex())
		err := SendTelegramChat(msg)

		slog.Info("SendTelegramChat", "msg", msg, "error", ut.ErrString(err))

		coll := ReportsCollection("alert")
		coll.InsertOne(context.TODO(), betlog)
	}
}

// 1.单款游戏，在10分钟区间内，总奖池输赢金额超10万会触发预警
// 总输赢后面是+表示系统赢
// 总输赢后面是-表示系统输
// 【这里是游戏的名字】，2023-11-20 09:00:00 - 2023-11-20 09:10:00,系统总输赢：±XXXXXX
func (a *Alert) checkPer10Min() {
	defer func() {
		if x := recover(); x != nil {
			os.Stdout.Write(debug.Stack())
			log.Println(x)
		}
	}()

	a.mtx.Lock()
	defer a.mtx.Unlock()

	now := time.Now()
	cfg := lazy.CommCfg()
	if cfg == nil {
		return
	}
	alert := cfg.Alert
	if alert == nil {
		return
	}
	if a.WinLossSum > alert.AlertGameThreshold {
		t1 := now.Add(-ClearDuration)
		msg := fmt.Sprintf("【%s】，%s - %s,系统总输赢：%+.2f", lazy.ServiceName, t1.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"), -a.WinLossSum)
		SendTelegramChat(msg)
	}
	a.WinLossSum = 0
}

func (a *Alert) Run() {
	for ; ; time.Sleep(ClearDuration) {
		a.checkPer10Min()
	}
}

*/
