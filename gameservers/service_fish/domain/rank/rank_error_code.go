package rank

import "serve/service_fish/models"

const (
	Rank_REDIS_CONNECTION_NOT_FOUND = models.Rank + "0"
	Rank_REDIS_ZADD_FAILED          = models.Rank + "1"
	Rank_REDIS_ZINCRBY_FAILED       = models.Rank + "2"
	Rank_REDIS_HSET_FAILED          = models.Rank + "3"
	Rank_REDIS_INCR_FAILED          = models.Rank + "4"
	Rank_REDIS_HGET_FAILED          = models.Rank + "5"
)
