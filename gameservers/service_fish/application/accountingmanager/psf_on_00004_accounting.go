package accountingmanager

import (
	"fmt"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	"serve/service_fish/domain/receipt"
	"serve/service_fish/domain/redenvelope"
	"serve/service_fish/domain/slot"
	"strconv"

	"serve/fish_comm/accounting"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/fish_comm/wallet"

	"github.com/google/uuid"
)

func psf_on_00004_HitBullet(action *flux.Action, am *AccountingManager) {
	if am.guestEnable {
		return
	}

	secWebSocketKey := action.Payload()[0].(string)
	hitFish := action.Payload()[1].(*fish.Fish)
	hitBullet := action.Payload()[2].(*bullet.Bullet)
	hitResult := action.Payload()[3].(*probability.Probability)
	accountingSn := action.Payload()[5].(uint64)

	if wallet.Service.WalletCategory(secWebSocketKey) != wallet.CategoryHost {
		accountingSn = am.accountingSn
	}

	// Process XiaoLongBao
	if hitResult.TriggerIconId != 100 && hitResult.TriggerIconId != 101 && hitResult.TriggerIconId != 102 &&
		hitResult.TriggerIconId > 0 && hitResult.TriggerIconId != 302 {
		logger.Service.Zap.Infow("Trigger Bullet Accounting",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", accountingSn,
		)
		return
	}

	gameResult := fmt.Sprintf("%d/%d/%d/%d/%d/%d",
		hitBullet.TypeId,
		hitBullet.Bullets,
		hitBullet.UsedBullets,
		hitFish.TypeId,
		hitFish.Scene,
		hitResult.Multiplier,
	)

	bet := hitBullet.Bet * hitBullet.Rate
	gameData := ""
	if hitBullet.TypeId == 201 || hitBullet.TypeId == 202 {
		gameData = strconv.Itoa(int(bet))
		bet = 0
	}

	go receipt.Service.New(
		hitBullet.Uuid,
		secWebSocketKey,
		am.hostExtId,
		gameResult,
		gameData,
		accountingSn,
		bet,
		uint64(hitResult.Pay)*hitBullet.Bet*hitBullet.Rate*uint64(hitResult.Multiplier),
	)

	logger.Service.Zap.Infow("Bullet Accounting",
		"GameUser", am.secWebSocketKey,
		"AccountingSn", accountingSn,
	)
}

func psf_on_00004_Bonus(action *flux.Action, am *AccountingManager) {
	if am.guestEnable {
		return
	}

	secWebSocketKey := action.Payload()[0].(string)
	accountingSn := am.accountingSn

	switch action.Payload()[1].(type) {
	case *bullet.Bullet:
		b := action.Payload()[1].(*bullet.Bullet)

		hitFish := b.ExtraData[0].(*fish.Fish)
		hitBullet := b.ExtraData[1].(*bullet.Bullet)
		hitResult := b.ExtraData[2].(*probability.Probability)

		gameResult := fmt.Sprintf("%d/%d/%d/%d/%d/%d",
			hitBullet.TypeId,
			hitBullet.Bullets,
			hitBullet.UsedBullets,
			hitFish.TypeId,
			hitFish.Scene,
			hitResult.Multiplier,
		)

		bet := b.Bet * b.Rate

		if hitBullet.TypeId == 201 || hitBullet.TypeId == 202 {
			bet = 0
		}

		gameData := fmt.Sprintf("%d/%d", b.TypeId, b.Bullets)

		go receipt.Service.New(
			b.BonusUuid,
			secWebSocketKey,
			am.hostExtId,
			gameResult,
			gameData,
			accountingSn,
			bet,
			uint64(hitResult.Pay)*b.Bet*b.Rate,
		)

		logger.Service.Zap.Infow("Bonus Bullet Accounting",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", accountingSn,
		)

	case *redenvelope.RedEnvelope:
		re := action.Payload()[1].(*redenvelope.RedEnvelope)

		hitFish := re.ExtraData[0].(*fish.Fish)
		hitBullet := re.ExtraData[1].(*bullet.Bullet)
		hitResult := re.ExtraData[2].(*probability.Probability)

		gameResult := fmt.Sprintf("%d/%d/%d/%d/%d/%d",
			hitBullet.TypeId,
			hitBullet.Bullets,
			hitBullet.UsedBullets,
			hitFish.TypeId,
			hitFish.Scene,
			hitResult.Multiplier,
		)

		bet := re.Bet * re.Rate

		if hitResult.TriggerIconId == 100 || hitResult.TriggerIconId == 102 || hitBullet.TypeId == 201 || hitBullet.TypeId == 202 {
			bet = 0
		}

		gameData := fmt.Sprintf("%d/%d/%d/%d/%d/%d",
			re.PlayerOptIndex, re.AllPay[0], re.AllPay[1], re.AllPay[2], re.AllPay[3], re.AllPay[4])

		go receipt.Service.New(
			re.Uuid,
			secWebSocketKey,
			am.hostExtId,
			gameResult,
			gameData,
			accountingSn,
			bet,
			//(re.Pay+uint64(hitResult.Pay))*re.Bet*re.Rate,
			re.Pay*re.Bet*re.Rate,
		)

		logger.Service.Zap.Infow("RedEnvelop Accounting",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", accountingSn,
		)

	case *slot.Slot:
		s := action.Payload()[1].(*slot.Slot)

		hitFish := s.ExtraData[0].(*fish.Fish)
		hitBullet := s.ExtraData[1].(*bullet.Bullet)
		hitResult := s.ExtraData[2].(*probability.Probability)

		gameResult := fmt.Sprintf("%d/%d/%d/%d/%d/%d",
			hitBullet.TypeId,
			hitBullet.Bullets,
			hitBullet.UsedBullets,
			hitFish.TypeId,
			hitFish.Scene,
			hitResult.Multiplier,
		)

		bet := s.Bet * s.Rate

		if hitResult.TriggerIconId == 101 || hitBullet.TypeId == 201 || hitBullet.TypeId == 202 {
			bet = 0
		}

		gameData := fmt.Sprintf("%d_%d/%d_%d/%d_%d",
			s.AllPay[0], s.Reels[1],
			s.AllPay[1], s.Reels[4],
			s.AllPay[2], s.Reels[7],
		)

		go receipt.Service.New(
			s.Uuid,
			secWebSocketKey,
			am.hostExtId,
			gameResult,
			gameData,
			accountingSn,
			bet,
			s.Pay*s.Bet*s.Rate,
		)

		logger.Service.Zap.Infow("Slot Accounting",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", accountingSn,
		)

	default:
		logger.Service.Zap.Errorw(accounting.Accounting_BONUS_TYPE_INVALID,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, accounting.Accounting_BONUS_TYPE_INVALID)
	}
}

func psf_on_00004_RefundBonusBullet(action *flux.Action, am *AccountingManager) {
	if am.guestEnable {
		return
	}

	secWebSocketKey := action.Payload()[0].(string)
	b := action.Payload()[1].(*bullet.Bullet)

	hitFish := b.ExtraData[0].(*fish.Fish)
	hitResult := b.ExtraData[2].(*probability.Probability)
	accountingSn := action.Payload()[3].(uint64)

	if wallet.Service.WalletCategory(secWebSocketKey) != wallet.CategoryHost {
		accountingSn = am.accountingSn
	}

	remainBullets := b.Bullets - b.UsedBullets

	gameResult := fmt.Sprintf("%d/%d/%d/%s/%d/%d",
		b.TypeId,
		b.Bullets,
		remainBullets,
		"-",
		hitFish.Scene,
		hitResult.Multiplier,
	)

	bet := b.Bet * b.Rate

	go receipt.Service.New(
		uuid.New().String(),
		secWebSocketKey,
		am.hostExtId,
		gameResult,
		"",
		accountingSn,
		0,
		bet*uint64(remainBullets),
	)

	logger.Service.Zap.Infow("Refund Bonus Bullet Accounting",
		"GameUser", am.secWebSocketKey,
		"AccountingSn", accountingSn,
	)
}

func psf_on_00004_CheckBulletTypeId(bulletId, restBullets int32,
	am *AccountingManager,
	b *bullet.Bullet,
	machineGunUuid, superMachineGunUuid map[string]map[uint64]int32,
) {
	switch bulletId {
	case 201:
		if machineGunUuid[b.BonusUuid][am.accountingSn] > 0 {
			// 2. 核銷機關槍 cff118d7
			machineGunUuid[b.BonusUuid][am.accountingSn] -= restBullets
		}

		if machineGunUuid[b.BonusUuid][am.accountingSn] <= 0 {
			delete(machineGunUuid, b.BonusUuid)
		}

	case 202:
		if superMachineGunUuid[b.BonusUuid][am.accountingSn] > 0 {
			superMachineGunUuid[b.BonusUuid][am.accountingSn] -= restBullets
		}

		if superMachineGunUuid[b.BonusUuid][am.accountingSn] <= 0 {
			delete(superMachineGunUuid, b.BonusUuid)
		}
	}
}

func psf_on_00004_BulletHit(bulletId int32,
	bulletUuid map[string]uint64,
	am *AccountingManager,
	hitBullet *bullet.Bullet,
	tempReceipts map[uint64]*tempReceipt,
	machineGunUuid, superMachineGunUuid map[string]map[uint64]int32,
	bonusBulletShooted map[string]map[string]map[uint64]int32,
	isBonusBulletShooted bool,
) {
	switch bulletId {
	case 0:
		// cause by async reason, see 1dd8ed7e
		if _, ok := bulletUuid[hitBullet.Uuid]; !ok {
			// 1. 紀錄子彈 35adc39d
			bulletUuid[hitBullet.Uuid] = am.accountingSn
			tempReceipts[am.accountingSn].betCent += hitBullet.Bet * hitBullet.Rate
			am.rounds++
		} else {
			// 2. 核銷子彈 35adc39d
			delete(bulletUuid, hitBullet.Uuid)
		}

	case 201:
		if !isBonusBulletShooted {
			if machineGunUuid[hitBullet.BonusUuid][am.accountingSn] > 0 {
				// 2. 核銷機關槍 cff118d7
				machineGunUuid[hitBullet.BonusUuid][am.accountingSn]--
			}

			if machineGunUuid[hitBullet.BonusUuid][am.accountingSn] <= 0 {
				delete(machineGunUuid, hitBullet.BonusUuid)
			}
		} else {
			logger.Service.Zap.Infow("BonusBulletShooted But Can't find fish",
				"BulletType", "MachineGun",
				"BonusUuid", hitBullet.BonusUuid,
				"AccountingSn", am.accountingSn,
				"Shooted", bonusBulletShooted["machineGun"][hitBullet.BonusUuid][am.accountingSn],
			)
		}

		if bonusBulletShooted["machineGun"][hitBullet.BonusUuid][am.accountingSn] > 0 {
			bonusBulletShooted["machineGun"][hitBullet.BonusUuid][am.accountingSn]--
		}

		if bonusBulletShooted["machineGun"][hitBullet.BonusUuid][am.accountingSn] <= 0 {
			delete(bonusBulletShooted["machineGun"], hitBullet.BonusUuid)
		}

	case 202:
		if !isBonusBulletShooted {
			if superMachineGunUuid[hitBullet.BonusUuid][am.accountingSn] > 0 {
				superMachineGunUuid[hitBullet.BonusUuid][am.accountingSn]--
			}

			if superMachineGunUuid[hitBullet.BonusUuid][am.accountingSn] <= 0 {
				delete(superMachineGunUuid, hitBullet.BonusUuid)
			}
		} else {
			logger.Service.Zap.Infow("BonusBulletShooted But Can't find fish",
				"BulletType", "SuperMachineGun",
				"BonusUuid", hitBullet.BonusUuid,
				"AccountingSn", am.accountingSn,
				"Shooted", bonusBulletShooted["superMachineGun"][hitBullet.BonusUuid][am.accountingSn],
			)
		}

		if bonusBulletShooted["superMachineGun"][hitBullet.BonusUuid][am.accountingSn] > 0 {
			bonusBulletShooted["superMachineGun"][hitBullet.BonusUuid][am.accountingSn]--
		}

		if bonusBulletShooted["superMachineGun"][hitBullet.BonusUuid][am.accountingSn] <= 0 {
			delete(bonusBulletShooted["superMachineGun"], hitBullet.BonusUuid)
		}
	}
}

func psf_on_00004_BulletHit_TriggerId(triggerIconId int,
	am *AccountingManager,
	hitBullet *bullet.Bullet,
	hitResult *probability.Probability,
	tempReceipts map[uint64]*tempReceipt,
	machineGunUuid, superMachineGunUuid map[string]map[uint64]int32,
	bonusBulletShooted map[string]map[string]map[uint64]int32,
	redEnvelopeUuid, slotUuid map[string]uint64,
) bool {
	switch triggerIconId {
	case 100, 102:
		// 1. 紀錄紅包 858cedf5
		redEnvelopeUuid[hitResult.BonusUuid] = am.accountingSn
		return true

	case 101:
		// 1. 紀錄老虎機 6e01e01a
		slotUuid[hitResult.BonusUuid] = am.accountingSn
		return true

	case 201:
		machineGun := make(map[uint64]int32)

		machineGun[am.accountingSn] = int32(hitResult.BonusPayload.(int))

		machineGunUuid[hitResult.BonusUuid] = machineGun

		tempMachineGun := make(map[uint64]int32)
		tempMachineGun[am.accountingSn] = int32(hitResult.BonusPayload.(int))
		tempMachineGunUuid := make(map[string]map[uint64]int32)
		tempMachineGunUuid[hitResult.BonusUuid] = tempMachineGun
		bonusBulletShooted["machineGun"] = tempMachineGunUuid

		tempReceipts[am.accountingSn].bonusWinCent += uint64(hitResult.Pay) * hitBullet.Bet * hitBullet.Rate

		return true

	case 202:
		superMachineGun := make(map[uint64]int32)

		superMachineGun[am.accountingSn] = int32(hitResult.BonusPayload.(int))

		superMachineGunUuid[hitResult.BonusUuid] = superMachineGun

		tempSuperMachineGun := make(map[uint64]int32)
		tempSuperMachineGun[am.accountingSn] = int32(hitResult.BonusPayload.(int))
		tempSuperMachineGunUuid := make(map[string]map[uint64]int32)
		tempSuperMachineGunUuid[hitResult.BonusUuid] = tempSuperMachineGun
		bonusBulletShooted["superMachineGun"] = tempSuperMachineGunUuid

		tempReceipts[am.accountingSn].bonusWinCent += uint64(hitResult.Pay) * hitBullet.Bet * hitBullet.Rate

		return true
	}
	return false
}

func psf_on_00004_BulletBonus(bonusBulletId int32, bonusBullet *bullet.Bullet, machineGunUuid, superMachineGunUuid map[string]map[uint64]int32, bonusBulletShooted map[string]map[string]map[uint64]int32) {
	switch bonusBulletId {
	case 201:
		// 2. 核銷機關槍 cff118d7

		// if free bullets reached 999 then got 0 in the new free bullets
		if bonusBullet.Bullets == 0 {
			delete(machineGunUuid, bonusBullet.BonusUuid)
			delete(bonusBulletShooted["machineGun"], bonusBullet.BonusUuid)
		}

	case 202:
		if bonusBullet.Bullets == 0 {
			delete(superMachineGunUuid, bonusBullet.BonusUuid)
			delete(bonusBulletShooted["superMachineGun"], bonusBullet.BonusUuid)
		}
	}
}
