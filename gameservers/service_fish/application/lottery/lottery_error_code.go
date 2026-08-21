package lottery

import (
	"serve/service_fish/models"
)

const (
	Lottery_RESULT_CALL_MAP_SIZE_INVALID = models.Lottery + "0"
	Lottery_RESULT_CALL_PROTO_INVALID    = models.Lottery + "1"
	Lottery_HIT_FISH_BULLET_NOT_ALLOWED  = models.Lottery + "2"
	Lottery_GAME_ID_INVALID              = models.Lottery + "3"
	Lottery_RESULT_RECALL_INVALID        = models.Lottery + "4"
	Lottery_WALLET_DECREASE_FAILED       = models.Lottery + "5"
	Lottery_HITFISH_NIL                  = models.Lottery + "6"
	Lottery_HITBULLET_NIL                = models.Lottery + "7"

	Lottery_MERCENARY_OPEN_RECALL_INVALID = models.Lottery + "8"
)
