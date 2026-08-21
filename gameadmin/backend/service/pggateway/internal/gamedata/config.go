package gamedata

import (
	"game/comm/ut"
	"game/duck/lazy"
	"sync/atomic"
)

type Config struct {
	ApiUrlBase          string
	BoBetDetailsTempUrl string
	// PG                  PG
	ReverseProxy     map[string]string
	ReverseProxyUrls ut.ReverseProxyUrlMap `yaml:"-"`
}

var config atomic.Pointer[Config]

func InitConfig() {
	w := lazy.ConfigManager
	w.WatchUnmarshal("pggateway_config.yaml", loadConfig)

	w.WatchUnmarshal("game_config.yaml", loadgameConfig)
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
	config.Store(tmp)
	return nil
}

func Get() *Config {
	return config.Load()
}

// /////////////////////////////////////////////////////
type PG struct {
	PgSoftAPIDomain   string
	DataGrabAPIDomain string
	LaunchURL         string
	// 爬虫专用
	ClientApiURL      string
	OperatorToken     string
	SecretKey         string
	Lang              string
	OperatorLaunchUrl string
	PGLaunchUrl       string
}
