package lazy

import (
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

var (
	ServiceName string
	Line        int
)

// var PortProvider discovery.PortProvider
// var ServiceDiscovery discovery.Discovery
// var ServiceRegister discovery.Register

var GConfigManager *ConfigManager

// var GrpcClient *rpc1.ClientManager
// var GrpcServer *rpc1.Server

func SignalProc() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Block until a signal is received.
	s := <-c
	slog.Info("Got signal:", "signal", s)
	os.Exit(0)
}

// 启动服务
func Serve() {
	load()
	GConfigManager.LoadAll()
	go GConfigManager.Start()

	pid()
	slog.Info("服务启动", "ServiceName", ServiceName)

	// select {} // 阻塞住进程。调用os.Exit主动退出
	SignalProc()
}

// write pid
func pid() {
	pid := os.Getpid()
	executable, _ := os.Executable()
	processName := filepath.Base(executable)
	pidStr := strconv.Itoa(pid)

	if err := os.WriteFile(processName+".pid", []byte(pidStr), 0644); err != nil {
		slog.Info("pid file error")
		return
	}
}
