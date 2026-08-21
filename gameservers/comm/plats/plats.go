package main

import (
	"serve/comm/lazy"
	"serve/comm/mux"

	_ "serve/comm/plats/jili"
	_ "serve/comm/plats/pg"
	"serve/comm/plats/platcomm"
	_ "serve/comm/plats/pp"
)

func main() {
	// httpAddr, _ := lazy.RouteFile.Get("plats.http.api")
	// if httpAddr == "" {
	// 	return
	// }

	httpAddr := ":55555"

	platcomm.Start()

	mux.StartHttpServer(httpAddr)
	lazy.SignalProc()
}
