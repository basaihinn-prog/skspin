package staticproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/ut"
	"game/duck/lazy"
	"game/duck/ut2"
	"game/service/ppgateway/internal/gamedata"
	"io"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// excute before save
type contenthook = func([]byte) []byte

// excute before send
// type respcontenthook = func([]byte, http.ResponseWriter, *http.Request) []byte
type respcontenthook = func([]byte) []byte

func static_assets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ul := gamedata.Get().ReverseProxyUrls.Get(r.Host)
	if ul == nil {
		w.WriteHeader(404)
		return
	}

	ul.Path = r.URL.Path
	ul.RawQuery = ""

	var (
		file = "cache/" + ul.Host + r.URL.Path
	)

	var contenthooks = []contenthook{}
	var respcontenthooks = []respcontenthook{}

	pth := r.URL.Path
	switch {
	case strings.HasSuffix(pth, "/build.js"):
		contenthooks = append(contenthooks, buildjs_nop_host_valid)
		contenthooks = append(contenthooks, buildjs_gtag_hook)
		contenthooks = append(contenthooks, buildjs_Number_hook)
		contenthooks = append(contenthooks, buildjs_Lang_hook)

		if lazy.CommCfg().IsDev {
			respcontenthooks = append(respcontenthooks, buildjs_post)
		}
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
			// data = hook(data, w, r)
			data = hook(data)
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

	lg = lg.With("from", "REMOTE")

	for _, hook := range contenthooks {
		body = hook(body)
	}

	if !cache_nostore(resp.Header) {
		ut.Writefile(file, body)
	}

	for _, hook := range respcontenthooks {
		body = hook(body)
	}
	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
}

//func static_assets2(w http.ResponseWriter, r *http.Request) {
//	ul := gamedata.Get().ReverseProxyUrls.Get(r.Host)
//	if ul == nil {
//		w.WriteHeader(404)
//		return
//	}
//	ul.Host = `huidubet.prerelease-env.biz`
//	ul.Path = r.URL.Path
//	ul.RawQuery = ""
//
//	var (
//		file = "cache/" + ul.Host + r.URL.Path
//	)
//
//	var contenthooks = []contenthook{}
//	var respcontenthooks = []respcontenthook{}
//
//	pth := r.URL.Path
//	switch {
//	case strings.HasSuffix(pth, "/build.js"):
//		contenthooks = append(contenthooks, buildjs_nop_host_valid)
//		contenthooks = append(contenthooks, buildjs_gtag_hook)
//
//		if lazy.CommCfg().IsDev {
//			respcontenthooks = append(respcontenthooks, buildjs_post)
//		}
//	}
//
//	var err error
//	if file[len(file)-1] == '/' {
//		file += "index.html"
//	}
//
//	inuri := r.URL.RequestURI()
//	lg := slog.With(
//		"file", file,
//		"inuri", inuri,
//		"method", r.Method,
//	)
//
//	defer func() {
//		lg = lg.With(
//			"error", err,
//		)
//		lg.Info("req resources")
//	}()
//
//	if w.Header().Get("Cache-Control") == "" {
//		w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=3600")
//	}
//
//	//replay页面
//	if strings.Contains(r.Host, "replay.") && len(r.URL.Path) == 8 {
//		file = `service/ppgateway/ppcomm/html/ppReplayTemplate.html`
//		if fileinfo, err := os.Stat(file); err == nil {
//			data, err := os.ReadFile(file)
//			if err != nil {
//				http.Error(w, err.Error(), http.StatusNotFound)
//				return
//			}
//			data = replayTemplateDo_hook(data, w, r)
//			for _, hook := range respcontenthooks {
//				// data = hook(data, w, r)
//				data = hook(data)
//			}
//			http.ServeContent(w, r, file, fileinfo.ModTime(), bytes.NewReader(data))
//			return
//		} else {
//			slog.Error("ReplayTemplate is not exist")
//			// w.WriteHeader(http.Status)
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//	}
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	if fileinfo, err := os.Stat(file); err == nil {
//		// http.ServeFile(w, r, file)
//		data, err := os.ReadFile(file)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusNotFound)
//			return
//		}
//
//		for _, hook := range respcontenthooks {
//			// data = hook(data, w, r)
//			data = hook(data)
//		}
//		http.ServeContent(w, r, file, fileinfo.ModTime(), bytes.NewReader(data))
//		return
//	}
//
//	outuri := ul.String()
//	lg = lg.With("outuri", outuri)
//
//	req := lo.Must(http.NewRequest(http.MethodGet, outuri, nil))
//
//	resp, err := http.DefaultClient.Do(req)
//
//	if err != nil {
//		// w.WriteHeader(http.Status)
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer resp.Body.Close()
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		// w.WriteHeader(http.Status)
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	lg = lg.With("Status", resp.Status)
//
//	if resp.StatusCode != http.StatusOK {
//		w.WriteHeader(resp.StatusCode)
//		return
//	}
//
//	lg = lg.With("from", "REMOTE")
//
//	for _, hook := range contenthooks {
//		body = hook(body)
//	}
//
//	if !cache_nostore(resp.Header) {
//		ut.Writefile(file, body)
//	}
//
//	for _, hook := range respcontenthooks {
//		body = hook(body)
//	}
//	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(body))
//}

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

func replayTemplateDo_hook(content []byte, w http.ResponseWriter, r *http.Request, mgckey string, rpInfo ppReplayTokenMapBody) []byte {
	host := r.Host
	if strings.HasPrefix(host, "cdn.") {
		host = host[4:]
	}
	content = bytes.ReplaceAll(content, []byte("{{mainHots}}"), []byte(ut2.Domain(host)))
	content = bytes.ReplaceAll(content, []byte("{{token}}"), []byte(r.URL.Path))
	ps := r.URL.Query()
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	{

		var oGameConfig ReplayGameConfig
		coll2 := db.Collection2(rpInfo.Gid, "BetHistory")
		gid := bytes.ReplaceAll([]byte(rpInfo.Gid), []byte("pp_"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("{{gid}}"), gid)
		content = bytes.ReplaceAll(content, []byte("{{WebTitle}}"), []byte(comm.PPGameWebTitle[string(gid)]))
		param := bson.M{"betId": rpInfo.BetId}
		projection := bson.M{"gameConfig": 1, "_id": 0} // 只返回 "fieldName"，不返回 "_id"
		findOptions := options.FindOne().SetProjection(projection)
		err := coll2.FindOne(context.TODO(), param, findOptions).Decode(&oGameConfig)
		if err != nil {
			slog.Error("Get oGameConfig err", "err: ", err)
			return nil
		}
		// adjust gameconfig
		submatch := gameConfigRegexp.FindSubmatchIndex(content)

		i, j := submatch[2], submatch[3]
		jsondata := content[i:j]
		var gameconfig map[string]json.RawMessage
		json.Unmarshal(jsondata, &gameconfig)
		//todo 游戏中gameconfig替换现在的
		gameconfig["currency"] = ut.GetJsonRawMust(ps.Get("currency"))
		gameconfig["currencyOriginal"] = ut.GetJsonRawMust(ps.Get("currency"))
		gameconfig["ReplaySystemUrl"] = ut.GetJsonRawMust(ps.Get(""))
		gameconfig["ReplaySystemContextPath"] = ut.GetJsonRawMust(ps.Get(""))
		gameconfig["mgckey"] = ut.GetJsonRawMust(mgckey)
		gameconfig["replayRoundId"] = ut.GetJsonRawMust(rpInfo.BetId)
		lang := "en"
		if ps.Get("lang") != "" {
			lang = ps.Get("lang")
		} else if rpInfo.Lang != "" {
			lang = rpInfo.Lang
		}
		gameconfig["lang"] = ut.GetJsonRawMust(lang)
		//https://common-static.prerelease-env.biz/gs2c/common/v1/games-html5/games/vs/vs20olympx/
		host := ""
		game := ""
		gameconfig["Datapath"] = ut.GetJsonRawMust(fmt.Sprintf("https://%s/gs2c/common/v1/games-html5/games/vs/%s/", host, game))

		outdata, _ := json.Marshal(gameconfig)
		content = slices.Replace(content, i, j, outdata...)
	}
	return content
}

type ppReplayTokenMapBody struct {
	Token    string `json:"token"`
	Gid      string `json:"gid"`
	BetId    string `json:"betId"`
	Init     string `json:"init"`
	CreateAt int64  `json:"createAt"`
	Lang     string `json:"lang"`
}

type ReplayGameConfig struct {
	AmountType              string   `json:"amountType" bson:"amounttype"`
	EnvironmentId           int      `json:"environmentId" bson:"environmentid"`
	CurrencyOriginal        string   `json:"currencyOriginal" bson:"currencyoriginal"`
	ReplayMode              bool     `json:"replayMode" bson:"replaymode"`
	Jurisdiction            string   `json:"jurisdiction" bson:"jurisdiction"`
	Currency                string   `json:"currency" bson:"currency"`
	SessionTimeout          string   `json:"sessionTimeout" bson:"sessiontimeout"`
	StyleName               string   `json:"styleName" bson:"stylename"`
	Lang                    string   `json:"lang" bson:"lang"`
	Region                  string   `json:"region" bson:"region"`
	ReplaySystemUrl         string   `json:"replaySystemUrl" bson:"replaysystemurl"`
	BrandRequirements       string   `json:"brandRequirements" bson:"brandrequirements"`
	ReplaySystemContextPath string   `json:"replaySystemContextPath" bson:"replaysystemcontextpath"`
	ReplayRoundId           string   `json:"replayRoundId" bson:"replayroundid"`
	Mgckey                  string   `json:"mgckey" bson:"mgckey"`
	Datapath                string   `json:"datapath" bson:"datapath"`
	SessionKey              []string `json:"sessionKey" bson:"sessionkey"`
	SessionKeyV2            []string `json:"sessionKeyV2" bson:"sessionkeyv2"`
}
