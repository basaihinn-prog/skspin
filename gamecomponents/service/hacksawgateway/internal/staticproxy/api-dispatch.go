package staticproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game/comm/mq"
	"game/duck/ut2/jwtutil"
	"github.com/nats-io/nats.go"
	"io"
	"net/http"
	"time"
)

type RequestData struct {
	Seq          int    `json:"seq"`
	Partner      string `json:"partner"`
	GameID       int    `json:"gameId"`
	GameVersion  string `json:"gameVersion"`
	Currency     string `json:"currency"`
	LanguageCode string `json:"languageCode"`
	Branding     string `json:"branding"`
	Channel      int    `json:"channel"`
	Mode         int    `json:"mode"`
	Token        string `json:"token"`
	UserAgent    string `json:"userAgent"`
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	var (
		payload []byte
		err     error
	)
	payload, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//去除 { xx }
	//payload = payload[1:]
	//payload = payload[:len(payload)-1]
	r.Body = io.NopCloser(bytes.NewReader(payload))

	header := nats.Header(r.Header)
	header.Set("query", r.URL.RawQuery)
	//var reqData RequestData
	//err = json.NewDecoder(r.Body).Decode(&reqData)
	//if err != nil {
	//	http.Error(w, "Invalid JSON format", http.StatusBadRequest)
	//	return
	//}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {

	}
	var data RequestData
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		http.Error(w, "JSON 解析失败", http.StatusBadRequest)
		return
	}
	subj := fmt.Sprintf("hacksaw_%d.%s", data.GameID, "authenticate")
	resp, err := mq.NC().RequestMsg(&nats.Msg{
		Subject: subj,
		Data:    payload,
		Header:  header,
	}, time.Second*60)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errstr := resp.Header.Get("error"); errstr != "" {
		http.Error(w, errstr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有源，或指定特定源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// 设置缓存控制
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")

	w.Write(resp.Data)
}

// Bet 结构体表示 bets 数组中的单个对象
type Bet struct {
	BetAmount string `json:"betAmount"`
	BuyBonus  string `json:"buyBonus"`
}

// ContinueInstructions 结构体表示 continueInstructions 对象
type ContinueInstructions struct {
	Action string `json:"action"`
}

// SessionData 结构体表示整个 JSON 数据
type SessionData struct {
	Seq                  int                   `json:"seq"`
	SessionUUID          string                `json:"sessionUuid"`
	Bets                 []Bet                 `json:"bets,omitempty"` // omitempty 表示字段为空时不序列化
	OfferID              string                `json:"offerId"`        // 使用指针类型表示可以为 null
	PromotionID          string                `json:"promotionId"`    // 使用指针类型表示可以为 null
	Autoplay             bool                  `json:"autoplay"`
	RoundID              string                `json:"roundId,omitempty"`
	ContinueInstructions *ContinueInstructions `json:"continueInstructions,omitempty"`
}

func bet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	var (
		payload []byte
		err     error
	)
	payload, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var req SessionData
	err = json.Unmarshal(payload, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	action := "doSpin"
	pid, GameID, err := jwtutil.ParseTokenData(req.SessionUUID)
	fmt.Println(pid, GameID, err)
	if req.ContinueInstructions != nil && req.ContinueInstructions.Action == "win_presentation_complete" {
		action = "doCollect"
	}
	subj := fmt.Sprintf("%s.%s", GameID, action)
	resp, err := mq.NC().RequestMsg(&nats.Msg{
		Subject: subj,
		Data:    payload,
	}, time.Second*60)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errstr := resp.Header.Get("error"); errstr != "" {
		http.Error(w, errstr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有源，或指定特定源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// 设置缓存控制
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")

	w.Write(resp.Data)
}

func gameInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	ps := ParseVariables(r.URL.RawQuery)
	header := nats.Header(r.Header)
	header.Set("query", r.URL.RawQuery)
	subj := fmt.Sprintf("hacksaw_%s.%s", ps["gameId"], "gameInfo")
	resp, err := mq.NC().RequestMsg(&nats.Msg{
		Subject: subj,
		Header:  header,
	}, time.Second*60)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errstr := resp.Header.Get("error"); errstr != "" {
		http.Error(w, errstr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有源，或指定特定源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// 设置缓存控制
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")

	w.Write(resp.Data)
}

func keepAlive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有源，或指定特定源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// 设置缓存控制
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	json := `{
    "statusCode": 0,
    "statusMessage": "",
    "accountBalance": null,
    "statusData": null,
    "dialog": null,
    "customData": null,
    "serverTime": "%s"
}`
	json = fmt.Sprintf(json, time.Now().UTC().Format(time.RFC3339))
	w.Write([]byte(json))
}
