package staticproxy

import (
	"bytes"
	"game/comm/ut"
	"game/service/facaigateway/internal/gamedata"

	"io"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
)

// excute before save
type contenthook = func([]byte) []byte

// excute before send
// type respcontenthook = func([]byte, http.ResponseWriter, *http.Request) []byte
type respcontenthook = func([]byte, http.ResponseWriter, *http.Request) []byte

func static_assets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ul := gamedata.Get().ReverseProxyUrls.Get(r.Host)
	if ul == nil {
		w.WriteHeader(404)
		return
	}

	ul.Path = r.URL.Path
	ul.RawQuery = r.URL.RawQuery

	var (
		file = "cache/" + ul.Host + r.URL.Path
	)

	var contenthooks = []contenthook{}
	var respcontenthooks = []respcontenthook{}

	pth := r.URL.Path
	switch {
	case strings.HasSuffix(pth, ".js"):
		//contenthooks = append(contenthooks, buildjs_nop_host_valid)
	//contenthooks = append(contenthooks, buildjs_nop_host_valid2)
	case pth == "/index.html":
		ul.Path = "/index"
	}
	//respcontenthooks = append(respcontenthooks, html5gameDo_hook)

	var err error
	if file[len(file)-1] == '/' {
		file += "index.html"
	}

	inuri := r.URL.RequestURI()
	lg := slog.With(
		"file", file,
		"inuri", inuri,
		"method", r.Method,
	)

	defer func() {
		lg = lg.With(
			"error", err,
		)
		lg.Info("req resources")
	}()

	if w.Header().Get("Cache-Control") == "" {
		w.Header().Set("Cache-Control", "public, max-age=2592000, must-revalidate, proxy-revalidate")
	}

	if fileinfo, err := os.Stat(file); err == nil {
		// http.ServeFile(w, r, file)
		data, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		for _, hook := range contenthooks {
			data = hook(data)
		}
		for _, hook := range respcontenthooks {
			data = hook(data, w, r)
			//data = hook(data)
		}
		http.ServeContent(w, r, file, fileinfo.ModTime(), bytes.NewReader(data))
		return
	}

	outuri := ul.String()
	lg = lg.With("outuri", outuri)

	req := lo.Must(http.NewRequest(http.MethodGet, outuri, nil))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		// w.WriteHeader(http.Status)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// w.WriteHeader(http.Status)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lg = lg.With("Status", resp.Status)

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		return
	}
	body = bytes.ReplaceAll(body, []byte("https://dluqiiiaw.cnwzhy.com"), []byte("{{HOST}}"))
	////body = bytes.ReplaceAll(body, []byte("https://dlbase0pla.ebaenycho.net"), []byte("{{HOST2}}"))
	////body = bytes.ReplaceAll(body, []byte("https://locdev.jdb188.net"), []byte("{{HOST3}}"))
	////body = bytes.ReplaceAll(body, []byte("https://loc.jdb188.net"), []byte("{{HOST4}}"))
	body = bytes.ReplaceAll(body, []byte("https://dlbase0pla.ebaenycho.net"), []byte("{{HOST2}}")) //json1688        paytable图片
	//body = bytes.ReplaceAll(body, []byte("https://locdev.jdb188.net"), []byte("{{HOST}}"))
	//body = bytes.ReplaceAll(body, []byte("https://loc.jdb188.net"), []byte("{{HOST}}"))
	//body = bytes.ReplaceAll(body, []byte("player.jygrq.com"), []byte("{{HOST2}}")) //历史相关
	lg = lg.With("from", "REMOTE")

	if !cache_nostore(resp.Header) {
		ut.Writefile(file, body)
	}
	for _, hook := range contenthooks {
		body = hook(body)
	}

	for _, hook := range respcontenthooks {
		body = hook(body, w, r)
	}
	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
}

func cache_nostore(header http.Header) bool {
	cache_control := header.Values("Cache-Control")
	nostore := slices.ContainsFunc(cache_control, func(s string) bool {
		return strings.Contains(s, "no-store")
	})
	return nostore
}
