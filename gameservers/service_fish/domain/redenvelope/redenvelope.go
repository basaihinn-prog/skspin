package redenvelope

type RedEnvelope struct {
	Uuid            string
	SecWebSocketKey string
	MathModuleId    string // not use yet
	RoomUuid        string
	SeatId          int32
	FishTypeId      int32
	PlayerOptIndex  int32
	BetIndex        int32 // not use yet
	BetLevelIndex   int32 // not use yet
	RateIndex       int32 // not use yet
	Line            int32
	Bet             uint64
	Rate            uint64
	Pay             uint64
	AllPay          []int64
	BonusPayload    []int
	ExtraData       []interface{}
}

func New(uuid, secWebSocketKey, mathModuleId, roomUuid string, seatId, fishTypeId int32, bonusPayload []int, betIndex, betLevelIndex, line, rateIndex int32) *RedEnvelope {
	return &RedEnvelope{
		Uuid:            uuid,
		SecWebSocketKey: secWebSocketKey,
		MathModuleId:    mathModuleId,
		RoomUuid:        roomUuid,
		SeatId:          seatId,
		FishTypeId:      fishTypeId,
		PlayerOptIndex:  -1,
		BetIndex:        betIndex,
		BetLevelIndex:   betLevelIndex,
		RateIndex:       rateIndex,
		Line:            line,
		Pay:             0,
		AllPay:          make([]int64, 5),
		BonusPayload:    bonusPayload,
		ExtraData:       make([]interface{}, 3),
	}
}
