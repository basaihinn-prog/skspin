package define

const (
	Operator_Status_Normal = iota // 正常
	Operatar_Status_Delete        // 删除
)

const (
	User_Status_Normal = iota // 正常
	User_Status_Delete        // 删除
)

const (
	Day   = "Day"
	Month = "Month"
)

const (
	GameStatus_Open        = 0 // 正常
	GameStatus_Maintenance = 1 // 维护中（列表中出现）
	GameStatus_Hide        = 2 // 隐藏（列表中不出现）
	GameStatus_Stop        = 3 // 关闭
)

const (
	GameType_Slot    = 0 // 拉霸游戏
	GameType_Mini    = 1 // 小游戏
	GameType_Poker   = 3 // 棋牌游戏
	GameType_CaiPiao = 4 // 彩票游戏
)
