package main

import (
	"game/duck/lazy"
	"log/slog"
	"net"
	"slices"
	"sync/atomic"

	"github.com/phuslu/iploc"
)

func LoadCfg() {
	// cm := lazy.ConfigManager
	// cm.WatchAndLoad("admin_rpc.json", LoadRpcConfig)
	lazy.ConfigManager.WatchUnmarshal("admin_setting.yaml", loadSetting)
	lazy.ConfigManager.WatchUnmarshal("game_config.yaml", loadgameConfig)
	//lazy.ConfigManager.WatchUnmarshal("output.json", loadGameNow)
	//lazy.ConfigManager.WatchUnmarshal("GameBet.json", loadGameBet)
	lazy.ConfigManager.WatchUnmarshal("demoGame.json", loadDemoGame)
	lazy.ConfigManager.WatchUnmarshal("theme.json", loadDemoTheme)
	lazy.ConfigManager.WatchUnmarshal("lang.json", loadLange)
}

type Setting struct {
	OpenWhite           bool       `yaml:"OpenWhite"`
	DefaultPermissionId int64      `yaml:"DefaultPermissionId"`
	WhiteIps            []string   `yaml:"WhiteIps"`
	IconAddr            string     `yaml:"IconAddr"`
	ApiUrl              string     `yaml:"ApiUrl"`
	Alert               *AlertInfo `yaml:"Alert"`
}

type AlertInfo struct {
	Hour bool `yaml:"Hour"`
}

type GameConfig struct {
	// 禁用的区域代码
	BlackLocs []string
}
type RTPConfig struct {
	// 禁用的区域代码
	RTPList []RTPData
}
type RTPData struct {
	// 禁用的区域代码
	RTP            int64
	RewardPercent  int64
	NoAwardPercent int64
}
type GameNow struct {
	GameID        string `json:"gameId"`
	RewardPercent int64  `json:"reward_percent"`
}

type GameBet struct {
	GameID          string  `json:"Game_id"`
	Bet             string  `json:"Bet"`
	GameName        string  `json:"GameName"`
	DefaultBet      float64 `json:"DefaultBet"`
	DefaultBetLevel int64   `json:"DefaultBetLevel"`
}

var MapGameNow = map[string]int64{}
var MapGameBet = map[string]*GameBet{}
var Language []map[string]string
var DemoGameList []*DemoGame
var DemoThemeList []map[string]string

var setting *Setting
var gameconfig atomic.Pointer[GameConfig]

func loadSetting(setting_ *Setting) error {
	setting = setting_
	return nil
}

func loadGameNow(setting_ []*GameNow) error {
	for _, g := range setting_ {
		MapGameNow[g.GameID] = g.RewardPercent
	}
	return nil
}

func loadGameBet(settint_ []*GameBet) error {
	//map init
	MapGameBet = map[string]*GameBet{}
	for _, set := range settint_ {
		MapGameBet[set.GameID] = set
	}
	return nil
}

// 获取试玩站游戏列表
func loadDemoTheme(themeList []map[string]string) error {
	DemoThemeList = themeList
	return nil
}
func loadDemoGame(gameList []*DemoGame) error {
	DemoGameList = gameList
	return nil
}

func loadLange(systemLanguages []map[string]string) error {
	//map init
	Language = systemLanguages
	return nil
}

func loadgameConfig(tmp *GameConfig) (err error) {
	gameconfig.Store(tmp)
	return nil
}

func GetGameConfig() *GameConfig {
	return gameconfig.Load()
}

func (config *GameConfig) IsBlockLoc(ip string) (block bool) {
	lg := slog.With("config.BlackLocs", config.BlackLocs)
	defer func() {
		if block {
			lg.Info("IsBlockLoc", "block", block)
		}
	}()

	if len(config.BlackLocs) == 0 {
		return false
	}

	loc := string(iploc.Country(net.ParseIP(ip)))
	lg = lg.With("ip", ip, "loc", loc)
	return slices.Contains(config.BlackLocs, loc)
}

func IsBlockLoc(ip string) (block bool) {
	return GetGameConfig().IsBlockLoc(ip)
}

func initRTPConfig() {

}

/*
// 配置rpc走proxy转发，用于本地调试
var useProxyMap = ut2.NewSyncMap[string, bool]()

func LoadRpcConfig(buf []byte) error {

	mp := map[string]any{}
	err := json.Unmarshal(buf, &mp)
	if err != nil {
		return err
	}

	useProxyMap.Clear()
	for k := range mp {
		useProxyMap.Store(k, true)
	}

	return nil
}
*/
