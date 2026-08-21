package lobbysetting

type lobbySetting struct {
	secWebSocketKey string
	hostExtId       string
	remoteAddr      string
	userAgent       string
	availableGames  map[string]*hostFishGame
}

type hostFishGame struct {
	HostId           string
	GameId           string
	SubGameId        uint32
	MathId           string
	Langs            string
	BetList          string
	RateList         string
	AccountingPeriod uint32
	PlayerInfo       uint32 // ID, NAME, BALANCE, AVATAR (1:show, 0:hidden) All show:15 (1111)
	SessionTimeout   uint32
	JackpotGroup     int32
}

func newLobbySetting(secWebSocketKey, hostExtId, remoteAddr, userAgent string) *lobbySetting {
	return &lobbySetting{
		secWebSocketKey: secWebSocketKey,
		hostExtId:       hostExtId,
		remoteAddr:      remoteAddr,
		userAgent:       userAgent,
		availableGames:  make(map[string]*hostFishGame),
	}
}
