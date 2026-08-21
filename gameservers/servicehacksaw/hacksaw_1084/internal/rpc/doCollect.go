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
    "roundId": "122646606",
    "possibleActions": [],
    "events": [
      {
        "et": 2,
        "etn": "reveal",
        "en": "0",
        "ba": "0",
        "bc": "0",
        "wa": "20",
        "wc": "0",
        "awa": "20",
        "awc": "0",
        "c": {
          "actions": [
            {
              "at": "gridwin",
              "data": {
                "winAmount": "20",
                "symbol": "5",
                "mask": "000000000000010001100001100000",
                "count": "5"
              }
            },
            {
              "at": "remove",
              "data": {
                "mask": "111111111111101110011110011111"
              }
            },
            {
              "at": "updateColumnMultipliers",
              "data": {
                "value": [
                  "0",
                  "0",
                  "1",
                  "3",
                  "1"
                ]
              }
            }
          ],
          "reelSet": "default",
          "stops": [
            "8",
            "29",
            "16",
            "20",
            "52"
          ],
          "grid": "-.0),+*.--,*/0)-210--)1*.--2.).*"
        }
      },
      {
        "et": 2,
        "etn": "collapse",
        "en": "0",
        "ba": "0",
        "bc": "0",
        "wa": "100",
        "wc": "0",
        "awa": "120",
        "awc": "0",
        "c": {
          "actions": [
            {
              "at": "gridwin",
              "data": {
                "winAmount": "100",
                "symbol": "5",
                "mask": "000000000000000000100001000111",
                "count": "5",
                "baseWinAmount": "20",
                "winMultipliers": [
                  "1",
                  "3",
                  "1"
                ]
              }
            }
          ],
          "reelSet": "default",
          "stops": [
            "8",
            "29",
            "16",
            "20",
            "52"
          ],
          "grid": "-.4*000)**10))11)+**-)+**-)00---"
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
    "balance": "499600",
    "realBalance": null,
    "bonusBalance": null
  },
  "statusData": null,
  "dialog": null,
  "customData": null,
  "serverTime": "2025-04-17T08:57:01Z"
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
