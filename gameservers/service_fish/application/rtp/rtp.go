package rtp

func Builder() *builder {
	return &builder{
		id:               "",
		state:            0,
		budget:           0,
		guestEnable:      true,
		denominator:      0,
		mercenaryInfo:    make(map[uint64]uint64),
		bulletCollection: 0,
	}
}

type builder struct {
	id               string
	state            int
	budget           uint64
	guestEnable      bool
	denominator      uint64
	bulletCollection uint64
	mercenaryInfo    map[uint64]uint64
	netWinGroup      int
}

type rtp struct {
	id               string
	state            int // high:1 low:0 default:-1
	budget           uint64
	guestEnable      bool // bot user is also guest
	denominator      uint64
	bulletCollection uint64
	mercenaryInfo    map[uint64]uint64 // 0:Regular 1:ShotGun 2:Bazooka
	netWinGroup      int               // Simulation Used
}

func (b *builder) reset() {
	b.id = ""
	b.state = -1
	b.budget = 0
	b.guestEnable = true
	b.denominator = 0
	b.mercenaryInfo = nil
	b.bulletCollection = 0
	b.netWinGroup = -1
}

func (b *builder) setId(id string) *builder {
	b.id = id
	return b
}

func (b *builder) setState(state int) *builder {
	b.state = state
	return b
}

func (b *builder) setBudget(budget uint64) *builder {
	b.budget = budget
	return b
}

func (b *builder) setGuestEnable(guestEnable bool) *builder {
	b.guestEnable = guestEnable
	return b
}

func (b *builder) setDenominator(denominator uint64) *builder {
	b.denominator = denominator
	return b
}

func (b *builder) setMercenaryInfo(mercenaryInfo map[uint64]uint64) *builder {
	b.mercenaryInfo = mercenaryInfo
	return b
}

func (b *builder) setBulletCollection(bulletCollection uint64) *builder {
	b.bulletCollection = bulletCollection
	return b
}

func (b *builder) setNetWinGroup(netWinGroup int) *builder {
	b.netWinGroup = netWinGroup
	return b
}

func (b *builder) build() *rtp {
	v := &rtp{
		id:               b.id,
		state:            b.state,
		budget:           b.budget,
		guestEnable:      b.guestEnable,
		denominator:      b.denominator,
		bulletCollection: b.bulletCollection,
		mercenaryInfo:    b.mercenaryInfo,
		netWinGroup:      b.netWinGroup,
	}

	b.reset()

	return v
}
