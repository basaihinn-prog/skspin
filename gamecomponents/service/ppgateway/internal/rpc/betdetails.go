package rpc

import (
	"cmp"
	"game/comm/mux"
	"game/service/ppgateway/internal/gamedata"
	"net/url"
)

func init() {
	mux.RegRpc("/ppgateway/getBetDetailsUrl", "gameinfo", "game-api", getBetDetailsUrl, getBetDetailsUrlPs{
		Gid:   "pp_vs20olympx",
		BetID: "66ebbb292235a2b8b233d713",
		Lang:  "en",
	})
}

type getBetDetailsUrlPs struct {
	Gid   string
	BetID string
	Lang  string
	Token string
}
type getBetDetailsUrlRet struct {
	Url string
}

func getBetDetailsUrl(ps getBetDetailsUrlPs, ret *getBetDetailsUrlRet) (err error) {
	tmpl := gamedata.Get().BoBetDetailsTempUrl
	u, err := url.Parse(tmpl)
	if err != nil {
		return
	}
	q := u.Query()
	q.Set("playSessionId", ps.BetID)
	q.Set("symbol", ps.Gid[len("pp_"):])
	q.Set("lang", cmp.Or(ps.Lang, "en"))

	u.RawQuery = q.Encode()
	ret.Url = u.String()

	return
}
