package staticproxy

import (
	"log"
	"net/http"
)

func StartProxy(addr string) {
	go func() {
		httpmux := http.NewServeMux()
		httpmux.HandleFunc("GET /gs2c/html5Game.do", html5Game)
		httpmux.HandleFunc("GET /gs2c/lastGameHistory.do", lastGameHistory)
		httpmux.HandleFunc("GET /", static_assets)
		//httpmux.HandleFunc("POST /", static_assets2)
		httpmux.HandleFunc("POST  /gs2c/stats.do", func(w http.ResponseWriter, r *http.Request) {
			// Content-Type: application/javascript
			w.Header().Set("Content-Type", "application/javascript")
			w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Private-Network", "true")

			w.Write([]byte(`{"error":0,"description":"OK"}`))
		})
		httpmux.HandleFunc("POST /robots_center/create_task", robotsCenter)
		httpmux.HandleFunc("POST /robots_center_hgc/create_task", robotsCenter)
		httpmux.HandleFunc("POST /gs2c/ge/v3/gameService", gameService)
		httpmux.HandleFunc("POST /gs2c/ge/v4/gameService", gameService)
		httpmux.HandleFunc("POST /gs2c/saveSettings.do", saveSettings)
		httpmux.HandleFunc("GET /gs2c/reloadBalance.do", reloadBalance)

		httpmux.HandleFunc("POST /api/history/v2/settings/general", settings_general)
		httpmux.HandleFunc("POST /gs2c/api/history/v2/play-session/last-items", history_forward)
		httpmux.HandleFunc("POST /gs2c/api/history/v3/action/children", history_forward)
		httpmux.HandleFunc("POST /gs2c/api/history/v2/play-session/by-round", history_forward)
		httpmux.HandleFunc("POST /gs2c/playGame.do", playGame)
		httpmux.HandleFunc("GET /gs2c/announcements/unread/", Unread)
		httpmux.HandleFunc("GET /gs2c/promo/active/", PromoActive)
		httpmux.HandleFunc("GET /ReplayService/api/top/winnings/list", getWinList)
		httpmux.HandleFunc("GET /ReplayService/api/top/share/link", getShareLink)
		httpmux.HandleFunc("GET /ReplayService/api/top/replay/gameconfig", playGame)
		httpmux.HandleFunc("GET /ReplayServiceGlobal/api/replay/data", replayData)
		httpmux.HandleFunc("GET /ReplayService/replayGame.do", replayGame)
		httpmux.HandleFunc("POST /gs2c/api/history/v2/settings/general", general)
		httpmux.HandleFunc("GET /rpt/", replayGame)
		httpmux.HandleFunc("POST /game-rtp-api/", gamertpapi)

		// /gs2c/ge/v3/gameService
		// httpmux.HandleFunc("/", proxyhandle)
		// httpmux.Handle("/", newReverseProxy())

		//证书

		err := http.ListenAndServe(addr, httpmux)
		//err := http.ListenAndServe(addr, httpmux)
		log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}
