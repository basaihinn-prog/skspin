package slotsmongo

import (
	"context"
	"fmt"
	"game/comm/db"
	"game/duck/lazy"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"net/url"
)

func SendTelegramChat(msg string) (err error) {
	cfg := lazy.CommCfg().Alert
	// uri, _ := lazy.RouteFile.Get("telegramNotify")
	if cfg.TelegramNotify == "" {
		return
	}

	resp, err := http.Get(cfg.TelegramNotify + url.QueryEscape(msg))
	if err != nil {
		fmt.Printf("http request err:%v", err)
		return
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	return
}

func SendOperatorTelegramChat(appId, msg string) (err error) {
	var TelegramChat struct {
		Robot  string `bson:"Robot"`
		ChatID string `bson:"ChatID"`
	}

	coll := db.Collection2("GameAdmin", "AdminOperator")

	coll.FindOne(context.TODO(), db.D("AppID", appId), options.FindOne().SetProjection(db.D("Robot", 1, "ChatID", 1))).Decode(&TelegramChat)

	if TelegramChat.Robot == "" || TelegramChat.ChatID == "" {
		return
	}
	telegramChatUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=", TelegramChat.Robot, TelegramChat.ChatID)
	resp, err := http.Get(telegramChatUrl + url.QueryEscape(msg))
	if err != nil {
		fmt.Printf("http request err:%v", err)
		return
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	return
}

func SendOperatorBalanceTelegramChat(appId, msg string) (err error) {
	var TelegramChat struct {
		AdminRobot  string `bson:"AdminRobot"`
		AdminChatID string `bson:"AdminChatID"`
	}

	coll := db.Collection2("GameAdmin", "AdminOperator")

	coll.FindOne(context.TODO(), db.D("AppID", appId), options.FindOne().SetProjection(db.D("AdminRobot", 1, "AdminChatID", 1))).Decode(&TelegramChat)

	if TelegramChat.AdminRobot == "" || TelegramChat.AdminChatID == "" {
		return
	}
	telegramChatUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=", TelegramChat.AdminRobot, TelegramChat.AdminChatID)
	resp, err := http.Get(telegramChatUrl + url.QueryEscape(msg))
	if err != nil {
		fmt.Printf("http request err:%v", err)
		return
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	return
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
