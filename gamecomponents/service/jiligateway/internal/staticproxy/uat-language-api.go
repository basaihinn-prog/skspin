package staticproxy

import (
	"bytes"
	"game/comm/ut"
	"game/service/jiligateway/internal/gamedata"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/samber/lo"
)

func uat_language_api(w http.ResponseWriter, r *http.Request) {
	storeul := gamedata.Get().ReverseProxyUrls.Get(r.Host)
	if storeul == nil {
		w.WriteHeader(404)
		return
	}

	ul := *storeul

	r.URL.RawQuery = ""

	var err error
	file := path.Join("cache", ul.Host, r.URL.RequestURI())
	// file := "cache/" + ul.Host + r.URL.Path
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
			// "Header", w.Header(),
			"error", err,
		)
		lg.Info("req resources")
	}()

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	if is_index_html(file) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
	}

	// only test used
	if strings.HasSuffix(file, ".js") {
		w.Header().Set("Cache-Control", "no-cache")
	}

	if _, err := os.Stat(file); err == nil {
		http.ServeFile(w, r, file)
		// if strings.HasPrefix(r.Host, "uat-language-api.") && r.URL.Path == "/language/en-US" {
		// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// }
		return
	}

	// outuri := ul.JoinPath(r.URL.Path).String()

	ul.Path = r.URL.Path
	ul.RawQuery = r.URL.RawQuery

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

	lg = lg.With("from", "REMOTE")

	if !cache_nostore(resp.Header) {
		// if file[len(file)-1] == '/' {
		// 	file += "index.html"
		// }
		ut.Writefile(file, body)
	}
	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
}
