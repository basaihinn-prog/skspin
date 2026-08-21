package rpc

import (
	"fmt"
	"log/slog"
	"serve/servicejili/jili_2_csh/internal"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/servicejili/jili_2_csh/internal/message"
	"serve/servicejili/jili_2_csh/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/account/login", internal.GameShortName), login)
}

func login(ps *nats.Msg) (ret []byte, err error) {
	var login message.Csh_LoginDataReq
	err = proto.Unmarshal(ps.Data, &login)
	if err != nil {
		slog.Error("login::Unmarshal", "err", err)
		return
	}

	// fmt.Println(login.String())

	token := login.GetToken()
	slog.Info("login", "token", token)

	//TODO
	var pid int64
	pid, err = jwtutil.ParseToken(token)
	if err != nil {
		slog.Error("login", "token", token, "err", err)
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		aid := int32(pid)
		var ack = message.Csh_LoginDataAck{
			Aid: &aid,
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Csh_ResData{
			Type: proto.Int32(AckType["login"]),
			Data: []*message.Csh_InfoData{
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
