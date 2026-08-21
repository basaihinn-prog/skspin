package staticproxy

import (
	"log"
	"net/http"
)

func StartProxy(addr string) {
	go func() {
		httpmux := http.NewServeMux()
		httpmux.HandleFunc("/api/play/authenticate", authenticate)
		httpmux.HandleFunc("/api/meta/gameInfo", gameInfo)
		httpmux.HandleFunc("/api/play/keepAlive", keepAlive)
		httpmux.HandleFunc("/api/play/bet", bet)
		httpmux.HandleFunc("/api/play/gameLaunch", keepAlive)
		//httpmux.HandleFunc("/api//history/GetHistoryForPlayer", keepAlive)
		httpmux.HandleFunc("/", static_assets)

		err := http.ListenAndServe(addr, httpmux)
		//err := http.ListenAndServe(addr, httpmux)
		log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}
