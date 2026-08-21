package rtp

import (
	"serve/service_fish/domain/gamerecovery"
	"serve/service_fish/domain/probability"
	"serve/service_fish/models"
	"strconv"
	"sync"

	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/logger"
	"serve/fish_comm/wallet"
)

var Service = &service{
	Id:       "RtpService",
	usersRtp: make(map[string]*rtp),
	rwMutex:  sync.RWMutex{},
}

type service struct {
	Id       string
	usersRtp map[string]*rtp
	rwMutex  sync.RWMutex
}

func (s *service) Open(secWebSocketKey, gameId string, subgameId int, hostId, hostExtId, memberId string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	if _, ok := s.usersRtp[hashId]; ok {
		return
	}

	if v := gamerecovery.Service.Open(hashId, gameId, subgameId, hostId, hostExtId, memberId); v != nil {
		if v.Rtp != nil {
			switch gameId {
			case models.PSF_ON_00003:
				s.usersRtp[hashId] = Builder().
					setId(v.Rtp.Id).
					setState(int(v.Rtp.State)).
					setBudget(v.Rtp.Budget).
					setGuestEnable(wallet.Service.Guest(secWebSocketKey)).
					setDenominator(v.Rtp.Denominator).
					setBulletCollection(v.Mercenary.BulletCollection).
					setMercenaryInfo(v.Mercenary.MercenaryData).
					build()

			default:
				s.usersRtp[hashId] = Builder().
					setId(v.Rtp.Id).
					setState(int(v.Rtp.State)).
					setBudget(v.Rtp.Budget).
					setGuestEnable(wallet.Service.Guest(secWebSocketKey)).
					setDenominator(v.Rtp.Denominator).
					build()
			}

			logger.Service.Zap.Infow("Recovery RTP Budget",
				"GameUser", secWebSocketKey,
				"RtpId", v.Rtp.Id,
				"State", v.Rtp.State,
				"Budget", v.Rtp.Budget,
				"Denominator", v.Rtp.Denominator,
			)
		}
	} else {
		gamerecovery.Service.New(hashId, gameId, subgameId, hostId, hostExtId, memberId)

		logger.Service.Zap.Infow("Init RTP Budget",
			"GameUser", secWebSocketKey,
			"MemberId", memberId,
			"GameId", gameId,
			"HostId", hostId,
			"HostExtId", hostExtId,
		)
	}
}

func (s *service) Decrease(gameId string, subgameId int, mathModuleId, secWebSocketKey string, cent uint64) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	var budget uint64 = 0
	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	if v, ok := s.usersRtp[hashId]; ok {
		if cent > 0 && v.id != "" && v.state >= 0 {
			//if v.budget < cent {
			if v.budget <= 0 {

				rtpState, rtpId, rtpBullets, netWinGroup := probability.Service.Rtp(gameId, mathModuleId, v.state)

				if rtpBullets > 0 {
					budget = (rtpBullets - 1) * cent
				} else {
					budget = 0
				}

				v.id = rtpId
				v.state = rtpState
				v.budget = budget
				v.netWinGroup = netWinGroup

				logger.Service.Zap.Infow("Renew RTP Budget",
					"GameUser", secWebSocketKey,
					"GameId", gameId,
					"MathModuleId", mathModuleId,
					"RtpId", v.id,
					"State", v.state,
					"Bet", cent,
					"Bullets", rtpBullets,
					"Budget", v.budget,
				)
			} else {
				if v.budget < cent {
					v.budget = 0
				} else {
					v.budget -= cent
				}

				if !v.guestEnable {
					logger.Service.Zap.Infow("Decrease RTP Budget",
						"GameUser", secWebSocketKey,
						"GameId", gameId,
						"MathModuleId", mathModuleId,
						"RtpId", v.id,
						"State", v.state,
						"Bet", cent,
						"Budget", v.budget,
					)
				}
			}
			return
		}
	}

	rtpState, rtpId, rtpBullets, netWinGroup := probability.Service.Rtp(gameId, mathModuleId, -1)

	if rtpBullets > 0 {
		budget = (rtpBullets - 1) * cent
	} else {
		budget = 0
	}

	s.usersRtp[hashId] = Builder().
		setId(rtpId).
		setState(rtpState).
		setBudget(budget).
		setGuestEnable(wallet.Service.Guest(secWebSocketKey)).
		setMercenaryInfo(setMercenaryDataMap(3)).
		setNetWinGroup(netWinGroup).
		build()

	logger.Service.Zap.Infow("Init RTP Budget",
		"GameUser", secWebSocketKey,
		"GameId", gameId,
		"MathModuleId", mathModuleId,
		"RtpId", rtpId,
		"State", rtpState,
		"Bet", cent,
		"Budget", budget,
	)
}

func (s *service) Save(secWebSocketKey string, gameId string, subgameId int) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	if v, ok := s.usersRtp[hashId]; ok {
		if !v.guestEnable {
			gamerecovery.Service.Save(
				hashId, v.id, v.state, v.budget, v.denominator,
				v.bulletCollection, v.mercenaryInfo,
			)

			logger.Service.Zap.Infow("Save RTP Budget",
				"GameUser", secWebSocketKey,
				"GameId", gameId,
				"SubgameId", subgameId,
				"RtpId", v.id,
				"State", v.state,
				"Budget", v.budget,
				"Denominator", v.denominator,
			)
		}
	}
}

func (s *service) Close(secWebSocketKey, gameId string, subgameId int) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	if v, ok := s.usersRtp[hashId]; ok {
		if !v.guestEnable {
			gamerecovery.Service.Save(
				hashId, v.id, v.state, v.budget, v.denominator,
				v.bulletCollection, v.mercenaryInfo,
			)

			logger.Service.Zap.Infow("Save RTP Budget",
				"GameUser", secWebSocketKey,
				"RtpId", v.id,
				"State", v.state,
				"Budget", v.budget,
			)
		}
	}

	gamerecovery.Service.Close(hashId)
	delete(s.usersRtp, hashId)
}

func (s *service) RtpId(secWebSocketKey, gameId string, subgameId int) (rtpId string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	if r, ok := s.usersRtp[hashId]; ok {
		return r.id
	} else {
		logger.Service.Zap.Warnw(Rtp_RTP_ID_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		return "0"
	}
}

func (s *service) RtpState(secWebSocketKey, gameId string, subgameId int) (rtpState int) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)
	if r, ok := s.usersRtp[hashId]; ok {
		return r.state
	} else {
		logger.Service.Zap.Warnw(Rtp_RTP_STATE_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		return -1
	}
}

func (s *service) NetWinGroup(secWebSocketKey, gameId string, subgameId int) (netWinGroup int) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)
	if r, ok := s.usersRtp[hashId]; ok {
		return r.netWinGroup
	} else {
		return -1
	}
}

func (s *service) HashId(secWebSocketKey, gameId string, subgameId int) string {
	switch gameId {
	case models.PSF_ON_00003, models.PSF_ON_00004:
		return secWebSocketKey + gameId + strconv.Itoa(subgameId)

	default:
		return secWebSocketKey
	}
}

func (s *service) RtpBudget(secWebSocketKey, gameId string, subgameId int, mathModuleId string) uint64 {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	//if gameId != models.PSF_ON_00004 {
	//    return 0
	//}
	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	if r, ok := s.usersRtp[hashId]; ok {
		return r.budget
	} else {
		logger.Service.Zap.Warnw(Rtp_RTP_BUDGET_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)

		return 0
	}
}

func (s *service) DecreaseDenominator(gameId string, subgameId int, mathModuleId, secWebSocketKey string, fishId int32, rtpOder string, cent uint64, minBet int) (
	simulationDenominator, simulationMolecular int) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	// Set RtpId
	rtpId := probability.Service.Rtp_V2(gameId, mathModuleId, fishId, rtpOder)

	if v, ok := s.usersRtp[hashId]; ok {
		// 分母<=0
		if v.denominator <= 0 || v.denominator < cent {
			denominator, molecular, multiplier := probability.Service.Fraction(gameId, mathModuleId)
			v.denominator = v.denominator + (uint64(denominator) * cent)
			v.budget = v.budget + (uint64(molecular) * cent)

			v.denominator = v.denominator - cent
			v.id = rtpId

			logger.Service.Zap.Infow("Init RTP Denominator",
				"GameUser", secWebSocketKey,
				"GameId", gameId,
				"SubgameId", subgameId,
				"MathModuleId", mathModuleId,
				"RtpId", v.id,
				"Denominator", v.denominator,
				"Budget", v.budget,
			)

			simulationDenominator = denominator
			if denominator == 0 {
				simulationMolecular = -1
			} else {
				simulationMolecular = molecular * multiplier / denominator
			}

			return simulationDenominator, simulationMolecular
		} else {
			// 分母>0
			v.denominator = v.denominator - cent
			v.id = rtpId

			logger.Service.Zap.Infow("Decrease RTP Denominator",
				"GameUser", secWebSocketKey,
				"GameId", gameId,
				"SubgameId", subgameId,
				"MathModuleId", mathModuleId,
				"RtpId", v.id,
				"Denominator", v.denominator,
				"Budget", v.budget,
			)

			return 0, 0
		}
	}

	// First Time Play Game
	firstMultiplier := probability.Service.FirstMultiplier(gameId, subgameId, mathModuleId)
	denominator, molecular, multiplier := probability.Service.Fraction(gameId, mathModuleId)

	denominator = (denominator * int(cent)) - int(cent)
	initMolecular := uint64(firstMultiplier * minBet)
	budget := uint64(molecular) * cent

	s.usersRtp[hashId] = Builder().
		setId(rtpId).
		setBudget(budget + initMolecular).
		setDenominator(uint64(denominator)).
		setGuestEnable(wallet.Service.Guest(secWebSocketKey)).
		setBulletCollection(0).
		setMercenaryInfo(make(map[uint64]uint64)).
		build()

	logger.Service.Zap.Infow("Init RTP Denominator and Molecular",
		"GameUser", secWebSocketKey,
		"GameId", gameId,
		"SubgameId", subgameId,
		"MathModuleId", mathModuleId,
		"RtpId", rtpId,
		"Denominator", denominator,
		"Budget", budget+initMolecular,
		"Molecular", budget,
		"initMolecular", initMolecular,
	)

	simulationDenominator = denominator + int(cent)
	if denominator == 0 {
		simulationMolecular = -1
	} else {
		simulationMolecular = molecular * multiplier / denominator
	}

	return simulationDenominator, simulationMolecular
}

func (s *service) DecreaseMolecular(gameId string, subgameId int, mathModuleId, secWebSocketKey string,
	fishId int32, iconPay, bet int, triggerIconId int) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	if v, ok := s.usersRtp[hashId]; ok {
		// 處理有AvgPay的種類
		haveAvgPay, avgPay := probability.Service.AvgPay(gameId, mathModuleId, int(fishId))
		switch gameId {
		case models.PSF_ON_00004:
			iconPay = psf_on_00004_AvgPayCheck(fishId, triggerIconId, iconPay, avgPay, bet, haveAvgPay)
		}

		// 分母>0 魚倍率<=分子
		if v.denominator+uint64(bet) > 0 && uint64(iconPay) <= v.budget {
			v.budget = v.budget - uint64(iconPay)

			logger.Service.Zap.Infow("Decrease RTP Budget",
				"GameUser", secWebSocketKey,
				"GameId", gameId,
				"SubgameId", subgameId,
				"MathModuleId", mathModuleId,
				"Denominator", v.denominator,
				"Budget", v.budget,
			)
		}
	}
}

func (s *service) IncreaseBulletCollection(secWebSocketKey string, gameId string, subgameId int, bulletCollection uint64) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	switch gameId {
	case models.PSF_ON_00003:
		if v, ok := s.usersRtp[hashId]; ok {
			if bulletCollection > 0 {
				v.bulletCollection += bulletCollection

				logger.Service.Zap.Infow("Increase Mercenary Bullet Collection",
					"GameUser", secWebSocketKey,
					"GameId", gameId,
					"SubgameId", subgameId,
					"BulletCollection", v.bulletCollection,
				)
			}
		}
	}
}

func (s *service) MercenaryOpen(secWebSocketKey string, gameId string, subgameId int,
	mercenaryType, bullet uint64,
) (openResult bool) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)
	if v, ok := s.usersRtp[hashId]; ok {
		if v.bulletCollection >= bullet {
			// Check BulletCollection
			var mercenaryBullet uint64 = 0
			switch {
			case v.bulletCollection >= PSF_ON_00003_MERCENARY_BULLET_LEVEL_3:
				mercenaryBullet = PSF_ON_00003_MERCENARY_BULLET_LEVEL_3
			case v.bulletCollection >= PSF_ON_00003_MERCENARY_BULLET_LEVEL_2 && v.bulletCollection < PSF_ON_00003_MERCENARY_BULLET_LEVEL_3:
				mercenaryBullet = PSF_ON_00003_MERCENARY_BULLET_LEVEL_2
			case v.bulletCollection >= PSF_ON_00003_MERCENARY_BULLET_LEVEL_1 && v.bulletCollection < PSF_ON_00003_MERCENARY_BULLET_LEVEL_2:
				mercenaryBullet = PSF_ON_00003_MERCENARY_BULLET_LEVEL_1
			}
			v.bulletCollection -= bullet

			// Check Increase Mercenary Bullet
			switch mercenaryType {
			case PSF_ON_00003_MERCENARY_TYPE_RIFLE:
				v.mercenaryInfo[mercenaryType] += mercenaryBullet
			case PSF_ON_00003_MERCENARY_TYPE_SHOTGUN:
				v.mercenaryInfo[mercenaryType] += mercenaryBullet / PSF_ON_00003_MERCENARY_SHOTGUN
			case PSF_ON_00003_MERCENARY_TYPE_BAZOOKA:
				v.mercenaryInfo[mercenaryType] += mercenaryBullet / PSF_ON_00003_MERCENARY_BAZOOKA
			}

			logger.Service.Zap.Infow("Open Mercenary",
				"GameUser", secWebSocketKey,
				"GameId", gameId,
				"SubGameId", subgameId,
				"MercenaryCollection", v.bulletCollection,
				"MercenaryType", mercenaryType,
				"MercenaryBullets", v.mercenaryInfo[mercenaryType],
			)

			return true
		} else {
			logger.Service.Zap.Warnw(Rtp_DECREASE_COLLECTION_FAILED,
				"GameUser", secWebSocketKey,
				"GameId", gameId,
				"SubGameId", subgameId,
				"MercenaryCollection", v.bulletCollection,
				"MercenaryType", mercenaryType,
			)
			return false
		}
	}

	return false
}

func (s *service) GetMercenary(secWebSocketKey string, gameId string, subgameId int, mercenaryType uint64,
) (bulletCollection uint64, mercenaryId int32, mercenaryBullet uint64) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	if v, ok := s.usersRtp[hashId]; ok {
		if m, ok := v.mercenaryInfo[mercenaryType]; ok {
			return v.bulletCollection, int32(mercenaryType), m
		} else {
			logger.Service.Zap.Warnw(Rtp_MERCENARY_TYPE_NOT_FOUND,
				"GameUser", secWebSocketKey,
			)
			return 0, -1, 0
		}
	} else {
		logger.Service.Zap.Warnw(Rtp_MERCENARY_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)

		return 0, -1, 0
	}
}

func (s *service) MercenaryInfo(secWebSocketKey string, gameId string, subgameId int,
) (bulletCollection uint64, mercenaryInfo map[uint64]uint64) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	if v, ok := s.usersRtp[hashId]; ok {
		return v.bulletCollection, v.mercenaryInfo
	} else {
		logger.Service.Zap.Warnw(Rtp_MERCENARY_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)

		return 0, nil
	}
}

func (s *service) DecreaseMercenaryBullet(secWebSocketKey, gameId string, subgameId int, mercenaryType uint64) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	hashId := s.HashId(secWebSocketKey, gameId, subgameId)

	if v, ok := s.usersRtp[hashId]; ok {
		if v.mercenaryInfo[mercenaryType] > 0 {
			v.mercenaryInfo[mercenaryType]--

			logger.Service.Zap.Infow("Decrease Mercenary Bullet",
				"GameUser", secWebSocketKey,
				"GameId", gameId,
				"MercenaryType", mercenaryType,
				"MercenaryBullet", v.mercenaryInfo[mercenaryType],
			)
		} else {
			logger.Service.Zap.Warnw(Rtp_DECREASE_MERCENARY_BULLET_FAILED,
				"GameUser", secWebSocketKey,
				"MercenaryType", mercenaryType,
			)
			errorcode.Service.Fatal(secWebSocketKey, Rtp_DECREASE_MERCENARY_BULLET_FAILED)
		}
	} else {
		logger.Service.Zap.Warnw(Rtp_MERCENARY_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
	}
}
