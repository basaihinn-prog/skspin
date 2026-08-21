package staticproxy

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"game/comm/mux"
	"game/comm/ut"
	"game/service/jdbgateway/internal/gamedata"
	"time"

	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
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
	//ps := r.URL.Query()
	host := r.Host
	if strings.HasPrefix(host, "cdn.") {
		host = host[4:]
	}
	preStr := "jdbstatic."
	originHost := strings.Join(strings.Split(host, ".")[1:], ".")
	content = bytes.ReplaceAll(content, []byte("{{HOST}}"), []byte("https://"+preStr+originHost))

	// content = bytes.ReplaceAll(content, []byte("{{HOST}}"), []byte("http://"+strings.Split(host, ":")[0]+":8088"))

	return content
}

func html5Game(w http.ResponseWriter, r *http.Request) {
	ps := r.URL.Query()
	ul := gamedata.Get().ReverseProxyUrls.Get(r.Host)
	if ul == nil {
		w.WriteHeader(404)
		return
	}

	file := path.Join("cache", ul.Host, r.URL.Path)
	if fileinfo, err := os.Stat(file); err == nil {
		data, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		//data = html5gameDo_hook(data, w, r)
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

func Details(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\":\"0000\",\"data\":[]}"))
}

//go:embed html/gcp1688.json
var gcp1688Json []byte

func Gcp1688(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	json := JsonConfigHook(gcp1688Json, w, r)
	byets := base64.StdEncoding.EncodeToString(json)
	w.Write([]byte("j" + byets))
}

// 都是api
func JsonConfigHook(content []byte, w http.ResponseWriter, r *http.Request) []byte {
	host := r.Host
	if strings.HasPrefix(host, "cdn.") {
		host = host[4:]
	}
	preStr := "jdbapi."
	originHost := strings.Join(strings.Split(host, ".")[1:], ".")
	content = bytes.ReplaceAll(content, []byte("{{HOST}}"), []byte(preStr+originHost))
	//content = bytes.ReplaceAll(content, []byte("{{HOST}}"), []byte(strings.Split(host, ":")[0]+":8088"))
	return content
}

//go:embed html/loggerConfig.json
var loggerConfigJson []byte

func loggerConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	json := JsonConfigHook(loggerConfigJson, w, r)
	byets := base64.StdEncoding.EncodeToString(json)
	w.Write([]byte("j" + byets))
}
func downloadhtml(ps url.Values) (body []byte, err error) {
	var urlret struct {
		Url string
	}
	err = mux.HttpInvoke("https://192.168.1.193:55555/plat/PP/LaunchGame", map[string]any{
		"Game": ps.Get("symbol"),
		"UID":  "12345678",
		"Lang": "en",
	}, &urlret)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("https://192.168.1.193:55555/plat/PP/LaunchGame err", "err", err)
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
