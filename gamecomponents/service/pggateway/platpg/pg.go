package platpg

import (
	"fmt"
	"game/comm/define"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type pg struct{}

func init() {
	regPlat("pg", pg{})
}

func PG() IPlat {
	return pg{}
}

func (_ pg) Regist(uid string) (err error) {
	username := uid

	var result struct {
		Action_result int
	}
	err = invoke("/Player/v1/Create", define.M{
		"player_name": username,
		"nickname":    username,
		"currency":    "THB",
	}, &result)

	if err != nil {
		if define.CodeErrorEq(err, 1305) {
			err = nil
		}
		return
	}

	if result.Action_result == 1 {
		return nil
	}

	return fmt.Errorf("Action_result %d", result.Action_result)
}

func (_ pg) GetBalance(uid string) (balance float64, err error) {
	username := uid
	var result struct {
		CashBalance float64
	}
	err = invoke("/Cash/v3/GetPlayerWallet", define.M{
		"player_name": username,
	}, &result)
	if err == nil {
		balance = result.CashBalance
	}
	return
}

type game struct {
	GameId        int
	GameCode      string
	Category      int
	GameName      string
	ReleaseStatus int
	Status        int
}

type HotGames []*HotGame

type HotGame struct {
	Plat string
	ID   string
	Name string
	Type int
	Icon string
}

func (pg) GetGameList() (games HotGames, err error) {

	var gamelist []game
	err = invoke("Game/v2/Get", define.M{
		"currency": "THB",
		"language": "en-us",
		"status":   1,
	}, &gamelist)
	if err != nil {
		return
	}

	games = make(HotGames, 0, len(gamelist))

	for _, v := range gamelist {
		if v.Category == 1 && v.ReleaseStatus == 1 && v.Status == 1 {
			hg := &HotGame{
				Plat: "PG",
				ID:   strconv.Itoa(v.GameId),
				Name: v.GameName,
			}
			games = append(games, hg)
		}
	}

	return
}

func (p pg) LaunchGame(uid string, game, lang string) (url_ string, err error) {
	// return p.LaunchGameHtml(uid, game, lang)
	return p.LaunchGame_old(uid, game, lang)
}

func (pg) LaunchGame_old(uid string, game, lang string) (url_ string, err error) {
	cfg := GetConfig()
	// if err != nil {
	// 	return
	// }
	ul, err := url.Parse(cfg.LaunchURL)
	if err != nil {
		return
	}

	ul.Path = path.Join(ul.Path, game, "index.html")

	// token := "doudoutoken" + strconv.Itoa(uid)
	token, err := GenToken(uid)

	ps := url.Values{}

	ps.Set("l", lang)
	ps.Set("btt", "1")
	ps.Set("ot", cfg.OperatorToken)
	ps.Set("ops", token)
	ps.Set("or", "static-pg.kafa010.com")
	// ps.Set("or", "static.pgsoft-games.com")
	// ps.Set("__refer", "m.pg-redirect.net")

	// __refer=m.pg-redirect.net
	// or=static-pg.kafa010.com

	ul.RawQuery = ps.Encode()
	url_ = ul.String()

	// https://public.pg-redirect.net
	// 3.2 Launch URL (URL Scheme)
	// https://m.pg-redirect.net/{game_code}/index.html?language={0}&bet_type=1&operator_token={2}&operator_player_session={3}

	return
}

/*
5.2.4 充值与转出流程
PG 系统正在为所有 API 操作实现幂等操作。根据以下情况，要求运营商通过使用相同的
transfer_reference 重新发送请求：
• 从 PG 系统中收到错误的返回（错误的返回格式）。
• 收到除了 3001、3005、3100 的错误返回代码。
• 没有收到 PG 系统的返回
*/
func (pg) FundTransferIn(uid string, amount float64) (status string) {
	// username := strconv.Itoa(uid)
	orderID := uuid.NewString()

	var transfer = func() (retry bool) {
		var result struct {
			TransactionId       int
			BalanceAmountBefore float64
			BalanceAmount       float64
			Amount              float64
		}
		err := invoke("/Cash/v3/TransferIn", define.M{
			"player_name":        uid,
			"amount":             amount,
			"transfer_reference": orderID,
			"currency":           "THB",
		}, &result)

		status = GetTransStatus(err)

		retry = IsTimeoutErr(err) || define.CodeErrorEq(err, 3101)
		return
	}

	FundTransferWithRetry(transfer, transfer)
	return
}

func (pg) FundTransferOut(uid string) (amount float64, status string) {
	username := uid

	orderID := uuid.NewString()
	var kickRet struct {
		Action_result int
	}
	invoke("/Player/v1/Kick", define.M{
		"player_name": username,
	}, &kickRet)

	if kickRet.Action_result == 1 {
		time.Sleep(5 * time.Second)
	}

	var transfer = func() (retry bool) {
		var result struct {
			TransactionId       int
			BalanceAmountBefore float64
			BalanceAmount       float64
			Amount              float64
		}
		err := invoke("/Cash/v3/TransferAllOut", define.M{
			"player_name":        username,
			"transfer_reference": orderID,
			"currency":           "THB",
		}, &result)

		amount = result.Amount
		status = GetTransStatus(err)

		retry = IsTimeoutErr(err) || define.CodeErrorEq(err, 3101)
		return
	}
	FundTransferWithRetry(transfer, transfer)
	return
}
