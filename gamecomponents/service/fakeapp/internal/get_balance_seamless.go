package internal

import (
	"encoding/json"
	"errors"
	"game/comm/mux"
	"net/http"
	"os"
	"time"
)

func init() {
	// mux.RegHttpWithSample("/GetBalance", "get player balace", "player", getBalance, getBalancePs{"testuser1"})
	mux.RegHttpWithSample("/Cash/Get", "get player balace", "player", cashGet, cashGetPs{"testuser1", appId, appSecret})
	mux.RegHttpWithSample("/Cash/TransferInOut", "modify player balance", "player", cashTrans, cashTransPs{
		UserID:        "testuser1",
		Amount:        10000,
		AppID:         appId,
		AppSecret:     appSecret,
		TransactionID: "abc1234",
	})

	mux.RegHttpWithSample("/History/Spin", "accept history spin", "player", historySpin, nil)

}

func historySpin(_ *http.Request, ps json.RawMessage, _ *mux.EmptyResult) (err error) {
	os.Stdout.Write(ps)
	os.Stdout.Write([]byte("\n"))
	return
}

type cashGetPs struct {
	UserID    string
	AppID     string
	AppSecret string
}

type cashGetRet struct {
	Balance float64
}

func cashGet(_ *http.Request, ps cashGetPs, ret *cashGetRet) (err error) {
	if ps.AppID != appId || ps.AppSecret != appSecret {
		err = errors.New("Invalid App Params")
		return
	}
	ensureUserExists(ps.UserID)

	gold, ok := userMap.Load(ps.UserID)
	if !ok {
		return errors.New("用户不存在")
	}
	ret.Balance = gold
	return
}

type cashTransPs struct {
	UserID        string // 玩家ID
	AppID         string
	AppSecret     string
	TransactionID string
	Amount        float64
	RoundID       string
	GameID        string
	ReqTime       time.Time
	Reason        string
	IsEnd         bool
}

func cashTrans(_ *http.Request, ps cashTransPs, ret *cashGetRet) (err error) {
	if ps.AppID != appId || ps.AppSecret != appSecret {
		err = errors.New("Invalid App Params")
		return
	}
	ensureUserExists(ps.UserID)

	gold, ok := userMap.Load(ps.UserID)
	if !ok {
		return errors.New("用户不存在")
	}
	gold += ps.Amount
	if gold < 0 {
		return errors.New("余额不足")
	}
	userMap.Store(ps.UserID, gold)

	ret.Balance = gold
	return
}
