package lottery

import (
	"serve/service_fish/application/accountingmanager"
	"serve/service_fish/application/activity"
	"serve/service_fish/application/gamesetting"
	"serve/service_fish/application/rtp"
	"serve/service_fish/domain/auth"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/fish"
	game_recovery_proto "serve/service_fish/domain/gamerecovery/proto"
	"time"

	"serve/fish_comm/mediator"

	"serve/fish_comm/broadcaster"

	//game_recovery_proto "serve/service_fish/domain/gamerecovery/proto"
	"fmt"
	"os"
	"serve/service_fish/domain/probability"
	probability_proto "serve/service_fish/domain/probability/proto"
	"serve/service_fish/models"
	common_fish_proto "serve/service_fish/models/proto"
	"strconv"
	"sync"

	"serve/fish_comm/common"
	common_proto "serve/fish_comm/common/proto"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/fish_comm/wallet"

	"github.com/gogo/protobuf/proto"
)

const (
	BOT    = "BOT"
	PLAYER = "PLAYER"
	GUEST  = "GUEST"
)

type lottery struct {
	secWebSocketKey string
	gameId          string
	mathModuleId    string // not used yet
	poolSize        string
	kind            string
	in              chan *flux.Action
	rwMutex         *sync.RWMutex
	iHandler        IHandler
}

type IHandler interface {
	BulletHitHandler(secWebSocketKey string, hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability, accountingSn uint64)
	BulletBonusHandler(secWebSocketKey string, bonusBullet *bullet.Bullet, accountingSn uint64)
	BulletMultiHitHandler(secWebSocketKey string, hitFishies *fish.Fishies, hitBullet *bullet.Bullet, hitResults []*probability.Probability, oneBet, accountingSn uint64)
	BonusBulletShootedHandler(secWebSocketKey string, bonusBullet *bullet.Bullet)
}

func newLottery(secWebSocketKey, gameId, mathModuleId, kind string, h IHandler) *lottery {
	l := &lottery{
		secWebSocketKey: secWebSocketKey,
		gameId:          gameId,
		mathModuleId:    mathModuleId,
		poolSize:        os.Getenv("POOL_SIZE"),
		kind:            kind,
		in:              make(chan *flux.Action, common.Service.ChanSize),
		rwMutex:         &sync.RWMutex{},
		iHandler:        h,
	}
	l.run()
	logger.Service.Zap.Infow("Lottery Running",
		"GameUser", secWebSocketKey,
		"Chan", fmt.Sprintf("%p", l.in),
		"GameId", l.gameId,
		"MathModuleId", l.mathModuleId,
		"Kind", l.kind,
	)
	return l
}

func (l *lottery) run() {
	flux.Register(l.id(), l.in)
	hitFishTemp := make(map[*fish.Fish]*bullet.Bullet)
	hitBulletTemp := make(map[*bullet.Bullet]*fish.Fish)

	multiHitFishTemp := make(map[*fish.Fishies]*bullet.Bullet, 0)
	multiHitBulletTemp := make(map[*bullet.Bullet]*fish.Fishies, 0)

	poolSize := 10

	if l.poolSize != "" {
		if poolSize, _ = strconv.Atoi(l.poolSize); poolSize == 0 {
			poolSize = 10
		}
	}

	for i := 0; i < poolSize; i++ {
		go func() {
			for action := range l.in {
				l.handleAction(action, hitFishTemp, hitBulletTemp, multiHitFishTemp, multiHitBulletTemp)
			}
		}()
	}
}

func (l *lottery) handleAction(action *flux.Action,
	hitFishTemp map[*fish.Fish]*bullet.Bullet, hitBulletTemp map[*bullet.Bullet]*fish.Fish,
	multiHitFishTemp map[*fish.Fishies]*bullet.Bullet, multiHitBulletTemp map[*bullet.Bullet]*fish.Fishies,
) {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	if action.Key().To() != l.id() {
		return
	}

	switch action.Key().Name() {
	case flux.ActionFluxRegisterDone:
		// do nothing now

	case probability.EMSGID_eResultCall:

		data := action.Payload()[0].(*probability_proto.ResultCall)
		blacklistValue := action.Payload()[1].(uint32)
		depositMultiple := action.Payload()[2].(uint64)
		status := action.Payload()[3].(*auth.MemberStatus)

		// game user or bot user
		secWebSocketKey := action.Key().From()

		// prevent process after game user unregister
		if status.Alive {
			mediator.Service.Wait(secWebSocketKey, data.HitBullet.Uuid) // 4. 29b250c5, 78adf904
		} else {
			mediator.Service.Done(secWebSocketKey, data.HitBullet.Uuid)
		}

		if len(hitFishTemp) != len(hitBulletTemp) {
			logger.Service.Zap.Errorw(Lottery_RESULT_CALL_MAP_SIZE_INVALID,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, Lottery_RESULT_CALL_MAP_SIZE_INVALID)
			return
		}

		isTokenValid := make(chan bool, 1024)
		flux.Send(auth.ActionCheckToken, l.id(), auth.Service.Id, isTokenValid, data.Token)

		if !<-isTokenValid {
			logger.Service.Zap.Errorw(auth.Auth_TOKEN_EXPIRED,
				"GameUser", secWebSocketKey,
				"Token", data.Token,
			)
			errorcode.Service.Fatal(secWebSocketKey, auth.Auth_TOKEN_EXPIRED)
			return
		}

		if l.kind != BOT {
			logger.Service.Zap.Infow("Received Result Call",
				"GameUser", secWebSocketKey,
				"GameRoomUuid", data.RoomUuid,
				"FishType", data.HitFish.SymbolId,
				"FishUuid", data.HitFish.Uuid,
				"BulletType", data.HitBullet.TypeId,
				"BonusUuid", data.HitBullet.BonusUuid,
				"BulletUuid", data.HitBullet.Uuid,
				"MercenaryType", data.MercenaryType,
			)
		}

		gameId := gamesetting.Service.GameId(secWebSocketKey)

		switch gameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			if !psf_on_00001_Check(data.HitBullet.TypeId, data.HitFish.SymbolId) {
				logger.Service.Zap.Errorw(Lottery_HIT_FISH_BULLET_NOT_ALLOWED,
					"GameUser", secWebSocketKey,
					"BulletType", data.HitBullet.TypeId,
					"FishType", data.HitFish.SymbolId,
				)
				errorcode.Service.Fatal(secWebSocketKey, Lottery_HIT_FISH_BULLET_NOT_ALLOWED)
				return
			}

		case models.PSF_ON_00002, models.PSF_ON_20002:
			if !psf_on_00002_Check(data.HitBullet.TypeId, data.HitFish.SymbolId) {
				logger.Service.Zap.Errorw(Lottery_HIT_FISH_BULLET_NOT_ALLOWED,
					"GameUser", secWebSocketKey,
					"BulletType", data.HitBullet.TypeId,
					"FishType", data.HitFish.SymbolId,
				)
				errorcode.Service.Fatal(secWebSocketKey, Lottery_HIT_FISH_BULLET_NOT_ALLOWED)
				return
			}

		case models.PSF_ON_00003:
			if !psf_on_00003_Check(data.HitBullet.TypeId, data.HitFish.SymbolId) {
				logger.Service.Zap.Errorw(Lottery_HIT_FISH_BULLET_NOT_ALLOWED,
					"GameUser", secWebSocketKey,
					"BulletType", data.HitBullet.TypeId,
					"FishType", data.HitFish.SymbolId,
				)
				errorcode.Service.Fatal(secWebSocketKey, Lottery_HIT_FISH_BULLET_NOT_ALLOWED)
				return
			}

		case models.PSF_ON_00004:
			if !psf_on_00004_Check(data.HitBullet.TypeId, data.HitFish.SymbolId) {
				logger.Service.Zap.Errorw(Lottery_HIT_FISH_BULLET_NOT_ALLOWED,
					"GameUser", secWebSocketKey,
					"BulletType", data.HitBullet.TypeId,
					"FishType", data.HitFish.SymbolId,
				)
				errorcode.Service.Fatal(secWebSocketKey, Lottery_HIT_FISH_BULLET_NOT_ALLOWED)
				return
			}

		case models.PSF_ON_00005:
			if !psf_on_00005_Check(data.HitBullet.TypeId, data.HitFish.SymbolId) {
				logger.Service.Zap.Errorw(Lottery_HIT_FISH_BULLET_NOT_ALLOWED,
					"GameUser", secWebSocketKey,
					"BulletType", data.HitBullet.TypeId,
					"FishType", data.HitFish.SymbolId,
				)
				errorcode.Service.Fatal(secWebSocketKey, Lottery_HIT_FISH_BULLET_NOT_ALLOWED)
				return
			}

		case models.PSF_ON_00006:
			if !psf_on_00006_Check(data.HitBullet.TypeId, data.HitFish.SymbolId) {
				logger.Service.Zap.Errorw(Lottery_HIT_FISH_BULLET_NOT_ALLOWED,
					"GameUser", secWebSocketKey,
					"BulletType", data.HitBullet.TypeId,
					"FishType", data.HitFish.SymbolId,
				)
				errorcode.Service.Fatal(secWebSocketKey, Lottery_HIT_FISH_BULLET_NOT_ALLOWED)
				return
			}

		case models.PSF_ON_00007:
			if !psf_on_00007_Check(data.HitBullet.TypeId, data.HitFish.SymbolId) {
				logger.Service.Zap.Errorw(Lottery_HIT_FISH_BULLET_NOT_ALLOWED,
					"GameUser", secWebSocketKey,
					"BulletType", data.HitBullet.TypeId,
					"FishType", data.HitFish.SymbolId,
				)
				errorcode.Service.Fatal(secWebSocketKey, Lottery_HIT_FISH_BULLET_NOT_ALLOWED)
				return
			}

		case models.RKF_H5_00001:
			if !rkf_h5_00001_Check(data.HitBullet.TypeId, data.HitFish.SymbolId) {
				logger.Service.Zap.Errorw(Lottery_HIT_FISH_BULLET_NOT_ALLOWED,
					"GameUser", secWebSocketKey,
					"BulletType", data.HitBullet.TypeId,
					"FishType", data.HitFish.SymbolId,
				)
				errorcode.Service.Fatal(secWebSocketKey, Lottery_HIT_FISH_BULLET_NOT_ALLOWED)
				return
			}

		default:
			logger.Service.Zap.Errorw(Lottery_GAME_ID_INVALID,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, Lottery_GAME_ID_INVALID)
			return
		}

		hitFish := fish.New(
			gameId,
			data.RoomUuid,
			data.HitFish.Uuid,
			gamesetting.Service.MathModuleId(secWebSocketKey),
			data.SeatId,
			data.HitFish.SymbolId,
			-1,
			-1,
			data.HitFish.RtpId,
			nil,
			nil,
			"",
			-1,
		)

		hitBullet := bullet.New(
			data.HitBullet.BonusUuid,
			data.HitBullet.Uuid,
			secWebSocketKey,
			data.RoomUuid,
			-1,
			-1,
			-1,
			-1,
			data.HitBullet.TypeId,
			-1,
			-1,
			-1,
			data.HitBullet.ShootMode,
			data.HitBullet.Target,
		)

		hitFishTemp[hitFish] = hitBullet
		hitBulletTemp[hitBullet] = hitFish

		var mercenaryType int32 = -1
		switch gameId {
		case models.PSF_ON_00003:
			mercenaryType = data.MercenaryType
		default:
			mercenaryType = -1
		}

		flux.Send(fish.ActionFishExistCheck, l.id(), fish.Service.HashId(data.RoomUuid), hitFish, mercenaryType, blacklistValue, depositMultiple)

	case probability.EMSGID_eMultiResultCall:
		data := action.Payload()[0].(*probability_proto.MultiResultCall)
		blacklistValue := action.Payload()[1].(uint32)
		depositMultiple := action.Payload()[2].(uint64)
		status := action.Payload()[3].(*auth.MemberStatus)

		secWebSocketKey := action.Key().From()

		// prevent process after game user unregister
		if status.Alive {
			mediator.Service.Wait(secWebSocketKey, data.HitBullet.Uuid) // 4. 29b250c5, 78adf904
		} else {
			mediator.Service.Done(secWebSocketKey, data.HitBullet.Uuid)
		}

		if len(multiHitFishTemp) != len(multiHitBulletTemp) {
			logger.Service.Zap.Errorw(Lottery_RESULT_CALL_MAP_SIZE_INVALID,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, Lottery_RESULT_CALL_MAP_SIZE_INVALID)
			return
		}

		isTokenValid := make(chan bool, 1024)
		flux.Send(auth.ActionCheckToken, l.id(), auth.Service.Id, isTokenValid, data.Token)
		if !<-isTokenValid {
			logger.Service.Zap.Errorw(auth.Auth_TOKEN_EXPIRED,
				"GameUser", secWebSocketKey,
				"Token", data.Token,
			)
			errorcode.Service.Fatal(secWebSocketKey, auth.Auth_TOKEN_EXPIRED)
			return
		}

		fishTemp := make([]*fish.Fish, 0, 10)
		gameId := gamesetting.Service.GameId(secWebSocketKey)
		for i := 0; i < len(data.HitFishies); i++ {
			if l.kind != BOT {
				logger.Service.Zap.Infow("Received Result Call",
					"GameUser", secWebSocketKey,
					"GameRoomUuid", data.RoomUuid,
					"FishType", data.HitFishies[i].SymbolId,
					"FishUuid", data.HitFishies[i].Uuid,
					"BulletType", data.HitBullet.TypeId,
					"BonusUuid", data.HitBullet.BonusUuid,
					"BulletUuid", data.HitBullet.Uuid,
				)
			}

			switch gameId {
			case models.PSF_ON_00003:
				if !psf_on_00003_Check(data.HitBullet.TypeId,
					data.HitFishies[i].SymbolId) {
					logger.Service.Zap.Errorw(Lottery_HIT_FISH_BULLET_NOT_ALLOWED,
						"GameUser", secWebSocketKey,
						"BulletType", data.HitBullet.TypeId,
						"FishType", data.HitFishies[i].SymbolId,
					)
					errorcode.Service.Fatal(secWebSocketKey, Lottery_HIT_FISH_BULLET_NOT_ALLOWED)
					return
				}

			default:
				logger.Service.Zap.Errorw(Lottery_GAME_ID_INVALID,
					"GameUser", secWebSocketKey,
				)
				errorcode.Service.Fatal(secWebSocketKey, Lottery_GAME_ID_INVALID)
				return
			}

			hitFish := fish.New(
				gameId,
				data.RoomUuid,
				data.HitFishies[i].Uuid,
				gamesetting.Service.MathModuleId(secWebSocketKey),
				data.SeatId,
				data.HitFishies[i].SymbolId,
				-1,
				-1,
				data.HitFishies[i].RtpId,
				nil,
				nil,
				"",
				-1,
			)

			fishTemp = append(fishTemp, hitFish)
			//fishies = append(fishies, hitFish)
		}
		fishies := &fish.Fishies{
			Fishies: fishTemp,
		}

		hitBullet := bullet.New(
			data.HitBullet.BonusUuid,
			data.HitBullet.Uuid,
			secWebSocketKey,
			data.RoomUuid,
			-1,
			-1,
			-1,
			-1,
			data.HitBullet.TypeId,
			-1,
			-1,
			-1,
			data.HitBullet.ShootMode,
			data.HitBullet.Target,
		)

		multiHitFishTemp[fishies] = hitBullet
		multiHitBulletTemp[hitBullet] = fishies

		var mercenaryType int32 = -1
		switch gameId {
		case models.PSF_ON_00003:
			mercenaryType = data.MercenaryType
		default:
			mercenaryType = -1
		}
		flux.Send(fish.ActionMultiFishExistCheck, l.id(), fish.Service.HashId(data.RoomUuid), fishies, mercenaryType, blacklistValue, depositMultiple)

	case fish.ActionFishExist:
		hitFish := action.Payload()[0].(*fish.Fish)
		mercenaryType := action.Payload()[1].(int32)
		blacklistValue := action.Payload()[2].(uint32)
		depositMultiple := action.Payload()[3].(uint64)
		hitBullet := hitFishTemp[hitFish]

		flux.Send(bullet.ActionBulletExistCheck, l.id(), bullet.Service.HashId(hitBullet.SecWebSocketKey), hitBullet, mercenaryType, blacklistValue, depositMultiple, true)

	case fish.ActionFishNotExist:
		hitFish := action.Payload()[0].(*fish.Fish)
		mercenaryType := action.Payload()[1].(int32)
		blacklistValue := action.Payload()[2].(uint32)
		depositMultiple := action.Payload()[3].(uint64)
		hitBullet := hitFishTemp[hitFish]

		flux.Send(bullet.ActionBulletExistCheck, l.id(), bullet.Service.HashId(hitBullet.SecWebSocketKey), hitBullet, mercenaryType, blacklistValue, depositMultiple, false)

		delete(hitFishTemp, hitFish)
		delete(hitBulletTemp, hitBullet)
		l.sendResult(hitFish, hitBullet, nil, 0, mercenaryType)

	case fish.ActionMultiFishiesCheck:
		hitFishies := action.Payload()[0].(*fish.Fishies)
		hitBullet := multiHitFishTemp[hitFishies]
		checkFishiesExist := action.Payload()[1].([]bool)
		mercenaryType := action.Payload()[2].(int32)
		blacklistValue := action.Payload()[3].(uint32)
		depositMultiple := action.Payload()[4].(uint64)

		// Process All hitFishies is not exist
		if l.checkMultiFish(checkFishiesExist) {
			delete(multiHitFishTemp, hitFishies)
			delete(multiHitBulletTemp, hitBullet)

			bonusTimes := make([]int32, 0, len(hitFishies.Fishies))
			for i := 0; i < len(hitFishies.Fishies); i++ {
				bonusTimes = append(bonusTimes, 0)
			}
			l.sendMultiResult(hitFishies.Fishies, hitBullet, nil, bonusTimes, mercenaryType)
		} else {
			flux.Send(bullet.ActionBulletMultiExistCheck, l.id(), bullet.Service.HashId(hitBullet.SecWebSocketKey),
				hitBullet,
				checkFishiesExist,
				mercenaryType,
				blacklistValue,
				depositMultiple,
			)
		}

	case bullet.ActionBulletBonusExist, bullet.ActionBulletExist:
		hitBullet := action.Payload()[0].(*bullet.Bullet)
		mercenaryType := action.Payload()[1].(int32)
		blacklistValue := action.Payload()[2].(uint32)
		depositMultiple := action.Payload()[3].(uint64)
		secWebSocketKey := hitBullet.SecWebSocketKey
		subgameId := gamesetting.Service.SubgameId(secWebSocketKey)

		hitFish := hitBulletTemp[hitBullet]

		switch hitFish.GameId {
		case models.PSF_ON_00003:
			if mercenaryType != -1 {
				rtp.Service.DecreaseMercenaryBullet(secWebSocketKey, hitFish.GameId, subgameId, uint64(mercenaryType))
				rtp.Service.Save(secWebSocketKey, hitFish.GameId, subgameId)
				hitBullet.Bet = 1
			}
		}

		checkDecrease, hitResult, _, _, _, _, _, _, bulletSn := Service.MathProcess(
			secWebSocketKey, hitFish.GameId, hitFish.MathModuleId, hitBullet.Uuid, hitBullet.BonusUuid, hitFish.RtpId, hitBullet.RtpId,
			mercenaryType, hitFish.TypeId, hitBullet.TypeId,
			hitBullet.Bet, hitBullet.Rate, subgameId,
			hitBullet, hitFish,
			l.kind, l.id(),
			len(hitFishTemp), len(hitBulletTemp),
			false, blacklistValue, depositMultiple,
		)

		if !checkDecrease {
			delete(hitFishTemp, hitFish)
			delete(hitBulletTemp, hitBullet)
			l.sendResult(hitFish, hitBullet, nil, 0, mercenaryType)
			return
		}

		// DEBUG Process Weapon Type and Check Fish Is Died
		Service.WriteDebugLog(hitBullet, hitFish, nil, hitResult, nil, secWebSocketKey, hitFish.GameId, subgameId, mercenaryType)

		flux.Send(bullet.ActionBulletHitCall, l.id(), bullet.Service.HashId(hitBullet.SecWebSocketKey),
			hitFish,
			hitBullet,
			hitResult,
			hitResult.Pay,
			mercenaryType,
			hitBullet.Bet,
			bulletSn,
		)

	case bullet.ActionBulletBonusNotExist, bullet.ActionBulletNotExist:
		hitBullet := action.Payload()[0].(*bullet.Bullet)
		mercenaryType := action.Payload()[1].(int32)
		hitFish := hitBulletTemp[hitBullet]
		delete(hitFishTemp, hitFish)
		delete(hitBulletTemp, hitBullet)

		l.sendResult(hitFish, hitBullet, nil, 0, mercenaryType)

	case bullet.ActionBulletMultiBonusExist, bullet.ActionBulletMultiExist:
		hitBullet := action.Payload()[0].(*bullet.Bullet)
		checkFishiesExist := action.Payload()[1].([]bool)
		mercenaryType := action.Payload()[2].(int32)
		blacklistValue := action.Payload()[3].(uint32)
		depositMultiple := action.Payload()[4].(uint64)
		hitFishies := multiHitBulletTemp[hitBullet]
		var hitResults []*probability.Probability
		resultPay := 0

		var bulletSn uint64 = 0

		secWebSocketKey := hitBullet.SecWebSocketKey
		gameId := gamesetting.Service.GameId(secWebSocketKey)
		subgameId := gamesetting.Service.SubgameId(secWebSocketKey)

		bet := hitBullet.Bet
		switch gameId {
		case models.PSF_ON_00003:
			bet = psf_on_00003_processWeapon(hitBullet.TypeId, mercenaryType, len(hitFishies.Fishies), hitBullet.Bet)
			if mercenaryType != -1 {
				rtp.Service.DecreaseMercenaryBullet(secWebSocketKey, gameId, subgameId, uint64(mercenaryType))
				rtp.Service.Save(secWebSocketKey, gameId, subgameId)
				hitBullet.Bet = 1
			}
		}

		var decreaseBet uint64 = 0
		mathModuleId := ""
		for index, fishCheck := range checkFishiesExist {
			if fishCheck {
				mathModuleId = hitFishies.Fishies[index].MathModuleId
				decreaseBet += 1
			}
		}

		//Regular Bullet
		if hitBullet.BonusUuid == "" && mercenaryType == -1 {
			switch wallet.Service.WalletCategory(secWebSocketKey) {
			case wallet.CategoryHost:
				// FishWallet Cent is 0 then Deposit Bet x 20
				if cent, err := wallet.Service.CentCache(secWebSocketKey); err == nil {
					if cent <= 0 || cent < decreaseBet*hitBullet.Rate {
						accountingSn := accountingmanager.Service.GetAccountingSn(secWebSocketKey)
						am := accountingmanager.Service.Get(secWebSocketKey)
						betAmt, winAmt := am.GetWalletHostAmt(secWebSocketKey)
						//todo top one
						//wallet.Service.EndAccounting(secWebSocketKey, decreaseBet*hitBullet.Rate, accountingSn, false)
						wallet.Service.EndAccounting(secWebSocketKey, decreaseBet*hitBullet.Rate, false)
						time.Sleep(500 * time.Millisecond)
						if cent > 0 {
							wallet.Service.Withdraw(secWebSocketKey, accountingSn, cent, betAmt, winAmt)
						}

						balance, _ := wallet.Service.BalanceCache(secWebSocketKey)

						depositCent := hitBullet.Bet * hitBullet.Rate * depositMultiple
						if balance < depositCent {
							depositCent = balance
						}

						wallet.Service.Deposit(secWebSocketKey, depositCent, accountingSn)
					}

					if !wallet.Service.DecreaseCent(secWebSocketKey, hitBullet.Uuid, decreaseBet*hitBullet.Rate) {
						delete(multiHitFishTemp, hitFishies)
						delete(multiHitBulletTemp, hitBullet)

						bonusTimes := make([]int32, 0, len(hitFishies.Fishies))
						for i := 0; i < len(hitFishies.Fishies); i++ {
							bonusTimes = append(bonusTimes, 0)
						}

						l.sendMultiResult(hitFishies.Fishies, hitBullet, nil, bonusTimes, mercenaryType)
						return
					}

				} else {
					logger.Service.Zap.Errorw("Shared Wallet Get Cent Error",
						"secWebSocketKey", secWebSocketKey,
					)
				}
				bulletSn = accountingmanager.Service.GetAccountingSn(secWebSocketKey)

			default:
				if !wallet.Service.DecreaseBalance(secWebSocketKey, hitBullet.Uuid, decreaseBet*hitBullet.Rate) {
					delete(multiHitFishTemp, hitFishies)
					delete(multiHitBulletTemp, hitBullet)

					bonusTimes := make([]int32, 0, len(hitFishies.Fishies))
					for i := 0; i < len(hitFishies.Fishies); i++ {
						bonusTimes = append(bonusTimes, 0)
					}

					l.sendMultiResult(hitFishies.Fishies, hitBullet, nil, bonusTimes, mercenaryType)
					return
				}
			}

			rtp.Service.Decrease(gameId, subgameId, mathModuleId, secWebSocketKey, decreaseBet*hitBullet.Rate)
			rtp.Service.Save(secWebSocketKey, gameId, subgameId)
			hitBullet.RtpId = rtp.Service.RtpId(secWebSocketKey, gameId, subgameId)
		}
		hitBullet.RtpId = probability.Service.UseRtp(gameId, mathModuleId, hitBullet.RtpId, hitBullet.TypeId, mercenaryType)

		hitFreeBullet := 0
		hitPay := make([]uint64, 0)
		hitMultiplier := make([]uint64, 0)
		for i, hitFish := range hitFishies.Fishies {
			if checkFishiesExist[i] {
				hitResult := probability.Service.Calc(
					gameId,
					mathModuleId,
					hitBullet.RtpId,
					hitFish.TypeId,
					hitBullet.TypeId,
					rtp.Service.RtpBudget(secWebSocketKey, gameId, subgameId, mathModuleId),
					bet*hitBullet.Rate,
				)

				if probability.Service.Random(1000) < uint64(blacklistValue) {
					hitResult = probability.Service.GetNewSimplePay(int(hitFish.TypeId))
				}
				hitResults = append(hitResults, hitResult)

				resultPay += hitResult.Pay
				hitPay = append(hitPay, uint64(hitResult.Pay))
				hitMultiplier = append(hitMultiplier, uint64(hitResult.Multiplier))
				if hitResult.Bullet > 0 {
					hitFreeBullet += hitResult.Bullet
				}

				if l.kind == PLAYER {
					activity.Service.Record(secWebSocketKey, hitFish, hitBullet, hitResult)
				}
			} else {
				// Process Fish is Not Exist
				hitResults = append(hitResults, nil)
			}

		}

		rtp.Service.IncreaseBulletCollection(secWebSocketKey, gameId, subgameId, uint64(hitFreeBullet))
		rtp.Service.Save(secWebSocketKey, gameId, subgameId)

		var increaseCent uint64 = 0
		for index, pay := range hitPay {
			increaseCent += pay * hitMultiplier[index] * bet * hitBullet.Rate
		}
		switch wallet.Service.WalletCategory(secWebSocketKey) {
		case wallet.CategoryHost:
			/*accountingSn := accountingmanager.Service.GetAccountingSn(secWebSocketKey)
			if accountingSn == 0 {
				balance, _ := wallet.Service.BalanceCache(secWebSocketKey)

				depositCent := hitBullet.Bet * hitBullet.Rate * depositMultiple
				if balance < depositCent {
					depositCent = balance
				}

				wallet.Service.Deposit(secWebSocketKey, depositCent, accountingSn)
			}*/
			wallet.Service.IncreaseCent(secWebSocketKey, hitBullet.Uuid, increaseCent)

		default:
			wallet.Service.IncreaseBalance(secWebSocketKey, hitBullet.Uuid, increaseCent, 0, 0)
		}

		balance, err := wallet.Service.BalanceCache(secWebSocketKey)
		if err != nil {
			// TODO JOHNNY game user is already disconnect
			logger.Service.Zap.Errorw("BUG",
				"GameUser", secWebSocketKey,
				"HitBulletTemp", len(multiHitBulletTemp),
				"HitFishTemp", len(multiHitFishTemp),
				"Error", err,
			)
		}
		hitBullet.PlayerCent = balance

		// DEBUG Process Weapon Type and Check Fish Is Died And Free Bullet
		Service.WriteDebugLog(hitBullet, nil, hitFishies, nil, hitResults, secWebSocketKey, gameId, subgameId, mercenaryType)

		flux.Send(bullet.ActionBulletMultiCall, l.id(), bullet.Service.HashId(hitBullet.SecWebSocketKey),
			hitFishies,
			hitBullet,
			hitResults,
			resultPay,
			mercenaryType,
			bulletSn,
		)

	case bullet.ActionBulletMultiBonusNotExist, bullet.ActionBulletMultiNotExist:
		hitBullet := action.Payload()[0].(*bullet.Bullet)
		mercenaryType := action.Payload()[1].(int32)
		hitFishies := multiHitBulletTemp[hitBullet]

		bonusTimes := make([]int32, 0, len(hitFishies.Fishies))
		for i := 0; i < len(hitFishies.Fishies); i++ {
			bonusTimes = append(bonusTimes, 0)
		}

		delete(multiHitFishTemp, hitFishies)
		delete(multiHitBulletTemp, hitBullet)
		l.sendMultiResult(hitFishies.Fishies, hitBullet, nil, bonusTimes, mercenaryType)

	case bullet.ActionBulletBonusRecall:
		bonusBullet := action.Payload()[0].(*bullet.Bullet)
		accountingSn := action.Payload()[1].(uint64)
		hitFish := bonusBullet.ExtraData[0].(*fish.Fish)
		hitBullet := bonusBullet.ExtraData[1].(*bullet.Bullet)
		hitResult := bonusBullet.ExtraData[2].(*probability.Probability)

		// Only PSF-ON-00003 Have Mercenary
		var mercenaryType int32 = -1
		gameId := gamesetting.Service.GameId(hitBullet.SecWebSocketKey)
		switch gameId {
		case models.PSF_ON_00003:
			mercenaryType = action.Payload()[2].(int32)
		default:
			break
		}

		if l.kind == PLAYER && l.iHandler != nil {
			l.iHandler.BulletBonusHandler(hitBullet.SecWebSocketKey, bonusBullet, accountingSn)
			// Process Fish have pay and slot pay
			if gameId == models.PSF_ON_00003 && !accountingmanager.CheckBonusBullet(int32(hitResult.FishTypeId)) {
				l.iHandler.BulletHitHandler(hitBullet.SecWebSocketKey, hitFish, hitBullet, hitResult, accountingSn)
			}
		}

		l.sendResult(hitFish, hitBullet, hitResult, bonusBullet.Bullets, mercenaryType)

	case bullet.ActionBulletHitRecall:
		hitFish := action.Payload()[0].(*fish.Fish)
		hitBullet := action.Payload()[1].(*bullet.Bullet)
		hitResult := action.Payload()[2].(*probability.Probability)
		accountingSn := action.Payload()[3].(uint64)
		secWebSocketKey := hitBullet.SecWebSocketKey

		// Only PSF-ON-00003 Have Mercenary
		var mercenaryType int32 = -1
		switch gamesetting.Service.GameId(secWebSocketKey) {
		case models.PSF_ON_00003:
			mercenaryType = action.Payload()[4].(int32)
		default:
			break
		}

		if l.kind == PLAYER && l.iHandler != nil {
			l.iHandler.BulletHitHandler(hitBullet.SecWebSocketKey, hitFish, hitBullet, hitResult, accountingSn)
		}

		// Waiting BulletHitHandler Complete
		time.Sleep(100 * time.Millisecond)
		// Process Stored Value Update Balance
		switch wallet.Service.WalletCategory(secWebSocketKey) {
		case wallet.CategoryHost:
			balance, _ := wallet.Service.BalanceCache(hitBullet.SecWebSocketKey)
			cent, _ := wallet.Service.CentCache(hitBullet.SecWebSocketKey)
			hitBullet.PlayerCent = balance + cent
		default:
			hitBullet.PlayerCent, _ = wallet.Service.BalanceCache(hitBullet.SecWebSocketKey)
		}

		switch gamesetting.Service.GameId(secWebSocketKey) {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			psf_on_00001_TriggerBonus(hitFish, hitBullet, hitResult)
			isDied := psf_on_00001_Fish(hitFish, hitBullet, hitResult, accountingSn)

			if hitResult.TriggerIconId != 26 && hitResult.TriggerIconId != 27 {
				l.sendResult(hitFish, hitBullet, hitResult, 0, mercenaryType)
			}

			if isDied {
				flux.Send(fish.ActionFishDie, Service.HashId(hitBullet.SecWebSocketKey), fish.Service.HashId(hitFish.RoomUuid), hitFish)
			}

			flux.Send(broadcaster.ActionBroadcasterBigWin, l.id(), broadcaster.Service.Id, broadcaster.New(
				secWebSocketKey,
				gamesetting.Service.GameId(secWebSocketKey),
				wallet.Service.MemberId(secWebSocketKey),
				gamesetting.Service.BetList(secWebSocketKey),
				gamesetting.Service.RateList(secWebSocketKey),
				gamesetting.Service.MathModuleId(secWebSocketKey),
				int(gamesetting.Service.RateIndex(secWebSocketKey)),
				hitFish.TypeId,
				uint64(hitResult.Pay),
				uint64(hitResult.Pay)*hitBullet.Bet*hitBullet.Rate,
				gamesetting.Service.Rate(secWebSocketKey),
			))
			delete(hitFishTemp, hitFish)
			delete(hitBulletTemp, hitBullet)

		case models.PSF_ON_00002, models.PSF_ON_20002:
			psf_on_00002_TriggerBonus(hitFish, hitBullet, hitResult)
			psf_on_00002_ExtraTriggerBonus(hitFish, hitBullet, hitResult)
			isDied := psf_on_00002_Fish(hitFish, hitBullet, hitResult, accountingSn)

			if hitResult.TriggerIconId != 26 && hitResult.TriggerIconId != 27 {
				l.sendResult(hitFish, hitBullet, hitResult, 0, mercenaryType)
			}

			if isDied {
				flux.Send(fish.ActionFishDie, Service.HashId(hitBullet.SecWebSocketKey), fish.Service.HashId(hitFish.RoomUuid), hitFish)
			}

			flux.Send(broadcaster.ActionBroadcasterBigWin, l.id(), broadcaster.Service.Id, broadcaster.New(
				secWebSocketKey,
				gamesetting.Service.GameId(secWebSocketKey),
				wallet.Service.MemberId(secWebSocketKey),
				gamesetting.Service.BetList(secWebSocketKey),
				gamesetting.Service.RateList(secWebSocketKey),
				gamesetting.Service.MathModuleId(secWebSocketKey),
				int(gamesetting.Service.RateIndex(secWebSocketKey)),
				hitFish.TypeId,
				uint64(hitResult.Pay),
				uint64(hitResult.Pay)*hitBullet.Bet*hitBullet.Rate,
				gamesetting.Service.Rate(secWebSocketKey),
			))
			delete(hitFishTemp, hitFish)
			delete(hitBulletTemp, hitBullet)

		case models.PSF_ON_00003:
			psf_on_00003_TriggerBonus(hitFish, hitBullet, hitResult)
			isDied, _ := psf_on_00003_Fish(hitFish, hitBullet, hitResult, false, mercenaryType, 1, accountingSn)

			if checkTriggerIconId(hitResult.TriggerIconId) {
				l.sendResult(hitFish, hitBullet, hitResult, 0, mercenaryType)
			}

			if isDied {
				flux.Send(fish.ActionFishDie, Service.HashId(hitBullet.SecWebSocketKey), fish.Service.HashId(hitFish.RoomUuid), hitFish)
			}

			// Process Broadcast x80 x100
			hitFishId, bigWinShow := psf_on_00003_broadcasterCheck(hitFish.TypeId, isDied)

			flux.Send(broadcaster.ActionBroadcasterBigWin, l.id(), broadcaster.Service.Id, broadcaster.New(
				secWebSocketKey,
				gamesetting.Service.GameId(secWebSocketKey),
				wallet.Service.MemberId(secWebSocketKey),
				gamesetting.Service.BetList(secWebSocketKey),
				gamesetting.Service.RateList(secWebSocketKey),
				gamesetting.Service.MathModuleId(secWebSocketKey),
				int(gamesetting.Service.RateIndex(secWebSocketKey)),
				hitFishId,
				uint64(hitResult.Pay),
				uint64(hitResult.Pay)*hitBullet.Bet*hitBullet.Rate,
				gamesetting.Service.Rate(secWebSocketKey),
			), bigWinShow)
			delete(hitFishTemp, hitFish)
			delete(hitBulletTemp, hitBullet)

		case models.PSF_ON_00004:
			psf_on_00004_TriggerBonus(hitFish, hitBullet, hitResult)
			isDied := psf_on_00004_Fish(hitFish, hitBullet, hitResult, accountingSn)

			if hitResult.TriggerIconId != 201 && hitResult.TriggerIconId != 202 {
				l.sendResult(hitFish, hitBullet, hitResult, 0, mercenaryType)
			}

			if isDied {
				flux.Send(fish.ActionFishDie, Service.HashId(hitBullet.SecWebSocketKey), fish.Service.HashId(hitFish.RoomUuid), hitFish)
			}

			flux.Send(broadcaster.ActionBroadcasterBigWin, l.id(), broadcaster.Service.Id, broadcaster.New(
				secWebSocketKey,
				gamesetting.Service.GameId(secWebSocketKey),
				wallet.Service.MemberId(secWebSocketKey),
				gamesetting.Service.BetList(secWebSocketKey),
				gamesetting.Service.RateList(secWebSocketKey),
				gamesetting.Service.MathModuleId(secWebSocketKey),
				int(gamesetting.Service.RateIndex(secWebSocketKey)),
				hitFish.TypeId,
				uint64(hitResult.Pay),
				uint64(hitResult.Pay)*hitBullet.Bet*hitBullet.Rate,
				gamesetting.Service.Rate(secWebSocketKey),
			))
			delete(hitFishTemp, hitFish)
			delete(hitBulletTemp, hitBullet)

		case models.PSF_ON_00005:
			psf_on_00005_TriggerBonus(hitFish, hitBullet, hitResult)
			psf_on_00005_ExtraTriggerBonus(hitFish, hitBullet, hitResult)
			isDied := psf_on_00005_Fish(hitFish, hitBullet, hitResult, accountingSn)

			if hitResult.TriggerIconId != 201 && hitResult.TriggerIconId != 202 {
				l.sendResult(hitFish, hitBullet, hitResult, 0, mercenaryType)
			}

			if isDied {
				flux.Send(fish.ActionFishDie, Service.HashId(hitBullet.SecWebSocketKey), fish.Service.HashId(hitFish.RoomUuid), hitFish)
			}

			flux.Send(broadcaster.ActionBroadcasterBigWin, l.id(), broadcaster.Service.Id, broadcaster.New(
				secWebSocketKey,
				gamesetting.Service.GameId(secWebSocketKey),
				wallet.Service.MemberId(secWebSocketKey),
				gamesetting.Service.BetList(secWebSocketKey),
				gamesetting.Service.RateList(secWebSocketKey),
				gamesetting.Service.MathModuleId(secWebSocketKey),
				int(gamesetting.Service.RateIndex(secWebSocketKey)),
				hitFish.TypeId,
				uint64(hitResult.Pay),
				uint64(hitResult.Pay)*hitBullet.Bet*hitBullet.Rate,
				gamesetting.Service.Rate(secWebSocketKey),
			))
			delete(hitFishTemp, hitFish)
			delete(hitBulletTemp, hitBullet)

		case models.PSF_ON_00006:
			psf_on_00006_TriggerBonus(hitFish, hitBullet, hitResult)
			psf_on_00006_ExtraTriggerBonus(hitFish, hitBullet, hitResult)
			isDied := psf_on_00006_Fish(hitFish, hitBullet, hitResult, accountingSn)

			if hitResult.TriggerIconId != 201 && hitResult.TriggerIconId != 202 {
				l.sendResult(hitFish, hitBullet, hitResult, 0, mercenaryType)
			}

			if isDied {
				flux.Send(fish.ActionFishDie, Service.HashId(hitBullet.SecWebSocketKey), fish.Service.HashId(hitFish.RoomUuid), hitFish)
			}

			flux.Send(broadcaster.ActionBroadcasterBigWin, l.id(), broadcaster.Service.Id, broadcaster.New(
				secWebSocketKey,
				gamesetting.Service.GameId(secWebSocketKey),
				wallet.Service.MemberId(secWebSocketKey),
				gamesetting.Service.BetList(secWebSocketKey),
				gamesetting.Service.RateList(secWebSocketKey),
				gamesetting.Service.MathModuleId(secWebSocketKey),
				int(gamesetting.Service.RateIndex(secWebSocketKey)),
				hitFish.TypeId,
				uint64(hitResult.Pay),
				uint64(hitResult.Pay)*hitBullet.Bet*hitBullet.Rate,
				gamesetting.Service.Rate(secWebSocketKey),
			))
			delete(hitFishTemp, hitFish)
			delete(hitBulletTemp, hitBullet)

		case models.PSF_ON_00007:
			psf_on_00007_TriggerBonus(hitFish, hitBullet, hitResult)
			psf_on_00007_ExtraTriggerBonus(hitFish, hitBullet, hitResult)
			isDied := psf_on_00007_Fish(hitFish, hitBullet, hitResult, accountingSn)

			if hitResult.TriggerIconId != 201 && hitResult.TriggerIconId != 202 {
				l.sendResult(hitFish, hitBullet, hitResult, 0, mercenaryType)
			}

			if isDied {
				flux.Send(fish.ActionFishDie, Service.HashId(hitBullet.SecWebSocketKey), fish.Service.HashId(hitFish.RoomUuid), hitFish)
			}

			flux.Send(broadcaster.ActionBroadcasterBigWin, l.id(), broadcaster.Service.Id, broadcaster.New(
				secWebSocketKey,
				gamesetting.Service.GameId(secWebSocketKey),
				wallet.Service.MemberId(secWebSocketKey),
				gamesetting.Service.BetList(secWebSocketKey),
				gamesetting.Service.RateList(secWebSocketKey),
				gamesetting.Service.MathModuleId(secWebSocketKey),
				int(gamesetting.Service.RateIndex(secWebSocketKey)),
				hitFish.TypeId,
				uint64(hitResult.Pay),
				uint64(hitResult.Pay)*hitBullet.Bet*hitBullet.Rate,
				gamesetting.Service.Rate(secWebSocketKey),
			))
			delete(hitFishTemp, hitFish)
			delete(hitBulletTemp, hitBullet)

		case models.RKF_H5_00001:
			rkf_h5_00001_TriggerBonus(hitFish, hitBullet, hitResult)
			rkf_h5_00001_ExtraTriggerBonus(hitFish, hitBullet, hitResult)
			isDied := rkf_h5_00001_Fish(hitFish, hitBullet, hitResult, accountingSn)

			if hitResult.TriggerIconId != 201 && hitResult.TriggerIconId != 202 {
				l.sendResult(hitFish, hitBullet, hitResult, 0, mercenaryType)
			}

			if isDied {
				flux.Send(fish.ActionFishDie, Service.HashId(hitBullet.SecWebSocketKey), fish.Service.HashId(hitFish.RoomUuid), hitFish)
			}

			flux.Send(broadcaster.ActionBroadcasterBigWin, l.id(), broadcaster.Service.Id, broadcaster.New(
				secWebSocketKey,
				gamesetting.Service.GameId(secWebSocketKey),
				wallet.Service.MemberId(secWebSocketKey),
				gamesetting.Service.BetList(secWebSocketKey),
				gamesetting.Service.RateList(secWebSocketKey),
				gamesetting.Service.MathModuleId(secWebSocketKey),
				int(gamesetting.Service.RateIndex(secWebSocketKey)),
				hitFish.TypeId,
				uint64(hitResult.Pay),
				uint64(hitResult.Pay)*hitBullet.Bet*hitBullet.Rate,
				gamesetting.Service.Rate(secWebSocketKey),
			))
			delete(hitFishTemp, hitFish)
			delete(hitBulletTemp, hitBullet)

		default:
			delete(hitFishTemp, hitFish)
			delete(hitBulletTemp, hitBullet)

			logger.Service.Zap.Errorw(Lottery_GAME_ID_INVALID,
				"GameUser", secWebSocketKey,
			)
			errorcode.Service.Fatal(secWebSocketKey, Lottery_GAME_ID_INVALID)
		}

	case bullet.ActionBulletMultiRecall:
		hitFishies := action.Payload()[0].(*fish.Fishies)
		hitBullet := action.Payload()[1].(*bullet.Bullet)
		hitResults := action.Payload()[2].([]*probability.Probability)
		accountingSn := action.Payload()[3].(uint64)
		mercenaryType := action.Payload()[4].(int32)
		secWebSocketKey := hitBullet.SecWebSocketKey

		bonusBullets := make([]*bullet.Bullet, len(hitFishies.Fishies))
		bonusTimes := make([]int32, len(hitFishies.Fishies))

		bet := psf_on_00003_processWeapon(hitBullet.TypeId, mercenaryType, len(hitFishies.Fishies), hitBullet.Bet)

		if l.kind == PLAYER && l.iHandler != nil {
			l.iHandler.BulletMultiHitHandler(secWebSocketKey, hitFishies, hitBullet, hitResults, bet, accountingSn)
		}

		for i := 0; i < len(hitFishies.Fishies); i++ {
			switch wallet.Service.WalletCategory(secWebSocketKey) {
			case wallet.CategoryHost:
				balance, _ := wallet.Service.BalanceCache(secWebSocketKey)
				cent, _ := wallet.Service.CentCache(secWebSocketKey)
				hitBullet.PlayerCent = balance + cent
			default:
				hitBullet.PlayerCent, _ = wallet.Service.BalanceCache(secWebSocketKey)
			}

			if hitResults[i] != nil {
				switch gamesetting.Service.GameId(secWebSocketKey) {
				case models.PSF_ON_00003:
					psf_on_00003_TriggerBonus(hitFishies.Fishies[i], hitBullet, hitResults[i])
					isDied, bonusBullet := psf_on_00003_Fish(hitFishies.Fishies[i], hitBullet, hitResults[i], true, mercenaryType, len(hitFishies.Fishies), accountingSn)

					bonusBullets[i] = bonusBullet
					if bonusBullet != nil {
						bonusTimes[i] = bonusBullet.Bullets
					} else {
						bonusTimes[i] = 0
					}

					if isDied {
						flux.Send(fish.ActionFishDie, Service.HashId(hitBullet.SecWebSocketKey), fish.Service.HashId(hitFishies.Fishies[i].RoomUuid), hitFishies.Fishies[i])
					}

					// Process Broadcast x80 x100
					hitFishId, bigWinShow := psf_on_00003_broadcasterCheck(hitFishies.Fishies[i].TypeId, isDied)

					flux.Send(broadcaster.ActionBroadcasterBigWin, l.id(), broadcaster.Service.Id, broadcaster.New(
						secWebSocketKey,
						gamesetting.Service.GameId(secWebSocketKey),
						wallet.Service.MemberId(secWebSocketKey),
						gamesetting.Service.BetList(secWebSocketKey),
						gamesetting.Service.RateList(secWebSocketKey),
						gamesetting.Service.MathModuleId(secWebSocketKey),
						int(gamesetting.Service.RateIndex(secWebSocketKey)),
						hitFishId,
						uint64(hitResults[i].Pay),
						uint64(hitResults[i].Pay)*bet*hitBullet.Rate,
						gamesetting.Service.Rate(secWebSocketKey),
					), bigWinShow)

				default:
					delete(multiHitFishTemp, hitFishies)
					delete(multiHitBulletTemp, hitBullet)

					logger.Service.Zap.Errorw(Lottery_GAME_ID_INVALID,
						"GameUser", secWebSocketKey,
					)
					errorcode.Service.Fatal(secWebSocketKey, Lottery_GAME_ID_INVALID)
				}
			}
		}

		// 合併特殊武器及一般結果
		flux.Send(bullet.ActionBulletMultiBonusCall, l.id(), bullet.Service.HashId(hitBullet.SecWebSocketKey), bonusBullets)
		l.sendMultiResult(hitFishies.Fishies, hitBullet, hitResults, bonusTimes, mercenaryType)

		delete(multiHitFishTemp, hitFishies)
		delete(multiHitBulletTemp, hitBullet)

	case probability.EMSGID_eMercenaryOpenCall:
		seat_id := action.Payload()[0].(int32)
		bulletCollection := action.Payload()[1].(uint64)
		mercenaryType := action.Payload()[2].(int32)
		mercenaryBullet := action.Payload()[3].(uint64)
		openResult := action.Payload()[4].(bool)

		secWebSocketKey := action.Key().From()

		data := &probability_proto.MercenaryOpenRecall{
			Msgid: common_fish_proto.EMSGID_eMercenaryOpenRecall,
		}

		if openResult {
			data.StatusCode = common_proto.Status_kSuccess

		} else {
			data.StatusCode = common_proto.Status_kInvalid

			logger.Service.Zap.Warnw(Lottery_MERCENARY_OPEN_RECALL_INVALID,
				"GameUser", secWebSocketKey,
				"BulletCollection", bulletCollection,
				"MercenaryType", mercenaryType,
				"MercenaryBullet", mercenaryBullet,
			)
		}

		data.SeatId = seat_id
		data.BulletCollection = bulletCollection
		data.MercenaryType = uint64(mercenaryType)
		data.MercenaryBullet = mercenaryBullet

		dataBuffer, _ := proto.Marshal(data)
		flux.Send(probability.EMSGID_eMercenaryOpenRecall, l.id(), secWebSocketKey, dataBuffer)

	case bullet.ActionBonusBulletShooted:
		bonusBullet := action.Payload()[0].(*bullet.Bullet)
		if l.kind == PLAYER && l.iHandler != nil {
			l.iHandler.BonusBulletShootedHandler(bonusBullet.SecWebSocketKey, bonusBullet)
		} else if l.kind == PLAYER {
			logger.Service.Zap.Warnw("Bonus Bullet Shooted Warning",
				"GameUser", bonusBullet.SecWebSocketKey,
			)
		}
	case ActionUnprocessedBullet:
		secWebSocketKey := action.Payload()[0].(string)
		isDone := action.Payload()[1].(chan bool)

		if len(hitBulletTemp) == 0 && len(multiHitBulletTemp) == 0 {
			logger.Service.Zap.Infow("Processed Bullets Complete",
				"GameUser", secWebSocketKey,
			)
			isDone <- true
		} else {
			logger.Service.Zap.Infow("Have Unprocessed Bullets",
				"GameUser", secWebSocketKey,
				"HitBulletTemp", len(hitBulletTemp),
				"MultiHitBulletTemp", len(multiHitBulletTemp),
			)
			isDone <- false
		}
	}
}

func (l *lottery) sendResult(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability, times, mercenaryType int32) {
	data, dataBuffer := l.eResultRecall(hitFish, hitBullet, hitResult, times, mercenaryType)

	flux.Send(probability.EMSGID_eResultRecall, l.id(), hitBullet.SecWebSocketKey, dataBuffer)
	flux.Send(probability.EMSGID_eBroadcastResult, l.id(), hitFish.RoomUuid, l.eBroadcastResult(data))
}

func (l *lottery) sendMultiResult(hitFishies []*fish.Fish, hitBullet *bullet.Bullet, hitResult []*probability.Probability, times []int32, mercenaryType int32) {
	if (hitBullet.TypeId == 0 && len(hitFishies) == 5) || (hitBullet.TypeId == 999 && mercenaryType == 2) {
		_, _, data, dataBuffer := l.eMultiResultRecall(hitFishies, hitBullet, hitResult, times, mercenaryType)
		flux.Send(probability.EMSGID_eResultRecall, l.id(), hitBullet.SecWebSocketKey, dataBuffer)
		flux.Send(probability.EMSGID_eBroadcastResult, l.id(), hitFishies[0].RoomUuid, l.eBroadcastResult(data))
	} else {
		data, dataBuffer, _, _ := l.eMultiResultRecall(hitFishies, hitBullet, hitResult, times, mercenaryType)
		flux.Send(probability.EMSGID_eMultiResultRecall, l.id(), hitBullet.SecWebSocketKey, dataBuffer)
		flux.Send(probability.EMSGID_eBroadcastResult, l.id(), hitFishies[0].RoomUuid, l.eBroadcastMultiResult(data))
	}
}

func (l *lottery) eMultiResultRecall(hitFishies []*fish.Fish, hitBullet *bullet.Bullet, hitResult []*probability.Probability, times []int32, mercenaryType int32) (
	*probability_proto.MultiResultRecall, []byte, *probability_proto.ResultRecall, []byte,
) {
	bet := psf_on_00003_processWeapon(hitBullet.TypeId, mercenaryType, len(hitFishies), hitBullet.Bet)
	hitBullet.Bet = bet

	// Process Mercenary Bazooka and Player Bazooka
	if (hitBullet.TypeId == 0 && len(hitFishies) == 5) || (hitBullet.TypeId == 999 && mercenaryType == 2) {
		var iconPay uint64 = 0
		var bulletPay uint64 = 0

		var data *probability_proto.ResultRecall
		var triggerBonus []*probability_proto.TriggerBonus

		for i := 0; i < len(hitFishies); i++ {
			if hitResult != nil {
				iconPay += uint64(hitResult[i].Pay) * bet * hitBullet.Rate
				bulletPay += uint64(hitResult[i].Bullet)

				triggerBonus = append(triggerBonus, &probability_proto.TriggerBonus{
					Uuid:   hitResult[i].BonusUuid,
					TypeId: int32(hitResult[i].TriggerIconId),
					Times:  times[i],
				})
			}

			if i == len(hitFishies)-1 {
				if hitResult == nil {
					data, _ = l.eResultRecall(hitFishies[i], hitBullet, nil, times[i], mercenaryType)
				} else {
					data, _ = l.eResultRecall(hitFishies[i], hitBullet, hitResult[i], times[i], mercenaryType)
					data.HitResult.Pay = iconPay
					data.HitResult.BulletPay = bulletPay
					data.ExtraTriggerBonus = triggerBonus

					// Process Multi Fishes
					data.HitBullet.AccumulatedPay = hitBullet.ExtraData[3].(uint64) * bet * hitBullet.Rate
					data.HitResult.Pay = uint64(hitResult[i].Pay) * bet * hitBullet.Rate
				}
			}
		}

		dataBuffer, _ := proto.Marshal(data)
		return nil, nil, data, dataBuffer

	} else {
		multiResultRecall := &probability_proto.MultiResultRecall{
			Msgid: common_proto.EMSGID_eMultiResultRecall,
		}

		for i := 0; i < len(hitFishies); i++ {
			var data *probability_proto.ResultRecall

			if hitResult == nil {
				data, _ = l.eResultRecall(hitFishies[i], hitBullet, nil, times[i], mercenaryType)
			} else {
				data, _ = l.eResultRecall(hitFishies[i], hitBullet, hitResult[i], times[i], mercenaryType)
			}
			multiResultRecall.ResultRecall = append(multiResultRecall.ResultRecall, data)
		}

		dataBuffer, _ := proto.Marshal(multiResultRecall)
		return multiResultRecall, dataBuffer, nil, nil
	}
}

func (l *lottery) eResultRecall(hitFish *fish.Fish, hitBullet *bullet.Bullet, hitResult *probability.Probability, times, mercenaryType int32) (*probability_proto.ResultRecall, []byte) {
	secWebSocketKey := hitBullet.SecWebSocketKey
	subgameId := gamesetting.Service.SubgameId(secWebSocketKey)
	bulletCollection, mercenaryInfo := rtp.Service.MercenaryInfo(secWebSocketKey, hitFish.GameId, subgameId)

	data := &probability_proto.ResultRecall{
		Msgid:         common_proto.EMSGID_eResultRecall,
		RoomUuid:      hitFish.RoomUuid,
		SeatId:        hitFish.SeatId,
		MercenaryType: mercenaryType,
	}

	if hitResult != nil {
		data.StatusCode = common_proto.Status_kSuccess

		data.BetRateLine = &common_fish_proto.BetRateLine{
			BetLevelIndex: hitBullet.BetLevelIndex,
			BetIndex:      hitBullet.BetIndex,
			LineIndex:     hitBullet.LineIndex,
			RateIndex:     hitBullet.RateIndex,
		}

		data.HitFish = &probability_proto.HitFish{
			Uuid:     hitFish.FishUuid,
			SymbolId: hitFish.TypeId,
		}

		data.HitBullet = &probability_proto.HitBullet{
			BonusUuid:      hitBullet.BonusUuid,
			Uuid:           hitBullet.Uuid,
			TypeId:         hitBullet.TypeId,
			Target:         hitBullet.Target,
			ShootMode:      hitBullet.ShootMode,
			AccumulatedPay: hitBullet.ExtraData[3].(uint64) * hitBullet.Bet * hitBullet.Rate,
		}

		data.HitResult = &probability_proto.HitResult{
			Pay:        uint64(hitResult.Pay) * hitBullet.Bet * hitBullet.Rate,
			Multiplier: int32(hitResult.Multiplier),
			PlayerCent: hitBullet.PlayerCent,
		}

		switch hitFish.GameId {
		case models.PSF_ON_00003:
			data.TriggerBonus = &probability_proto.TriggerBonus{
				Uuid:   hitResult.BonusUuid,
				TypeId: int32(hitResult.TriggerIconId),
				Times:  times,
			}

			for _, v := range hitResult.ExtraTriggerBonus {
				data.ExtraTriggerBonus = append(data.ExtraTriggerBonus, &probability_proto.TriggerBonus{
					Uuid:   v.BonusUuid,
					TypeId: int32(v.TriggerIconId),
					Times:  times,
				})
			}

			// 傭兵的資訊
			data.HitResult.BulletPay = uint64(hitResult.Bullet)
			data.HitResult.BulletCollection = bulletCollection
			data.MercenaryInfo = &game_recovery_proto.Mercenary{
				MercenaryData: mercenaryInfo,
			}

		case models.PSF_ON_00004, models.PSF_ON_00006, models.PSF_ON_00007:
			var bonusPays []int32
			var bonusLocation []int32

			if len(hitResult.ExtraData) > 0 {
				if l.checkType_302(hitResult.ExtraData[0]) {
					bonusPays = l.convert302(hitResult.ExtraData[0].([]int), hitBullet.Bet, hitBullet.Rate)
				}
				if l.checkType_302(hitResult.ExtraData[1]) {
					bonusLocation = l.convert302Location(hitResult.ExtraData[1].([]int))
				}
			}

			data.TriggerBonus = &probability_proto.TriggerBonus{
				Uuid:   hitResult.BonusUuid,
				TypeId: int32(hitResult.TriggerIconId),
				Times:  times,
				BonusExtraInfo: &probability_proto.ExtraInfo{
					BonusPays:     bonusPays,
					BonusLocation: bonusLocation,
				},
			}

		default:
			data.TriggerBonus = &probability_proto.TriggerBonus{
				Uuid:   hitResult.BonusUuid,
				TypeId: int32(hitResult.TriggerIconId),
				Times:  times,
			}

			for _, v := range hitResult.ExtraTriggerBonus {
				data.ExtraTriggerBonus = append(data.ExtraTriggerBonus, &probability_proto.TriggerBonus{
					Uuid:   v.BonusUuid,
					TypeId: int32(v.TriggerIconId),
					Times:  times,
				})
			}
		}

		if l.kind != BOT {
			logger.Service.Zap.Infow("Send Result Recall",
				"GameUser", secWebSocketKey,
				"GameRoomUuid", data.RoomUuid,
				"FishType", data.HitFish.SymbolId,
				"FishUuid", data.HitFish.Uuid,
				"BulletType", data.HitBullet.TypeId,
				"BonusUuid", data.HitBullet.BonusUuid,
				"BulletUuid", data.HitBullet.Uuid,
				"AccumulatedPay", data.HitBullet.AccumulatedPay,
				"Cent", data.HitResult.Pay,
				"PlayerCent", data.HitResult.PlayerCent,
				"StatusCode", data.StatusCode,
				"TriggerBonus", data.TriggerBonus,
			)
		}
	} else {
		data.StatusCode = common_proto.Status_kInvalid

		data.HitFish = &probability_proto.HitFish{
			Uuid:     hitFish.FishUuid,
			SymbolId: hitFish.TypeId,
		}

		data.HitBullet = &probability_proto.HitBullet{
			BonusUuid:      hitBullet.BonusUuid,
			Uuid:           hitBullet.Uuid,
			TypeId:         hitBullet.TypeId,
			Target:         hitBullet.Target,
			ShootMode:      hitBullet.ShootMode,
			AccumulatedPay: 0,
		}

		if l.kind != BOT {
			logger.Service.Zap.Warnw(Lottery_RESULT_RECALL_INVALID,
				"GameUser", secWebSocketKey,
				"GameRoomUuid", data.RoomUuid,
				"FishType", hitFish.TypeId,
				"FishUuid", hitFish.FishUuid,
				"BulletType", hitBullet.TypeId,
				"BonusUuid", hitBullet.BonusUuid,
				"BulletUuid", hitBullet.Uuid,
				"StatusCode", data.StatusCode,
			)
		}
	}

	dataBuffer, _ := proto.Marshal(data)
	return data, dataBuffer
}

func (l *lottery) eBroadcastResult(resultRecall *probability_proto.ResultRecall) []byte {
	data := &probability_proto.BroadcastResult{
		Msgid:        common_proto.EMSGID_eBroadcastResult,
		ResultRecall: resultRecall,
	}

	dataByte, _ := proto.Marshal(data)
	return dataByte
}

func (l *lottery) eBroadcastMultiResult(resultMultiRecall *probability_proto.MultiResultRecall) []byte {
	data := &probability_proto.BroadcastMultiResult{
		Msgid:             common_proto.EMSGID_eBroadcastMultiResult,
		MultiResultRecall: resultMultiRecall,
	}

	dataByte, _ := proto.Marshal(data)
	return dataByte
}

func (l *lottery) id() string {
	return Service.HashId(l.secWebSocketKey)
}

func (l *lottery) destroy() {
	l.rwMutex.RLock()
	defer l.rwMutex.RUnlock()

	logger.Service.Zap.Infow("Lottery Unregister",
		"GameUser", l.secWebSocketKey,
		"Chan", fmt.Sprintf("%p", l.in),
		"GameId", l.gameId,
		"MathModuleId", l.mathModuleId,
	)
	flux.UnRegister(l.id(), l.in)
}

func (l *lottery) checkMultiFish(checkFishiesList []bool) bool {
	for i := 0; i < len(checkFishiesList); i++ {
		if checkFishiesList[i] == true {
			return false
		}
	}

	return true
}

func (l *lottery) checkType_302(input interface{}) bool {
	switch input.(type) {
	case []int:
		return true
	default:
		return false
	}
}

func (l *lottery) convert302Location(input []int) []int32 {
	var output []int32
	for _, v := range input {
		output = append(output, int32(v))
	}
	return output
}

func (l *lottery) convert302(input []int, bet, rate uint64) []int32 {
	var output []int32
	for _, v := range input {
		output = append(output, int32(v)*int32(bet)*int32(rate))
	}
	return output
}
