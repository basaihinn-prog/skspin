package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"log"
	"os"
	"os/signal"
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"
	"serve/comm/redisx"
	"serve/serviceSpribeGame/aviator/comm"
	"serve/serviceSpribeGame/aviator/gamedata"
	"serve/serviceSpribeGame/aviator/room"
	"serve/serviceSpribeGame/aviator/server"
)

func main() {
	lazy.Init(comm.GameID)
	gamedata.Load()
	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	mux.RegistRpcToMQ(mqconn)
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)
	// 新增redis连接
	redisAddr := lo.Must(lazy.RouteFile.Get("redis"))
	redisx.RegSubscribe(lazy.ServiceName, redisAddr)
	clickHouseAddr := lo.Must(lazy.RouteFile.Get("clickhouse"))
	db.DialToClickHouse(clickHouseAddr)
	s := server.NewServer()
	err := s.Start()
	if err != nil {
		log.Printf("nbio.Start failed: %v\n", err)
		return
	}

	DataDeal()
	defer s.Stop()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	fmt.Println("你好")
	<-interrupt
	log.Println("exit")
}

func DataDeal() {

	//msg := YesterdayTotal()
	//slog.Info("YesterDayTotal", "data", msg)
	//err := slotsmongo.SendTelegramChat(msg)
	//if err != nil {
	//	slog.Warn("SendTelegramChat", "error", ut.ErrString(err))
	//}

	//支持秒级
	//c := cron.New(cron.WithSeconds())

	c := cron.New()

	// 数据每天0点更新一次
	c.AddFunc("0 0 * * *", func() {
		rooms, _ := server.GetAppIds()
		deleteDays(rooms)
	})
	// 每月1号0点执行（非秒级精度）
	c.AddFunc("0 0 1 * *", func() {
		rooms, _ := server.GetAppIds()
		deleteMonths(rooms)
	})
	// 每年1月1日0点执行
	c.AddFunc("0 0 1 1 *", func() {
		rooms, _ := server.GetAppIds()
		deleteYears(rooms)
	})
	c.Start()
}

func deleteDays(rooms []string) {
	// 删除前一天数据
	room.DeleteTopData(rooms, "day")

}

func deleteMonths(rooms []string) {
	// 删除前一个月数据
	room.DeleteTopData(rooms, "month")
}

func deleteYears(rooms []string) {
	// 删除去年数据
	room.DeleteTopData(rooms, "year")
}
