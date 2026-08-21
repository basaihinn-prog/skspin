package receipt

import (
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/logger"
	"serve/fish_comm/mysql"
	"sync"
)

var Service = &service{
	id:    "ReceiptService",
	mutex: sync.Mutex{},
}

type service struct {
	id    string
	mutex sync.Mutex
}

func (s *service) New(bulletUuid, secWebSocketKey, hostExtId, gameResult, gameData string,
	accountingSn, bet, pay uint64) {
	//s.mutex.Lock()
	r := newReceipt(bulletUuid, secWebSocketKey, hostExtId, gameResult, gameData, accountingSn, bet, pay)

	logger.Service.Zap.Infow("Receipt",
		"GameUser", secWebSocketKey,
		"BulletUuid", bulletUuid,
		"AccountingSn", accountingSn,
		"Bet", bet,
		"Win", pay,
		"GameData", gameData,
		"GameResult", gameResult,
	)

	if db, err := mysql.Repository.GameDB(hostExtId); err == nil {
		if ok := db.Table("accounting_fish").Create(r.fish).RowsAffected; ok != 1 {
			logger.Service.Zap.Errorw(Receipt_INSERT_ACCOUNTING_FISH_FAILED,
				"GameUser", secWebSocketKey,
				"BulletUuid", bulletUuid,
				"AccountingSn", accountingSn,
				"Bet", bet,
				"Win", pay,
				"GameData", gameData,
				"GameResult", gameResult,
			)
			errorcode.Service.Fatal(secWebSocketKey, Receipt_INSERT_ACCOUNTING_FISH_FAILED)
		}
	} else {
		logger.Service.Zap.Errorw(Receipt_GAME_DB_NOT_FOUND,
			"GameUser", secWebSocketKey,
			"HostExtId", hostExtId,
			"Error", err,
		)
		errorcode.Service.Fatal(secWebSocketKey, Receipt_GAME_DB_NOT_FOUND)
	}
	//s.mutex.Unlock()
}
