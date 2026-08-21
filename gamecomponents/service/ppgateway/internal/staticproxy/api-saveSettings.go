package staticproxy

import (
	"context"
	"fmt"
	"game/comm/db"
	"game/duck/ut2/jwtutil"
	"io"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func saveSettings(w http.ResponseWriter, r *http.Request) {
	// < method=load&id=vs20olympx&mgckey=AUTHTOKEN@17f8076abf3edb3c651be065210634dcd0b4ba1231c206e4405c0fa4455b9ec0~stylename@hllgd_hollygod~SESSION@5df8d3d1-dfa3-4cff-983f-6bd0853c162f~SN@273d3dc7
	// > SoundState=true_true_true_false_false;FastPlay=false;Intro=true;StopMsg=0;TurboSpinMsg=0;BetInfo=9_0;BatterySaver=false;ShowCCH=false;ShowFPH=true;CustomGameStoredData=;Coins=false;Volume=1;InitialScreen=11,8,8,9,9_3,8,8,5,5_7,4,4,11,11_5,11,11,10,10_9,4,4,10,10_10,8,8,9,9;SBPLock=true

	// < id=vs20olympx&settings=SoundState=true_true_true_false_false;FastPlay=false;Intro=true;StopMsg=0;TurboSpinMsg=0;BetInfo=9_0;BatterySaver=false;ShowCCH=false;ShowFPH=true;CustomGameStoredData=;Coins=false;Volume=1;InitialScreen=11,8,8,9,9_3,8,8,5,5_7,4,4,11,11_5,11,11,10,10_9,4,4,10,10_10,8,8,9,9;SBPLock=true&mgckey=AUTHTOKEN@17f8076abf3edb3c651be065210634dcd0b4ba1231c206e4405c0fa4455b9ec0~stylename@hllgd_hollygod~SESSION@5df8d3d1-dfa3-4cff-983f-6bd0853c162f~SN@273d3dc7
	// > SoundState=true_true_true_false_false;FastPlay=false;Intro=true;StopMsg=0;TurboSpinMsg=0;BetInfo=9_0;BatterySaver=false;ShowCCH=false;ShowFPH=true;CustomGameStoredData=;Coins=false;Volume=1;InitialScreen=11,8,8,9,9_3,8,8,5,5_7,4,4,11,11_5,11,11,10,10_9,4,4,10,10_10,8,8,9,9;SBPLock=true

	// < method=load&id=vsCommon&mgckey=AUTHTOKEN@17f8076abf3edb3c651be065210634dcd0b4ba1231c206e4405c0fa4455b9ec0~stylename@hllgd_hollygod~SESSION@5df8d3d1-dfa3-4cff-983f-6bd0853c162f~SN@273d3dc7
	// > {"MinimizedNotificationTypes":""}

	// < id=vsCommon&settings={"MinimizedNotificationTypes":""}&mgckey=AUTHTOKEN@17f8076abf3edb3c651be065210634dcd0b4ba1231c206e4405c0fa4455b9ec0~stylename@hllgd_hollygod~SESSION@5df8d3d1-dfa3-4cff-983f-6bd0853c162f~SN@273d3dc7
	// > {"MinimizedNotificationTypes":""}

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ps := ParseVariables(string(payload))

	method := ps.Str("method")
	id := ps.Str("id")
	settings := ps.Str("settings")
	mgckey := ps.Str("mgckey")

	pid, err := jwtutil.ParseToken(mgckey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key := fmt.Sprintf("%d:%s", pid, id)

	coll := db.Collection("settings")
	if method == "" && settings != "" {
		// Save settings
		coll.UpdateByID(context.TODO(), key, db.D("$set", db.D("settings", settings)), options.Update().SetUpsert(true))
	} else if method == "load" {
		// Load settings
		var doc struct {
			Settings string
		}
		coll.FindOne(context.TODO(), db.ID(key)).Decode(&doc)
		settings = doc.Settings
	}

	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有源，或指定特定源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	w.Write([]byte(settings))
}

// /api/history/v2/settings/general
func settings_general(w http.ResponseWriter, r *http.Request) {
	data := `{"language":"en","jurisdiction":"99","jurisdictionRequirements":[],"brandRequirements":["MINSPINBET:80000~IDR"]}`

	io.WriteString(w, data)
}

func playGame(w http.ResponseWriter, r *http.Request) { //todo rsp待替换
	w.Header().Set("content-type", "application/json;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有源，或指定特定源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// 设置缓存控制
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	otk := ""
	if r.Method == "POST" {
		ps := ParseVariables(string(payload))
		otk = ps.Str("otk")
	} else {
		err = r.ParseForm()
		if err != nil {
			return
		}
		otk = r.FormValue("token")
	}

	settings := fmt.Sprintf(`{
   "error": 0,
   "description": "OK",
   "currency": "$",
   "currencyOriginal": "$",
   "gameName": "vs20olympx",
   "mgckey": "%v",
   "jurisdictionRequirements": "",
   "amountType": "COIN",
   "brandRequirements": "BBB",
   "gameConfig": {
       "region": "Asia"
   }
}`, otk)
	w.Write([]byte(settings))
}

func Unread(w http.ResponseWriter, r *http.Request) { // 先写死
	buf := `<EmptyGetUnreadAnnouncementsResponseDTO><error>0</error><description>OK</description><announcements/></EmptyGetUnreadAnnouncementsResponseDTO>`
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/xml;charset=utf-8")
	w.Write([]byte(buf))
}

func PromoActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"error":0,"description":"OK","serverTime":1732759802}`))
}

func general(w http.ResponseWriter, r *http.Request) { // todo 先写死 待改成动态
	referer := r.Header.Get("Referer")
	// 解析 URL
	parsedUrl, err := url.Parse(referer)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	queryParams := parsedUrl.Query()
	language := queryParams.Get("lang")
	buf := fmt.Sprintf(`{
    "language": "%s",
    "jurisdiction": "99",
    "jurisdictionRequirements": [],
    "brandRequirements": []
}`, language)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/xml;charset=utf-8")
	// 设置缓存控制
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	w.Write([]byte(buf))
}
