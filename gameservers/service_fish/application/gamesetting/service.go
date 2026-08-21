package gamesetting

import (
	"serve/fish_comm/mysql"
	game_setting_proto "serve/service_fish/application/gamesetting/proto"
	"serve/service_fish/domain/auth"
	auth_proto "serve/service_fish/domain/auth/proto"
	game_recovery_proto "serve/service_fish/domain/gamerecovery/proto"
	"serve/service_fish/models"
	"strconv"
	"strings"

	sync "github.com/sasha-s/go-deadlock"

	common_proto "serve/fish_comm/common/proto"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/idle"
	"serve/fish_comm/logger"
	"serve/fish_comm/session"
	sessionmanager "serve/fish_comm/session-manager"
	"serve/fish_comm/wallet"

	"github.com/gogo/protobuf/proto"
)

const (
	EMSGID_eGameConfigCall   = "EMSGID_eGameConfigCall"
	EMSGID_eGameConfigRecall = "EMSGID_eGameConfigRecall"
	EMSGID_eGameStripsCall   = "EMSGID_eGameStripsCall"
	EMSGID_eGameStripsRecall = "EMSGID_eGameStripsRecall"
)

var Service = &service{
	Id:              "GameSettingService",
	allUsersSetting: make(map[string]*gameSetting),
	rwMutex:         sync.RWMutex{},
}

type service struct {
	Id              string
	allUsersSetting map[string]*gameSetting
	rwMutex         sync.RWMutex
}

func (s *service) EMSGID_eGameConfigCall(secWebSocketKey, memberId string, guestEnable int, msg []byte) {
	s.rwMutex.RLock()

	data := &game_setting_proto.ConfigCall{}

	if err := proto.Unmarshal(msg, data); err != nil {
		logger.Service.Zap.Errorw(GameSetting_CONFIG_CALL_PROTO_INVALID,
			"GameUser", secWebSocketKey,
			"MemberId", memberId,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_CONFIG_CALL_PROTO_INVALID)
		return
	}

	isTokenValid := make(chan bool, 1024)
	flux.Send(auth.ActionCheckToken, s.Id, auth.Service.Id, isTokenValid, data.Token)

	if v, ok := s.allUsersSetting[secWebSocketKey]; ok {
		s.rwMutex.RUnlock()

		s.parseGameSetting(data.GameId, data.SubgameId, v)

		session.Service.New(
			secWebSocketKey,
			v.hostExtId,
			v.selectedGame.hostId,
			memberId,
			v.selectedGame.gameId,
			data.SubgameId,
			v.remoteAddr,
			v.userAgent,
			data.Version,
			data.GameVersion,
		)

		wallet.Service.SetGameIdAndSubGameId(secWebSocketKey, v.selectedGame.gameId, int(data.SubgameId))

		if guestEnable == 0 {
			go sessionmanager.Service.Start(
				secWebSocketKey,
				v.hostExtId,
				v.selectedGame.hostId,
				memberId,
				v.selectedGame.gameId,
			)
		}

		flux.Send(EMSGID_eGameConfigRecall, s.Id, secWebSocketKey,
			s.eGameConfigRecall(secWebSocketKey,
				data.SubgameId,
				<-isTokenValid,
				v,
			),
		)

		flux.Send(idle.ActionIdleTimeoutStart, s.Id, idle.Service.Id,
			secWebSocketKey,
			v.availableGames[data.GameId].SessionTimeout,
		)
	} else {
		s.rwMutex.RUnlock()

		logger.Service.Zap.Errorw(GameSetting_GAME_SETTING_NOT_FOUND,
			"GameUser", secWebSocketKey,
			"MemberId", memberId,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_GAME_SETTING_NOT_FOUND)
	}
}

func (s *service) parseGameSetting(gameId string, subGameId uint32, st *gameSetting) {
	if game, ok := st.availableGames[gameId]; ok {
		st.selectedGame.gameId = game.GameId
		st.selectedGame.subgameId = int32(subGameId)
		st.selectedGame.hostId = game.HostId
		st.selectedGame.mathModuleId = game.MathId
		st.selectedGame.betList = game.BetList
		st.selectedGame.rateList = game.RateList
		st.selectedGame.rateIndex = subGameId
		st.selectedGame.accountingPeriod = game.AccountingPeriod
		st.selectedGame.playerInfo = game.PlayerInfo
		st.selectedGame.sessionTimeout = game.SessionTimeout
		st.selectedGame.rate = nil
		st.selectedGame.langs = nil
		st.selectedGame.jackpotGroup = game.JackpotGroup

		roomsBetList := strings.Split(game.BetList, "|")
		betLevels := strings.Split(roomsBetList[subGameId], "/")

		switch gameId {
		case models.PSF_ON_00003:
			for k, bet := range betLevels {
				atoiBet, _ := strconv.Atoi(bet)
				st.selectedGame.bets[k] = nil
				st.selectedGame.bets[k] = append(st.selectedGame.bets[k], uint32(atoiBet))
			}

		default:
			for k, bets := range betLevels {
				st.selectedGame.bets[k] = nil

				for _, bet := range strings.Split(bets, ",") {
					atoiBet, _ := strconv.Atoi(bet)
					st.selectedGame.bets[k] = append(st.selectedGame.bets[k], uint32(atoiBet))
				}
			}
		}

		rates := strings.Split(game.RateList, ",")
		atoiRate, _ := strconv.Atoi(rates[subGameId])
		st.selectedGame.rate = []uint32{uint32(atoiRate)}

		for _, v := range strings.Split(game.Langs, ",") {
			st.selectedGame.langs = append(st.selectedGame.langs, v)
		}
	} else {
		logger.Service.Zap.Errorw(GameSetting_GAME_ID_NOT_FOUND,
			"GameUser", st.secWebSocketKey,
			"GameId", gameId,
			"SubGameId", subGameId,
		)
		errorcode.Service.Fatal(st.secWebSocketKey, GameSetting_GAME_ID_NOT_FOUND)
	}
}

func (s *service) eGameConfigRecall(secWebSocketKey string, subGameId uint32, isTokenValid bool, v *gameSetting) *game_setting_proto.ConfigRecall {
	data := &game_setting_proto.ConfigRecall{
		Msgid:                 common_proto.EMSGID_eGameConfigRecall,
		RatesDefaultIndex:     subGameId,
		LanguagesDefaultIndex: 0,
	}

	if !isTokenValid {
		s.cleanData(secWebSocketKey, auth.Auth_TOKEN_EXPIRED)
		data.StatusCode = common_proto.Status_kInvalid
		return data
	}

	if !wallet.Service.Existed(secWebSocketKey) {
		s.cleanData(secWebSocketKey, GameSetting_CONFIG_INVALID)
		data.StatusCode = common_proto.Status_kInvalid
		return data
	}

	go session.Service.Start(secWebSocketKey)

	data.StatusCode = common_proto.Status_kSuccess

	data.Languages = v.selectedGame.langs

	switch v.selectedGame.gameId {
	// PSF-ON-00003 Bet have 4 type
	case models.PSF_ON_00003:
		for _, v := range v.selectedGame.bets {
			b := &game_setting_proto.ConfigRecall_Bets{}

			for _, vv := range v {
				b.Bet = append(b.Bet, vv)
			}

			data.Levels = append(data.Levels, b)
		}

	default:
		for i := 0; i < len(v.selectedGame.bets)-1; i++ {
			b := &game_setting_proto.ConfigRecall_Bets{}

			for _, vv := range v.selectedGame.bets[i] {
				b.Bet = append(b.Bet, vv)
			}

			data.Levels = append(data.Levels, b)
		}
	}

	data.Line = []uint32{1}

	data.Rates = v.selectedGame.rate

	balance, err := wallet.Service.Balance(secWebSocketKey)

	if err != nil {
		logger.Service.Zap.Errorw("BUG",
			"GameUser", secWebSocketKey,
			"Error", err,
		)
		//panic(secWebSocketKey + err.Error())
	}

	data.PlayerInfo = &auth_proto.UserInfo{
		PlayerId:   wallet.Service.MemberId(secWebSocketKey),
		PlayerCent: balance,
		PlayerName: wallet.Service.MemberName(secWebSocketKey),
		Visible:    auth.Service.ColumnVisible(v.selectedGame.playerInfo),
	}

	data.MercenaryInfo = &game_recovery_proto.Mercenary{
		BulletCollection: 0,
		MercenaryData:    make(map[uint64]uint64),
	}

	return data
}

func (s *service) EMSGID_eGameStripsCall(secWebSocketKey string, msg []byte) {
	data := &game_setting_proto.StripsCall{}

	if err := proto.Unmarshal(msg, data); err != nil {
		s.cleanData(secWebSocketKey, GameSetting_STRIPS_CALL_PROTO_INVALID)
		return
	}

	isTokenValid := make(chan bool, 1024)
	flux.Send(auth.ActionCheckToken, s.Id, auth.Service.Id, isTokenValid, data.Token)

	flux.Send(EMSGID_eGameStripsRecall, s.Id, secWebSocketKey, s.eGameStripsRecall(secWebSocketKey, <-isTokenValid))
}

func (s *service) eGameStripsRecall(secWebSocketKey string, isTokenValid bool) *game_setting_proto.StripsRecall {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	data := &game_setting_proto.StripsRecall{
		Msgid:          common_proto.EMSGID_eGameStripsRecall,
		SymbolMaxLayer: 29,
	}

	if !isTokenValid {
		s.cleanData(secWebSocketKey, auth.Auth_TOKEN_EXPIRED)
		data.StatusCode = common_proto.Status_kInvalid
		return data
	}

	if v, ok := s.allUsersSetting[secWebSocketKey]; ok {
		data.StatusCode = common_proto.Status_kSuccess

		checkMath := true
		switch v.selectedGame.gameId {
		case models.PSF_ON_00001, models.PGF_ON_00001:
			if v.selectedGame.mathModuleId == models.PSFM_00001_98_1 {
				checkMath = psfm_00001_1_PayTable(data, v.selectedGame.mathModuleId)
			} else {
				checkMath = psfm_00002_1_PayTable(data, v.selectedGame.mathModuleId)
			}
		case models.PSF_ON_00002, models.PSF_ON_20002:
			checkMath = psfm_00003_1_PayTable(data, v.selectedGame.mathModuleId)
		case models.PSF_ON_00003:
			checkMath = psfm_00004_1_PayTable(data, v.selectedGame.mathModuleId)
		case models.PSF_ON_00004:
			checkMath = psfm_00005_1_PayTable(data, v.selectedGame.mathModuleId)
		case models.PSF_ON_00005:
			checkMath = psfm_00006_1_PayTable(data, v.selectedGame.mathModuleId)
		case models.PSF_ON_00006:
			checkMath = psfm_00007_1_PayTable(data, v.selectedGame.mathModuleId)
		case models.PSF_ON_00007:
			checkMath = psfm_00008_1_PayTable(data, v.selectedGame.mathModuleId)
		case models.RKF_H5_00001:
			checkMath = psfm_00013_1_PayTable(data, v.selectedGame.mathModuleId)
		}

		if !checkMath {
			s.cleanData(secWebSocketKey, GameSetting_MATH_MODULE_ID_NOT_FOUND)
			data.StatusCode = common_proto.Status_kInvalid
		}
	} else {
		s.cleanData(secWebSocketKey, GameSetting_SELECTED_GAME_SETTING_NOT_FOUND)
		data.StatusCode = common_proto.Status_kInvalid
	}

	return data
}

func (s *service) cleanData(secWebSocketKey, errorCode string) {
	logger.Service.Zap.Errorw(errorCode,
		"GameUser", secWebSocketKey,
	)
	errorcode.Service.Fatal(secWebSocketKey, errorCode)
}

func (s *service) Clone(secWebSocketKey, uuid string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		s.allUsersSetting[uuid] = data
	} else {
		logger.Service.Zap.Errorw(GameSetting_DUPLICATE_FAILED,
			"GameUser", secWebSocketKey,
			"UUID", uuid,
		)
		errorcode.Service.Fatal(uuid, GameSetting_DUPLICATE_FAILED)
	}
}

func (s *service) Bet(secWebSocketKey, gameId string, betLevelIndex, betIndex int32) (bet uint64) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if v, ok := s.allUsersSetting[secWebSocketKey]; ok {
		betLevelSize := int32(len(v.selectedGame.bets))

		if betLevelIndex >= betLevelSize {
			betLevelIndex = betLevelSize - 1
		}

		if betLevelIndex < 0 {
			betLevelIndex = 0
		}

		betSize := int32(len(v.selectedGame.bets[betLevelIndex]))

		if betIndex >= betSize {
			betIndex = betSize - 1
		}

		if betIndex < 0 {
			betIndex = 0
		}

		switch gameId {
		case models.PSF_ON_00003:
			bet = uint64(v.selectedGame.bets[betLevelIndex][0])
		default:
			bet = uint64(v.selectedGame.bets[betLevelIndex][betIndex])
		}
	} else {
		logger.Service.Zap.Errorw(GameSetting_BET_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_BET_NOT_FOUND)
		bet = 0
	}

	if bet == 0 {
		logger.Service.Zap.Errorw(GameSetting_BET_ZERO,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_BET_ZERO)
	}
	return bet
}

func (s *service) MaxBet(secWebSocketKey string) (bet uint64) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if v, ok := s.allUsersSetting[secWebSocketKey]; ok {
		betLevelIndex := len(v.selectedGame.bets) - 2
		betIndex := len(v.selectedGame.bets[betLevelIndex]) - 1
		bet = uint64(v.selectedGame.bets[betLevelIndex][betIndex])
	} else {
		logger.Service.Zap.Errorw(GameSetting_BET_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_BET_NOT_FOUND)
		bet = 0
	}

	if bet == 0 {
		logger.Service.Zap.Errorw(GameSetting_BET_ZERO,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_BET_ZERO)
	}
	return bet
}

func (s *service) Rate(secWebSocketKey string) (rate uint64) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if v, ok := s.allUsersSetting[secWebSocketKey]; ok {
		rate = uint64(v.selectedGame.rate[0])
	} else {
		logger.Service.Zap.Errorw(GameSetting_RATE_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_RATE_NOT_FOUND)
		rate = 0
	}

	if rate == 0 {
		logger.Service.Zap.Errorw(GameSetting_RATE_ZERO,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_RATE_ZERO)
	}
	return rate
}

func (s *service) HostId(secWebSocketKey string) (hostId string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		hostId = data.selectedGame.hostId
	} else {
		logger.Service.Zap.Errorw(GameSetting_HOST_ID_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_HOST_ID_NOT_FOUND)
		hostId = ""
	}

	if hostId == "" {
		logger.Service.Zap.Errorw(GameSetting_HOST_ID_EMPTY,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_HOST_ID_EMPTY)
	}
	return hostId
}

func (s *service) GameId(secWebSocketKey string) (gameId string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		gameId = data.selectedGame.gameId
	} else {
		logger.Service.Zap.Errorw(GameSetting_GAME_ID_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_GAME_ID_NOT_FOUND)
		gameId = ""
	}

	if gameId == "" {
		logger.Service.Zap.Errorw(GameSetting_GAME_ID_EMPTY,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_GAME_ID_EMPTY)
	}
	return gameId
}

func (s *service) RateIndex(secWebSocketKey string) (rateIndex uint32) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		rateIndex = data.selectedGame.rateIndex
	} else {
		logger.Service.Zap.Errorw(GameSetting_RATE_INDEX_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_RATE_INDEX_NOT_FOUND)
		rateIndex = 0
	}
	return rateIndex
}

func (s *service) MathModuleId(secWebSocketKey string) (mathModuleId string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		mathModuleId = data.selectedGame.mathModuleId
	} else {
		logger.Service.Zap.Errorw(GameSetting_MATH_MODULE_ID_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_MATH_MODULE_ID_NOT_FOUND)
		mathModuleId = ""
	}

	if mathModuleId == "" {
		logger.Service.Zap.Errorw(GameSetting_MATH_MODULE_ID_EMPTY,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_MATH_MODULE_ID_EMPTY)
	}
	return mathModuleId
}

func (s *service) BetList(secWebSocketKey string) (betList string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		betList = data.selectedGame.betList
	} else {
		logger.Service.Zap.Errorw(GameSetting_BET_LIST_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_BET_LIST_NOT_FOUND)
		betList = ""
	}

	if betList == "" {
		logger.Service.Zap.Errorw(GameSetting_BET_LIST_EMPTY,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_BET_LIST_EMPTY)
	}
	return betList
}

func (s *service) RateList(secWebSocketKey string) (rateList string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		rateList = data.selectedGame.rateList
	} else {
		logger.Service.Zap.Errorw(GameSetting_RATE_LIST_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_RATE_LIST_NOT_FOUND)
		rateList = ""
	}

	if rateList == "" {
		logger.Service.Zap.Errorw(GameSetting_RATE_LIST_EMPTY,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_RATE_LIST_EMPTY)
	}
	return rateList
}

func (s *service) AccountingPeriod(secWebSocketKey string) (accountingPeriod uint32) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		accountingPeriod = data.selectedGame.accountingPeriod
	} else {
		logger.Service.Zap.Errorw(GameSetting_ACCOUNTING_PERIOD_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_ACCOUNTING_PERIOD_NOT_FOUND)
		accountingPeriod = 0
	}
	return accountingPeriod
}

func (s *service) Idle(secWebSocketKey string) (sessionTimeout uint32) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		sessionTimeout = data.selectedGame.sessionTimeout
	} else {
		logger.Service.Zap.Errorw(GameSetting_SESSION_TIMEOUT_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_SESSION_TIMEOUT_NOT_FOUND)
		sessionTimeout = 0
	}
	return sessionTimeout
}

func (s *service) Delete(secWebSocketKey string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	delete(s.allUsersSetting, secWebSocketKey)
}

func (s *service) SubgameId(secWebSocketKey string) (subgameId int) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		subgameId = int(data.selectedGame.subgameId)
	} else {
		logger.Service.Zap.Errorw(GameSetting_SUBGAME_ID_NOT_FOUND,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_SUBGAME_ID_NOT_FOUND)
		subgameId = -1
	}

	if subgameId < 0 {
		logger.Service.Zap.Errorw(GameSetting_SUBGAME_ID_EMPTY,
			"GameUser", secWebSocketKey,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_SUBGAME_ID_EMPTY)
	}
	return subgameId
}

func (s *service) JackpotGroup(secWebSocketKey string) (jackpotGroup int32) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if data, ok := s.allUsersSetting[secWebSocketKey]; ok {
		jackpotGroup = data.selectedGame.jackpotGroup
	} else {
		logger.Service.Zap.Errorw(GameSetting_JACKPOT_GROUP_NOT_FOUND, "GameUser", secWebSocketKey)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_JACKPOT_GROUP_NOT_FOUND)
		jackpotGroup = -1
	}

	if jackpotGroup < 0 {
		logger.Service.Zap.Errorw(GameSetting_JACKPOT_GROUP_EMPTY, "GameUser", secWebSocketKey)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_JACKPOT_GROUP_EMPTY)
	}

	return jackpotGroup
}

func (s *service) New(secWebSocketKey, hostExtId, remoteAddr, userAgent string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	st := newGameSetting(secWebSocketKey, hostExtId, remoteAddr, userAgent)

	if db, err := mysql.Repository.GameDB(hostExtId); err == nil {
		if rows, err := db.
			Table("game").
			Select("IF(host_fish_game.math_id <> '',host_fish_game.math_id,game.math_id) math_id, game.langs, host_fish_game.host_id, host_fish_game.game_id, host_fish_game.subgame_id, host_fish_game.bet_list, host_fish_game.rate_list, host_fish_game.accounting_period, host_fish_game.player_info, host_id.session_timeout, jackpot_group").
			Joins("INNER JOIN host_fish_game ON game.enable = 1 AND host_fish_game.enable = 1 AND game.id = host_fish_game.game_id").
			Joins("INNER JOIN host_id ON host_id.host_id = host_fish_game.host_id AND host_id.host_ext_id = ?", hostExtId).
			Rows(); err == nil {

			defer rows.Close()

			for rows.Next() {
				game := &hostFishGame{}

				if err := rows.Scan(
					&game.MathId,
					&game.Langs,
					&game.HostId,
					&game.GameId,
					&game.SubGameId,
					&game.BetList,
					&game.RateList,
					&game.AccountingPeriod,
					&game.PlayerInfo,
					&game.SessionTimeout,
					&game.JackpotGroup,
				); err == nil {
					logger.Service.Zap.Infow("read host fish game settings",
						"GameUser", secWebSocketKey,
						"HostExtId", hostExtId,
						"HostId", game.HostId,
						"GameId", game.GameId,
						"MathModuleId", game.MathId,
						"SubGameId", game.SubGameId,
						"BetList", game.BetList,
						"RateList", game.RateList,
						"Languages", game.Langs,
						"AccountingPeriod", game.AccountingPeriod,
						"PlayerInfo", game.PlayerInfo,
						"SessionTimeout", game.SessionTimeout,
						"JackpotGroup", game.JackpotGroup,
					)
					st.availableGames[game.GameId] = game
				} else {
					logger.Service.Zap.Errorw(GameSetting_QUERY_HOST_FISH_GAME_FAILED,
						"GameUser", secWebSocketKey,
						"HostExtId", hostExtId,
						"Error", err,
					)
					errorcode.Service.Fatal(secWebSocketKey, GameSetting_QUERY_HOST_FISH_GAME_FAILED)
					return
				}
			}
			s.allUsersSetting[secWebSocketKey] = st
		} else {
			logger.Service.Zap.Errorw(GameSetting_QUERY_HOST_FISH_GAME_FAILED,
				"GameUser", secWebSocketKey,
				"HostExtId", hostExtId,
				"Error", err,
			)
			errorcode.Service.Fatal(secWebSocketKey, GameSetting_QUERY_HOST_FISH_GAME_FAILED)
		}
	} else {
		logger.Service.Zap.Errorw(GameSetting_GAME_DB_NOT_FOUND,
			"GameUser", secWebSocketKey,
			"HostExtId", hostExtId,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameSetting_GAME_DB_NOT_FOUND)
	}
}
