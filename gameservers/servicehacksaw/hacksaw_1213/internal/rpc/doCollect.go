package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"serve/comm/db"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicehacksaw/hacksawcomm"
	"time"
)

func init() {
	hacksawcomm.RegRpc("doCollect", doCollect)
}

var (
	sampleCollect = `{
  "round": {
    "status": "completed",
    "jackpotWin": null,
    "roundId": "158105077",
    "possibleActions": [],
    "events": [
      {
        "et": 2,
        "etn": "reveal",
        "en": "0",
        "ba": "0",
        "bc": "0",
        "wa": "0",
        "wc": "0",
        "awa": "0",
        "awc": "0",
        "c": {
          "reelSet": "default",
          "stops": [
            "20",
            "27",
            "60",
            "11",
            "33"
          ],
          "grid": "-,122/11///*2,01*),02-"
        }
      }
    ]
  },
  "promotionNoLongerAvailable": false,
  "promotionWin": null,
  "offer": null,
  "freeRoundOffer": null,
  "statusCode": 0,
  "statusMessage": "",
  "accountBalance": {
    "currencyCode": "EUR",
    "balance": "499800",
    "realBalance": null,
    "bonusBalance": null
  },
  "statusData": null,
  "dialog": null,
  "customData": null,
  "serverTime": "2025-04-18T06:55:29Z"
}`
)

func doCollect(msg *nats.Msg) (ret []byte, err error) {

	pid, sessionData, err := hacksawcomm.ParseHackSawReq(msg.Data)
	if err != nil {
		return nil, err
	}
	err = db.CallWithPlayer(pid, func(plr *hacksawcomm.Player) error {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(sampleCollect), &data); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return err
		}

		balance, err := slotsmongo.GetBalance(pid)
		if err != nil {
			return err
		}
		data["serverTime"] = time.Now().UTC().Format(time.RFC3339)
		data["round"].(map[string]interface{})["roundId"] = sessionData.RoundID
		data["accountBalance"].(map[string]interface{})["balance"] = ut.HackGold2Money(balance)
		fmt.Println(sessionData.RoundID)
		modifiedJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return err
		}

		ret = modifiedJSON

		//ret = info.Bytes()
		return nil
	})

	return
}
