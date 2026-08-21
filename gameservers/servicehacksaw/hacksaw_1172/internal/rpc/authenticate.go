package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/servicehacksaw/hacksaw_1172/internal"
	"serve/servicehacksaw/hacksawcomm"
	"strconv"
	"time"
)

func init() {
	hacksawcomm.RegRpc("authenticate", authenticate)
}

var (
	authenticateTemplate = `{
  "gameId": "1172",
  "partnerId": "0",
  "roundId": "0",
  "roundStatus": null,
  "gameState": null,
  "freeRoundOffer": null,
  "pendingWin": "0",
  "events": null,
  "progressionData": null,
  "playerId": "62a2d7a2-4560-4d2f-a848-fdf1be87d34f",
  "name": "Demo",
  "languageCode": "en",
  "sessionUuid": "6c6dd6c5-793d-419e-8507-e4d6047a28bd",
  "jurisdiction": "curacao",
  "cheatsEnabled": false,
  "betLevels": [
    "10",
    "20",
    "40",
    "60",
    "80",
    "100",
    "120",
    "140",
    "160",
    "180",
    "200",
    "300",
    "400",
    "500",
    "600",
    "700",
    "800",
    "900",
    "1000",
    "1500",
    "2000",
    "2500",
    "3000",
    "3500",
    "4000",
    "4500",
    "5000",
    "7500",
    "10000"
  ],
  "defaultBetLevel": "200",
  "autoPlayAlternatives": [
    "10",
    "25",
    "50",
    "75",
    "100",
    "500",
    "1000"
  ],
  "disableRoundHistory": false,
  "minimumRoundDuration": "0",
  "autoplayLossLimitRequired": false,
  "autoplayWinLimitRequired": false,
  "autoplayDisabled": false,
  "turboDisabled": false,
  "superTurboDisabled": false,
  "slamStopDisabled": false,
  "sessionRescueEnabled": true,
  "keepAliveInterval": "300",
  "rm": "96",
  "autoCollectAfter": "86400",
  "rollbackAfter": "86400",
  "clearOldRoundImmediatelyOnNewRound": false,
  "bonusGames": [
    {
      "bonusGameId": "mod_boxes",
      "betCostMultiplier": "3",
      "expectedRtp": 96.41
    },
    {
      "bonusGameId": "bonus_boxes",
      "betCostMultiplier": "100",
      "expectedRtp": 96.38
    },
    {
      "bonusGameId": "bonus_reels",
      "betCostMultiplier": "200",
      "expectedRtp": 96.35
    }
  ],
  "stopAutoplayOnFeatureWin": false,
  "displayRtp": false,
  "displaySessionTimer": false,
  "displayNetPosition": false,
  "disableBetWhenScreensAreOpen": false,
  "spacebarDisabled": false,
  "sessionTimeoutSeconds": "1800",
  "maxFeatureCost": "0",
  "maxFeatureSpinCost": "0",
  "backendGameVersion": "1.0.2",
  "serverVersion": "2.0.209",
  "maxExposure": "0",
  "rememberBetLevel": true,
  "hideGameInfoRtp": false,
  "disableWinHistory": false,
  "displayMaxWinOdds": false,
  "displayMaxWinMultiplier": false,
  "hideGameInfoInterrupted": false,
  "availableTournament": null,
  "availableMission": null,
  "availableMysteryPrize": null,
  "offlinePromotionWins": null,
  "replayLinkDisabled": false,
  "displayGameInfoRtpRange": false,
  "parallelRoundsSupportDisabled": false,
  "hideGameInfoDate": false,
  "displayPayoutTableOnGameLaunch": false,
  "displayPayoutTableAsMultipliers": false,
  "disableMidRoundFullScreenMenus": false,
  "disableExternalLinks": false,
  "disableKeybinds": false,
  "roundInProgressCurrency": null,
  "statusCode": 0,
  "statusMessage": "",
  "accountBalance": {
    "currencyCode": "EUR",
    "balance": "500000",
    "realBalance": null,
    "bonusBalance": null
  },
  "statusData": null,
  "dialog": null,
  "customData": null,
  "serverTime": "2025-04-18T06:14:15Z"
}`
)

func authenticate(msg *nats.Msg) (ret []byte, err error) {
	var ps hacksawcomm.Variables
	json.Unmarshal(msg.Data, &ps)
	//fmt.Println(ps)

	var pid int64
	pid, err = jwtutil.ParseToken(ps.Str("token"))
	if err != nil {
		return
	}

	err = db.CallWithPlayer(pid, func(plr *hacksawcomm.Player) error {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(authenticateTemplate), &data); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return err
		}
		info, err := redisx.LoadAppIdCache(plr.AppID)
		//修改token
		token, err := jwtutil.NewTokenWithData(pid, time.Now().Add(time.Hour*12), internal.GameID)
		data["sessionUuid"] = token
		data["languageCode"] = ps.Str("languageCode")
		data["serverTime"] = time.Now().UTC().Format(time.RFC3339)
		data["playerId"] = token
		balance, err := slotsmongo.GetBalance(pid)
		if err != nil {
			return err
		}
		betLevels := []string{}
		for i := range info.Cs {
			betLevels = append(betLevels, strconv.FormatFloat(info.Cs[i], 'f', -1, 64))
		}
		data["betLevels"] = betLevels
		data["defaultBetLevel"] = strconv.FormatFloat(info.DefaultCs, 'f', -1, 64)
		data["accountBalance"].(map[string]interface{})["balance"] = balance / 100
		data["accountBalance"].(map[string]interface{})["currencyCode"] = ps.Str("currency")
		modifiedJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return err
		}

		ret = modifiedJSON
		plr.SpinCountOfThisEnter = 0
		return nil
	})

	return
}
