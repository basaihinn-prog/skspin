package staticproxy

import (
	"bytes"
	"game/duck/ut2"
	"game/service/pggateway/internal/gamedata"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

// 待优化 func startwith static ....
// var matcher = map[string]*url.URL{
// 	"static-pg.kafa010.com": lo.Must(url.Parse("https://static.pg-demo.com")),
// 	"m-pg.kafa010.com":      lo.Must(url.Parse("https://m.pgsoft-games.com")),
// 	"public-pg.kafa010.com": lo.Must(url.Parse("https://public.pgsoft-games.com")),
// }

var (
	indexRegexp = regexp.MustCompile(`/\d+/index.html`)
)

func StartProxy(addr string) {
	// proxy := newProxy()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ul := gamedata.Get().ReverseProxyUrls[r.Host]
		if ul == nil {
			w.WriteHeader(404)
			return
		}
		r.URL.RawQuery = ""

		var err error
		file := path.Join("cache", ul.Host, r.URL.RequestURI())

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
			// 2024/04/16 08:32:11 INFO req resources file=cache/m.pgsoft-games.com/98/index.html inuri=/98/index.html method=GET from=cache error=<nil>
		}()

		if strings.HasSuffix(file, ".js.map") {
			w.WriteHeader(404)
			return
		}

		if true {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if strings.HasSuffix(file, "sw.js") {
				w.Header().Set("Service-Worker-Allowed", "/")
				w.Header().Add("Content-Type", "application/javascript")
				w.Header().Add("Content-Type", "charset=UTF-8")
			}

			if strings.HasSuffix(file, "/index.js") ||
				strings.HasSuffix(file, "/index.json") {
				w.Header().Set("Cache-Control", "public, max-age=28800, must-revalidate, proxy-revalidate")
			} else {
				// max-age=31536000
				w.Header().Set("Cache-Control", "public, max-age=86400, must-revalidate, proxy-revalidate")
			}

			if data, err := os.ReadFile(file); err == nil {
				var modtime time.Time
				if d, e := os.Stat(file); e == nil {
					modtime = d.ModTime()
				}

				if strings.Contains(inuri, "03fc0777ec") && strings.HasSuffix(inuri, ".js") {
					mainHost := ut2.Domain(r.Host)
					if host := gamedata.Get().ChangeVerifyMainHost; host != "" {
						mainHost = host
					}
					//todo 替换pg游戏前端资源后，需要修改这里
					data = bytes.Replace(data, []byte(`this[qE(0x688)]`), []byte(`'https://verify.`+mainHost+`'`), 1)
					data = bytes.Replace(data, []byte(`this[mi(0x71c)]`), []byte(`'https://verify.`+mainHost+`'`), 1)
				}
				if indexRegexp.MatchString(inuri) {
					mainHost := ut2.Domain(r.Host)
					data = bytes.Replace(data, []byte(`<meta name="apple-mobile-web-app-capable" content="yes">`), []byte(`<meta name="mobile-web-app-capable" content="yes"><link rel="preconnect" href="https://static.`+mainHost+`"><link rel="preconnect" href="https://api.`+mainHost+`">`), 1)
					str := `<script>appClhild_back=document.head.appendChild;document.head.appendChild=function(a){if(a&&a.src&&a.src.indexOf("www.googletagmanager.com")!=-1){return null;};return appClhild_back.apply(this,arguments);};</script>`
					// 查找字符 'o' 的位置
					index := bytes.Index(data, []byte(`<script id="main-script"`))
					if index != -1 {
						data = append(data[:index], append([]byte(str), data[index:]...)...)
					} else {
						slog.Info("未找到对应script标签")
					}
				}
				http.ServeContent(w, r, file, modtime, bytes.NewReader(data))
				lg = lg.With("from", "cache")

				// http.ServeFile()
				return
			}

			outuri := ul.JoinPath(r.URL.Path).String()
			lg = lg.With("outuri", outuri)
			req := lo.Must(http.NewRequest(http.MethodGet, outuri, nil))
			req.Header.Set("Refer", "https://m.pgsoft-games.com/")
			req.Header.Set("Origin", "https://m.pgsoft-games.com")

			resp, err := http.DefaultClient.Do(req)
			// resp, err := http.Get(ul.JoinPath(r.URL.Path).String())
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

			if resp.StatusCode == http.StatusForbidden {
				//todo 去其他平台获取资源

			}

			if resp.StatusCode != http.StatusOK {
				w.WriteHeader(resp.StatusCode)
				return
			}

			lg = lg.With("from", "REMOTE")
			if indexRegexp.MatchString(inuri) {
				body = bytes.Replace(body, []byte("</script>"), []byte(`</script><script>window.gtag=console.log;window.dataLayer=[];</script>`), 1)
			}
			writefile(file, body)
			http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
		}
	})

	go func() {
		////证书
		//certfile, _ := lazy.RouteFile.Get("certfile")
		//keyfile, _ := lazy.RouteFile.Get("keyfile")
		//err := http.ListenAndServeTLS(addr, certfile, keyfile, serveMux)
		err := http.ListenAndServe(addr, serveMux)
		log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}

func StartProxyVerify(addr string) {
	// proxy := newProxy()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ul := gamedata.Get().ReverseProxyUrls[r.Host]
		if ul == nil {
			w.WriteHeader(404)
			return
		}
		r.URL.RawQuery = ""

		var err error
		file := path.Join("cache", ul.Host, r.URL.RequestURI())

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
			// 2024/04/16 08:32:11 INFO req resources file=cache/m.pgsoft-games.com/98/index.html inuri=/98/index.html method=GET from=cache error=<nil>
		}()

		if strings.HasSuffix(file, ".js.map") {
			w.WriteHeader(404)
			return
		}

		if true {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if strings.HasSuffix(file, "sw.js") {
				w.Header().Set("Service-Worker-Allowed", "/")
				w.Header().Add("Content-Type", "application/javascript")
				w.Header().Add("Content-Type", "charset=UTF-8")
			}

			if strings.HasSuffix(file, "/index.js") ||
				strings.HasSuffix(file, "/index.json") ||
				strings.HasSuffix(file, "/index.html") ||
				false {
				// w.Header().Set("Cache-Control", "public, max-age=120, s-maxage=604800")
			} else {
				// max-age=31536000
				w.Header().Set("Cache-Control", "max-age=31536000")
			}

			if data, err := os.ReadFile(file); err == nil {
				var modtime time.Time
				if d, e := os.Stat(file); e == nil {
					modtime = d.ModTime()
				}
				if strings.Contains(inuri, "index-CS1eGBey.js") {
					data = bytes.Replace(data, []byte(`{link}`), []byte(`<span style='color:#ffffff'>https://verify.pgsoft.com</span>`), -1)
					// 匹配 "=数字;" 或 "= 数字;" 的正则表达式，支持有空格或没有空格
					regex := regexp.MustCompile(`let\s+pe\s*=\s*(\d+);`)

					// 使用 ReplaceAllFunc 对字节切片进行替换
					data = regex.ReplaceAllFunc(data, func(match []byte) []byte {
						// 提取数字部分并转换为整数
						numStr := regex.FindStringSubmatch(string(match))[1]
						num, _ := strconv.Atoi(numStr)

						// 如果数字大于2，替换为 "= 2;"，否则保留原值
						if num > 2 {
							return []byte("let pe = 2;")
						}
						return match
					})
				}
				if indexRegexp.MatchString(inuri) {
					data = bytes.Replace(data, []byte(`<meta name="apple-mobile-web-app-capable" content="yes">`), []byte(`<meta name="mobile-web-app-capable" content="yes"><link rel="preconnect" href="https://static.pgsofts-games.com"><link rel="preconnect" href="https://api.pgsofts-games.com">`), 1)
				}
				http.ServeContent(w, r, file, modtime, bytes.NewReader(data))
				lg = lg.With("from", "cache")

				// http.ServeFile()
				return
			}

			outuri := ul.JoinPath(r.URL.Path).String()
			lg = lg.With("outuri", outuri)
			req := lo.Must(http.NewRequest(http.MethodGet, outuri, nil))
			req.Header.Set("Refer", "https://m.pgsoft-games.com/")
			req.Header.Set("Origin", "https://m.pgsoft-games.com")

			resp, err := http.DefaultClient.Do(req)
			// resp, err := http.Get(ul.JoinPath(r.URL.Path).String())
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
			if indexRegexp.MatchString(inuri) {
				body = bytes.Replace(body, []byte("</script>"), []byte(`</script><script>window.gtag=console.log;window.dataLayer=[];</script>`), 1)
			}
			writefile(file, body)
			http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
		}
	})

	go func() {
		////证书
		//certfile, _ := lazy.RouteFile.Get("certfile")
		//keyfile, _ := lazy.RouteFile.Get("keyfile")
		//err := http.ListenAndServeTLS(addr, certfile, keyfile, serveMux)
		err := http.ListenAndServe(addr, serveMux)
		log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}

func writefile(file string, data []byte) {
	dir := path.Dir(file)
	if e := os.MkdirAll(dir, 0755); e == nil || os.IsExist(e) {
		os.WriteFile(file, data, 0644)
	}
}

//func Verify(w http.ResponseWriter, r *http.Request) {
//	ps := r.URL.Query()
//	ul := gamedata.Get().ReverseProxyUrls.Get(r.Host)
//	if ul == nil {
//		w.WriteHeader(404)
//		return
//	}
//
//	file := path.Join("cache", ul.Host, r.URL.Path+"-"+ps.Get("symbol")+".html")
//	if fileinfo, err := os.Stat(file); err == nil {
//		data, err := os.ReadFile(file)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusNotFound)
//			return
//		}
//
//		data = html5gameDo_hook(data, w, r)
//		http.ServeContent(w, r, file, fileinfo.ModTime(), bytes.NewReader(data))
//		return
//	}
//}

type PGError struct {
	Cd  string `json:"cd"`
	Msg string `json:"msg"`
	Tid string `json:"tid"`
}

func (e *PGError) Error() string {
	return e.Msg
}
