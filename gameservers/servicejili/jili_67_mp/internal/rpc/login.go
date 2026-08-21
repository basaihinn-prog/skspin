package rpc

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/servicejili/jili_67_mp/internal"
	"serve/servicejili/jili_67_mp/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/account/login", internal.GameShortName), login)
}

func login(ps *nats.Msg) (ret []byte, err error) {
	var loginPs map[string]interface{}
	err = json.Unmarshal(ps.Data, &loginPs)
	if err != nil {
		return
	}

	// fmt.Println(login.String())

	token := loginPs["token"].(string)
	slog.Info("login", "token", token)

	//TODO
	var pid int64
	pid, err = jwtutil.ParseToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		aid := int32(pid)

		resdata := map[string]interface{}{
			"info": map[string]interface{}{
				"aid":    aid,
				"result": 0,
			},
			"ret":   0,
			"token": token,
			"type":  0,
		}
		// resdata.

		ret, _ = json.Marshal(&resdata)

		return err
	})
	return
}
