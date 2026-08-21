package main

import (
	"log/slog"
	"sync"
	"time"

	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/comm/mux"

	"github.com/samber/lo"
)

func main() {
	lazy.Init("ip2world")
	// var m = map[string]Endpoint{}
	// var regions = []string{"id", "th", "in", "mm", "qa", "hk", "my", "vn", "sg"}
	// for i, r := range regions {
	// 	if i != 0 {
	// 		time.Sleep(time.Second)
	// 	}
	// 	endpoints, err := getEndpoints(r)
	// 	slog.Info("getEndpoints", "region", r, "err", err, "size", len(endpoints))
	// 	for _, ep := range endpoints {
	// 		m[ep.IP] = ep
	// 	}
	// }

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	mux.RegistRpcToMQ(mqconn)

	for {
		regions := "th"
		endpoints, err := getEndpoints(regions)
		slog.Info("getEndpoints", "region", regions, "err", err, "size", len(endpoints))

		if len(endpoints) == 0 {
			time.Sleep(time.Second)
			continue
		}

		mtx.Lock()
		g_endpoints = endpoints
		mtx.Unlock()

		time.Sleep(10 * time.Minute)
	}
}

func init() {
	mux.RegRpc("/ip2world/getEndpoint", "获取一个ip代理端点", "services", getEndpoint, nil)
}

var (
	g_endpoints []Endpoint
	mtx         sync.Mutex
)

func getEndpoint(_ mux.EmptyParams, ret *Endpoint) (err error) {
	mtx.Lock()
	defer mtx.Unlock()

	*ret = lo.Sample(g_endpoints)
	return
}
