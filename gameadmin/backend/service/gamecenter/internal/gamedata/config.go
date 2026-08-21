package gamedata

import (
	"game/duck/lazy"
	"sync/atomic"
)

func InitConfig() {
	w := lazy.ConfigManager
	w.WatchUnmarshal("game_config.yaml", loadConfig)
	w.WatchUnmarshal("admin_setting.yaml", loadSetting)
	w.WatchUnmarshal("game_replaces.yaml", loadReplaces)
	w.WatchUnmarshal("gamecenter_transferlimit.yaml", loadTransferLimitMap)
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
	PGLaunchUrl     string
	PGLaunchUrlHttp string
	JILILaunchUrl   string
	PPLaunchUrl     string
	LaunchUrl       string
	Websocket       string
	GameHistoryUrl  string
	// ReplaceGames    map[string]string
	// ReplaceLangs    map[string]string
	// 禁用的区域代码
	// BlackLocs []string
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
