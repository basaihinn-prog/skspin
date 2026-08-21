package rtp

import (
	"serve/service_fish/models"
)

const (
	Rtp_RTP_ID_NOT_FOUND     = models.Rtp + "0"
	Rtp_RTP_STATE_NOT_FOUND  = models.Rtp + "1"
	Rtp_RTP_BUDGET_NOT_FOUND = models.Rtp + "2"

	Rtp_DECREASE_COLLECTION_FAILED       = models.Rtp + "3"
	Rtp_MERCENARY_NOT_FOUND              = models.Rtp + "4"
	Rtp_MERCENARY_TYPE_NOT_FOUND         = models.Rtp + "5"
	Rtp_DECREASE_MERCENARY_BULLET_FAILED = models.Rtp + "6"
)
