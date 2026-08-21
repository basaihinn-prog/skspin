package internal

const (
	GameID   = "hacksaw_1067"
	FetchBet = 60
	MinBet   = 1

	// TRAIN 购买模式
	TRAIN = 160
	// DUEL 购买模式
	DUEL = 400
	// DEAD 购买模式
	DEAD = 800

	ModCount = 5
)

// GetFreeMultiply 兜底为1 返回常量数组
func GetFreeMultiply() []int {
	return []int{MinBet, TRAIN, DUEL, DEAD}
}

// FindBuy 找对应关系
func FindBuy(mod string) int {
	switch mod {
	case "freespins_train":
		return 1
	case "freespins_duel":
		return 2
	case "freespins_dead":
		return 3
	default:
		return 0
	}
}
