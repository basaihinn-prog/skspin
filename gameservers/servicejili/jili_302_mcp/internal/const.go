package internal

const (
	GameID        = "jili_302_mcp"
	GameShortName = "mcp"
	GameNo        = 302
)

const (
	GameTypeBet1    = 0 // 投注是1
	GameTypeBet5    = 1 // 投注是5
	GameTypeBet10   = 2 // 投注是10
	GameTypeBet50   = 3 // 投注是50
	GameTypeBet100  = 4 // 投注是100
	GameTypeBet500  = 5 // 投注是100
	GameTypeBet1000 = 6 // 投注是100
)

var BetMap = map[float64]int{
	1:    GameTypeBet1,
	5:    GameTypeBet5,
	10:   GameTypeBet10,
	50:   GameTypeBet50,
	100:  GameTypeBet100,
	500:  GameTypeBet500,
	1000: GameTypeBet1000,
}
