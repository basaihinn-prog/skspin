package staticproxy

import (
	"context"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/ut2/httputil"
	"game/duck/ut2/jwtutil"
	"game/service/jiligateway/jilicomm"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/phuslu/iploc"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	// https://uat-wbslot-fd-jlfafafa1.kafa010.com/sso-login.api?key=19d8eede5f336e61bbaf2764fdd0b6a72ec61d19&lang=en-US
	// https://wbslot-fd-jlfafafa1.rpfafafa33cdn.com/sso-login.api
	uat_wbslot_fd_mux.HandleFunc("/sso-login.api", sso_login)
}

func sso_login(w http.ResponseWriter, r *http.Request) {
	// key=7a1e467dbd8c49ba962acf414db005346098a8ba&lang=en-US
	// 	{
	//     "ShopingMall": 1,
	//     "SmartMessage": 2,
	//     "Promotion": 3,
	//     "TreasureChest": 4,
	//     "GameHistory": 5,
	//     "FeaturedDisplay": 6,
	//     "CoinAccuracy": 7,
	//     "CurrencySymbol": 8,
	//     "ItemBoxImport": 9,
	//     "AllChangeGameImport": 10,
	//     "AutoPlay": 11,
	//     "IsDelay": 12,
	//     "CloseSpeedUp": 13,
	//     "CloseBackpack": 14,
	//     "Trail": 15,
	//     "ShowAutoSetting": 16,
	//     "NoSoundUnder1": 17,
	//     "NoQuickSpin": 18,
	//     "RatioOnView": 19,
	//     "ShowPlateformVer": 20,
	//     "ShowTime": 21,
	//     "ShowPlayTime": 22,
	//     "BlockLobbyOff": 23,
	//     "NoRedSpot": 24,
	//     "CloseVip": 25,
	//     "CloseBuyBonusInfo": 26,
	//     "CloseSettingInfo": 27,
	//     "UseKilo": 28,
	//     "ClickAutoSetting": 29,
	//     "CloseBuyBonusAdd": 30,
	//     "RemoveDecimal": 31,
	//     "ShowBuyBonusBetInfo": 32,
	//     "CloseWinTxtWithZero": 33,
	//     "ShowNetWin": 34,
	//     "CloseSideFeatures": 35,
	//     "RealityCheck": 36,
	//     "CloseJPList": 37,
	//     "DisableSettingInfo": 38,
	//     "CloseManual": 39,
	//     "CloseHotChilli": 40,
	//     "CloseFreeSpin": 41,
	//     "CloseTiggerRank": 42,
	//     "OpenTurbo": 43,
	//     "CloseAutoShowEventWebView": 44
	// }
	ip := httputil.GetIPFromRequest(r)
	loc := string(iploc.Country(net.ParseIP(ip)))
	key := r.URL.Query().Get("key")
	pid, gameid, err := jwtutil.ParseTokenData(key)
	fmt.Printf("userid=%d,gameid=%s,ip=%s,loc=%s\n", pid, gameid, ip, loc)
	insertLoginDetail(pid, gameid, ip, loc)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("error key [%s]", key)))
		return
	}

	w.Header().Set("Content-Tye", "application/json; charset=utf-8")

	// lo.Range()
	//switchOffs := []int{26, 27, 42, 14, 37, 25, 37, 24, 40, 3, 2}
	//switchOffs := []int{26, 27, 42, 43, 14}
	params := comm.GetEXParams(pid, gameid)
	switchOffs := []int{26, 27, 43, 4}
	if params.BackPackOff == 0 {
		switchOffs = append(switchOffs, 14)
	}
	if params.OpenScreenOff == 0 {
		switchOffs = append(switchOffs, 49)
	}
	if params.SidebarOff == 0 {
		switchOffs = append(switchOffs, 53)
	}

	// sort.Ints(switchOffs)
	if _, ok := jilicomm.ChangeGameMap2[gameid]; ok {
		switchOffs = []int{26, 27, 42, 14, 37, 25, 37, 24, 40, 3, 2}
		//if params.OpenScreenOff == 0 {
		//	switchOffs = append(switchOffs, 49)
		//}
	}
	if _, ok := jilicomm.ChangeGameMap33[gameid]; ok {
		switchOffs = []int{26, 27, 42, 14, 37, 25, 37, 24, 40, 3, 2}
		//if params.OpenScreenOff == 0 {
		//	switchOffs = append(switchOffs, 49)
		//}
	}
	if _, ok := jilicomm.ChangeGameMap44[gameid]; ok {
		switchOffs = []int{26, 27, 42, 14, 37, 25, 37, 24, 40, 3, 2}
		//if params.OpenScreenOff == 0 {
		//	switchOffs = append(switchOffs, 49)
		//}
	}
	// switchOffs := lo.Range(43)

	msg := M{"homeUrl": "", "linecode": 0, "profile": M{"id": "", "aid": pid, "apiId": 555, "transactionMode": 0, "subAgentCode": 0, "isLobbyOpen": true, "meta": M{"agentAccount": "rp_Online@api-555.game"}, "platform": "", "lobbyMode": 2, "switchOffs": switchOffs, "wallets": nil, "nickname": strconv.Itoa(int(pid)), "newNickname": "", "siteId": pid, "account": "123456@api-555.game", "coin": 0, "isJPEnabled": 0, "linecode": 1, "prefix": "", "betLevel": 0, "license": 0, "isGiftCodeOpen": 0, "freeSpinBetValue": 0, "apiType": 0, "walletType": 2}, "token": key, "response": M{"error": 0, "message": "", "time": time.Now().Unix()}, "platformVersion": "uat.2.0.13", "lobbyMode": 0}
	//msg := M{"homeUrl": "", "linecode": 0, "profile": M{"id": "", "aid": pid, "apiId": 555, "transactionMode": 0, "subAgentCode": 0, "isLobbyOpen": false, "meta": M{"agentAccount": "rp_Online@api-555.game"}, "platform": "", "lobbyMode": 2, "switchOffs": switchOffs, "wallets": nil, "nickname": strconv.Itoa(int(pid)), "newNickname": "", "siteId": pid, "account": "123456@api-555.game", "coin": 0, "isJPEnabled": 0, "linecode": 1, "prefix": "", "betLevel": 0, "license": 0, "isGiftCodeOpen": 0, "freeSpinBetValue": 0, "apiType": 0, "walletType": 2}, "token": key, "response": M{"error": 0, "message": "", "time": time.Now().Unix()}, "platformVersion": "uat.2.0.13", "lobbyMode": 0}

	// game-4: {"homeUrl":"","linecode":0,"profile":{"id":"","aid":1189004,"apiId":555,"transactionMode":0,"subAgentCode":0,"isLobbyOpen":false,"meta":{"agentAccount":"rp_Online@api-555.game"},"platform":"web","lobbyMode":0,"switchOffs":[26,27,42,43],"wallets":null,"nickname":"1234567","newNickname":"","siteId":114298358,"account":"1234567@api-555.game","coin":0,"isJPEnabled":0,"linecode":0,"prefix":"","betLevel":-1,"license":0,"isGiftCodeOpen":false,"freeSpinBetValue":0,"apiType":0,"walletType":1},"token":"cf3ff118fd8060fa8eb74ad2c9151cb97e88895b","response":{"error":0,"message":"","time":1716790241},"platformVersion":"uat.2.0.13","lobbyMode":0}

	w.Write(lo.Must(ut.GetJsonRaw(msg)))
}

func insertLoginDetail(pid int64, gameid, ip, loc string) {
	appId, uid, err := slotsmongo.GetPlayerInfo(pid)
	if err == nil {
		ld := comm.GameLoginDetail{
			ID:        primitive.NewObjectID(),
			Pid:       appId,
			UserID:    uid,
			GameID:    gameid,
			Ip:        ip,
			Loc:       ut.CountryNameByCode(loc),
			LoginTime: time.Now().Unix(),
		}
		CollGameLoginDetail := db.Collection2("GameAdmin", "GameLoginDetail")
		_, err := CollGameLoginDetail.InsertOne(context.TODO(), ld)
		if err != nil {
			log.Printf("loginDetail.InsertOne occured an error => %s", err.Error())
		}
	} else {
		log.Printf("Players.FindId occured an error => %s", err.Error())
	}
}
