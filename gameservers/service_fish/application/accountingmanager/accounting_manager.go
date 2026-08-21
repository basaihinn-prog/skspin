package accountingmanager

import (
	"log/slog"
	"serve/service_fish/application/gamesetting"
	"serve/service_fish/application/refund"
	"serve/service_fish/domain/bet"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	"serve/service_fish/domain/redenvelope"
	"serve/service_fish/domain/slot"
	"serve/service_fish/models"
	"strconv"
	"strings"
	"sync"
	"time"

	"serve/fish_comm/accounting"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/jackpot-client/jackpot"
	"serve/fish_comm/logger"
	"serve/fish_comm/vip"
	"serve/fish_comm/wallet"
)

const (
	actionAccountingBonus             = "actionAccountingBonus"
	actionAccountingHitBullet         = "actionAccountingHitBullet"
	actionAccountingHitBulletMulti    = "actionAccountingHitBulletMulti"
	actionAccountingRefundBonusBullet = "actionAccountingRefundBonusBullet"
	actionRefund                      = "actionRefund"
	actionRotation                    = "actionRotation"
	actionDecrease                    = "actionDecrease"
	actionIncrease                    = "actionIncrease"
	actionBulletHit                   = "actionBulletHit"
	actionBulletBonus                 = "actionBulletBonus"
	actionBulletMultiHit              = "actionBulletMultiHit"
	actionPickBonus                   = "actionPickBonus"
	actionDisconnect                  = "actionDisconnect"
	actionLeaveGameRoom               = "actionLeaveGameRoom"
	actionGetAccountingSn             = "actionGetAccountingSn"
	actionEndAccounting               = "actionEndAccounting"
	actionJackpotHit                  = "actionJackpotHit"
	actionWalletHostAmt               = "actionWalletHostAmt"
	actionBonusBulletShooted          = "actionBonusBulletShooted"
	ActionUpdateBlacklist             = "ActionUpdateBlacklist"
)

func Builder() *builder {
	return &builder{
		secWebSocketKey: "",
		memberId:        "",
		hostId:          "",
		hostExtId:       "",
		guestEnable:     false,
		billCurrency:    "",
		depositMultiple: 0,
	}
}

type builder struct {
	secWebSocketKey string
	memberId        string
	hostId          string
	hostExtId       string
	guestEnable     bool
	category        int
	billCurrency    string
	depositMultiple uint64
}

type AccountingManager struct {
	secWebSocketKey string
	memberId        string
	hostId          string
	hostExtId       string
	guestEnable     bool
	accountingSn    uint64
	isLeaveGameRoom bool // if true then need refund
	isDisconnect    bool // if true then need refund
	timerPause      bool
	timerStop       bool
	rounds          uint64
	increaseSn      map[string]uint64
	decreaseSn      map[string]uint64
	walletType      int
	depositCent     uint64
	billCurrency    string
	depositMultiple uint64
	in              chan *flux.Action
	rwMutex         *sync.RWMutex
	RotatiomRwMutex *sync.RWMutex
}

type tempReceipt struct {
	accountingSn uint64
	initCent     uint64
	endCent      uint64
	betCent      uint64
	winCent      uint64
	bonusWinCent uint64
	jackpotWin   uint64
	isDone       bool
}

func (b *builder) SetSecWebSocketKey(secWebSocketKey string) *builder {
	b.secWebSocketKey = secWebSocketKey
	return b
}

func (b *builder) SetMemberId(memberId string) *builder {
	b.memberId = memberId
	return b
}

func (b *builder) SetHostId(hostId string) *builder {
	b.hostId = hostId
	return b
}

func (b *builder) SetHostExtId(hostExtId string) *builder {
	b.hostExtId = hostExtId
	return b
}

func (b *builder) SetGuestEnable(guestEnable bool) *builder {
	b.guestEnable = guestEnable
	return b
}

func (b *builder) SetCategory(category int) *builder {
	b.category = category
	return b
}

func (b *builder) SetBillCurrency(billCurrency string) *builder {
	b.billCurrency = billCurrency
	return b
}

func (b *builder) SetDepositMultiple(depositMultiple uint64) *builder {
	b.depositMultiple = depositMultiple
	return b
}

func (b *builder) Build() *AccountingManager {
	if b.guestEnable {
		return nil
	}

	v := &AccountingManager{
		secWebSocketKey: b.secWebSocketKey,
		memberId:        b.memberId,
		hostId:          b.hostId,
		hostExtId:       b.hostExtId,
		guestEnable:     b.guestEnable,
		accountingSn:    0,
		isLeaveGameRoom: false,
		isDisconnect:    false,
		timerPause:      false,
		timerStop:       false,
		rounds:          0,
		increaseSn:      make(map[string]uint64),
		decreaseSn:      make(map[string]uint64),
		walletType:      b.category,
		billCurrency:    b.billCurrency,
		depositMultiple: b.depositMultiple,
		in:              make(chan *flux.Action, 1024),
		rwMutex:         &sync.RWMutex{},
		RotatiomRwMutex: &sync.RWMutex{},
	}

	go v.run()

	return v
}

func (am *AccountingManager) run() {
	flux.Register(am.hashId(am.secWebSocketKey), am.in)

	bulletUuid := make(map[string]uint64)
	redEnvelopeUuid := make(map[string]uint64)
	slotUuid := make(map[string]uint64)
	drillUuid := make(map[string]map[uint64]int32)
	machineGunUuid := make(map[string]map[uint64]int32)
	tempReceipts := make(map[uint64]*tempReceipt)
	superMachineGunUuid := make(map[string]map[uint64]int32)
	bonusBulletShooted := make(map[string]map[string]map[uint64]int32)

	for action := range am.in {
		am.handleAction(action, bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, bonusBulletShooted, tempReceipts)
	}
}

func (am *AccountingManager) handleAction(
	action *flux.Action,
	bulletUuid, redEnvelopeUuid, slotUuid map[string]uint64,
	drillUuid, machineGunUuid, superMachineGunUuid map[string]map[uint64]int32,
	bonusBulletShooted map[string]map[string]map[uint64]int32,
	tempReceipts map[uint64]*tempReceipt) {

	logger.Service.Zap.Infow(action.Key().Name(),
		"GameUser", am.secWebSocketKey,
	)

	switch action.Key().Name() {
	case flux.ActionFluxRegisterDone:
		// do nothing

	case actionRefund:
		bonus := action.Payload()[0]
		accountingSn := action.Payload()[1].(uint64)

		am.checkStart(0, bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, bonusBulletShooted, tempReceipts)
		if accountingSn == 0 {
			accountingSn = am.accountingSn
		}

		logger.Service.Zap.Infow("Accounting Manager Refund",
			"GameUser", am.secWebSocketKey,
			"am.AccountingSn", am.accountingSn,
			"AccountingSn", accountingSn,
			"Bullets", len(bulletUuid),
			"RedEnvelopes", len(redEnvelopeUuid),
			"Slots", len(slotUuid),
			"Drills", len(drillUuid),
			"MachineGuns", len(machineGunUuid),
			"SuperMachineGuns", len(superMachineGunUuid),
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
			"InitCent", tempReceipts[accountingSn].initCent,
			"BetCent", tempReceipts[accountingSn].betCent,
			"WinCent", tempReceipts[accountingSn].winCent,
			"BonusWinCent", tempReceipts[accountingSn].bonusWinCent,
		)

		switch bonus.(type) {
		case *bullet.Bullet:
			b := bonus.(*bullet.Bullet)

			restBullets := b.Bullets - b.UsedBullets

			tempReceipts[accountingSn].bonusWinCent += b.Bet * b.Rate * uint64(restBullets)

			isDone := make(chan bool)

			flux.Send(actionAccountingRefundBonusBullet, am.hashId(am.secWebSocketKey), am.to(am.accountingSn),
				am.secWebSocketKey, b, isDone, accountingSn)

			<-isDone

			gameId := gamesetting.Service.GameId(b.SecWebSocketKey)
			switch gameId {
			case models.PSF_ON_00001, models.PGF_ON_00001:
				psf_on_00001_CheckBulletTypeId(b.TypeId, restBullets, am, b, drillUuid, machineGunUuid)
			case models.PSF_ON_00002, models.PSF_ON_20002:
				psf_on_00002_CheckBulletTypeId(b.TypeId, restBullets, am, b, drillUuid, machineGunUuid)
			case models.PSF_ON_00003:
				psf_on_00003_CheckBulletTypeId(b.TypeId, restBullets, am, b, drillUuid)
			case models.PSF_ON_00004:
				psf_on_00004_CheckBulletTypeId(b.TypeId, restBullets, am, b, machineGunUuid, superMachineGunUuid)
			case models.PSF_ON_00005:
				psf_on_00005_CheckBulletTypeId(b.TypeId, restBullets, am, b, machineGunUuid, superMachineGunUuid)
			case models.PSF_ON_00006:
				psf_on_00006_CheckBulletTypeId(b.TypeId, restBullets, am, b, machineGunUuid, superMachineGunUuid)
			case models.PSF_ON_00007:
				psf_on_00007_CheckBulletTypeId(b.TypeId, restBullets, am, b, machineGunUuid, superMachineGunUuid)
			case models.RKF_H5_00001:
				rkf_h5_00001_CheckBulletTypeId(b.TypeId, restBullets, am, b, machineGunUuid, superMachineGunUuid)
			}

		case *redenvelope.RedEnvelope:
			// 2. 核銷紅包 858cedf5
			re := bonus.(*redenvelope.RedEnvelope)

			tempReceipts[accountingSn].bonusWinCent += re.Bet * re.Pay * re.Rate

			isDone := make(chan bool)

			flux.Send(actionAccountingBonus, Service.id, am.to(am.accountingSn),
				am.secWebSocketKey, re, isDone, accountingSn)

			<-isDone

			delete(redEnvelopeUuid, re.Uuid)

		case *slot.Slot:
			// 2. 核銷老虎機 6e01e01a
			sl := bonus.(*slot.Slot)

			tempReceipts[accountingSn].bonusWinCent += sl.Bet * sl.Pay * sl.Rate

			// Process Multi Slot Write Off
			gameId := gamesetting.Service.GameId(sl.SecWebSocketKey)
			if gameId != models.PSF_ON_00003 {
				isDone := make(chan bool)
				flux.Send(actionAccountingBonus, Service.id, am.to(am.accountingSn),
					am.secWebSocketKey, sl, isDone, accountingSn)
				<-isDone
			}

			delete(slotUuid, sl.Uuid)
		}

		if am.checkEnd(bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, tempReceipts, false) {
			if am.walletType == wallet.CategoryHost {
				accountingSn := am.accountingSn
				betAmt, winAmt := am.getWalletHostAmt(accountingSn, tempReceipts)
				slog.Info("should be do ", betAmt, winAmt)
				flux.Send(accounting.ActionAccountingStop, am.hashId(am.secWebSocketKey), am.to(accountingSn), "")

				//cent, _ := wallet.Service.CentV2(am.secWebSocketKey)
				//if cent == 0 {
				//	wallet.Service.IncreaseBalanceV2(am.secWebSocketKey, "", 0, betAmt, winAmt)
				//} else {
				//	wallet.Service.WithdrawV2(am.secWebSocketKey, accountingSn, cent, betAmt, winAmt)
				//}
			}
			am.isLeaveGameRoom = false
			am.accountingSn = 0
			am.timerPause = false
		}

	case actionRotation:
		accountingSn := action.Payload()[0].(uint64)
		isDone := action.Payload()[1].(chan bool)

		defer am.RotatiomRwMutex.Unlock()

		if len(bulletUuid) > 0 {
			isDone <- false
			return
		}

		tempReceipts[accountingSn].isDone = true

		totalBet := tempReceipts[accountingSn].betCent
		totalBonusWin := tempReceipts[accountingSn].bonusWinCent
		totalWin := tempReceipts[accountingSn].winCent
		jackpotWin := tempReceipts[accountingSn].jackpotWin

		logger.Service.Zap.Infow("Accounting Manager Rotation",
			"GameUser", am.secWebSocketKey,
			"AccountingSn1", am.accountingSn,
			"AccountingSn2", accountingSn,
			"Bullets", len(bulletUuid),
			"RedEnvelopes", len(redEnvelopeUuid),
			"Slots", len(slotUuid),
			"Drills", len(drillUuid),
			"MachineGuns", len(machineGunUuid),
			"SuperMachineGuns", len(superMachineGunUuid),
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
			"InitCent", tempReceipts[accountingSn].initCent,
			"BetCent", tempReceipts[accountingSn].betCent,
			"WinCent", tempReceipts[accountingSn].winCent,
			"BonusWinCent", tempReceipts[accountingSn].bonusWinCent,
			"JackpotWin", tempReceipts[accountingSn].jackpotWin,
		)

		if am.walletType == wallet.CategoryHost {
			//winAmt := totalWin + jackpotWin
			//cent, _ := wallet.Service.CentCacheV2(am.secWebSocketKey)
			//if cent > 0 {
			//	wallet.Service.WithdrawV2(am.secWebSocketKey, 0, cent, totalBet, winAmt)
			//} else {
			//	wallet.Service.IncreaseBalanceV2(am.secWebSocketKey, "", 0, totalBet, winAmt)
			//}
		}

		flux.Send(accounting.ActionAccountingEnd, am.hashId(am.secWebSocketKey), am.to(accountingSn),
			totalBet, totalBonusWin, totalWin, jackpotWin, accountingSn)

		if am.accountingSn == accountingSn {
			am.accountingSn = 0
			am.timerPause = false
		}

		isDone <- true

	case actionGetAccountingSn:
		balance := action.Payload()[0].(uint64)
		sn := action.Payload()[1].(uint64)

		am.accountingSn = sn

		if len(redEnvelopeUuid) != 0 || len(slotUuid) != 0 || len(bonusBulletShooted["drill"]) != 0 || len(bonusBulletShooted["machineGun"]) != 0 || len(bonusBulletShooted["superMachineGun"]) != 0 {
			flux.Send(accounting.ActionAccountingPause, am.hashId(am.secWebSocketKey), am.to(am.accountingSn))
		}

		tempReceipts[am.accountingSn] = &tempReceipt{
			accountingSn: am.accountingSn,
			initCent:     balance,
			endCent:      0,
			betCent:      0,
			winCent:      0,
			bonusWinCent: 0,
			isDone:       false,
		}

		// updated new accounting sn
		for k := range bulletUuid {
			bulletUuid[k] = am.accountingSn
		}

		for k := range redEnvelopeUuid {
			redEnvelopeUuid[k] = am.accountingSn
		}

		for k := range slotUuid {
			slotUuid[k] = am.accountingSn
		}

		for _, v := range drillUuid {
			for accountingSn, bullets := range v {
				delete(v, accountingSn)
				v[am.accountingSn] = bullets
			}
		}

		for _, v := range machineGunUuid {
			for accountingSn, bullets := range v {
				delete(v, accountingSn)
				v[am.accountingSn] = bullets
			}
		}

		for _, v := range superMachineGunUuid {
			for accountingSn, bullets := range v {
				delete(v, accountingSn)
				v[am.accountingSn] = bullets
			}
		}

		for _, b := range bonusBulletShooted {
			for _, v := range b {
				for accountingSn, bullets := range v {
					delete(v, accountingSn)
					v[am.accountingSn] = bullets
				}
			}
		}

	case actionEndAccounting:
		//check := action.Payload()[0].(bool)
		walletType := action.Payload()[1].(int)
		cent := action.Payload()[2].(uint64)
		accountingSn := action.Payload()[3].(uint64)

		if accountingSn != am.accountingSn {
			return
		}

		if am.accountingSn > 0 {
			flux.Send(accounting.ActionAccountingStop, am.hashId(am.secWebSocketKey), am.to(am.accountingSn), "")

			if walletType == wallet.CategoryHost && am.accountingSn != 0 {
				if cent == 0 {
					betAmt, winAmt := am.getWalletHostAmt(am.accountingSn, tempReceipts)
					wallet.Service.IncreaseBalance(am.secWebSocketKey, "", 0, betAmt, winAmt)
				}
			}
		}

		if am.accountingSn != 0 {
			if _, ok := tempReceipts[am.accountingSn]; ok {
				tempReceipts[am.accountingSn].isDone = true
				totalBet := tempReceipts[am.accountingSn].betCent
				totalBonusWin := tempReceipts[am.accountingSn].bonusWinCent
				totalWin := tempReceipts[am.accountingSn].winCent
				jackpotWin := tempReceipts[am.accountingSn].jackpotWin

				logger.Service.Zap.Infow("Accounting Manager Done",
					"GameUser", am.secWebSocketKey,
					"AccountingSn", am.accountingSn,
					"Bullets", len(bulletUuid),
					"RedEnvelopes", len(redEnvelopeUuid),
					"Slots", len(slotUuid),
					"Drills", len(drillUuid),
					"MachineGuns", len(machineGunUuid),
					"SuperMachineGuns", len(superMachineGunUuid),
					"IsLeaveGameRoom", am.isLeaveGameRoom,
					"IsDisconnect", am.isDisconnect,
					"InitCent", tempReceipts[am.accountingSn].initCent,
					"BetCent", tempReceipts[am.accountingSn].betCent,
					"WinCent", tempReceipts[am.accountingSn].winCent,
					"BonusWinCent", tempReceipts[am.accountingSn].bonusWinCent,
					"JackpotWin", tempReceipts[am.accountingSn].jackpotWin,
				)

				accountingSn := am.accountingSn
				am.accountingSn = 0
				am.timerPause = false

				flux.Send(accounting.ActionAccountingEnd, am.hashId(am.secWebSocketKey), am.to(accountingSn),
					totalBet, totalBonusWin, totalWin, jackpotWin, accountingSn)
			}
		}

	case actionDecrease:
		uuid := action.Payload()[0].(string)
		// this balance is decreased cent
		balance := action.Payload()[1].(uint64)
		cent := action.Payload()[2].(uint64)

		am.checkStart(balance+cent, bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, bonusBulletShooted, tempReceipts)

		accountingSn := am.accountingSn

		// cause by async reason, see 1dd8ed7e
		if _, ok := bulletUuid[uuid]; ok {
			// 2. 核銷子彈 35adc39d
			delete(bulletUuid, uuid)
		} else {
			// 1. 紀錄子彈 35adc39d
			tempReceipts[accountingSn].betCent += cent
			bulletUuid[uuid] = accountingSn
			am.rounds++
		}

		if uuid != "" {
			am.decreaseSn[uuid] = accountingSn
		}

		am.contributeJackpot(cent, false)

		// Process Leave GameRoom and rejoin room
		if am.walletType == wallet.CategoryHost && accountingSn != 0 {
			am.isLeaveGameRoom = false
		}

		logger.Service.Zap.Infow("Accounting Manager Decrease",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", accountingSn,
			"BulletUuid", uuid,
			"Balance", balance+cent,
			"Cent", cent,
			"Bullets", len(bulletUuid),
			"RedEnvelopes", len(redEnvelopeUuid),
			"Slots", len(slotUuid),
			"Drills", len(drillUuid),
			"MachineGuns", len(machineGunUuid),
			"SuperMachineGuns", len(superMachineGunUuid),
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
			"InitCent", tempReceipts[accountingSn].initCent,
			"BetCent", tempReceipts[accountingSn].betCent,
			"WinCent", tempReceipts[accountingSn].winCent,
			"BonusWinCent", tempReceipts[accountingSn].bonusWinCent,
		)
		if am.timerPause && len(redEnvelopeUuid) == 0 && len(slotUuid) == 0 && len(bonusBulletShooted["drill"]) == 0 && len(bonusBulletShooted["machineGun"]) == 0 && len(bonusBulletShooted["superMachineGun"]) == 0 {
			flux.Send(accounting.ActionAccountingContinue, am.hashId(am.secWebSocketKey), am.to(am.accountingSn))
		}

	case actionIncrease:
		uuid := action.Payload()[0].(string)
		// this balance is increased cent
		balance := action.Payload()[1].(uint64)
		cent := action.Payload()[2].(uint64)

		am.checkStart(balance-cent, bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, bonusBulletShooted, tempReceipts)

		accountingSn := am.accountingSn

		tempReceipts[accountingSn].winCent += cent

		if uuid != "" {
			am.increaseSn[uuid] = accountingSn
		}

		logger.Service.Zap.Infow("Accounting Manager Increase",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", accountingSn,
			"BulletUuid", uuid,
			"Balance", balance-cent,
			"Cent", cent,
			"Bullets", len(bulletUuid),
			"RedEnvelopes", len(redEnvelopeUuid),
			"Slots", len(slotUuid),
			"Drills", len(drillUuid),
			"MachineGuns", len(machineGunUuid),
			"SuperMachineGuns", len(superMachineGunUuid),
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
			"InitCent", tempReceipts[accountingSn].initCent,
			"BetCent", tempReceipts[accountingSn].betCent,
			"WinCent", tempReceipts[accountingSn].winCent,
			"BonusWinCent", tempReceipts[accountingSn].bonusWinCent,
		)

		if am.checkEnd(bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, tempReceipts, false) {
			if am.walletType == wallet.CategoryHost {
				am.requestResult(accountingSn, tempReceipts)
			}
			am.isLeaveGameRoom = false
			am.accountingSn = 0
			am.timerPause = false
		}

	case actionBulletHit:
		hitFish := action.Payload()[0].(*fish.Fish)
		hitBullet := action.Payload()[1].(*bullet.Bullet)
		hitResult := action.Payload()[2].(*probability.Probability)
		accountingSn := action.Payload()[3].(uint64)

		am.checkStart(0, bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, bonusBulletShooted, tempReceipts)

		if sn1, ok1 := am.increaseSn[hitBullet.Uuid]; ok1 {
			accountingSn = sn1
			delete(am.increaseSn, hitBullet.Uuid)
			delete(am.decreaseSn, hitBullet.Uuid)
		} else if sn2, ok2 := am.decreaseSn[hitBullet.Uuid]; ok2 {
			accountingSn = sn2
			delete(am.increaseSn, hitBullet.Uuid)
			delete(am.decreaseSn, hitBullet.Uuid)
		} else if accountingSn == 0 {
			accountingSn = am.accountingSn
		}

		isDone := make(chan bool)

		flux.Send(actionAccountingHitBullet, am.hashId(am.secWebSocketKey), am.to(am.accountingSn),
			am.secWebSocketKey,
			hitFish,
			hitBullet,
			hitResult,
			isDone,
			accountingSn,
		)

		<-isDone

		isPause := false //accounting rotation timer

		gameId := gamesetting.Service.GameId(am.secWebSocketKey)
		switch gameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, drillUuid, machineGunUuid, bonusBulletShooted, false)
			isPause = psf_on_00001_BulletHit_TriggerId(hitResult.TriggerIconId, am, hitBullet, hitResult, tempReceipts,
				drillUuid, machineGunUuid, bonusBulletShooted, redEnvelopeUuid, slotUuid,
			)
		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, drillUuid, machineGunUuid, bonusBulletShooted, false)
			isPause = psf_on_00002_BulletHit_TriggerId(hitResult.TriggerIconId, am, hitBullet, hitResult, tempReceipts,
				drillUuid, machineGunUuid, bonusBulletShooted, redEnvelopeUuid, slotUuid,
			)
		case models.PSF_ON_00003:
			psf_on_00003_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, drillUuid, bonusBulletShooted, hitBullet.Bet, false)
			isPause = psf_on_00003_BulletHit_TriggerId(hitResult.TriggerIconId, am, hitBullet, hitResult, tempReceipts, drillUuid, bonusBulletShooted, slotUuid)

		case models.PSF_ON_00004:
			psf_on_00004_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, false)
			isPause = psf_on_00004_BulletHit_TriggerId(hitResult.TriggerIconId, am, hitBullet, hitResult, tempReceipts,
				machineGunUuid, superMachineGunUuid, bonusBulletShooted, redEnvelopeUuid, slotUuid,
			)
		case models.PSF_ON_00005:
			psf_on_00005_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, false)
			isPause = psf_on_00005_BulletHit_TriggerId(hitResult.TriggerIconId, am, hitBullet, hitResult, tempReceipts,
				machineGunUuid, superMachineGunUuid, bonusBulletShooted, redEnvelopeUuid, slotUuid,
			)
		case models.PSF_ON_00006:
			psf_on_00006_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, false)
			isPause = psf_on_00006_BulletHit_TriggerId(hitResult.TriggerIconId, am, hitBullet, hitResult, tempReceipts,
				machineGunUuid, superMachineGunUuid, bonusBulletShooted, redEnvelopeUuid, slotUuid,
			)
		case models.PSF_ON_00007:
			psf_on_00007_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, false)
			isPause = psf_on_00007_BulletHit_TriggerId(hitResult.TriggerIconId, am, hitBullet, hitResult, tempReceipts,
				machineGunUuid, superMachineGunUuid, bonusBulletShooted, redEnvelopeUuid, slotUuid,
			)
		case models.RKF_H5_00001:
			rkf_h5_00001_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, false)
			isPause = rkf_h5_00001_BulletHit_TriggerId(hitResult.TriggerIconId, am, hitBullet, hitResult, tempReceipts,
				machineGunUuid, superMachineGunUuid, bonusBulletShooted, redEnvelopeUuid, slotUuid,
			)
		}

		// Fix tempReceipts[am.accountingSn] is nil Error.
		logger.Service.Zap.Infow("Accounting Manager Hit",
			"GameUser", am.secWebSocketKey,
			"am.AccountingSn", am.accountingSn,
			"AccountingSn", accountingSn,
			"BulletUuid", hitBullet.Uuid,
			"BulletType", hitBullet.TypeId,
			"BulletBonusUuid", hitBullet.BonusUuid,
			"TriggerId", hitResult.TriggerIconId,
			"TriggerBonusUuid", hitResult.BonusUuid,
			"Bullets", len(bulletUuid),
			"RedEnvelopes", len(redEnvelopeUuid),
			"Slots", len(slotUuid),
			"Drills", len(drillUuid),
			"MachineGuns", len(machineGunUuid),
			"SuperMachineGuns", len(superMachineGunUuid),
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
			"InitCent", tempReceipts[accountingSn].initCent,
			"BetCent", tempReceipts[accountingSn].betCent,
			"WinCent", tempReceipts[accountingSn].winCent,
			"BonusWinCent", tempReceipts[accountingSn].bonusWinCent,
		)

		if am.checkEnd(bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, tempReceipts, false) {
			if am.walletType == wallet.CategoryHost {
				am.requestResult(accountingSn, tempReceipts)
			}
			am.isLeaveGameRoom = false
			am.accountingSn = 0
			am.timerPause = false
		} else if isPause && !am.timerPause {
			flux.Send(accounting.ActionAccountingPause, am.hashId(am.secWebSocketKey), am.to(am.accountingSn))
		}

	case actionBulletMultiHit:
		hitFishies := action.Payload()[0].(*fish.Fishies)
		hitBullet := action.Payload()[1].(*bullet.Bullet)
		hitResults := action.Payload()[2].([]*probability.Probability)
		accountingSn := action.Payload()[4].(uint64)

		am.checkStart(0, bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, bonusBulletShooted, tempReceipts)
		if accountingSn == 0 {
			accountingSn = am.accountingSn
		}
		isDone := make(chan bool)
		gameId := gamesetting.Service.GameId(am.secWebSocketKey)

		flux.Send(actionAccountingHitBulletMulti, am.hashId(am.secWebSocketKey), am.to(am.accountingSn),
			am.secWebSocketKey,
			hitFishies,
			hitBullet,
			hitResults,
			isDone,
			accountingSn,
		)
		<-isDone

		isPause := false

		switch gameId {
		case models.PSF_ON_00003:
			psf_on_00003_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, drillUuid, bonusBulletShooted, hitBullet.Bet, false)
			isPause = psf_on_00003_BulletMultiHit_TriggerId(am, hitBullet, hitResults, tempReceipts, drillUuid, bonusBulletShooted, slotUuid)
		}

		logger.Service.Zap.Infow("Accounting Manager Hit",
			"GameUser", am.secWebSocketKey,
			"am.AccountingSn", am.accountingSn,
			"AccountingSn", accountingSn,
			"BulletUuid", hitBullet.Uuid,
			"BulletType", hitBullet.TypeId,
			"BulletBonusUuid", hitBullet.BonusUuid,
			"Bullets", len(bulletUuid),
			"RedEnvelopes", len(redEnvelopeUuid),
			"Slots", len(slotUuid),
			"Drills", len(drillUuid),
			"MachineGuns", len(machineGunUuid),
			"SuperMachineGuns", len(superMachineGunUuid),
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
			"InitCent", tempReceipts[am.accountingSn].initCent,
			"BetCent", tempReceipts[am.accountingSn].betCent,
			"WinCent", tempReceipts[am.accountingSn].winCent,
			"BonusWinCent", tempReceipts[am.accountingSn].bonusWinCent,
		)

		if am.checkEnd(bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, tempReceipts, false) {
			if am.walletType == wallet.CategoryHost {
				am.requestResult(accountingSn, tempReceipts)
			}
			am.isLeaveGameRoom = false
			am.accountingSn = 0
			am.timerPause = false
		} else if isPause && !am.timerPause {
			flux.Send(accounting.ActionAccountingPause, am.hashId(am.secWebSocketKey), am.to(am.accountingSn))
		}

	case actionBulletBonus:
		bonusBullet := action.Payload()[0].(*bullet.Bullet)
		accountingSn := action.Payload()[1].(uint64)

		am.checkStart(0, bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, bonusBulletShooted, tempReceipts)

		gameId := gamesetting.Service.GameId(am.secWebSocketKey)
		switch gameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_BulletBonus(bonusBullet.TypeId, bonusBullet, machineGunUuid, bonusBulletShooted)
		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_BulletBonus(bonusBullet.TypeId, bonusBullet, machineGunUuid, bonusBulletShooted)
		case models.PSF_ON_00003:
			psf_on_00003_BulletBonus(bonusBullet.TypeId)
		case models.PSF_ON_00004:
			psf_on_00004_BulletBonus(bonusBullet.TypeId, bonusBullet, machineGunUuid, superMachineGunUuid, bonusBulletShooted)
		case models.PSF_ON_00005:
			psf_on_00005_BulletBonus(bonusBullet.TypeId, bonusBullet, machineGunUuid, superMachineGunUuid, bonusBulletShooted)
		case models.PSF_ON_00006:
			psf_on_00006_BulletBonus(bonusBullet.TypeId, bonusBullet, machineGunUuid, superMachineGunUuid, bonusBulletShooted)
		case models.PSF_ON_00007:
			psf_on_00007_BulletBonus(bonusBullet.TypeId, bonusBullet, machineGunUuid, superMachineGunUuid, bonusBulletShooted)
		case models.RKF_H5_00001:
			rkf_h5_00001_BulletBonus(bonusBullet.TypeId, bonusBullet, machineGunUuid, superMachineGunUuid, bonusBulletShooted)
		}

		isDone := make(chan bool)

		flux.Send(actionAccountingBonus, am.hashId(am.secWebSocketKey), am.to(am.accountingSn),
			am.secWebSocketKey, bonusBullet, isDone, accountingSn)

		<-isDone

	case actionPickBonus:
		bonus := action.Payload()[0]
		var accountingSn uint64 = 0

		am.checkStart(0, bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, bonusBulletShooted, tempReceipts)

		switch bonus.(type) {

		case *redenvelope.RedEnvelope:
			// 2. 核銷紅包 858cedf5
			re := bonus.(*redenvelope.RedEnvelope)

			if sn, ok := am.increaseSn[re.Uuid]; ok {
				accountingSn = sn
				delete(am.increaseSn, re.Uuid)
			} else if accountingSn == 0 {
				accountingSn = am.accountingSn
			}

			tempReceipts[accountingSn].bonusWinCent += re.Bet * re.Pay * re.Rate

			isDone := make(chan bool)

			flux.Send(actionAccountingBonus, am.hashId(am.secWebSocketKey), am.to(am.accountingSn),
				am.secWebSocketKey, re, isDone, accountingSn)

			<-isDone

			delete(redEnvelopeUuid, re.Uuid)

		case *slot.Slot:
			// 2. 核銷老虎機 6e01e01a
			s := bonus.(*slot.Slot)

			if sn, ok := am.increaseSn[s.Uuid]; ok {
				accountingSn = sn
				delete(am.increaseSn, s.Uuid)
			} else if accountingSn == 0 {
				accountingSn = am.accountingSn
			}

			tempReceipts[accountingSn].bonusWinCent += s.Bet * s.Pay * s.Rate

			// Process Multi Slot Write Off
			gameId := gamesetting.Service.GameId(s.SecWebSocketKey)
			if gameId != models.PSF_ON_00003 {
				isDone := make(chan bool)
				flux.Send(actionAccountingBonus, am.hashId(am.secWebSocketKey), am.to(am.accountingSn),
					am.secWebSocketKey, s, isDone, accountingSn)
				<-isDone
			}

			delete(slotUuid, s.Uuid)
		}

		if am.checkEnd(bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, tempReceipts, false) {
			if am.walletType == wallet.CategoryHost {
				am.requestResult(accountingSn, tempReceipts)
			}
			am.isLeaveGameRoom = false
			am.accountingSn = 0
			am.timerPause = false
		} else if am.timerPause && len(redEnvelopeUuid) == 0 && len(slotUuid) == 0 && len(bonusBulletShooted["drill"]) == 0 && len(bonusBulletShooted["machineGun"]) == 0 && len(bonusBulletShooted["superMachineGun"]) == 0 {
			flux.Send(accounting.ActionAccountingContinue, am.hashId(am.secWebSocketKey), am.to(am.accountingSn))
		}

	case actionDisconnect:
		done := action.Payload()[0].(chan bool)

		logger.Service.Zap.Infow("Accounting Manager Disconnect",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", am.accountingSn,
			"Bullets", len(bulletUuid),
			"RedEnvelopes", len(redEnvelopeUuid),
			"Slots", len(slotUuid),
			"Drills", len(drillUuid),
			"MachineGuns", len(machineGunUuid),
			"SuperMachineGuns", len(superMachineGunUuid),
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
		)

		// waiting lottery to finish ActionBulletHitRecall
		if len(bulletUuid) > 0 {
			done <- false
			return
		}

		am.isDisconnect = true

		refundDone := make(chan bool, 2014)
		flux.Send(refund.ActionRefundCall, am.hashId(am.secWebSocketKey), refund.Service.Id,
			am.secWebSocketKey, true, refundDone, am.accountingSn)

		if <-refundDone {
			if am.accountingSn > 0 {
				flux.Send(accounting.ActionAccountingStop, am.hashId(am.secWebSocketKey), am.to(am.accountingSn), "")
			}

			if len(bulletUuid) == 0 &&
				len(redEnvelopeUuid) == 0 &&
				len(slotUuid) == 0 &&
				len(drillUuid) == 0 &&
				len(machineGunUuid) == 0 &&
				len(superMachineGunUuid) == 0 &&
				(am.isLeaveGameRoom || am.isDisconnect) &&
				am.accountingSn == 0 {
				done <- true
				return
			}

			accountingSn := am.accountingSn
			if am.checkEnd(bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, tempReceipts, false) {
				if am.walletType == wallet.CategoryHost {
					am.requestResult(accountingSn, tempReceipts)
				}

				am.isLeaveGameRoom = false
				am.accountingSn = 0
				am.timerPause = false
				done <- true
				return
			}

			done <- false
		}

	case actionLeaveGameRoom:
		done := action.Payload()[0].(chan bool)
		refundDone := make(chan bool, 2014)

		flux.Send(refund.ActionRefundCall, am.hashId(am.secWebSocketKey), refund.Service.Id,
			am.secWebSocketKey, false, refundDone, am.accountingSn)

		if <-refundDone {
			// didn't shoot any bullet yet
			if am.accountingSn == 0 && len(bulletUuid) == 0 {
				done <- true
				return
			}

			if am.accountingSn > 0 {
				flux.Send(accounting.ActionAccountingStop, am.hashId(am.secWebSocketKey), am.to(am.accountingSn), "")
			}

			accountingSn := am.accountingSn
			if am.checkEnd(bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, tempReceipts, false) {
				if am.walletType == wallet.CategoryHost {
					am.requestResult(accountingSn, tempReceipts)
				}

				am.isLeaveGameRoom = false
				am.accountingSn = 0
				am.timerPause = false
				done <- true
				return
			}

			done <- false
		}

	case actionJackpotHit:
		jackpotWin := action.Payload()[0].(uint64)

		am.checkStart(0, bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, bonusBulletShooted, tempReceipts)

		tempReceipts[am.accountingSn].jackpotWin += jackpotWin

		logger.Service.Zap.Infow("Accounting Jackpot Hit",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", am.accountingSn,
			"Bullets", len(bulletUuid),
			"RedEnvelopes", len(redEnvelopeUuid),
			"Slots", len(slotUuid),
			"Drills", len(drillUuid),
			"MachineGuns", len(machineGunUuid),
			"SuperMachineGuns", len(superMachineGunUuid),
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
			"InitCent", tempReceipts[am.accountingSn].initCent,
			"BetCent", tempReceipts[am.accountingSn].betCent,
			"WinCent", tempReceipts[am.accountingSn].winCent,
			"BonusWinCent", tempReceipts[am.accountingSn].bonusWinCent,
			"JackpotWin", tempReceipts[am.accountingSn].jackpotWin,
		)

		if am.checkEnd(bulletUuid, redEnvelopeUuid, slotUuid, drillUuid, machineGunUuid, superMachineGunUuid, tempReceipts, false) {
			if am.walletType == wallet.CategoryHost {
				am.requestResult(am.accountingSn, tempReceipts)
			}
			am.isLeaveGameRoom = false
			am.accountingSn = 0
			am.timerPause = false
		}

	case actionWalletHostAmt:
		ch := action.Payload()[0].(chan string)

		betAmt, winAmt := am.getWalletHostAmt(am.accountingSn, tempReceipts)
		output := strconv.FormatUint(betAmt, 10) + " " + strconv.FormatUint(winAmt, 10)

		ch <- output

	case actionBonusBulletShooted:
		hitBullet := action.Payload()[0].(*bullet.Bullet)

		gameId := gamesetting.Service.GameId(am.secWebSocketKey)
		switch gameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, drillUuid, machineGunUuid, bonusBulletShooted, true)

		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, drillUuid, machineGunUuid, bonusBulletShooted, true)

		case models.PSF_ON_00003:
			psf_on_00003_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, drillUuid, bonusBulletShooted, hitBullet.Bet, true)

		case models.PSF_ON_00004:
			psf_on_00004_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, true)

		case models.PSF_ON_00005:
			psf_on_00005_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, true)

		case models.PSF_ON_00006:
			psf_on_00006_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, true)

		case models.PSF_ON_00007:
			psf_on_00007_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, true)

		case models.RKF_H5_00001:
			rkf_h5_00001_BulletHit(hitBullet.TypeId, bulletUuid, am, hitBullet, tempReceipts, machineGunUuid, superMachineGunUuid, bonusBulletShooted, true)

		}
		logger.Service.Zap.Infow("Accounting Manager Bonus Bullet Shooted",
			"GameUser", am.secWebSocketKey,
			"am.AccountingSn", am.accountingSn,
			"BulletUuid", hitBullet.Uuid,
			"BulletType", hitBullet.TypeId,
			"BulletBonusUuid", hitBullet.BonusUuid,
			"Bullets", len(bulletUuid),
			"RedEnvelopes", len(redEnvelopeUuid),
			"Slots", len(slotUuid),
			"Drills", len(drillUuid),
			"MachineGuns", len(machineGunUuid),
			"SuperMachineGuns", len(superMachineGunUuid),
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
		)
	}
}

func (am *AccountingManager) FluxActionHandler(action *flux.Action, a *accounting.Accounting, isOk chan bool) {
	if action.Key().To() != a.HashId() {
		return
	}

	switch action.Key().Name() {
	case flux.ActionFluxRegisterDone:
		isOk <- true

	case accounting.ActionAccountingStop:
		if !a.IsStop {
			a.IsStop = true
			a.CtrlRotation <- true
			close(a.CtrlRotation)
		}
	case accounting.ActionAccountingPause:
		if !a.TimerStop {
			logger.Service.Zap.Infow("Accounting Manager Timer Pause",
				"GameUser", am.secWebSocketKey,
				"AccountingSn", am.accountingSn,
			)
			am.timerPause = true
			a.PauseTimer <- true
		}

	case accounting.ActionAccountingContinue:
		if !a.TimerStop {
			logger.Service.Zap.Infow("Accounting Manager Timer Continue",
				"GameUser", am.secWebSocketKey,
				"AccountingSn", am.accountingSn,
			)
			am.timerPause = false
			a.ContinueTimer <- true
		}

	case accounting.ActionAccountingEnd:
		totalBet := action.Payload()[0].(uint64)
		totalBonusWin := action.Payload()[1].(uint64)
		totalWin := action.Payload()[2].(uint64)
		jackpotWin := action.Payload()[3].(uint64)
		accountingSn := action.Payload()[4].(uint64)
		gameId := gamesetting.Service.GameId(am.secWebSocketKey)

		if a.End(totalBet, totalBonusWin, int64(totalWin), jackpotWin, "", "") {
			if wallet.Service.WalletCategory(a.SecWebSocketKey) == wallet.CategoryHost {
				a.DeleteRefundLog()
			}
			logger.Service.Zap.Infow("QA log ----->",
				"GameUser", am.secWebSocketKey,
				"AccountingSn", accountingSn,
				"TotalWin", int(am.depositCent)-int(totalBet)+int(totalWin),
			)
			value := vip.Service.CheckVip(true, am.secWebSocketKey, am.hostExtId, am.hostId, am.memberId, gameId, am.billCurrency, 1, totalBet, totalWin, accountingSn)
			flux.Send(ActionUpdateBlacklist, am.hashId(am.secWebSocketKey), am.secWebSocketKey, value)
		}

	case accounting.ActionAccountingDelete:
		accountingSn := action.Payload()[0].(uint64)

		a.Delete(accountingSn)

	case actionAccountingHitBullet:
		secWebSocketKey := action.Payload()[0].(string)
		isDone := action.Payload()[4].(chan bool)

		switch a.Start.GameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_HitBullet(action, am)

		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_HitBullet(action, am)

		case models.PSF_ON_00003:
			psf_on_00003_HitBullet(action, am)

		case models.PSF_ON_00004:
			psf_on_00004_HitBullet(action, am)

		case models.PSF_ON_00005:
			psf_on_00005_HitBullet(action, am)

		case models.PSF_ON_00006:
			psf_on_00006_HitBullet(action, am)

		case models.PSF_ON_00007:
			psf_on_00007_HitBullet(action, am)

		case models.RKF_H5_00001:
			rkf_h5_00001_HitBullet(action, am)

		default:
			logger.Service.Zap.Errorw(accounting.Accounting_GAME_ID_INVALID,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, accounting.Accounting_GAME_ID_INVALID)
		}

		isDone <- true
		close(isDone)

	case actionAccountingHitBulletMulti:
		secWebSocketKey := action.Payload()[0].(string)
		isDone := action.Payload()[4].(chan bool)

		switch a.Start.GameId {
		case models.PSF_ON_00003:
			psf_on_00003_MultiHitBullet(action, am)

		default:
			logger.Service.Zap.Errorw(accounting.Accounting_GAME_ID_INVALID,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, accounting.Accounting_GAME_ID_INVALID)
		}

		isDone <- true
		close(isDone)

	case actionAccountingBonus:
		secWebSocketKey := action.Payload()[0].(string)
		isDone := action.Payload()[2].(chan bool)

		switch a.Start.GameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_Bonus(action, am)

		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_Bonus(action, am)

		case models.PSF_ON_00003:
			psf_on_00003_Bonus(action, am)

		case models.PSF_ON_00004:
			psf_on_00004_Bonus(action, am)

		case models.PSF_ON_00005:
			psf_on_00005_Bonus(action, am)

		case models.PSF_ON_00006:
			psf_on_00006_Bonus(action, am)

		case models.PSF_ON_00007:
			psf_on_00007_Bonus(action, am)

		case models.RKF_H5_00001:
			rkf_h5_00001_Bonus(action, am)

		default:
			logger.Service.Zap.Errorw(accounting.Accounting_GAME_ID_INVALID,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, accounting.Accounting_GAME_ID_INVALID)
		}

		isDone <- true
		close(isDone)

	case actionAccountingRefundBonusBullet:
		secWebSocketKey := action.Payload()[0].(string)
		isDone := action.Payload()[2].(chan bool)

		switch a.Start.GameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_RefundBonusBullet(action, am)

		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_RefundBonusBullet(action, am)

		case models.PSF_ON_00003:
			psf_on_00003_RefundBonusBullet(action, am)

		case models.PSF_ON_00004:
			psf_on_00004_RefundBonusBullet(action, am)

		case models.PSF_ON_00005:
			psf_on_00005_RefundBonusBullet(action, am)

		case models.PSF_ON_00006:
			psf_on_00006_RefundBonusBullet(action, am)

		case models.PSF_ON_00007:
			psf_on_00007_RefundBonusBullet(action, am)

		case models.RKF_H5_00001:
			rkf_h5_00001_RefundBonusBullet(action, am)

		default:
			logger.Service.Zap.Errorw(accounting.Accounting_GAME_ID_INVALID,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, accounting.Accounting_GAME_ID_INVALID)
		}

		isDone <- true
		close(isDone)
	}
}

func (am *AccountingManager) RefundHandler(bonus interface{}, accountingSn uint64) {
	if accountingSn == 0 {
		accountingSn = am.accountingSn
	}

	flux.Send(actionRefund, am.hashId(am.secWebSocketKey), am.hashId(am.secWebSocketKey), bonus, accountingSn)
}

func (am *AccountingManager) RotationHandler(a *accounting.Accounting) {
	isDone := make(chan bool)
	am.RotatiomRwMutex.Lock()

	flux.Send(actionRotation, am.hashId(am.secWebSocketKey), am.hashId(am.secWebSocketKey), a.Start.Sn, isDone)

	if !<-isDone {
		for {
			time.Sleep(2 * time.Second)
			if a.IsStop {
				break
			}
			logger.Service.Zap.Infow("Accounting Manager Rotation Again",
				"GameUser", am.secWebSocketKey)

			am.RotatiomRwMutex.Lock()
			flux.Send(actionRotation, am.hashId(am.secWebSocketKey), am.hashId(am.secWebSocketKey), a.Start.Sn, isDone)

			if <-isDone {
				break
			}
		}
	}
	close(isDone)
}

func (am *AccountingManager) DecreaseHandler(secWebSocketKey, uuid string, balance, cent, accountingSn uint64) {
	if accountingSn == 0 {
		accountingSn = am.accountingSn
	}

	flux.Send(actionDecrease, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey),
		uuid,
		balance,
		cent,
		accountingSn,
	)
}

func (am *AccountingManager) GetAccountingSnHandler(secWebSocketKey string, balance, cent, sn uint64, depositCheck bool) (accountingSn uint64) {
	if depositCheck {
		am.depositCent = cent
	}
	if am.accountingSn == 0 || (depositCheck && am.accountingSn == sn) {
		newSn := am.renew(balance)

		flux.Send(actionGetAccountingSn, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey),
			balance,
			newSn,
		)

		return newSn
	}

	return am.accountingSn
}

// todo top one
//
//	func (am *AccountingManager) EndAccountingHandler(secWebSocketKey string, checkAccountingEnd bool, walletType int, cent, sn uint64) {
//		flux.Send(actionEndAccounting, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey), checkAccountingEnd, walletType, cent, sn)
//	}
func (am *AccountingManager) EndAccountingHandler(secWebSocketKey string, checkAccountingEnd bool, walletType int, cent uint64) {
	flux.Send(actionEndAccounting, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey), checkAccountingEnd, walletType, cent)
}

func (am *AccountingManager) IncreaseHandler(secWebSocketKey, uuid string, balance, cent, accountingSn uint64) {
	if accountingSn == 0 {
		accountingSn = am.accountingSn
	}

	flux.Send(actionIncrease, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey),
		uuid,
		balance,
		cent,
		accountingSn,
	)
}

func (am *AccountingManager) BulletHitHandler(secWebSocketKey string, hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability, accountingSn uint64) {
	flux.Send(actionBulletHit, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey),
		hitFish,
		hitBullet,
		hitResult,
		accountingSn,
	)
}

func (am *AccountingManager) BulletMultiHitHandler(secWebSocketKey string, hitFishies *fish.Fishies, hitBullet *bullet.Bullet, hitResults []*probability.Probability, oneBet, accountingSn uint64) {
	if accountingSn == 0 {
		accountingSn = am.accountingSn
	}

	flux.Send(actionBulletMultiHit, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey),
		hitFishies,
		hitBullet,
		hitResults,
		oneBet,
		accountingSn,
	)
}

func (am *AccountingManager) BulletBonusHandler(secWebSocketKey string, bonusBullet *bullet.Bullet, accountingSn uint64) {
	flux.Send(actionBulletBonus, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey), bonusBullet, accountingSn)
}

func (am *AccountingManager) BonusBulletShootedHandler(secWebSocketKey string, bonusBullet *bullet.Bullet) {
	flux.Send(actionBonusBulletShooted, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey), bonusBullet)
}

func (am *AccountingManager) BetErrorHandler(secWebSocketKey string, accountingSn uint64) {
	if am.accountingSn > 0 {
		flux.Send(accounting.ActionAccountingStop, am.hashId(secWebSocketKey), am.to(accountingSn), "")
		flux.Send(accounting.ActionAccountingDelete, am.hashId(secWebSocketKey), am.to(accountingSn), accountingSn)
		am.accountingSn = 0
		am.timerPause = false
	}
}

func (am *AccountingManager) PickBonus(secWebSocketKey string, bonus interface{}) {
	flux.Send(actionPickBonus, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey), bonus)
}

func (am *AccountingManager) LeaveGameRoom(done chan bool) {
	am.rwMutex.Lock()
	defer am.rwMutex.Unlock()

	am.isLeaveGameRoom = true

	flux.Send(actionLeaveGameRoom, am.hashId(am.secWebSocketKey), am.hashId(am.secWebSocketKey), done)
}

func (am *AccountingManager) Disconnect(done chan bool) {
	am.rwMutex.Lock()
	defer am.rwMutex.Unlock()

	am.contributeJackpot(0, true)

	flux.Send(actionDisconnect, am.hashId(am.secWebSocketKey), am.hashId(am.secWebSocketKey), done)
}

func (am *AccountingManager) Rounds() uint64 {
	am.rwMutex.RLock()
	defer am.rwMutex.RUnlock()

	return am.rounds
}

func (am *AccountingManager) JackpotHit(secWebSocketKey string, jackpotWin uint64) {
	flux.Send(actionJackpotHit, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey), jackpotWin)
}

func (am *AccountingManager) GetWalletHostAmt(secWebSocketKey string) (betAmt, winAmt uint64) {
	ch := make(chan string, 0)
	flux.Send(actionWalletHostAmt, am.hashId(secWebSocketKey), am.hashId(secWebSocketKey), ch)

	data := <-ch
	results := strings.Split(data, " ")
	betAmt, _ = strconv.ParseUint(results[0], 10, 64)
	winAmt, _ = strconv.ParseUint(results[1], 10, 64)

	return betAmt, winAmt
}

func (am *AccountingManager) checkStart(
	balance uint64,
	bulletUuid, redEnvelopeUuid, slotUuid map[string]uint64,
	drillUuid, machineGunUuid, superMachineGunUuid map[string]map[uint64]int32, bonusBulletShooted map[string]map[string]map[uint64]int32,
	tempReceipts map[uint64]*tempReceipt) {

	if am.accountingSn == 0 {
		// Get Balance
		if balance == 0 {
			balance = am.balance()
		}
		am.accountingSn = am.renew(balance)
		if len(redEnvelopeUuid) != 0 || len(slotUuid) != 0 || len(bonusBulletShooted["drill"]) != 0 || len(bonusBulletShooted["machineGun"]) != 0 || len(bonusBulletShooted["superMachineGun"]) != 0 {
			flux.Send(accounting.ActionAccountingPause, am.hashId(am.secWebSocketKey), am.to(am.accountingSn))
		}

		if am.walletType == wallet.CategoryHost { //理論上:換單->斷線->Increase會進
			betInfo := bet.Service.Get(am.secWebSocketKey)
			gameId := gamesetting.Service.GameId(am.secWebSocketKey)
			bet := gamesetting.Service.Bet(am.secWebSocketKey, gameId, betInfo.BetLevelIndex, betInfo.BetIndex)
			rate := gamesetting.Service.Rate(am.secWebSocketKey)

			balance, _ := wallet.Service.BalanceCache(am.secWebSocketKey)

			depositCent := bet * rate * am.depositMultiple
			if balance < depositCent {
				depositCent = balance
			}
			wallet.Service.Deposit(am.secWebSocketKey, depositCent, 0)
		}

		logger.Service.Zap.Infow("Accounting Manager CheckStart",
			"GameUser", am.secWebSocketKey,
			"NewAccountingSn", am.accountingSn,
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
		)

		tempReceipts[am.accountingSn] = &tempReceipt{
			accountingSn: am.accountingSn,
			initCent:     balance,
			endCent:      0,
			betCent:      0,
			winCent:      0,
			bonusWinCent: 0,
			isDone:       false,
		}

		// updated new accounting sn
		for k := range bulletUuid {
			bulletUuid[k] = am.accountingSn
		}

		for k := range redEnvelopeUuid {
			redEnvelopeUuid[k] = am.accountingSn
		}

		for k := range slotUuid {
			slotUuid[k] = am.accountingSn
		}

		for _, v := range drillUuid {
			for accountingSn, bullets := range v {
				delete(v, accountingSn)
				v[am.accountingSn] = bullets
			}
		}

		for _, v := range machineGunUuid {
			for accountingSn, bullets := range v {
				delete(v, accountingSn)
				v[am.accountingSn] = bullets
			}
		}

		for _, v := range superMachineGunUuid {
			for accountingSn, bullets := range v {
				delete(v, accountingSn)
				v[am.accountingSn] = bullets
			}
		}

		for _, b := range bonusBulletShooted {
			for _, v := range b {
				for accountingSn, bullets := range v {
					delete(v, accountingSn)
					v[am.accountingSn] = bullets
				}
			}
		}
	}
}

func (am *AccountingManager) checkEnd(
	bulletUuid, redEnvelopeUuid, slotUuid map[string]uint64,
	drillUuid, machineGunUuid, superMachineGunUuid map[string]map[uint64]int32,
	tempReceipts map[uint64]*tempReceipt, checkWalletHost bool) bool {

	if len(bulletUuid) == 0 &&
		len(redEnvelopeUuid) == 0 &&
		len(slotUuid) == 0 &&
		len(drillUuid) == 0 &&
		len(machineGunUuid) == 0 &&
		len(superMachineGunUuid) == 0 &&
		(((am.isLeaveGameRoom || am.isDisconnect) || checkWalletHost) && am.accountingSn != 0) {

		tempReceipts[am.accountingSn].isDone = true

		totalBet := tempReceipts[am.accountingSn].betCent
		totalBonusWin := tempReceipts[am.accountingSn].bonusWinCent
		totalWin := tempReceipts[am.accountingSn].winCent
		jackpotWin := tempReceipts[am.accountingSn].jackpotWin

		logger.Service.Zap.Infow("Accounting Manager Done",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", am.accountingSn,
			"Bullets", len(bulletUuid),
			"RedEnvelopes", len(redEnvelopeUuid),
			"Slots", len(slotUuid),
			"Drills", len(drillUuid),
			"MachineGuns", len(machineGunUuid),
			"SuperMachineGuns", len(superMachineGunUuid),
			"IsLeaveGameRoom", am.isLeaveGameRoom,
			"IsDisconnect", am.isDisconnect,
			"InitCent", tempReceipts[am.accountingSn].initCent,
			"BetCent", tempReceipts[am.accountingSn].betCent,
			"WinCent", tempReceipts[am.accountingSn].winCent,
			"BonusWinCent", tempReceipts[am.accountingSn].bonusWinCent,
			"JackpotWin", tempReceipts[am.accountingSn].jackpotWin,
		)

		flux.Send(accounting.ActionAccountingEnd, am.hashId(am.secWebSocketKey), am.to(am.accountingSn),
			totalBet, totalBonusWin, totalWin, jackpotWin, am.accountingSn)

		return true
	}
	return false
}

func (am *AccountingManager) balance() uint64 {
	balance, err := wallet.Service.Balance(am.secWebSocketKey)

	if err != nil {
		logger.Service.Zap.Errorw(AccountingManager_GET_BALANCE_FAILED,
			"GameUser", am.secWebSocketKey,
			"AccountingSn", am.accountingSn,
			"Error", err,
		)
		errorcode.Service.Fatal(am.secWebSocketKey, AccountingManager_GET_BALANCE_FAILED)
		return 0
	}

	return balance
}

func (am *AccountingManager) renew(balance uint64) uint64 {
	if a := accounting.Builder().
		SetSecWebSocketKey(am.secWebSocketKey).
		SetHostExtId(am.hostExtId).
		SetGuestEnable(am.guestEnable).
		SetPeriod(gamesetting.Service.AccountingPeriod(am.secWebSocketKey)).
		SetHostId(am.hostId).
		SetGameId(gamesetting.Service.GameId(am.secWebSocketKey)).
		SetSubgameId(gamesetting.Service.RateIndex(am.secWebSocketKey)).
		SetMemberId(am.memberId).
		SetInitCent(balance).
		SetDenom(gamesetting.Service.Rate(am.secWebSocketKey)).
		SetGameData(strconv.Itoa(int(gamesetting.Service.Rate(am.secWebSocketKey)))).
		SetGameModule(gamesetting.Service.MathModuleId(am.secWebSocketKey)).
		SetDB(am.hostExtId).
		SetKind(accounting.PLAYER).
		SetTestAccount(wallet.Service.AccountingType(am.secWebSocketKey)).
		SetHandler(am).
		Build(); a != nil {

		return a.Start.Sn
	}
	return 0
}

func (am *AccountingManager) contributeJackpot(cent uint64, offline bool) {
	gameId := gamesetting.Service.GameId(am.secWebSocketKey)

	if !am.guestEnable {
		switch gameId {
		case models.PSF_ON_20002:
			flux.Send(jackpot.ActionJackpotContribute, am.hashId(am.secWebSocketKey), jackpot.Service.Id,
				cent,
				am.hostId,
				am.hostExtId,
				am.memberId,
				am.secWebSocketKey,
				am.accountingSn,
				offline,
				gameId,
				gamesetting.Service.SubgameId(am.secWebSocketKey),
				gamesetting.Service.JackpotGroup(am.secWebSocketKey),
			)
		}
	}
}

func (am *AccountingManager) getWalletHostAmt(accountingSn uint64, tempReceipts map[uint64]*tempReceipt) (betAmt, winAmt uint64) {
	betAmt = 0
	winAmt = 0

	if _, ok := tempReceipts[accountingSn]; ok {
		betAmt = tempReceipts[accountingSn].betCent
		winAmt = tempReceipts[accountingSn].winCent
	}

	return betAmt, winAmt
}

func (am *AccountingManager) requestResult(accountingSn uint64, tempReceipts map[uint64]*tempReceipt) {
	betAmt, winAmt := am.getWalletHostAmt(accountingSn, tempReceipts)
	// those func should be V2
	cent, _ := wallet.Service.CentCache(am.secWebSocketKey)
	if cent == 0 {
		wallet.Service.IncreaseBalance(am.secWebSocketKey, "", 0, betAmt, winAmt)
	} else {
		wallet.Service.Withdraw(am.secWebSocketKey, accountingSn, cent, betAmt, winAmt)
	}
}

func (am *AccountingManager) to(accountingSn uint64) string {
	if accountingSn == 0 {
		panic("BUG SN == 0" + am.secWebSocketKey)
	}

	return accounting.Service.HashId(am.secWebSocketKey, accountingSn)
}

func (am *AccountingManager) hashId(secWebSocketKey string) string {
	return Service.hashId(secWebSocketKey)
}

func (am *AccountingManager) destroy() {
	flux.UnRegister(am.hashId(am.secWebSocketKey), am.in)
}
