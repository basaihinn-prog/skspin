package internal

const (
	GameID   = "hacksaw_1233"
	FetchBet = 60
	MinBet   = 1

	// MOD_BONUS 购买模式
	MOD_BONUS = 6

	// MOD_2_EXPAND 购买模式
	MOD_2_EXPAND = 40

	// MOD_3_EXPAND 购买模式
	MOD_3_EXPAND = 100

	// FREESPINS 购买模式
	FREESPINS = 220

	// PROGRESSIVE_FS 购买模式
	PROGRESSIVE_FS = 400

	ModCount = 5
)

// GetFreeMultiply 兜底为1 返回常量数组
func GetFreeMultiply() []int {
	return []int{MinBet, MOD_BONUS, MOD_2_EXPAND, MOD_3_EXPAND, FREESPINS, PROGRESSIVE_FS}
}

// FindBuy 找对应关系
func FindBuy(mod string) int {
	switch mod {
	case "mod_bonus":
		return 1
	case "mod_2_expand":
		return 2
	case "mod_3_expand":
		return 3
	case "freespins":
		return 4
	case "progressive_fs":
		return 5
	default:
		return 0
	}
}
