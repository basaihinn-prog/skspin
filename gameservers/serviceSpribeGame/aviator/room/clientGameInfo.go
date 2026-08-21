package room

import (
	"encoding/json"
	"log/slog"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"slices"
	"strconv"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

func ClientGameInfo(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	//req := ClientGameInfoReq{}
	//err := json.Unmarshal(data, &req)
	//if err != nil {
	//	slog.Error("ClientLogin Err", "json.Unmarshal", err)
	//}
	session, err := comm.GetSession(c)
	if err != nil {
		slog.Error("ClientGameInfo Err", "GetSession err", err)
		rsp := GetErrMsg("", comm.ErrNoPlr, "NOT Found This Plr")
		rspMarshal, _ := json.Marshal(rsp)
		c.WriteMessage(messageType, rspMarshal)
	}
	plr := session.Plr
	info, err := redisx.LoadAppIdCache(plr.AppID)
	if err != nil {
		slog.Error("ClientGameInfo Err", "LoadAppIdCache err", err)
		rsp := GetErrMsg("", comm.ErrNoPlr, "NOT Found This Plr")
		rspMarshal, _ := json.Marshal(rsp)
		c.WriteMessage(messageType, rspMarshal)
	}
	playerIdStr := strconv.Itoa(int(plr.PID))
	plrInfo := r.Players[playerIdStr]
	//rsp := ClientGameInfoRsp{
	//	Id:               1,
	//	TargetController: 0,
	//	Content: &Content{
	//		Pi: 0,
	//		Rl: [][]interface{}{
	//			{"1", "game_state", "default", false, false, 0, 20, []interface{}{}},
	//		}, //[1, "game_state", "default", false, false, 0, 20, []]//todo 待替换
	//		Rs: 0,
	//		Un: `67590&&release`,     //todo 待替换
	//		Zn: "aviator_core_inst4", //todo 待替换
	//	},
	//}
	//marshal, _ := json.Marshal(rsp)
	marshal, _ := GetClientGameInfoRsp().ToBinary()
	initRsp := &CommBody{}
	//json.Unmarshal([]byte(initStr), initRsp)

	chatBody, err := GetHashAll(redisx.GetClient(), "room:"+r.Name+":message")
	messages := make([]*ChatHistory, 0)
	for _, chatInfo := range chatBody {
		message := ChatHistory{}
		err = json.Unmarshal([]byte(chatInfo), &message)
		if err != nil {
			slog.Error("ClientGameInfo Err", "ClientGameInfo err", err)
		}
		if message.Likes != nil {
			for _, val := range message.Likes.UsersWhoLike {
				if val == playerIdStr {
					message.Likes.IsMeLiked = true
					break
				}
			}
		}
		messages = append(messages, &message)
	}
	config := &Config{}
	json.Unmarshal([]byte(configg), config)
	config.BetOptions = info.Cs
	config.MinBet = info.OnlineUpNum
	config.MaxBet = info.OnlineDownNum
	config.BetInputStep = info.Scale
	config.DefaultBetValue = info.DefaultCs
	config.Currency = r.Currency
	roundsInfo := make([]*RoundsInfo, len(r.RoundHistories))
	copy(roundsInfo, r.RoundHistories)
	slices.Reverse(roundsInfo)
	initContent := &Content{
		C: "init",
		P: P{
			Code: comm.Succ,

			RoundsInfo:    roundsInfo,
			ActiveBets:    make([]interface{}, 0),
			OnlinePlayers: len(r.Players) + 4923, //todo
			ChatSettings: &ChatSettings{
				Top10Gifs: Top10Gifs,
			},
			ChatHistory:        messages,
			Config:             config,
			RoundId:            r.RoundId,
			StageID:            1,
			ActiveFreeBetsInfo: make([]interface{}, 0),
		},
	}
	if r.State == ROOMSTATE_WAIT_BET {
		initContent.P.StageID = 1
		psTime := int32(time.Now().UnixMilli() - r.RoundStartDate)
		initContent.P.FullBetTime = int32(TIme * Duration)
		initContent.P.BetTimeLeft = int32(TIme*Duration) - psTime
	} else if r.State == ROOMSTATE_FLY {
		initContent.P.StageID = 2
		initContent.P.CurrentMultiplier = r.XArr[r.XIndex]
	} else {
		initContent.P.StageID = 3
	}
	initRsp = &CommBody{
		Id:               13,
		TargetController: 1,
		Content:          initContent,
	}
	//initRsp.Content.P.RoundsInfo = r.RoundHistories
	//initRsp.Content.P.ChatHistory = chat
	//initRsp.Content.P.OnlinePlayers = len(r.Players)
	//initRsp.Content.P.RoundId = r.RoundId
	balance, err := slotsmongo.GetBalance(plr.PID)
	if err != nil {
		slog.Error("ClientGameInfo Err", "Player Has No Balance", err)
		rsp := GetErrMsg("init", comm.ErrNoBalance, "Player Has No Balance")
		rspMarshal, _ := json.Marshal(rsp)
		c.WriteMessage(messageType, rspMarshal)
		return
	}
	initRsp.Content.P.User = &User{
		UserID:       playerIdStr,
		Username:     plrInfo.PlayerName,
		Balance:      ut.Gold2Money(balance),
		ProfileImage: plrInfo.PlayerIcon,
		Settings:     NewSettings(),
	}
	//initMarshal, err := json.Marshal(initRsp)
	//if err != nil {
	//}
	//initMarshal, _ := GetClientGameInfoInitRsp(initRsp).ToBinary()
	ss := GetClientGameInfoInitRsp(initRsp)
	initMarshal, err := ss.ToBinary()
	if err != nil {

	}
	c.WriteMessage(messageType, marshal)
	c.WriteMessage(messageType, initMarshal)
}

type ClientGameInfoReq struct {
	Channel          string   `json:"channel"`
	TargetController int      `json:"targetController"`
	Content          *Content `json:"content"`
}

type ClientGameInfoRsp struct {
	Id               int      `json:"id"`
	TargetController int      `json:"targetController"`
	Content          *Content `json:"content"`
}

func GetClientGameInfoRsp() *ut.SFSObject {
	//自己做一个
	so := &ut.SFSObject{}
	so.Init()
	p := &ut.SFSObject{}
	p.Init()
	p.PutShort("rs", 0)
	p.PutString("zn", "aviator_inst1")
	p.PutString("un", `xckSpribetest&&huidustg`)
	p.PutShort("pi", 0)

	rl := ut.NewSFSArray()
	inrl := ut.NewSFSArray()
	inrl.Add(int32(421), ut.INT, true)
	inrl.Add("game_state", ut.UTF_STRING, true)
	inrl.Add("default", ut.UTF_STRING, true)
	inrl.Add(false, ut.BOOL, true)
	inrl.Add(false, ut.BOOL, true)
	inrl.Add(false, ut.BOOL, true)
	inrl.Add(int16(0), ut.SHORT, true)
	inrl.Add(int16(20), ut.SHORT, true)
	inrl.Add(&ut.SFSArray{}, ut.SFS_ARRAY, true)
	rl.Add(inrl, ut.SFS_ARRAY, true)
	p.PutSFSArray("rl", rl)

	p.PutInt("id", 1937) //todo 不知道干什么的id

	so.PutSFSObject("p", p)
	so.PutShort("a", 1)
	so.PutByte("c", 0)
	return so
}

func GetClientGameInfoInitRsp(initRsp *CommBody) *ut.SFSObject {
	so := &ut.SFSObject{}
	so.Init()
	p := &ut.SFSObject{}
	p.Init()

	pp := &ut.SFSObject{}
	pp.Init()
	roundsInfo := ut.NewSFSArray()
	if len(initRsp.Content.P.RoundsInfo) > 0 {
		for _, val := range initRsp.Content.P.RoundsInfo {
			if val != nil {
				temp := ut.NewSFSObject()
				temp.PutDouble("maxMultiplier", val.MaxMultiplier)
				temp.PutInt("roundId", int32(val.RoundID))
				roundsInfo.Add(temp, ut.SFS_OBJECT, true)
			}
		}
	}
	pp.PutSFSArray("roundsInfo", roundsInfo)
	pp.PutInt("onlinePlayers", int32(initRsp.Content.P.OnlinePlayers))
	chatSettings := ut.NewSFSObject()
	top10Gifs := ut.NewSFSArray()
	topGifs := getDefaultGif()
	for _, val := range topGifs {
		temp := ut.NewSFSObject()
		temp.PutString("id", val.Id)
		temp.PutString("url", val.Url)
		dims := ut.NewSFSArray()
		dims.Add(val.Dims[0], ut.INT, true)
		dims.Add(val.Dims[1], ut.INT, true)
		temp.PutSFSArray("dims", dims)
		top10Gifs.Add(temp, ut.SFS_OBJECT, true)
	}
	chatSettings.PutSFSArray("top10Gifs", top10Gifs)
	pp.PutSFSObject("chatSettings", chatSettings)
	user := ut.NewSFSObject()
	settings := ut.NewSFSObject()
	settings.PutBool("music", true)
	settings.PutBool("sound", true)
	settings.PutBool("secondBet", true)
	settings.PutBool("animation", true)
	user.PutDouble("balance", initRsp.Content.P.User.Balance)
	user.PutString("profileImage", initRsp.Content.P.User.ProfileImage)
	user.PutString("userId", initRsp.Content.P.User.UserID)
	user.PutString("username", initRsp.Content.P.User.Username)
	user.PutSFSObject("settings", settings)
	pp.PutSFSObject("user", user)
	pp.PutInt("roundId", int32(initRsp.Content.P.RoundId))
	pp.PutInt("code", 200)
	pp.PutInt("fullBetTime", 5000)
	pp.PutSFSArray("activeBets", ut.NewSFSArray())
	chatHistory := ut.NewSFSArray()
	chatHistory = GetChatHistory(chatHistory, initRsp.Content.P.ChatHistory)
	pp.PutSFSArray("chatHistory", chatHistory)
	pp.PutSFSArray("activeFreeBets", ut.NewSFSArray())
	pp.PutInt("betTimeLeft", 9)
	config := getConfig(initRsp.Content.P.Config)
	pp.PutSFSObject("config", config)
	pp.PutInt("stageId", int32(initRsp.Content.P.StageID))
	if initRsp.Content.P.StageID == 1 {
		pp.PutInt("fullBetTime", initRsp.Content.P.FullBetTime)
		pp.PutInt("betTimeLeft", initRsp.Content.P.BetTimeLeft)
	} else if initRsp.Content.P.StageID == 2 {
		pp.PutDouble("currentMultiplier", initRsp.Content.P.CurrentMultiplier)
	}

	p.PutSFSObject("p", pp)
	p.PutString("c", "init")

	so.PutSFSObject("p", p)
	so.PutShort("a", 13)
	so.PutByte("c", 1)
	return so
}

func getConfig(config *Config) *ut.SFSObject {
	res := ut.NewSFSObject()
	// 填充根对象字段
	res.PutBool("isAutoBetFeatureEnabled", config.IsAutoBetFeatureEnabled)
	res.PutBool("isBetsHistoryEndBalanceEnabled", config.IsBetsHistoryStartBalanceEnabled)
	res.PutInt("betPrecision", int32(config.BetPrecision))
	res.PutDouble("maxBet", float64(config.MaxBet))
	res.PutInt("multiplierPrecision", int32(config.MultiplierPrecision))
	res.PutString("accountHistoryActionType", config.AccountHistoryActionType)

	// 填充 autoBetOptions 对象
	autoBetOptions := ut.NewSFSObject()
	autoBetOptions.PutBool("decreaseOrExceedStopPointReq", config.AutoBetOptions.DecreaseOrExceedStopPointReq)

	//numberOfRounds := ut.NewSFSArray()
	//for _, round := range config.AutoBetOptions.NumberOfRounds {
	//	numberOfRounds.Add()(round)
	//}
	//autoBetOptions.PutSFSArray("numberOfRounds", numberOfRounds)
	//res.PutSFSObject("autoBetOptions", autoBetOptions)

	// 填充 autoCashOut 对象
	autoCashOut := ut.NewSFSObject()
	autoCashOut.PutDouble("minValue", config.AutoCashOut.MinValue)
	autoCashOut.PutDouble("defaultValue", config.AutoCashOut.DefaultValue)
	autoCashOut.PutDouble("maxValue", float64(config.AutoCashOut.MaxValue))
	res.PutSFSObject("autoCashOut", autoCashOut)

	// 填充其他字段
	res.PutString("backToHomeActionType", config.BackToHomeActionType)
	res.PutDouble("betInputStep", config.BetInputStep)

	// 填充 betOptions 数组
	betOptions := ut.NewSFSArray()
	for _, bet := range config.BetOptions {
		betOptions.Add(bet, ut.DOUBLE, true)
	}
	res.PutSFSArray("betOptions", betOptions)

	// 填充 chat 对象
	chat := ut.NewSFSObject()

	promo := ut.NewSFSObject()
	promo.PutBool("isEnabled", config.Chat.Promo.IsEnabled)
	chat.PutSFSObject("promo", promo)

	rain := ut.NewSFSObject()
	rain.PutBool("isEnabled", config.Chat.Rain.IsEnabled)
	rain.PutDouble("rainMinBet", config.Chat.Rain.RainMinBet)
	rain.PutShort("defaultNumOfUsers", int16(config.Chat.Rain.DefaultNumOfUsers))
	rain.PutShort("minNumOfUsers", int16(config.Chat.Rain.MinNumOfUsers))
	rain.PutShort("maxNumOfUsers", int16(config.Chat.Rain.MaxNumOfUsers))
	rain.PutDouble("rainMaxBet", float64(config.Chat.Rain.RainMaxBet))
	chat.PutSFSObject("rain", rain)

	chat.PutBool("isEnabled", config.Chat.IsEnabled)
	chat.PutBool("isGifsEnabled", config.Chat.IsGifsEnabled)
	chat.PutShort("maxMessageLength", int16(config.Chat.MaxMessageLength))
	chat.PutShort("maxMessages", int16(config.Chat.MaxMessages))
	chat.PutInt("sendMessageDelay", int32(config.Chat.SendMessageDelay))
	res.PutSFSObject("chat", chat)

	// 填充其他字段
	res.PutString("currency", config.Currency)
	res.PutDouble("defaultBetValue", float64(config.DefaultBetValue))

	// 填充 engagementTools 对象（空对象）
	engagementTools := ut.NewSFSObject()
	engagementTools.PutBool("isExternalChatEnabled", false)
	res.PutSFSObject("engagementTools", engagementTools)

	res.PutInt("fullBetTime", int32(config.FullBetTime))
	res.PutString("gameRulesAutoCashOutType", config.GameRulesAutoCashOutType)
	res.PutInt("inactivityTimeForDisconnect", int32(config.InactivityTimeForDisconnect))
	res.PutString("ircDisplayType", config.IrcDisplayType)
	res.PutBool("isActiveGameFocused", config.IsActiveGameFocused)
	res.PutBool("isAlderneyModalShownOnInit", config.IsAlderneyModalShownOnInit)
	res.PutBool("isBalanceValidationEnabled", config.IsBalanceValidationEnabled)
	res.PutBool("isBetsHistoryStartBalanceEnabled", config.IsBetsHistoryStartBalanceEnabled)
	res.PutBool("isClockVisible", config.IsClockVisible)
	res.PutBool("isCurrencyNameHidden", config.IsCurrencyNameHidden)
	res.PutBool("isFreeBetsEnabled", config.IsFreeBetsEnabled)
	res.PutBool("isGameRulesHaveMaxWin", config.IsGameRulesHaveMaxWin)
	res.PutBool("isGameRulesHaveMinimumBankValue", config.IsGameRulesHaveMinimumBankValue)
	res.PutBool("isHolidayTheme", config.IsHolidayTheme)
	res.PutBool("isLiveBetsAndStatisticsHidden", config.IsLiveBetsAndStatisticsHidden)
	res.PutBool("isLoginTimer", config.IsLoginTimer)
	res.PutBool("isMaxUserMultiplierEnabled", config.IsMaxUserMultiplierEnabled)
	res.PutBool("isMultipleBetsEnabled", config.IsMultipleBetsEnabled)
	res.PutBool("isNetSessionEnabled", config.IsNetSessionEnabled)
	res.PutBool("isUseMaskedUsername", config.IsUseMaskedUsername)
	res.PutInt("maxUserWin", int32(config.MaxUserWin))
	res.PutDouble("minBet", float64(config.MinBet))
	res.PutString("onLockUIActions", config.OnLockUIActions)
	res.PutBool("quickHelp", config.QuickHelp)
	res.PutShort("returnToPlayer", int16(config.ReturnToPlayer))

	// 填充 newChat 对象
	newChat := ut.NewSFSObject()
	newChat.PutNull("promo")
	newChat.PutNull("rain")
	newChat.PutBool("isEnabled", config.NewChat.IsEnabled)
	newChat.PutBool("isGifsEnabled", config.NewChat.IsGifsEnabled)
	newChat.PutShort("maxMessageLength", int16(config.NewChat.MaxMessageLength))
	newChat.PutShort("maxMessages", int16(config.NewChat.MaxMessages))
	newChat.PutInt("sendMessageDelay", int32(config.NewChat.SendMessageDelay))
	res.PutSFSObject("newChat", newChat)

	// 填充其他字段
	res.PutInt("pingIntervalMs", int32(config.PingIntervalMs))
	res.PutBool("hideHeader", config.HideHeader)
	res.PutBool("isLogoUrlHidden", false)
	//
	res.PutBool("isShowWinAmountUntilNextRound", false)
	res.PutBool("isShowBetControlNumber", false)
	res.PutString("modalShownOnInit", "none")
	res.PutBool("isGameRulesHaveMultiplierFormula", false)
	res.PutBool("isEmbeddedVideoHidden", false)
	res.PutInt("chatApiVersion", 2)
	res.PutBool("isBetTimerBranded", true)
	res.PutBool("showCrashExampleInRules", false)

	return res
}

// -------------------------------------------------------------------------------------------------------
type ClientGameInfoRspInit struct {
	ID               int      `json:"id"`
	TargetController int      `json:"targetController"`
	Content          *Content `json:"content"`
}
type RoundsInfo struct {
	MaxMultiplier float64 `json:"maxMultiplier"`
	RoundID       int64   `json:"roundId"`
}
type RoundInfo struct {
	RoundID        int64   `json:"roundId"`
	RoundStartDate int64   `json:"roundStartDate"`
	RoundEndDate   int64   `json:"roundEndDate"`
	Multiplier     float64 `json:"multiplier"`
}
type RoundCount struct {
	RoundCount int64 `json:"roundCount"`
}
type ChatSettings struct {
	Top10Gifs string `json:"top10Gifs"`
}
type Likes struct {
	UsersWhoLike     []string `json:"usersWhoLike"`
	IsMeLiked        bool     `json:"isMeLiked"`
	UsersLikesNumber int      `json:"usersLikesNumber"`
}
type ChatHistory struct {
	MessageType  string  `json:"messageType"`
	MessageID    int     `json:"messageId"`
	WinInfo      WinInfo `json:"winInfo"`
	ProfileImage string  `json:"profileImage"`
	Message      string  `json:"message"`
	PlayerName   string  `json:"playerName"`
	PlayerID     int     `json:"PlayerId"`
	Likes        *Likes  `json:"likes"`
	GifInfo      string  `json:"gifInfo,omitempty"`
}
type Settings struct {
	Music     bool `json:"music"`
	Sound     bool `json:"sound"`
	SecondBet bool `json:"secondBet"`
	Animation bool `json:"animation"`
}
type User struct {
	UserID       string    `json:"userId"`
	Username     string    `json:"username"`
	Balance      float64   `json:"balance"`
	ProfileImage string    `json:"profileImage"`
	Settings     *Settings `json:"settings"`
}
type AutoBetOptions struct {
	DecreaseOrExceedStopPointReq bool      `json:"decreaseOrExceedStopPointReq"`
	NumberOfRounds               []float64 `json:"numberOfRounds"`
}
type AutoCashOut struct {
	MinValue     float64 `json:"minValue"`
	DefaultValue float64 `json:"defaultValue"`
	MaxValue     int     `json:"maxValue"`
}
type Promo struct {
	IsEnabled bool `json:"isEnabled"`
}
type Rain struct {
	IsEnabled         bool    `json:"isEnabled"`
	RainMinBet        float64 `json:"rainMinBet"`
	DefaultNumOfUsers int     `json:"defaultNumOfUsers"`
	MinNumOfUsers     int     `json:"minNumOfUsers"`
	MaxNumOfUsers     int     `json:"maxNumOfUsers"`
	RainMaxBet        int     `json:"rainMaxBet"`
}
type Chat struct {
	Promo            *Promo `json:"promo"`
	Rain             *Rain  `json:"rain"`
	IsEnabled        bool   `json:"isEnabled"`
	IsGifsEnabled    bool   `json:"isGifsEnabled"`
	MaxMessageLength int    `json:"maxMessageLength"`
	MaxMessages      int    `json:"maxMessages"`
	SendMessageDelay int    `json:"sendMessageDelay"`
}
type EngagementTools struct {
}
type NewChat struct {
	Promo            interface{} `json:"promo"`
	Rain             interface{} `json:"rain"`
	IsEnabled        bool        `json:"isEnabled"`
	IsGifsEnabled    bool        `json:"isGifsEnabled"`
	MaxMessageLength int         `json:"maxMessageLength"`
	MaxMessages      int         `json:"maxMessages"`
	SendMessageDelay int         `json:"sendMessageDelay"`
}
type Config struct {
	IsAutoBetFeatureEnabled          bool             `json:"isAutoBetFeatureEnabled"`
	BetPrecision                     int              `json:"betPrecision"`
	MaxBet                           float64          `json:"maxBet"`
	MultiplierPrecision              int              `json:"multiplierPrecision"`
	AccountHistoryActionType         string           `json:"accountHistoryActionType"`
	AutoBetOptions                   *AutoBetOptions  `json:"autoBetOptions"`
	AutoCashOut                      *AutoCashOut     `json:"autoCashOut"`
	BackToHomeActionType             string           `json:"backToHomeActionType"`
	BetInputStep                     float64          `json:"betInputStep"`
	BetOptions                       []float64        `json:"betOptions"`
	Chat                             *Chat            `json:"chat"`
	Currency                         string           `json:"currency"`
	DefaultBetValue                  float64          `json:"defaultBetValue"`
	EngagementTools                  *EngagementTools `json:"engagementTools"`
	FullBetTime                      int              `json:"fullBetTime"`
	GameRulesAutoCashOutType         string           `json:"gameRulesAutoCashOutType"`
	InactivityTimeForDisconnect      int              `json:"inactivityTimeForDisconnect"`
	IrcDisplayType                   string           `json:"ircDisplayType"`
	IsActiveGameFocused              bool             `json:"isActiveGameFocused"`
	IsAlderneyModalShownOnInit       bool             `json:"isAlderneyModalShownOnInit"`
	IsBalanceValidationEnabled       bool             `json:"isBalanceValidationEnabled"`
	IsBetsHistoryStartBalanceEnabled bool             `json:"isBetsHistoryStartBalanceEnabled"`
	IsClockVisible                   bool             `json:"isClockVisible"`
	IsCurrencyNameHidden             bool             `json:"isCurrencyNameHidden"`
	IsFreeBetsEnabled                bool             `json:"isFreeBetsEnabled"`
	IsGameRulesHaveMaxWin            bool             `json:"isGameRulesHaveMaxWin"`
	IsGameRulesHaveMinimumBankValue  bool             `json:"isGameRulesHaveMinimumBankValue"`
	IsHolidayTheme                   bool             `json:"isHolidayTheme"`
	IsLiveBetsAndStatisticsHidden    bool             `json:"isLiveBetsAndStatisticsHidden"`
	IsLoginTimer                     bool             `json:"isLoginTimer"`
	IsMaxUserMultiplierEnabled       bool             `json:"isMaxUserMultiplierEnabled"`
	IsMultipleBetsEnabled            bool             `json:"isMultipleBetsEnabled"`
	IsNetSessionEnabled              bool             `json:"isNetSessionEnabled"`
	IsUseMaskedUsername              bool             `json:"isUseMaskedUsername"`
	MaxUserWin                       int              `json:"maxUserWin"`
	MinBet                           float64          `json:"minBet"`
	OnLockUIActions                  string           `json:"onLockUIActions"`
	QuickHelp                        bool             `json:"quickHelp"`
	ReturnToPlayer                   int              `json:"returnToPlayer"`
	NewChat                          *NewChat         `json:"newChat"`
	PingIntervalMs                   int              `json:"pingIntervalMs"`
	HideHeader                       bool             `json:"hideHeader"`
}

type WinInfo struct {
	Bet                float64 `json:"bet"`
	WinAmount          float64 `json:"winAmount"`
	PlayerName         string  `json:"playerName"`
	Multiplier         float64 `json:"multiplier"`
	Currency           string  `json:"currency"`
	ProfileImage       string  `json:"profileImage"`
	RoundId            int     `json:"roundId"`
	RoundMaxMultiplier float64 `json:"roundMaxMultiplier"`
	CashOutDate        int64   `json:"cashOutDate"`
	IsFreeBet          bool    `json:"isFreeBet"`
}

type TopGif struct {
	Id   string
	Dims []int32
	Url  string
}

func NewSettings() *Settings {
	return &Settings{
		Music:     true,
		Sound:     true,
		SecondBet: true,
		Animation: true,
	}
}

func getDefaultGif() []TopGif {
	res := make([]TopGif, 0)
	gifs := make([]interface{}, 0)
	err := json.Unmarshal([]byte(Top10Gifs), &gifs)
	if err != nil {
		panic(err)
	}
	for _, gif := range gifs {
		o := gif.(map[string]interface{})
		media := o["media"].([]interface{})
		g := media[0].(map[string]interface{})
		gg := g["gif"].(map[string]interface{})
		dims := make([]int32, 0)
		for _, val := range gg["dims"].([]interface{}) {
			dims = append(dims, int32(val.(float64)))
		}
		temp := TopGif{
			Id:   o["id"].(string),
			Dims: dims,
			Url:  gg["url"].(string),
		}
		res = append(res, temp)
	}

	return res
}

var configg = `{
    "isAutoBetFeatureEnabled": true,
    "betPrecision": 2,
    "maxBet": 1000,
    "multiplierPrecision": 2,
    "accountHistoryActionType": "navigate",
    "autoBetOptions": {
        "decreaseOrExceedStopPointReq": true,
        "numberOfRounds": [
            10,
            20,
            50,
            100
        ]
    },
    "autoCashOut": {
        "minValue": 1.01,
        "defaultValue": 1.1,
        "maxValue": 100
    },
    "backToHomeActionType": "navigate",
    "betInputStep": 0.1,
    "betOptions": [
        3,
        5,
        10,
        20
    ],
    "chat": {
        "promo": {
            "isEnabled": true
        },
        "rain": {
            "isEnabled": false,
            "rainMinBet": 0.1,
            "defaultNumOfUsers": 5,
            "minNumOfUsers": 3,
            "maxNumOfUsers": 10,
            "rainMaxBet": 10
        },
        "isEnabled": true,
        "isGifsEnabled": true,
        "maxMessageLength": 160,
        "maxMessages": 70,
        "sendMessageDelay": 5000
    },
    "currency": "BRL",
    "defaultBetValue": 3,
    "engagementTools": {},
    "fullBetTime": 5000,
    "gameRulesAutoCashOutType": "default",
    "inactivityTimeForDisconnect": 0,
    "ircDisplayType": "modal",
    "isActiveGameFocused": false,
    "isAlderneyModalShownOnInit": false,
    "isBalanceValidationEnabled": false,
    "isBetsHistoryStartBalanceEnabled": false,
    "isClockVisible": false,
    "isCurrencyNameHidden": false,
    "isFreeBetsEnabled": false,
    "isGameRulesHaveMaxWin": false,
    "isGameRulesHaveMinimumBankValue": false,
    "isHolidayTheme": false,
    "isLiveBetsAndStatisticsHidden": false,
    "isLoginTimer": false,
    "isMaxUserMultiplierEnabled": false,
    "isMultipleBetsEnabled": true,
    "isNetSessionEnabled": false,
    "isUseMaskedUsername": true,
    "maxUserWin": 10000,
    "minBet": 10,
    "onLockUIActions": "cancelBet",
    "quickHelp": true,
    "returnToPlayer": 97,
    "newChat": {
        "promo": null,
        "rain": null,
        "isEnabled": false,
        "isGifsEnabled": false,
        "maxMessageLength": 0,
        "maxMessages": 0,
        "sendMessageDelay": 0
    },
    "pingIntervalMs": 15000,
    "hideHeader": false
}`

var Top10Gifs = `[{"created":1.704530758487682E9,"flags":[],"media":[{"mp4":{"duration":3.5,"preview":"https://media.tenor.com/DwH1qQbIoCoAAAAD/aidan-hutchinson-detroit-lions.png","size":608439,"dims":[384,480],"url":"https://media.tenor.com/DwH1qQbIoCoAAAPo/aidan-hutchinson-detroit-lions.mp4"},"tinygif":{"duration":3.6,"preview":"https://media.tenor.com/DwH1qQbIoCoAAAAN/aidan-hutchinson-detroit-lions.png","dims":[165,206],"size":501453,"url":"https://media.tenor.com/DwH1qQbIoCoAAAAM/aidan-hutchinson-detroit-lions.gif"},"gif":{"preview":"https://media.tenor.com/DwH1qQbIoCoAAAAD/aidan-hutchinson-detroit-lions.png","duration":3.5,"size":5675474,"dims":[384,480],"url":"https://media.tenor.com/DwH1qQbIoCoAAAAC/aidan-hutchinson-detroit-lions.gif"}}],"title":"","url":"https://tenor.com/br23CkQNMus.gif","content_description":"a group of football players wearing lions uniforms","tags":[],"shares":1,"bg_color":"","hascaption":false,"composite":null,"content_rating":"","id":"1081415491857719338","source_id":"","hasaudio":false,"h1_title":"","itemurl":"https://tenor.com/view/aidan-hutchinson-detroit-lions-nfl-lions-gif-1081415491857719338"},{"created":1.574978865659721E9,"flags":[],"media":[{"mp4":{"preview":"https://media.tenor.com/QvCr4GMBbeQAAAAD/happy-turkey-day-charlie-brown.png","duration":0.8,"dims":[640,536],"size":90378,"url":"https://media.tenor.com/QvCr4GMBbeQAAAPo/happy-turkey-day-charlie-brown.mp4"},"tinygif":{"preview":"https://media.tenor.com/QvCr4GMBbeQAAAAN/happy-turkey-day-charlie-brown.png","duration":0.8,"size":15978,"dims":[220,184],"url":"https://media.tenor.com/QvCr4GMBbeQAAAAM/happy-turkey-day-charlie-brown.gif"},"gif":{"preview":"https://media.tenor.com/QvCr4GMBbeQAAAAD/happy-turkey-day-charlie-brown.png","duration":0.8,"dims":[498,417],"size":250181,"url":"https://media.tenor.com/QvCr4GMBbeQAAAAC/happy-turkey-day-charlie-brown.gif"}}],"title":"","url":"https://tenor.com/bdXkG.gif","content_description":"a cartoon of charlie brown holding a turkey and snoopy holding a tray with the words happy thanksgiving","tags":[],"shares":1,"bg_color":"","hascaption":false,"composite":null,"content_rating":"","id":"4823544181135863268","source_id":"","hasaudio":false,"h1_title":"","itemurl":"https://tenor.com/view/happy-turkey-day-charlie-brown-turkey-snoopy-happy-thanksgiving-gif-15680328"},{"created":1.700747178875564E9,"flags":[],"media":[{"mp4":{"preview":"https://media.tenor.com/5wWkKuyIQHoAAAAD/happy-thankgiving-thanksgiving.png","duration":0.8,"size":44531,"dims":[352,498],"url":"https://media.tenor.com/5wWkKuyIQHoAAAPo/happy-thankgiving-thanksgiving.mp4"},"tinygif":{"duration":0.8,"preview":"https://media.tenor.com/5wWkKuyIQHoAAAAN/happy-thankgiving-thanksgiving.png","dims":[220,312],"size":96839,"url":"https://media.tenor.com/5wWkKuyIQHoAAAAM/happy-thankgiving-thanksgiving.gif"},"gif":{"preview":"https://media.tenor.com/5wWkKuyIQHoAAAAD/happy-thankgiving-thanksgiving.png","duration":0.8,"dims":[352,498],"size":245347,"url":"https://media.tenor.com/5wWkKuyIQHoAAAAC/happy-thankgiving-thanksgiving.gif"}}],"title":"","url":"https://tenor.com/tZS6ATpsaH8.gif","content_description":"a greeting card that says celebrating the season appreciating friends and family wishing the best to all this beautiful season and happy thanksgiving","tags":[],"shares":1,"bg_color":"","hascaption":false,"composite":null,"content_rating":"","id":"16646892101908840570","source_id":"","hasaudio":false,"h1_title":"","itemurl":"https://tenor.com/view/happy-thankgiving-thanksgiving-happy-turkey-gif-16646892101908840570"},{"created":1.567988595008469E9,"flags":[],"media":[{"mp4":{"duration":4.1,"preview":"https://media.tenor.com/YA0lAlNoED4AAAAD/baltimore-ravens-lamar-jackson.png","dims":[640,554],"size":586725,"url":"https://media.tenor.com/YA0lAlNoED4AAAPo/baltimore-ravens-lamar-jackson.mp4"},"tinygif":{"duration":4.1,"preview":"https://media.tenor.com/YA0lAlNoED4AAAAN/baltimore-ravens-lamar-jackson.png","size":169187,"dims":[220,190],"url":"https://media.tenor.com/YA0lAlNoED4AAAAM/baltimore-ravens-lamar-jackson.gif"},"gif":{"preview":"https://media.tenor.com/YA0lAlNoED4AAAAD/baltimore-ravens-lamar-jackson.png","duration":4.1,"size":1306918,"dims":[498,431],"url":"https://media.tenor.com/YA0lAlNoED4AAAAC/baltimore-ravens-lamar-jackson.gif"}}],"title":"","url":"https://tenor.com/ba5Mv.gif","content_description":"a man in a purple jersey with the number 8 on it is dancing","tags":[],"shares":1,"bg_color":"","hascaption":false,"composite":null,"content_rating":"","id":"6921228894257811518","source_id":"","hasaudio":false,"h1_title":"","itemurl":"https://tenor.com/view/baltimore-ravens-lamar-jackson-football-ravens-gif-14997821"},{"created":1.723000074425894E9,"flags":[],"media":[{"mp4":{"preview":"https://media.tenor.com/QZGbU2xJE7AAAAAD/da-bears-caleb-williams.png","duration":1.8,"dims":[498,358],"size":162057,"url":"https://media.tenor.com/QZGbU2xJE7AAAAPo/da-bears-caleb-williams.mp4"},"tinygif":{"duration":1.8,"preview":"https://media.tenor.com/QZGbU2xJE7AAAAAN/da-bears-caleb-williams.png","dims":[220,158],"size":232087,"url":"https://media.tenor.com/QZGbU2xJE7AAAAAM/da-bears-caleb-williams.gif"},"gif":{"duration":1.8,"preview":"https://media.tenor.com/QZGbU2xJE7AAAAAD/da-bears-caleb-williams.png","size":2301602,"dims":[498,358],"url":"https://media.tenor.com/QZGbU2xJE7AAAAAC/da-bears-caleb-williams.gif"}}],"title":"","url":"https://tenor.com/fNbs6P84RJQ.gif","content_description":"a man wearing a denver broncos hat stands in front of a nfl logo","tags":[],"shares":1,"bg_color":"","hascaption":false,"composite":null,"content_rating":"","id":"4724728266689680304","source_id":"","hasaudio":false,"h1_title":"","itemurl":"https://tenor.com/view/da-bears-caleb-williams-chicago-bears-gif-4724728266689680304"},{"created":1.660532621099304E9,"flags":[],"media":[{"mp4":{"preview":"https://media.tenor.com/5Xgt3Phtx64AAAAD/thank-you-sticker-thanks-sticker.png","duration":2,"size":283062,"dims":[640,548],"url":"https://media.tenor.com/5Xgt3Phtx64AAAPo/thank-you-sticker-thanks-sticker.mp4"},"tinygif":{"preview":"https://media.tenor.com/5Xgt3Phtx64AAAAN/thank-you-sticker-thanks-sticker.png","duration":2,"size":59841,"dims":[220,188],"url":"https://media.tenor.com/5Xgt3Phtx64AAAAM/thank-you-sticker-thanks-sticker.gif"},"gif":{"preview":"https://media.tenor.com/5Xgt3Phtx64AAAAD/thank-you-sticker-thanks-sticker.png","duration":2,"size":962795,"dims":[498,426],"url":"https://media.tenor.com/5Xgt3Phtx64AAAAC/thank-you-sticker-thanks-sticker.gif"}}],"title":"","url":"https://tenor.com/bXfXr.gif","content_description":"a cartoon cat is holding a thank you sign","tags":[],"shares":1,"bg_color":"","hascaption":false,"composite":null,"content_rating":"","id":"16535016458974775214","source_id":"","hasaudio":false,"h1_title":"","itemurl":"https://tenor.com/view/thank-you-sticker-thanks-sticker-line-sticker-cat-sticker-orange-cat-gif-26476683"},{"created":1.690894985682449E9,"flags":[],"media":[{"mp4":{"duration":0.9,"preview":"https://media.tenor.com/NmvXJsE5-xsAAAAD/thank-you-grateful.png","dims":[480,480],"size":71003,"url":"https://media.tenor.com/NmvXJsE5-xsAAAPo/thank-you-grateful.mp4"},"tinygif":{"duration":0.9,"preview":"https://media.tenor.com/NmvXJsE5-xsAAAAN/thank-you-grateful.png","size":88482,"dims":[220,220],"url":"https://media.tenor.com/NmvXJsE5-xsAAAAM/thank-you-grateful.gif"},"gif":{"duration":0.9,"preview":"https://media.tenor.com/NmvXJsE5-xsAAAAD/thank-you-grateful.png","dims":[480,480],"size":1283586,"url":"https://media.tenor.com/NmvXJsE5-xsAAAAC/thank-you-grateful.gif"}}],"title":"","url":"https://tenor.com/ePQvM6q7uCB.gif","content_description":"a neon sign that says thank you with a lucas & friends logo","tags":[],"shares":1,"bg_color":"","hascaption":false,"composite":null,"content_rating":"","id":"3921464462006680347","source_id":"","hasaudio":false,"h1_title":"","itemurl":"https://tenor.com/view/thank-you-grateful-reaction-thank-you-very-much-muchas-gracias-gif-3921464462006680347"},{"created":1.608064913506202E9,"flags":[],"media":[{"mp4":{"duration":2,"preview":"https://media.tenor.com/5q1Ly0am-MAAAAAD/mmmkay-mr-mackey.png","dims":[640,360],"size":167889,"url":"https://media.tenor.com/5q1Ly0am-MAAAAPo/mmmkay-mr-mackey.mp4"},"tinygif":{"duration":2,"preview":"https://media.tenor.com/5q1Ly0am-MAAAAAN/mmmkay-mr-mackey.png","size":26780,"dims":[220,124],"url":"https://media.tenor.com/5q1Ly0am-MAAAAAM/mmmkay-mr-mackey.gif"},"gif":{"preview":"https://media.tenor.com/5q1Ly0am-MAAAAAD/mmmkay-mr-mackey.png","duration":2,"dims":[498,280],"size":2287181,"url":"https://media.tenor.com/5q1Ly0am-MAAAAAC/mmmkay-mr-mackey.gif"}}],"title":"","url":"https://tenor.com/bujU3.gif","content_description":"a cartoon character from south park sits at a desk with a computer","tags":[],"shares":1,"bg_color":"","hascaption":false,"composite":null,"content_rating":"","id":"16622025136130160832","source_id":"","hasaudio":false,"h1_title":"","itemurl":"https://tenor.com/view/mmmkay-mr-mackey-south-park-alright-okay-gif-19580399"},{"created":1.69916244278455E9,"flags":[],"media":[{"mp4":{"duration":1.7,"preview":"https://media.tenor.com/6nto_FBZrmwAAAAD/sunday-blessings.png","dims":[498,498],"size":281242,"url":"https://media.tenor.com/6nto_FBZrmwAAAPo/sunday-blessings.mp4"},"tinygif":{"preview":"https://media.tenor.com/6nto_FBZrmwAAAAN/sunday-blessings.png","duration":1.7,"dims":[220,220],"size":412783,"url":"https://media.tenor.com/6nto_FBZrmwAAAAM/sunday-blessings.gif"},"gif":{"preview":"https://media.tenor.com/6nto_FBZrmwAAAAD/sunday-blessings.png","duration":1.7,"size":2200176,"dims":[498,498],"url":"https://media.tenor.com/6nto_FBZrmwAAAAC/sunday-blessings.gif"}}],"title":"","url":"https://tenor.com/uii0af0N1Qe.gif","content_description":"a picture of a tree with leaves and the words sunday blessings on it","tags":[],"shares":1,"bg_color":"","hascaption":false,"composite":null,"content_rating":"","id":"16896213859899649644","source_id":"","hasaudio":false,"h1_title":"","itemurl":"https://tenor.com/view/sunday-blessings-gif-16896213859899649644"},{"created":1.699114937477244E9,"flags":[],"media":[{"mp4":{"duration":0.9,"preview":"https://media.tenor.com/YNdy7jIDu9QAAAAD/happy-sunday.png","size":24083,"dims":[372,480],"url":"https://media.tenor.com/YNdy7jIDu9QAAAPo/happy-sunday.mp4"},"tinygif":{"duration":1,"preview":"https://media.tenor.com/YNdy7jIDu9QAAAAN/happy-sunday.png","size":34781,"dims":[220,284],"url":"https://media.tenor.com/YNdy7jIDu9QAAAAM/happy-sunday.gif"},"gif":{"preview":"https://media.tenor.com/YNdy7jIDu9QAAAAD/happy-sunday.png","duration":0.9,"dims":[372,480],"size":110758,"url":"https://media.tenor.com/YNdy7jIDu9QAAAAC/happy-sunday.gif"}}],"title":"","url":"https://tenor.com/itEgCwrD5oK.gif","content_description":"a tweety bird is sitting next to a cup of coffee and says \" happy sunday \"","tags":[],"shares":1,"bg_color":"","hascaption":false,"composite":null,"content_rating":"","id":"6978172515000761300","source_id":"","hasaudio":false,"h1_title":"","itemurl":"https://tenor.com/view/happy-sunday-gif-6978172515000761300"}]`

func GetChatHistory(chatHistory *ut.SFSArray, ochatHistory []*ChatHistory) *ut.SFSArray {
	for _, val := range ochatHistory {
		temp := ut.NewSFSObject()
		if val.MessageType == "gif" {
			var gif struct {
				Id       string  `json:"id"`
				Preview  string  `json:"preview"`
				Duration float64 `json:"duration"`
				Size     int64   `json:"size"`
				Dims     []int32 `json:"dims"`
				Url      string  `json:"url"`
			}
			json.Unmarshal([]byte(val.GifInfo), &gif)
			if len(gif.Dims) > 0 {
				gifInfo := ut.NewSFSObject()
				dims := ut.NewSFSArray()
				dims.Add(gif.Dims[0], ut.INT, true) //?
				dims.Add(gif.Dims[1], ut.INT, true)
				gifInfo.PutSFSArray("dims", dims)
				gifInfo.PutString("id", gif.Id)
				gifInfo.PutString("url", gif.Url)
				temp.PutSFSObject("gifInfo", gifInfo)
			}
		} else if val.MessageType == "win_info" {
			winInfo := ut.NewSFSObject()
			winInfo.PutDouble("bet", val.WinInfo.Bet)
			winInfo.PutString("playerName", val.PlayerName)
			winInfo.PutDouble("multiplier", val.WinInfo.Multiplier)
			winInfo.PutString("currency", val.WinInfo.Currency)
			winInfo.PutLong("cashOutDate", val.WinInfo.CashOutDate)
			winInfo.PutDouble("winAmount", val.WinInfo.WinAmount)
			winInfo.PutBool("isFreeBet", val.WinInfo.IsFreeBet)
			winInfo.PutString("profileImage", val.WinInfo.ProfileImage)
			winInfo.PutInt("roundId", int32(val.WinInfo.RoundId))
			winInfo.PutDouble("roundMaxMultiplier", val.WinInfo.RoundMaxMultiplier)
			temp.PutSFSObject("winInfo", winInfo)
		}
		temp.PutString("profileImage", val.ProfileImage)
		temp.PutString("messageType", val.MessageType)
		temp.PutInt("messageId", int32(val.MessageID))
		temp.PutString("message", val.Message)
		temp.PutString("playerName", val.PlayerName)
		temp.PutInt("playerId", int32(val.PlayerID)) //?

		if val.Likes != nil && val.Likes.UsersLikesNumber > 0 {
			likes := ut.NewSFSObject()
			usersWhoLike := ut.NewSFSArray()
			for _, value := range val.Likes.UsersWhoLike {
				usersWhoLike.Add(value, ut.UTF_STRING, true)
			}
			likes.PutSFSArray("usersWhoLike", usersWhoLike)
			likes.PutBool("isMeLiked", val.Likes.IsMeLiked)
			likes.PutInt("usersLikesNumber", int32(val.Likes.UsersLikesNumber))
			temp.PutSFSObject("likes", likes)
		}
		chatHistory.Add(temp, ut.SFS_OBJECT, true)
	}
	return chatHistory
}
