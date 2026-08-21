package operator

import (
	"errors"
	"game/comm"
	"game/comm/ut"
	"game/duck/ut2"
	"game/service/gamecenter/msgdef"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemApp struct {
	AppID     string `bson:"_id"`
	AppSecret string
	Address   string // 请求地址
	// Comment    string // 备注
	WalletMode     int
	PublishHistory bool
	Status         int
	CurrencyKey    string

	LastPullHisTime time.Time

	// api     *HttpApi
	regLock sync.Mutex

	// uid-> plr
	players *ut2.SyncMap[string, *MemPlr]
}

// 来自内部游戏服的rpc调用
func (memapp *MemApp) GetBalance(plr *MemPlr) (balance int64, err error) {
	switch memapp.WalletMode {
	case comm.WalletModeOld:
		balance, err = NewHttpApi(memapp.Address).GetBalance(plr)
	case comm.WalletModeTransfer:
		balance, err = plr.Balance()
	case comm.WalletModeSeamless:
		balance, err = memapp.seamlessGetBalance(plr)
	default:
		err = errors.New("Not Support Wallet Mode!")
	}

	return
}

// 来自内部游戏服的rpc调用
// func (memapp *MemApp) ModifyGold(plr *MemPlr, change int64, comment string) (balance int64, err error) {
func (memapp *MemApp) ModifyGold(plr *MemPlr, ps *msgdef.ModifyGoldPs) (balance int64, err error) {
	// before := time.Now()
	logid := primitive.NewObjectID()
	change := ps.Change
	comment := ps.Comment

	t0 := time.Now()

	switch memapp.WalletMode {
	case comm.WalletModeOld:
		// balance, err = memapp.api.ModifyGold(plr, change, comment)
		balance, err = NewHttpApi(memapp.Address).ModifyGold(plr, ps.Change, comment)
	case comm.WalletModeTransfer:
		if change > 0 {
			balance, err = plr.TransferIn(change)
		} else {
			balance, err = plr.TransferOut(-change)
		}
	case comm.WalletModeSeamless:
		// balance, err = memapp.seamlessModifyGold(plr, change, logid.Hex())
		balance, err = memapp.seamlessModifyGold(plr, logid.Hex(), ps)
	default:
		err = errors.New("Not Support Wallet Mode")
	}

	// InsertModifyLog(logid, plr, change, comment, before, balance, err)
	now := time.Now()
	InsertModifyLog(&DocModifyGoldLog{
		ID:         logid,
		InsertTime: now,
		ElapsedMS:  now.Sub(t0) / time.Millisecond,
		Change:     change,
		Error:      ut.ErrString(err),
		Balance:    balance,
		Action:     "ModifyGold",
		ReqData:    ps,
	}, plr)

	return
}
