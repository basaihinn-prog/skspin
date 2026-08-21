package internal

const (
	GameID   = "hacksaw_1237"
	FetchBet = 60
	MinBet   = 1

	// MOD_BONUS 购买模式
	MOD_BONUS = 6

	// MOD_EXPAND 购买模式
	MOD_EXPAND = 100

	// FS_STICKY 购买模式
	FS_STICKY = 200

	// FS_MULT 购买模式
	FS_MULT = 400

	ModCount = 5
)

// GetFreeMultiply 兜底为1 返回常量数组
func GetFreeMultiply() []int {
	return []int{MinBet, MOD_BONUS, MOD_EXPAND, FS_STICKY, FS_MULT}
}

// FindBuy 找对应关系
func FindBuy(mod string) int {
	switch mod {
	case "mod_bonus":
		return 1
	case "mod_expand":
		return 2
	case "fs_sticky":
		return 3
	case "fs_mult":
		return 4
	default:
		return 0
	}
}
