package auth

import (
	"serve/service_fish/models"
)

const (
	Auth_LOGIN_FAILED             = models.Auth + "0"
	Auth_LOGIN_CALL_PROTO_INVALID = models.Auth + "1"
	Auth_TOKEN_EMPTY              = models.Auth + "2"
	Auth_TOKEN_EXPIRED            = models.Auth + "3"
	Auth_TOKEN_INVALID            = models.Auth + "4"
	Auth_HOST_ID_NOT_FOUND        = models.Auth + "5"
	Auth_JSON_API_INFO_INVALID    = models.Auth + "6"
	Auth_JSON_MEMBER_INFO_INVALID = models.Auth + "7"
	Auth_API_URL_GET_FAILED       = models.Auth + "8"
	Auth_TOKEN_NOT_FOUND          = models.Auth + "9"
	Auth_GAME_DB_NOT_FOUND        = models.Auth + "10"
)
