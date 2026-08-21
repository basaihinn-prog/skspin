package rpc

import (
	"fmt"
	"log/slog"

	"serve/servicejili/jili_225_ctk/internal"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/servicejili/jili_225_ctk/internal/message"
	"serve/servicejili/jili_225_ctk/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/account/login", internal.GameShortName), login)
}

func login(ps *nats.Msg) (ret []byte, err error) {
	var login message.Ctk_LoginDataReq
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
		var ack = message.Ctk_LoginDataAck{
			Aid: &aid, //1234843
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Ctk_ResData{
			Type: proto.Int32(AckType["login"]),
			Ret:  proto.Int32(0),
			Data: []*message.Ctk_InfoData{
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
