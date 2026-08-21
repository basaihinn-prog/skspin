package room

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"serve/comm/redisx"
	"time"
)

type Robot struct {
	Id         string
	PlayerIcon string
	PlayerName string
	BetInfo    []*Bets
	RoomId     string
	BetNum     int
	Seed       string
}

func (rb *Robot) RobotBet() {
	var err error
	for _, betInfo := range rb.BetInfo {
		bet := Bets{
			Bet:          betInfo.Bet,
			BetID:        betInfo.BetID,
			Currency:     betInfo.Currency,
			PlayerID:     rb.Id,
			ProfileImage: rb.PlayerIcon,
			Username:     rb.PlayerName,
			WinAmount:    0,
		}
		err = StoreBet(redisx.GetClient(), rb.RoomId, bet)
		if err != nil {
			slog.Error("RobotBet::StoreBet Err", "err", err)
			return
		}
	}
}

func RobotBet2(bet2 *Bets, r *Room) {
	var err error
	bet := Bets{
		Bet:          bet2.Bet,
		BetID:        bet2.BetID,
		Currency:     bet2.Currency,
		PlayerID:     bet2.PlayerID,
		ProfileImage: bet2.ProfileImage,
		Username:     bet2.Username,
		WinAmount:    0,
	}
	err = StoreBet(redisx.GetClient(), r.Name, bet)
	if err != nil {
		slog.Error("RobotBet::StoreBet Err", "err", err)
		return
	}
	bet2.Effective = true
}

func RobotCashOut2(bet2 *Bets, r *Room) {
	var err error
	winAmount := math.Round(bet2.Bet*bet2.Multiplier*100) / 100
	bet := Bets{
		Bet:          bet2.Bet,
		BetID:        bet2.BetID,
		Currency:     bet2.Currency,
		PlayerID:     bet2.PlayerID,
		ProfileImage: bet2.ProfileImage,
		Username:     bet2.Username,
		WinAmount:    winAmount,
		Multiplier:   bet2.Multiplier,
	}
	err = UpdateBet(redisx.GetClient(), r.Name, bet)
	if err != nil {
		slog.Error("RobotCashOut::UpdateBet Err", "err", err)
		return
	}
	//min, err := UpdateBet(redisx.GetClient(), r.Name, bet)
	//if err != nil {
	//	slog.Error("RobotCashOut::UpdateBet Err", "err", err)
	//	return
	//}
	//if min > 0 {
	//	r.MinBet = min
	//}
	topWin := TopWin{
		Bet:                     bet2.Bet,
		EndDate:                 time.Now().UTC().UnixMilli(),
		Currency:                bet2.Currency,
		MaxMultiplier:           r.XArr[len(r.XArr)-1],
		Payout:                  bet2.Multiplier,
		PlayerId:                bet2.PlayerID,
		ProfileImage:            bet2.ProfileImage,
		RoundBetId:              "9999999999",
		RoundId:                 r.RoundId,
		Username:                bet2.Username,
		WinAmountInMainCurrency: winAmount,
		WinAmount:               winAmount,
		Zone:                    "aviator_core_inst2_demo1",
	}
	if topWin.Payout >= 4.0 {
		for _, date := range DateList {
			keyTopWin := fmt.Sprintf("%s_%s_%s", r.Name, "topwin", date)
			StoreTopWins(redisx.GetClient(), keyTopWin, topWin)
			keyHugeWin := fmt.Sprintf("%s_%s_%s", r.Name, "huge", date)
			StoreHugeWins(redisx.GetClient(), keyHugeWin, topWin)
		}
	}
}

//func (rb *Robot) RobotCashOut(betId int) {
//	var err error
//	for i := range rb.BetInfo {
//		if rb.BetInfo[i].BetID != betId || rb.BetInfo[i].RobotBetFinishFlag {
//			return
//		}
//		bet := Bets{
//			Bet:          rb.BetInfo[i].Bet,
//			BetID:        rb.BetInfo[i].BetID,
//			Currency:     rb.BetInfo[i].Currency,
//			PlayerID:     rb.Id,
//			ProfileImage: rb.PlayerIcon,
//			Username:     rb.PlayerName,
//			WinAmount:    rb.BetInfo[i].Bet * rb.BetInfo[i].Multiplier,
//			Multiplier:   rb.BetInfo[i].Multiplier,
//		}
//		err = UpdateBet(redisx.GetClient(), rb.RoomId, bet)
//		if err != nil {
//			slog.Error("RobotCashOut::UpdateBet Err", "OB", "")
//			return
//		}
//		rb.BetInfo[i].RobotBetFinishFlag = true
//		rb.BetNum = rb.BetNum - 1
//	}
//}

func (rb *Robot) GetMultiplier(betId int) float64 {
	multiplier := 0.0
	num := rand.IntN(100)
	if num < 15 { //15概率是4-7之间
		randFloat := float64(GetRandInt(400, 700))
		multiplier = randFloat / 100
	} else if num < 25 { //25概率是2-4之间
		randFloat := float64(GetRandInt(200, 400))
		multiplier = randFloat / 100
	} else if num < 60 { //60概率是1.1-2之间
		randFloat := float64(GetRandInt(110, 200))
		multiplier = randFloat / 100
	}
	return multiplier
}

func GetRandInt(min, max int) int {
	if max <= min {
		slog.Info("GetRandInt::因为入参错误即将赋予默认值1", "min", min, "max", max)
		return 1
	}
	return min + rand.IntN(max-min+1)
}
