package staticproxy

import (
	"bytes"
	"fmt"
	"game/comm/ut"
	"game/service/minigateway/internal/gamedata"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/samber/lo"
)

var (
	mgckeyRegexp = regexp.MustCompile(`AUTHTOKEN@[^"]+"`)
)

func html5gameDo_hook(content []byte, w http.ResponseWriter, r *http.Request) []byte {
	//ps := r.URL.Query()
	host := r.Host
	//content = bytes.ReplaceAll(content, []byte("{{HOST}}"), []byte(host))
	//content = bytes.ReplaceAll(content, []byte("{{MGCKEY}}"), []byte(ps.Get("mgckey")))
	content = bytes.ReplaceAll(content, []byte("{{app-config.spribe.dev}}"), []byte(`https://`+host))
	//content = bytes.ReplaceAll(content, []byte("https://{{app-config.spribe.dev}}"), []byte(host))

	return content
}

func staticAssets(w http.ResponseWriter, r *http.Request) {
	var err error
	ps := r.URL.Query()
	ul := gamedata.Get().ReverseProxyUrls.Get(r.Host)
	if ul == nil {
		w.WriteHeader(404)
		return
	}
	file := path.Join("cache", ul.Host, r.URL.Path)
	if r.URL.Path == "/" {
		file += "/index.html"
	}
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
	//没有文件 需要下载
	outuri := ul.JoinPath(r.URL.Path).String()
	body, err := downloadhtml(ps, outuri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	ut.Writefile(file, body)
	body = html5gameDo_hook(body, w, r)

	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
}

var contenthooks = []contenthook{}

func downloadhtml(ps url.Values, outuri string) (body []byte, err error) {
	lg := slog.With()
	lg = lg.With("outuri", outuri)
	req := lo.Must(http.NewRequest(http.MethodGet, outuri, nil))
	//req.Header.Set("Refer", "https://m.pgsoft-games.com/")
	//req.Header.Set("Origin", "https://m.pgsoft-games.com")
	//根据文件类型设置文件处理方式
	if strings.Contains(outuri, "main.d2d985e22e14c8ac.js") {
		contenthooks = append(contenthooks, buildjs_replace_appconfigSpribeDev)
		contenthooks = append(contenthooks, buildjs_JavaScript2Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("downloadhtml::http.DefaultClient.Do Err", "err", err.Error(), "code", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("downloadhtml::io.ReadAll Err", "err", err.Error(), "code", http.StatusInternalServerError)
		return
	}
	lg = lg.With("Status", resp.Status)
	if resp.StatusCode == http.StatusForbidden {
		//todo 去其他平台获取资源
		slog.Info("去其他平台获取资源")
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("downloadhtml::resp.StatusCode != http.StatusOK")
		err = fmt.Errorf("downloadhtml::resp.StatusCode == %d", resp.StatusCode)
		return
	}
	//处理下载的源文件
	for _, hook := range contenthooks {
		body = hook(body)
	}
	return
}

func huidustg(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	addr := gamedata.Get().WsAddr
	out := fmt.Sprintf(`{
		"brandLogo": "",
		"brandName": "huidustg",
		"ws": {
		"host": "%v",
			"port": 443,
			"useSSL": true,
			"zone": "aviator-inst1",
			"debug": false
	},
		"userWithOperator": true
	}`, addr)
	w.Write([]byte(out))
}

func buildjs_replace_appconfigSpribeDev(content []byte) []byte {
	return bytes.ReplaceAll(content, []byte(`https://app-config.spribe.dev`), []byte(`{{app-config.spribe.dev}}`))
}

func buildjs_JavaScript2Token(content []byte) []byte {
	return bytes.ReplaceAll(content, []byte(`V._clientDetails="JavaScript"`), []byte(`V._clientDetails=window["session"]["token"]`))
}
