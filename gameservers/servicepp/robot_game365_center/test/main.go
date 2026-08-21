package main

// Response 是整个响应的结构体
type Response struct {
	Status    int      `json:"status"`    // 状态码
	Msg       string   `json:"msg"`       // 消息
	Data      GameData `json:"data"`      // 游戏数据
	Timestamp int64    `json:"timestamp"` // 时间戳
}

// GameData 是游戏数据的结构体
type GameData struct {
	ID               int64  `json:"id"`                         // 游戏ID
	Status           int    `json:"status"`                     // 状态
	IDs              []int  `json:"ids,omitempty"`              // ID列表，可选
	AddUserId        int64  `json:"addUserId,omitempty"`        // 添加用户ID，可选
	ModUserId        int64  `json:"modUserId,omitempty"`        // 修改用户ID，可选
	CompanyId        int64  `json:"companyId,omitempty"`        // 公司ID，可选
	CompanyName      string `json:"companyName"`                // 公司名称
	GameName         string `json:"gameName"`                   // 游戏名称
	GameNameZh       string `json:"gameNameZh"`                 // 游戏名称（中文）
	GameIcon         string `json:"gameIcon"`                   // 游戏图标URL
	GameIconPath     string `json:"gameIconPath,omitempty"`     // 游戏图标路径，可选
	GameCategory     string `json:"gameCategory,omitempty"`     // 游戏分类，可选
	GamePlatform     string `json:"gamePlatform"`               // 游戏平台
	GameType         string `json:"gameType"`                   // 游戏类型
	GameKey          string `json:"gameKey,omitempty"`          // 游戏键，可选
	GameId           string `json:"gameId"`                     // 游戏ID
	Jackpot          string `json:"jackpot,omitempty"`          // 彩池，可选
	GameRanking      string `json:"gameRanking,omitempty"`      // 游戏排名，可选
	GameRtp          string `json:"gameRtp"`                    // 返回玩家的百分比
	FeatureBuy       string `json:"featureBuy,omitempty"`       // 特性购买，可选
	Volatility       string `json:"volatility,omitempty"`       // 波动性，可选
	MaxExposure      string `json:"maxExposure,omitempty"`      // 最大暴露，可选
	MultiLanguage    string `json:"multiLanguage,omitempty"`    // 多语言，可选
	Remark           string `json:"remark"`                     // 备注
	Address          string `json:"address,omitempty"`          // 地址，可选
	GameNum          int    `json:"gameNum,omitempty"`          // 游戏编号，可选
	GameCategoryCode string `json:"gameCategoryCode,omitempty"` // 游戏分类代码，可选
	Gameplat         string `json:"gameplat,omitempty"`         // 游戏平台（备用），可选
	Lang             string `json:"lang"`                       // 支持的语言
	Currency         string `json:"currency"`                   // 支持的货币
	Rate             string `json:"rate,omitempty"`             // 汇率，可选
	GameScreenshot   string `json:"gameScreenshot,omitempty"`   // 游戏截图，可选
	GameFile         string `json:"gameFile,omitempty"`         // 游戏文件，可选
	MaxWin           string `json:"maxWin,omitempty"`           // 最大获胜，可选
	IconGameKey      string `json:"iconGameKey,omitempty"`      // 图标游戏键，可选
	BrandName        string `json:"brandName,omitempty"`        // 品牌名称，可选
	GameEncrypt      string `json:"gameEncrypt"`                // 游戏加密
	Support          string `json:"support,omitempty"`          // 支持，可选
}
type ResponseGame struct {
	Data      DataGame `json:"data"`
	Msg       string   `json:"msg"`
	Status    int      `json:"status"`
	Timestamp int64    `json:"timestamp"`
}

type DataGame struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

//func main() {
//	client := &http.Client{}
//	internal.GetHotGameClubNewmgcKey("", client)
//}
