package models

import (
	"serve/servicejili/jili_193_df/internal/message"
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

	DfData                  *message.Df_SpinAck
	UtData                  *jiliUtMessage.Server_SpinResponse
	HistoryRecord           *HistoryRecord
	SingleRoundLogSummaries []*SingleRoundLogSummary
	LogPlateInfos           map[string][]*LogPlateInfo
}
