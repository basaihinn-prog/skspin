package activity

import "serve/service_fish/models"

const (
	Activity_REDIS_HGET_ERROR           = models.Activity + "0"
	Activity_REDIS_HGET_NOT_FOUND       = models.Activity + "1"
	Activity_JSON_DECODE_INVALID        = models.Activity + "2"
	Activity_REDIS_CONNECTION_NOT_FOUND = models.Activity + "3"
	Activity_REDIS_DATA_TYPE_INVALID    = models.Activity + "4"
)
