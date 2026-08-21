package staticproxy

import (
	"game/service/jiligateway/internal/gamedata"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"path"
)

func newReverseProxy() *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			// ul := matcher[r.In.Host]
			ul := gamedata.Get().ReverseProxyUrls.Get(r.In.Host)
			r.SetURL(ul)
			// r.Out.Header.Set("Refer", "https://m.pgsoft-games.com/")
			// r.Out.Header.Set("Origin", "https://m.pgsoft-games.com")
		},
		ModifyResponse: func(resp *http.Response) (err error) {

			file := path.Join(resp.Request.URL.Host, resp.Request.URL.RequestURI())
			slog.Info("save response",
				"file", file,
				"method", resp.Request.Method,
				"status", resp.Status,
				"content-type", resp.Header.Get("Content-Type"),
				"content-length", resp.ContentLength,
			)
			return nil
		},
	}
	return proxy
}
