package lazy

import (
	"game/duck/cfgmgr"
	"game/duck/lang"
	"game/duck/rpc1/discovery"
	"log/slog"
)

var RouteFile *discovery.FileStaticDiscovery

func InitWithoutGrpc(serviceName string) {
	ServiceName = serviceName

	// var conf = env.Default

	slog.Info("服务初始化", "ServiceName", ServiceName)

	fw := cfgmgr.NewFileWatcher("config")
	ConfigManager = cfgmgr.New(fw)

	/////// /////// /////// /////// ///////
	langfile := "lang.csv"
	slog.Info("使用多语言文件", "langfile", langfile)
	lang.Init(langfile, ConfigManager)
	/////// /////// /////// /////// ///////

	p := discovery.NewFileStaticDiscovery()
	// p.GetAddr()
	// p.GetPort()
	// PortProvider = p
	// ServiceDiscovery = p

	RouteFile = p
}

// func Currency() (cc, cs string) {
// 	cc, _ = RouteFile.Get("Currency")
// 	cs, _ = RouteFile.Get("CurrencySymbol")

// 	if cc == "" {
// 		cc = "THB"
// 	}

// 	if cs == "" {
// 		cs = "฿"
// 	}
// 	return
// }

/*
func Currency() (cc string) {
	cc, _ = RouteFile.Get("Currency")
	// cs, _ = RouteFile.Get("CurrencySymbol")

	if cc == "" {
		cc = "THB"
	}

	// if cs == "" {
	// 	cs = "฿"
	// }
	return
}

func CurrencySymbol() (cs string) {
	// cc, _ = RouteFile.Get("Currency")
	cs, _ = RouteFile.Get("CurrencySymbol")

	// if cc == "" {
	// 	cc = "THB"
	// }

	if cs == "" {
		cs = "฿"
	}
	return
}
*/
