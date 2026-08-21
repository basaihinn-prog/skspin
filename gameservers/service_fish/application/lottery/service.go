package lottery

import (
	"fmt"
	"serve/service_fish/application/accountingmanager"
	"serve/service_fish/application/activity"
	"serve/service_fish/application/lobbysetting"
	"serve/service_fish/application/rtp"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	"serve/service_fish/domain/probability"
	"serve/service_fish/models"
	"sync"
	"time"

	"serve/fish_comm/common"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/fish_comm/wallet"
)

const (
	ActionLotteryStart      = "ActionLotteryStart"
	ActionLotteryStop       = "ActionLotteryStop"
	ActionLotteryDelete     = "ActionLotteryDelete"
	ActionUnprocessedBullet = "ActionUnprocessedBullet"
)

var Service = &service{
	Id:      "LotteryService",
	in:      make(chan *flux.Action, common.Service.ChanSize),
	rwMutex: &sync.RWMutex{},
}

type service struct {
	Id      string
	in      chan *flux.Action
	rwMutex *sync.RWMutex
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
	lotteries := make(map[string]*lottery)

	for action := range s.in {
		s.handleAction(action, lotteries)
	}
}

func (s *service) handleAction(action *flux.Action, lotteries map[string]*lottery) {
	s.rwMutex.RLock()
	logger.Service.Zap.Debugw("All Lotteries",
		"Size", len(lotteries),
		"Map", lotteries,
	)
	s.rwMutex.RUnlock()

	switch action.Key().Name() {
	case ActionLotteryStart:
		secWebSocketKey := action.Key().From()
		gameId := action.Payload()[0].(string)
		mathModuleId := action.Payload()[1].(string)
		kind := action.Payload()[2].(string)

		var h IHandler = nil

		if v, ok := action.Payload()[3].(*accountingmanager.AccountingManager); ok && v != nil {
			h = v
		}

		// See lottery's id()
		id := s.HashId(secWebSocketKey)

		if _, ok := lotteries[id]; !ok {
			s.rwMutex.Lock()
			lotteries[id] = newLottery(secWebSocketKey, gameId, mathModuleId, kind, h)
			s.rwMutex.Unlock()
		}

	case ActionLotteryStop:
		secWebSocketKey := action.Payload()[0].(string)

		//See lottery's id()
		id := s.HashId(secWebSocketKey)

		if v, ok := lotteries[id]; ok {
			v.destroy()

			s.rwMutex.Lock()
			delete(lotteries, id)
			s.rwMutex.Unlock()
		}

	case ActionLotteryDelete:
		secWebSocketKey := action.Key().From()

		for k, v := range lotteries {
			if v.secWebSocketKey == secWebSocketKey {
				v.destroy()

				s.rwMutex.Lock()
				delete(lotteries, k)
				s.rwMutex.Unlock()
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

func (s *service) MathProcess(
	secWebSocketKey, gameId, mathModuleId, bulletId, bonusBulletId, fishRtpId, bulletRtpId string,
	mercenaryType, fishTypeId, bulletTypeId int32,
	bet, rate uint64, subgameId int,
	hitBullet *bullet.Bullet, hitFish *fish.Fish,
	playerKind, lotteryId string,
	lenFishTemp, lenBulletTemp int,
	isSimulation bool, blacklistValue uint32, depositMultiple uint64) (bool, *probability.Probability, int, string, int, uint64, int, int, uint64) {

	am := accountingmanager.Service.Get(secWebSocketKey)
	if am != nil {
		am.RotatiomRwMutex.RLock()
		defer am.RotatiomRwMutex.RUnlock()
	}

	rtpState := -1
	rtpId := ""
	netWinGroup := -1
	var rtpBullets uint64 = 0

	sn := accountingmanager.Service.GetAccountingSn(secWebSocketKey)

	denominator := 0
	molecular := 0

	if !isSimulation && hitBullet != nil && hitFish != nil {
		rate = hitBullet.Rate
		fishTypeId = hitFish.TypeId
		bulletTypeId = hitBullet.TypeId
		bulletId = hitBullet.Uuid
		bonusBulletId = hitBullet.BonusUuid
		fishRtpId = hitFish.RtpId
	}

	// Regular Bullet
	if bonusBulletId == "" && mercenaryType == -1 {
		if !isSimulation {
			switch wallet.Service.WalletCategory(secWebSocketKey) {
			case wallet.CategoryHost:
				// FishWallet Cent is 0 then Deposit Bet x 20
				if cent, err := wallet.Service.CentCache(secWebSocketKey); err == nil {
					// Check hitBullet is not nil
					if hitBullet != nil {
						checkAccountingEnd := true
						// Process WalletHost Use MultiWeapon
						if gameId == models.PSF_ON_00003 && hitBullet.Bet > 1 && cent >= 1*rate {
							checkAccountingEnd = false
						}

						if cent <= 0 || cent < hitBullet.Bet*rate {
							accountingSn := accountingmanager.Service.GetAccountingSn(secWebSocketKey)
							//am := accountingmanager.Service.Get(secWebSocketKey)
							betAmt, winAmt := am.GetWalletHostAmt(secWebSocketKey)
							//todo top one
							//wallet.Service.EndAccounting(secWebSocketKey, hitBullet.Bet*rate, accountingSn, checkAccountingEnd)
							wallet.Service.EndAccounting(secWebSocketKey, hitBullet.Bet*rate, checkAccountingEnd)
							time.Sleep(500 * time.Millisecond)
							if cent > 0 {
								wallet.Service.Withdraw(secWebSocketKey, accountingSn, cent, betAmt, winAmt)
							}

							balance, err := wallet.Service.BalanceCache(secWebSocketKey)
							if err != nil {
								logger.Service.Zap.Errorw("Shared Wallet Get Balance Error",
									"secWebSocketKey", secWebSocketKey,
								)
							}

							depositCent := hitBullet.Bet * rate * depositMultiple
							if balance < depositCent {
								depositCent = balance
							}

							wallet.Service.Deposit(secWebSocketKey, depositCent, accountingSn)
						}
					}

					if !wallet.Service.DecreaseCent(secWebSocketKey, bulletId, bet*rate) {
						return false, nil, rtpState, rtpId, netWinGroup, rtpBullets, denominator, molecular, sn
					}
				} else {
					logger.Service.Zap.Errorw("Shared Wallet Get Cent Error",
						"secWebSocketKey", secWebSocketKey,
					)
				}

			default:
				if !wallet.Service.DecreaseBalance(secWebSocketKey, bulletId, bet*rate) {
					return false, nil, rtpState, rtpId, netWinGroup, rtpBullets, denominator, molecular, sn
				}
			}
		}

		switch gameId {
		case models.PSF_ON_00004:
			minBet, minRate := lobbysetting.Service.MinBetAndRate(secWebSocketKey, gameId)
			denominator, molecular = rtp.Service.DecreaseDenominator(gameId, subgameId, mathModuleId, secWebSocketKey,
				fishTypeId,
				fishRtpId,
				bet*rate,
				minBet*minRate,
			)

			if !isSimulation {
				rtp.Service.Save(secWebSocketKey, gameId, subgameId)
			}

		default:
			rtp.Service.Decrease(gameId, subgameId, mathModuleId, secWebSocketKey, bet*rate)
			if !isSimulation {
				rtp.Service.Save(secWebSocketKey, gameId, subgameId)
			}
		}

		hitBullet.RtpId = rtp.Service.RtpId(secWebSocketKey, gameId, subgameId)
		bulletRtpId = hitBullet.RtpId
	}

	// Process MachineGun and SuperMachineGun have UseRTP Value
	hitBullet.RtpId = probability.Service.UseRtp(gameId, mathModuleId, bulletRtpId, bulletTypeId, mercenaryType)

	if !isSimulation && hitBullet != nil && hitFish != nil {
		bulletRtpId = hitBullet.RtpId
	}

	hitResult := probability.Service.Calc(
		gameId,
		mathModuleId,
		bulletRtpId,
		fishTypeId,
		bulletTypeId,
		rtp.Service.RtpBudget(secWebSocketKey, gameId, subgameId, mathModuleId),
		bet*rate,
	)

	switch gameId {
	case models.PSF_ON_00003:
		rtp.Service.IncreaseBulletCollection(secWebSocketKey, gameId, subgameId, uint64(hitResult.Bullet))
		if !isSimulation {
			rtp.Service.Save(secWebSocketKey, gameId, subgameId)
		}

	case models.PSF_ON_00004:
		rtp.Service.DecreaseMolecular(gameId, subgameId, mathModuleId, secWebSocketKey,
			fishTypeId,
			int(uint64(hitResult.Pay*hitResult.Multiplier)),
			int(bet*rate),
			hitResult.TriggerIconId,
		)
	}

	// BlackList Check
	if probability.Service.Random(1000) < uint64(blacklistValue) {
		hitResult = probability.Service.GetNewSimplePay(int(fishTypeId))
	}

	if !isSimulation {
		switch wallet.Service.WalletCategory(secWebSocketKey) {
		case wallet.CategoryHost:
			/*accountingSn := accountingmanager.Service.GetAccountingSn(secWebSocketKey)
			if accountingSn == 0 {
				balance, _ := wallet.Service.BalanceCache(secWebSocketKey)

				depositCent := bet * rate * depositMultiple
				if balance < depositCent {
					depositCent = balance
				}

				wallet.Service.Deposit(secWebSocketKey, depositCent, accountingSn)
			}*/

			wallet.Service.IncreaseCent(
				secWebSocketKey,
				bulletId,
				uint64(hitResult.Pay)*bet*rate*uint64(hitResult.Multiplier),
			)

		default:
			wallet.Service.IncreaseBalance(
				secWebSocketKey,
				bulletId,
				uint64(hitResult.Pay)*bet*rate*uint64(hitResult.Multiplier),
				0,
				0,
			)
		}

		if playerKind == PLAYER {
			activity.Service.Record(secWebSocketKey, hitFish, hitBullet, hitResult)
		}

		balance, err := wallet.Service.BalanceCache(secWebSocketKey)

		if err != nil {
			// TODO JOHNNY game user is already disconnect
			logger.Service.Zap.Errorw("BUG",
				"GameUser", secWebSocketKey,
				"HitBulletTemp", lenBulletTemp,
				"HitFishTemp", lenFishTemp,
				"Error", err,
			)
		}
		hitBullet.PlayerCent = balance

	} else {
		// Simulation Use
		rtpState = rtp.Service.RtpState(secWebSocketKey, gameId, subgameId)
		rtpId = rtp.Service.RtpId(secWebSocketKey, gameId, subgameId)
		netWinGroup = rtp.Service.NetWinGroup(secWebSocketKey, gameId, subgameId)
		rtpBullets = rtp.Service.RtpBudget(secWebSocketKey, gameId, subgameId, mathModuleId)
	}

	return true, hitResult, rtpState, rtpId, netWinGroup, rtpBullets, denominator, molecular, sn
}

func (s *service) WriteDebugLog(
	hitBullet *bullet.Bullet,
	hitFish *fish.Fish, hitFishies *fish.Fishies,
	hitResult *probability.Probability, hitResults []*probability.Probability,
	secWebSocketKey, gameId string, subGameId int, mercenaryType int32,
) {
	if hitFish != nil && hitResult != nil {
		var weaponType = -1
		switch {
		case hitBullet.BonusUuid == "" && hitBullet.TypeId == 999 && mercenaryType == 0:
			weaponType = 4
		case hitBullet.TypeId == 200:
			weaponType = 5
		default:
			weaponType = 0
		}

		fishIsDied := false
		totalPay := hitResult.Pay * hitResult.Multiplier
		switch hitFish.GameId {
		case models.PSF_ON_00003:
			if (hitFish.TypeId == 303 && (hitResult.Pay > 0 || hitResult.Bullet > 0)) ||
				(hitFish.TypeId != 303 && hitResult.Pay > 0) ||
				hitResult.TriggerIconId > 0 {
				fishIsDied = true
			}
			// Get Slot Pay
			if hitResult.Pay > 0 {
				switch hitFish.TypeId {
				case 6, 7, 8, 9, 10, 11:
					totalPay += hitResult.ExtraData[0].(int)
				}
			}

		}

		logger.Service.Zap.Debugw("Probability Verification",
			"GameId", hitFish.GameId,
			"FishId", hitFish.TypeId,
			"WeaponType", weaponType,
			"RtpId", hitBullet.RtpId,
			"FishIsDied", fishIsDied,
			"RtpBudget", rtp.Service.RtpBudget(secWebSocketKey, hitFish.GameId, subGameId, hitFish.MathModuleId),
			"FreeBullet", hitResult.Bullet,
			"Pay", totalPay,
		)
	}

	if hitFishies != nil && hitResults != nil {
		var weaponType = -1
		var fishIdList []int32
		var fishIsDied []bool
		var freeBulletList []int
		var payList []uint64

		for i := 0; i < len(hitFishies.Fishies); i++ {
			fishIdList = append(fishIdList, hitFishies.Fishies[i].TypeId)
			if hitResults[i] != nil {
				freeBulletList = append(freeBulletList, hitResults[i].Bullet)
			} else {
				freeBulletList = append(freeBulletList, 0)
			}

			switch gameId {
			case models.PSF_ON_00003:
				if hitResults[i] != nil {
					tempPay := hitResults[i].Pay * hitResults[i].Multiplier
					if (hitFishies.Fishies[i].TypeId == 303 && (hitResults[i].Pay > 0 || hitResults[i].Bullet > 0)) ||
						(hitFishies.Fishies[i].TypeId != 303 && hitResults[i].Pay > 0) ||
						hitResults[i].TriggerIconId > 0 {
						fishIsDied = append(fishIsDied, true)
					} else {
						fishIsDied = append(fishIsDied, false)
					}

					// Get Slot Pay
					if hitResults[i].Pay > 0 {
						switch hitFishies.Fishies[i].TypeId {
						case 6, 7, 8, 9, 10, 11:
							tempPay += hitResults[i].ExtraData[0].(int)
						}
					}
					payList = append(payList, uint64(tempPay))

				} else {
					fishIsDied = append(fishIsDied, false)
					payList = append(payList, 0)
				}

			}
		}
		switch {
		case len(hitFishies.Fishies) == 3 && mercenaryType == -1:
			weaponType = 1
		case len(hitFishies.Fishies) == 5 && mercenaryType == -1:
			weaponType = 2
		case len(hitFishies.Fishies) == 10 && mercenaryType == -1:
			weaponType = 3
		case (len(hitFishies.Fishies) == 3 && mercenaryType == 1) ||
			len(hitFishies.Fishies) == 5 && mercenaryType == 2:
			weaponType = 4
		}

		logger.Service.Zap.Debugw("Probability Verification",
			"GameId", gameId,
			"FishId", fishIdList,
			"WeaponType", weaponType,
			"RtpId", hitBullet.RtpId,
			"FishIsDied", fishIsDied,
			"RtpBudget", rtp.Service.RtpBudget(secWebSocketKey, hitFishies.Fishies[0].GameId, subGameId, hitFishies.Fishies[0].MathModuleId),
			"FreeBullet", freeBulletList,
			"Pay", payList,
		)
	}
}
