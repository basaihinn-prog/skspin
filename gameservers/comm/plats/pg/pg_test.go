package pg

import (
	"fmt"
	"net/url"
	"testing"

	"serve/comm/define"
)

func TestLaunch(t *testing.T) {
	// p := pg{}
	// p.LaunchGameHtml()

	cfg := GetConfig()
	uid := "123456"
	token, _ := GenToken(uid)

	extra_args := url.Values{}
	extra_args.Set("l", cfg.Lang)
	extra_args.Set("btt", "1")
	extra_args.Set("ops", token)

	game := "39"
	var result []byte
	invoke("/external-game-launcher/api/v1/GetLaunchURLHTML", define.M{
		"path":       fmt.Sprintf("/%s/index.html", game),
		"extra_args": extra_args.Encode(),
		"url_type":   "game-entry",
		// "client_ip":  ut.GetIPFromRequest(r),
		"client_ip": "127.0.0.1",
	}, &result)

	fmt.Println(string(result))

}
