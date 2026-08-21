package staticproxy

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"game/duck/ut2/jwtutil"
	"game/service/facaigateway/internal/gamedata"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//go:embed html/lineSettings.js
var lineSettings []byte

func LineSetting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	if r.Method == "OPTIONS" {
		return
	}
	w.Write(JsonConfigHook(lineSettings, w, r))
}

//go:embed html/ExtensionLineSettings.js
var extensionLineSettings []byte

func ExtensionLineSetting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	if r.Method == "OPTIONS" {
		return
	}
	w.Write(JsonConfigHook(extensionLineSettings, w, r))
}
func GetGamePageVersion(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	w.Write([]byte(`{"isSuccess":true,"errorMessage":null,"returnObject":"25.4.1"}`))
}
func GetJwtToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	w.Write([]byte(`{
    "isSuccess": true,
    "errorMessage": null,
    "returnObject": {
        "token": "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJkZW1vMTMxOSIsImV4cCI6MTc0NTYxNTU3MCwiaWF0IjoxNzQ0NjE1NTcxLCJqdGkiOiJGQ1RfREVNX1RfZ3Vlc3QifQ.U9d1vd30byHmVCoQrLmoMNpuWtOGmEUVdE6rT84pwHg"
    }
}`))
}

type Version struct {
	Product string `json:"product"`
	Version string `json:"version"`
}

type BoundaryRatioSettings struct {
	Judgment string  `json:"judgment"`
	Value    float64 `json:"value"`
}

type RestrictBoundaryRatioSettings struct {
	BoundaryX BoundaryRatioSettings `json:"boundaryX"`
	BoundaryY BoundaryRatioSettings `json:"boundaryY"`
}

type DefaultPositionRatioSettings struct {
	AlignPosition string  `json:"alignPosition"`
	PositionX     float64 `json:"positionX"`
	PositionY     float64 `json:"positionY"`
}

type CustomSettings struct {
	DefaultPositionRatioSettings  DefaultPositionRatioSettings  `json:"defaultPositionRatioSettings"`
	RestrictBoundaryRatioSettings RestrictBoundaryRatioSettings `json:"restrictBoundaryRatioSettings"`
}

type ReturnObject struct {
	Currency                 string         `json:"currency"`
	Language                 string         `json:"language"`
	GameServerSocket         string         `json:"gameServerSocket"`
	GameId                   string         `json:"gameId"`
	HomeUrl                  string         `json:"homeUrl"`
	TimeZone                 string         `json:"timeZone"`
	DecimalPoint             int            `json:"decimalPoint"`
	ReduceDisplay            string         `json:"reduceDisplay"`
	EventBtn                 bool           `json:"eventBtn"`
	EventPopup               bool           `json:"eventPopup"`
	Token                    string         `json:"token"`
	LoginName                string         `json:"loginName"`
	GameZone                 string         `json:"gameZone"`
	RedirectUrl              string         `json:"redirectUrl"`
	ReportBtn                bool           `json:"reportBtn"`
	Version                  []Version      `json:"version"`
	IsUS                     bool           `json:"isUS"`
	IsXinStar                bool           `json:"isXinStar"`
	OutsideHomeBtn           bool           `json:"outsideHomeBtn"`
	PrizeViewer              bool           `json:"prizeViewer"`
	JackpotStatus            bool           `json:"jackpotStatus"`
	CloseFeatureBuy          bool           `json:"closeFeatureBuy"`
	LevelStatus              bool           `json:"levelStatus"`
	MissionStatus            bool           `json:"missionStatus"`
	ActivityStatus           bool           `json:"activityStatus"`
	IconPopup                bool           `json:"iconPopup"`
	FishDefaultRoomIndex     int            `json:"fishDefaultRoomIndex"`
	SupportFullScreenStatus  bool           `json:"supportFullScreenStatus"`
	ShowCurrencySymbolStatus bool           `json:"showCurrencySymbolStatus"`
	ShowLogoStatus           bool           `json:"showLogoStatus"`
	IsFollowBet              bool           `json:"isFollowBet"`
	IsExtra                  bool           `json:"isExtra"`
	CanChangeGame            bool           `json:"canChangeGame"`
	IsLandscapeGame          bool           `json:"isLandscapeGame"`
	ShowAutoPlay             bool           `json:"showAutoPlay"`
	ShowGameName             bool           `json:"showGameName"`
	IsHelpPageShown          bool           `json:"isHelpPageShown"`
	LogoUrl                  interface{}    `json:"logoUrl"`
	HelpVersion              interface{}    `json:"helpVersion"`
	CustomSettings           CustomSettings `json:"customSettings"`
}

type EnterGameResponse struct {
	IsSuccess    bool         `json:"isSuccess"`
	ErrorMessage interface{}  `json:"errorMessage"`
	ReturnObject ReturnObject `json:"returnObject"`
}

func EnterGame(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	jsonStr := `{
    "isSuccess": true,
    "errorMessage": null,
    "returnObject": {
        "currency": "DEM",
        "language": "1",
        "gameServerSocket": "192.168.1.123:8089",
        "gameId": "22043",
        "homeUrl": "/close",
        "timeZone": "Etc/GMT+4",
        "decimalPoint": 2,
        "reduceDisplay": "false",
        "eventBtn": false,
        "eventPopup": false,
        "token": "%s",
        "loginName": "demo1319",
        "gameZone": "FC_GAME_ZONE",
        "redirectUrl": "",
        "reportBtn": false,
        "version": [
            {
                "product": "gamePageExtension",
                "version": "25.4.1"
            },
            {
                "product": "cocos2d",
                "version": "2.0.10-20211004"
            },
            {
                "product": "jszip",
                "version": "fe1e4"
            },
            {
                "product": "22043",
                "version": "0c4a212"
            }
        ],
        "isUS": false,
        "isXinStar": false,
        "outsideHomeBtn": false,
        "prizeViewer": false,
        "jackpotStatus": false,
        "closeFeatureBuy": false,
        "levelStatus": false,
        "missionStatus": false,
        "activityStatus": false,
        "iconPopup": false,
        "fishDefaultRoomIndex": -1,
        "supportFullScreenStatus": true,
        "showCurrencySymbolStatus": false,
        "showLogoStatus": false,
        "isFollowBet": false,
        "isExtra": false,
        "canChangeGame": true,
        "isLandscapeGame": false,
        "showAutoPlay": true,
        "showGameName": false,
        "isHelpPageShown": false,
        "logoUrl": null,
        "helpVersion": null,
        "customSettings": {
            "defaultPositionRatioSettings": {
                "alignPosition": "left",
                "positionX": 0,
                "positionY": 0.097
            },
            "restrictBoundaryRatioSettings": {
                "boundaryX": {
                    "judgment": "none",
                    "value": 0
                },
                "boundaryY": {
                    "judgment": "less",
                    "value": 0.845
                }
            }
        }
    }
}`
	// 解析 JSON
	var response EnterGameResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	m := map[string]string{}
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(payload, &m)
	if err != nil {
		return
	}
	pid, gameId, err := jwtutil.ParseTokenData(m["token"])
	response.ReturnObject.Token = m["token"]
	response.ReturnObject.LoginName, response.ReturnObject.GameId = strconv.Itoa(int(pid)), strings.Split(gameId, "_")[1]
	response.ReturnObject.GameServerSocket = gamedata.Get().SocketConn
	bytes, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(bytes)
}

func GetEventStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	w.Write([]byte(`{
    "isSuccess": false,
    "errorMessage": "Error.AgentAccountNotFound",
    "returnObject": null
}`))
}
func GetAffiliatedStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	w.Write([]byte(`{
    "isSuccess": true,
    "errorMessage": null,
    "returnObject": {
        "hasDailyMission": false,
        "cardMissionStatus": false,
        "propsStatus": false
    }
}`))
}
func GetFreeSpinStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	w.Write([]byte(`{
    "isSuccess": true,
    "errorMessage": null,
    "returnObject": {
        "freeSpinStatus": false,
        "freeSpinPopup": false,
        "freeSpinNeedStartEventID": null,
        "hasFreeSpinGiftGiftCode": false,
        "freeSpinGiftCodeIds": [],
        "freeSpinGiftCodeExpireTime": null,
        "freeSpinGiftCodeCreateTime": null
    }
}`))
}
