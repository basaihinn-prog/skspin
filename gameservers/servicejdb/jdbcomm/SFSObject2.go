package jdbcomm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ResponseAction101 struct {
	Status string `json:"status"`
	Data   struct {
		Ots string `json:"ots"`
	} `json:"data"`
}
type UserInfo struct {
	UID        string `json:"uid"`
	UserName   string `json:"userName"`
	Lvl        int    `json:"lvl"`
	UserStatus int    `json:"userStatus"`
	Currency   string `json:"currency"`
}

type Url struct {
	Url string `json:"url"`
}

type ResponseAction6 struct {
	Status string     `json:"status"`
	Data   []UserInfo `json:"data"`
}

type ResponseAction20 struct {
	Status string `json:"status"`
	Data   Url    `json:"data"`
}

type ResponseAction19 struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	IsShowAutoPlay bool     `json:"isShowAutoPlay"`
	Result4        Result4  `json:"result4"`
	Result6        Result6  `json:"result6"`
	Result10       Result10 `json:"result10"`
}

type Result4 struct {
	Currency         string `json:"currency"`
	IsDemoAccount    bool   `json:"isDemoAccount"`
	IsApiAccount     bool   `json:"isApiAccount"`
	IsShowJackpot    bool   `json:"isShowJackpot"`
	IsShowCurrency   bool   `json:"isShowCurrency"`
	IsShowDollarSign bool   `json:"isShowDollarSign"`
	DecimalPoint     int    `json:"decimalPoint"`
	GameGroup        []int  `json:"gameGroup"`
	FunctionList     []int  `json:"functionList"`
}

type Result6 struct {
	UID        string `json:"uid"`
	UserName   string `json:"userName"`
	Lvl        int    `json:"lvl"`
	UserStatus int    `json:"userStatus"`
	Currency   string `json:"currency"`
}

type Result10 struct {
	Status               string   `json:"status"`
	SessionID            []string `json:"sessionID"`
	Zone                 string   `json:"zone"`
	GsInfo               string   `json:"gsInfo"`
	GameType             int      `json:"gameType"`
	MachineType          int      `json:"machineType"`
	IsRecovery           bool     `json:"isRecovery"`
	S0                   string   `json:"s0"`
	S1                   string   `json:"s1"`
	S2                   string   `json:"s2"`
	S3                   string   `json:"s3"`
	S4                   string   `json:"s4"`
	GameUid              string   `json:"gameUid"`
	GamePass             string   `json:"gamePass"`
	UseSSL               bool     `json:"useSSL"`
	StreamingUrl         struct{} `json:"streamingUrl"`
	AchievementServerUrl string   `json:"achievementServerUrl"`
	ChatServerUrl        string   `json:"chatServerUrl"`
	IsWSBinary           bool     `json:"isWSBinary"`
}

// getTKLocation 获取 tk 和 location
func GetTKLocation() (string, *url.URL, error) {
	tk, location, err := GetHuiDuNewmgcKey(&http.Client{}, "")
	if err != nil {
		return "", nil, err
	}
	return tk, location, nil
}

// sendRequest 发送请求
func sendRequest(method, url string, payload io.Reader, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return client.Do(req)
}

// Action101 发送 Action 101 请求
func Action101(tk string) (*ResponseAction101, error) {
	url2 := "https://eweb03.js-mingyi.com/frontendAPI.do"
	method := "POST"
	data := url.Values{}
	data.Set("action", "101")
	data.Set("x", tk)
	payload := strings.NewReader(data.Encode())
	headers := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Connection":         "keep-alive",
		"Content-Type":       "application/x-www-form-urlencoded",
		"Origin":             "https://jifjie.h9vo10aqz.com",
		"Referer":            "https://jifjie.h9vo10aqz.com/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "cross-site",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
		"jots":               "892e6599-c1ff-497d-aed9-59ab5e2cd956",
		"sec-ch-ua":          "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
	}
	res, err := sendRequest(method, url2, payload, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp ResponseAction101
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Action6 发送 Action 6 请求
func Action6(tk string) (*ResponseAction6, error) {
	url2 := "https://eweb03.js-mingyi.com/frontendAPI.do"
	method := "POST"
	data := url.Values{}
	data.Set("action", "6")
	data.Set("x", tk)
	payload := strings.NewReader(data.Encode())
	headers := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Connection":         "keep-alive",
		"Content-Type":       "application/x-www-form-urlencoded",
		"Origin":             "https://jifjie.h9vo10aqz.com",
		"Referer":            "https://jifjie.h9vo10aqz.com/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "cross-site",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
		"jots":               "892e6599-c1ff-497d-aed9-59ab5e2cd956",
		"sec-ch-ua":          "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
	}
	res, err := sendRequest(method, url2, payload, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp ResponseAction6
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Action19 发送 Action 19 请求
func Action19(tk string) (*ResponseAction19, error) {
	url2 := "https://eweb03.js-mingyi.com/frontendAPI.do"
	method := "POST"
	data := url.Values{}
	data.Set("action", "19")
	data.Set("x", tk)
	data.Set("gameType", "14")
	data.Set("mType", "14027")
	data.Set("gName", "LuckySeven_f11946f")
	data.Set("clientType", "web")
	data.Set("gameLine", "est03.js-mingyi.com_443_0")
	payload := strings.NewReader(data.Encode())
	headers := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Connection":         "keep-alive",
		"Content-Type":       "application/x-www-form-urlencoded",
		"Origin":             "https://jifjie.h9vo10aqz.com",
		"Referer":            "https://jifjie.h9vo10aqz.com/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "cross-site",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
		"jots":               "892e6599-c1ff-497d-aed9-59ab5e2cd956",
		"sec-ch-ua":          "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
	}
	res, err := sendRequest(method, url2, payload, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp ResponseAction19
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Action23 发送 Action 23 请求
func Action23(tk string) (string, error) {
	url2 := "https://eweb03.js-mingyi.com/frontendAPI.do"
	method := "POST"
	data := url.Values{}
	data.Set("action", "23")
	data.Set("x", tk)
	data.Set("mType", "14027")
	payload := strings.NewReader(data.Encode())
	headers := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Connection":         "keep-alive",
		"Content-Type":       "application/x-www-form-urlencoded",
		"Origin":             "https://jifjie.h9vo10aqz.com",
		"Referer":            "https://jifjie.h9vo10aqz.com/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "cross-site",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
		"jots":               "892e6599-c1ff-497d-aed9-59ab5e2cd956",
		"sec-ch-ua":          "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
	}
	res, err := sendRequest(method, url2, payload, headers)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Action24 发送 Action 24 请求
func Action24(tk string) (string, error) {
	url2 := "https://eweb03.js-mingyi.com/frontendAPI.do"
	method := "POST"
	data := url.Values{}
	data.Set("action", "24")
	data.Set("x", tk)
	data.Set("mType", "14027")
	payload := strings.NewReader(data.Encode())
	headers := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Connection":         "keep-alive",
		"Content-Type":       "application/x-www-form-urlencoded",
		"Origin":             "https://jifjie.h9vo10aqz.com",
		"Referer":            "https://jifjie.h9vo10aqz.com/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "cross-site",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
		"jots":               "892e6599-c1ff-497d-aed9-59ab5e2cd956",
		"sec-ch-ua":          "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
	}
	res, err := sendRequest(method, url2, payload, headers)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Action13 发送 Action 13 请求
func Action13(tk string) (string, error) {
	url2 := "https://eweb03.js-mingyi.com/frontendAPI.do"
	method := "POST"
	str := `{"level":"GAME","event":"EGRET_READY","message":{"message":"GAME IS READY."},"accessToken":"%s","apiServer":"https://eweb03.js-mingyi.com","gName":"LuckySeven_f11946f","gameType":"14","mType":"14027","ui_version":"3.143.1","uniqueKey":"%s","userAgent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36","userName":"demo001652@XX"}`
	data := url.Values{}
	data.Set("action", "13")
	data.Set("x", tk)
	data.Set("data", fmt.Sprintf(str, tk, "1741146017009@eb98bf56-ff5f-4dc8-896a-1920cb612498 demo001652@XX"))
	payload := strings.NewReader(data.Encode())
	headers := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Connection":         "keep-alive",
		"Content-Type":       "application/x-www-form-urlencoded",
		"Origin":             "https://jifjie.h9vo10aqz.com",
		"Referer":            "https://jifjie.h9vo10aqz.com/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "cross-site",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
		"jots":               "892e6599-c1ff-497d-aed9-59ab5e2cd956",
		"sec-ch-ua":          "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
	}
	res, err := sendRequest(method, url2, payload, headers)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// HeartBeat 发送心跳请求
func HeartBeat(tk string) {
	for {
		_, err := Action24(tk)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Millisecond * 15000)
	}
}
