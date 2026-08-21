package rpc

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"serve/servicejili/jili_33_pp/internal"
	tmp "serve/servicejili/jili_33_pp/internal/message"
	"serve/servicejili/jiliut"
)

func init() {
	jiliut.RegRpc(fmt.Sprintf("/%s/jpinfo", internal.GameShortName), jpinfo)
}

func jpinfo(ps *nats.Msg) (ret []byte, err error) {

	var req serverOfficial.Request

	err = proto.Unmarshal(ps.Data, &req)
	if err != nil {
		return
	}

	token := ps.Header.Get("Token")

	var (
		retData protoreflect.ProtoMessage
	)

	retData = &tmp.JPInfo{
		Fronts:   []float64{7.12, 10.24, 26.34, 33.44, 53.54},
		OddPrize: []float64{15, 40, 120, 250, 1000},
		BaseBet:  100,
	}
	if err != nil {
		return
	}
	data, err := jiliut.ProtoEncryption(token, retData)
	if err != nil {
		return
	}
	var resp = &serverOfficial.GaiaResponse{
		Type: req.GetAck(),
		Data: data,
	}

	ret = jiliut.ProtoEncode(resp)
	return
}
