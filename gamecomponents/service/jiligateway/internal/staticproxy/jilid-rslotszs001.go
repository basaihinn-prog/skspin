package staticproxy

import (
	"bytes"
	_ "embed"
	"game/comm/ut"
	"game/duck/ut2"
	"game/duck/ut2/jwtutil"
	"game/service/jiligateway/internal/gamedata"
	"game/service/jiligateway/jilicomm"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
)

func init() {
	jilid_rslotszs001_mux.Handle("/", http.HandlerFunc(static_assets))
}

//go:embed html/action.html
var actionHtml []byte

//go:embed html/ingame.html
var inGameHtml []byte

//go:embed html/ingameOfficial.html
var ingameOfficial []byte

//go:embed html/test.html
var test []byte

var (
	regIngame = regexp.MustCompile(`^(/\w\w-\w\w)?/(ingame|intro)`)
	// regIngame2 = regexp.MustCompile(`^/ingame`)
	// regIntro = regexp.MustCompile(`^/\w\w-\w\w/intro$`)
	regLang         = regexp.MustCompile(`^/language/\w\w-\w\w$`)
	regJILI         = regexp.MustCompile(`^/\w+/JILI_NEW_\w+\.png$`)
	regOtherIndexJS = regexp.MustCompile(`^/\w+/assets/other/index.\w+.js$`)

	regGTM         = regexp.MustCompile(`GTM-\w{7}`)
	regAction      = regexp.MustCompile(`^(/\w\w-\w\w)?/(promote/action|landingpage)`)
	reqJscc        = regexp.MustCompile(`cc.*.js`)
	reqJsccPattern = regexp.MustCompile(`if\s*\(\s*"\.webp"\s*===\s*[\s\S]*?\)\s*continue\s*;`)
)

type contenthook = func([]byte) []byte

func static_assets(w http.ResponseWriter, r *http.Request) {
	var storeul *url.URL
	gname := jiliUrlGameNameGetter(*r.URL)
	//if regLang.MatchString(r.URL.Path) { //获取国际化接口的游戏名
	//	query := r.URL.Query()
	//	groupIds := query.Get(`GroupIds[]`)
	//	if len(groupIds) > 0 {
	//		groupIdsSplit := strings.Split(groupIds, "_")
	//		if len(groupIdsSplit) == 2 {
	//			gname = strings.ToLower(groupIdsSplit[1])
	//		}
	//	}
	//}
	//if gname == `_nuxt` && r.Header.Get(`Referer`) != `` {// todo 为了区分游戏历史的资源使用反代的资源，暂时不需要了 都使用灰度还是rp来着的历史页面
	//	rawURL := r.Header.Get(`Referer`)
	//	// 解析 URL
	//	parsedURL, err := url.Parse(rawURL)
	//	if err != nil {
	//		fmt.Println("URL 解析失败:", err)
	//		return
	//	}
	//	// 获取查询参数（Query Parameters）
	//	queryParams := parsedURL.Query()
	//	// 提取特定参数
	//	token := queryParams.Get("token")
	//	_, game, err := jwtutil.ParseTokenData(token)
	//	if err != nil {
	//		return
	//	}
	//	gname = GetJiliGameName(game)
	//}
	if _, ok := jilicomm.ChangeGameMap[gname]; ok && !strings.HasPrefix(r.Host, "tada") {
		storeul = gamedata.Get().ReverseProxyUrls2.Get(r.Host)
		if storeul == nil {
			w.WriteHeader(404)
			return
		}
	} else if _, ok := jilicomm.ChangeGameMap3[gname]; ok && !strings.HasPrefix(r.Host, "tada") {
		storeul = gamedata.Get().ReverseProxyUrls3.Get(r.Host)
		if storeul == nil {
			w.WriteHeader(404)
			return
		}
	} else if _, ok := jilicomm.ChangeGameMap4[gname]; ok && !strings.HasPrefix(r.Host, "tada") {
		storeul = gamedata.Get().ReverseProxyUrls4.Get(r.Host)
		if storeul == nil {
			w.WriteHeader(404)
			return
		}
	} else {
		storeul = gamedata.Get().ReverseProxyUrls.Get(r.Host)
		if storeul == nil {
			w.WriteHeader(404)
			return
		}
	}

	ul := *storeul
	ul.Path = r.URL.Path
	ul.RawQuery = ""
	// ul.RawQuery = r.URL.RawQuery

	var (
		file string
	)

	var contenthooks []contenthook
	// var respcontenthooks = []contenthook{}

	switch {
	// case false && strings.HasPrefix(r.Host, "jilid-rslotszs001.") && regJILI.Match([]byte(r.URL.Path)):
	// 	// https://jilid-rslotszs001.kafa010.com/tks/JILI_NEW_1717655961.png
	// 	file = path.Join("cache", ul.Host, "_custom", "JILI_NEW.png")
	case ut.HasPrefix(r.Host, "uat-language-api.", "language-api.", "cdn.uat-language-api.", "cdn.language-api."):
		// https://uat-language-api.kafa010.com/language/zh-CN?t=202406281551
		// https://uat-language-api.kafa010.com/language/zh-CN?t=202406281551&GroupIds[]=INTRO_CSH&GroupIds[]=INTRO_SHARE

		query := r.URL.Query()
		query.Del("t")
		// r.URL.RawQuery = query.Encode()
		ul.RawQuery = query.Encode()
		file = path.Join("cache", ul.Host, r.URL.RequestURI())
		if regLang.MatchString(r.URL.Path) {
			groupIds := query[`GroupIds[]`]
			gname = `/comm`
			if len(groupIds) > 0 {
				groupIdsSplit := strings.Split(groupIds[0], "_")
				if len(groupIdsSplit) == 2 {
					gname = `/` + strings.ToLower(groupIdsSplit[1])
				}
			}
			file = path.Join("cache", ul.Host, r.URL.Path, gname)
			contenthooks = append(contenthooks, func(content []byte) []byte {
				// 定义正则表达式
				re := regexp.MustCompile(`(?i)as bet\s*=\s*{(\w+)}|when Bet\s*=\s*{(\w+)}`)
				// 替换为你想要的内容
				return re.ReplaceAll(content, []byte("current bet"))
				// return bytes.ReplaceAll(content, []byte("GTM-P8B36RZ"), []byte("GTM-1234567"))
			})
		}
		// w.Header().Set("Cache-Control", "private, max-age=60")

	case regOtherIndexJS.MatchString((r.URL.Path)):
		// /data/game/bin/cache/jilid.rslotszs001.com/samba/assets/other/index.76138.js
		// /data/game/bin/cache/wbgame.bd33fgabh.com/ge/assets/other/index.c8bbb.js

		ul.RawQuery = ""
		file = "cache/" + ul.Host + r.URL.Path

		contenthooks = append(contenthooks, func(content []byte) []byte {
			return bytes.ReplaceAll(content, []byte("m_discountBtn.node.active=t"), []byte("m_discountBtn.node.active=false"))
		})
	case ut.HasPrefix(r.Host, "uat-history.", "history.", "cdn.uat-history.", "cdn.history.") && (false ||
		regIngame.MatchString(r.URL.Path) ||
		// regIntro.MatchString(r.URL.Path) ||
		false):
		// https://uat-history.jlfafafa3.com/en-US/ingame?token=6ec798030314f4a3b4de39c5a75e3419b221e6fb&game=2&posthost=uat-wbgame.jlfafafa3.com&&sac=0
		// https://uat-history.kafa010.com/en-US/intro?token=eyJQIjoxMDA3OTQsIkUiOjE3MTk2MDI0OTgsIlMiOjEwMDMsIkQiOiJqaWxpXzJfY3NoIn0.b6tyoxhN7onaBCsF6iHR52nWm2q20URgInANht-33eI&game=2&posthost=jilid-rslotszs001.kafa010.com&&sac=0
		query := r.URL.Query()
		file = path.Join("cache", ul.Host, "ingame-"+query.Get("game")+".html")
		file = path.Join("cache", ul.Host, "ingame"+".html")

		ul.Path = "/ingame/gamehistory"
		ul.RawQuery = ""

		contenthooks = append(contenthooks, func(content []byte) []byte {
			return regGTM.ReplaceAllLiteral(content, []byte("GTM-1234567"))
			// return bytes.ReplaceAll(content, []byte("GTM-P8B36RZ"), []byte("GTM-1234567"))
		})
		contenthooks = append(contenthooks, func(content []byte) []byte {
			// return bytes.ReplaceAll(content, []byte("GTM-P8B36RZ"), []byte("GTM-1234567"))
			s := `<script>
let intervalID = 0
function checkLayoutMenu20240715() {
	const nodes=document.getElementsByClassName('layout-menu');
	if(nodes.length){
		node=nodes[0];
		node.style.display='none';
		clearInterval(intervalID);
		console.log("clearInterval", intervalID);
	}
}
intervalID = window.setInterval(checkLayoutMenu20240715,100);
console.log("intervalID", intervalID);
</script></body>`

			return bytes.Replace(content, []byte("</body>"), []byte(s), 1)
		})
	case ut.HasPrefix(r.Host, "uat-history.", "history.", "cdn.uat-history.", "cdn.history.") &&
		regAction.MatchString(r.URL.Path):
		http.ServeContent(w, r, "action", time.Now(), bytes.NewReader(actionHtml))
		return
	case reqJscc.MatchString(r.URL.Path):
		file = "cache/" + ul.Host + r.URL.Path
		contenthooks = append(contenthooks, func(content []byte) []byte {
			return reqJsccPattern.ReplaceAllLiteral(content, []byte(""))
		})
	default:
		ul.RawQuery = ""
		file = "cache/" + ul.Host + r.URL.Path
	}

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
			// "Header", w.Header(),
			"error", err,
		)
		lg.Info("req resources")
	}()

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// if is_index_html(file) {
	// 	w.Header().Set("Cache-Control", "public, max-age=3600")
	// } else {
	// 	w.Header().Set("Cache-Control", "public, max-age=2592000")
	// }

	if w.Header().Get("Cache-Control") == "" {
		w.Header().Set("Cache-Control", "public, max-age=2592000, must-revalidate, proxy-revalidate")
	}

	// only dev used
	// if strings.HasSuffix(file, ".js") {
	// w.Header().Set("Cache-Control", "no-cache")
	// }
	if regIngame.MatchString(r.URL.String()) {
		ps := r.URL.Query()
		token := ps.Get("token")
		_, game, err := jwtutil.ParseTokenData(token)
		if err != nil {
			return
		}
		body := ingameOfficial
		if _, ok := jilicomm.ChangeGameMap44[game]; ok {
			//body = ingameOfficial
		}
		mainHost := ut2.Domain(r.Host)
		body = bytes.ReplaceAll(body, []byte("{{mainHost}}"), []byte(mainHost))

		if strings.HasPrefix(game, "tada") {
			body = bytes.ReplaceAll(body, []byte(`GAME_ICON_BRAND_FILENAME: "jili",`), []byte(`GAME_ICON_BRAND_FILENAME: "tada",`))
		}
		// 设置缓存控制
		w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
		return
	}
	if _, err := os.Stat(file); err == nil {
		http.ServeFile(w, r, file)
		return
	}

	// ul.Path = r.URL.Path
	// ul.RawQuery = r.URL.RawQuery

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

	for _, hook := range contenthooks {
		body = hook(body)
	}

	if !cache_nostore(resp.Header) {
		// if file[len(file)-1] == '/' {
		// 	file += "index.html"
		// }
		ut.Writefile(file, body)
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

func is_index_html(pth string) bool {
	return strings.HasSuffix(pth, "/") || strings.HasSuffix(pth, "/index.html")
}

//func ingameDo_hook(content []byte, w http.ResponseWriter, r *http.Request) []byte {
//	{
//		//re := regexp.MustCompile(`config\s*=\s*(\{.*?\});`)
//		//matches := re.FindStringSubmatch(htmlContent)
//		// adjust cashierUrl
//		submatch := configRegexp.FindSubmatchIndex(content)
//		if submatch != nil {
//			i, j := submatch[2], submatch[3]
//			content = slices.Replace(content, i, j, []byte(``)...)
//		}
//		i, j := submatch[2], submatch[3]
//		jsondata := content[i:j]
//		var gameconfig map[string]json.RawMessage
//		json.Unmarshal(jsondata, &gameconfig)
//		delete(gameconfig, "IOS_DOWNLOAD_LINK")
//		delete(gameconfig, "IOS_MOBILE_PROVISION_LINK")
//		delete(gameconfig, "ANDROID_DOWNLOAD_LINK")
//		mainHost := ut2.Domain(r.Host)
//		gameconfig["NEW_WEB_API_URL"] = ut.GetJsonRawMust(fmt.Sprintf("https://uat-history-api.%s", mainHost))
//		gameconfig["LANG_API_URL"] = ut.GetJsonRawMust(fmt.Sprintf("https://uat-language-api.%s", mainHost))
//		gameconfig["LOG_API_URL"] = ut.GetJsonRawMust(fmt.Sprintf("https://uat-log-api.%s", mainHost))
//		gameconfig["WEB_RESOURCE_CDN_URL"] = ut.GetJsonRawMust(fmt.Sprintf("https://uat-web-cdn.%s/static", mainHost))
//		outdata, _ := json.Marshal(gameconfig)
//		content = slices.Replace(content, i, j, outdata...)
//	}
//	return content
//}

func jiliUrlGameNameGetter(url2 url.URL) string {
	gname := ""
	uurl := url2.String()
	// 查找第一个 "/"
	start := strings.Index(uurl, "/")
	if start != -1 {
		// 从第一个 "/" 开始查找下一个 "/"
		end := strings.Index(uurl[start+1:], "/")
		if end != -1 {
			// 提取从第一个 "/" 到下一个 "/" 之间的内容
			gname = uurl[start+1 : start+1+end]
		} else {
			// 如果没有找到下一个 "/"
			gname = uurl[start:]
		}
	}
	return gname
}

func GetJiliGameName(game string) string {
	res := game
	if strings.Contains(game, "jili_") {
		if _, ok := jilicomm.AllChangeGameNameMap[game]; ok {
			res = jilicomm.AllChangeGameNameMap[game]
		}
	} else {
		if _, ok := jilicomm.AllChangeGameMap[game]; ok {
			res = jilicomm.AllChangeGameMap[game]
		}
	}
	return res
}
