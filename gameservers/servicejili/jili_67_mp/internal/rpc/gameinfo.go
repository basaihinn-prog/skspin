package rpc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/jwtutil"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_67_mp/internal"
	"serve/servicejili/jili_67_mp/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/game/info", internal.GameShortName), gameinfo)
}

var BetGrade = []float64{1, 2, 3, 5, 8, 10, 20, 50, 100, 200, 300, 400, 500, 700, 1000}

func gameinfo(ps *nats.Msg) (ret []byte, err error) {
	query := lo.Must(url.ParseQuery(ps.Header.Get("query")))
	token := query.Get("token")
	aid := query.Get("aid")
	pid, err := jwtutil.ParseToken(token)
	if err != nil {
		err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	gold, err := slotsmongo.GetBalance(pid)
	if err != nil {
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		info, err := redisx.LoadAppIdCache(plr.AppID)
		if err != nil {
			return err
		}
		curItem := lazy.GetCurrencyItem(plr.CurrencyKey)
		resdata := map[string]interface{}{
			"info": map[string]interface{}{
				"ApiType":        0,
				"ArcadeNo":       aid,
				"CanIntoArcade":  true,
				"FreeSpinRemain": 0,
				"Prefer":         nil,
				"Test":           true,
				"WalletInfo": []interface{}{
					map[string]interface{}{
						"currencyNumber": 0,
						"currencyName":   "",
						"currencySymbol": "",
						"Exponent":       4,
						"coin":           ut.Gold2Money(gold),
						"bet":            ut.FloatArrMul(info.Cs, curItem.Multi),
						"unit":           1,
						"ratio":          1,
						"rate":           1,
						"upper":          0,
						"lower":          0,
						"cycle":          0,
					},
				},
			},
			"ret":   0,
			"token": token,
			"type":  11,
		}
		ret, _ = json.Marshal(resdata)
		plr.SpinCountOfThisEnter = 0
		return nil
	})
	return
}
