package define

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PPBHItem struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Tid        string             `json:"tid" bson:"tid"`
	CC         string             `json:"cc" bson:"cc"`
	AgentCode  string             `json:"agentCode" bson:"agentCode"`
	UserCode   string             `json:"userCode" bson:"userCode"`
	GameCode   string             `json:"gameCode" bson:"gameCode"`
	RoundID    string             `json:"roundID" bson:"roundID"`
	Bet        float64            `json:"bet" bson:"bet"`
	Win        float64            `json:"win" bson:"win"`
	Rtp        float64            `json:"rtp" bson:"rtp"`
	PlayedDate int64              `json:"playedDate" bson:"playedDate"`
	Data       []*PPBDItem        `json:"data" bson:"data"`
	SharedLink string             `json:"sharedLink" bson:"sharedLink"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}

type PPBDItem struct {
	CR string `json:"cr" bson:"cr"`
	SR string `json:"sr" bson:"sr"`
}
