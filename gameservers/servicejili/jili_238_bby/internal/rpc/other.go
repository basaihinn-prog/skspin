package rpc

import (
	"fmt"
	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/servicejili/jili_238_bby/internal"
	"serve/servicejili/jili_238_bby/internal/message"
	"serve/servicejili/jili_238_bby/internal/models"

	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/mall/checkmallproto", internal.GameShortName), checkmallproto)

	jiliut.RegRpc(fmt.Sprintf("/%s/fulljp/JPInfoProto", internal.GameShortName), unionjp_JPInfoProto)

	jiliut.RegRpc(fmt.Sprintf("/%s/fulljp/JPInfoAllProto", internal.GameShortName), unionjp_JPBlockProto)

	jiliut.RegRpc(fmt.Sprintf("/%s/account/heart", internal.GameShortName), accountHeart)

	jiliut.RegRpc(fmt.Sprintf("/%s/unionjp/JPBlockProto", internal.GameShortName), unionjp_JPBlockProto)

}

// /csh/account/heart
func accountHeart(ps *nats.Msg) (ret []byte, err error) {
	var ack = message.Bb_HeartAck{
		// Aid: &aid,
		Message: proto.String("testtest"),
	}

	encode, _ := proto.Marshal(&ack)

	var resdata = message.Bb_ResData{
		Type: proto.Int32(AckType["heartBeat"]),
		Data: []*message.Bb_InfoData{
			{
				Encode: encode,
			},
		},
	}
	// resdata.

	ret, _ = proto.Marshal(&resdata)
	return
}

// /csh/mall/checkmallproto
func checkmallproto(ps *nats.Msg) (ret []byte, err error) {
	return
}

func unionjp_JPInfoProto(ps *nats.Msg) (ret []byte, err error) {
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
		zero := 0.0
		zi := int32(0)
		var ack = message.Gaia_FullJPInfoAck{
			Value:  &zero,
			Full:   &zero,
			Minvip: &zi,
			Minbet: &zero,
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Bb_ResData{
			Type: proto.Int32(AckType["login"]),
			Data: []*message.Bb_InfoData{
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
		temp := make([]float64, 0)
		var ack = message.Gaia_UnionJPBlockAck{
			List: temp,
		}

		encode, _ := proto.Marshal(&ack)

		var resdata = message.Bb_ResData{
			Type: proto.Int32(AckType["login"]),
			Data: []*message.Bb_InfoData{
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
