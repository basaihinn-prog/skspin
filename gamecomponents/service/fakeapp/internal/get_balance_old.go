package internal

/*
import (
	"errors"
	"game/comm/mux"
	"net/http"
)

func init() {
	mux.RegHttpWithSample("/GetBalance", "get player balace", "player", getBalance, getBalancePs{"testuser1"})
	mux.RegHttpWithSample("/ModifyGold", "modify player balace", "player", modifyGold, modifyGoldPs{"testuser1", 10000})
}

type getBalancePs struct {
	UserID string
}

type getBalanceRet struct {
	Balance int64
}

func getBalance(_ *http.Request, ps getBalancePs, ret *getBalanceRet) (err error) {
	ensureUserExists(ps.UserID)

	gold, ok := userMap.Load(ps.UserID)
	if !ok {
		return errors.New("用户不存在")
	}
	ret.Balance = gold
	return
}

type modifyGoldPs struct {
	UserID string // 玩家ID
	Change int64  // 金币变化：>0 加钱， <0 扣钱
}

func modifyGold(_ *http.Request, ps modifyGoldPs, ret *getBalanceRet) (err error) {
	ensureUserExists(ps.UserID)

	gold, ok := userMap.Load(ps.UserID)
	if !ok {
		return errors.New("用户不存在")
	}
	gold += ps.Change
	if gold < 0 {
		return errors.New("余额不足")
	}
	userMap.Store(ps.UserID, gold)

	ret.Balance = gold
	return
}
*/
