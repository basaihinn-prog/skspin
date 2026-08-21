package jiliDemo

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"log/slog"
	"net/http"
	"net/url"
	"serve/comm/ut"
	"strings"
	"time"
)

var agentId = "Royalebet_RaniBaji_INR"
var agentKey = "4d33f06a363e33234601368325f7b8bfaa63e8b4"

const (
	apiUrl1 = `https://uat-wb-api-2.jiscc88.com/api1/`
	apiUrl2 = `https://uat-wb-api-2.kijl788du.com/api1/`
	apiUrl3 = `https://uat-wb-api-2.huuykk865s.com/api1/`
	apiUrl4 = `https://uat-wb-api-2.jiscc77s.com/api1/`
	apiUrl5 = `https://uat-wb-api-2.jismk2u.com/api1/`
)

func GetDemoUrl(name string, gameId int) (string, error) {
	var err error
	rspUrl := ""
	api := `LoginWithoutRedirect`
	apiUrl := apiUrl1 + api
	if name == "" {
		name = `xiaoxiang`
	}
	if gameId == 0 {
		gameId = 1
	}
	param := fmt.Sprintf(`Account=%s&GameId=%d&Lang=%s&`, name, gameId, `zh-CN`)
	reqstr := fmt.Sprintf("%s&Key=%s", GetQueryString(param), GetKey(GetQueryString(param), GetKeyG()))
	reqBody := strings.NewReader(reqstr)
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, apiUrl, reqBody)
	if err != nil {
		slog.Error("http.NewRequest err", "err", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	// http.DefaultClient.Do(req)
	body, _ := lo.Must2(ut.DoHttpReq(httpClient, req))
	slog.Info("rsp", "body", string(body))

	temp := make(map[string]interface{})
	json.Unmarshal(body, &temp)
	if _, ok := temp["Data"]; ok {
		rspUrl = temp["Data"].(string)
	}
	slog.Info("rsp", "rspUrl", rspUrl)
	//	ssoKey := rspUrl[]
	//
	//https: //uat-wbwebapi.jlfafafa2.com/sso-login.api?key=07dadd5f7b2c5b35b1a3a7a8e9e9625f7e72ef3b&lang=zh-CN(ssoKey)
	//
	//	`token:"7495aad47d8656e1b9c2e048c07abae572dd71e7"`
	//
	//	req0 -> `token:"7495aad47d8656e1b9c2e048c07abae572dd71e7"`
	//	rsp
	return rspUrl, nil
}

func GetGameList() (string, error) {
	var err error
	rspUrl := ""
	apistr := `GetGameList`
	apiUrl := apiUrl1 + apistr

	reqstr := fmt.Sprintf("%s&Key=%s", GetQueryString(""), GetKey(GetQueryString(""), GetKeyG()))
	reqBody := strings.NewReader(reqstr)
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, apiUrl, reqBody)
	if err != nil {
		slog.Error("http.NewRequest err", "err", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	// http.DefaultClient.Do(req)
	body, _ := lo.Must2(ut.DoHttpReq(httpClient, req))
	slog.Info("rsp", "body", string(body))

	jilirsp, err := url.Parse(string(body))
	if err != nil {
		return "", err
	}
	query := jilirsp.Query()
	rspUrl = query.Get("")

	return rspUrl, nil
}

func GetAccount(accountName string) (string, error) {
	var err error
	rspUrl := ""
	api := `CreateMember`
	apiUrl := apiUrl1 + api
	reqstr := fmt.Sprintf("%s&Key=%s", GetQueryString(accountName), GetKey(GetQueryString(accountName), GetKeyG()))
	reqBody := strings.NewReader(reqstr)
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, apiUrl, reqBody)
	if err != nil {
		slog.Error("http.NewRequest err", "err", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	// http.DefaultClient.Do(req)
	body, _ := lo.Must2(ut.DoHttpReq(httpClient, req))
	slog.Info("rsp", "body", string(body))

	jilirsp, err := url.Parse(string(body))
	if err != nil {
		return "", err
	}
	query := jilirsp.Query()
	rspUrl = query.Get("")

	return rspUrl, nil
}

func ExchangeJili(account string, amount float64) (string, error) {
	var err error
	rspUrl := ""
	api := `ExchangeTransferByAgentId`
	apiUrl := apiUrl1 + api
	param := fmt.Sprintf(`Account=%s&TransactionId=%s&Amount=%v&TransferType=%d&`, account, `17`, amount, 2)
	reqstr := fmt.Sprintf("%s&Key=%s", GetQueryString(param), GetKey(GetQueryString(param), GetKeyG()))
	reqBody := strings.NewReader(reqstr)
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, apiUrl, reqBody)
	if err != nil {
		slog.Error("http.NewRequest err", "err", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	// http.DefaultClient.Do(req)
	body, _ := lo.Must2(ut.DoHttpReq(httpClient, req))
	slog.Info("rsp", "body", string(body))

	jilirsp, err := url.Parse(string(body))
	if err != nil {
		return "", err
	}
	query := jilirsp.Query()
	rspUrl = query.Get("")

	return rspUrl, nil
}

func GetQueryString(querystring string) string {
	//string querystring = "Account=Test006&GameId=1&Lang=zh-CN&AgentId=abcd_RMB";
	querystring += fmt.Sprintf(`AgentId=%s`, agentId)
	return querystring
}

func GetKeyG() string {
	now := FormatUTCMinus4Time()
	keyG := GetMD5String(now + agentId + agentKey)
	return keyG
}

func GetKey(querystring, KeyG string) string {
	randomText1 := ut.GenerateRandomString2(6)
	randomText2 := ut.GenerateRandomString2(6)
	return randomText1 + GetMD5String(querystring+KeyG) + randomText2
}

func GetMD5String(str string) string {
	h := md5.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		slog.Error("GetMD5String", "err", err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func FormatUTCMinus4Time() string {
	// 1. 获取当前 UTC 时间
	utcNow := time.Now().UTC()

	// 2. 转换为 UTC-4 时区
	utcMinus4 := utcNow.Add(-4 * time.Hour)

	// 3. 提取年月日
	year := utcMinus4.Year() % 100 // 取后两位
	month := int(utcMinus4.Month())
	day := utcMinus4.Day()

	// 4. 格式化
	// 年份：2位，不足补零（如 2024 → 24）
	yy := fmt.Sprintf("%02d", year)
	// 月份：2位，不足补零（如 3 → 03）
	MM := fmt.Sprintf("%02d", month)
	// 日：1-9 不补零，10+ 保持不变
	d := fmt.Sprintf("%d", day)

	// 组合成 yyMMd 格式
	return yy + MM + d
}
