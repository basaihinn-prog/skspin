package rank

import (
	"fmt"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/logger"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
)

var Service = &service{
	Id: "RankService",
}

type service struct {
	Id string
}

func (s *service) New(r *rank) {
	defer func() {
		if r != nil && r.redisConn != nil {
			r.redisConn.Close()
		}
	}()

	if r == nil || r.redisConn == nil {
		logger.Service.Zap.Errorw(Rank_REDIS_CONNECTION_NOT_FOUND,
			"GameUser", r.secWebSocketKey,
			"HostId", r.hostId,
			"HostExtId", r.hostExtId,
		)
		errorcode.Service.Fatal(r.secWebSocketKey, Rank_REDIS_CONNECTION_NOT_FOUND)
		return
	}

	if r.kind == KILL_RANKING {
		r.redisScore = 1
	}

	score := int64(0)

	// total record count
	count, err := redigo.Int64(r.redisConn.Do("INCR", r.activityUuid))

	if err != nil {
		logger.Service.Zap.Errorw(Rank_REDIS_INCR_FAILED,
			"GameUser", r.secWebSocketKey,
			"HostId", r.hostId,
			"HostExtId", r.hostExtId,
			"SortedSetKey", r.redisSortedSetKey,
			"HashKey", r.redisHashKey,
			"Value", r.redisValue,
			"Kind", r.kind,
			"Error", err,
		)
		errorcode.Service.Fatal(r.secWebSocketKey, Rank_REDIS_INCR_FAILED)
		return
	}

	reply, err := r.redisConn.Do("HGET", r.redisHashKey, r.redisValue)

	if err != nil {
		logger.Service.Zap.Errorw(Rank_REDIS_HGET_FAILED,
			"GameUser", r.secWebSocketKey,
			"HostId", r.hostId,
			"HostExtId", r.hostExtId,
			"SortedSetKey", r.redisSortedSetKey,
			"HashKey", r.redisHashKey,
			"Value", r.redisValue,
			"Kind", r.kind,
			"Error", err,
		)
		errorcode.Service.Fatal(r.secWebSocketKey, Rank_REDIS_HGET_FAILED)
		return
	}

	if reply == nil {

		reply, err := redigo.Int64(r.redisConn.Do("HSET", r.redisHashKey, r.redisValue, r.redisScore))

		if reply == 0 || err != nil {
			logger.Service.Zap.Errorw(Rank_REDIS_HSET_FAILED,
				"GameUser", r.secWebSocketKey,
				"HostId", r.hostId,
				"HostExtId", r.hostExtId,
				"SortedSetKey", r.redisSortedSetKey,
				"HashKey", r.redisHashKey,
				"Value", r.redisValue,
				"Kind", r.kind,
				"Error", err,
			)
			errorcode.Service.Fatal(r.secWebSocketKey, Rank_REDIS_HSET_FAILED)
			return
		}

		score = int64(r.redisScore)

		logger.Service.Zap.Infow("HINCRBY New Score",
			"GameUser", r.secWebSocketKey,
			"HostId", r.hostId,
			"HostExtId", r.hostExtId,
			"HashKey", r.redisHashKey,
			"Score", score,
			"Value", r.redisValue,
			"Kind", r.kind,
		)

	} else {
		score, _ = redigo.Int64(r.redisConn.Do("HINCRBY", r.redisHashKey, r.redisValue, r.redisScore))

		logger.Service.Zap.Infow("HINCRBY New Score",
			"GameUser", r.secWebSocketKey,
			"HostId", r.hostId,
			"HostExtId", r.hostExtId,
			"HashKey", r.redisHashKey,
			"Score", score,
			"Value", r.redisValue,
			"Kind", r.kind,
		)
	}

	scoreWithCount, _ := strconv.Atoi(fmt.Sprintf("%d%08d", score, count))

	logger.Service.Zap.Infow("ZADD New Score",
		"GameUser", r.secWebSocketKey,
		"HostId", r.hostId,
		"HostExtId", r.hostExtId,
		"SortedSetKey", r.redisSortedSetKey,
		"Score", scoreWithCount,
		"Value", r.redisValue,
		"Kind", r.kind,
	)

	reply, err = redigo.Int64(r.redisConn.Do("ZADD", r.redisSortedSetKey, scoreWithCount, r.redisValue))

	if reply == 0 || err != nil {
		logger.Service.Zap.Errorw(Rank_REDIS_ZADD_FAILED,
			"GameUser", r.secWebSocketKey,
			"HostId", r.hostId,
			"HostExtId", r.hostExtId,
			"SortedSetKey", r.redisSortedSetKey,
			"HashKey", r.redisHashKey,
			"Value", r.redisValue,
			"Kind", r.kind,
			"Error", err,
		)
		errorcode.Service.Fatal(r.secWebSocketKey, Rank_REDIS_ZADD_FAILED)
	}
}
