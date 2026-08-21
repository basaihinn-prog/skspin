package internal

const (
	GameID        = "jili_51_mc"
	GameShortName = "mc"
	GameNo        = 51
)

const (
	GameTypeBet1   = 0 // 投注是1
	GameTypeBet5   = 1 // 投注是5
	GameTypeBet10  = 2 // 投注是10
	GameTypeBet50  = 3 // 投注是50
	GameTypeBet100 = 4 // 投注是100
)

var BetMap = map[float64]int{
	1:   GameTypeBet1,
	5:   GameTypeBet5,
	10:  GameTypeBet10,
	50:  GameTypeBet50,
	100: GameTypeBet100,
}
