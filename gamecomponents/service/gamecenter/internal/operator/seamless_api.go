package operator

import (
	"context"
	"game/comm"
	"game/comm/define"
	"game/comm/ut"
	"game/duck/lazy"
	"game/service/gamecenter/msgdef"
	"log/slog"
	"net/url"
	"time"
)

func IsTimeoutErr(err error) bool {
	if err == nil {
		return false
	}
	urlerr, ok := err.(*url.Error)
	return ok && urlerr.Timeout()
}

// func (h *MemApp) seamlessModifyGold(memplr *MemPlr, change int64, transId string) (int64, error) {
func (h *MemApp) seamlessModifyGold(memplr *MemPlr, transId string, ps *msgdef.ModifyGoldPs) (int64, error) {
	apiUrl, err := url.JoinPath(h.Address, "/Cash/TransferInOut")
	if err != nil {
		return 0, err
	}
	realAmount := ut.Gold2Money(ps.Change)
	item := lazy.GetCurrencyItemNew(h.CurrencyKey)
	if h.CurrencyCtrlStatus == 1 {
		realAmount = ut.Gold2Money(ps.Change) * float64(item.Multi)
	}

	req := map[string]any{
		"TransactionID": transId,
		"UserID":        memplr.Uid,
		"AppID":         h.AppID,
		"AppSecret":     h.AppSecret,
		"Amount":        ut.Gold2Money(ps.Change),
		"RoundID":       ps.RoundID,
		"GameID":        ps.GameID,
		"ReqTime":       time.Now(),
		"Reason":        ps.Reason,
		"IsEnd":         ps.IsEnd,
		"RealAmount":    realAmount,
	}

	var resp struct {
		Balance    float64
		RealAmount float64
	}

	var transfer = func() (retry bool) {
		err = comm.PostJsonCode(context.TODO(), apiUrl, req, &resp, nil)
		retry = IsTimeoutErr(err)
		return
	}

	ut.FundTransferWithRetryAndInterval(transfer, 10*time.Second, 5)

	// err = comm.PostJsonCode(context.TODO(), apiUrl, req, &resp, nil)
	if err != nil {
		return 0, err
	}
	if resp.RealAmount != 0 {
		if h.CurrencyCtrlStatus == 1 {
			if resp.RealAmount != realAmount {
				err = define.NewErrCode("Incorrect transaction amount", 1016)
				return 0, err
			}
		} else {
			if resp.RealAmount != realAmount*float64(item.Multi) {
				err = define.NewErrCode("Incorrect transaction amount", 1016)
				return 0, err
			}
		}
	}
	if resp.Balance != 0 {
		if h.CurrencyCtrlStatus == 1 {
			resp.Balance = resp.Balance / float64(item.Multi)
		}
	}

	// return resp.Balance, nil
	return ut.Money2Gold(resp.Balance), nil
}

func (h *MemApp) seamlessGetBalance(memplr *MemPlr) (int64, error) {
	lg := slog.With("memplr", memplr)

	var err error
	defer func() {
		lg.Info("seamlessGetBalance", "error", ut.ErrString(err))
	}()

	apiUrl, err := url.JoinPath(h.Address, "/Cash/Get")
	if err != nil {
		return 0, err
	}
	lg = lg.With("apiUrl", apiUrl)

	var resp struct {
		Balance float64
	}
	req := map[string]any{
		"UserID":    memplr.Uid,
		"AppID":     h.AppID,
		"AppSecret": h.AppSecret,
	}

	lg = lg.With("req", req)

	err = comm.PostJsonCode(context.TODO(), apiUrl, req, &resp, nil)
	if err != nil {
		return 0, err
	}
	if resp.Balance != 0 {
		if h.CurrencyCtrlStatus == 1 {
			item := lazy.GetCurrencyItemNew(h.CurrencyKey)
			resp.Balance = resp.Balance / float64(item.Multi)
		}
	}

	lg = lg.With("balance", resp.Balance)

	return ut.Money2Gold(resp.Balance), nil
}
