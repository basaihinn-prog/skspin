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

// func InitConfig() {
// 	w := lazy.ConfigManager
// 	w.WatchUnmarshal("game_config.yaml", loadConfig)
// }

////////

type GameConfig struct {
	// 禁用的区域代码
	BlackLocs []string
}

func (config *GameConfig) IsBlockLoc(req *http.Request) (block bool, ip, loc string) {
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

var gameconfig atomic.Pointer[GameConfig]

func loadgameConfig(tmp *GameConfig) (err error) {
	gameconfig.Store(tmp)
	return nil
}

func GetGameConfig() *GameConfig {
	return gameconfig.Load()
}

func IsBlockLoc(req *http.Request) (block bool, ip, loc string) {
	return GetGameConfig().IsBlockLoc(req)
}
