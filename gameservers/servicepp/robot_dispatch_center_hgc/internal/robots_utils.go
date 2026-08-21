package internal

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"serve/comm/db"
	"serve/comm/mq"
	"serve/comm/plats/pp"
	"serve/comm/ut"
	ip2worldpub "serve/servicepp/ip2world/pub"
	"serve/servicepp/ppcomm"
	"strconv"
	"strings"
	"time"

	"github.com/phuslu/iploc"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/html"
)

// Response 是整个响应的结构体
type Response struct {
	Status    int      `json:"status"`    // 状态码
	Msg       string   `json:"msg"`       // 消息
	Data      GameData `json:"data"`      // 游戏数据
	Timestamp int64    `json:"timestamp"` // 时间戳
}

// GameData 是游戏数据的结构体
type GameData struct {
	ID               int64  `json:"id"`                         // 游戏ID
	Status           int    `json:"status"`                     // 状态
	IDs              []int  `json:"ids,omitempty"`              // ID列表，可选
	AddUserId        int64  `json:"addUserId,omitempty"`        // 添加用户ID，可选
	ModUserId        int64  `json:"modUserId,omitempty"`        // 修改用户ID，可选
	CompanyId        int64  `json:"companyId,omitempty"`        // 公司ID，可选
	CompanyName      string `json:"companyName"`                // 公司名称
	GameName         string `json:"gameName"`                   // 游戏名称
	GameNameZh       string `json:"gameNameZh"`                 // 游戏名称（中文）
	GameIcon         string `json:"gameIcon"`                   // 游戏图标URL
	GameIconPath     string `json:"gameIconPath,omitempty"`     // 游戏图标路径，可选
	GameCategory     string `json:"gameCategory,omitempty"`     // 游戏分类，可选
	GamePlatform     string `json:"gamePlatform"`               // 游戏平台
	GameType         string `json:"gameType"`                   // 游戏类型
	GameKey          string `json:"gameKey,omitempty"`          // 游戏键，可选
	GameId           string `json:"gameId"`                     // 游戏ID
	Jackpot          string `json:"jackpot,omitempty"`          // 彩池，可选
	GameRanking      string `json:"gameRanking,omitempty"`      // 游戏排名，可选
	GameRtp          string `json:"gameRtp"`                    // 返回玩家的百分比
	FeatureBuy       string `json:"featureBuy,omitempty"`       // 特性购买，可选
	Volatility       string `json:"volatility,omitempty"`       // 波动性，可选
	MaxExposure      string `json:"maxExposure,omitempty"`      // 最大暴露，可选
	MultiLanguage    string `json:"multiLanguage,omitempty"`    // 多语言，可选
	Remark           string `json:"remark"`                     // 备注
	Address          string `json:"address,omitempty"`          // 地址，可选
	GameNum          int    `json:"gameNum,omitempty"`          // 游戏编号，可选
	GameCategoryCode string `json:"gameCategoryCode,omitempty"` // 游戏分类代码，可选
	Gameplat         string `json:"gameplat,omitempty"`         // 游戏平台（备用），可选
	Lang             string `json:"lang"`                       // 支持的语言
	Currency         string `json:"currency"`                   // 支持的货币
	Rate             string `json:"rate,omitempty"`             // 汇率，可选
	GameScreenshot   string `json:"gameScreenshot,omitempty"`   // 游戏截图，可选
	GameFile         string `json:"gameFile,omitempty"`         // 游戏文件，可选
	MaxWin           string `json:"maxWin,omitempty"`           // 最大获胜，可选
	IconGameKey      string `json:"iconGameKey,omitempty"`      // 图标游戏键，可选
	BrandName        string `json:"brandName,omitempty"`        // 品牌名称，可选
	GameEncrypt      string `json:"gameEncrypt"`                // 游戏加密
	Support          string `json:"support,omitempty"`          // 支持，可选
}
type ResponseGame struct {
	Data      DataGame `json:"data"`
	Msg       string   `json:"msg"`
	Status    int      `json:"status"`
	Timestamp int64    `json:"timestamp"`
}

type DataGame struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func ToMapFI(data string) map[float64]int {
	// 使用逗号分割字符串
	decoded, _ := url.PathUnescape(data)
	pairs := strings.Split(decoded, ",")

	result := make(map[float64]int)

	for _, pair := range pairs {
		// 使用冒号分割键值对
		kv := strings.Split(pair, ":")

		// 检查键值对格式是否正确
		if len(kv) != 2 {
			fmt.Errorf("invalid pair format: %s", pair)
			return nil
		}

		// 转换键为 float64
		key, err := strconv.ParseFloat(kv[0], 64)
		if err != nil {
			fmt.Errorf("invalid key format: %s", kv[0])
			return nil
		}

		// 转换值为 int
		value, err := strconv.Atoi(kv[1])
		if err != nil {
			fmt.Errorf("invalid value format: %s", kv[1])
			return nil
		}

		// 存入 map
		result[key] = value
	}

	return result
}

func getHuiDuLocation(dataGameSrc string, client *http.Client) (string, *url.URL) {
	req2, err := http.NewRequest(http.MethodGet, dataGameSrc, nil)
	if err != nil {
		slog.Error(dataGameSrc+" err", "err", err)
		return "", nil
	}

	// time.Sleep(time.Second)

	req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	trans := cmp.Or(client.Transport, http.DefaultTransport)
	resp2, err := trans.RoundTrip(req2)
	if err != nil {
		slog.Error("http.DefaultTransport.RoundTrip err", "err", err)
		return "", nil
	}

	io.Copy(os.Stdout, resp2.Body)

	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusFound {
		return "", nil
	}

	location := resp2.Header.Get("Location")
	location2, err := GetRedirect(location, client)
	if err != nil {
		err = errors.New("重定向失败")
		return "", nil
	}
	query := location2.Query()
	mgckey := query.Get("mgckey")

	return mgckey, location2
}
func getHotGameClubLocation(dataGameSrc string, client *http.Client) (string, *url.URL) {
	req2, err := http.NewRequest(http.MethodGet, dataGameSrc, nil)
	if err != nil {
		slog.Error(dataGameSrc+" err", "err", err)
		return "", nil
	}

	// time.Sleep(time.Second)

	req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	trans := cmp.Or(client.Transport, http.DefaultTransport)
	resp2, err := trans.RoundTrip(req2)
	if err != nil {
		slog.Error("http.DefaultTransport.RoundTrip err", "err", err)
		return "", nil
	}

	io.Copy(os.Stdout, resp2.Body)

	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusFound {
		return "", nil
	}

	location := resp2.Header.Get("Location")
	location2, err := url.Parse(location)
	if err != nil {
		return "", nil
	}
	query := location2.Query()
	mgckey := query.Get("mgckey")

	return mgckey, location2
}

// 提取 src 的值
func extractIframeSrc(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var dataGameSrc string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "iframe" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					dataGameSrc = attr.Val
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	if dataGameSrc == "" {
		fmt.Println(htmlContent)
		return "", fmt.Errorf("未找到 src 属性")
	}

	return dataGameSrc, nil
}

// 提取 data-game-src 的值
func extractDataGameSrc(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var dataGameSrc string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "iframe" {
			for _, attr := range n.Attr {
				if attr.Key == "data-game-src" {
					dataGameSrc = attr.Val
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	if dataGameSrc == "" {
		return "", fmt.Errorf("未找到 data-game-src 属性")
	}

	return dataGameSrc, nil
}
func getPubIP(c *http.Client) string {
	req, _ := http.NewRequest("GET", "http://ip.sb/", nil)
	req.Header.Set("User-Agent", "curl/7.81.0")
	body, _, _ := ut.DoHttpReq(c, req)
	//os.Stdout.Write(body)
	body = bytes.TrimSpace(body)
	return string(body)
}
func GetMgckey(game, user string, client *http.Client) (mgckey string, location *url.URL, err error) {
	// client req the rest
	uid, guid, cash := "ceshi1938", "4ae52ed2e1a8c353878ba65ed7791ac4", "5000000"
	uid = user
	guid = ppcomm.GuidMap[game]
	uurl, _ := pp.GetKey(uid, guid, cash)
	//req, err := http.NewRequest(http.MethodGet, uurl, nil)
	//
	//if err != nil {
	//	return
	//}
	//
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	//
	//trans := cmp.Or(client.Transport, http.DefaultTransport)
	//resp, err := trans.RoundTrip(req)
	//if err != nil {
	//	return
	//}
	//
	//defer resp.Body.Close()
	//if resp.StatusCode != http.StatusFound {
	//	err = errors.New(resp.Status)
	//	return
	//}
	//
	//location, err = url.Parse(resp.Header.Get("Location"))
	//if err != nil {
	//	return
	//}
	location, err = GetRedirect(uurl, client)
	if err != nil {
		err = errors.New("重定向失败")
		return
	}
	if location.String() == "" {
		err = errors.New("重定向失败")
		return
	}
	location, err = GetRedirect(location.String(), client)
	if err != nil {
		err = errors.New("重定向失败")
		return
	}
	query := location.Query()
	mgckey = query.Get("mgckey")
	location.RawQuery = ""

	return
}

// 重定向方法
func GetRedirect(uurl string, client *http.Client) (*url.URL, error) {
	var location *url.URL
	req, err := http.NewRequest(http.MethodGet, uurl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	trans := cmp.Or(client.Transport, http.DefaultTransport)
	resp, err := trans.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		err = errors.New(resp.Status)
		return nil, err
	}

	location, err = url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return nil, err
	}
	return location, nil
}

// 查看数据是否足够
func CheckGameEnough(gameId string, dataType int) bool {
	var mission GameInfo
	coll := db.Collection2("robot", "game")
	err := coll.FindOne(context.TODO(), bson.M{"_id": gameId}).Decode(&mission)
	if err == mongo.ErrNoDocuments {
		return false
	}
	var normalCnt = mission.NormalCount // 需要拉取普通盘的次数
	var gameCnt = mission.BuyCount      // 需要拉取游戏的次数
	needMap := map[int]int64{}
	if normalCnt > 0 && dataType == ppcomm.GameTypeNormal {
		needMap[ppcomm.GameTypeNormal] = int64(normalCnt)
	}
	if gameCnt > 0 && dataType == ppcomm.GameTypeGame {
		needMap[ppcomm.GameTypeGame] = int64(gameCnt)
	}
	if mission.SuperBuy > 0 && dataType == ppcomm.GameTypeSuperGame1 {
		needMap[ppcomm.GameTypeSuperGame1] = int64(mission.SuperBuy)
	}
	if mission.SuperBuy2 > 0 && dataType == ppcomm.GameTypeSuperGame2 {
		needMap[ppcomm.GameTypeSuperGame2] = int64(mission.SuperBuy2)
	}
	//查看是否需要结束
	coll2 := db.Collection2(gameId, "simulate")
	isEnough := true
	for t, c := range needMap {
		cnt, _ := coll2.CountDocuments(context.TODO(), bson.M{"type": t})
		if cnt < c {
			isEnough = false
			break
		}
	}
	if isEnough {
		fmt.Printf("the game:%v ,type is %v enough\n", gameId, dataType)
	}
	return isEnough
}
func CheckGameEnoughHGC(gameId string, dataType int) bool {
	var mission GameInfo
	coll := db.Collection2("robot_hgc", "game")
	err := coll.FindOne(context.TODO(), bson.M{"_id": gameId}).Decode(&mission)
	if err == mongo.ErrNoDocuments {
		return false
	}
	var normalCnt = mission.NormalCount // 需要拉取普通盘的次数
	var gameCnt = mission.BuyCount      // 需要拉取游戏的次数
	needMap := map[int]int64{}
	if normalCnt > 0 && dataType == ppcomm.GameTypeNormal {
		needMap[ppcomm.GameTypeNormal] = int64(normalCnt)
	}
	if gameCnt > 0 && dataType == ppcomm.GameTypeGame {
		needMap[ppcomm.GameTypeGame] = int64(gameCnt)
	}
	if mission.SuperBuy > 0 && dataType == ppcomm.GameTypeSuperGame1 {
		needMap[ppcomm.GameTypeSuperGame1] = int64(mission.SuperBuy)
	}
	if mission.SuperBuy2 > 0 && dataType == ppcomm.GameTypeSuperGame2 {
		needMap[ppcomm.GameTypeSuperGame2] = int64(mission.SuperBuy2)
	}
	//查看是否需要结束
	coll2 := db.Collection2(gameId, "simulate")
	isEnough := true
	for t, c := range needMap {
		cnt, _ := coll2.CountDocuments(context.TODO(), bson.M{"type": t})
		if cnt < c {
			isEnough = false
			break
		}
	}
	if isEnough {
		fmt.Printf("the game:%v ,type is %v enough\n", gameId, dataType)
	}
	return isEnough
}

func NewFetcher(game, apiVersion, huiDuSerial string, bonusTy int) *Fetcher {
	ip := getPubIP(http.DefaultClient)
	loc := string(iploc.Country(net.ParseIP(ip)))
	c := http.DefaultClient
	game = strings.TrimPrefix(game, "pp_")
	// game = strings.Split(game, "_")[1]
	if loc != "CN" {
		if mq.NC() == nil {
			mq.ConnectServerMust("127.0.0.1:11002")
		}
		endpoint := lo.Must(ip2worldpub.GetEndpoint())
		urlproxy, _ := url.Parse(fmt.Sprintf("http://%s:%d", endpoint.IP, endpoint.Port))
		c = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlproxy),
			},
			Timeout: 30 * time.Second,
		}
		//c = &http.Client{
		//	//Transport: &http.Transport{
		//	//	Proxy: http.ProxyURL(urlproxy),
		//	//},
		//	Timeout: 30 * time.Second,
		//}
	}
	mgckey, location, err := getHuiDuNewmgcKey(game, huiDuSerial, c)
	if err != nil || location == nil {
		slog.Error("new fetcher,", " err", err, "location", location)
		return nil
	}
	location.Path = fmt.Sprintf("/gs2c/ge/%v/gameService", apiVersion)
	//location.Path = "/gs2c/ge/v3/gameService"
	//location.Path = "/gs2c/ge/v4/gameService"
	fetcher := &Fetcher{
		game:       game,
		httpClient: c,
		mgckey:     mgckey,
		ul:         location.String(),
		bonusTy:    bonusTy,
	}
	return fetcher

}
func NewFetcherHGC(game, apiVersion, HotGameClubSerial string, bonusTy int) *Fetcher {
	ip := getPubIP(http.DefaultClient)
	loc := string(iploc.Country(net.ParseIP(ip)))
	c := http.DefaultClient
	game = strings.TrimPrefix(game, "pp_")
	// game = strings.Split(game, "_")[1]
	if loc != "CN" {
		if mq.NC() == nil {
			mq.ConnectServerMust("127.0.0.1:11002")
		}
		endpoint := lo.Must(ip2worldpub.GetEndpoint())
		urlproxy, _ := url.Parse(fmt.Sprintf("http://%s:%d", endpoint.IP, endpoint.Port))
		c = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlproxy),
			},
			Timeout: 30 * time.Second,
		}
		//c = &http.Client{
		//	//Transport: &http.Transport{
		//	//	Proxy: http.ProxyURL(urlproxy),
		//	//},
		//	Timeout: 30 * time.Second,
		//}
	}
	mgckey, location, err := GetHotGameClubNewmgcKey(HotGameClubSerial, c)
	if err != nil || location == nil {
		slog.Error("new fetcher,", " err", err, "location", location)
		return nil
	}
	location.Path = fmt.Sprintf("/gs2c/ge/%v/gameService", apiVersion)
	//location.Path = "/gs2c/ge/v3/gameService"
	//location.Path = "/gs2c/ge/v4/gameService"
	fetcher := &Fetcher{
		game:       game,
		httpClient: c,
		mgckey:     mgckey,
		ul:         location.String(),
		bonusTy:    bonusTy,
	}
	return fetcher

}
func getLocation(dataGameSrc string) string {
	req2, err := http.NewRequest(http.MethodGet, dataGameSrc, nil)
	if err != nil {
		slog.Error(dataGameSrc+" err", "err", err)
		return ""
	}

	// time.Sleep(time.Second)

	req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	resp2, err := http.DefaultTransport.RoundTrip(req2)
	if err != nil {
		slog.Error("http.DefaultTransport.RoundTrip err", "err", err)
		return ""
	}

	io.Copy(os.Stdout, resp2.Body)

	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusFound {
		return ""
	}

	location := resp2.Header.Get("Location")
	//location = "https://demogamesfree.pragmaticplay.net/gs2c/openGame.do?gameSymbol=vs20superlanche&websiteUrl=https%3A%2F%2Fdemogamesfree.pragmaticplay.net&jurisdiction=99&lobby_url=https%3A%2F%2Fwww.pragmaticplay.com%2Fen%2F"
	return strings.ReplaceAll(location, " ", "%20")
}

func getHuiDuNewmgcKey(game, huiDuSerial string, client *http.Client) (mgckey string, location *url.URL, err error) {
	//game = "1583"
	// 定义要请求的 URL
	url2 := "https://www.huidu.io/en/gameapi/%v/"
	url2 = fmt.Sprintf(url2, huiDuSerial)
	//url2 := "https://www.huidu.io/en/gameapi/1430/"

	// Define the URL
	cmd := exec.Command("curl", "-i", "-A", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36", url2)
	// 捕获命令的输出
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// 执行命令
	err = cmd.Run()
	if err != nil {
		log.Fatalf("执行命令失败: %v\n错误输出: %s", err, stderr.String())
	}
	//
	////url2 = fmt.Sprintf(url2, PPGameRelation[game])
	//// 创建一个新的 HTTP 请求
	//req, err := http.NewRequest("GET", url2, nil)
	//if err != nil {
	//	fmt.Println("创建请求失败:", err)
	//	return
	//}
	//
	//// 设置请求头
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	//req.Header.Add("Cookie", "i18n_redirected=en")
	//resp, err := client.Do(req)
	//if err != nil {
	//	fmt.Println("发送请求失败:", err)
	//	return
	//}
	//defer resp.Body.Close()
	// 读取响应体
	//body, err := io.ReadAll()
	//if err != nil {
	//	fmt.Println(string(body))
	//	fmt.Println("读取响应体失败:", err)
	//	return
	//}
	// 解析 HTML 并提取 data-game-src 的值
	dataGameSrc, err := extractIframeSrc(out.String())
	if err != nil {
		//fmt.Println(out.String())
		fmt.Printf("提取 ifreame-src 失败: %s\n", err)
		return
	}
	//获取游戏链接
	mgckey, location = getHuiDuLocation(dataGameSrc, client)
	if mgckey == "" || location == nil {
		fmt.Println("mgckey or location is nil")
		return
	}
	location.RawQuery = ""
	return

}

func GetHotGameClubNewmgcKey(HotGameClubSerial string, client *http.Client) (mgckey string, location *url.URL, err error) {
	// 定义要请求的 URL
	url2 := "https://admin.hotgameclub.com/ctgameclub/manager/home/V2/game/get/%v"
	url2 = fmt.Sprintf(url2, HotGameClubSerial)
	//url2 := "https://admin.hotgameclub.com/ctgameclub/manager/home/V2/game/get/126"

	payload := strings.NewReader(`{}`)

	req, err := http.NewRequest(http.MethodPost, url2, payload)

	if err != nil {
		//fmt.Println(err)
		return
	}
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Add("access-control-allow-origin", "*")
	req.Header.Add("content-type", "application/json;charset=UTF-8")
	req.Header.Add("language", "zh")
	req.Header.Add("origin", "https://hotgameclub.com")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://hotgameclub.com/")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Add("x-token", "")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(body))
	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}
	//fmt.Println(response.Data.GameEncrypt)

	url3 := "https://www.hotgameclub.com/game-test/startGame-demo"

	payload2 := strings.NewReader(fmt.Sprintf("gameEncrypt=%v", response.Data.GameEncrypt))

	req2, err := http.NewRequest(http.MethodPost, url3, payload2)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Add("access-control-allow-origin", "*")
	req.Header.Add("content-type", "application/json;charset=UTF-8")
	req.Header.Add("origin", "https://hotgameclub.com")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://hotgameclub.com/")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	res2, err := client.Do(req2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body2, err := io.ReadAll(res2.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(body2))

	var response2 ResponseGame
	if err = json.Unmarshal(body2, &response2); err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}
	//fmt.Println(response2.Data.Data)
	//获取游戏链接
	mgckey, location = getHotGameClubLocation(response2.Data.Data, client)
	if mgckey == "" || location == nil {
		fmt.Println("mgckey or location is nil")
		return
	}
	location.RawQuery = ""
	return
}
func getNewmgcKey(game string, client *http.Client) (mgckey string, location *url.URL, err error) {

	// 定义要请求的 URL
	//url2 := "https://www.pragmaticplay.com/en/games/monster-superlanche/?gamelang=en&cur=USD#"

	url2 := "https://www.pragmaticplay.com/en/games/%v/?gamelang=en&cur=USD#"
	//
	url2 = fmt.Sprintf(url2, PPGameRelation[game])
	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 设置请求头
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(string(body))
		fmt.Println("读取响应体失败:", err)
		return
	}
	// 解析 HTML 并提取 data-game-src 的值
	dataGameSrc, err := extractDataGameSrc(string(body))
	if err != nil {
		fmt.Println(string(body))
		fmt.Printf("提取 data-game-src 失败: %s\n", err)
		return
	}
	//获取游戏链接
	location2 := getLocation(dataGameSrc)
	//进入游戏
	enterGame(location2)
	location, err = url.Parse(location2)
	if err != nil {
		return
	}
	query := location.Query()
	mgckey = query.Get("mgckey")
	location.RawQuery = ""
	return

}
func enterGame(url string) {
	//fmt.Println(url)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Host", "demogamesfree.pragmaticplay.net")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(body))
}
