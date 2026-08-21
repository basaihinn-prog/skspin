package models

import (
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	"serve/servicejili/jili_171_sc/internal/message"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameType int

type RawSpin struct {
	ID       primitive.ObjectID `bson:"_id"`
	Times    float64
	BucketId int
	Type     int
	HasGame  bool
	Selected bool

	ScData                  *message.SpinAck
	UtData                  *serverOfficial.SpinResponse
	HistoryRecord           *HistoryRecord
	SingleRoundLogSummaries []*SingleRoundLogSummary
	LogPlateInfos           map[string][]*LogPlateInfo
}
