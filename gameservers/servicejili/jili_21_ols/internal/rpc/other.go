package rpc

import (
	"serve/servicejili/jiliut/AckType"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	reqMux[AckType.JpShowData] = jpShowData
	reqMux[AckType.JpInfo] = jpInfo
	reqMux[AckType.FullJpInfo] = fullJpInfo
	reqMux[AckType.FullJpInfoAll] = fullJpInfoAll
	reqMux[AckType.VipSignInfo] = vipSignInfo
	reqMux[AckType.GetBackPack] = getBackPack
	reqMux[AckType.GetAllItemCanUse] = getAllItemCanUse
	reqMux[AckType.GetMail] = getMail
	reqMux[AckType.DebrisInfo] = debrisInfo
	reqMux[AckType.GetNowMission] = getNowMission
	reqMux[AckType.Notice] = notice
}

func jpShowData(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}

func jpInfo(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}

func fullJpInfo(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}

func fullJpInfoAll(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}

func vipSignInfo(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}

func getBackPack(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}

func getAllItemCanUse(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}

func getMail(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}

func debrisInfo(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}

func getNowMission(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}
func notice(pid int64, data []byte, ps *nats.Msg) (ret protoreflect.ProtoMessage, err error) {
	return
}
