package slotscall

import (
	"context"
	"game/duck/lazy"
	"game/duck/logger"
	"game/pb/_gen/pb/slots"
)

// slots 游戏常用rpc

var (
	// GameRpcClient gamepb.GameRpcClient

	SlotsPoolRpcClient              slots.SlotsPoolRpcClient
	SlotsRuntimeRpcClient           slots.SlotsRuntimeRpcClient
	SlotsConfigRpcClient            slots.SlotsConfigRpcClient
	SlotsPlayerGameDayStatRpcClient slots.SlotsPlayerGameDayStatRpcClient
	SlotsGameDayStatRpcClient       slots.SlotsGameDayStatRpcClient
	SlotsStatsBuyGameRpcClient      slots.SlotsStatsBuyGameRpcClient
	SlotsPlayerRotateDataRpcClient  slots.SlotsPlayerRotateDataRpcClient
)

func RpcCallInit() {
	// GameRpcClient = lazy.NewGrpcClient("gamecenter", gamepb.NewGameRpcClient)

	SlotsPoolRpcClient = lazy.NewGrpcClient("slots", slots.NewSlotsPoolRpcClient)
	SlotsRuntimeRpcClient = lazy.NewGrpcClient("slots", slots.NewSlotsRuntimeRpcClient)
	SlotsConfigRpcClient = lazy.NewGrpcClient("slots", slots.NewSlotsConfigRpcClient)
	SlotsPlayerGameDayStatRpcClient = lazy.NewGrpcClient("slots", slots.NewSlotsPlayerGameDayStatRpcClient)
	SlotsGameDayStatRpcClient = lazy.NewGrpcClient("slots", slots.NewSlotsGameDayStatRpcClient)
	SlotsStatsBuyGameRpcClient = lazy.NewGrpcClient("slots", slots.NewSlotsStatsBuyGameRpcClient)
	SlotsPlayerRotateDataRpcClient = lazy.NewGrpcClient("slots", slots.NewSlotsPlayerRotateDataRpcClient)
}

// func SendData(pid int, cbid int, data interface{}) error {
// 	mq.Send2ps([]int64{int64(pid)}, cbid, data)
// 	return nil
// }

// 修改金币。并产生一条金币修改日志
// 会自动给comment拼上游戏名前缀
// func ModifyGold(ctx context.Context, pid int, change int, comment string) (int, error) {
// 	resp, err := GameRpcClient.ModifyGold(ctx, &gamepb.ModifyGoldReq{
// 		Pid:     int64(pid),
// 		Change:  int64(change),
// 		Comment: lazy.ServiceName + ":" + comment,
// 	})
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int(resp.Balance), err
// }

// func GetGold(ctx context.Context, pid int) (int, error) {
// 	resp, err := GameRpcClient.GetBalance(ctx, &gamepb.GetBalanceReq{
// 		Pid: int64(pid),
// 	})

// 	if err != nil {
// 		return 0, err
// 	}

// 	return int(resp.Balance), nil
// }

func GetYeJiRatio(ctx context.Context) (float64, error) {
	return 1.0, nil
	// resp, err := SlotsConfigRpcClient.GetYeJiRatio(ctx, &slots.GetYeJiRatioReq{Game: slots.GetYeJiRatioReq_Slots})
	// if err != nil {
	// 	logger.Err("GetYeJiRatio err:", err)
	// 	return 0, err
	// }

	// return resp.Ratio, err
}

func GetSelfSlotsPool(ctx context.Context, pid int, typ slots.PoolType) (int, error) {
	resp, err := SlotsPoolRpcClient.GetSelfSlotsPool(ctx, &slots.GetSelfSlotsPoolReq{
		Pid:  int64(pid),
		Type: typ,
	})

	if err != nil {
		logger.Err("GetSelfSlotsPool err:", err)
		return 0, err
	}

	return int(resp.Gold), nil
}

func IncSelfSlotsPool(ctx context.Context, pid int, gold int, typ slots.PoolType) error {
	_, err := SlotsPoolRpcClient.IncSelfSlotsPool(ctx, &slots.IncSelfSlotsPoolReq{
		Pid:  int64(pid),
		Gold: int64(gold),
		Type: typ,
	})

	return err
}

func AddRunTimeData(ctx context.Context, pid, expense, income, selfPoolReward int, isSmallGame bool, isFreeGameTurn bool) {
	// req := &slots.AddRunTimeDataReq{
	// 	Game:           lazy.ServiceName,
	// 	Pid:            int64(pid),
	// 	Expense:        int64(expense),
	// 	Income:         int64(income),
	// 	SelfPoolReward: int64(selfPoolReward),
	// 	IsSmallGame:    isSmallGame,
	// 	IsFreeGameTurn: isFreeGameTurn,
	// }

	// _, err := SlotsRuntimeRpcClient.AddRunTimeData(ctx, req)
	// if err != nil {
	// 	logger.Err("AddRunTimeData err:", err, req)
	// }
}

func AddSlotsStatus(expand, income int, comment string) {
	// game := lazy.ServiceName
	// winBigAwardStatistics(expand, income, comment)
}

func InsertTodayAward(expend float64, times float64, comment string) {
	// player.InsertTodayAward(gameType.NiuBi, float64(bet), res.Times, res.UniqueStr)
}

func PlayerGameDayStatIncBoth(data *slots.SlotsPlayerGameDayStatRecord) error {
	// _, err := SlotsPlayerGameDayStatRpcClient.IncBoth(context.Background(), data)
	// if err != nil {
	// 	logger.Err("PlayerGameDayStatIncBoth err:", err, data)
	// }
	return nil
}

func PlayerGameDayStatIncAll(data *slots.SlotsPlayerGameDayStatRecord) error {
	// _, err := SlotsPlayerGameDayStatRpcClient.IncAll(context.Background(), data)
	// if err != nil {
	// 	logger.Err("PlayerGameDayStatIncAll err:", err, data)
	// }
	return nil
}

func GameDayStatIncAll(data *slots.SlotsGameDayStatRecord) error {
	_, err := SlotsGameDayStatRpcClient.IncAll(context.Background(), data)
	if err != nil {
		logger.Err("GameDayStatIncAll err:", err, data)
	}
	return err
}

func RuntimeGetTotal(ctx context.Context) (int64, int64, error) {
	return RuntimeGetTotalByGame(ctx, lazy.ServiceName)
}

func RuntimeGetTotalByGame(ctx context.Context, game string) (int64, int64, error) {
	resp, err := SlotsRuntimeRpcClient.GetTotal(ctx, &slots.GetTotalReq{Game: game})
	if err != nil {
		logger.Err("SlotsRuntimeRpcClient.GetTotal err:", err)
		return 0, 0, err
	}
	return resp.Income, resp.Expense, nil
}

func FixRunTimeData(ctx context.Context, day string) (*slots.ThreeMonkeyRunTimeData, error) {
	return FixRunTimeDataByGame(ctx, day, lazy.ServiceName)
}

func FixRunTimeDataByGame(ctx context.Context, day string, game string) (*slots.ThreeMonkeyRunTimeData, error) {
	resp, err := SlotsRuntimeRpcClient.FixRunTimeData(ctx, &slots.FixRunTimeDataReq{Day: day, Game: game})
	if err != nil {
		logger.Err("SlotsRuntimeRpcClient.GetTotal err:", err)
		return nil, err
	}
	return resp, nil
}

func GameDayStatGetDocs(ctx context.Context, req *slots.SlotsGameDayStatDocsReq) ([]*slots.SlotsGameDayStatRecord, error) {
	resp, err := SlotsGameDayStatRpcClient.Docs(ctx, req)
	if err != nil {
		logger.Err("SlotsGameDayStatRpcClient.Docs(ctx, req) err:", err, req)
		return nil, err
	}

	return resp.List, nil
}

func GetGameDayWinCount(ctx context.Context, game string, day string) (int64, error) {
	resp, err := SlotsPlayerGameDayStatRpcClient.GetGameDayWinCount(ctx, &slots.GetGameDayWinCountReq{Game: game, Day: day})
	if err != nil {
		logger.Err("SlotsGameDayStatRpcClient.Docs(ctx, req) err:", err, game, day)
		return 0, err
	}

	return resp.Count, nil
}

func StatsBuyGameAdd(ctx context.Context, game string, bet int, win int) error {
	// _, err := SlotsStatsBuyGameRpcClient.Add(ctx, &slots.SlotsStatsBuyGame_Doc{
	// 	Day:  time.Now().Format("2006-01-02"),
	// 	Game: game,
	// 	Bet:  int64(bet),
	// 	Win:  int64(win),
	// })

	// if err != nil {
	// 	logger.Err("SlotsStatsBuyGameRpcClient.Add err:", game, bet, win, err)
	// }

	return nil
}

func StatsBuyGameQuery(ctx context.Context, days []string) ([]*slots.SlotsStatsBuyGame_Doc, error) {
	return StatsBuyGameQueryByGame(ctx, days, lazy.ServiceName)
}

func StatsBuyGameQueryByGame(ctx context.Context, days []string, game string) ([]*slots.SlotsStatsBuyGame_Doc, error) {
	resp, err := SlotsStatsBuyGameRpcClient.Query(ctx, &slots.SlotsStatsBuyGame_QueryReq{
		Days: days,
		Game: game,
	})

	if err != nil {
		logger.Err("SlotsStatsBuyGameRpcClient.Query err:", days, err)
	}

	return resp.List, err
}

func SlotsPlayerRotateData_IncRotate(ctx context.Context, pid int, game string) error {
	_, err := SlotsPlayerRotateDataRpcClient.IncRotate(ctx, &slots.SlotsPlayerRotateData_IncRotateReq{
		Pid:  int64(pid),
		Game: game,
	})

	if err != nil {
		logger.Err("SlotsPlayerRotateDataRpcClient.IncRotate err:", pid, game, err)
		return err
	}

	return err
}

func SlotsPlayerRotateData_GetRotate(ctx context.Context, pid int, game string) (int, error) {
	resp, err := SlotsPlayerRotateDataRpcClient.GetRotate(ctx, &slots.SlotsPlayerRotateData_GetRotateReq{
		Pid:  int64(pid),
		Game: game,
	})
	if err != nil {
		logger.Err("SlotsPlayerRotateDataRpcClient.GetRotate err:", pid, game, err)
		return 0, err
	}
	return int(resp.Count), err
}
