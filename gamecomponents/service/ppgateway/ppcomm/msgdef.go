package pgcomm

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//todo 修改为pp游戏通用字段,也可能用不上

type BHItem struct {
	ID    primitive.ObjectID `bson:"_id" json:"-"`
	Pid   int64              `json:"-" bson:"pid"`
	Tid   string             `json:"tid" bson:"tid"`
	Gid   int                `json:"gid" bson:"gid"`
	CC    string             `json:"cc" bson:"cc"`
	Gtba  float64            `json:"gtba" bson:"gtba"`
	Gtwla float64            `json:"gtwla" bson:"gtwla"`
	Bt    int64              `json:"bt" bson:"bt"`
	Ge    any                `json:"ge" bson:"ge"`
	Bd    []*BDItem          `json:"bd" bson:"bd"`
	Mgcc  int                `json:"mgcc" bson:"mgcc"`
	Fscc  int                `json:"fscc" bson:"fscc"`
}

type BDItem struct {
	Tid  string          `json:"tid" bson:"tid"`
	Tba  float64         `json:"tba" bson:"tba"`
	Twla float64         `json:"twla" bson:"twla"`
	Bl   float64         `json:"bl" bson:"bl"`
	Bt   int64           `json:"bt" bson:"bt"`
	Gd   json.RawMessage `json:"gd" bson:"gd"`
}
