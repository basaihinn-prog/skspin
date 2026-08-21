package lazy

import (
	"game/comm/db"
	"game/duck/cfgmgr"
	"game/duck/etcd"
	"game/duck/exit"
	"game/duck/lang"
	"game/duck/lazy/env"
	"game/duck/logger"
	"game/duck/rpc1"
	"game/duck/rpc1/discovery"
	"game/duck/ut2"
	"game/duck/ut2/jwtutil"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"google.golang.org/grpc"
)

var ServiceName string

var PortProvider discovery.PortProvider
var ServiceDiscovery discovery.Discovery
var ServiceRegister discovery.Register

var ConfigManager *cfgmgr.ConfigManager
var GrpcClient *rpc1.ClientManager
var GrpcServer *rpc1.Server

// func init() {
// 	rand.Seed(time.Now().UnixMilli())
// }

func parseKV(kv string) (string, string) {
	cfgArr := strings.Split(kv, "://")
	return cfgArr[0], cfgArr[1]
}

var etcdMgr = ut2.NewSyncMap[string, *etcd.Etcd]()

func getEtcd(addr string) *etcd.Etcd {

	cli, ok := etcdMgr.Load(addr)
	if ok {
		return cli
	}

	cli, err := etcd.NewEtcd(addr)
	if err != nil {
		logger.Info("连接etcd失败", addr)
	}

	etcdMgr.Store(addr, cli)

	return cli
}

func initLogger() {
	switch env.Default.Logger {
	case "daily":
		logger.SetDaily("log", ServiceName)
	case "console":
		// 默认值，不做处理
	case "docker":
		// docker 这边日志 特殊定制
		l := logger.NewLogger()
		l.Use(
			logger.PartLevel(),
			logger.PartCaller(4, true),
			logger.PartMessage(),
		)
		logger.DefaultLogger = l
	}
}

// 初始化
func Init(serviceName string) {
	ServiceName = serviceName

	var conf = env.Default

	initLogger()

	logger.Info(ServiceName, "服务初始化")

	/////// /////// /////// /////// ///////
	logger.Info("使用配置中心", conf.ConfigWatcher)

	var fw cfgmgr.Watcher
	cfgK, cfgV := parseKV(conf.ConfigWatcher)
	if cfgK == "file" {
		fw = cfgmgr.NewFileWatcher(cfgV)
	} else if cfgK == "etcd" {
		fw = cfgmgr.NewEtcdWatcher(getEtcd(cfgV))
	}

	ConfigManager = cfgmgr.New(fw)
	exit.Close("关闭配置中心", ConfigManager)

	/////// /////// /////// /////// ///////
	langfile := "lang.csv"
	logger.Info("使用多语言文件", langfile)
	lang.Init(langfile, ConfigManager)
	/////// /////// /////// /////// ///////

	if conf.JwtKey != "" {
		jwtutil.SetKey(conf.JwtKey)
	}

	p := discovery.NewFileStaticDiscovery()
	RouteFile = p
	PortProvider = p
	ServiceDiscovery = p
	ServiceRegister = discovery.NewNullRegister()

	GrpcClient = rpc1.NewClientManager(ServiceDiscovery)
	// 与其他服务的链接最后关闭。可能仍需要调用一些东西
	exit.CloseWithPriority("关闭GRPC客户端", GrpcClient, -1)

	GrpcServer = rpc1.NewServer(&rpc1.ServerOption{Name: ServiceName}, PortProvider, ServiceRegister)
	// 服务入口首先关闭。避免新的请求进入
	exit.CloseWithPriority("关闭GRPC服务器", GrpcServer, 1)
}

func ServeFn(fn func()) {
	load()
	ConfigManager.LoadAll()
	go ConfigManager.Start()

	slog.Info("服务启动", "ServiceName", ServiceName)
	printInfo()

	if fn != nil {
		fn()
	}

	// select {} // 阻塞住进程。调用os.Exit主动退出
	SignalProc()
}

func SignalProc() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Block until a signal is received.
	s := <-c
	slog.Info("Got signal:", "signal", s)
	db.ClickHouseClose()
	os.Exit(0)
}

// 启动服务
func ServeWithoutGrpc() {
	ServeFn(nil)
}

// 启动服务
func Serve() {
	ServeFn(func() {
		err := GrpcServer.Serve()
		if err != nil {
			logger.Err(err)
		}
	})
}

func SetClientCodec(name string, codec grpc.DialOption) {
	_, err := GrpcClient.GetClient(name, codec)
	if err != nil {
		logger.Info("GRPC 连接出错", name, err)
		return
	}
}

// 创建并连接Grpc服务
func NewGrpcClient[T any](name string, fn func(cc grpc.ClientConnInterface) T) (ret T) {
	conn, err := GrpcClient.GetClient(name)
	if err != nil {
		logger.Info("GRPC 连接出错", name, err)
		return
	}

	return fn(conn)
}

// 获取地址
func GetAddr(sv string) string {
	if RouteFile != nil {
		addr, _ := RouteFile.Get(sv)
		return addr
	}

	addr, err := ServiceDiscovery.GetAddr(sv)
	if err != nil {
		return ""
	}
	if len(addr) > 0 {
		return addr[0]
	}
	return ""
}

func GetPortMust(sv string) int {

	port, err := PortProvider.GetPort(sv)
	if err != nil {
		panic(err)
	}

	return port
}
