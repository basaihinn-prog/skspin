package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"serve/servicejili/jiliOfficialProto/serverOfficial"
	tmp "serve/servicejili/jili_17_sweetheart/internal/message"
)

type GameType int

type RawSpin struct {
	ID       primitive.ObjectID `bson:"_id"`
	Times    float64
	BucketId int
	Type     int
	HasGame  bool
	Selected bool

	SweetheartData          *tmp.AllPlate
	UtData                  *serverOfficial.SpinResponse
	HistoryRecord           *HistoryRecord
	SingleRoundLogSummaries []*SingleRoundLogSummary
	LogPlateInfos           map[string][]*LogPlateInfo
}
