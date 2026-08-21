package gamerecovery

func Builder() *builder {
	return &builder{
		hostExtId: "",
		hostId:    "",
		memberId:  "",
		gameId:    "",
		subgameId: -1,
	}
}

type builder struct {
	hostExtId string
	hostId    string
	memberId  string
	gameId    string
	subgameId int
}

type gameRecovery struct {
	hostExtId string
	recovery  *recovery
}

// DB table
type recovery struct {
	HostId       string
	MemberId     string
	GameId       string
	SubgameId    int
	GameData     []byte
	AccountingSn uint64 // not used yet
	GameEnd      int    // not used yet
}

func (b *builder) setHostExtId(hostExtId string) *builder {
	b.hostExtId = hostExtId
	return b
}

func (b *builder) setHostId(hostId string) *builder {
	b.hostId = hostId
	return b
}

func (b *builder) setGameId(gameId string) *builder {
	b.gameId = gameId
	return b
}

func (b *builder) setMemberId(memberId string) *builder {
	b.memberId = memberId
	return b
}

func (b *builder) setSubgameId(subgameId int) *builder {
	b.subgameId = subgameId
	return b
}

func (b *builder) build() *gameRecovery {
	return &gameRecovery{
		hostExtId: b.hostExtId,
		recovery: &recovery{
			HostId:       b.hostId,
			MemberId:     b.memberId,
			GameId:       b.gameId,
			SubgameId:    b.subgameId,
			GameData:     nil,
			AccountingSn: 0,
			GameEnd:      -1,
		},
	}
}
