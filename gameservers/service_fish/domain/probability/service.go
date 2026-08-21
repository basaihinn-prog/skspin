package probability

import (
	"crypto/rand"
	"math/big"
	"serve/service_fish/models"

	"serve/fish_comm/logger"
)

const (
	EMSGID_eResultCall        = "EMSGID_eResultCall"
	EMSGID_eResultRecall      = "EMSGID_eResultRecall"
	EMSGID_eMultiResultCall   = "EMSGID_eMultiResultCall"
	EMSGID_eMultiResultRecall = "EMSGID_eMultiResultRecall"
	EMSGID_eBroadcastResult   = "EMSGID_eBroadcastResult"
	EMSGID_eOptionCall        = "EMSGID_eOptionCall"
	EMSGID_eOptionRecall      = "EMSGID_eOptionRecall"
	EMSGID_eBroadcastOption   = "EMSGID_eBroadcastOption"

	EMSGID_eMercenaryOpenCall   = "EMSGID_eMercenaryOpen"
	EMSGID_eMercenaryOpenRecall = "EMSGID_eMercenaryOpenRecall"
)

var Service = &service{
	id: "ProbabilityService",
}

type service struct {
	id string
}

func init() {
	logger.Service.Zap.Infow(
		"Service created.",
		"Service", Service.id,
	)
}

func (s *service) Calc(gameId, mathModuleId, rtpId string, fishId, bulletId int32, budget uint64, cent uint64) *Probability {
	switch gameId {
	case models.PSF_ON_00001, models.PGF_ON_00001:
		if mathModuleId == models.PSFM_00001_98_1 {
			bsMath, fsMath := chooseMath_psfm00001(mathModuleId)
			if bsMath != nil && fsMath != nil {
				return psfm_00001_1(fishId, bsMath, fsMath)
			}
		} else {
			bsMath, fsMath, _ := chooseMath_psfm00002(mathModuleId)
			if bsMath != nil && fsMath != nil {
				return psfm_00002_1(fishId, rtpId, bsMath, fsMath)
			}
		}

	case models.PSF_ON_00002, models.PSF_ON_20002:
		bsMath, fsMath, _ := chooseMath_psfm00003(mathModuleId)
		if bsMath != nil && fsMath != nil {
			return psfm_00003_1(fishId, rtpId, bsMath, fsMath)
		}

	case models.PSF_ON_00003:
		bsMath, fsMath, _ := chooseMath_psfm00004(mathModuleId)
		if bsMath != nil && fsMath != nil {
			return psfm_00004_1(fishId, rtpId, bsMath, fsMath)
		}

	case models.PSF_ON_00004:
		bsMath, fsMath, _ := chooseMath_psfm00005(mathModuleId)
		if bsMath != nil && fsMath != nil {
			return psfm_00005_1(fishId, bulletId, rtpId, budget, cent, bsMath, fsMath)
		}

	case models.PSF_ON_00005:
		bsMath, fsMath, _ := chooseMath_psfm00006(mathModuleId)
		if bsMath != nil && fsMath != nil {
			return psfm_00006_1(fishId, rtpId, bsMath, fsMath)
		}

	case models.PSF_ON_00006:
		bsMath, fsMath, _ := chooseMath_psfm00007(mathModuleId)
		if bsMath != nil && fsMath != nil {
			return psfm_00007_1(fishId, rtpId, bsMath, fsMath)
		}

	case models.PSF_ON_00007:
		bsMath, fsMath, _ := chooseMath_psfm00008(mathModuleId)
		if bsMath != nil && fsMath != nil {
			return psfm_00008_1(fishId, rtpId, bsMath, fsMath)
		}

	case models.RKF_H5_00001:
		bsMath, fsMath, _ := chooseMath_psfm00013(mathModuleId)
		if bsMath != nil && fsMath != nil {
			return psfm_00013_1(fishId, rtpId, bsMath, fsMath)
		}
	}

	return newSimplePay(int(fishId), 0, 1)
}

func (s *service) Rtp(gameId, mathModuleId string, state int) (rtpState int, rtpId string, rtpBullets uint64, netWinGroup int) {
	switch gameId {
	case models.PSF_ON_00001, models.PGF_ON_00001:
		if mathModuleId == models.PSFM_00001_98_1 {
			return -1, "0", 0, -1
		} else {
			_, _, drbMath := chooseMath_psfm00002(mathModuleId)
			if drbMath != nil {
				rtpState, rtpId, rtpBullets = psfm_00002_1_RngRtp(state, drbMath)
				return rtpState, rtpId, rtpBullets, -1
			}
		}

	case models.PSF_ON_00002, models.PSF_ON_20002:
		_, _, drbMath := chooseMath_psfm00003(mathModuleId)
		if drbMath != nil {
			rtpState, rtpId, rtpBullets = psfm_00003_1_RngRtp(state, drbMath)
			return rtpState, rtpId, rtpBullets, -1
		}

	case models.PSF_ON_00003:
		_, _, drbMath := chooseMath_psfm00004(mathModuleId)
		if drbMath != nil {
			return psfm_00004_1_RngRtp(state, drbMath)
		}

	case models.PSF_ON_00005:
		_, _, drbMath := chooseMath_psfm00006(mathModuleId)
		if drbMath != nil {
			rtpState, rtpId, rtpBullets = psfm_00006_1_RngRtp(state, drbMath)
			return rtpState, rtpId, rtpBullets, -1
		}

	case models.PSF_ON_00006:
		_, _, drbMath := chooseMath_psfm00007(mathModuleId)
		if drbMath != nil {
			rtpState, rtpId, rtpBullets = psfm_00007_1_RngRtp(state, drbMath)
			return rtpState, rtpId, rtpBullets, -1
		}

	case models.PSF_ON_00007:
		_, _, drbMath := chooseMath_psfm00008(mathModuleId)
		if drbMath != nil {
			rtpState, rtpId, rtpBullets = psfm_00008_1_RngRtp(state, drbMath)
			return rtpState, rtpId, rtpBullets, -1
		}

	case models.RKF_H5_00001:
		_, _, drbMath := chooseMath_psfm00013(mathModuleId)
		if drbMath != nil {
			rtpState, rtpId, rtpBullets = psfm_00013_1_RngRtp(state, drbMath)
			return rtpState, rtpId, rtpBullets, -1
		}
	}

	return -1, "0", 0, -1
}

func (s *service) Rtp_V2(gameId, mathModuleId string, fishId int32, rtpOrder string) (rtpId string) {
	switch gameId {
	case models.PSF_ON_00004:
		bsMath, _, _ := chooseMath_psfm00005(mathModuleId)
		if bsMath != nil {
			return psfm_00005_1_RngRtp(fishId, rtpOrder, bsMath)
		}
	}

	return ""
}

func (s *service) Fraction(gameId, mathModuleId string) (denominator, molecular, multiplier int) {
	switch gameId {
	case models.PSF_ON_00004:
		_, _, drbMath := chooseMath_psfm00005(mathModuleId)
		if drbMath != nil {
			return psfm_00005_1_RngFraction(drbMath)
		}
	}

	return 0, 0, 0
}

func (s *service) UseRtp(gameId, mathModuleId, rtpId string, bulletType, mercenaryType int32) string {
	switch gameId {
	case models.PSF_ON_00003:
		bsMath, _, _ := chooseMath_psfm00004(mathModuleId)
		if bsMath != nil {
			return psfm_00004_1_useRtp(bulletType, mercenaryType, rtpId, bsMath)
		}

	case models.PSF_ON_00004:
		bsMath, _, _ := chooseMath_psfm00005(mathModuleId)
		if bsMath != nil {
			return psfm_00005_1_useRtp(bulletType, rtpId, bsMath)
		}

	case models.PSF_ON_00005:
		bsMath, _, _ := chooseMath_psfm00006(mathModuleId)
		if bsMath != nil {
			return psfm_00006_1_useRtp(bulletType, rtpId, bsMath)
		}

	case models.PSF_ON_00006:
		bsMath, _, _ := chooseMath_psfm00007(mathModuleId)
		if bsMath != nil {
			return psfm_00007_1_useRtp(bulletType, rtpId, bsMath)
		}

	case models.PSF_ON_00007:
		bsMath, _, _ := chooseMath_psfm00008(mathModuleId)
		if bsMath != nil {
			return psfm_00008_1_useRtp(bulletType, rtpId, bsMath)
		}

	case models.RKF_H5_00001:
		bsMath, _, _ := chooseMath_psfm00013(mathModuleId)
		if bsMath != nil {
			return psfm_00013_1_useRtp(bulletType, rtpId, bsMath)
		}
	}

	return rtpId
}

func (s *service) AvgPay(gameId, mathModuleId string, fishId int) (haveAvgPay bool, avgPay int) {
	switch gameId {
	case models.PSF_ON_00004:
		bsMath, _, _ := chooseMath_psfm00005(mathModuleId)
		if bsMath != nil {
			return psfm_00005_1_avgPay(fishId, bsMath)
		}
	}

	return false, 0
}

func (s *service) FirstMultiplier(gameId string, subgameId int, mathModuleId string) int {
	switch gameId {
	case models.PSF_ON_00004:
		_, _, drbMath := chooseMath_psfm00005(mathModuleId)
		if drbMath != nil {
			return psfm_00005_1_firstMultiplier(subgameId, drbMath)
		}
	}

	return 0
}

func (s *service) GetNewSimplePay(fishId int) *Probability {
	return newSimplePay(fishId, 0, 1)
}

func (s *service) Random(value int64) uint64 {
	n, err := rand.Int(rand.Reader, big.NewInt(value))
	if err != nil {
		panic(err)
	}

	return n.Uint64()
}
