package internal

const (
	GameID   = "hacksaw_1168"
	FetchBet = 60
	MinBet   = 1

	// BONUS_1 购买模式
	BONUS_1 = 258
	// BONUS_2 购买模式
	BONUS_2 = 400

	ModCount = 5
)

// GetFreeMultiply 兜底为1 返回常量数组
func GetFreeMultiply() []int {
	return []int{MinBet, BONUS_1, BONUS_2}
}

// FindBuy 找对应关系
func FindBuy(mod string) int {
	switch mod {
	case "bonus_1":
		return 1
	case "bonus_2":
		return 2
	default:
		return 0
	}
}
