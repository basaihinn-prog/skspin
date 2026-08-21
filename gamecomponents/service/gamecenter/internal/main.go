package internal

import (
	"context"
	"game/comm/db"
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/lazy"
	"game/service/gamecenter/internal/gamedata"
	"game/service/gamecenter/internal/gcdb"
	"game/service/gamecenter/internal/handlers"
	"game/service/gamecenter/internal/operator"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// var appMgr *AppManager
// var taskMgr *ut.TaskMgr

func Main() {
	lazy.Init("gamecenter")

	gamedata.InitConfig()

	gcdb.ConnectMongoByConfig()
	// 新增redis连接
	gcdb.InitRedis()
	gcdb.EnsureMongoIndex()

	operator.NewAppManagerAndInit()
	// taskMgr = ut.NewTaskMgr()

	// gamepb.RegisterGameRpcServer(lazy.GrpcServer, &GameRpcServer{})
	// gamepb.RegisterAdminGameCenterServer(lazy.GrpcServer, &AdminRpcServer{})
	// gamepb.RegisterAdminAnalysisServer(lazy.GrpcServer, &AdminAnalysisServer{})

	nconn := lo.Must(mq.ConnectServerMust())
	mux.DefaultRpcMux.RegistRpcToMQ(nconn)
	mux.SubscribeWSMsg(lazy.ServiceName, mq.GobNC())

	publicMux := mux.NewMux()
	publicMux.ServeMux.HandleFunc("/ping", ping)
	publicMux.ServeMux.HandleFunc("/config/token.json", handlers.HttpGetWSEndpoint)
	for _, h := range mux.DefaultRpcMux.ToArr() {
		if h.Class == "operator" || h.Class == "http" || h.Path == "/config/token.json" {
			publicMux.Add(h)
		}
	}

	addr := lo.Must(lazy.RouteFile.Get("gamecenter.http"))
	if len(addr) != 0 {
		fs := http.FileServer(http.Dir("./apimock/"))
		publicMux.ServeMux.Handle("/apimock/", http.StripPrefix("/apimock/", fs))
		publicMux.StartHttpServer(addr)
	}

	if addr, ok := lazy.RouteFile.Get("gamecenter.http.test"); ok {
		mux.StartHttpServer(addr)
	}

	handlers.RegSub()

	c := cron.New()
	c.AddFunc("30 5 * * *", deleteOrders)
	c.Start()
	// go initHttp()
	lazy.Serve()
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func deleteOrders() {
	coll := db.Collection2("game", "orders")
	before := time.Now().Add(-7 * 24 * time.Hour)
	fid := primitive.NewObjectIDFromTimestamp(before)

	coll.DeleteMany(context.TODO(), db.D(
		"_id", db.D("$lt", fid),
	))
}
