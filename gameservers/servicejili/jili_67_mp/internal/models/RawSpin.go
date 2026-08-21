package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameType int

type RawSpin struct {
	ID                      primitive.ObjectID `bson:"_id"`
	Times                   float64
	BucketId                int
	Type                    int
	HasGame                 bool
	Selected                bool
	Data                    map[string]interface{}
	HistoryRecord           *HistoryRecord
	SingleRoundLogSummaries []*SingleRoundLogSummary
	LogPlateInfos           map[string][]*LogPlateInfo
}
