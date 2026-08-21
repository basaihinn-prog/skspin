//go:build staging || prod
// +build staging prod

package lobbysetting

import (
	errorcode "serve/fish_comm/flux/error-code"
	"serve/fish_comm/flux/logger"
	"serve/fish_comm/flux/mysql"
)

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
