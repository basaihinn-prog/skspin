package models

import (
	"serve/servicejili/jili_144_cai/internal/message"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RawSpin struct {
	ID       primitive.ObjectID `bson:"_id"`
	Times    float64
	BucketId int
	Type     int
	HasGame  bool
	Selected bool

	Data                    *message.Cai_SpinAck
	HistoryRecord           *HistoryRecord
	SingleRoundLogSummaries []*SingleRoundLogSummary
	LogPlateInfos           map[string][]*LogPlateInfo
}
