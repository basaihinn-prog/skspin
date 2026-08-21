package slot

type Slot struct {
	Uuid            string
	SecWebSocketKey string
	MathModuleId    string // not use yet
	RoomUuid        string
	FishTypeId      int32
	SeatId          int32
	BetIndex        int32
	BetLevelIndex   int32
	RateIndex       int32
	Line            int32
	Bet             uint64
	Rate            uint64
	Pay             uint64
	AllPay          []int64
	Reels           []int32
	ExtraData       []interface{}
}

func New(uuid, secWebSocketKey, mathModuleId, roomUuid string,
	seatId, fishTypeId, betIndex, betLevelIndex, line, rateIndex int32) *Slot {
	return &Slot{
		Uuid:            uuid,
		SecWebSocketKey: secWebSocketKey,
		MathModuleId:    mathModuleId,
		RoomUuid:        roomUuid,
		SeatId:          seatId,
		FishTypeId:      fishTypeId,
		BetIndex:        betIndex,
		BetLevelIndex:   betLevelIndex,
		RateIndex:       rateIndex,
		Line:            line,
		Pay:             0,
		ExtraData:       make([]interface{}, 3),
	}
}
