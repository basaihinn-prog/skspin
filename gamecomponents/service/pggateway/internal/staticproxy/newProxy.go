package staticproxy

/*
import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"time"

	"github.com/samber/lo"
)

func newProxy() *httputil.ReverseProxy {
	// rpURL := lo.Must(url.Parse("https://m.pgsoft-games.com/"))
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			ul := matcher[r.In.Host]
			r.SetURL(ul)
			r.Out.Header.Set("Refer", "https://m.pgsoft-games.com/")
			r.Out.Header.Set("Origin", "https://m.pgsoft-games.com")
		},
		ModifyResponse: func(resp *http.Response) (err error) {
			file := path.Join("cache", resp.Request.URL.Host, resp.Request.URL.RequestURI())

			slog.Info("save response",
				"file", file,
				"method", resp.Request.Method,
				"status", resp.Status,
				"content-type", resp.Header.Get("Content-Type"),
				"content-length", resp.ContentLength,
			)

			if resp.Request.Method != "GET" || resp.StatusCode != http.StatusOK || resp.ContentLength == 0 {
				return
			}
			// querys := resp.Request.URL.Query()
			resp.Header.Set("Access-Control-Allow-Origin", "*")

			resp.Header.Set("Cache-Control", "max-age=315360000")
			resp.Header.Set("Expires", time.Now().Add(time.Second*315360000).UTC().Format(time.RFC1123))

			res, _ := httputil.DumpResponse(resp, true)
			writefile(file, res)

			if false {
				body := lo.Must(io.ReadAll(resp.Body))
				resp.Body.Close()
				resp.Body = io.NopCloser(bytes.NewReader(body))

				// zr, err := gzip.NewReader(req.Body)
				if resp.Header.Get("Content-Encoding") == "gzip" {
					zr := lo.Must(gzip.NewReader(bytes.NewReader(body)))
					body = lo.Must(io.ReadAll(zr))
					zr.Close()
				}
				writefile(file, body)
			}

			return
		},
	}
	return proxy
}


*/
