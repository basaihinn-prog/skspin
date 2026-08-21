package staticproxy

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"game/comm/ut"
	"game/service/pggateway/internal/gamedata"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/samber/lo"
)

var (
	historyHtmlExp = regexp.MustCompile(`^/history/\d+.html$`)
	indexRegexp    = regexp.MustCompile(`^/\d+/index.html$`)
	// https://static-pg.rpgamestest.com/shared/3c4695a542/index-100.json
	sharedIndexExp = regexp.MustCompile(`^/shared/\w+/index-\d+\.(json|js)$`)
)

type contenthook = func([]byte) []byte

func static_assets(w http.ResponseWriter, r *http.Request) {
	storeul := gamedata.Get().ReverseProxyUrls.Get(r.Host)
	if storeul == nil {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	ul := *storeul
	ul.Path = r.URL.Path
	ul.RawQuery = ""

	var (
		file             = path.Join("cache", ul.Host, r.URL.Path)
		err              error
		contenthooks     = []contenthook{}
		respcontenthooks = []contenthook{}
		pth              = r.URL.Path
	)

	// inuri := r.URL.RequestURI()
	lg := slog.With(
		"file", file,
		"inuri", pth,
		"method", r.Method,
	)

	defer func() {
		lg.Info("static_assets!", "error", ut.ErrString(err))
	}()

	switch {
	case historyHtmlExp.MatchString(pth):
		w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")
	case indexRegexp.MatchString(pth):
		contenthooks = append(contenthooks, hookhtml)
		contenthooks = append(contenthooks, hookfetch)
		respcontenthooks = append(respcontenthooks, hookSkipInsecurePmt)
		// respcontenthooks = append(respcontenthooks, hookNopSW)
		w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")
	case sharedIndexExp.MatchString(pth):
		// /shared/3c4695a542/index-100.js
		data := []byte(pth)
		l := bytes.LastIndexByte(data, '-')
		r := bytes.LastIndexByte(data, '.')

		data = slices.Delete(data, l, r)
		ul.Path = string(data)
		w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")
	case strings.HasSuffix(pth, ".js.map"):
		w.WriteHeader(404)
		return
	case strings.HasSuffix(pth, "sw.js"):
		w.Header().Set("Service-Worker-Allowed", "/")
		w.Header().Add("Content-Type", "application/javascript")
		w.Header().Add("Content-Type", "charset=UTF-8")
	case strings.HasSuffix(file, "/index.js") ||
		strings.HasSuffix(file, "/index.json") ||
		strings.HasSuffix(file, "/index.html") ||
		false:
		// w.Header().Set("Cache-Control", "public, max-age=120, s-maxage=604800")
		w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")

	default:
		w.Header().Set("Cache-Control", "public, max-age=86400")
	}

	if data, err := os.ReadFile(file); err == nil {
		var modtime time.Time
		if d, e := os.Stat(file); e == nil {
			modtime = d.ModTime()
		}
		for _, hook := range respcontenthooks {
			data = hook(data)
		}
		http.ServeContent(w, r, file, modtime, bytes.NewReader(data))
		lg = lg.With("from", "cache")
		return
	}

	outuri := ul.String()
	lg = lg.With("outuri", outuri)

	req := lo.Must(http.NewRequest(http.MethodGet, outuri, nil))
	// req.Header.Set("Refer", "https://m.pgsoft-games.com/")
	if strings.HasSuffix(file, ".json") {
		req.Header.Set("accept-encoding", "gzip, deflate, br")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// w.WriteHeader(http.Status)
		w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	lg = lg.With("Status", resp.Status)
	if resp.StatusCode != http.StatusOK {
		w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")
		w.WriteHeader(resp.StatusCode)
		return
	}

	var body []byte

	contentEncoding := resp.Header.Get("Content-Encoding")
	if contentEncoding == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer reader.Close()
		body, err = io.ReadAll(reader)
		if err != nil {
			w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if contentEncoding == "deflate" {
		reader := flate.NewReader(resp.Body)
		defer reader.Close()
		body, err = io.ReadAll(reader)
		if err != nil {
			w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if contentEncoding == "br" {
		reader := brotli.NewReader(resp.Body)
		body, err = io.ReadAll(reader)
		if err != nil {
			w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			w.Header().Set("Cache-Control", "public, max-age=600, s-maxage=600")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	lg = lg.With("from", "REMOTE")

	for _, hook := range contenthooks {
		body = hook(body)
	}

	ut.Writefile(file, body)

	for _, hook := range respcontenthooks {
		body = hook(body)
	}
	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
}

func StartProxy(addr string) {

	go func() {
		err := http.ListenAndServe(addr, http.HandlerFunc(static_assets))
		log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}
