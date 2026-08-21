package staticproxy

import (
	"bytes"
	"cmp"
	"game/comm/mux"
	"game/comm/ut"
	"game/service/ppgateway/internal/gamedata"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

func lastGameHistory(w http.ResponseWriter, r *http.Request) {
	// https://5g6kpi7kjf.uapuqhki.net/gs2c/lastGameHistory.do?symbol=vs20olympx&tz=%2B08%3A00&mgckey=AUTHTOKEN@91c576f33153e16efbf90a9d89d37355f301ac43735173aa1564143d82bfe50c~stylename@hllgd_hollygod~SESSION@296051b8-4a10-45ee-9954-7ea4137e0a06~SN@638c2343

	// https://backoffice-sg53.ppgames.net/admin/gameHistoryDetails.do?playSessionId=39297197956091&styleName=hllgd_hollygod&memberID=48297256

	// tz: +08:00

	// ps := r.URL.Query()
	host := r.Host
	if strings.HasPrefix(host, "cdn.") {
		host = host[4:]
	}
	ul := gamedata.Get().ReverseProxyUrls.Get(host)
	if ul == nil {
		w.WriteHeader(404)
		return
	}

	file := path.Join("cache", ul.Host, r.URL.Path+".html")
	if fileinfo, err := os.Stat(file); err == nil {
		data, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data = lastGameHistory_hook(data, w, r)
		http.ServeContent(w, r, file, fileinfo.ModTime(), bytes.NewReader(data))
		return
	}

	// symbol := ps.Get("symbol")
	body, err := downloadGameHistoryHtml()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ut.Writefile(file, body)

	body = lastGameHistory_hook(body, w, r)
	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
}

var (
	tokenRegexp = regexp.MustCompile(`window\.token\s*=\s*"\w+"`)
	langRegexp  = regexp.MustCompile(`window\.language\s*=\s*"\w\w"`)
)

func downloadGameHistoryHtml() (body []byte, err error) {
	symbol := "vs20olympx"
	mgckey, loc, err := fetchMgckey(symbol, "12345678", http.DefaultClient)
	if err != nil {
		return
	}

	loc.Path = "/gs2c/lastGameHistory.do"
	q := url.Values{
		"symbol": {symbol},
		"tz":     {"+08:00"},
		"mgckey": {mgckey},
	}
	loc.RawQuery = q.Encode()

	// req, _ := http.NewRequest(http.MethodGet, loc.String(), nil)
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	body, err = ut.HttpGetBody(loc.String())
	if err != nil {
		return
	}

	body = bytes.ReplaceAll(body, []byte(loc.Host), []byte("{{HOST}}"))
	body = bytes.ReplaceAll(body, []byte("static.prerelease-env.biz"), []byte("{{HOST}}"))

	body = tokenRegexp.ReplaceAll(body, []byte(`window.token="{{TOKEN}}"`))
	body = langRegexp.ReplaceAll(body, []byte(`window.language="{{LANG}}"`))

	return
}

func fetchMgckey(game, user string, client *http.Client) (mgckey string, location *url.URL, err error) {
	var urlret struct {
		Url string
	}
	err = mux.HttpInvoke("http://127.0.0.1:55555/plat/PP/LaunchGame", map[string]any{
		"Game": game,
		"UID":  user,
		"Lang": "en",
	}, &urlret)

	if err != nil {
		return
	}

	// client req the rest
	req, err := http.NewRequest(http.MethodGet, urlret.Url, nil)

	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	// resp, err := http.DefaultTransport.RoundTrip(req)
	trans := cmp.Or(client.Transport, http.DefaultTransport)
	resp, err := trans.RoundTrip(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		return
	}

	location, err = url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return
	}

	// client req the rest
	req, err = http.NewRequest(http.MethodGet, location.String(), nil)

	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	// resp, err := http.DefaultTransport.RoundTrip(req)
	trans = cmp.Or(client.Transport, http.DefaultTransport)
	resp, err = trans.RoundTrip(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		return
	}
	location, err = url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return
	}

	query := location.Query()
	mgckey = query.Get("mgckey")
	location.RawQuery = ""

	return
}

func lastGameHistory_hook(content []byte, w http.ResponseWriter, r *http.Request) []byte {
	ps := r.URL.Query()
	content = bytes.Replace(content, []byte("</head>"), []byte(`<style>
            .Table__Cell--win{
                text-align: left!important;
            }
            .Table__Cell--roundId{
                font-size: 1rem;
            }
        </style>
    </head>`), 1)
	host := r.Host
	if strings.HasPrefix(host, "cdn.") {
		host = host[4:]
	}
	content = bytes.ReplaceAll(content, []byte("{{HOST}}"), []byte(host))

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	content = bytes.ReplaceAll(content, []byte("{{TOKEN}}"), []byte(ps.Get("mgckey")))

	lang := cmp.Or(ps.Get("lang"), "en")
	content = bytes.ReplaceAll(content, []byte("{{LANG}}"), []byte(lang))
	return content
}
