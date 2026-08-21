package rank

import (
	"fmt"
	"serve/fish_comm/redis"

	redigo "github.com/gomodule/redigo/redis"
)

const (
	KILL_RANKING = "KILL_RANKING"
	WIN_RANKING  = "WIN_RANKING"
)

func Builder() *builder {
	return &builder{
		activityUuid:    "",
		secWebSocketKey: "",
		memberId:        "",
		hostId:          "",
		hostExtId:       "",
		kind:            "",
		iconId:          -1,
		cent:            0,
		conn:            nil,
	}
}

type builder struct {
	activityUuid    string
	secWebSocketKey string
	memberId        string
	hostId          string
	hostExtId       string
	kind            string
	iconId          int16
	cent            uint64
	conn            redigo.Conn
}

type rank struct {
	activityUuid      string
	secWebSocketKey   string
	memberId          string
	hostId            string
	hostExtId         string
	kind              string
	redisHashKey      string
	redisSortedSetKey string
	redisScore        uint64
	redisValue        string
	redisConn         redigo.Conn
}

func (b *builder) SetActivityUuid(uuid string) *builder {
	b.activityUuid = uuid
	return b
}

func (b *builder) SetSecWebSocketKey(secWebSocketKey string) *builder {
	b.secWebSocketKey = secWebSocketKey
	return b
}

func (b *builder) SetMemberId(memberId string) *builder {
	b.memberId = memberId
	return b
}

func (b *builder) SetHostId(hostId string) *builder {
	b.hostId = hostId
	return b
}

func (b *builder) SetHostExtId(hostExtId string) *builder {
	b.hostExtId = hostExtId
	return b
}

func (b *builder) SetKind(kind string) *builder {
	b.kind = kind
	return b
}

func (b *builder) SetCent(cent uint64) *builder {
	b.cent = cent
	return b
}

func (b *builder) SetConn(hostExtId string) *builder {
	b.conn = redis.Repository.Conn(hostExtId)
	return b
}

func (b *builder) Build() *rank {
	return &rank{
		activityUuid:      b.activityUuid,
		secWebSocketKey:   b.secWebSocketKey,
		memberId:          b.memberId,
		hostId:            b.hostId,
		hostExtId:         b.hostExtId,
		kind:              b.kind,
		redisHashKey:      fmt.Sprintf("%s:%s", b.activityUuid, b.kind),
		redisSortedSetKey: fmt.Sprintf("%s.%s", b.activityUuid, b.kind),
		redisScore:        b.cent,
		redisValue:        fmt.Sprintf("%s:%s", b.hostId, b.memberId),
		redisConn:         b.conn,
	}
}
