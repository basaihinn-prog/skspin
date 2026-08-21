package pgcomm

import "serve/comm/db"

type PGPlayer struct {
	db.DocPlayer `bson:"inline"`
	LS           string
	Cs           float64
	Ml           float64
	LastSid      string
	BdRecords    []map[string]any
	IsBuy        bool //是否是购买小游戏
	BigReward    int64
	FpChoose     []string //记录玩家小游戏的选择顺序
}
