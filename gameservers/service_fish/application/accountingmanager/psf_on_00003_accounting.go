package accountingmanager

import (
	"fmt"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	"serve/service_fish/domain/receipt"
	"serve/service_fish/domain/slot"
	"strconv"

	"serve/fish_comm/accounting"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/fish_comm/wallet"

	"github.com/google/uuid"
)

const (
	SlotTypeId            = 101
	DrillTypeId           = 200
	NormalBulletTypeId    = 0
	MercenaryBulletTypeId = 999
)

func psf_on_00003_HitBullet(action *flux.Action, am *AccountingManager) {
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

	// Process Drill don't accounting normal
	if hitResult.TriggerIconId == 200 {
		logger.Service.Zap.Infow("Trigger Bullet Accounting",
			"GameUser", am.secWebSocketKey,
			"AccountingSn", accountingSn,
		)
		return
	}

	bulletUuid := uuid.New().String()

	// Check Drill, Drill had only one UUID
	if !CheckBonusBullet(hitBullet.TypeId) {
		bulletUuid = hitBullet.Uuid
	}

	pay := uint64(hitResult.Pay)
	bet := hitBullet.Bet * hitBullet.Rate
	gameData := ""

	if CheckBonusBullet(hitBullet.TypeId) {
		gameData = strconv.Itoa(int(bet))
		bet = 0
	}

	switch hitFish.TypeId {
	case 6, 7, 8, 9, 10, 11:
		if len(hitResult.ExtraData) == 2 {
			icons := hitResult.ExtraData[1].([]int32)
			pays := hitResult.BonusPayload.([]int64)
			totalPay := hitResult.ExtraData[0].(int)

			gameData = fmt.Sprintf("%d_%d/%d_%d/%d_%d",
				pays[0], icons[1],
				pays[1], icons[4],
				pays[2], icons[7],
			)

			pay += uint64(totalPay)
		}
	}

	gameResult := fmt.Sprintf("%d/%d/%d/%d/%d/%d/%d/%d",
		hitBullet.TypeId,
		hitBullet.Bullets,
		hitBullet.UsedBullets,
		hitFish.TypeId,
		hitFish.Scene,
		pay*uint64(hitResult.Multiplier),
		pay*hitBullet.Bet*hitBullet.Rate*uint64(hitResult.Multiplier),
		hitResult.Bullet,
	)

	go receipt.Service.New(
		bulletUuid,
		secWebSocketKey,
		am.hostExtId,
		gameResult,
		gameData,
		accountingSn,
		bet,
		pay*hitBullet.Bet*hitBullet.Rate*uint64(hitResult.Multiplier),
	)

	logger.Service.Zap.Infow("Bullet Accounting",
		"GameUser", am.secWebSocketKey,
		"AccountingSn", accountingSn,
	)
}

func psf_on_00003_MultiHitBullet(action *flux.Action, am *AccountingManager) {
	if am.guestEnable {
		return
	}

	secWebSocketKey := action.Payload()[0].(string)
	hitFishies := action.Payload()[1].(*fish.Fishies)
	hitBullet := action.Payload()[2].(*bullet.Bullet)
	hitResults := action.Payload()[3].([]*probability.Probability)
	accountingSn := action.Payload()[5].(uint64)

	if wallet.Service.WalletCategory(secWebSocketKey) != wallet.CategoryHost {
		accountingSn = am.accountingSn
	}

	// Check Drill, Drill had only one UUID
	if CheckBonusBullet(hitBullet.TypeId) {
		hitBullet.Uuid = uuid.New().String()
	}

	gameResult := ""
	gameData := ""
	var totalPay uint64 = 0
	bet := 1 * hitBullet.Rate

	if CheckBonusBullet(hitBullet.TypeId) {
		gameData = strconv.Itoa(int(bet))
		bet = 0
	}

	for i := 0; i < len(hitResults); i++ {
		var pay uint64 = 0
		bullet := 0
		if hitResults[i] != nil {
			pay = uint64(hitResults[i].Pay * hitResults[i].Multiplier)
			bullet = hitResults[i].Bullet

			// Process GameData
			switch hitFishies.Fishies[i].TypeId {
			case 6, 7, 8, 9, 10, 11:
				if len(hitResults[i].ExtraData) == 2 {
					icons := hitResults[i].ExtraData[1].([]int32)
					pays := hitResults[i].BonusPayload.([]int64)
					tempPay := hitResults[i].ExtraData[0].(int)

					tempGameData := fmt.Sprintf("%d_%d/%d_%d/%d_%d",
						pays[0], icons[1],
						pays[1], icons[4],
						pays[2], icons[7],
					)

					if gameData != "" {
						gameData = gameData + "|" + tempGameData
					} else {
						gameData = tempGameData
					}

					totalPay += uint64(tempPay) * bet
				}
			}
		}

		tempResult := fmt.Sprintf("%d/%d/%d/%d/%d/%d/%d/%d",
			hitBullet.TypeId,
			hitBullet.Bullets,
			hitBullet.UsedBullets,
			hitFishies.Fishies[i].TypeId,
			hitFishies.Fishies[i].Scene,
			pay,
			pay*1*hitBullet.Rate,
			bullet,
		)

		if gameResult != "" {
			gameResult = gameResult + "|" + tempResult
		} else {
			gameResult = tempResult
		}

		totalPay = totalPay + (pay * bet)
	}

	go receipt.Service.New(
		hitBullet.Uuid,
		secWebSocketKey,
		am.hostExtId,
		gameResult,
		gameData,
		accountingSn,
		bet,
		totalPay,
	)

	logger.Service.Zap.Infow("Bullet Accounting",
		"GameUser", am.secWebSocketKey,
		"AccountingSn", accountingSn,
	)
}

func psf_on_00003_Bonus(action *flux.Action, am *AccountingManager) {
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

		gameResult := fmt.Sprintf("%d/%d/%d/%d/%d/%d/%d/%d",
			hitBullet.TypeId,
			hitBullet.Bullets,
			hitBullet.UsedBullets,
			hitFish.TypeId,
			hitFish.Scene,
			hitResult.Pay*hitResult.Multiplier,
			uint64(hitResult.Pay)*b.Bet*b.Rate,
			hitResult.Bullet,
		)

		bet := b.Bet * b.Rate

		if CheckBonusBullet(hitBullet.TypeId) {
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

		if CheckBonusBullet(hitBullet.TypeId) {
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

func psf_on_00003_RefundBonusBullet(action *flux.Action, am *AccountingManager) {
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

func psf_on_00003_CheckBulletTypeId(bulletId, restBullets int32,
	am *AccountingManager,
	b *bullet.Bullet,
	drillUuid map[string]map[uint64]int32,
) {
	switch bulletId {
	//case DrillTypeId:
	case 12, 13, 14, 15, 16, 17, 200:
		if drillUuid[b.BonusUuid][am.accountingSn] > 0 {
			// 2. 核銷鑽頭炮 50dd6879
			drillUuid[b.BonusUuid][am.accountingSn] -= restBullets
		}

		if drillUuid[b.BonusUuid][am.accountingSn] <= 0 {
			delete(drillUuid, b.BonusUuid)
		}
	}
}

func psf_on_00003_BulletHit(bulletId int32,
	bulletUuid map[string]uint64,
	am *AccountingManager,
	hitBullet *bullet.Bullet,
	tempReceipts map[uint64]*tempReceipt,
	drillUuid map[string]map[uint64]int32,
	bonusBulletShooted map[string]map[string]map[uint64]int32,
	oneBet uint64,
	isBonusBulletShooted bool,
) {
	switch bulletId {
	case NormalBulletTypeId:
		// cause by async reason, see 1dd8ed7e
		if _, ok := bulletUuid[hitBullet.Uuid]; !ok {
			// 1. 紀錄子彈 35adc39d
			bulletUuid[hitBullet.Uuid] = am.accountingSn
			//tempReceipts[am.accountingSn].betCent += hitBullet.Bet * hitBullet.Rate
			tempReceipts[am.accountingSn].betCent += oneBet * hitBullet.Rate
			am.rounds++
		} else {
			// 2. 核銷子彈 35adc39d
			delete(bulletUuid, hitBullet.Uuid)
		}

	//case DrillTypeId:
	case 12, 13, 14, 15, 16, 17, 200:
		if !isBonusBulletShooted {
			if drillUuid[hitBullet.BonusUuid][am.accountingSn] > 0 {
				// 2. 核銷鑽頭炮 50dd6879
				drillUuid[hitBullet.BonusUuid][am.accountingSn]--
			}

			if drillUuid[hitBullet.BonusUuid][am.accountingSn] <= 0 {
				delete(drillUuid, hitBullet.BonusUuid)
			}
		} else {
			logger.Service.Zap.Infow("BonusBulletShooted But Can't find fish",
				"BulletType", "Drill",
				"BonusUuid", hitBullet.BonusUuid,
				"AccountingSn", am.accountingSn,
				"Shooted", bonusBulletShooted["drill"][hitBullet.BonusUuid][am.accountingSn],
			)
		}

		if bonusBulletShooted["drill"][hitBullet.BonusUuid][am.accountingSn] > 0 {
			bonusBulletShooted["drill"][hitBullet.BonusUuid][am.accountingSn]--
		}

		if bonusBulletShooted["drill"][hitBullet.BonusUuid][am.accountingSn] <= 0 {
			delete(bonusBulletShooted["drill"], hitBullet.BonusUuid)
		}

	//傭兵子彈
	case MercenaryBulletTypeId:
		//do nothing

	}
}

func psf_on_00003_BulletHit_TriggerId(triggerIconId int,
	am *AccountingManager,
	hitBullet *bullet.Bullet,
	hitResult *probability.Probability,
	tempReceipts map[uint64]*tempReceipt,
	drillUuid map[string]map[uint64]int32,
	bonusBulletShooted map[string]map[string]map[uint64]int32,
	slotUuid map[string]uint64,
) bool {
	switch triggerIconId {
	case SlotTypeId:
		// 1. 紀錄老虎機 6e01e01a
		slotUuid[hitResult.BonusUuid] = am.accountingSn

		return true

	//case DrillTypeId:
	case 12, 13, 14, 15, 16, 17, 200:
		// 1. 紀錄鑽頭炮 50dd6879
		drill := make(map[uint64]int32)

		drill[am.accountingSn] = int32(hitResult.BonusPayload.(int))

		drillUuid[hitResult.BonusUuid] = drill

		tempDrill := make(map[uint64]int32)
		tempDrill[am.accountingSn] = int32(hitResult.BonusPayload.(int))
		tempDrillUuid := make(map[string]map[uint64]int32)
		tempDrillUuid[hitResult.BonusUuid] = tempDrill
		bonusBulletShooted["drill"] = tempDrillUuid

		tempReceipts[am.accountingSn].bonusWinCent += uint64(hitResult.Pay) * hitBullet.Bet * hitBullet.Rate

		return true
	}
	return false
}

func psf_on_00003_BulletMultiHit_TriggerId(am *AccountingManager,
	hitBullet *bullet.Bullet,
	hitResults []*probability.Probability,
	tempReceipts map[uint64]*tempReceipt,
	drillUuid map[string]map[uint64]int32,
	bonusBulletShooted map[string]map[string]map[uint64]int32,
	slotUuid map[string]uint64,
) (isPause bool) {
	for i := 0; i < len(hitResults); i++ {
		if hitResults[i] != nil {
			switch hitResults[i].TriggerIconId {
			case SlotTypeId:
				// 1. 紀錄老虎機 6e01e01a
				slotUuid[hitResults[i].BonusUuid] = am.accountingSn

				isPause = true

			//case DrillTypeId:
			case 12, 13, 14, 15, 16, 17, 200:
				// 1. 紀錄鑽頭炮 50dd6879
				drill := make(map[uint64]int32)

				drill[am.accountingSn] = int32(hitResults[i].BonusPayload.(int))

				drillUuid[hitResults[i].BonusUuid] = drill

				tempDrill := make(map[uint64]int32)
				tempDrill[am.accountingSn] = int32(hitResults[i].BonusPayload.(int))
				tempDrillUuid := make(map[string]map[uint64]int32)
				tempDrillUuid[hitResults[i].BonusUuid] = tempDrill
				bonusBulletShooted["drill"] = tempDrillUuid

				tempReceipts[am.accountingSn].bonusWinCent += uint64(hitResults[i].Pay) * hitBullet.Bet * hitBullet.Rate

				isPause = true
			}
		}
	}
	return
}

func psf_on_00003_BulletBonus(bonusBulletId int32) {
	switch bonusBulletId {
	case DrillTypeId:
		// do nothing now
	}
}

func CheckBonusBullet(a int32) bool {
	//bulletTypeList := []int32 {DrillTypeId}
	bulletTypeList := []int32{12, 13, 14, 15, 16, 17, 200}

	for _, b := range bulletTypeList {
		if b == a {
			return true
		}
	}
	return false
}
