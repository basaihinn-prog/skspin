package gamedata

import (
	"game/duck/lazy"
	"game/duck/ut2/httputil"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"sync/atomic"

	"github.com/phuslu/iploc"
)

func InitConfig() {
	w := lazy.ConfigManager
	w.WatchUnmarshal("game_config.yaml", loadConfig)
	w.WatchUnmarshal("admin_setting.yaml", loadSetting)
}

type Setting struct {
	OpenWhite bool       `yaml:"OpenWhite"`
	WhiteIps  []string   `yaml:"WhiteIps"`
	IconAddr  string     `yaml:"IconAddr"`
	ApiUrl    string     `yaml:"ApiUrl"`
	Alert     *AlertInfo `yaml:"Alert"`
}

type AlertInfo struct {
	Hour bool `yaml:"Hour"`
}

var Settings *Setting

type OneGameConfig struct {
	ID   string
	Name string
	Type string // Fish,Slot
}

////////

type Config struct {
	PGLaunchUrl      string
	JILILaunchUrl    string
	PPLaunchUrl      string
	TADALaunchUrl    string
	SPRIBELaunchUrl  string
	JDBLaunchUrl     string
	HACKRAWLaunchUrl string
	BaseUrl          string
	// PGLaunchUrlObj *url.URL `yaml:"-" json:"-" bson:"-"`
	LaunchUrl      string
	Websocket      string
	GameHistoryUrl string
	ReplaceGames   map[string]string
	// 禁用的区域代码
	BlackLocs []string
	// ReplaceDoMain by Appid
	ReplaceAppid  []string
	ReplaceDoMain map[string]string
}

func (config *Config) IsBlockLoc(req *http.Request) (block bool, ip, loc string) {
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

func (config *Config) GetGameId(old string) string {
	if config.ReplaceGames == nil {
		return old
	}

	if new := config.ReplaceGames[old]; new != "" {
		return new
	}
	return old
}

var config atomic.Pointer[Config]

func loadConfig(tmp *Config) (err error) {
	// tmp.PGLaunchUrlObj, err = url.Parse(tmp.PGLaunchUrl)
	// if err != nil {
	// 	return
	// }
	config.Store(tmp)
	return nil
}
func loadSetting(setting_ *Setting) error {
	Settings = setting_
	return nil
}

func Get() *Config {
	return config.Load()
}
