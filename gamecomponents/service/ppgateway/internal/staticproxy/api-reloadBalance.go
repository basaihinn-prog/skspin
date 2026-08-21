package staticproxy

import (
	"fmt"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/ut2/jwtutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func reloadBalance(w http.ResponseWriter, r *http.Request) {
	// mgckey=AUTHTOKEN@17f8076abf3edb3c651be065210634dcd0b4ba1231c206e4405c0fa4455b9ec0~stylename@hllgd_hollygod~SESSION@5df8d3d1-dfa3-4cff-983f-6bd0853c162f~SN@273d3dc7

	// balance_bonus=0.00&balance=999980.00&balance_cash=999980.00&stime=1724918145560

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	mgckey := r.URL.Query().Get("mgckey")

	pid, err := jwtutil.ParseToken(mgckey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gold, err := slotsmongo.GetBalance(pid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance := ut.Gold2Money(gold)

	vs := url.Values{}
	vs.Set("balance_bonus", "0.00")
	vs.Set("balance", fmt.Sprintf("%.2f", balance))
	vs.Set("balance_cash", fmt.Sprintf("%.2f", balance))
	vs.Set("stime", strconv.Itoa(int(time.Now().UnixMilli())))

	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有源，或指定特定源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	w.Write([]byte(vs.Encode()))
}
