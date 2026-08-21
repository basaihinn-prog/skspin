package rpc

import (
	"fmt"
	"strconv"

	"serve/servicejili/jili_137_ge/internal"

	"time"

	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_137_ge/internal/message"
	"serve/servicejili/jili_137_ge/internal/models"
	"serve/servicejili/jiliut"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/account/exchangechips", internal.GameShortName), exchangechips)
}

// /csh/account/exchangechips
func exchangechips(ps *nats.Msg) (ret []byte, err error) {
	var req message.Ge_ExchangeReq
	err = proto.Unmarshal(ps.Data, &req)
	if err != nil {
		return
	}

	token := req.GetToken()
	// slog.Info("login", "token", token)
	pid, err := jwtutil.ParseToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	gold, err := slotsmongo.GetBalance(pid)
	if err != nil {
		return
	}

	err = db.CallWithPlayer(pid, func(plr *models.Player) error {

		stamp := time.Now().UnixMilli()
		var ack = message.Ge_ExchangeAck{
			Stamp: proto.String(strconv.Itoa(int(stamp))),
			Wallet: []*message.Ge_Wallet{
				{
					Coin:  proto.Float64(ut.Gold2Money(gold)),
					Ratio: proto.Float64(1),
					Rate:  proto.Float64(1),
					Unit:  proto.Float64(1),
				},
			},
		}
		encode, _ := proto.Marshal(&ack)

		var resData = message.Ge_ResData{
			Ret:  proto.Int32(0),
			Type: proto.Int32(AckType["exchangeChips"]),
			Data: []*message.Ge_InfoData{
				{
					Encode: encode,
				},
			},
		}

		ret, _ = proto.Marshal(&resData)
		return nil
	})
	return
}
