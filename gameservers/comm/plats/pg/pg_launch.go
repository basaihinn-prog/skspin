package pg

import (
	"fmt"
	"net/http"
	"net/url"

	"serve/comm/define"
)

func init() {
	http.HandleFunc("/pg/launch.html", onLaunchGame)
}
func onLaunchGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Cache-Control", "no-store")
	w.Header().Add("Cache-Control", "must-revalidate")

	qs := r.URL.Query()
	var (
		token = qs.Get("token")
		game  = qs.Get("game")
		lang  = qs.Get("lang")
		err   error
	)

	// cfg := GetConfig()

	extra_args := url.Values{}
	extra_args.Set("l", lang)
	extra_args.Set("btt", "1")
	extra_args.Set("ops", token)

	var result []byte
	err = invoke("/external-game-launcher/api/v1/GetLaunchURLHTML", define.M{
		"path":       fmt.Sprintf("/%s/index.html", game),
		"extra_args": extra_args.Encode(),
		"url_type":   "game-entry",
		// "client_ip":  ut.GetIPFromRequest(r),
		"client_ip": "127.0.0.1",
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

	token, _ := GenToken(uid)

	qs := url.Values{}
	qs.Add("token", token)
	qs.Add("uid", uid)
	qs.Add("game", game)
	qs.Add("lang", lang)
	url_ = cfg.OperatorLaunchUrl + "?" + qs.Encode()

	return
}
