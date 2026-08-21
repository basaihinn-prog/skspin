package internal

const (
	GameID   = "hacksaw_1243"
	FetchBet = 60
	MinBet   = 1

	// MOD_BONUS 购买模式
	MOD_BONUS = 6

	// MOD_2_EXPAND 购买模式
	MOD_2_EXPAND = 40

	// MOD_3_EXPAND 购买模式
	MOD_3_EXPAND = 100

	// FS_PORTAL 购买模式
	FS_PORTAL = 200

	// FS_PORTAL_2 购买模式
	FS_PORTAL_2 = 400

	// FS_REEL 购买模式
	FS_REEL = 400

	ModCount = 5
)

// GetFreeMultiply 兜底为1 返回常量数组
func GetFreeMultiply() []int {
	return []int{MinBet, MOD_BONUS, MOD_2_EXPAND, MOD_3_EXPAND, FS_PORTAL, FS_PORTAL_2, FS_REEL}
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
	case "fs_portal":
		return 4
	case "fs_portal_2":
		return 5
	case "fs_reel":
		return 6
	default:
		return 0
	}
}
