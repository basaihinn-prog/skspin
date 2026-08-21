package lobbysetting

import (
	"serve/fish_comm/broadcaster"
	common_proto "serve/fish_comm/common/proto"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/idle"
	"serve/fish_comm/logger"
	"serve/fish_comm/mysql"
	"serve/fish_comm/session"
	sessionmanager "serve/fish_comm/session-manager"
	"serve/fish_comm/wallet"
	lobby_setting_proto "serve/service_fish/application/lobbysetting/proto"
	"serve/service_fish/domain/auth"
	auth_proto "serve/service_fish/domain/auth/proto"
	"serve/service_fish/models"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	sync "github.com/sasha-s/go-deadlock"
)

const (
	EMSGID_eLobbyConfigCall   = "EMSGID_eLobbyConfigCall"
	EMSGID_eLobbyConfigRecall = "EMSGID_eLobbyConfigRecall"
)

var Service = &service{
	Id:              "LobbySettingService",
	allGamesSetting: make(map[string]*lobbySetting),
	rwMutex:         sync.RWMutex{},
}

type service struct {
	Id              string
	allGamesSetting map[string]*lobbySetting
	rwMutex         sync.RWMutex
}

func (s *service) EMSGID_eLobbyConfigCall(secWebSocketKey, memberId string, guestEnable int, msg []byte) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	data := &lobby_setting_proto.ConfigCall{}

	if err := proto.Unmarshal(msg, data); err != nil {
		logger.Service.Zap.Errorw(LobbySetting_CONFIG_CALL_PROTO_INVALID,
			"GameUser", secWebSocketKey,
			"MemberId", memberId,
		)
		errorcode.Service.Fatal(secWebSocketKey, LobbySetting_CONFIG_CALL_PROTO_INVALID)
		return
	}

	isTokenValid := make(chan bool, 1024)
	flux.Send(auth.ActionCheckToken, s.Id, auth.Service.Id, isTokenValid, data.Token)

	if v, ok := s.allGamesSetting[secWebSocketKey]; ok {
		if _, ok := v.availableGames[data.GameId]; !ok {
			logger.Service.Zap.Errorw(LobbySetting_GAME_ID_NOT_FOUND,
				"GameUser", secWebSocketKey,
				"MemberId", memberId,
				"GameId", data.GameId,
			)
			errorcode.Service.Fatal(secWebSocketKey, LobbySetting_GAME_ID_NOT_FOUND)
		} else {
			session.Service.New(
				secWebSocketKey,
				v.hostExtId,
				v.availableGames[data.GameId].HostId,
				memberId,
				data.GameId,
				data.SubgameId,
				v.remoteAddr,
				v.userAgent,
				data.Version,
				data.GameVersion,
			)

			// regular player
			if guestEnable == 0 {
				go sessionmanager.Service.Start(
					secWebSocketKey,
					v.hostExtId,
					v.availableGames[data.GameId].HostId,
					memberId,
					data.GameId,
				)
			}

			flux.Send(EMSGID_eLobbyConfigRecall, s.Id, secWebSocketKey,
				s.eLobbyConfigRecall(secWebSocketKey, data.GameId, <-isTokenValid),
				v.availableGames[data.GameId].JackpotGroup,
			)

			flux.Send(idle.ActionIdleTimeoutStart, s.Id, idle.Service.Id,
				secWebSocketKey,
				v.availableGames[data.GameId].SessionTimeout,
			)
		}
	} else {
		logger.Service.Zap.Errorw(LobbySetting_GAME_SETTING_NOT_FOUND,
			"GameUser", secWebSocketKey,
			"MemberId", memberId,
		)
		errorcode.Service.Fatal(secWebSocketKey, LobbySetting_GAME_SETTING_NOT_FOUND)
	}
}

func (s *service) eLobbyConfigRecall(secWebSocketKey, gameId string, isTokenValid bool) *lobby_setting_proto.ConfigRecall {
	data := &lobby_setting_proto.ConfigRecall{
		Msgid: common_proto.EMSGID_eLobbyConfigRecall,
	}

	if !isTokenValid {
		s.cleanData(secWebSocketKey, auth.Auth_TOKEN_EXPIRED)
		data.StatusCode = common_proto.Status_kInvalid
		return data
	}

	if !wallet.Service.Existed(secWebSocketKey) {
		s.cleanData(secWebSocketKey, LobbySetting_CONFIG_INVALID)
		data.StatusCode = common_proto.Status_kInvalid
		return data
	}

	if v, ok := s.allGamesSetting[secWebSocketKey]; ok {
		go session.Service.Start(secWebSocketKey)

		data.StatusCode = common_proto.Status_kSuccess

		games := v.availableGames

		for _, game := range games {
			if gameId == game.GameId {
				for k, rate := range strings.Split(game.RateList, ",") {
					atoiRate, _ := strconv.Atoi(rate)

					r := &lobby_setting_proto.ConfigRecall_RoomInfo{
						GameId:                game.GameId,
						SubgameId:             uint32(k),
						Rates:                 []uint32{uint32(atoiRate)},
						Line:                  []uint32{1},
						LanguagesDefaultIndex: 0,
					}
					s.parseGameSetting(game, r)
					data.RoomInfo = append(data.RoomInfo, r)
				}
			}
		}

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
			Visible:    auth.Service.ColumnVisible(v.availableGames[gameId].PlayerInfo),
		}
	} else {
		s.cleanData(secWebSocketKey, LobbySetting_GAME_SETTING_NOT_FOUND)
		data.StatusCode = common_proto.Status_kInvalid
	}
	return data
}

func (s *service) parseGameSetting(game *hostFishGame, r *lobby_setting_proto.ConfigRecall_RoomInfo) {
	roomsBetList := strings.Split(game.BetList, "|")

	switch game.GameId {
	case models.PSF_ON_00003:
		b := &lobby_setting_proto.ConfigRecall_Bets{}
		for betsIndex, bet := range strings.Split(roomsBetList[r.SubgameId], "/") {
			atoiBet, _ := strconv.Atoi(bet)
			b.Bet = append(b.Bet, uint32(atoiBet))

			if betsIndex == 0 && r.SubgameId == 1 {
				flux.Send(
					broadcaster.ActionBroadcasterConfig,
					s.Id,
					broadcaster.Service.Id,
					broadcaster.NewConfigBroadcaster(
						game.GameId,
						game.BetList,
						game.RateList,
						game.MathId,
						uint64(atoiBet),
						uint64(r.Rates[0]),
					),
				)
			}
		}
		r.Levels = append(r.Levels, b)

	default:
		betLevels := strings.Split(roomsBetList[r.SubgameId], "/")

		for betLevelsIndex, bets := range betLevels {
			b := &lobby_setting_proto.ConfigRecall_Bets{}

			for betsIndex, bet := range strings.Split(bets, ",") {
				atoiBet, _ := strconv.Atoi(bet)
				b.Bet = append(b.Bet, uint32(atoiBet))

				if betLevelsIndex == 0 && betsIndex == 0 && r.SubgameId == 1 {
					flux.Send(
						broadcaster.ActionBroadcasterConfig,
						s.Id,
						broadcaster.Service.Id,
						broadcaster.NewConfigBroadcaster(
							game.GameId,
							game.BetList,
							game.RateList,
							game.MathId,
							uint64(atoiBet),
							uint64(r.Rates[0]),
						),
					)
				}
			}
			r.Levels = append(r.Levels, b)
		}
	}

	for _, v := range strings.Split(game.Langs, ",") {
		r.Languages = append(r.Languages, v)
	}
}

func (s *service) MinBetAndRate(secWebSocketKey, gameId string) (bet, rate int) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if v, ok := s.allGamesSetting[secWebSocketKey]; ok {
		games := v.availableGames
		betList := games[gameId].BetList
		betIndex := strings.Split(betList, "|")[0]
		bet, _ = strconv.Atoi(strings.Split(betIndex, ",")[0])

		rateList := games[gameId].RateList
		rate, _ = strconv.Atoi(strings.Split(rateList, ",")[0])

		return bet, rate
	}

	return 0, 0
}

func (s *service) cleanData(secWebSocketKey, errorCode string) {
	logger.Service.Zap.Errorw(errorCode,
		"GameUser", secWebSocketKey,
	)
	errorcode.Service.Fatal(secWebSocketKey, errorCode)
}

func (s *service) Delete(secWebSocketKey string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	delete(s.allGamesSetting, secWebSocketKey)
}

func (s *service) GetPurposeByGameId(gameId string) int {
	switch gameId {
	case models.PSF_ON_00001:
		return models.Purpose_PSF_ON_00001
	case models.PSF_ON_00002:
		return models.Purpose_PSF_ON_00002
	case models.PSF_ON_00003:
		return models.Purpose_PSF_ON_00003
	case models.PSF_ON_00004:
		return models.Purpose_PSF_ON_00004
	case models.PSF_ON_00005:
		return models.Purpose_PSF_ON_00005
	case models.PSF_ON_00007:
		return models.Purpose_PSF_ON_00007
	case models.PSF_ON_20002:
		return models.Purpose_PSF_ON_20002
	case models.RKF_H5_00001:
		return models.Purpose_RKF_H5_00001
	}

	return 1
}
func (s *service) New(secWebSocketKey, hostExtId, remoteAddr, userAgent string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	st := newLobbySetting(secWebSocketKey, hostExtId, remoteAddr, userAgent)

	if db, err := mysql.Repository.GameDB(hostExtId); err == nil {
		if rows, err := db.
			Table("game").
			Select("IF(host_fish_game.math_id <> '',host_fish_game.math_id,game.math_id) math_id, game.langs, host_fish_game.host_id, host_fish_game.game_id, host_fish_game.subgame_id, host_fish_game.bet_list, host_fish_game.rate_list, host_fish_game.accounting_period, host_fish_game.player_info ,host_id.session_timeout, jackpot_group").
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
					)
					st.availableGames[game.GameId] = game
				} else {
					logger.Service.Zap.Errorw(LobbySetting_QUERY_HOST_FISH_GAME_FAILED,
						"GameUser", secWebSocketKey,
						"HostExtId", hostExtId,
						"Error", err,
					)
					errorcode.Service.Fatal(secWebSocketKey, LobbySetting_QUERY_HOST_FISH_GAME_FAILED)
					return
				}
			}
			s.allGamesSetting[secWebSocketKey] = st
		} else {
			logger.Service.Zap.Errorw(LobbySetting_QUERY_HOST_FISH_GAME_FAILED,
				"GameUser", secWebSocketKey,
				"HostExtId", hostExtId,
				"Error", err,
			)
			errorcode.Service.Fatal(secWebSocketKey, LobbySetting_QUERY_HOST_FISH_GAME_FAILED)
		}
	} else {
		logger.Service.Zap.Errorw(LobbySetting_GAME_DB_NOT_FOUND,
			"GameUser", secWebSocketKey,
			"HostExtId", hostExtId,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, LobbySetting_GAME_DB_NOT_FOUND)
	}
}
