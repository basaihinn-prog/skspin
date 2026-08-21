package internal

const (
	GameID   = "hacksaw_1217"
	FetchBet = 60
	MinBet   = 1

	// MOD_BONUS 购买模式
	MOD_BONUS = 6

	// MOD_TREATS 购买模式
	MOD_TREATS = 4

	// BONUS_TREATS 购买模式
	BONUS_TREATS = 200

	// BONUS_TOTAL_BAR 购买模式
	BONUS_TOTAL_BAR = 400

	ModCount = 5
)

// GetFreeMultiply 兜底为1 返回常量数组
func GetFreeMultiply() []int {
	return []int{MinBet, MOD_BONUS, MOD_TREATS, BONUS_TREATS, BONUS_TOTAL_BAR}
}

// FindBuy 找对应关系
func FindBuy(mod string) int {
	switch mod {
	case "mod_bonus":
		return 1
	case "mod_treats":
		return 2
	case "bonus_treats":
		return 3
	case "bonus_total_bar":
		return 4
	default:
		return 0
	}
}
