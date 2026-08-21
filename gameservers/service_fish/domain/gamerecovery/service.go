package gamerecovery

import (
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/logger"
	"serve/fish_comm/mysql"
	game_recovery_proto "serve/service_fish/domain/gamerecovery/proto"
	"serve/service_fish/models"

	"github.com/gogo/protobuf/proto"
)

var Service = &service{
	id:           "GameRecoveryService",
	gameRecovery: make(map[string]*gameRecovery),
}

type service struct {
	id           string
	gameRecovery map[string]*gameRecovery
}

func (s *service) New(secWebSocketKey, gameId string, subgameId int, hostId, hostExtId, memberId string) {
	if db, err := mysql.Repository.GameDB(hostExtId); err == nil {

		g := Builder().
			setGameId(gameId).
			setSubgameId(subgameId).
			setHostId(hostId).
			setHostExtId(hostExtId).
			setMemberId(memberId).
			build()

		if ok := db.Table("game_recovery").Create(g.recovery).RowsAffected; ok == 1 {
			g.hostExtId = hostExtId
			g.recovery.HostId = hostId
			g.recovery.GameId = gameId
			g.recovery.SubgameId = subgameId
			g.recovery.MemberId = memberId
			s.gameRecovery[secWebSocketKey] = g
			return
		}

		logger.Service.Zap.Errorw(GameRecovery_INSERT_GAME_DATA_FAILED,
			"GameUser", secWebSocketKey,
			"MemberId", memberId,
			"GameId", gameId,
			"HostId", hostId,
			"HostExtId", hostExtId,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameRecovery_INSERT_GAME_DATA_FAILED)
	} else {
		logger.Service.Zap.Errorw(GameRecovery_GAME_DB_NOT_FOUND,
			"GameUser", secWebSocketKey,
			"MemberId", memberId,
			"GameId", gameId,
			"HostId", hostId,
			"HostExtId", hostExtId,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameRecovery_GAME_DB_NOT_FOUND)
	}
}

func (s *service) Open(secWebSocketKey, gameId string, subgameId int, hostId, hostExtId, memberId string) *game_recovery_proto.GameRecovery {
	if db, err := mysql.Repository.GameDB(hostExtId); err == nil {

		dbResult := Builder().setHostExtId(hostExtId).build()

		switch gameId {
		case models.PSF_ON_00003:
			if ok := db.
				Table("game_recovery").
				Select("host_id, member_id, game_id, subgame_id, game_data, accounting_sn, game_end").
				Where("member_id = ? AND host_id = ? AND game_id = ? AND subgame_id = ?", memberId, hostId, gameId, subgameId).
				Scan(dbResult.recovery).
				RowsAffected; ok != 1 {

				logger.Service.Zap.Warnw(GameRecovery_QUERY_GAME_DATA_FAILED,
					"GameUser", secWebSocketKey,
					"MemberId", memberId,
					"GameId", gameId,
					"SubgameId", subgameId,
					"HostId", hostId,
					"HostExtId", hostExtId,
				)

				return nil
			}

		default:
			if ok := db.
				Table("game_recovery").
				Select("host_id, member_id, game_id, subgame_id, game_data, accounting_sn, game_end").
				Where("member_id = ? AND host_id = ? AND game_id = ?", memberId, hostId, gameId).
				Scan(dbResult.recovery).
				RowsAffected; ok != 1 {

				logger.Service.Zap.Warnw(GameRecovery_QUERY_GAME_DATA_FAILED,
					"GameUser", secWebSocketKey,
					"MemberId", memberId,
					"GameId", gameId,
					"SubgameId", subgameId,
					"HostId", hostId,
					"HostExtId", hostExtId,
				)

				return nil
			}
		}

		data := &game_recovery_proto.GameRecovery{}

		if err := proto.Unmarshal(dbResult.recovery.GameData, data); err != nil {
			logger.Service.Zap.Errorw(GameRecovery_PROTO_INVALID,
				"GameUser", secWebSocketKey,
				"MemberId", memberId,
				"GameId", gameId,
				"SubgameId", subgameId,
				"HostId", hostId,
				"HostExtId", hostExtId,
			)
			errorcode.Service.Fatal(secWebSocketKey, GameRecovery_PROTO_INVALID)
			return nil
		}
		s.gameRecovery[secWebSocketKey] = dbResult
		return data

	} else {
		logger.Service.Zap.Errorw(GameRecovery_GAME_DB_NOT_FOUND,
			"GameUser", secWebSocketKey,
			"MemberId", memberId,
			"GameId", gameId,
			"SubgameId", subgameId,
			"HostId", hostId,
			"HostExtId", hostExtId,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameRecovery_GAME_DB_NOT_FOUND)
		return nil
	}
}

func (s *service) Save(secWebSocketKey, rtpId string, rtpState int, rtpBudget, denominator,
	bulletCollection uint64, mercenaryData map[uint64]uint64) {
	if v, ok := s.gameRecovery[secWebSocketKey]; ok {
		if db, err := mysql.Repository.GameDB(v.hostExtId); err == nil {
			data := &game_recovery_proto.GameRecovery{
				Rtp: &game_recovery_proto.Rtp{
					Id:          rtpId,
					State:       int32(rtpState),
					Budget:      rtpBudget,
					Denominator: denominator,
				},
				Mercenary: &game_recovery_proto.Mercenary{
					BulletCollection: bulletCollection,
					MercenaryData:    mercenaryData,
				},
			}

			dataBuffer, _ := proto.Marshal(data)

			if ok := db.
				Table("game_recovery").
				Where("member_id = ? AND host_id = ? AND game_id = ? AND subgame_id = ?",
					v.recovery.MemberId, v.recovery.HostId, v.recovery.GameId, v.recovery.SubgameId).
				UpdateColumn("game_data", dataBuffer).
				RowsAffected; ok > 1 || ok < 0 {

				logger.Service.Zap.Errorw(GameRecovery_SAVE_GAME_DATA_FAILED,
					"GameUser", secWebSocketKey,
					"GameId", v.recovery.GameId,
					"SubgameId", v.recovery.SubgameId,
					"MemberId", v.recovery.MemberId,
					"HostId", v.recovery.HostId,
					"HostExtId", v.hostExtId,
				)
				errorcode.Service.Fatal(secWebSocketKey, GameRecovery_SAVE_GAME_DATA_FAILED)
			}
		} else {
			logger.Service.Zap.Errorw(GameRecovery_GAME_DB_NOT_FOUND,
				"GameUser", secWebSocketKey,
				"MemberId", v.recovery.MemberId,
				"GameId", v.recovery.GameId,
				"SubgameId", v.recovery.SubgameId,
				"HostId", v.recovery.HostId,
				"HostExtId", v.hostExtId,
				"Error", err,
			)
			errorcode.Service.Fatal(secWebSocketKey, GameRecovery_GAME_DB_NOT_FOUND)
		}
	} else {
		logger.Service.Zap.Errorw(GameRecovery_DATA_NOT_FOUND,
			"GameUser", secWebSocketKey,
			"RtpId", rtpId,
			"RtpState", rtpState,
			"RtpBudget", rtpBudget,
			"Denominator", denominator,
		)
		errorcode.Service.Fatal(secWebSocketKey, GameRecovery_DATA_NOT_FOUND)
	}
}

func (s *service) Close(secWebSocketKey string) {
	delete(s.gameRecovery, secWebSocketKey)
}
