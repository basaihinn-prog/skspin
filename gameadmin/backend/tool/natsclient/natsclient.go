package main

import (
	"flag"
	"game/comm/mq"
	"game/comm/slotsmongo"
	"game/service/gamecenter/msgdef"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

type addPs struct {
	X, Y int
}

type addRet struct {
	Sum int
}

func main() {
	// lazy.Init("natsclient")
	// mq.ConnectServerMust()

	// var addret addRet
	// err := mq.Invoke("admin.rpc.Add", addPs{rand.Intn(100), rand.Intn(100)}, &addret)
	// slog.Info("div",
	// 	"error", err,
	// 	"sum", addret.Sum,
	// )

	// var ans float64
	// err = mq.Invoke("admin.rpc.Div", [2]float64{1, 2}, &ans)
	// slog.Info("div",
	// 	"error", err,
	// 	"ans", ans,
	// )

	var pid int64
	flag.Int64Var(&pid, "pid", 100001, "pid in rp")
	flag.Parse()

	addr := "nats://127.0.0.1:11002"
	// addr = "nats://doudou-test:11002"
	nc := lo.Must(nats.Connect(addr))
	mq.SetNC(nc)

	// slotsmongo.GetBalance(100050)

	// nc.Subscribe("test.hello", func(msg *nats.Msg) {
	// 	msg.Respond([]byte("hello from" + time.Now().Format(time.RFC3339Nano)))
	// })

	for i := 0; i < 100; i++ {
		t1 := time.Now()
		var balance int64
		// wg := sync.WaitGroup{}
		for j := 0; j < 100; j++ {
			lo.Must(nc.Request("/gamecenter/player/getBalance", nil, 10*time.Second))
			balance = lo.Must(slotsmongo.ModifyGold(&msgdef.ModifyGoldPs{
				GameID:  "nats-test",
				Pid:     pid,
				Change:  1,
				RoundID: t1.String(),
				Reason:  "win",
			}))
			// wg.Add(1)

			// go func() {
			// 	slotsmongo.GetBalance(pid)
			// 	wg.Done()
			// }()

		}
		// wg.Wait()
		// balance = lo.Must(slotsmongo.GetBalance(pid))
		slog.Info("nats performance test", "i", i, "elapsed", time.Since(t1), "balance", balance)
	}

}
