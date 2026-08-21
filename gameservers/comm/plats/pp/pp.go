package pp

import (
	"log/slog"
	"serve/comm/plats/platcomm"
)

func init() {
	platcomm.Plats["PP"] = pp{}
}

type pp struct{}

// {"UID":"123456","Game":"vs20olympx","Lang":"en"}
func (p pp) LaunchGame(username, game, lang string, useProxy bool) (url string, err error) {
	//err = p.Regist(username)
	//if err != nil {
	//	return
	//}
	//balance, _ := p.GetBalance(username)
	//if balance < 10000 {
	//	p.FundTransferIn(username, 100e4)
	//}
	//
	//var ret struct {
	//	GameURL string
	//}
	//err = invoke("/game/start/", map[string]string{
	//	"externalPlayerId": username,
	//	"gameId":           game,
	//	"language":         lang,
	//	"platform":         "MOBILE",
	//}, &ret)
	//url = ret.GameURL
	//url := ""

	//reqstr := `{"agency_uid": "810f7f8748d39cdb8265fb95fa0ad462","timestamp": "1728721656324","payload": "aMUncI1FOKvh4O75N+rOLnkUx1CVdiOMc8FHTAtWpmu7G41P9vqwzTNfnPw34GVwRgIJYeiRZVjFCvV6Wbj/WzNv7PeTGMcXr/w5almnal86lP0v/tK3MnKX6VUrqZAkD3NYuDNJE4ePAd1LMNLPJEHk2BuQE/i25NahBToBu78dGf1I0EEWzhSy6z3+k2Of9vzyQugFY0Sr42gzVhVbZzdPW1/8uJBCHx8xCef+J3onQhVRY9bQssOFb1Dt8Je3ljemiyxbi4l9KaEY4Iok5A=="}`
	//var data map[string]string
	//err = json.Unmarshal([]byte(reqstr), &data)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//var ps = u.Values{}
	//// 将 map[string]string 转换为 url.Values
	//for k, v := range data {
	//	ps.Set(k, v)
	//}
	//data["timestamp"] = strconv.Itoa(int(time.Now().UnixMilli()))
	//jsonData, err := json.Marshal(data)
	//req, err := http.NewRequest("POST", "https://jsgame.live/game/v1", bytes.NewBuffer(jsonData))
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Accept", "*/*")
	//req.Header.Set("User-Agent", "PostmanRuntime/7.28.4")
	//req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	//req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("Cache-Control", "no-cache")
	//req.Header.Set("Pragma", "no-cache")
	//// 发送请求
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//defer resp.Body.Close()
	//rsp := make(map[string]interface{})
	//err = json.NewDecoder(resp.Body).Decode(&rsp)
	//if err != nil {
	//}
	//if rsp != nil && rsp["payload"] != nil {
	//	url = rsp["payload"].(map[string]interface{})["game_launch_url"].(string)
	//}
	//vs20olympx -> ae52ed2e1a8c353878ba65ed7791ac4
	uid, guid, cash := "hde0bdceshi1937", game, "5000000"
	slog.Info("game: " + game)
	url, _ = GetKey(uid, guid, cash)

	//if useProxy {
	//	url = strings.Replace(url, "5g6kpi7kjf.uapuqhki.net", "ppproxy.rpgamestest.com", 1)
	//}
	return
}

func (_ pp) GetBalance(username string) (balance float64, err error) {
	var ret transferRet
	err = invoke("/balance/current/", map[string]string{
		"externalPlayerId": username,
	}, &ret)

	balance = ret.Balance
	return
}

func (_ pp) Regist(username string) (err error) {
	currency := "THB"
	ps := map[string]string{
		"externalPlayerId": username,
		"currency":         currency,
		// "currency": "IDR",
	}
	var result struct {
		PlayerID int
	}
	err = invoke("/player/account/create/", ps, &result)
	return
}
