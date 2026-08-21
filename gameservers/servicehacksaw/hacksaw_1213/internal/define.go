package internal

const (
	GameID   = "hacksaw_1213"
	FetchBet = 60
	MinBet   = 1

	// MOD_BONUS 购买模式
	MOD_BONUS = 6

	// FREESPINS_BOOK 购买模式
	FREESPINS_BOOK = 200

	// FREESPINS_BOOK_2 购买模式
	FREESPINS_BOOK_2 = 400

	// FREESPINS_BOOK_3 购买模式
	FREESPINS_BOOK_3 = 800

	// FREESPINS_CLOCK 购买模式
	FREESPINS_CLOCK = 400

	ModCount = 5
)

// GetFreeMultiply 兜底为1 返回常量数组
func GetFreeMultiply() []int {
	return []int{MinBet, MOD_BONUS, FREESPINS_BOOK, FREESPINS_BOOK_2, FREESPINS_BOOK_3, FREESPINS_CLOCK}
}

// FindBuy 找对应关系
func FindBuy(mod string) int {
	switch mod {
	case "mod_bonus":
		return 1
	case "freespins_book":
		return 2
	case "freespins_book2":
		return 3
	case "freespins_book3":
		return 4
	case "freespins_clock":
		return 5
	default:
		return 0
	}
}
