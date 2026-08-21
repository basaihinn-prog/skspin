package hacksawcomm

import (
	"cmp"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"serve/comm/plats/pp"
	"strings"

	"serve/comm/mux"
)

func FetchMgckey(game, user string, client *http.Client) (mgckey string, location *url.URL, err error) {
	var urlret struct {
		Url string
	}
	err = mux.HttpInvoke("https://plats.rpgamestest.com/plat/PP/LaunchGame", map[string]any{
		"Game": game,
		"UID":  user,
		"Lang": "en",
	}, &urlret)

	if err != nil {
		return
	}

	// client req the rest
	req, err := http.NewRequest(http.MethodGet, urlret.Url, nil)

	if err != nil {
		return
	}

	fmt.Println(urlret.Url)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	// resp, err := http.DefaultTransport.RoundTrip(req)
	trans := cmp.Or(client.Transport, http.DefaultTransport)
	resp, err := trans.RoundTrip(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		err = errors.New(resp.Status)
		return
	}

	location, err = url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return
	}

	query := location.Query()
	mgckey = query.Get("mgckey")
	location.RawQuery = ""

	return
}
func GetHuiDuNewmgcKey(client *http.Client, id string) (mgckey string, location *url.URL, err error) {
	////game = "1583"
	//// 定义要请求的 URL
	////url2 := "https://www.huidu.io/en/gameapi/%v/"
	////url2 = fmt.Sprintf(url2, huiDuSerial)
	//endpoint := lo.Must(ip2worldpub.GetEndpoint())
	//u := fmt.Sprintf("http://%s:%d", endpoint.IP, endpoint.Port)
	//urlproxy, _ := url.Parse(u)
	//
	//// 创建自定义的 Transport，设置代理
	//transport := &http.Transport{
	//	Proxy: http.ProxyURL(urlproxy),
	//}
	//
	//// 创建 HTTP 客户端
	//client.Transport = transport
	//fmt.Println(transport)
	//url2 := "https://www.huidu.io/en/gameapi/135/"
	////url2 := "https://api.jdbgaming.com/game/demo?lang=zh-CN&amp;id=13"
	//
	//// Define the URL
	//cmd := exec.Command("curl",
	//	"-i",                                                                                                                    // 包含响应头
	//	"-A", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36", // 设置 User-Agent
	//	url2,
	//)
	//// 捕获命令的输出
	//var out bytes.Buffer
	//var stderr bytes.Buffer
	//cmd.Stdout = &out
	//cmd.Stderr = &stderr
	//
	//// 执行命令
	//err = cmd.Run()
	//if err != nil {
	//	log.Fatalf("执行命令失败: %v\n错误输出: %s", err, stderr.String())
	//}
	//
	////url2 = fmt.Sprintf(url2, PPGameRelation[game])
	//// 创建一个新的 HTTP 请求
	////req, err := http.NewRequest("GET", url2, nil)
	////if err != nil {
	////	fmt.Println("创建请求失败:", err)
	////	return
	////}
	//
	////// 设置请求头
	////req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	//////req.Header.Add("Cookie", "i18n_redirected=en")
	////resp, err := client.Do(req)
	////if err != nil {
	////	fmt.Println("发送请求失败:", err)
	////	return
	////}
	////defer resp.Body.Close()
	//////读取响应体
	////body, err := io.ReadAll(resp.Body)
	////fmt.Println(string(body))
	////if err != nil {
	////	fmt.Println("读取响应体失败:", err)
	////	return
	////}
	//// 解析 HTML 并提取 data-game-src 的值
	//fmt.Println(out.String())
	//dataGameSrc, err := extractIframeSrc(out.String())
	//if err != nil {
	//	//fmt.Println(out.String())
	//	fmt.Printf("提取 ifreame-src 失败: %s\n", err)
	//	return
	//}
	//https://api.jdbgaming.com/game/demo?lang=en-US&id=94
	//获取游戏链接
	//mgckey, location = getHuiDuLocation(dataGameSrc, client)
	gameID := fmt.Sprintf("https://api.jdbgaming.com/game/demo?lang=en-US&id=%s", id)
	mgckey, location = getHuiDuLocation(gameID, client)
	if mgckey == "" || location == nil {
		fmt.Println("mgckey or location is nil")
		return
	}
	location.RawQuery = ""
	return

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
	if resp2.StatusCode != http.StatusSeeOther {
		return "", nil
	}

	location := resp2.Header.Get("Location")
	//location2, err := GetRedirect(location, client)
	//if err != nil {
	//	err = errors.New("重定向失败")
	//	return "", nil
	//}
	location2, err := url.Parse(location)
	if err != nil {
		return "", location2
	}
	query := location2.Query()
	mgckey := query.Get("x")

	return mgckey, location2
}

// 提取 src 的值
func extractIframeSrc(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}
	fmt.Println(doc)
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
func GetMgckey(game, user string, client *http.Client) (mgckey string, location *url.URL, err error) {
	// client req the rest
	uid, guid, cash := "ceshi1938", "4ae52ed2e1a8c353878ba65ed7791ac4", "5000000"
	uid = user
	guid = GuidMap[game]
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
