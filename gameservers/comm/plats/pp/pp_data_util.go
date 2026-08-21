package pp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	u "net/url"
	"strconv"
	"time"
)

var (
	aes_key    = "8325734e945b4361b175274d6263d996"
	agency_uid = "720eed0a990c252782252af6a5a4b5dd"
)

type Payload struct {
	AgencyUID     string `json:"agency_uid"`
	MemberAccount string `json:"member_account"`
	GameUID       string `json:"game_uid"`
	Timestamp     string `json:"timestamp"`
	CreditAmount  string `json:"credit_amount"`
	CurrencyCode  string `json:"currency_code"`
}

func NewPayload(uid, guid, cash string) Payload {
	timestampStr := strconv.FormatInt(time.Now().UnixMilli(), 10)
	return Payload{
		AgencyUID:     agency_uid,
		MemberAccount: uid,
		GameUID:       guid,
		Timestamp:     timestampStr,
		CreditAmount:  cash,
		CurrencyCode:  "AED",
	}
}

func GetKey(uid, guid, cash string) (string, error) {
	url := ""
	//uid, guid, cash := "ceshi1937", "4ae52ed2e1a8c353878ba65ed7791ac4", "5000000"
	payload := NewPayload(uid, guid, cash)
	marshal, _ := json.Marshal(payload)
	epayload, _ := Encrypt(marshal, []byte(aes_key))
	reqstr := `{"agency_uid": "720eed0a990c252782252af6a5a4b5dd","timestamp": "1728721656324","payload": "puyntgzJZKHbB63TotTeUZsMzCtTvYdcaZWfogRRF1DgZ+A30/681vgyZ/TzUw1DBwrktQUqWsoI3XRnpO4bfBkOdCQnbnn4wFGbhKlLlNm2GBjQZWndbZk4wNikUFxYp0u5eO4xvsS2UQ0mL48nqOx8a5zrXT4KkZXFMTGUECPm/PCYVPMiD59Eo/Kp+1ChNAd8rK88KSKoFQpE4P/Gt8JjATl5w4qdGo7NNV/pKKMTxAYIssCAN7zGRgG1K8yeHDDEPqyEFr+VIMCDpIprq8wZyRHc5e3xNXrrEE+Ki0YiYaa+ZZF+aoYamqaCPHuJ3S9NYsqIjakGH78RgPIHQA=="}`
	var data map[string]string
	err := json.Unmarshal([]byte(reqstr), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return url, err
	}
	data["timestamp"] = payload.Timestamp
	data["payload"] = epayload
	var ps = u.Values{}
	// 将 map[string]string 转换为 url.Values
	for k, v := range data {
		ps.Set(k, v)
	}
	jsonData, err := json.Marshal(data)
	req, err := http.NewRequest("POST", "https://jsgame.live/game/v1", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return url, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "PostmanRuntime/7.28.4")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return url, err
	}
	defer resp.Body.Close()
	rsp := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&rsp)
	if err != nil {
		slog.Error("Error:", err)
		return "", err
	}
	slog.Info("jsgame.live/game/v1 rsp: ", "data", data, "resp", rsp)
	if rsp != nil && rsp["payload"] != nil {
		url = rsp["payload"].(map[string]interface{})["game_launch_url"].(string)
	}
	return url, nil
}

//func GetMgckey(game, user string, client *http.Client) (mgckey string, location *url.URL, err error) {
//	uid, guid, cash := "ceshi1937", "4ae52ed2e1a8c353878ba65ed7791ac4", "5000000"
//	payload := NewPayload(uid, guid, cash)
//	marshal, _ := json.Marshal(payload)
//	epayload, _ := Encrypt(marshal, []byte(aes_key))
//	reqstr := `{"agency_uid": "810f7f8748d39cdb8265fb95fa0ad462","timestamp": "1728721656324","payload": "aMUncI1FOKvh4O75N+rOLnkUx1CVdiOMc8FHTAtWpmu7G41P9vqwzTNfnPw34GVwRgIJYeiRZVjFCvV6Wbj/WzNv7PeTGMcXr/w5almnal86lP0v/tK3MnKX6VUrqZAkD3NYuDNJE4ePAd1LMNLPJEHk2BuQE/i25NahBToBu78dGf1I0EEWzhSy6z3+k2Of9vzyQugFY0Sr42gzVhVbZzdPW1/8uJBCHx8xCef+J3onQhVRY9bQssOFb1Dt8Je3ljemiyxbi4l9KaEY4Iok5A=="}`
//	var data map[string]string
//	err = json.Unmarshal([]byte(reqstr), &data)
//	if err != nil {
//		fmt.Println("Error:", err)
//		//return url, err
//	}
//	data["timestamp"] = payload.Timestamp
//	data["payload"] = epayload
//	var ps = u.Values{}
//	// 将 map[string]string 转换为 url.Values
//	for k, v := range data {
//		ps.Set(k, v)
//	}
//	jsonData, err := json.Marshal(data)
//	req, err := http.NewRequest("POST", "https://jsgame.live/game/v1", bytes.NewBuffer(jsonData))
//	if err != nil {
//		fmt.Println("Error:", err)
//		//return url, err
//	}
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Accept", "*/*")
//	req.Header.Set("User-Agent", "PostmanRuntime/7.28.4")
//	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
//	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
//	req.Header.Set("Connection", "keep-alive")
//	req.Header.Set("Cache-Control", "no-cache")
//	req.Header.Set("Pragma", "no-cache")
//	// resp, err := http.DefaultTransport.RoundTrip(req)
//	trans := cmp.Or(client.Transport, http.DefaultTransport)
//	resp, err := trans.RoundTrip(req)
//	if err != nil {
//		return
//	}
//
//	defer resp.Body.Close()
//	if resp.StatusCode != http.StatusFound {
//		err = errors.New(resp.Status)
//		return
//	}
//
//	location, err = url.Parse(resp.Header.Get("Location"))
//	if err != nil {
//		return
//	}
//
//	query := location.Query()
//	mgckey = query.Get("mgckey")
//	location.RawQuery = ""
//
//	return
//}

func GetReq() (*http.Request, error) {
	uid, guid, cash := "ceshi1937", "4ae52ed2e1a8c353878ba65ed7791ac4", "5000000"
	payload := NewPayload(uid, guid, cash)
	marshal, _ := json.Marshal(payload)
	epayload, _ := Encrypt(marshal, []byte(aes_key))
	reqstr := `{"agency_uid": "810f7f8748d39cdb8265fb95fa0ad462","timestamp": "1728721656324","payload": "aMUncI1FOKvh4O75N+rOLnkUx1CVdiOMc8FHTAtWpmu7G41P9vqwzTNfnPw34GVwRgIJYeiRZVjFCvV6Wbj/WzNv7PeTGMcXr/w5almnal86lP0v/tK3MnKX6VUrqZAkD3NYuDNJE4ePAd1LMNLPJEHk2BuQE/i25NahBToBu78dGf1I0EEWzhSy6z3+k2Of9vzyQugFY0Sr42gzVhVbZzdPW1/8uJBCHx8xCef+J3onQhVRY9bQssOFb1Dt8Je3ljemiyxbi4l9KaEY4Iok5A=="}`
	var data map[string]string
	err := json.Unmarshal([]byte(reqstr), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	data["timestamp"] = payload.Timestamp
	data["payload"] = epayload
	var ps = u.Values{}
	// 将 map[string]string 转换为 url.Values
	for k, v := range data {
		ps.Set(k, v)
	}
	jsonData, err := json.Marshal(data)
	req, err := http.NewRequest("POST", "https://jsgame.live/game/v1", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "PostmanRuntime/7.28.4")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	return req, nil
}
