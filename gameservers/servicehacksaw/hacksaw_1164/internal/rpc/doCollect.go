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
    "status": "wfwpc",
    "jackpotWin": null,
    "roundId": "151540888",
    "possibleActions": [],
    "events": [
      {
        "et": 2,
        "etn": "reveal",
        "en": "0",
        "ba": "0",
        "bc": "0",
        "wa": "200",
        "wc": "0",
        "awa": "200",
        "awc": "0",
        "c": {
          "actions": [
            {
              "at": "gridwin",
              "data": {
                "winAmount": "200",
                "symbol": "10",
                "mask": "000000000100011000100001000000",
                "count": "5"
              }
            },
            {
              "at": "orbValues",
              "data": {
                "red": "5",
                "blue": "1"
              }
            },
            {
              "at": "remove",
              "data": {
                "mask": "000000000100011000100001000000"
              }
            }
          ],
          "reelSet": "default",
          "stops": [
            "58",
            "36",
            "0",
            "19",
            "39"
          ],
          "grid": "-.1-/,-2)002)-022)-*21)0.21-0+*1"
        }
      },
      {
        "et": 2,
        "etn": "collapse",
        "en": "0",
        "ba": "0",
        "bc": "0",
        "wa": "0",
        "wc": "0",
        "awa": "200",
        "awc": "0",
        "c": {
          "reelSet": "default",
          "grid": "-.1-/.*2)0.-)-0,-)-*,1)0.01-0+*1"
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
  "serverTime": "2025-04-18T01:26:01Z"
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
