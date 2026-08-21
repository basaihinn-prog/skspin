package gamerecovery

import (
	"serve/service_fish/models"
)

const (
	GameRecovery_PROTO_INVALID           = models.GameRecovery + "0"
	GameRecovery_INSERT_GAME_DATA_FAILED = models.GameRecovery + "1"
	GameRecovery_QUERY_GAME_DATA_FAILED  = models.GameRecovery + "2"
	GameRecovery_SAVE_GAME_DATA_FAILED   = models.GameRecovery + "3"
	GameRecovery_DATA_NOT_FOUND          = models.GameRecovery + "4"
	GameRecovery_GAME_DB_NOT_FOUND       = models.GameRecovery + "5"
)
