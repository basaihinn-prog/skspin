package staticproxy

import (
	"bytes"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/samber/lo"
)

var (
	MyHost   = "ppproxy.rpgamestest.com"
	Upstream = lo.Must(url.Parse("https://5g6kpi7kjf.uapuqhki.net"))
)

func newReverseProxy() *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(Upstream)
			r.Out.Header.Del("Accept-Encoding")
		},
		ModifyResponse: func(resp *http.Response) (err error) {
			pth := resp.Request.URL.Path
			slog.Info("ModifyResponse", "path", pth)

			var contenthooks = []contenthook{}
			switch {
			case pth == "/gs2c/openGame.do" && resp.StatusCode == 302:
				location := resp.Header.Get("Location")
				lurl := lo.Must(url.Parse(location))
				lurl.Host = MyHost
				resp.Header.Set("Location", lurl.String())
			case strings.HasSuffix(pth, "/build.js"):
				contenthooks = append(contenthooks, buildjs_nop_host_valid)
				contenthooks = append(contenthooks, buildjs_gtag_hook)
				contenthooks = append(contenthooks, buildjs_post)
			case pth == "/gs2c/html5Game.do" && resp.StatusCode == 200:

				contenthooks = append(contenthooks, func(body []byte) []byte {
					body = bytes.ReplaceAll(body, []byte(Upstream.Host), []byte(MyHost))

					body = bytes.Replace(body, []byte("<head>"), []byte(`<head>
<script>
window.rawlog=console.log;
</script>`), 1)
					return body
				})
				// respBody, _ := io.ReadAll(resp.Body)
				// zr, _ := gzip.NewReader(bytes.NewReader(respBody))
				// content, _ := io.ReadAll(zr)
				// content = bytes.ReplaceAll(content, []byte("763a46cfc8.uympierc.net"), []byte(MyHost))

				// payload := &bytes.Buffer{}
				// zw := gzip.NewWriter(payload)
				// zw.Write(content)
				// zw.Close()

				// resp.Body = io.NopCloser(bytes.NewReader(payload.Bytes()))
			}

			if len(contenthooks) != 0 {
				// "/gs2c/html5Game.do?jackpotid=0&gname=Gates%20of%20Olympus%201000&extGame=1&ext=0&cb_target=exist_tab&symbol=vs20olympx&jurisdictionID=99&minilobby=false&mgckey=AUTHTOKEN@828d7a5e47f344f29a669379e709440e65c6ed4fdba0f88cf6a7080c348b8e95~stylename@hllgd_hollygod~SESSION@8ee5ac5d-9d40-4539-8007-04b43a8d4caa~SN@98e44375&tabName="

				// reqU := resp.Request.URL
				// query := reqU.Query()
				// maps.DeleteFunc(func(k string, v []string) bool{
				// 	return k != "symbol" && k != ""
				// })
				// file := path.Join(reqU.Hostname())

				body, _ := io.ReadAll(resp.Body)
				defer resp.Body.Close()

				for _, hook := range contenthooks {
					body = hook(body)
				}

				resp.Body = io.NopCloser(bytes.NewReader(body))
			}

			return nil
		},
	}
	return proxy
}

func StartReverseProxy(addr string) {
	go func() {
		httpmux := http.NewServeMux()
		// /gs2c/ge/v3/gameService
		// httpmux.HandleFunc("/", proxyhandle)
		httpmux.Handle("/", newReverseProxy())

		//证书
		err := http.ListenAndServe(addr, httpmux)
		log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}
