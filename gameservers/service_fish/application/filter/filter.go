package filter

import (
	"errors"
	"fmt"
	common_proto "serve/fish_comm/common/proto"
	heartbeat "serve/fish_comm/heart-beat"
	"serve/fish_comm/jackpot-client/jackpot"
	"serve/service_fish/application/gamesetting"
	"serve/service_fish/application/lobbysetting"
	"serve/service_fish/domain/auth"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	common_fish_proto "serve/service_fish/models/proto"

	"github.com/gogo/protobuf/proto"
)

func parseHeader(msg []byte) *common_proto.Header {
	header := &common_proto.Header{}

	if err := proto.Unmarshal(msg, header); err != nil {
		return nil
	}
	return header
}

func parseHeaderFish(msg []byte) *common_fish_proto.Header {
	header := &common_fish_proto.Header{}

	if err := proto.Unmarshal(msg, header); err != nil {
		return nil
	}
	return header
}

func FindActionName(msg []byte) (string, error) {
	if header := parseHeader(msg); header != nil {
		switch header.Msgid {
		case common_proto.EMSGID_eLoginCall:
			return auth.EMSGID_eLoginCall, nil

		case common_proto.EMSGID_eResultCall:
			return probability.EMSGID_eResultCall, nil

		case common_proto.EMSGID_eMultiResultCall:
			return probability.EMSGID_eMultiResultCall, nil

		case common_proto.EMSGID_eLobbyConfigCall:
			return lobbysetting.EMSGID_eLobbyConfigCall, nil

		case common_proto.EMSGID_eGameConfigCall:
			return gamesetting.EMSGID_eGameConfigCall, nil

		case common_proto.EMSGID_eGameStripsCall:
			return gamesetting.EMSGID_eGameStripsCall, nil

		case common_proto.EMSGID_eOptionCall:
			return probability.EMSGID_eOptionCall, nil

		case common_proto.EMSGID_eHeartbeatCall:
			return heartbeat.ActionHeartBeatCall, nil

		case common_proto.EMSGID_eCentInReask:
			return jackpot.EMSGID_eCentInReask, nil

		}
	}

	if headerFish := parseHeaderFish(msg); headerFish != nil {
		switch headerFish.Msgid {

		case common_fish_proto.EMSGID_eShoot:
			return bullet.EMSGID_eShoot, nil

		case common_fish_proto.EMSGID_eBetChange:
			return bullet.EMSGID_eBetChange, nil

		case common_fish_proto.EMSGID_eAllFishCall:
			return fish.EMSGID_eAllFishCall, nil

		case common_fish_proto.EMSGID_eMercenaryOpenCall:
			return probability.EMSGID_eMercenaryOpenCall, nil

		default:
			return "", errors.New(fmt.Sprintf("Unknown MSGID : %s", headerFish.Msgid))
		}
	}

	return "", errors.New(fmt.Sprintf("Unknown header format %T.", msg))
}
