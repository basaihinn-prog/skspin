package robot

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RobotStatus_Normal  = iota // 准备就绪
	RobotStatus_Working        // 游戏中
	RobotStatus_Rest           // 休息
)

type Robot struct {
	ID         primitive.ObjectID `bson:"_id"`
	Balance    int64              `bson:"balance"`
	Status     int64              `bson:"status"`
	EndTime    int64              `bson:"endTime"`
	GameId     string             `bson:"gameId"`     // 游戏id，游戏中才有效
	CreateTime int64              `bson:"createTime"` // 机器人创建的时间
	OverTime   int64              `bson:"overTime"`   // 机器人到期时间
}

type GameInfo struct {
	Id        string `bson:"_id" json:"id"`              // 游戏Id
	Weight    int64  `bson:"weight" json:"weight"`       // 拉取游戏权重
	IsEnd     bool   `bson:"isEnd" json:"isEnd"`         // 是否已结束
	StartTime int64  `bson:"startTime" json:"startTime"` // 开始时间

	EndTime     int64   `bson:"endTime" json:"endTime"` // 结束时间
	Line        int     `bson:"line" json:"line"`
	InsertC     float64 `bson:"insertC" json:"insertC"`
	NormalCount int     `bson:"normalCount" json:"normalCount"`
	BuyCount    int     `bson:"buyCount" json:"buyCount"`
	ArrC        string  `bson:"arrC" json:"arrC"`           //自行解析的数据
	Status      int     `bson:"status" json:"status"`       //状态 1 等待调度 2调度中
	SuperBuy    int     `bson:"superBuy" json:"superBuy"`   //状态 超级购买
	SuperBuy2   int     `bson:"superBuy2" json:"superBuy2"` //状态 超级购买
	APIVersion  string  `bson:"APIVersion" json:"APIVersion"`
	HuiDuSerial string  `bson:"HuiDuSerial" json:"HuiDuSerial"` //1430(大鲈鱼)

}

func (r *Robot) Run(fn func()) {
	fn()
}
