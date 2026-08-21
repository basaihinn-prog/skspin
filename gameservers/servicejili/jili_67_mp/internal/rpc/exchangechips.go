package rpc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/jwtutil"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_67_mp/internal"
	"serve/servicejili/jili_67_mp/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/account/exchangechips", internal.GameShortName), exchangechips)
}

// /csh/account/exchangechips
func exchangechips(ps *nats.Msg) (ret []byte, err error) {
	query := lo.Must(url.ParseQuery(ps.Header.Get("query")))
	token := query.Get("token")
	stamp := query.Get("stamp")
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
		resdata := map[string]interface{}{
			"info": []map[string]interface{}{
				{
					"Exponent":       4,
					"Symbol":         "",
					"coin":           ut.Gold2Money(gold),
					"currencyNumber": 0,
					"rate":           1,
					"ratio":          1,
					"unit":           1,
				},
			},
			"ret":   0,
			"stamp": stamp,
			"type":  2,
		}
		ret, _ = json.Marshal(resdata)
		return nil
	})
	return
}
