package pp

import (
	"strconv"

	"serve/comm/plats/platcomm"

	"github.com/google/uuid"
)

type transferRet struct {
	Balance float64
}

func trans(username string, orderID string, amount float64) (balance float64, err error) {
	var ret transferRet
	err = invoke("/balance/transfer/", map[string]string{
		"externalPlayerId":      username,
		"externalTransactionId": orderID,
		"amount":                strconv.FormatFloat(amount, 'f', -1, 64),
	}, &ret)

	balance = ret.Balance
	return
}

func (_ pp) FundTransferIn(username string, amount float64) (status string) {
	// username := strconv.Itoa(uid)

	orderID := uuid.NewString()
	_, err := trans(username, orderID, amount)
	status = platcomm.GetTransStatus(err)
	return

}

func (this pp) FundTransferOut(username string) (amount float64, status string) {
	invoke("/game/session/terminate/", map[string]string{
		// "externalTransactionId": orderID,
		"externalPlayerId": username,
	}, nil)

	amount, err := this.GetBalance(username)
	if err != nil {
		status = "ERROR:" + err.Error()
		return
	}
	if amount <= 0 {
		return
	}

	orderID := uuid.NewString()
	_, err = trans(username, orderID, -amount)
	status = platcomm.GetTransStatus(err)
	return
}
