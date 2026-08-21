package internal

import (
	"context"
	"fmt"
	"serve/comm/db"
	"serve/servicepp/ppcomm"

	"github.com/nats-io/nats.go"
)

func init() {
	ppcomm.RegRpc("create_task", Create)
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

func Create(msg *nats.Msg) (ret []byte, err error) {
	ps := ppcomm.ParseVariables(string(msg.Data))
	var (
		game        = ps.Str("game")        // 游戏ID
		line        = ps.Int("line")        // 游戏line, 参照游戏发送doSpin请求里面的‘l’值
		insertC     = ps.Float("insertC")   // 插入数据库时的投注，参照游戏发送doSpin请求里面的‘c’值
		normalCount = ps.Int("normalCount") // 需要拉取普通盘的次数
		buyCount    = ps.Int("buyCount")    // 需要拉取游戏的次数
		// 使用最低3个挡位
		arrC       = ps.Str("arrC")
		gameWeight = int64(ps.Int("weight")) // 如果需要更多机器人拉数据，可以把这个值改大点，建议不修改
	)

	coll := db.Collection2("robot_hgc", "game")
	_, err = coll.InsertOne(context.TODO(), GameInfo{
		Id:          game,
		StartTime:   0,
		IsEnd:       false,
		Line:        line,
		InsertC:     insertC,
		NormalCount: normalCount,
		BuyCount:    buyCount,
		Weight:      gameWeight,
		Status:      2, //等待任务调度
		ArrC:        arrC,
		APIVersion:  ps.Str("APIVersion"),
		SuperBuy:    ps.Int("superBuy"),
		SuperBuy2:   ps.Int("superBuy2"),
		HuiDuSerial: ps.Str("HuiDuSerial"),
	})
	fmt.Println(err)
	return
}
