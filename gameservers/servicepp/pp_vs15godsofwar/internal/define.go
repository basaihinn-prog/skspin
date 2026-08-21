package internal

import (
	"serve/comm/db"
	"serve/servicepp/ppcomm"
)

const (
	GameID        = "pp_vs15godsofwar" //游戏ID
	Line          = 10                 //游戏固定倍数，1/20
	Double        = 1                  //开启了倍率下注，这个游戏不存在倍率下注，部分游戏是1.25
	ZuseInd       = 1                  //开启了倍率下注，这个游戏不存在倍率下注，部分游戏是1.25
	HadesInd      = 0
	BetMin        = 1     //最小下注金额
	BetMax        = 3000  //最大下注金额
	MinBet        = 1     //用于计算的值，该值是采集出来的
	DoubleBuy     = false //是否允许双倍购买
	BuyBetMulti   = 1     //购买模式，0为不能够买免费，1普通欧股买，2加倍购买，3超级购买
	FreeMultiply1 = 150
	FreeMultiply2 = 300
	FreeMultiply3 = 75
	FreeMultiply4 = 300

	PurHadesBuy      = 0
	PurHadesSuperBuy = 1

	PurZuseBuy      = 2
	PurZuseSuperBuy = 3

	gameType = 0
)

// GetFreeMultiply 返回常量数组
func GetFreeMultiply() []float64 {
	return []float64{FreeMultiply1, FreeMultiply2, FreeMultiply3, FreeMultiply4}
}

func GetWitchGod() []int {
	return []int{HadesInd, ZuseInd}
}

func GodFindBucket(god, gameType int) db.BoundType {
	var t db.BoundType
	g := GetWitchGod()[god]
	if g == ZuseInd {
		if gameType == ppcomm.GameTypeNormal {
			t = 5
		} else if gameType == ppcomm.GameTypeGame {
			t = 6
		} else if gameType == ppcomm.GameTypeSuperGame1 {
			t = 7
		}
	} else if g == HadesInd {
		if gameType == ppcomm.GameTypeNormal {
			t = 0
		} else if gameType == ppcomm.GameTypeGame {
			t = 1
		} else if gameType == ppcomm.GameTypeSuperGame1 {
			t = 2
		}
	}
	return t
}

// 游戏部署好后，准备测试时下面全部改成false
const (
	NotCheckBigReward = false
	RandPlayResp      = false
)
