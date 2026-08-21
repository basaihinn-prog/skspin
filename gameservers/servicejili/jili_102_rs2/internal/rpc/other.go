package rpc

import (
	"fmt"
	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/servicejili/jili_102_rs2/internal"
	"serve/servicejili/jili_102_rs2/internal/message"
	"serve/servicejili/jili_102_rs2/internal/models"

	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc("/fulljp/JPInfoProto", unionjp_JPBlockProto)

	jiliut.RegRpc(fmt.Sprintf("/%s/account/heart", internal.GameShortName), accountHeart)

}

// /csh/account/heart
func accountHeart(ps *nats.Msg) (ret []byte, err error) {
	var ack = message.Rs2_HeartAck{
		// Aid: &aid,
		Message: proto.String("testtest"),
	}

	encode, _ := proto.Marshal(&ack)

	var resdata = message.Rs2_ResData{
		Type: proto.Int32(AckType["heartBeat"]),
		Data: []*message.Rs2_InfoData{
			{
				Encode: encode,
			},
		},
	}
	// resdata.

	ret, _ = proto.Marshal(&resdata)
	return
}

func unionjp_JPBlockProto(ps *nats.Msg) (ret []byte, err error) {
	var preq message.Gaia_FullJPInfoReq

	err = proto.Unmarshal(ps.Data, &preq)
	if err != nil {
		return
	}

	token := preq.GetToken()
	pid, err := jwtutil.ParseToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {
		aid := int32(pid)
		var ack = message.Rs2_LoginDataAck{
			Aid: &aid,
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Rs2_ResData{
			Type: proto.Int32(AckType["login"]),
			Data: []*message.Rs2_InfoData{
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
