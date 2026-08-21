package bonuslottery

import (
	"fmt"
	"serve/service_fish/application/gamesetting"
	"serve/service_fish/domain/auth"
	"serve/service_fish/domain/probability"
	probability_proto "serve/service_fish/domain/probability/proto"
	"serve/service_fish/domain/redenvelope"
	"serve/service_fish/domain/slot"
	"serve/service_fish/models"

	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
)

var Service = &service{
	Id: "BonusLotteryService",
	in: make(chan *flux.Action, 1024),
}

type service struct {
	Id string
	in chan *flux.Action
}

func init() {
	go Service.run()
	logger.Service.Zap.Infow("Service created.",
		"Service", Service.Id,
		"Chan", fmt.Sprintf("%p", Service.in),
	)
}

func (s *service) run() {
	flux.Register(s.Id, s.in)

	for {
		select {
		case action, ok := <-s.in:
			if !ok {
				return
			}
			s.handleAction(action)
		}
	}
}

func (s *service) handleAction(action *flux.Action) {
	switch action.Key().Name() {
	case flux.ActionFluxRegisterDone:
		// do nothing

	case probability.EMSGID_eOptionCall:
		secWebSocketKey := action.Key().From()
		data := action.Payload()[0].(*probability_proto.OptionCall)
		depositMultiple := action.Payload()[1].(uint64)

		isTokenValid := make(chan bool, 1024)
		flux.Send(auth.ActionCheckToken, s.Id, auth.Service.Id, isTokenValid, data.Token)

		if <-isTokenValid {
			switch data.OptMode {
			case probability_proto.OPTION_MODE_eGetSelectRedEnvelopeOption:
				re := redenvelope.New(data.BonusUuid, secWebSocketKey, "", data.RoomUuid, data.SeatId, 0, nil, 0, 0, 0, 0)
				re.PlayerOptIndex = data.PlayerOptIndex
				flux.Send(redenvelope.ActionRedEnvelopeExistCheck, s.Id, redenvelope.Service.Id, re, depositMultiple)

			case probability_proto.OPTION_MODE_eGetSelectSlot777Option:
				st := slot.New(data.BonusUuid, secWebSocketKey, "", data.RoomUuid, data.SeatId, -1, 0, 0, 0, 0)
				flux.Send(slot.ActionSlotExistCheck, s.Id, slot.Service.Id, st, depositMultiple)

			default:
				logger.Service.Zap.Errorw(BonusLottery_OPTION_MODE_PROTO_INVALID,
					"GameUser", secWebSocketKey,
					"Error", data.OptMode,
				)
				errorcode.Service.Fatal(secWebSocketKey, BonusLottery_OPTION_MODE_PROTO_INVALID)
			}
		} else {
			logger.Service.Zap.Errorw(auth.Auth_TOKEN_EXPIRED,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, auth.Auth_TOKEN_EXPIRED)
		}

	default:
		secWebSocketKey := action.Payload()[0].(string)
		depositMultiple := action.Payload()[2].(uint64)

		switch gamesetting.Service.GameId(secWebSocketKey) {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_RedEnvelope(action, depositMultiple)
			psf_on_00001_Slot(action, depositMultiple)

		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_RedEnvelope(action, depositMultiple)
			psf_on_00002_Slot(action, depositMultiple)

		case models.PSF_ON_00003:
			psf_on_00003_Slot(action, depositMultiple)

		case models.PSF_ON_00004:
			psf_on_00004_RedEnvelope(action, depositMultiple)
			psf_on_00004_Slot(action, depositMultiple)

		case models.PSF_ON_00005:
			psf_on_00005_RedEnvelope(action, depositMultiple)
			psf_on_00005_Slot(action, depositMultiple)

		case models.PSF_ON_00006:
			psf_on_00006_RedEnvelope(action, depositMultiple)
			psf_on_00006_Slot(action, depositMultiple)

		case models.PSF_ON_00007:
			psf_on_00007_RedEnvelope(action, depositMultiple)
			psf_on_00007_Slot(action, depositMultiple)

		case models.RKF_H5_00001:
			rkf_h5_00001_RedEnvelope(action, depositMultiple)
			rkf_h5_00001_Slot(action, depositMultiple)

		default:
			logger.Service.Zap.Errorw(BonusLottery_GAME_ID_INVALID,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, BonusLottery_GAME_ID_INVALID)
		}
	}
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}
