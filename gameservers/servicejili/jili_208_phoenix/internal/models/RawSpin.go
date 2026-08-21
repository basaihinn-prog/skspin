package models

import (
	tmp "serve/servicejili/jili_208_phoenix/internal/message"
	"serve/servicejili/jiliut/jiliUtMessage"

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

	Phoenix_SpinAck         *tmp.Phoenix_SpinAck
	UtData                  *jiliUtMessage.Server_SpinResponse
	HistoryRecord           *HistoryRecord
	SingleRoundLogSummaries []*SingleRoundLogSummary
	LogPlateInfos           map[string][]*LogPlateInfo
}
