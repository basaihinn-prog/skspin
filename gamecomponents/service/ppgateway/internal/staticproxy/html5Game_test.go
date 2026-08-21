package staticproxy

import (
	"encoding/json"
	"game/comm/ut"
	"os"
	"slices"
	"testing"
)

func TestDownload(t *testing.T) {
	// downloadhtml("vs20olympx")
	// u := "http://5g6kpi7kjf.uapuqhki.net/gs2c/openGame.do?tc=GV1TkfeL9olIn5AqKRcT71mqucarwRSaaz1T3xLZXUDc1G0CEDABncVJ54352F71&stylename=hllgd_hollygod&dummy="

	// http.DefaultTransport.RoundTrip()
	// resp, err := http.Get(u)
	// fmt.Println(resp, err)

	data, _ := os.ReadFile("/data/game/bin/cache/5g6kpi7kjf.uapuqhki.net/gs2c/html5Game.do-vs20olympx.html")

	submatch := gameConfigRegexp.FindSubmatchIndex(data)

	i, j := submatch[2], submatch[3]
	jsondata := data[i:j]
	var gameconfig map[string]json.RawMessage
	json.Unmarshal(jsondata, &gameconfig)
	gameconfig["currency"] = ut.GetJsonRawMust("USD")
	gameconfig["currencyOriginal"] = ut.GetJsonRawMust("USD")
	gameconfig["lang"] = ut.GetJsonRawMust("th")
	gameconfig["sessionTimeout"] = ut.GetJsonRawMust("300")
	gameconfig["styleName"] = ut.GetJsonRawMust("abcde_fekfdksk")
	gameconfig["brandRequirements"] = ut.GetJsonRawMust("NORP,NOGA")

	outdata, _ := json.Marshal(gameconfig)
	os.Stdout.Write(outdata)

	data = slices.Replace(data, i, j, outdata...)

	// bytes.Replace()

	// gameConfigRegexp.ReplaceAll()

	os.Stdout.Write(data)
	// fmt.Println(gameConfigRegexp.Match(data))
}

func TestGameHistoryHtml(t *testing.T) {
	body, _ := downloadGameHistoryHtml()
	os.Stdout.Write(body)

	// tokenRegexp
}
