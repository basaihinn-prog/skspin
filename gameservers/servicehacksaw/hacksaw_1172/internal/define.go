package internal

const (
	GameID   = "hacksaw_1172"
	FetchBet = 60
	MinBet   = 1

	// MOD_BOXES 购买模式
	MOD_BOXES = 6

	// BONUS_BOXES 购买模式
	BONUS_BOXES = 200

	// BONUS_REELS 购买模式
	BONUS_REELS = 400

	ModCount = 5
)

// GetFreeMultiply 兜底为1 返回常量数组
func GetFreeMultiply() []int {
	return []int{MinBet, MOD_BOXES, BONUS_BOXES, BONUS_REELS}
}

// FindBuy 找对应关系
func FindBuy(mod string) int {
	switch mod {
	case "mod_boxes":
		return 1
	case "bonus_boxes":
		return 2
	case "bonus_reels":
		return 3
	default:
		return 0
	}
}
