package rpc

import (
	"fmt"
	"log/slog"

	"serve/servicejili/jili_45_cbt/internal"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/servicejili/jili_45_cbt/internal/message"
	"serve/servicejili/jili_45_cbt/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/account/login", internal.GameShortName), login)
}

func login(ps *nats.Msg) (ret []byte, err error) {
	var login message.Cbt_LoginDataReq
	err = proto.Unmarshal(ps.Data, &login)
	if err != nil {
		return
	}

	// fmt.Println(login.String())

	token := login.GetToken()
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
		var ack = message.Cbt_LoginDataAck{
			Aid: &aid,
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Cbt_ResData{
			Type: proto.Int32(AckType["login"]),
			Data: []*message.Cbt_InfoData{
				{
					Encode: encode,
				},
			},
		}
		// resdata.

		ret, _ = proto.Marshal(&resdata)

		return err
	})
	return
}
