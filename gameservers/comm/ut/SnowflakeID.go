package ut

import (
	"math/rand/v2"
	"sync"
	"time"
)

const (
	workerIDBits     = 5
	dataCenterIDBits = 5
	sequenceBits     = 12

	maxWorkerID     = -1 ^ (-1 << workerIDBits)
	maxDataCenterID = -1 ^ (-1 << dataCenterIDBits)
	sequenceMask    = -1 ^ (-1 << sequenceBits)

	workerIDShift      = sequenceBits
	dataCenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + dataCenterIDBits
	twepoch            = int64(1288834974657)
)

type Snowflake struct {
	sync.Mutex
	timestamp    int64
	workerID     int64
	dataCenterID int64
	sequence     int64
}

func NewSnowflake() *Snowflake {
	return &Snowflake{
		workerID:     int64(rand.IntN(maxWorkerID)),
		dataCenterID: int64(rand.IntN(maxDataCenterID)),
	}
}

func (s *Snowflake) NextID() int64 {
	s.Lock()
	defer s.Unlock()

	timestamp := time.Now().UnixNano() / 1e6

	if s.timestamp == timestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for timestamp <= s.timestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = timestamp
	id := ((timestamp - twepoch) << timestampLeftShift) |
		(s.dataCenterID << dataCenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}

func (s *Snowflake) NextIDWithTime(timestamp int64) int64 {
	s.Lock()
	defer s.Unlock()

	if s.timestamp == timestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			timestamp += 1
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = timestamp
	id := ((timestamp - twepoch) << timestampLeftShift) |
		(s.dataCenterID << dataCenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}

func ParseSnowflakeID(id int64) (timestamp, dataCenterID, workerID, sequence int64) {
	sequence = id & sequenceMask
	workerID = (id >> workerIDShift) & maxWorkerID
	dataCenterID = (id >> dataCenterIDShift) & maxDataCenterID
	timestamp = (id >> timestampLeftShift) + twepoch
	return
}

func ParseSnowflakeUnix(id int64) (timestamp int64) {
	timestamp = (id >> timestampLeftShift) + twepoch
	return
}
