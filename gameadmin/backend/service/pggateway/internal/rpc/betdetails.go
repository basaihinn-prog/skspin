package rpc

import (
	"game/comm/mux"
	"game/service/pggateway/internal/gamedata"
	"net/url"
	"path"
	"strings"

	"github.com/samber/lo"
)

func init() {
	mux.RegRpc("/pggateway/getBetDetailsUrl", "gameinfo", "game-api", getBetDetailsUrl, getBetDetailsUrlPs{
		Gid:   "1489936",
		BetID: "660cb855859579949c5a0402",
		Lang:  "en",
		Token: "admin-token",
	})
}

type getBetDetailsUrlPs struct {
	Gid   string
	BetID string
	SID   string
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
	gid := strings.TrimPrefix(ps.Gid, "pg_")

	u.Path = path.Join("/history", gid+".html")

	q := u.Query()
	q.Set("psid", ps.BetID)
	q.Set("sid", lo.Ternary(ps.SID == "", "1", ps.SID))
	q.Set("gid", gid)
	q.Set("lang", ps.Lang)
	q.Set("t", ps.Token)

	u.RawQuery = q.Encode()

	ret.Url = u.String()
	return
}
