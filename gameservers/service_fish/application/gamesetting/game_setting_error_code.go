package gamesetting

import (
	"serve/service_fish/models"
)

const (
	GameSetting_CONFIG_CALL_PROTO_INVALID       = models.GameSetting + "0"
	GameSetting_CONFIG_INVALID                  = models.GameSetting + "1"
	GameSetting_BET_ZERO                        = models.GameSetting + "2"
	GameSetting_RATE_ZERO                       = models.GameSetting + "3"
	GameSetting_GAME_ID_NOT_FOUND               = models.GameSetting + "4"
	GameSetting_HOST_ID_NOT_FOUND               = models.GameSetting + "5"
	GameSetting_RATE_INDEX_NOT_FOUND            = models.GameSetting + "6"
	GameSetting_MATH_MODULE_ID_NOT_FOUND        = models.GameSetting + "7"
	GameSetting_BET_LIST_NOT_FOUND              = models.GameSetting + "8"
	GameSetting_RATE_LIST_NOT_FOUND             = models.GameSetting + "9"
	GameSetting_BET_NOT_FOUND                   = models.GameSetting + "10"
	GameSetting_RATE_NOT_FOUND                  = models.GameSetting + "11"
	GameSetting_HOST_ID_EMPTY                   = models.GameSetting + "12"
	GameSetting_GAME_ID_EMPTY                   = models.GameSetting + "13"
	GameSetting_MATH_MODULE_ID_EMPTY            = models.GameSetting + "14"
	GameSetting_BET_LIST_EMPTY                  = models.GameSetting + "15"
	GameSetting_RATE_LIST_EMPTY                 = models.GameSetting + "16"
	GameSetting_STRIPS_CALL_PROTO_INVALID       = models.GameSetting + "17"
	GameSetting_DUPLICATE_FAILED                = models.GameSetting + "18"
	GameSetting_QUERY_HOST_FISH_GAME_FAILED     = models.GameSetting + "19"
	GameSetting_ACCOUNTING_PERIOD_NOT_FOUND     = models.GameSetting + "20"
	GameSetting_SELECTED_GAME_SETTING_NOT_FOUND = models.GameSetting + "21"
	GameSetting_GAME_SETTING_NOT_FOUND          = models.GameSetting + "22"
	GameSetting_SESSION_TIMEOUT_NOT_FOUND       = models.GameSetting + "23"
	GameSetting_GAME_DB_NOT_FOUND               = models.GameSetting + "24"
	GameSetting_SUBGAME_ID_NOT_FOUND            = models.GameSetting + "25"
	GameSetting_SUBGAME_ID_EMPTY                = models.GameSetting + "26"
	GameSetting_JACKPOT_GROUP_NOT_FOUND         = models.GameSetting + "27"
	GameSetting_JACKPOT_GROUP_EMPTY             = models.GameSetting + "28"
)
