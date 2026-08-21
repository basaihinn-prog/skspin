package broadcaster

import (
	"fmt"
	"serve/fish_comm/broadcaster"
	"serve/fish_comm/broadcaster/proto"
	common_proto "serve/fish_comm/common/proto"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/format"
	"serve/fish_comm/logger"
	"serve/service_fish/application/gameuser"
	"serve/service_fish/models"

	"github.com/gogo/protobuf/proto"
)

var Service = &service{
	Id: "BroadcasterServiceFish",
}

type service struct {
	Id string
}

func (s *service) BigWinHandler(action *flux.Action, allGames map[string]*broadcaster.RegisterBroadcaster, allGamesConfig map[string]*broadcaster.ConfigBroadcaster) {
	b := action.Payload()[0].(*broadcaster.Broadcaster)

	checkBigShow := false
	if len(action.Payload()) > 1 {
		checkBigShow = action.Payload()[1].(bool)
	}

	if config, ok := allGamesConfig[b.Id()]; ok {
		if (b.Pay >= broadcaster.BigWinPay && b.Cent >= config.MinCent() && b.Cent > 0) || checkBigShow {
			data := &broadcaster_proto.Broadcaster{
				Msgid:     common_proto.EMSGID_eBroadcaster,
				BcType:    broadcaster_proto.BROADCASTER_TYPE_BIG_WIN,
				GameId:    b.GameId,
				SubgameId: uint32(b.SubGameId),
				MemberId:  format.HideWithStar(b.MemberId),
				Message:   fmt.Sprintf("%d/%d/%d", b.Pay, b.FishId, b.Cent),
				Time:      2,
			}

			dataBuffer, _ := proto.Marshal(data)

			for k, v := range allGames {
				switch {
				case v.Rate == 0 && v.GameId == b.GameId:
					// Lobby Room Rate is 0
					data.DvType = broadcaster_proto.DELIVER_TYPE_LOBBY
					flux.Send(broadcaster.EMSGID_eBroadcaster, s.Id, k, dataBuffer)

				// Process PSF-ON-00003 Broadcaster All Room.
				case b.Id() == v.Id() && v.GameId == models.PSF_ON_00003:
					data.DvType = broadcaster_proto.DELIVER_TYPE_GAME
					flux.Send(broadcaster.EMSGID_eBroadcaster, s.Id, k, dataBuffer)

				case b.Id() == v.Id() && b.Rate >= v.Rate:
					data.DvType = broadcaster_proto.DELIVER_TYPE_GAME
					flux.Send(broadcaster.EMSGID_eBroadcaster, s.Id, k, dataBuffer)
				}
			}
		}
	} else {
		logger.Service.Zap.Errorw(Broadcaster_CONFIG_NOT_FOUND,
			"GameUser", b.SecWebSocketKey,
			"GameId", b.GameId,
			"BetList", b.BetList,
			"RateList", b.RateList,
			"Rate", b.Rate,
		)
		errorcode.Service.Fatal(b.SecWebSocketKey, Broadcaster_CONFIG_NOT_FOUND)
	}
}

func (s *service) JackpotHandler(action *flux.Action, allGames map[string]*broadcaster.RegisterBroadcaster) {
	msgId := action.Payload()[0].(string)
	dataBuffer := action.Payload()[1]
	groupId := action.Payload()[2].(int32)

	for _, v := range allGames {
		switch {
		case v.GameId == models.PSF_ON_20002:
			// Broadcast All GameUsers, Not By Use Lobby room and Game room
			allGameUsers := gameuser.Service.GetAllGameUsers()

			for userKey, userValue := range allGameUsers {
				roomUuid := gameuser.Service.GetRoomUuid(userValue.SecWebSocketKey)

				if v.RoomUuid == roomUuid && userValue.JackpotGroupId == groupId {
					flux.Send(msgId, s.Id, userKey, dataBuffer)
				}
			}

		}
	}
}
