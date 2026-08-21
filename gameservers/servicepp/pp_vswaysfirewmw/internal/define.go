package internal

const (
	GameID        = "pp_vswaysfirewmw" //游戏ID
	Line          = 20                 //游戏固定倍数，1/20
	Double        = 1                  //开启了倍率下注，这个游戏不存在倍率下注，部分游戏是1.25
	BetMin        = 1                  //最小下注金额
	BetMax        = 3000               //最大下注金额
	MinBet        = 1                  //用于计算的值，该值是采集出来的
	DoubleBuy     = false              //是否允许双倍购买
	BuyBetMulti   = 1                  //购买模式，0为不能够买免费，1普通欧股买，2加倍购买，3超级购买
	FreeMultiply1 = 100
)

// GetFreeMultiply 返回常量数组
func GetFreeMultiply() []int {
	return []int{FreeMultiply1}
}

// 游戏部署好后，准备测试时下面全部改成false
const (
	NotCheckBigReward = false
	RandPlayResp      = false
)
