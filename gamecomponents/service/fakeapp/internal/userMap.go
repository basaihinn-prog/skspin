package internal

import "game/duck/ut2"

var userMap = ut2.NewSyncMap[string, float64]()

func ensureUserExists(uid string) {
	gold, ok := userMap.Load(uid)
	if !ok {
		gold = 10000000
		userMap.Store(uid, gold)
	}
}

var (
	gamecenterUrl = "http://127.0.0.1:11004"
	appId         = "faketrans"
	appSecret     = "11c5d190-3add-4482-b6d4-ee990903f981"
)
