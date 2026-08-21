package gamedata

import (
	"game/duck/ut2/httputil"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"sync/atomic"

	"github.com/phuslu/iploc"
)

type Replaces struct {
	ReplaceGames map[string]string
	ReplaceLangs map[string]string
	// 禁用的区域代码
	BlackLocs []string
}

func (config *Replaces) IsBlockLoc(req *http.Request) (block bool, ip, loc string) {
	lg := slog.With("config.BlackLocs", config.BlackLocs)
	defer func() {
		if block {
			lg.Info("IsBlockLoc", "block", block)
		}
	}()

	ip = httputil.GetIPFromRequest(req)
	loc = string(iploc.Country(net.ParseIP(ip)))
	lg = lg.With("ip", ip, "loc", loc)
	if len(config.BlackLocs) == 0 {
		return false, ip, loc
	}
	return slices.Contains(config.BlackLocs, loc), ip, loc
}

func (config *Replaces) GetGameId(old string) string {
	if config.ReplaceGames == nil {
		return old
	}

	if new := config.ReplaceGames[old]; new != "" {
		return new
	}
	return old
}

func (config *Replaces) GetLang(gid, lang string) string {
	if lang == "zh" || lang == "" {
		return "en"
	}

	if config.ReplaceLangs == nil {
		return lang
	}

	if new := config.ReplaceLangs[gid+":"+lang]; new != "" {
		return new
	}
	return lang
}

var replaces atomic.Pointer[Replaces]

func loadReplaces(tmp *Replaces) (err error) {
	replaces.Store(tmp)
	return nil
}

func GetReplaces() *Replaces {
	return replaces.Load()
}
