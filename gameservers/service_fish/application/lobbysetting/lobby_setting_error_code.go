package lobbysetting

import "serve/service_fish/models"

const (
	LobbySetting_CONFIG_INVALID              = models.LobbySetting + "0"
	LobbySetting_CONFIG_CALL_PROTO_INVALID   = models.LobbySetting + "1"
	LobbySetting_QUERY_HOST_FISH_GAME_FAILED = models.LobbySetting + "2"
	LobbySetting_GAME_SETTING_NOT_FOUND      = models.LobbySetting + "3"
	LobbySetting_GAME_DB_NOT_FOUND           = models.LobbySetting + "4"
	LobbySetting_GAME_ID_NOT_FOUND           = models.LobbySetting + "5"
)
