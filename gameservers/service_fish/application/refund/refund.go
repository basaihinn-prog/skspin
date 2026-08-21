package refund

type refund struct {
	controllerId    string
	secWebSocketKey string
	gameRoomUuid    string // not used
	refundStatus    map[string]bool
	isDisconnect    bool
	iHandler        IHandler
}

type IHandler interface {
	RefundHandler(bonus interface{}, accountingSn uint64)
}

func newRefund(controllerId, secWebSocketKey string, refundStatus map[string]bool, h IHandler) *refund {
	return &refund{
		controllerId:    controllerId,
		secWebSocketKey: secWebSocketKey,
		refundStatus:    refundStatus,
		isDisconnect:    false,
		iHandler:        h,
	}
}

func (r *refund) reset() {
	for k := range r.refundStatus {
		r.refundStatus[k] = false
	}
}
