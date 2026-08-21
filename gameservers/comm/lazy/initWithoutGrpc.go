package lazy

import (
	"log/slog"
	"serve/comm/lang"
)

var RouteFile *FilePortProvider

func Init(serviceName string) {
	ServiceName = serviceName

	slog.Info("服务初始化", "ServiceName", ServiceName)

	fw := NewFileWatcher("config")
	GConfigManager = NewConfigManager(fw)

	langfile := "lang.csv"
	slog.Info("使用多语言文件", "langfile", langfile)
	GConfigManager.WatchAndLoad(langfile, lang.LoadLang)

	RouteFile = NewFilePortProvider()
}

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
