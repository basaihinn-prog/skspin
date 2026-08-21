package pp

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"serve/comm/ut"

	"github.com/stretchr/testify/assert"
)

func TestInvoke(t *testing.T) {
	// POST /IntegrationService/v3/http/FreeRoundsBonusAPI/v2/bonus/create

	// secureLogin=username&bonusCode=421&startDate=1482588510&expirationDate=1482598510&rounds=10&

	payload := json.RawMessage(`
{"gameList": [
	{"gameId": "vs20olympx", "betValues": [
		{"totalBet": 20, "currency": "THB"}
	]}
]}`)

	/*
			payload = json.RawMessage(`{
		"gameList":[
		{
		"gameId":"vs50pixie",
		"betValues":[
		{"betPerLine":1.00,"currency":"USD"}
		]
		},
		{
		"gameId":"vs50kingkong",
		"betValues":[
		{"betPerLine":0.15,"currency":"USD"}
		]
		}]
		}`)
	*/

	bonusCode := strconv.Itoa(int(time.Now().Unix()))

	startDate := time.Now().Unix()
	err := invokeFRB("/v2/bonus/create", map[string]string{
		"bonusCode":      bonusCode,
		"startDate":      strconv.Itoa(int(startDate)),
		"expirationDate": strconv.Itoa(int(startDate + 30*3600*24)),
		"rounds":         "10",
	}, payload, nil)

	assert.Nil(t, err)

	err = invokeFRB("/v2/players/add", map[string]string{
		"bonusCode": bonusCode,
	}, json.RawMessage(`{"playerList": ["atestuser1"]}`), nil)
	assert.Nil(t, err)

	var info json.RawMessage
	err = invokeFRB("/getPlayersFRB", map[string]string{
		"playerId": "atestuser1",
	}, nil, &info)
	assert.Nil(t, err)
	os.Stdout.Write(info)
}
func TestFRBAddPlayer(t *testing.T) {
	// POST /IntegrationService/v3/http/FreeRoundsBonusAPI/v2/players/add HTTP/1.1
	// Host: api.prerelease-env.biz
	// Content-Type: application/json
	// Cache-Control: no-cache
	// secureLogin=username&bonusCode=421&hash=39554fed4f41132eb8fe75e9a7ba3df6
	// {"playerList": ["449986","450013","450509","437070"]}

	err := invokeFRB("/v2/players/add", map[string]string{
		"bonusCode": "25",
	}, json.RawMessage(`{"playerList": ["atestuser1"]}`), nil)
	assert.Nil(t, err)

	var info json.RawMessage
	err = invokeFRB("/getPlayersFRB", map[string]string{
		"playerId": "atestuser1",
	}, nil, &info)
	assert.Nil(t, err)
	os.Stdout.Write(info)
}

func TestGetGameList(t *testing.T) {
	var pp pp
	games, _ := pp.GetGameList()
	ut.PrintJson(games)

	// user := "123456"
	// pp.Regist(user)
	// pp.FundTransferIn(user, 1000000)
	// fmt.Println(pp.LaunchGame(user, "vs20olympx", "en"))

	// fmt.Println(strconv.FormatFloat(1000000.0, 'f', -1, 64))
}
