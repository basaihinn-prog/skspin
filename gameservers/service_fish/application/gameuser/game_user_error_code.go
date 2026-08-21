package gameuser

import (
	"serve/service_fish/models"
)

const (
	GameUser_GAME_ROOM_UUID_EMPTY      = models.GameUser + "0"
	GameUser_MSGID_PROTO_INVALID       = models.GameUser + "1"
	GameUser_PING_FAILED               = models.GameUser + "2"
	GameUser_NOT_FOUND                 = models.GameUser + "3"
	GameUser_RESULT_CALL_PROTO_INVALID = models.GameUser + "4"
	GameUser_OPTION_CALL_PROTO_INVALID = models.GameUser + "5"
	GameUser_GET_BALANCE_FAILED        = models.GameUser + "6"
	GameUser_PROFILE_NIL               = models.GameUser + "7"
	GameUser_ACCOUNTING_SN_ZERO        = models.GameUser + "8"
)
