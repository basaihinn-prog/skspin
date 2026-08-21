package rpc

import (
	"log/slog"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/servicejili/jili_166_utaxi/internal/message"
	"serve/servicejili/jili_166_utaxi/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc("/utaxi/account/login", login)
}

func login(ps *nats.Msg) (ret []byte, err error) {
	var login message.Utaxi_LoginDataReq
	err = proto.Unmarshal(ps.Data, &login)
	if err != nil {
		return
	}

	// fmt.Println(login.String())

	token := login.GetToken()
	slog.Info("login", "token", token)

	//TODO
	pid := int64(123456)

	pid, err = jwtutil.ParseToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		aid := int32(pid)
		var ack = message.Utaxi_LoginDataAck{
			Aid: &aid,
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Utaxi_ResData{
			Type: proto.Int32(AckType["login"]),
			Data: []*message.Utaxi_InfoData{
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
