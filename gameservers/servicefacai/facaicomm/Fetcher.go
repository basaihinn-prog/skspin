package facaicomm

import (
	"bytes"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/phuslu/iploc"
	"serve/comm/mq"
	"serve/comm/ut"
)

type Fetcher struct {
	game, user string

	mgckey string
	// location *url.URL
	ul         string
	httpClient *http.Client
	lastResp   Variables
	c          float64
	l          int
	ty         int
}

func initMQ() {
	//var m map[string]string
	//data := lo.Must(os.ReadFile("/config/grpc_route.yaml"))
	//yaml.Unmarshal(data, &m)
	//
	//addr := m["proxy.mq"]
	mq.ConnectServerMust("127.0.0.1:11002")
}

func getPubIP(c *http.Client) string {
	req, _ := http.NewRequest("GET", "http://ip.sb/", nil)
	req.Header.Set("User-Agent", "curl/7.81.0")
	body, _, _ := ut.DoHttpReq(c, req)
	os.Stdout.Write(body)
	body = bytes.TrimSpace(body)
	return string(body)
}

func NewFetcher(game, user string) *Fetcher {
	ip := getPubIP(http.DefaultClient)
	loc := string(iploc.Country(net.ParseIP(ip)))
	c := http.DefaultClient

	game = strings.TrimPrefix(game, "pp_")
	// game = strings.Split(game, "_")[1]
	if loc != "CN" {
		if mq.NC() == nil {
			initMQ()
		}
		//endpoint := lo.Must(ip2worldpub.GetEndpoint())
		//if endpoint.IP == "" {
		//	endpoint.IP = "54.251.234.111"
		//	endpoint.Port = 443
		//}
		//urlproxy, _ := url.Parse(fmt.Sprintf("http://%s:%d", endpoint.IP, endpoint.Port))
		//c = &http.Client{
		//	Transport: &http.Transport{
		//		Proxy: http.ProxyURL(urlproxy),
		//	},
		//	Timeout: 30 * time.Second,
		//}
		//fmt.Sprintf("http://%s:%d", endpoint.IP, endpoint.Port))
		c = &http.Client{
			//Transport: &http.Transport{
			//	Proxy: http.ProxyURL(urlproxy),
			//},
			Timeout: 30 * time.Second,
		}
	}

	mgckey, location, err := GetHuiDuNewmgcKey(c, "")
	if err != nil || location == nil {
		slog.Error("new fetcher,", " err", err, "location", location, "name", user)
		return nil
	}

	location.Path = "/gs2c/ge/v3/gameService"
	fetcher := &Fetcher{
		game:       game,
		httpClient: c,
		user:       user,
		mgckey:     mgckey,
		ul:         location.String(),
	}
	return fetcher
}

//
//func (f *Fetcher) Do(ps Variables) Variables {
//	// form, _ := url.ParseQuery("action=doSpin&symbol=vs20olympx&c=0.5&l=20&bl=0&index=3&counter=5&repeat=0&mgckey=123")
//	lo.Must0(ps.Str("action") != "")
//
//	if ps.Str("action") == "doSpin" {
//		ps.SetFloat("c", f.c)
//		ps.SetInt("l", f.l)
//		ps.Set("bl", "0")
//	}
//	if f.ty == GameTypeGame {
//		ps.SetInt("pur", 0)
//	}
//	if f.ty == GameTypeSuperGame1 {
//		ps.SetInt("pur", 1)
//	}
//	ps.Set("symbol", f.game)
//	ps.SetInt("index", f.lastResp.Int("index")+1)
//	ps.SetInt("counter", f.lastResp.Int("counter")+1)
//	ps.Set("mgckey", f.mgckey)
//	ps.Set("repeat", "0")
//
//	req, _ := http.NewRequest(http.MethodPost, f.ul, strings.NewReader(ps.Encode()))
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
//	// http.DefaultClient.Do(req)
//	body, _ := lo.Must2(ut.DoHttpReq(f.httpClient, req))
//
//	vars := ParseVariables(string(body))
//	f.lastResp = vars
//
//	slog.Info("jl_robot_center.DO", "ps", ps.Encode(),
//		"user", f.user,
//		"na", vars["na"],
//		"resp", mux.TruncMsg(body, 256),
//	)
//	return vars
//}
//
//func (f *Fetcher) NextAction() string {
//	return f.lastResp.Str("na")
//}

func (f *Fetcher) LastResp() Variables {
	return f.lastResp
}
