package gamedata

import (
	"game/comm/ut"
	"game/duck/lazy"
	"sync/atomic"
)

type Config struct {
	ReverseProxy        map[string]string
	ReverseProxyUrls    ut.ReverseProxyUrlMap `yaml:"-"`
	ReverseProxy2       map[string]string
	ReverseProxyUrls2   ut.ReverseProxyUrlMap `yaml:"-"`
	ReverseProxy3       map[string]string
	ReverseProxyUrls3   ut.ReverseProxyUrlMap `yaml:"-"`
	ReverseProxy4       map[string]string
	ReverseProxyUrls4   ut.ReverseProxyUrlMap `yaml:"-"`
	BoBetDetailsTempUrl string
}

var config atomic.Pointer[Config]

func InitConfig() {
	w := lazy.ConfigManager
	w.WatchUnmarshal("jiligateway_config.yaml", loadConfig)
	w.WatchUnmarshal("jiligateway_games.json", loadGames)
	w.WatchUnmarshal("game_config.yaml", loadgameConfig)

	// w.WatchUnmarshal("game_config.yaml", loadgameConfig)
}

func loadConfig(tmp *Config) (err error) {
	// tmp.ReverseProxyUrls = make(map[string]*url.URL, len(tmp.ReverseProxy))
	// for host, remote := range tmp.ReverseProxy {
	// 	var u *url.URL
	// 	u, err = url.Parse(remote)
	// 	if err != nil {
	// 		return
	// 	}

	// 	tmp.ReverseProxyUrls[host] = u
	// }
	tmp.ReverseProxyUrls, err = ut.NewReverseProxyUrlMap(tmp.ReverseProxy)
	if err != nil {
		return
	}
	tmp.ReverseProxyUrls2, err = ut.NewReverseProxyUrlMap(tmp.ReverseProxy2)
	if err != nil {
		return
	}
	tmp.ReverseProxyUrls3, err = ut.NewReverseProxyUrlMap(tmp.ReverseProxy3)
	if err != nil {
		return
	}
	tmp.ReverseProxyUrls4, err = ut.NewReverseProxyUrlMap(tmp.ReverseProxy4)
	if err != nil {
		return
	}
	config.Store(tmp)
	return nil
}

func Get() *Config {
	return config.Load()
}

///////////////////////////////////////////////////////

type Game struct {
	No          int
	Id          string
	Name        string
	HistoryType string
	Multiple    float64
	Jackpot     bool
}

var gameMap atomic.Pointer[map[int]*Game]

func loadGames(tmp []*Game) (err error) {
	m := make(map[int]*Game, len(tmp))

	for _, v := range tmp {
		m[v.No] = v
	}

	gameMap.Store(&m)
	return
}

func GameMap() map[int]*Game {
	return *gameMap.Load()
}
