package redenvelope

import (
	"fmt"

	"serve/fish_comm/common"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/fish_comm/rng"
)

const (
	ActionRedEnvelopeBonus        = "ActionRedEnvelopeBonus"
	ActionRedEnvelopeExistCheck   = "ActionRedEnvelopeExistCheck"
	ActionRedEnvelopeNotExist     = "ActionRedEnvelopeNotExist"
	ActionRedEnvelopePick         = "ActionRedEnvelopePick"
	ActionRedEnvelopeRefundCall   = "ActionRedEnvelopeRefundCall"
	ActionRedEnvelopeRefundRecall = "ActionRedEnvelopeRefundRecall"
	ActionRedEnvelopeRefundDone   = "ActionRedEnvelopeRefundDone"
)

var Service = &service{
	Id: "RedEnvelopeService",
	in: make(chan *flux.Action, common.Service.ChanSize),
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
	hitRedEnvelopes := make(map[string]*RedEnvelope)

	for {
		select {
		case action, ok := <-s.in:
			if !ok {
				return
			}

			switch action.Key().Name() {
			case ActionRedEnvelopeBonus:
				re := action.Payload()[0].(*RedEnvelope)
				isOk := action.Payload()[1].(chan bool)

				hitRedEnvelopes[re.Uuid] = re

				logger.Service.Zap.Infow("Got Red Envelope",
					"GameUser", re.SecWebSocketKey,
					"GameRoomUuid", re.RoomUuid,
					"FishType", re.FishTypeId,
					"BonusUuid", re.Uuid,
					"BonusPayload", re.BonusPayload,
					"Bet", re.Bet,
					"Rate", re.Rate,
				)

				isOk <- true
				close(isOk)

			case ActionRedEnvelopeExistCheck:
				re := action.Payload()[0].(*RedEnvelope)
				depositMultiple := action.Payload()[1].(uint64)

				if v, ok := hitRedEnvelopes[re.Uuid]; ok {
					if v.SecWebSocketKey == re.SecWebSocketKey {
						v.PlayerOptIndex = re.PlayerOptIndex

						s.rngRedEnvelope(v)

						delete(hitRedEnvelopes, v.Uuid)

						flux.Send(ActionRedEnvelopePick, s.Id, action.Key().From(), v.SecWebSocketKey, v, depositMultiple)

						logger.Service.Zap.Infow("Check Red Envelope",
							"GameUser", v.SecWebSocketKey,
							"GameRoomUuid", v.RoomUuid,
							"FishType", v.FishTypeId,
							"PlayerOptIndex", v.PlayerOptIndex,
							"BonusUuid", v.Uuid,
							"Pay", v.Pay,
							"BonusPayload", v.BonusPayload,
							"AllPay", v.AllPay,
							"Bet", v.Bet,
							"Rate", v.Rate,
						)
					} else {
						flux.Send(ActionRedEnvelopeNotExist, s.Id, action.Key().From(), v.SecWebSocketKey, re, depositMultiple)
					}
				} else {
					flux.Send(ActionRedEnvelopeNotExist, s.Id, action.Key().From(), re.SecWebSocketKey, re, depositMultiple)
				}

			case ActionRedEnvelopeRefundCall:
				secWebSocketKey := action.Payload()[0].(string)
				accountingSn := action.Payload()[1].(uint64)

				for k, v := range hitRedEnvelopes {
					if v.SecWebSocketKey == secWebSocketKey {

						v.PlayerOptIndex = 0

						s.rngRedEnvelope(v)

						delete(hitRedEnvelopes, k)

						logger.Service.Zap.Infow("Refund Red Envelope",
							"GameUser", v.SecWebSocketKey,
							"GameRoomUuid", v.RoomUuid,
							"FishType", v.FishTypeId,
							"PlayerOptIndex", v.PlayerOptIndex,
							"BonusUuid", v.Uuid,
							"Pay", v.Pay,
							"BonusPayload", v.BonusPayload,
							"AllPay", v.AllPay,
							"Bet", v.Bet,
							"Rate", v.Rate,
						)
						flux.Send(ActionRedEnvelopeRefundRecall, s.Id, action.Key().From(), secWebSocketKey, v, accountingSn)
					}
				}
				flux.Send(ActionRedEnvelopeRefundDone, s.Id, action.Key().From(), secWebSocketKey)
			}
		}
	}
}

func (s *service) rngRedEnvelope(r *RedEnvelope) *RedEnvelope {
	indexOptions0 := s.remove(nil, -1)
	rngIndex0 := s.rngIndex(indexOptions0)
	s.pickIndex(r, rngIndex0, 0)
	r.AllPay[rngIndex0] = int64(r.BonusPayload[0])

	indexOptions1 := s.remove(indexOptions0, rngIndex0)
	rngIndex1 := s.rngIndex(indexOptions1)
	s.pickIndex(r, rngIndex1, 1)
	r.AllPay[rngIndex1] = int64(r.BonusPayload[1])

	indexOptions2 := s.remove(indexOptions1, rngIndex1)
	rngIndex2 := s.rngIndex(indexOptions2)
	s.pickIndex(r, rngIndex2, 2)
	r.AllPay[rngIndex2] = int64(r.BonusPayload[2])

	indexOptions3 := s.remove(indexOptions2, rngIndex2)
	rngIndex3 := s.rngIndex(indexOptions3)
	s.pickIndex(r, rngIndex3, 3)
	r.AllPay[rngIndex3] = int64(r.BonusPayload[3])

	indexOptions4 := s.remove(indexOptions3, rngIndex3)
	rngIndex4 := s.rngIndex(indexOptions4)
	s.pickIndex(r, rngIndex4, 4)
	r.AllPay[rngIndex4] = int64(r.BonusPayload[4])

	return r
}

func (s *service) pickIndex(r *RedEnvelope, rngIndex, bonusPayloadIndex int) {
	if int(r.PlayerOptIndex) == rngIndex {
		r.Pay = uint64(r.BonusPayload[bonusPayloadIndex])

		logger.Service.Zap.Infow("Pick Red Envelope",
			"GameUser", r.SecWebSocketKey,
			"GameRoomUuid", r.RoomUuid,
			"FishType", r.FishTypeId,
			"BonusUuid", r.Uuid,
			"BonusPayload", r.BonusPayload,
			"Bet", r.Bet,
			"Rate", r.Rate,
			"Pay", r.Pay,
			"PlayerOptIndex", r.PlayerOptIndex,
		)
	}
}

func (s *service) rngIndex(options []rng.Option) (index int) {
	return rng.New(options).Item.(int)
}

func (s *service) remove(indexOptions []rng.Option, indexOptionItem int) []rng.Option {
	if indexOptionItem == -1 {
		var indexOptions []rng.Option
		indexOptions = append(indexOptions, rng.Option{1, 0})
		indexOptions = append(indexOptions, rng.Option{1, 1})
		indexOptions = append(indexOptions, rng.Option{1, 2})
		indexOptions = append(indexOptions, rng.Option{1, 3})
		indexOptions = append(indexOptions, rng.Option{1, 4})
		return indexOptions
	}

	var newIndexOptions []rng.Option

	for _, v := range indexOptions {
		if v.Item.(int) != indexOptionItem {
			newIndexOptions = append(newIndexOptions, v)
		}
	}

	return newIndexOptions
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}
