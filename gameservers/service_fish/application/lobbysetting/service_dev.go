//go:build dev
// +build dev

package lobbysetting

import (
	"serve/service_fish/models"
)

const (
	DEFAULT_HOST_ID           = "DEV"
	DEFAULT_RATE_LIST         = "10,100,1000"
	DEFAULT_BET_LIST          = "1,2,3,4,5,6,7/8,9,10,20,30,40/50,60,70,80,90,100|1,2,3,4,5,6,7/8,9,10,20,30,40/50,60,70,80,90,100|1,2,3,4/5,6,7/8,9,10"
	DEFAULT_LANGS             = "en-US,zh-CN,zh-TW"
	DEFAULT_ACCOUNTING_PERIOD = 1
	DEFAULT_PALYERINFO        = 15
	DEFAULT_SESSION_TIME_OUT  = 1200

	ZOMBIE_BET_LIST  = "1/3/5/10|1/3/5/10|1/3/5/10|1/3/5/10|1/3/5/10|1/3/5/10"
	ZOMBIE_RATE_LIST = "10,100,500,1000,5000,10000"
)

func (s *service) New(secWebSocketKey, hostExtId, remoteAddr, userAgent string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	st := newLobbySetting(secWebSocketKey, hostExtId, remoteAddr, userAgent)
	s.setDefaultAvailableGame(models.PSF_ON_00001, models.PSFM_00002_98_1, st)
	s.setDefaultAvailableGame(models.PSF_ON_00002, models.PSFM_00003_98_1, st)
	s.setBetAndRateAvailableGame(models.PSF_ON_00003, models.PSFM_00004_97_1, ZOMBIE_RATE_LIST, ZOMBIE_BET_LIST, st)
	s.setDefaultAvailableGame(models.PSF_ON_00004, models.PSFM_00005_97_1, st)
	s.setDefaultAvailableGame(models.PSF_ON_00005, models.PSFM_00006_98_1, st)
	s.setDefaultAvailableGame(models.PSF_ON_00006, models.PSFM_00007_97_1, st)
	s.setDefaultAvailableGame(models.PSF_ON_00007, models.PSFM_00008_97_1, st)

	s.setDefaultAvailableGame(models.PGF_ON_00001, models.PSFM_00002_98_1, st)

	s.setDefaultAvailableGame(models.PSF_ON_20002, models.PSFM_00003_93_1, st)

	s.setDefaultAvailableGame(models.RKF_H5_00001, models.PSFM_00013_97_1, st)

	s.allGamesSetting[secWebSocketKey] = st
}

func (s *service) setBetAndRateAvailableGame(gameId, mathId, rate, bet string, st *lobbySetting) {
	game := &hostFishGame{
		MathId:           mathId,
		GameId:           gameId,
		SubGameId:        0,
		HostId:           DEFAULT_HOST_ID,
		RateList:         rate,
		BetList:          bet,
		Langs:            DEFAULT_LANGS,
		AccountingPeriod: DEFAULT_ACCOUNTING_PERIOD,
		PlayerInfo:       DEFAULT_PALYERINFO,
		SessionTimeout:   DEFAULT_SESSION_TIME_OUT,
		JackpotGroup:     1,
	}

	st.availableGames[game.GameId] = game
}

func (s *service) setDefaultAvailableGame(gameId, mathId string, st *lobbySetting) {
	game := &hostFishGame{
		MathId:           mathId,
		GameId:           gameId,
		SubGameId:        0,
		HostId:           DEFAULT_HOST_ID,
		RateList:         DEFAULT_RATE_LIST,
		BetList:          DEFAULT_BET_LIST,
		Langs:            DEFAULT_LANGS,
		AccountingPeriod: DEFAULT_ACCOUNTING_PERIOD,
		PlayerInfo:       DEFAULT_PALYERINFO,
		SessionTimeout:   DEFAULT_SESSION_TIME_OUT,
	}

	st.availableGames[game.GameId] = game
}
