package rpc

import (
	"serve/comm"
	"serve/comm/slotsmongo"
	"serve/service/pg_39/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 添加记录
func addRecord(player *models.Player, bet, win, balance int64 /*roundId string,*/, hasBonusGame bool, isBuy bool, grade int, pgbetid primitive.ObjectID, bigReward int64, id string) {
	player.OnSpinFinish(bet, win, isBuy, hasBonusGame, 1)

	// slotsmongo.AddBetLog(player.PID, bet, win, balance, roundId, true, win-bet, isBuy, grade,
	// 	pgbetid, bigReward, comm.GameType_Slot)

	slotsmongo.AddBetLog(&slotsmongo.AddBetLogParams{
		UserName:     "",
		CurrencyKey:  player.CurrencyKey,
		ID:           id,
		Pid:          player.PID,
		Bet:          bet,
		Win:          win,
		Balance:      balance,
		RoundID:      id,
		Completed:    true,
		TotalWinLoss: win - bet,
		IsBuy:        false,
		Grade:        grade,
		PGBetID:      pgbetid,
		BigReward:    bigReward,
		GameType:     comm.GameType_Slot,
	})
}
