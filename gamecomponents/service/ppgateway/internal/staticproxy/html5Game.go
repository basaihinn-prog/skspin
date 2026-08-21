package staticproxy

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/define"
	"game/comm/mux"
	"game/comm/ut"
	"game/duck/lazy"
	"game/duck/ut2"
	"game/duck/ut2/jwtutil"
	"game/service/ppgateway/internal/gamedata"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mgckeyRegexp = regexp.MustCompile(`AUTHTOKEN@[^"]+"`)
	// currencyRegexp          = regexp.MustCompile(`"currency":"\w\w\w"`)
	// brandRequirementsRegexp = regexp.MustCompile(`"brandRequirements":"[^"]*"`)

	gameConfigRegexp = regexp.MustCompile(`gameConfig:\s*'([^']+)'`)
	cashierUrlRegexp = regexp.MustCompile(`cashierUrl:\s*"([^"]+)"`)
	lobbyUrlRegexp   = regexp.MustCompile(`lobbyUrl:\s*"([^"]+)"`)
)

func html5gameDo_hook(content []byte, w http.ResponseWriter, r *http.Request) []byte {
	ps := r.URL.Query()
	host := r.Host
	if strings.HasPrefix(host, "cdn.") {
		host = host[4:]
	}
	content = bytes.ReplaceAll(content, []byte("{{HOST}}"), []byte(host))
	content = bytes.ReplaceAll(content, []byte("{{MGCKEY}}"), []byte(ps.Get("mgckey")))

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	pid, err := jwtutil.ParseToken(ps.Get("mgckey"))
	if err != nil {
		err = define.NewErrCode("Invalid player session", 1302)
		return nil
	}
	{
		// adjust gameconfig
		submatch := gameConfigRegexp.FindSubmatchIndex(content)

		i, j := submatch[2], submatch[3]
		jsondata := content[i:j]
		var gameconfig map[string]json.RawMessage
		json.Unmarshal(jsondata, &gameconfig)
		var doc struct {
			AppID string `bson:"AppID"` // 所属产品
		}
		// var appid string
		coll := db.Collection2("game", "Players")
		coll.FindOne(context.TODO(), db.ID(pid), options.FindOne().SetProjection(db.D("AppID", 1))).Decode(&doc)
		gameconfig["currency"] = ut.GetJsonRawMust("")
		gameconfig["currencyOriginal"] = ut.GetJsonRawMust("")
		operatorInfo := comm.GetOperatorInfo(doc.AppID)
		if operatorInfo.CurrencyManufactureVisibleOff != nil {
			if _, ok := operatorInfo.CurrencyManufactureVisibleOff[comm.PP]; ok && operatorInfo.CurrencyManufactureVisibleOff[comm.PP] == 1 {
				item := lazy.GetCurrencyItem(operatorInfo.CurrencyKey)
				gameconfig["currency"] = ut.GetJsonRawMust(item.Symbol)
				gameconfig["currencyOriginal"] = ut.GetJsonRawMust(item.Key)
			}
		}
		lang := cmp.Or(ps.Get("lang"), "en")
		gameconfig["lang"] = ut.GetJsonRawMust(lang)
		gameconfig["sessionTimeout"] = ut.GetJsonRawMust("300")
		gameconfig["styleName"] = ut.GetJsonRawMust("abcde_fekfdksk")
		//gameconfig["brandRequirements"] = ut.GetJsonRawMust("NORP")
		//gameconfig["brandRequirements"] = ut.GetJsonRawMust("BBB")
		gameID := "pp_" + ps.Get("symbol")
		params := comm.GetEXParams(pid, gameID)
		if params.ShowNameAndTimeOff == 1 {
			gameconfig["jurisdictionRequirements"] = ut.GetJsonRawMust("GT,CLK,PORC,NOSTW")
		}
		gameconfig["replaySystemUrl"] = ut.GetJsonRawMust(fmt.Sprintf("https://replay.%s", ut2.Domain(host)))
		gameconfig["lobbyLaunched"] = ut.GetJsonRawMust(false)
		delete(gameconfig, "lobbyVersion")

		//gameconfig["openHistoryInWindow"] = ut.GetJsonRawMust("true")
		//gameconfig["openHistoryInTab"] = ut.GetJsonRawMust("true")

		// "HISTORY":"https://ppgames.rpgamestest.com/gs2c/lastGameHistory.do?symbol\u003dvs20olympx\u0026mgckey\u003deyJQIjoxMDAyMjUsIkUiOjE3Mjc0NDYxOTUsIlMiOjEwMDIsIkQiOiJwcF92czIwb2x5bXB4In0.Sxgf-S8be3Q5WwyMQH74COxrQNhqPugquxVeRxPc5zU"

		var history string
		json.Unmarshal(gameconfig["HISTORY"], &history)
		hisU, _ := url.Parse(history)
		if hisU != nil {
			q := hisU.Query()
			q.Add("lang", lang)
			hisU.RawQuery = q.Encode()

			gameconfig["HISTORY"] = ut.GetJsonRawMust(hisU.String())
		}

		outdata, _ := json.Marshal(gameconfig)
		content = slices.Replace(content, i, j, outdata...)
	}
	{
		// adjust cashierUrl
		submatch := cashierUrlRegexp.FindSubmatchIndex(content)
		if submatch != nil {
			i, j := submatch[2], submatch[3]
			content = slices.Replace(content, i, j, []byte(``)...)
		}

	}
	{ // adjust lobbyUrl
		submatch := lobbyUrlRegexp.FindSubmatchIndex(content)
		if submatch != nil {
			i, j := submatch[2], submatch[3]
			content = slices.Replace(content, i, j, []byte(``)...)
		}

	}

	// outdata, _ := json.Marshal(jsonmap)
	// os.Stdout.Write(outdata)

	return content
}

func html5Game(w http.ResponseWriter, r *http.Request) {
	// https://ppgames.rpgamestest.com/gs2c/html5Game.do?cb_target=exist_tab&ext=0&extGame=1&gname=Gates+of+Olympus+1000&jackpotid=0&jurisdictionID=99&lang=en&mgckey=eyJQIjoxMDAwMDEsIkUiOjE3MjQ4NTYwOTgsIlMiOjEwMDUsIkQiOiJwcF92czIwb2x5bXB4In0.NIyYV3wIdaUeuGzm7NQg0rQxGpihJBAiqZ4WQD7BttY&minilobby=false&symbol=vs20olympx&tabName=

	// {"UID":"123456","Game":"vs20olympx","Lang":"en"}
	ps := r.URL.Query()
	ul := gamedata.Get().ReverseProxyUrls.Get(r.Host)
	if ul == nil {
		w.WriteHeader(404)
		return
	}

	file := path.Join("cache", ul.Host, r.URL.Path+"-"+ps.Get("symbol")+".html")
	if fileinfo, err := os.Stat(file); err == nil {
		data, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data = html5gameDo_hook(data, w, r)
		http.ServeContent(w, r, file, fileinfo.ModTime(), bytes.NewReader(data))
		return
	}

	// symbol := ps.Get("symbol")
	body, err := downloadhtml(ps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ut.Writefile(file, body)

	body = html5gameDo_hook(body, w, r)
	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
}

func downloadhtml(ps url.Values) (body []byte, err error) {
	var urlret struct {
		Url string
	}
	err = mux.HttpInvoke("http://127.0.0.1:55555/plat/PP/LaunchGame", map[string]any{
		"Game": ps.Get("symbol"),
		"UID":  "12345678",
		"Lang": "en",
	}, &urlret)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("http://127.0.0.1:55555/plat/PP/LaunchGame err", "err", err)
		return
	}

	// urlret.Url = "https://plats.rpgamestest.com/api/"

	req, err := http.NewRequest(http.MethodGet, urlret.Url, nil)

	if err != nil {
		slog.Error(urlret.Url+" err", "err", err)
		return
	}

	// time.Sleep(time.Second)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		slog.Error("http.DefaultTransport.RoundTrip err", "err", err)
		return
	}

	io.Copy(os.Stdout, resp.Body)

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		return
	}

	location := resp.Header.Get("Location")

	req2, err := http.NewRequest(http.MethodGet, location, nil)
	if err != nil {
		slog.Error(location+" err", "err", err)
		return
	}

	// time.Sleep(time.Second)

	req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	response, err := http.DefaultTransport.RoundTrip(req2)
	if err != nil {
		slog.Error("http.DefaultTransport.RoundTrip err", "err", err)
		return
	}

	defer response.Body.Close()
	location2 := response.Header.Get("Location")

	locU, err := url.Parse(location2)
	if err != nil {
		slog.Error(location2+" err", "err", err)
		return
	}

	query := locU.Query()
	locU.RawQuery = query.Encode()
	body, err = ut.HttpGetBody(locU.String())
	if err != nil {
		return
	}

	body = bytes.ReplaceAll(body, []byte(locU.Host), []byte("{{HOST}}"))
	body = bytes.ReplaceAll(body, []byte("common-static.prerelease-env.biz"), []byte("{{HOST}}"))

	body = mgckeyRegexp.ReplaceAll(body, []byte("{{MGCKEY}}\""))

	body = bytes.Replace(body, []byte("<head>"), []byte(`<head>
<script>
window.rawlog=console.log;
</script>`), 1)

	// body = currencyRegexp.ReplaceAll(body, []byte(`"currency":"{{CURRENCY}}"`))

	// body = bytes.ReplaceAll(body, []byte(`"lang":"en"`), []byte(`"lang":"{{LANG}}"`))

	return
}
