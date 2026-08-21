package internal

const (
	GameID = "pp_vs40bigjuan"
	Line   = 40
	BuyMul = 100
	Double = 1.25
	BetMin = 1
	BetMax = 3000
)

// 游戏部署好后，准备测试时下面全部改成false
const (
	NotCheckBigReward = false
	RandPlayResp      = false
)

// GetFreeMultiply 返回常量数组
func GetFreeMultiply() []int {
	return []int{BuyMul}
}
