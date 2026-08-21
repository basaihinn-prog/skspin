package bullet

import (
	"fmt"

	"serve/fish_comm/common"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
)

const (
	ActionBulletStart  = "ActionBulletStart"
	ActionBulletStop   = "ActionBulletStop"
	ActionBulletDelete = "ActionBulletDelete"
	PLAYER             = "PLAYER"
	BOT                = "BOT"
	GUEST              = "GUEST"
)

var Service = &service{
	Id: "BulletService",
	in: make(chan *flux.Action, common.Service.ChanSize),
}

type service struct {
	Id     string
	gameId string
	in     chan *flux.Action
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
	bulletPools := make(map[string]*bulletPool)

	for {
		logger.Service.Zap.Debugw("AllBulletPools",
			"Size", len(bulletPools),
			"Map", bulletPools,
		)

		select {
		case action, ok := <-s.in:
			if !ok {
				for _, v := range bulletPools {
					v.destroy()
					delete(bulletPools, v.hashId())
				}
				return
			}

			switch action.Key().Name() {
			case ActionBulletStart:
				secWebSocketKey := action.Key().From()
				gameId := action.Payload()[0].(string)
				mathModuleId := action.Payload()[1].(string)
				kind := action.Payload()[2].(string)

				// See bulletPool's hashId()
				id := s.HashId(secWebSocketKey)

				if _, ok := bulletPools[id]; !ok {
					bulletPools[id] = newBulletPool(secWebSocketKey, gameId, mathModuleId, kind)
				}

			case ActionBulletStop:
				//secWebSocketKey := action.Key().From()
				//refundServiceId := action.Payload()[0]
				//gameRoomUuid := action.Payload()[1].(string)

				// See bulletPool's hashId()
				//hashId := s.HashId(secWebSocketKey)

				//if v, ok := bulletPools[hashId]; ok {
				//flux.Send(actionBulletDestroyCall, s.Id, v.hashId(), refundServiceId)
				//v.destroy()
				//delete(bulletPools, v.hashId())
				//}

			case ActionBulletDelete:
				secWebSocketKey := action.Key().From()
				refundServiceId := action.Payload()[0]

				for _, v := range bulletPools {
					if v.secWebSocketKey == secWebSocketKey {
						flux.Send(actionBulletDestroyCall, s.Id, v.hashId(), refundServiceId)
					}
				}

			case actionBulletDestroyRecall:
				refundServiceId := action.Payload()[0].(string)
				secWebSocketKey := action.Payload()[1].(string)

				flux.Send(ActionBulletDestroyDone, s.Id, refundServiceId, secWebSocketKey)

				// See bulletPool's hashId()
				id := s.HashId(secWebSocketKey)

				if v, ok := bulletPools[id]; ok {
					v.destroy()
					delete(bulletPools, v.hashId())
				}
			}
		}
	}
}

func (s *service) HashId(secWebSocketKey string) string {
	return secWebSocketKey + s.Id
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}
