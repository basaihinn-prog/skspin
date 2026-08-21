package gamedata

import (
	"game/comm/ut"
	"game/duck/lazy"
	"sync/atomic"
)

type Config struct {
	// MyHost           string
	ReverseProxy        map[string]string
	ReverseProxyUrls    ut.ReverseProxyUrlMap `yaml:"-"`
	BoBetDetailsTempUrl string
}

var config atomic.Pointer[Config]

func InitConfig() {
	w := lazy.ConfigManager
	w.WatchUnmarshal("ppgateway_config.yaml", loadConfig)

	// w.WatchUnmarshal("game_config.yaml", loadgameConfig)
}

func loadConfig(tmp *Config) (err error) {
	tmp.ReverseProxyUrls, err = ut.NewReverseProxyUrlMap(tmp.ReverseProxy)
	if err != nil {
		return
	}
	config.Store(tmp)
	return nil
}

func Get() *Config {
	return config.Load()
}
