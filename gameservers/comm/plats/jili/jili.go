package jili

import (
	"slices"
	"strconv"

	"serve/comm/define"
	"serve/comm/plats/platcomm"

	"github.com/google/uuid"
)

type jili struct{}

func init() {
	// regPlat("pg", pg{})
	platcomm.Plats["JILI"] = jili{}
}

func JILI() jili {
	return jili{}
}

func (jl jili) LaunchGame(username, game, lang string, useProxy bool) (url string, err error) {
	gameid, err := strconv.Atoi(game)
	if err != nil {
		return
	}
	err = jl.Regist(username)
	if err != nil {
		return
	}
	balance, _ := jl.GetBalance(username)
	if balance < 10000 {
		jl.FundTransferIn(username, 100e4)
	}

	err = invoke("/LoginWithoutRedirect", D(
		"Account", username,
		"GameId", gameid,
		"Lang", lang,
	), &url)

	if err != nil {
		return
	}
	// url = strings.Replace(url, "https://uat-wbgame.jlfafafa3.com", "https://jilid-rslotszs001.kafa010.com", 1)

	return
}

func (jili) GetBalance(username string) (balance float64, err error) {
	// {"ErrorCode":0,"Message":"","Data":[{"Account":"1159574","Balance":0,"Status":2}]}

	var ret []struct {
		Account string
		Balance float64
		Status  int
	}
	err = invoke("/GetMemberInfo", D(
		"Accounts", username,
	), &ret)
	if err != nil {
		return
	}

	balance = ret[0].Balance
	return
}

type getGamesResult = []GameInfo

// Jili>>>返回:{"ErrorCode":0,"Message":"","Data":[{"GameId":1,"name":{"en-US":"Royal Fishing","
// zh-CN":"\u94b1\u9f99\u6355\u9c7c","zh-TW":"\u9322\u9f8d\u6355\u9b5a"},"GameCategoryId":5},{"G
// ameId":2,"name":{"en-US":"Chin Shi Huang","zh-CN":"\u79e6\u7687\u4f20\u8bf4","zh-TW":"\u79e6\
// u7687\u50b3\u8aaa"},"GameCategoryId":1},{"GameId":4,"name":{"en-US":"God Of Martial","zh-CN":
// "\u5173\u4e91\u957f","zh-TW":"\u95dc\u96f2\u9577"},"GameCategoryId":1},{"GameId":5,"name":{"e
// n-US":"Hot Chilli","zh-CN":"\u706b\u70ed\u8fa3\u6912","zh-T +[6498]more
type GameInfo struct {
	GameId         int
	Name           map[string]string
	GameCategoryId int
}

func (jili) GetGameList() (games platcomm.HotGames, err error) {
	gameType := "RNG"
	var category []int
	switch gameType {
	case "RNG":
		category = []int{1}
	case "FISH":
		category = []int{5}
	case "PVP":
		category = []int{2, 8}
	default:
		return
	}

	var li getGamesResult
	err = invoke("/GetGameList", D(), &li)
	// err = get_games(nil, D(), &li)
	if err != nil {
		return
	}

	// http://a.kky188.com:8082/icon/goldenf/PG0052.png
	// GameID_20_EN.png

	for _, v := range li {
		// 1: 电子
		// 5: 捕鱼
		if slices.Contains(category, v.GameCategoryId) {
			game := &platcomm.HotGame{
				Plat: "JILI",
				ID:   strconv.Itoa(v.GameId),
				Name: v.Name["en-US"],
				Type: define.GameType_Slot,
			}

			if game.Name == "Dice" {
				continue
			}

			games = append(games, game)
		}
	}

	// gamedata.Data().GamesSort.Sort("JILI", games)

	return
}

func (jili) Regist(username string) (err error) {
	err = invoke("/CreateMember", D("Account", username), nil)
	if define.CodeErrorEq(err, 101) {
		err = nil
	}
	return
}

type transResult struct {
	TransactionId  string
	CoinBefore     float64
	CoinAfter      float64
	CurrencyBefore float64
	CurrencyAfter  float64
	Status         int
}

func (jili) FundTransferIn(username string, amount float64) (status string) {
	status = "SUCCESS"

	orderID := uuid.NewString()
	// /ExchangeTransferByAgentId

	var result transResult
	err := invoke("/ExchangeTransferByAgentId", D(
		"Account", username,
		"TransactionId", orderID,
		"Amount", amount,
		"TransferType", 2,
	), &result)

	status = platcomm.GetTransStatus(err)

	return
}

func (jili) FundTransferOut(uid string) (amount float64, status string) {
	status = "Fail: the method not implemented!!"
	return
}
