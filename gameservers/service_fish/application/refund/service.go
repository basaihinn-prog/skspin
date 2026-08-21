package refund

import (
	"fmt"
	"os"
	"serve/service_fish/application/activity"
	"serve/service_fish/application/gamesetting"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/redenvelope"
	"serve/service_fish/domain/slot"
	"serve/service_fish/models"
	"strconv"
	"sync"
	"time"

	"serve/fish_comm/common"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/fish_comm/wallet"
)

const (
	ActionRefundInit   = "ActionRefundInit"
	ActionRefundCall   = "ActionRefundCall"
	ActionRefundRecall = "ActionRefundRecall"
	ActionRefundDelete = "ActionRefundDelete"
)

var Service = &service{
	Id:       "RefundService",
	poolSize: os.Getenv("POOL_SIZE"),
	in:       make(chan *flux.Action, common.Service.ChanSize),
	rwMutex:  &sync.RWMutex{},
}

type service struct {
	Id       string
	poolSize string
	in       chan *flux.Action
	rwMutex  *sync.RWMutex
}

func init() {
	Service.run()
	logger.Service.Zap.Infow("Service created.",
		"Service", Service.Id,
		"Chan", fmt.Sprintf("%p", Service.in),
	)
}

func (s *service) run() {
	flux.Register(s.Id, s.in)
	allRefund := make(map[string]*refund)

	poolSize := 10

	if s.poolSize != "" {
		if poolSize, _ = strconv.Atoi(s.poolSize); poolSize == 0 {
			poolSize = 10
		}
	}

	for i := 0; i < poolSize; i++ {
		go func() {
			for action := range s.in {
				s.handleAction(action, allRefund)
			}
		}()
	}
}

func (s *service) handleAction(action *flux.Action, allRefund map[string]*refund) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	switch action.Key().Name() {
	case ActionRefundInit:
		controllerId := action.Key().From()
		secWebSocketKey := action.Payload()[0].(string)

		if v, ok := allRefund[secWebSocketKey]; ok {
			v.controllerId = controllerId
			v.reset()
		} else {
			switch gamesetting.Service.GameId(secWebSocketKey) {
			case models.PSF_ON_00001, models.PGF_ON_00001:
				psf_on_00001_init(action, allRefund)

			case models.PSF_ON_00002, models.PSF_ON_20002:
				psf_on_00002_init(action, allRefund)

			case models.PSF_ON_00003:
				psf_on_00003_init(action, allRefund)

			case models.PSF_ON_00004:
				psf_on_00004_init(action, allRefund)

			case models.PSF_ON_00005:
				psf_on_00005_init(action, allRefund)

			case models.PSF_ON_00006:
				psf_on_00006_init(action, allRefund)

			case models.PSF_ON_00007:
				psf_on_00007_init(action, allRefund)

			case models.RKF_H5_00001:
				rkf_h5_00001_init(action, allRefund)

			default:
				logger.Service.Zap.Errorw(Refund_GAME_ID_INVALID,
					"GameUser", secWebSocketKey,
				)
				errorcode.Service.Fatal(secWebSocketKey, Refund_GAME_ID_INVALID)
			}
		}

	case ActionRefundCall:
		secWebSocketKey := action.Payload()[0].(string)
		isDisconnect := action.Payload()[1].(bool)
		refundDone := action.Payload()[2].(chan bool)
		accountingSn := action.Payload()[3].(uint64)

		if v, ok := allRefund[secWebSocketKey]; ok {
			v.isDisconnect = isDisconnect

			flux.Send(bullet.ActionBulletRefundCall, s.Id, bullet.Service.HashId(secWebSocketKey), secWebSocketKey, accountingSn)
			flux.Send(redenvelope.ActionRedEnvelopeRefundCall, s.Id, redenvelope.Service.Id, secWebSocketKey, accountingSn)
			flux.Send(slot.ActionSlotRefundCall, s.Id, slot.Service.Id, secWebSocketKey, accountingSn)

			refundDone <- true

		} else {
			logger.Service.Zap.Errorw(Refund_NOT_FOUND,
				"GameUser", secWebSocketKey,
				"IsDisconnect", isDisconnect,
			)
			errorcode.Service.Fatal(secWebSocketKey, Refund_NOT_FOUND)
		}

	case bullet.ActionBulletRefundRecall:
		secWebSocketKey := action.Payload()[0].(string)
		b := action.Payload()[1].(*bullet.Bullet)
		accountingSn := action.Payload()[2].(uint64)

		// bonus bullets
		if b.Bullets > 0 && b.UsedBullets >= 0 {
			restBullets := uint64(b.Bullets - b.UsedBullets)

			switch wallet.Service.WalletCategory(secWebSocketKey) {
			case wallet.CategoryHost:
				wallet.Service.IncreaseCent(secWebSocketKey, b.BonusUuid, b.Bet*b.Rate*restBullets)
			default:
				wallet.Service.IncreaseBalance(secWebSocketKey, b.BonusUuid, b.Bet*b.Rate*restBullets, 0, 0)
			}

			if v, ok := allRefund[secWebSocketKey]; ok && v.iHandler != nil {
				time.Sleep(500 * time.Millisecond)
				v.iHandler.RefundHandler(b, accountingSn)
			}
		}

	case bullet.ActionBulletRefundDone:
		secWebSocketKey := action.Payload()[0].(string)
		s.refund(secWebSocketKey, bullet.Service.Id, allRefund)

	case redenvelope.ActionRedEnvelopeRefundRecall:
		secWebSocketKey := action.Payload()[0].(string)
		re := action.Payload()[1].(*redenvelope.RedEnvelope)
		accountingSn := action.Payload()[2].(uint64)

		switch wallet.Service.WalletCategory(secWebSocketKey) {
		case wallet.CategoryHost:
			wallet.Service.IncreaseCent(secWebSocketKey, re.Uuid, re.Bet*re.Pay*re.Rate)
		default:
			wallet.Service.IncreaseBalance(secWebSocketKey, re.Uuid, re.Bet*re.Pay*re.Rate, 0, 0)
		}

		activity.Service.RecordBonus(secWebSocketKey, re)

		if v, ok := allRefund[secWebSocketKey]; ok && v.iHandler != nil {
			time.Sleep(500 * time.Millisecond)
			v.iHandler.RefundHandler(re, accountingSn)
		}

	case redenvelope.ActionRedEnvelopeRefundDone:
		secWebSocketKey := action.Payload()[0].(string)
		s.refund(secWebSocketKey, redenvelope.Service.Id, allRefund)

	case slot.ActionSlotRefundRecall:
		secWebSocketKey := action.Payload()[0].(string)
		sl := action.Payload()[1].(*slot.Slot)
		accountingSn := action.Payload()[2].(uint64)

		switch wallet.Service.WalletCategory(secWebSocketKey) {
		case wallet.CategoryHost:
			wallet.Service.IncreaseCent(secWebSocketKey, sl.Uuid, sl.Bet*sl.Pay*sl.Rate)
		default:
			wallet.Service.IncreaseBalance(secWebSocketKey, sl.Uuid, sl.Bet*sl.Pay*sl.Rate, 0, 0)
		}

		activity.Service.RecordBonus(secWebSocketKey, sl)

		if v, ok := allRefund[secWebSocketKey]; ok && v.iHandler != nil {
			time.Sleep(500 * time.Millisecond)
			v.iHandler.RefundHandler(sl, accountingSn)
		}

	case slot.ActionSlotRefundDone:
		secWebSocketKey := action.Payload()[0].(string)
		s.refund(secWebSocketKey, slot.Service.Id, allRefund)

	case bullet.ActionBulletDestroyDone:
		secWebSocketKey := action.Payload()[0].(string)

		s.bulletDestroyRefund(secWebSocketKey, allRefund)

	case ActionRefundDelete:
		secWebSocketKey := action.Payload()[0].(string)

		delete(allRefund, secWebSocketKey)
	}
}

func (s *service) bulletDestroyRefund(secWebSocketKey string, allRefund map[string]*refund) {
	for k, v := range allRefund {
		if v.secWebSocketKey == secWebSocketKey {
			s.refund(k, bullet.Service.Id, allRefund)
		}
	}
}

func (s *service) refund(secWebSocketKey, serviceId string, allRefund map[string]*refund) {
	if r, ok := allRefund[secWebSocketKey]; ok {
		if !r.refundStatus[serviceId] {
			r.refundStatus[serviceId] = true

			logger.Service.Zap.Infow("Refund Bonus Status",
				"GameUser", secWebSocketKey,
				"RefundType", serviceId,
				"RefundStatus", r.refundStatus,
			)

			isAllDone := true

			for _, v := range r.refundStatus {
				if !v {
					isAllDone = false
					break
				}
			}

			if isAllDone {
				flux.Send(ActionRefundRecall, s.Id, r.controllerId, secWebSocketKey, r.isDisconnect)
			}
		}
	}
}

func (s *service) HashId(secWebSocketKey string, accountingSn uint64) string {
	return secWebSocketKey + s.Id + strconv.FormatUint(accountingSn, 10)
}

func (s *service) Destroy() {
	flux.UnRegister(s.Id, s.in)
}
