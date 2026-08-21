package platpg

import (
	"fmt"
	"game/comm/define"
	"game/duck/ut2/httputil"
	"net/http"
	"net/url"
)

func init() {
	http.HandleFunc("/pg/launch.html", onLaunchGame)
}
func onLaunchGame(w http.ResponseWriter, r *http.Request) {
	// Note: Operator should include the following (headers) to prevent the response from being stored in
	// the web browser cache:
	// Cache-Control: no-cache, no-store, must-revalidate
	// w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Cache-Control", "no-store")
	w.Header().Add("Cache-Control", "must-revalidate")

	// s := "visit by " + r.URL.String()
	// w.Write([]byte(s))

	// visit by /pg/launch.html?game=1529867&token=0001020304293b2808090b80c4ab602b&uid=123456

	qs := r.URL.Query()
	var (
		token = qs.Get("token")
		game  = qs.Get("game")
		err   error
	)

	cfg := GetConfig()

	extra_args := url.Values{}
	extra_args.Set("l", cfg.Lang)
	extra_args.Set("btt", "1")
	// extra_args.Set("ot", cfg.OperatorToken)
	extra_args.Set("ops", token)

	var result []byte
	ip := httputil.GetIPFromRequest(r)
	err = invoke("/external-game-launcher/api/v1/GetLaunchURLHTML", define.M{
		"path":       fmt.Sprintf("/%s/index.html", game),
		"extra_args": extra_args.Encode(),
		"url_type":   "game-entry",
		"client_ip":  ip,
	}, &result)

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(result)
}

func (pg) LaunchGameHtml(uid string, game, lang string) (url_ string, err error) {
	cfg := GetConfig()

	// token := platcomm.GenAuthToken(uid)
	token, err := GenToken(uid)

	qs := url.Values{}
	qs.Add("token", token)
	qs.Add("game", game)
	qs.Add("lang", lang)
	url_ = cfg.OperatorLaunchUrl + "?" + qs.Encode()

	return
}
