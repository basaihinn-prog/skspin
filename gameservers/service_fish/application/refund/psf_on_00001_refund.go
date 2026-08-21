package refund

import (
	"serve/fish_comm/flux"
	"serve/service_fish/domain/bullet"
	"serve/service_fish/domain/redenvelope"
	"serve/service_fish/domain/slot"
)

func psf_on_00001_init(action *flux.Action, allRefund map[string]*refund) {
	controllerId := action.Key().From()
	secWebSocketKey := action.Payload()[0].(string)
	h := action.Payload()[1].(IHandler)

	refundStatus := map[string]bool{
		bullet.Service.Id:      false,
		redenvelope.Service.Id: false,
		slot.Service.Id:        false,
	}

	allRefund[secWebSocketKey] = newRefund(controllerId, secWebSocketKey, refundStatus, h)
}
