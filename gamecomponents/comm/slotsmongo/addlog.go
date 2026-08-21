package slotsmongo

type Log struct {
	Pid     int64  `bson:"Pid"`     // 玩家ID
	Bet     int64  `bson:"Bet"`     // 下注
	Win     int64  `bson:"Win"`     // 输赢（是指产出，不是净输赢）
	GameID  string `bson:"GameID"`  // 游戏ID，用于标识是哪个游戏
	RoundID string `bson:"RoundID"` // 对局ID，方便查询
	Comment string `bson:"Comment"` // 注释
	Balance int64  `bson:"Balance"` // 余额
}

// var (
// 	DB = mongodb.NewDB("game")

// 	CollPlayers       = DB.Collection("Players")       // 玩家表
// 	CollApps          = DB.Collection("Apps")          // 接入的app
// 	CollStore         = DB.Collection("Store")         // kv存储
// 	CollBetLog        = DB.Collection("BetLog")        // 下注日志
// 	CollModifyGoldLog = DB.Collection("ModifyGoldLog") // 金币修改日志
// 	CollGames         = DB.Collection("Games")         // 游戏列表

//	CollAppInfoMonth  = DB.Collection("AppInfoMonth")  // 产品统计，月
//	CollAppInfoDay    = DB.Collection("AppInfoDay")    // 产品统计，日
//	CollGameInfoMonth = DB.Collection("GameInfoMonth") // 游戏统计，月
//	CollGameInfoDay   = DB.Collection("GameInfoDay")   // 游戏统计，日
//	CollPlayerInfoDay = DB.Collection("PlayerInfoDay") // 玩家统计，日
//
// )
/*
func day(t time.Time) string {
	return t.Format("2006-01-02")
}

func mon(t time.Time) string {
	return t.Format("2006-01")
}

func AddLog_noused(v *Log) (err error) {
	plr, err := players.Get(v.Pid)
	if err != nil {
		return
	}
	now := time.Now()
	oneLog := &DocBetLog{
		ID:         primitive.NewObjectIDFromTimestamp(now),
		InsertTime: now,
		UserID:     plr.Uid,
		Pid:        plr.Pid,
		AppID:      plr.AppID,
		GameID:     v.GameID,
		RoundID:    v.RoundID,
		Bet:        v.Bet,
		Win:        v.Win,
		Comment:    v.Comment,
		Balance:    v.Balance,
	}

	_, err = GameCollection("BetLog").InsertOne(context.TODO(), oneLog)
	if err != nil {
		return
	}

	opts := options.Update().SetUpsert(true)
	{
		update := db.D(
			"$set", db.D(
				"AppID", oneLog.AppID,
			),
			"$inc", db.D(
				"Bet", v.Bet,
				"Win", v.Win,
			),
		)

		id := fmt.Sprintf("%v_%v", day(now), oneLog.AppID)
		GameCollection("AppInfoDay").UpdateByID(context.TODO(), id, update, opts)

		id = fmt.Sprintf("%v_%v", mon(now), oneLog.AppID)
		GameCollection("AppInfoMonth").UpdateByID(context.TODO(), id, update, opts)
	}

	{
		update := db.D(
			"$set", db.D(
				"GameID", oneLog.GameID,
			),
			"$inc", db.D(
				"Bet", v.Bet,
				"Win", v.Win,
			),
		)
		id := fmt.Sprintf("%v_%v", day(now), oneLog.GameID)
		GameCollection("GameInfoDay").UpdateByID(context.TODO(), id, update, opts)

		id = fmt.Sprintf("%v_%v", mon(now), oneLog.GameID)
		GameCollection("GameInfoMonth").UpdateByID(context.TODO(), id, update, opts)
	}

	{
		update := db.D(
			"$set", db.D(
				"Pid", oneLog.Pid,
			),
			"$inc", db.D(
				"Bet", v.Bet,
				"Win", v.Win,
			),
		)
		id := fmt.Sprintf("%v_%v", day(now), oneLog.Pid)
		GameCollection("PlayerInfoDay").UpdateByID(context.TODO(), id, update, opts)

	}

	// TODO
	// func (a *AlertMgr) SingleWin(log *gamepb.DocBetLog)

	return
}
*/
