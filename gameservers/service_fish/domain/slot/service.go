package slot

import (
	"fmt"

	"serve/fish_comm/common"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
)

const (
	ActionSlotBonus        = "ActionSlotBonus"
	ActionSlotExistCheck   = "ActionSlotExistCheck"
	ActionSlotNotExist     = "ActionSlotNotExist"
	ActionSlotPick         = "ActionSlotPick"
	ActionSlotRefundCall   = "ActionSlotRefundCall"
	ActionSlotRefundRecall = "ActionSlotRefundRecall"
	ActionSlotRefundDone   = "ActionSlotRefundDone"
)

var Service = &service{
	Id: "SlotService",
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
	hitSlots := make(map[string]*Slot)

	for {
		select {
		case action, ok := <-s.in:
			if !ok {
				return
			}

			switch action.Key().Name() {
			case ActionSlotBonus:
				slot := action.Payload()[0].(*Slot)
				isOk := action.Payload()[1].(chan bool)

				hitSlots[slot.Uuid] = slot

				logger.Service.Zap.Infow("Got Slot",
					"GameUser", slot.SecWebSocketKey,
					"GameRoomUuid", slot.RoomUuid,
					"FishType", slot.FishTypeId,
					"BonusUuid", slot.Uuid,
					"AllPay", slot.AllPay,
					"Bet", slot.Bet,
					"Rate", slot.Rate,
					"Pay", slot.Pay,
				)

				isOk <- true
				close(isOk)

			case ActionSlotExistCheck:
				slot := action.Payload()[0].(*Slot)
				depositMultiple := action.Payload()[1].(uint64)

				if v, ok := hitSlots[slot.Uuid]; ok {
					if v.SecWebSocketKey == slot.SecWebSocketKey {

						delete(hitSlots, v.Uuid)

						flux.Send(ActionSlotPick, s.Id, action.Key().From(), v.SecWebSocketKey, v, depositMultiple)
					} else {
						flux.Send(ActionSlotNotExist, s.Id, action.Key().From(), v.SecWebSocketKey, slot, depositMultiple)
					}
				} else {
					flux.Send(ActionSlotNotExist, s.Id, action.Key().From(), slot.SecWebSocketKey, slot, depositMultiple)
				}

			case ActionSlotRefundCall:
				secWebSocketKey := action.Payload()[0].(string)
				accountingSn := action.Payload()[1].(uint64)

				for k, v := range hitSlots {
					if v.SecWebSocketKey == secWebSocketKey {

						delete(hitSlots, k)

						logger.Service.Zap.Infow("Refund Slot",
							"GameUser", secWebSocketKey,
							"GameRoomUuid", v.RoomUuid,
							"FishType", v.FishTypeId,
							"Pay", v.Pay,
							"AllPay", v.AllPay,
							"Bet", v.Bet,
							"Rate", v.Rate,
						)

						flux.Send(ActionSlotRefundRecall, s.Id, action.Key().From(), secWebSocketKey, v, accountingSn)
					}
				}
				flux.Send(ActionSlotRefundDone, s.Id, action.Key().From(), secWebSocketKey)
			}
		}
	}
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}
