package ut

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
)

//bbgt游戏平台操作工具

var UserList = []string{
	user1,
	user2,
	user3,
	user4,
	user5,
	user6,
	user7,
	user8,
}

var EXUserList = []string{
	user9,
	user10,
	user11,
	user12,
	user13,
	user14,
	user15,
	user16,
}

var BFUserList = []string{
	user17,
	user18,
	user19,
	user20,
	user21,
	user22,
	user23,
	user24,
}

const (
	user1  = "uid-m2mocz8a1qn52q84v"
	user2  = "uid-m2mocz8a1qn52q8vb"
	user3  = "uid-1uycl0l4lo9vp1585"
	user4  = "uid-hik09ecikqgpxdmya"
	user5  = "uid-y6x17oc9vgc4i155m"
	user6  = "uid-ibipn0g7q562rmo6t"
	user7  = "uid-84fswtb2xfexmu2qo"
	user8  = "uid-lvivetkqwy79gwpby"
	user9  = "uid-j932er9k5kc28wqyb"
	user10 = "uid-7xcgz2emwrquk9rr4"
	user11 = "uid-sn577q25dnyrp75d4"
	user12 = "uid-64ajlpin6gmolp85i"
	user13 = "uid-xg76nvvgxjm4buf8v"
	user14 = "uid-dpjvd237nnjje1er0"
	user15 = "uid-rs15ks8tpb363b0e4"
	user16 = "uid-6ek0gkbkjlafg5whr"
	user17 = "uid-3adz5v3d3ptsfvz2z"
	user18 = "uid-woj3l7psj3nruovo6"
	user19 = "uid-053kcr2ur491ql76y"
	user20 = "uid-ru0pm6zba1vg49e69"
	user21 = "uid-1t4wxkv7qbc4rlk8g"
	user22 = "uid-0hqr77wh2teaegitg"
	user23 = "uid-1dcumyc1v3f7iuulm"
	user24 = "uid-ogvw9eghbuhoht8gg"
)

type BbgtRegisterBody struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	Confirmpassword string `json:"confirm_password"`
	Prefix          string `json:"prefix"`
	Phone           string `json:"phone"`
	Type            string `json:"type,omitempty"`
}

// BbgtRegister 注册 入参只给username即可，不给就随机默认值
func BbgtRegister(info *BbgtRegisterBody) (string, error) {
	url := "https://bbgtgame.com/gameapi/api/register"
	info.Password = "123456"
	info.Confirmpassword = "123456"
	info.Prefix = "55"
	info.Phone = ""
	if info.Username == "" {
		info.Username = fmt.Sprintf("uid-%v", GenerateRandomString(17))
	}
	marshal, _ := json.Marshal(info)
	data := string(marshal)
	// 构建 curl 命令
	cmd := exec.Command("curl", "-s", url,
		"-H", "accept: */*",
		"-H", "accept-language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"-H", "app-token:",
		"-H", "cache-control: no-cache",
		"-H", "content-type: application/json",
		"-H", "origin: https://bbgtgame.com",
		"-H", "pragma: no-cache",
		"-H", "priority: u=1, i",
		"-H", "referer: https://bbgtgame.com/",
		"-H", `sec-ch-ua: "Chromium";v="130", "Microsoft Edge";v="130", "Not?A_Brand";v="99"`,
		"-H", "sec-ch-ua-mobile: ?0",
		"-H", `sec-ch-ua-platform: "Windows"`,
		"-H", "sec-fetch-dest: empty",
		"-H", "sec-fetch-mode: cors",
		"-H", "sec-fetch-site: same-origin",
		"-H", "think-lang: pt",
		"-H", "user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.0.0",
		"--data-raw", data,
	)
	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing curl:", err)
		return "", err
	}
	mmmap := make(map[string]interface{})
	json.Unmarshal(output, &mmmap)
	appToken := ""
	if _, ok := mmmap["data"]; ok {
		da := mmmap["data"].(map[string]interface{})
		appToken = da["token"].(string)
	} else {
		slog.Error("获取token失败")
		return appToken, err
	}

	return appToken, nil
}

// BbgtLogin 登陆 入参info只给username即可，必传
func BbgtLogin(info *BbgtRegisterBody, appToken string) (string, error) {
	url := "https://bbgtgame.com/gameapi/api/login"
	info.Password = "123456"
	info.Confirmpassword = "123456"
	info.Prefix = "55"
	info.Phone = ""
	info.Type = "1"
	if info.Username == "" {
		slog.Error("username必传")
		return "", errors.New("登陆失败，username必传")
	}
	marshal, _ := json.Marshal(info)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshal))
	if err != nil {
		slog.Error(fmt.Sprintf("请求登陆失败 err: %v", err))
		return "", err
	}
	// 构建 curl 命令
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	if appToken != "" {
		req.Header.Set("App-Token", appToken)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://bbgtgame.com")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://bbgtgame.com/")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Think-Lang", "pt")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error(fmt.Sprintf("请求登陆失败 err: %v", err))
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	mmmap := make(map[string]interface{})
	json.Unmarshal(body, &mmmap)
	if _, ok := mmmap["data"]; ok {
		da := mmmap["data"].(map[string]interface{})
		appToken = da["token"].(string)
	} else {
		slog.Error("获取token失败")
		return appToken, err
	}

	return appToken, nil
}

// BbgtSetBalance 给用户设置余额
func BbgtSetBalance(uid, balance, appToken string) error {
	// 基础 URL
	baseURL := "https://bbgtgame.com/gameapi/api/games/wwwSetUserBalance"
	// 创建 URL 对象
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		slog.Error(fmt.Sprintf("生成url失败 err: %v", err))
		return err
	}
	// 设置查询参数
	query := reqURL.Query()
	query.Set("uid", uid)
	query.Set("balance", balance)
	reqURL.RawQuery = query.Encode() // 生成查询字符串

	// 创建新的请求
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		slog.Error(fmt.Sprintf("构建请求失败 err: %v", err))
		return err
	}

	// 设置请求头
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("App-Token", appToken)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://bbgtgame.com/")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Think-Lang", "pt")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error(fmt.Sprintf("请求失败 err: %v", err))
		return err
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	mmmap := make(map[string]interface{})
	json.Unmarshal(body, &mmmap)
	if mmmap["code"].(float64) != 1 {
		slog.Error("修改余额失败 err: ", mmmap["data"])
		return err
	} else {
		return nil
	}
}

// BbgtGetUserInfo 获取游客信息
func BbgtGetUserInfo(appToken string) (map[string]interface{}, error) {
	// 定义请求 URL
	url := "https://bbgtgame.com/gameapi/api/info"

	// 创建请求体
	jsonData := []byte("{}") // 使用空 JSON 对象

	// 创建新的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error(fmt.Sprintf("构建请求失败 err: %v", err))
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("App-Token", appToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://bbgtgame.com")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://bbgtgame.com/")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Think-Lang", "pt")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error(fmt.Sprintf("请求失败 err: %v", err))
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	mmmap := make(map[string]interface{})
	json.Unmarshal(body, &mmmap)
	return mmmap, nil
}

// GetNewUser 注册新用户并给20000钱
func GetNewUser() (string, error) {
	username := ""
	//1.注册
	userInfo := &BbgtRegisterBody{}
	token, err := BbgtRegister(userInfo)
	if err != nil {
		slog.Error(fmt.Sprintf("注册用户失败，err: %v", err))
		return "", err
	}
	//2.查询信息获取uid
	info, err := BbgtGetUserInfo(token)
	if err != nil {
		slog.Error(fmt.Sprintf("获取用户信息失败，err: %v", err))
		return "", err
	}
	data := info["data"].(map[string]interface{})
	tempInfo := data["info"].(map[string]interface{})
	username = tempInfo["username"].(string)
	//3.给钱
	uid := strconv.FormatFloat(tempInfo["id"].(float64), 'f', 4, 64)
	err = BbgtSetBalance(uid, "20000", token)
	if err != nil {
		slog.Error(fmt.Sprintf("给钱失败，err: %v", err))
		return "", err
	}
	return username, nil
}

// UserLogin 用户登陆，如果不足一定额度的钱就新加20000钱
func UserLogin(uid string, cash float64) (string, error) {
	//登陆
	userInfo := &BbgtRegisterBody{
		Username: uid,
	}
	token, err := BbgtLogin(userInfo, "")
	if err != nil {
		slog.Error(fmt.Sprintf("用户登陆失败，err: %v", err))
		return "", err
	}
	//2.查询信息获取uid
	info, err := BbgtGetUserInfo(token)
	if err != nil {
		slog.Error(fmt.Sprintf("获取用户信息失败，err: %v", err))
		return "", err
	}
	data := info["data"].(map[string]interface{})
	tempInfo := data["info"].(map[string]interface{})
	uid = strconv.FormatFloat(tempInfo["id"].(float64), 'f', 4, 64)
	balance, err := strconv.ParseFloat(tempInfo["balance"].(string), 64)
	if err != nil {
		slog.Error(fmt.Sprintf("转换失败，err: %v", err))
		return "", err
	}
	if balance < cash {
		//3.给钱
		err = BbgtSetBalance(uid, "20000", token)
		if err != nil {
			slog.Error(fmt.Sprintf("给钱失败，err: %v", err))
			return "", err
		}
	}

	return token, nil
}

func ChargeCash(token, cash string) error {
	//3.给钱
	//2.查询信息获取uid
	info, err := BbgtGetUserInfo(token)
	if err != nil {
		slog.Error(fmt.Sprintf("获取用户信息失败，err: %v", err))
		return err
	}
	data := info["data"].(map[string]interface{})
	tempInfo := data["info"].(map[string]interface{})
	uid := strconv.FormatFloat(tempInfo["id"].(float64), 'f', 4, 64)
	err = BbgtSetBalance(uid, cash, token)
	if err != nil {
		slog.Error(fmt.Sprintf("给钱失败，err: %v", err))
		return err
	}
	return nil
}
