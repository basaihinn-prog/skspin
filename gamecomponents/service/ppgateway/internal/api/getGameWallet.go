package api

import (
	"game/comm/slotsmongo"
	"game/comm/ut"
)

// https://api.pg-demo.com/web-api/game-proxy/v2/GameWallet/Get?traceId=WBADFR10
// btt=1&wk=0_C&atk=0A2C1AE9-FAFD-4FCF-A907-A102FED5DC4E&pf=1&gid=39
// {"dt":{"cc":"THB","tb":115060.50,"pb":0.00,"cb":115060.50,"tbb":0.00,"tfgb":0.00,"rfgc":0,"inbe":false,"infge":false,"iebe":false,"iefge":false,"ch":{"k":"0_C","cid":0,"cb":115060.50},"p":null,"ocr":null},"err":null}
func getGameWallet(ps *PGParams, ret *M) (err error) {
	gold, err := slotsmongo.GetBalance(ps.Pid)
	if err != nil {
		return
	}

	balance := ut.Gold2Money(gold)

	*ret = M{"cc": "THB", "tb": balance, "pb": 0.00, "cb": balance, "tbb": 0.00, "tfgb": 0.00, "rfgc": 0, "inbe": false, "infge": false, "iebe": false, "iefge": false, "ch": M{"k": "0_C", "cid": 0, "cb": balance}, "p": nil, "ocr": nil}

	return
}
