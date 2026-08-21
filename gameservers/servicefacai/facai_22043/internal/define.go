package internal

const (
	GameID = "facai_22043" //游戏ID
	Line   = 9             //游戏固定倍数，1/20

	BuyBetMulti   = 1 //购买模式，0为不能够买免费，1普通欧股买，2加倍购买，3超级购买
	FreeMultiply1 = 100
	FreeMultiply2 = 500
	ModCount      = 3
	originBetNor  = 40
	originBetEX   = 60
	originBetBuy  = 2000
	originBetBuy2 = 3000

	InitExMul    = 1.5
	InitBuyMul   = 50
	InitExBuyMul = 75
	InitBuyRange = 8
)

// GetFreeMultiply 返回常量数组
func GetFreeMultiply() []int {
	return []int{FreeMultiply1, FreeMultiply2}
}

// 游戏部署好后，准备测试时下面全部改成false
const (
	NotCheckBigReward = false
	RandPlayResp      = false
)

// 找对应关系
func FindBuy(mod int32) (int, float64) {
	switch mod {
	case 3:
		return 10, originBetEX
	case 4:
		return 1, originBetBuy
	case 6:
		return 2, originBetBuy2
	default:
		return 0, originBetNor
	}
}
