package staticproxy

import (
	"encoding/json"
	"fmt"
	"game/comm"
	"game/duck/ut2/jwtutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ResponseAction101 struct {
	Status string `json:"status"`
	Data   struct {
		Ots string `json:"ots"`
	} `json:"data"`
}
type ResponseAction6 struct {
	Status string    `json:"status"`
	Data   []Result6 `json:"data"`
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
	Currency         string         `json:"currency"`
	IsDemoAccount    bool           `json:"isDemoAccount"`
	IsApiAccount     bool           `json:"isApiAccount"`
	IsShowJackpot    bool           `json:"isShowJackpot"`
	IsShowCurrency   bool           `json:"isShowCurrency"`
	IsShowDollarSign bool           `json:"isShowDollarSign"`
	ShowDemoFeatures bool           `json:"showDemoFeatures"`
	DecimalPoint     int            `json:"decimalPoint"`
	GameGroup        []int          `json:"gameGroup"`
	FunctionList     []functionList `json:"functionList"`
}
type functionList struct {
	ItemId     string `json:"itemId"`
	FunctionId string `json:"functionId"`
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
	GameType             int32    `json:"gameType"`
	MachineType          int64    `json:"machineType"`
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

type ResponseAction5 struct {
	Status string    `json:"status"`
	Data   []Result5 `json:"data"`
}

type Result5 struct {
	UID        string        `json:"uid"`
	UserStatus int           `json:"userStatus"`
	Ts         int64         `json:"ts"`
	TimeZone   string        `json:"timeZone"`
	HitJackpot []interface{} `json:"hitJackpot"`
}

// LogEntry 结构体，用于匹配 JSON 数据中的每个日志项
type LogEntry struct {
	ID         string `json:"id"`
	Namespace  string `json:"namespace"`
	Level      string `json:"level"`
	Descriptor string `json:"descriptor"`
	Env        string `json:"env"`
	Time       string `json:"time"`
	Data       struct {
		Level   string `json:"level"`
		Event   string `json:"event"`
		Message struct {
			PerformanceGrade struct {
				Grade string  `json:"grade"`
				Time  float64 `json:"time"`
			} `json:"performanceGrade"`
			// 其他可能的字段可以在这里添加
		} `json:"message"`
		AccessToken string `json:"accessToken"`
		ApiServer   string `json:"apiServer"`
		GName       string `json:"gName"`
		GameType    string `json:"gameType"`
		MType       string `json:"mType"`
		UIVersion   string `json:"ui_version"`
		UniqueKey   string `json:"uniqueKey"`
		UserAgent   string `json:"userAgent"`
		UserName    string `json:"userName"`
	} `json:"data"`
}

const (
	success  = "0000"
	zone     = "JDB_ZONE_GAME"
	slotGame = 14
)

func batchLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}

	// 解析 JSON 数据
	var logs []LogEntry
	err := json.NewDecoder(r.Body).Decode(&logs)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	str := `{"data":"%d has been created."}`
	fmt.Println(fmt.Sprintf(str, len(logs)))
	w.Write([]byte(fmt.Sprintf(str, len(logs))))
	//	w.Write([]byte(`{
	//    "data": [
	//        null,
	//        null,
	//        null,
	//        {
	//            "name": "error",
	//            "length": 209,
	//            "severity": "ERROR",
	//            "code": "23505",
	//            "detail": "Key (id)=(ddb0c9a6-feac-4299-9b02-84407ceaeee2) already exists.",
	//            "schema": "public",
	//            "table": "log",
	//            "constraint": "log_pkey",
	//            "file": "nbtinsert.c",
	//            "line": "673",
	//            "routine": "_bt_check_unique"
	//        },
	//        null,
	//        null,
	//        null,
	//        {
	//            "name": "error",
	//            "length": 209,
	//            "severity": "ERROR",
	//            "code": "23505",
	//            "detail": "Key (id)=(ed1f6d82-e43e-4eec-b5ed-2ae38c1a0722) already exists.",
	//            "schema": "public",
	//            "table": "log",
	//            "constraint": "log_pkey",
	//            "file": "nbtinsert.c",
	//            "line": "673",
	//            "routine": "_bt_check_unique"
	//        },
	//        null,
	//        {
	//            "name": "error",
	//            "length": 209,
	//            "severity": "ERROR",
	//            "code": "23505",
	//            "detail": "Key (id)=(bee31a73-26e0-478d-aabb-bbce3041d5f2) already exists.",
	//            "schema": "public",
	//            "table": "log",
	//            "constraint": "log_pkey",
	//            "file": "nbtinsert.c",
	//            "line": "673",
	//            "routine": "_bt_check_unique"
	//        },
	//        {
	//            "name": "error",
	//            "length": 209,
	//            "severity": "ERROR",
	//            "code": "23505",
	//            "detail": "Key (id)=(f5bd2e1c-b53a-474a-b60e-8f24003236f4) already exists.",
	//            "schema": "public",
	//            "table": "log",
	//            "constraint": "log_pkey",
	//            "file": "nbtinsert.c",
	//            "line": "673",
	//            "routine": "_bt_check_unique"
	//        },
	//        {
	//            "name": "error",
	//            "length": 209,
	//            "severity": "ERROR",
	//            "code": "23505",
	//            "detail": "Key (id)=(abb6a090-5f57-418d-aa9f-51b88ac820f8) already exists.",
	//            "schema": "public",
	//            "table": "log",
	//            "constraint": "log_pkey",
	//            "file": "nbtinsert.c",
	//            "line": "673",
	//            "routine": "_bt_check_unique"
	//        },
	//        {
	//            "name": "error",
	//            "length": 209,
	//            "severity": "ERROR",
	//            "code": "23505",
	//            "detail": "Key (id)=(6ef71e25-0b29-4670-b39b-1e669db320ce) already exists.",
	//            "schema": "public",
	//            "table": "log",
	//            "constraint": "log_pkey",
	//            "file": "nbtinsert.c",
	//            "line": "673",
	//            "routine": "_bt_check_unique"
	//        },
	//        {
	//            "name": "error",
	//            "length": 209,
	//            "severity": "ERROR",
	//            "code": "23505",
	//            "detail": "Key (id)=(b646c272-2c5b-4e89-8ee6-0739b88b89f7) already exists.",
	//            "schema": "public",
	//            "table": "log",
	//            "constraint": "log_pkey",
	//            "file": "nbtinsert.c",
	//            "line": "673",
	//            "routine": "_bt_check_unique"
	//        }
	//    ]
	//}`))
}

func frontendAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}

	r.ParseForm()

	action := r.FormValue("action")
	x := r.FormValue("x")
	if action == "101" {
		var a101 ResponseAction101
		a101.Status = success
		a101.Data.Ots = x
		//bytes, err := json.Marshal(a101)
		//if err != nil {
		//	w.WriteHeader(404)
		//	return
		//}
		////w.Write(bytes)
		w.Write([]byte(`{"status":"0000","data":{"ots":"06a7f691-a8d1-440c-9239-c9bf028db180"}}`))
		return
	}
	if action == "6" {
		var a6 ResponseAction6
		var a6Data Result6
		a6.Status = "0000"
		pid, err := jwtutil.ParseToken(x)
		if err != nil {
			return
		}
		item := comm.GetJDBCurrentItem(pid)
		a6Data.Currency = item.Key
		a6Data.Lvl = 0
		a6Data.UID = strconv.FormatInt(pid, 10) + "@" + item.Key
		a6Data.UserStatus = 0
		a6Data.UserName = strconv.FormatInt(pid, 10)
		a6.Data = append(a6.Data, a6Data)
		bytes, err := json.Marshal(a6)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		w.Write(bytes)
		//		w.Write([]byte(`{
		//    "status": "0000",
		//    "data": [
		//        {
		//            "uid": "demo000902@XX",
		//            "userName": "000902",
		//            "lvl": 0,
		//            "userStatus": 0,
		//            "currency": "XX"
		//        }
		//    ]
		//}`))
		return
	}
	//灰度换资源新添加的
	if action == "20" {
		//baseURL := "//jdbweb." + ut2.Domain(r.Host) + "/?tpg2tl=1&"
		baseURL := "http://192.168.1.123:8089/?tpg2tl=1&"
		// 拼接 URL 参数
		u, err := url.Parse(baseURL)
		if err != nil {
			log.Fatalf("Failed to parse base URL: %v", err)
		}
		query := u.Query()
		for key := range r.Form {
			query.Set(key, r.FormValue(key))
		}
		baseURL += query.Encode()
		jsonString := `{
	  "status": "0000",
	  "data": {
	      "url": "%s"
	  }
	}`
		w.Write([]byte(fmt.Sprintf(jsonString, baseURL)))
	}
	if action == "19" {
		var a19 ResponseAction19
		a19.Status = success
		a19.Data.IsShowAutoPlay = true
		pid, err := jwtutil.ParseToken(x)
		if err != nil {
			return
		}
		item := comm.GetJDBCurrentItem(pid)
		var a4Data Result4
		a4Data.Currency = item.Key
		a4Data.DecimalPoint = 2
		a4Data.FunctionList = make([]functionList, 0)
		// var flist functionList
		// flist.ItemId = "3"
		// flist.FunctionId = "306"
		// a4Data.FunctionList = append(a4Data.FunctionList, flist)
		//copy from origin
		//a4Data.GameGroup = []int{130, 66, 131, 67, 70, 7, 0, 9, 75, 12, 0, 0, 80, 81, 18, 22, 90, 92, 93, 30, 31, 32, 50, 55, 120, 56, 57, 58, 59, 60}
		a4Data.GameGroup = []int{7, 0, 9, 12, 0, 0, 18}
		a4Data.IsDemoAccount = false
		a4Data.IsApiAccount = true
		a4Data.IsShowJackpot = false
		a4Data.IsShowCurrency = true
		a4Data.IsShowDollarSign = true
		a4Data.ShowDemoFeatures = false
		a19.Data.Result4 = a4Data

		var a6Data Result6
		a6Data.Currency = item.Key
		a6Data.Lvl = 0
		a6Data.UID = strconv.FormatInt(pid, 10) + "@" + item.Key
		a6Data.UserStatus = 0
		a6Data.UserName = strconv.FormatInt(pid, 10)

		a19.Data.Result6 = a6Data

		var a10Data Result10
		a10Data.Status = "0000"
		a10Data.SessionID = []string{
			"",
			"",
			"",
			"CD4414C0DB1C7818D6A8D4C519F33D513F588D96FBD72A4F8E4D82530BDAFE4B2AAC7574517657CC2E3BFDAC22E9CA758B6098055CB940322695A3801434BC6AE1D5CE6A6F5356ED1A75EC32A7C4733370E9E79272A0BC2C455B29E2CE458E2F17A6212C0BD3EC11359573D7B5DC8B84082A3B720E2A821E618C1D3E2D6940186A13B05DED7B6B7839BE375CE939BC54903E1518C5BE1272CC126D06C010ADE7E0AAE5A2E46949CB8597D475270D72A3E67F641EA9DCDCD6EF170DC77F778D885DAE1AE7B5A2E00BAA8C4F2314B2F63E0ED63C341085CFE3CEF60BD336CD560F4CA1BB9A06E16AE842FFBA1E6F55FBA39DFB9C98A4D8F73E5684F967791097F996ADB935661155670165558D92206083298B74C096A22D0F4A86489698CEF198B2D3EEC04EDFCF174D097A1FBAB81B3A066CD635CACB87AC9F0857C6DB65375AE6603C57F226D6ADC781FD894343923AA85CA4A33E50D5DCCB3395A68BA0278E8CD99149A7A12B8A0082593D869179A2FB0C4C9645A3217B3A27BCAD0766F4D588800E5E0B2F695C5A6E0E0B0E5398C25CF054B9D5E0B7F9417BDD97001E8DDE3E2F24238409A26012525CD03A3C7C085788CE2461CC374C9378BC48F1D71C35384E3B355747D2C85010EC2C6765AB07D645F8AB80CE66E64B3ACBEA82062420A73BE53696190AC2",
			""}
		a10Data.Zone = zone
		a10Data.GsInfo = "192.168.1.123_443_0"
		a10Data.GsInfo = "jdb1688.net_443_0"
		//a10Data.GameType = slotGame
		GameType, _ := strconv.Atoi(r.FormValue("gameType"))
		a10Data.GameType = int32(GameType)
		a10Data.MachineType, _ = strconv.ParseInt(r.FormValue("mType"), 10, 64)
		a10Data.IsRecovery = false
		a10Data.S0 = ""
		a10Data.S1 = ""
		a10Data.S2 = ""
		a10Data.S3 = "CD4414C0DB1C7818D6A8D4C519F33D513F588D96FBD72A4F8E4D82530BDAFE4B2AAC7574517657CC2E3BFDAC22E9CA758B6098055CB940322695A3801434BC6AE1D5CE6A6F5356ED1A75EC32A7C4733370E9E79272A0BC2C455B29E2CE458E2F17A6212C0BD3EC11359573D7B5DC8B84082A3B720E2A821E618C1D3E2D6940186A13B05DED7B6B7839BE375CE939BC54903E1518C5BE1272CC126D06C010ADE7E0AAE5A2E46949CB8597D475270D72A3E67F641EA9DCDCD6EF170DC77F778D885DAE1AE7B5A2E00BAA8C4F2314B2F63E0ED63C341085CFE3CEF60BD336CD560F4CA1BB9A06E16AE842FFBA1E6F55FBA39DFB9C98A4D8F73E5684F967791097F996ADB935661155670165558D92206083298B74C096A22D0F4A86489698CEF198B2D3EEC04EDFCF174D097A1FBAB81B3A066CD635CACB87AC9F0857C6DB65375AE6603C57F226D6ADC781FD894343923AA85CA4A33E50D5DCCB3395A68BA0278E8CD99149A7A12B8A0082593D869179A2FB0C4C9645A3217B3A27BCAD0766F4D588800E5E0B2F695C5A6E0E0B0E5398C25CF054B9D5E0B7F9417BDD97001E8DDE3E2F24238409A26012525CD03A3C7C085788CE2461CC374C9378BC48F1D71C35384E3B355747D2C85010EC2C6765AB07D645F8AB80CE66E64B3ACBEA82062420A73BE53696190AC2"
		a10Data.S4 = ""
		a10Data.GameUid = strconv.FormatInt(pid, 10) + "@" + item.Key
		a10Data.GamePass = "8158025"
		a10Data.UseSSL = false
		a10Data.StreamingUrl = struct{}{}
		a10Data.AchievementServerUrl = ""
		a10Data.ChatServerUrl = ""
		a10Data.IsWSBinary = false
		a19.Data.Result10 = a10Data

		bytes, err := json.Marshal(a19)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		w.Write(bytes)

		return
	}
	if action == "23" {
		w.Write([]byte("{\"status\":\"0000\",\"data\":[]}"))
		return
	}
	if action == "24" {
		w.Write([]byte("{\"status\":\"0000\",\"data\":{\"enable\":false}}"))
		return
	}
	//EGRET_READY
	if action == "13" {
		w.Write([]byte("{\"status\":\"0000\"}"))
		return
	}
	//heartBeat
	if action == "5" {
		var a5 ResponseAction5
		var a5Data Result5
		a5.Status = success
		pid, err := jwtutil.ParseToken(x)
		if err != nil {
			return
		}
		item := comm.GetJDBCurrentItem(pid)
		a5Data.UID = strconv.FormatInt(pid, 10) + "@" + item.Key
		a5Data.UserStatus = 0
		a5Data.Ts = time.Now().UnixNano()
		a5Data.TimeZone = time.Local.String()
		a5.Data = append(a5.Data, a5Data)
		bytes, err := json.Marshal(a5)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		w.Write(bytes)
		return
	}
	fmt.Println("none")
}
