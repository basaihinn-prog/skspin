package gamesetting

import game_setting_proto "serve/service_fish/application/gamesetting/proto"

type gameSetting struct {
	secWebSocketKey string
	hostExtId       string
	remoteAddr      string
	userAgent       string
	selectedGame    *selectedGame
	availableGames  map[string]*hostFishGame
}

type selectedGame struct {
	hostId           string
	gameId           string
	subgameId        int32
	mathModuleId     string
	betList          string      // for game room used
	rateList         string      // for game room used
	bets             [4][]uint32 // Change for PSF-ON-00003
	rate             []uint32
	langs            []string
	rateIndex        uint32
	accountingPeriod uint32
	playerInfo       uint32 // ID, NAME, BALANCE, AVATAR (1:show, 0:hidden) All show:15 (1111)
	sessionTimeout   uint32
	jackpotGroup     int32
}

type hostFishGame struct {
	HostId           string
	GameId           string
	SubGameId        int32
	MathId           string
	Langs            string
	BetList          string
	RateList         string
	AccountingPeriod uint32
	PlayerInfo       uint32 // ID, NAME, BALANCE, AVATAR (1:show, 0:hidden) All show:15 (1111)
	SessionTimeout   uint32
	JackpotGroup     int32
}

type symbol struct {
	symbolType game_setting_proto.StripsRecall_SymbolDef_SymbolType
	payType    game_setting_proto.StripsRecall_SymbolDef_PayType
}

func newGameSetting(secWebSocketKey, hostExtId, remoteAddr, userAgent string) *gameSetting {
	return &gameSetting{
		secWebSocketKey: secWebSocketKey,
		hostExtId:       hostExtId,
		remoteAddr:      remoteAddr,
		userAgent:       userAgent,
		selectedGame:    &selectedGame{},
		availableGames:  make(map[string]*hostFishGame),
	}
}
