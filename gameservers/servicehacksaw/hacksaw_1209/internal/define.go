package internal

const (
	GameID   = "hacksaw_1209"
	FetchBet = 60
	MinBet   = 1

	// MOD_BONUS 购买模式
	MOD_BONUS = 6

	// MOD_CLONER 购买模式
	MOD_CLONER = 40

	// FREESPINS 购买模式
	FREESPINS = 220

	// PROGRESSIVE_FS 购买模式
	PROGRESSIVE_FS = 500

	ModCount = 5
)

// GetFreeMultiply 兜底为1 返回常量数组
func GetFreeMultiply() []int {
	return []int{MinBet, MOD_BONUS, MOD_CLONER, FREESPINS, PROGRESSIVE_FS}
}

// FindBuy 找对应关系
func FindBuy(mod string) int {
	switch mod {
	case "mod_bonus":
		return 1
	case "mod_cloner":
		return 2
	case "freespins":
		return 3
	case "progressive_fs":
		return 4
	default:
		return 0
	}
}
