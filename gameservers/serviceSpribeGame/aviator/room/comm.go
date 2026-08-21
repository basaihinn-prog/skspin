package room

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"log/slog"
	"math/rand"
	"serve/comm/redisx"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"sort"
	"strconv"
	"time"
)

type CommBody struct {
	Id               int      `json:"id"`
	Channel          string   `json:"channel,omitempty"`
	TargetController int      `json:"targetController"`
	Content          *Content `json:"content"`
}

type Content struct {
	Api string `json:"api,omitempty"`
	Cl  string `json:"cl,omitempty"`
	Ct  int    `json:"ct,omitempty"`
	Ms  int    `json:"ms,omitempty"`
	Tk  string `json:"tk,omitempty"`

	Zn string `json:"zn,omitempty"`
	Un string `json:"un,omitempty"`
	Pw string `json:"pw,omitempty"`
	P  P      `json:"p,omitempty"` //需要调查有没有同名但是字段不同的

	Pi int             `json:"pi,omitempty"`
	Rl [][]interface{} `json:"rl,omitempty"`
	Rs int             `json:"rs,omitempty"`

	C string `json:"c,omitempty"`
	R int    `json:"r,omitempty"`
}

type P struct { //猜测是通用字段，内部包含很多字段
	Token            string            `json:"token,omitempty"`
	Currency         string            `json:"currency,omitempty"`
	Lang             string            `json:"lang,omitempty"`
	SessionToken     string            `json:"sessionToken,omitempty"`
	Platform         *Platform         `json:"platform,omitempty"`
	Version          string            `json:"version,omitempty"`
	AdditionalParams *AdditionalParams `json:"additionalParams,omitempty"`

	RoundsInfo    []*RoundsInfo  `json:"roundsInfo,omitempty"`
	Code          int            `json:"code,omitempty"`
	ActiveBets    []interface{}  `json:"activeBets,omitempty"`
	OnlinePlayers int            `json:"onlinePlayers,omitempty"`
	ChatSettings  *ChatSettings  `json:"chatSettings,omitempty"`
	ChatHistory   []*ChatHistory `json:"chatHistory,omitempty"`

	Message          string         `json:"message,omitempty"`
	Messages         []*ChatHistory `json:"messages,omitempty"`
	MessageId        int64          `json:"messageId,omitempty"`
	UsersLikesNumber int64          `json:"usersLikesNumber,omitempty"`

	SetLike            bool          `json:"setLike,omitempty"`
	User               *User         `json:"user,omitempty"`
	Config             *Config       `json:"config,omitempty"`
	RoundID            int64         `json:"roundID,omitempty"`
	StageID            int           `json:"stageId,omitempty"`
	ActiveFreeBetsInfo []interface{} `json:"activeFreeBetsInfo,omitempty"`

	Bets      []*Bets     `json:"bets,omitempty"`
	BetsCount int         `json:"betsCount,omitempty"`
	CashOuts  []*CashOuts `json:"cashOuts,omitempty"`

	Bet         float64 `json:"bet,omitempty"`
	ClientSeed  string  `json:"clientSeed,omitempty"`
	BetID       int     `json:"betId,omitempty"`
	FreeBet     bool    `json:"freeBet,omitempty"`
	AutoCashOut float64 `json:"autoCashOut,omitempty"`
	Multiplying float64 `json:"multiplying,omitempty"`

	PlayerID     string `json:"player_id,omitempty"`
	ProfileImage string `json:"profileImage,omitempty"`
	Username     string `json:"username,omitempty"`

	Cashouts    []*Cashouts `json:"cashouts,omitempty"`
	Multiplier  float64     `json:"multiplier,omitempty"`
	OperatorKey string      `json:"operatorKey,omitempty"`

	NewBalance float64 `json:"newBalance,omitempty"`

	NewStateID int `json:"newStateId,omitempty"`
	TimeLeft   int `json:"timeLeft,omitempty"`

	X      float64 `json:"x,omitempty"`
	CrashX float64 `json:"crashX,omitempty"`
	Crash  bool    `json:"crash,omitempty"`

	MaxMultiplier float64 `json:"maxMultiplier,omitempty"`
	RoundId       int64   `json:"roundId,omitempty"`

	LastBetId int32 `json:"lastBetId,omitempty"`

	TopWins      []*TopWin     `json:"topWins,omitempty"`
	TopRounds    []*TopRound   `json:"topRounds,omitempty"`
	Period       string        `json:"period,omitempty"`
	ErrorMessage string        `json:"errorMessage,omitempty"`
	Fairness     *FairnessData `json:"fairness,omitempty"`

	BetTimeLeft       int32   `json:"betTimeLeft,omitempty"`
	FullBetTime       int32   `json:"fullBetTime,omitempty"`
	CurrentMultiplier float64 `json:"currentMultiplier,omitempty"`
}

type Platform struct {
	DeviceInfo string `json:"deviceInfo,omitempty"`
	UserAgent  string `json:"userAgent,omitempty"`
	DeviceType string `json:"deviceType,omitempty"`
}

type AdditionalParams struct {
}

type Bets struct {
	Bet          float64 `json:"bet,omitempty"`
	BetID        int     `json:"betId,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	Payout       float64 `json:"payout,omitempty"`
	PlayerID     string  `json:"player_id,omitempty"`
	ProfileImage string  `json:"profileImage,omitempty"`
	Username     string  `json:"username,omitempty"`
	WinAmount    float64 `json:"winAmount,omitempty"`
	Multiplier   float64 `json:"multiplier,omitempty"`
	BetAmount    float64 `json:"betAmount,omitempty"`

	//历史记录中Bets字段
	//Bet           int     `json:"bet"`
	CashOutDate int64 `json:"cashOutDate,omitempty"`
	//Currency      string  `json:"currency"`
	//Payout        float64 `json:"payout"`
	Profit             float64 `json:"profit,omitempty"` // 利润
	RoundBetId         string  `json:"roundBetId,omitempty"`
	RoundId            int64   `json:"roundId,omitempty"`
	Win                bool    `json:"win,omitempty"`
	MaxMultiplier      float64 `json:"maxMultiplier,omitempty"`
	HistoryFlag        bool    `json:"-"`
	RobotBetFinishFlag bool    `json:"-"`
	Effective          bool    `json:"-"`
}
type FairnessData struct {
	Number                float64      `json:"number"`                // 数值型结果 (如 126.8372)
	PartSeedDecimalNumber int64        `json:"partSeedDecimalNumber"` // 十进制种子数 (需注意 int64 最大值限制)
	PartSeedHexNumber     string       `json:"partSeedHexNumber"`     // 十六进制种子字符串
	PlayerSeeds           []PlayerSeed `json:"playerSeeds"`           // 玩家种子列表
	Result                float64      `json:"result"`                // 最终结果值 (如 49.35)
	RoundID               int64        `json:"roundId"`               // 回合唯一标识
	RoundStartDate        int64        `json:"roundStartDate"`        // 回合开始时间戳 (毫秒级)
	SeedSHA256            string       `json:"seedSHA256"`            // SHA256 哈希值 (注意重复拼接特征)
	ServerSeed            string       `json:"serverSeed"`            // 服务器种子值
}

type FairnessDataBson struct {
	RoomID                string       `json:"roomId" bson:"roomID"`
	Number                float64      `json:"number" bson:"number"`                               // 数值型结果 (如 126.8372)
	PartSeedDecimalNumber int64        `json:"partSeedDecimalNumber" bson:"partSeedDecimalNumber"` // 十进制种子数 (需注意 int64 最大值限制)
	PartSeedHexNumber     string       `json:"partSeedHexNumber" bson:"partSeedHexNumber"`         // 十六进制种子字符串
	PlayerSeeds           []PlayerSeed `json:"playerSeeds" bson:"playerSeeds"`                     // 玩家种子列表
	Result                float64      `json:"result" bson:"result"`                               // 最终结果值 (如 49.35)
	RoundID               int64        `json:"roundId" bson:"roundID"`                             // 回合唯一标识
	RoundStartDate        int64        `json:"roundStartDate" bson:"roundStartDate"`               // 回合开始时间戳 (毫秒级)
	SeedSHA256            string       `json:"seedSHA256" bson:"seedSHA256"`                       // SHA256 哈希值 (注意重复拼接特征)
	ServerSeed            string       `json:"serverSeed" bson:"serverSeed"`                       // 服务器种子值
	CreateAt              time.Time    `json:"createAt" bson:"createAt"`                           // 创建时间
}

type PlayerSeed struct {
	ProfileImage string `json:"profileImage"` // 玩家头像路径 (如 "av-16.png")
	Seed         string `json:"seed"`         // 玩家种子字符串 (含混合字符)
	Username     string `json:"username"`     // 玩家用户名 (如 "d7r48")
}
type CashOuts struct { //currentBetsInfo请求
	BetID         int     `json:"betId,omitempty"`
	PlayerID      string  `json:"player_id,omitempty"`
	Currency      string  `json:"currency,omitempty"`
	Multiplier    float64 `json:"multiplier,omitempty"`
	BetAmount     float64 `json:"betAmount,omitempty"`
	WinAmount     float64 `json:"winAmount,omitempty"`
	MaxMultiplier float64 `json:"maxMultiplier"`
}

type Cashouts struct { //cashOut、updateCurrentCashOuts请求
	BetID         int     `json:"betId,omitempty"`
	PlayerID      string  `json:"player_id,omitempty"`
	Currency      string  `json:"currency,omitempty"`
	Multiplier    float64 `json:"multiplier,omitempty"`
	BetAmount     float64 `json:"betAmount,omitempty"`
	WinAmount     float64 `json:"winAmount,omitempty"`
	Win           bool    `json:"win,omitempty"`
	MaxMultiplier float64 `json:"maxMultiplier"`
}

type TopWin struct {
	Bet                     float64 `json:"bet,omitempty"`
	EndDate                 int64   `json:"endDate,omitempty"`
	Currency                string  `json:"currency,omitempty"`
	MaxMultiplier           float64 `json:"maxMultiplier,omitempty"`
	Payout                  float64 `json:"payout,omitempty"`
	PlayerId                string  `json:"playerId,omitempty"`
	ProfileImage            string  `json:"profileImage,omitempty"`
	RoundBetId              string  `json:"roundBetId,omitempty"`
	RoundId                 int64   `json:"roundId,omitempty"`
	Username                string  `json:"username,omitempty"`
	WinAmountInMainCurrency float64 `json:"winAmountInMainCurrency,omitempty"`
	WinAmount               float64 `json:"winAmount,omitempty"`
	Zone                    string  `json:"zone,omitempty"`
}

type TopRound struct {
	EndDate        int64   `json:"endDate"`
	MaxMultiplier  float64 `json:"maxMultiplier"`
	RoundId        int64   `json:"roundId"`
	RoundStartDate int64   `json:"roundStartDate"`
	ServerSeed     string  `json:"serverSeed"`
	Zone           string  `json:"zone"`
}

var maxLength = int64(20)

func GetErrMsg(c string, code int, errorMessage string, p ...P) CommBody {
	rspp := P{
		Code:         code,
		OperatorKey:  "release",
		ErrorMessage: errorMessage,
	}
	if len(p) > 0 {
		rspp = p[0]
	}
	return CommBody{
		Id:               13,
		TargetController: 1,
		Content: &Content{
			C: c,
			P: rspp,
		},
	}
}

//rsp.Content = &Content{
//C: "bet",
//P: P{
//Cashouts: []*Cashouts{{
//BetID:    req.Content.P.BetID,
//PlayerID: playerIdStr,
//}},
//Code:         comm.ErrTimeout,
//OperatorKey:  "release",
//ErrorMessage: "Stage time out",
//},
//}

// 存储数据
func StoreBet(client *redis.Client, key string, bet Bets) error {
	mapBets := make(map[string]Bets)
	mapBets[bet.PlayerID+"_"+strconv.Itoa(bet.BetID)] = bet
	jsonData, _ := json.Marshal(mapBets)
	return client.ZAdd(context.TODO(), key, &redis.Z{
		Score:  bet.Bet,
		Member: jsonData,
	}).Err()
}

// 查询排序结果
func GetSortedBets(client *redis.Client, key string) ([]*Bets, error) {
	exists, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("检查键是否存在失败: %v", err)
	}
	if exists == 0 {
		return nil, nil
	}
	result, err := client.ZRevRangeWithScores(context.TODO(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	bets := []*Bets{}
	for _, z := range result {
		mapBets := map[string]Bets{}
		if err := json.Unmarshal([]byte(z.Member.(string)), &mapBets); err == nil {
			for _, v := range mapBets {
				bets = append(bets, &v)
			}
		}
	}
	return bets, nil
}

func GetSortedCashOuts(client *redis.Client, key string) ([]*Cashouts, error) {
	result, err := client.ZRevRangeWithScores(context.TODO(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	bets := []*Cashouts{}
	for _, z := range result {
		mapBets := map[string]Cashouts{}
		if err := json.Unmarshal([]byte(z.Member.(string)), &mapBets); err == nil {
			for _, v := range mapBets {
				bets = append(bets, &v)
			}
		}
	}
	return bets, nil
}

// 更新数据
//func UpdateBet(client *redis.Client, key string, bet Bets) error {
//	keyStr := bet.PlayerID + "_" + strconv.Itoa(bet.BetID)
//	iter := client.ZScan(context.TODO(), key, 0, "*\""+keyStr+"\"*", 1).Iterator()
//	var foundMembers string
//	for iter.Next(context.TODO()) {
//		foundMembers = iter.Val()
//		break
//	}
//	betOld := make(map[string]*Bets)
//	json.Unmarshal([]byte(foundMembers), &betOld)
//	betOld[keyStr].WinAmount = bet.WinAmount
//	betOld[keyStr].Multiplier = bet.Multiplier
//	betOld[keyStr].Payout = bet.Multiplier
//	updatedJSON, _ := json.Marshal(betOld)
//	tx := client.TxPipeline()
//	tx.ZRem(context.TODO(), key, foundMembers)
//	tx.ZAdd(context.TODO(), key, &redis.Z{
//		Score:  bet.Bet, // 如果需要更新分数
//		Member: updatedJSON,
//	})
//	_, err := tx.Exec(context.TODO())
//	if err != nil {
//		return err
//	}
//	return err
//}

func GetBet(key string) (*Bets, error) {
	client := redisx.GetClient()
	result, err := client.ZRange(context.TODO(), key, 0, 0).Result()
	if err == nil && len(result) > 0 {
		firstElement := result[0]
		// 反序列化并更新赌注信息
		bet := make(map[string]*Bets)
		if err := json.Unmarshal([]byte(firstElement), &bet); err != nil {
			return nil, fmt.Errorf("failed to unmarshal bet: %v", err)
		}
		for _, v := range bet {
			return v, nil
		}
	}
	return nil, nil
}

func UpdateBet(client *redis.Client, key string, bet Bets) error {
	keyStr := bet.PlayerID + "_" + strconv.Itoa(bet.BetID)
	memberPattern := "*\"" + keyStr + "\"*"

	// 查找目标成员
	iter := client.ZScan(context.TODO(), key, 0, memberPattern, 0).Iterator()
	var foundMembers string
	for iter.Next(context.TODO()) {
		foundMembers = iter.Val()
		break
	}
	if foundMembers == "" {
		return fmt.Errorf("bet not found: %s", keyStr)
	}

	// 反序列化并更新赌注信息
	betOld := make(map[string]*Bets)
	if err := json.Unmarshal([]byte(foundMembers), &betOld); err != nil {
		return fmt.Errorf("failed to unmarshal bet: %v", err)
	}
	betOld[keyStr].WinAmount = bet.WinAmount
	betOld[keyStr].Multiplier = bet.Multiplier
	betOld[keyStr].Payout = bet.Multiplier

	// 序列化更新后的赌注信息
	updatedJSON, err := json.Marshal(betOld)
	if err != nil {
		return fmt.Errorf("failed to marshal updated bet: %v", err)
	}

	// 使用 Redis 事务更新数据
	tx := client.TxPipeline()
	tx.ZRem(context.TODO(), key, foundMembers)
	tx.ZAdd(context.TODO(), key, &redis.Z{
		Score:  bet.Bet, // 如果需要更新分数
		Member: updatedJSON,
	})
	if _, err := tx.Exec(context.TODO()); err != nil {
		return fmt.Errorf("failed to execute transaction: %v", err)
	}
	//slog.Info(fmt.Sprintf("这个key更新了 %s", keyStr))

	return nil
}

//func UpdateBet(client *redis.Client, key string, bet Bets) (float64,error) {
//	keyStr := bet.PlayerID + "_" + strconv.Itoa(bet.BetID)
//	memberPattern := "*\"" + keyStr + "\"*"
//
//	// 查找目标成员
//	iter := client.ZScan(context.TODO(), key, 0, memberPattern, 0).Iterator()
//	var foundMembers string
//	for iter.Next(context.TODO()) {
//		foundMembers = iter.Val()
//		break
//	}
//	if foundMembers == "" {
//		return 0, fmt.Errorf("bet not found: %s", keyStr)
//	}
//
//	// 反序列化并更新赌注信息
//	betOld := make(map[string]*Bets)
//	if err := json.Unmarshal([]byte(foundMembers), &betOld); err != nil {
//		return 0, fmt.Errorf("failed to unmarshal bet: %v", err)
//	}
//	betOld[keyStr].WinAmount = bet.WinAmount
//	betOld[keyStr].Multiplier = bet.Multiplier
//	betOld[keyStr].Payout = bet.Multiplier
//
//	// 序列化更新后的赌注信息
//	updatedJSON, err := json.Marshal(betOld)
//	if err != nil {
//		return 0, fmt.Errorf("failed to marshal updated bet: %v", err)
//	}
//
//	// 使用 Redis 事务更新数据
//	tx := client.TxPipeline()
//	tx.ZRem(context.TODO(), key, foundMembers)
//	tx.ZAdd(context.TODO(), key, &redis.Z{
//		Score:  bet.Bet, // 如果需要更新分数
//		Member: updatedJSON,
//	})
//	if _, err := tx.Exec(context.TODO()); err != nil {
//		return 0, fmt.Errorf("failed to execute transaction: %v", err)
//	}
//	//slog.Info(fmt.Sprintf("这个key更新了 %s", keyStr))
//	// 获取更新后的最小值
//	min := 0.0
//	minMember, err := client.ZRangeWithScores(context.TODO(), key, 0, 0).Result()
//	if err != nil {
//		return 0, fmt.Errorf("failed to get min member: %v", err)
//	}
//	if len(minMember) > 0 {
//		tempBet := Bets{}
//		json.Unmarshal([]byte(minMember[0].Member.(string)), &tempBet)
//		min = tempBet.Bet
//	}
//
//	return min, nil
//}

// 删除数据
func DeleteKeyAsync(client *redis.Client, key string) error {
	exists, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		return fmt.Errorf("检查键是否存在失败: %v", err)
	}
	if exists == 0 {
		return nil
	}
	// 异步删除不阻塞服务
	deleted, err := client.Unlink(context.Background(), key).Result()
	if err != nil {
		return fmt.Errorf("异步删除失败: %v", err)
	}
	if deleted == 0 {
		return nil
	}
	return nil
}

// 添加数据到上一手的key里
func AddPreviousBetInfo(client *redis.Client, key string) error {
	ctx := context.Background()
	// 添加数据前先清空数据
	slog.Info(fmt.Sprintf("房间%v插入前先清空数据", key))
	exists, err := client.Exists(ctx, key+"_previous").Result()
	if err != nil {
		slog.Error("插入上一手失败::AddPreviousBetInfo::client.Exists", "Err", err)
		return fmt.Errorf("检查键是否存在失败: %v", err)
	}
	if exists == 1 {
		_, err = client.Del(context.Background(), key+"_previous").Result()
		if err != nil {
			slog.Error("插入上一手失败::AddPreviousBetInfo::client.Del", "Err", err)
			return fmt.Errorf("删除失败: %v", err)
		}
	}
	exists, err = client.Exists(ctx, key).Result()
	if err != nil {
		slog.Error("插入上一手失败::AddPreviousBetInfo::client.Exists", "Err", err)
		return fmt.Errorf("检查键是否存在失败: %v", err)
	}
	if exists == 0 {
		slog.Info(fmt.Sprintf("房间%v没有查询到上一手", key))
		return nil
	}
	result, err := client.ZRevRangeWithScores(context.TODO(), key, 0, -1).Result()
	if err != nil {
		slog.Error("插入上一手失败::AddPreviousBetInfo::client.ZRevRangeWithScores", "Err", err)
		log.Fatalf("获取源键失败: %v", err)
		return nil
	}
	slog.Info(fmt.Sprintf("房间%v插入上一手查询结果", key), "result.len", len(result))
	if len(result) > 0 {
		for _, z := range result {
			err = client.ZAdd(ctx, key+"_previous", &redis.Z{
				Score:  z.Score,
				Member: z.Member,
			}).Err()
			if err != nil {
				slog.Error("插入上一手失败::AddPreviousBetInfo::client.ZAdd", "Err", err)
				log.Fatalf("添加元素到目标有序集合失败: %v", err)
			}
		}
	}
	return err
}

// 添加数据到上一手的回合信息
func AddPreviousRoundInfo(key string, info *RoundInfo) error {
	ctx := context.Background()
	client := redisx.GetClient()
	jsonStore, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = client.SetEX(ctx, key+"_previous:"+"roundInfo", jsonStore, time.Second*60*60*24).Err()
	if err != nil {
		slog.Error("插入上一手失败::addPreviousRoundInfo::client.Exists", "Err", err)
		return fmt.Errorf("检查键是否存在失败: %v", err)
	}

	return err
}

// 获取上一回合信息
func GetPreviousRoundInfo(key string) (*RoundInfo, error) {
	ctx := context.Background()
	client := redisx.GetClient()
	exists, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, nil
	}
	jsonStore, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var info RoundInfo
	if err := json.Unmarshal([]byte(jsonStore), &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func IncrementKey(client *redis.Client, key string) (int64, error) {
	return client.Incr(context.TODO(), key).Result()
}

func StoreMessage(client *redis.Client, key string, Score float64, message []byte) error {
	return client.ZAdd(context.TODO(), key, &redis.Z{
		Score:  Score,
		Member: message,
	}).Err()
}

func GetMessage(client *redis.Client, key string) ([]*ChatHistory, error) {
	exists, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("检查键是否存在失败: %v", err)
	}
	if exists == 0 {
		return nil, nil
	}
	result, err := client.ZRangeWithScores(context.TODO(), key, 0, -1).Result()

	chatHistorys := []*ChatHistory{}
	for _, z := range result {
		commBody := ChatHistory{}
		if err := json.Unmarshal([]byte(z.Member.(string)), &commBody); err == nil {
			chatHistorys = append(chatHistorys, &commBody)
		}
	}
	return chatHistorys, nil
}

func GetOneMessage(client *redis.Client, key string, index int64) (ChatHistory, error) {
	exists, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		return ChatHistory{}, fmt.Errorf("检查键是否存在失败: %v", err)
	}
	if exists == 0 {
		return ChatHistory{}, nil
	}

	result, err := client.ZRange(context.TODO(), key, index-1, index-1).Result()
	if err != nil {
		return ChatHistory{}, err
	}

	commBody := ChatHistory{}
	err = json.Unmarshal([]byte(result[0]), &commBody)
	if err != nil {
		return ChatHistory{}, err
	}

	return commBody, nil
}

func UpdateMessage(client *redis.Client, key string, Score float64, message []byte) error {

	script := `
    local key = KEYS[1]
    local score = tonumber(ARGV[1])
	local rank = score - 1
    local message = ARGV[2]
    redis.call('ZREMRANGEBYRANK', key, rank, rank)
    redis.call('ZADD', key, score, message)
    return 1
    `

	_, err := client.Eval(context.TODO(), script, []string{key}, Score, string(message)).Result()
	return err
}

// 存储数据最高奖金
func StoreTopWins(client *redis.Client, key string, win TopWin) error {
	// 添加新元素到有序集合
	jsonData, _ := json.Marshal(win)
	ctx := context.Background()
	err := client.ZAdd(ctx, key, &redis.Z{
		Score:  win.WinAmount,
		Member: jsonData,
	}).Err()
	if err != nil {
		log.Fatalf("添加元素失败: %v", err)
	}

	// 检查有序集合的当前长度
	currentLength, err := client.ZCard(ctx, key).Result()
	if err != nil {
		log.Fatalf("获取有序集合长度失败: %v", err)
	}

	// 如果长度超过最大限制，移除分数最低的元素
	if currentLength > maxLength {
		_, err = client.ZRemRangeByRank(ctx, key, 0, 0).Result() // 移除分数最低的元素
		if err != nil {
			log.Fatalf("移除分数最低的元素失败: %v", err)
		}
		//fmt.Printf("超出长度限制，已移除分数最低的元素\n")
	}

	return err
}

// 存储大额奖金
func StoreHugeWins(client *redis.Client, key string, win TopWin) error {
	// 添加新元素到有序集合
	jsonData, _ := json.Marshal(win)
	ctx := context.Background()
	err := client.ZAdd(ctx, key, &redis.Z{
		Score:  win.Payout,
		Member: jsonData,
	}).Err()
	if err != nil {
		log.Fatalf("添加元素失败: %v", err)
	}

	// 检查有序集合的当前长度
	currentLength, err := client.ZCard(ctx, key).Result()
	if err != nil {
		log.Fatalf("获取有序集合长度失败: %v", err)
	}

	// 如果长度超过最大限制，移除分数最低的元素
	if currentLength > maxLength {
		_, err = client.ZRemRangeByRank(ctx, key, 0, 0).Result() // 移除分数最低的元素
		if err != nil {
			log.Fatalf("移除分数最低的元素失败: %v", err)
		}
		//fmt.Printf("超出长度限制，已移除分数最低的元素\n")
	}

	return err
}

// 查询排序结果
func GetSortedTopWins(client *redis.Client, key string) ([]*TopWin, error) {
	exists, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("检查键是否存在失败: %v", err)
	}
	if exists == 0 {
		return nil, nil
	}
	result, err := client.ZRevRangeWithScores(context.TODO(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	topWins := []*TopWin{}
	for _, z := range result {
		topWin := TopWin{}
		if err := json.Unmarshal([]byte(z.Member.(string)), &topWin); err == nil {
			topWins = append(topWins, &topWin)
		}
	}
	return topWins, nil
}

// 查询排序结果
func GetSortedTopRounds(client *redis.Client, key string) ([]*TopRound, error) {
	exists, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("检查键是否存在失败: %v", err)
	}
	if exists == 0 {
		return nil, nil
	}
	result, err := client.ZRevRangeWithScores(context.TODO(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	topRounds := []*TopRound{}
	for _, z := range result {
		topRound := TopRound{}
		if err := json.Unmarshal([]byte(z.Member.(string)), &topRound); err == nil {
			topRounds = append(topRounds, &topRound)
		}
	}
	return topRounds, nil
}

// 存储大额奖金
func StoreTopRounds(client *redis.Client, key string, round TopRound) error {
	// 添加新元素到有序集合
	jsonData, _ := json.Marshal(round)
	ctx := context.Background()
	err := client.ZAdd(ctx, key, &redis.Z{
		Score:  round.MaxMultiplier,
		Member: jsonData,
	}).Err()
	if err != nil {
		log.Fatalf("添加元素失败: %v", err)
	}

	// 检查有序集合的当前长度
	currentLength, err := client.ZCard(ctx, key).Result()
	if err != nil {
		log.Fatalf("获取有序集合长度失败: %v", err)
	}

	// 如果长度超过最大限制，移除分数最低的元素
	if currentLength > maxLength {
		_, err = client.ZRemRangeByRank(ctx, key, 0, 0).Result() // 移除分数最低的元素
		if err != nil {
			log.Fatalf("移除分数最低的元素失败: %v", err)
		}
		//fmt.Printf("超出长度限制，已移除分数最低的元素\n")
	}

	return err
}

func DeleteTopData(rooms []string, date string) {
	for _, room := range rooms {
		keyHugeWin := fmt.Sprintf("%s_%s_%s", room, "huge", date)
		err := DeleteKeyAsync(redisx.GetClient(), keyHugeWin)
		if err != nil {
			slog.Error(fmt.Sprintf("DeleteTopData Err"), "keyHugeWin", keyHugeWin, "err", err)
		}
		keyTopWin := fmt.Sprintf("%s_%s_%s", room, "topwin", date)
		err = DeleteKeyAsync(redisx.GetClient(), keyTopWin)
		if err != nil {
			slog.Error(fmt.Sprintf("DeleteTopData Err"), "keyTopWin", keyTopWin, "err", err)
		}
		keyTopRound := fmt.Sprintf("%s_%s_%s", room, "topround", date)
		err = DeleteKeyAsync(redisx.GetClient(), keyTopRound)
		if err != nil {
			slog.Error(fmt.Sprintf("DeleteTopData Err"), "keyTopRound", keyTopRound, "err", err)
		}
	}
}

// 同步房间的状态到redis中
func RoomBackup2Redis(r *Room) error {
	var err error
	client := redisx.GetClient()

	jsonData, err := json.Marshal(r)
	if err != nil {
		slog.Error(fmt.Sprintf("RoomBackup2Redis: Marshal Err"), "err", err)
	}
	// 存储到 Redis
	roomKey := fmt.Sprintf("backup:%s", r.Name)
	err = client.Set(context.TODO(), roomKey, jsonData, 0).Err()
	if err != nil {
		slog.Error("RoomBackup2Redis Err", "roomKey", roomKey, "err", err)
	}
	return nil
}

func GetRoomInfoFromRedis(roomName string) (*Room, error) {
	var room *Room
	var err error
	client := redisx.GetClient()
	roomKey := fmt.Sprintf("backup:%s", roomName)
	result, err := client.Get(context.TODO(), roomKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		slog.Error("GetRoomInfoFromRedis::client.Get Err", "roomKey", roomKey, "err", err)
	}
	err = json.Unmarshal([]byte(result), &room)
	if err != nil {
		slog.Error("GetRoomInfoFromRedis::json.Unmarshal Err", "result", result, "err", err)
	}
	return room, nil
}

var ActionMap = map[any]string{
	0:                           "AviatorLoginIdReq",
	1:                           "AviatorGameInfoIdReq",
	"currentBetsInfoHandler":    "AviatorCurrentBetsInfoIdReq",
	"betHandler":                "AviatorBetIdReq",
	"cancelBetHandler":          "AviatorCancelBetIdReq",
	"cashOutHandler":            "AviatorCashOutIdReq",
	"betHistoryHandler":         "AviatorBetHistoryIdReq",
	"PING_REQUEST":              "AviatorGameStatePingIdReq",
	"AddChatMessageV2Handler":   "AviatorAddChatMessageIdReq",
	"likeHandler":               "AviatorLikeIdReq",
	"searchGifsHandler":         "ClientSearchGifs",
	"previousRoundInfoHandler":  "AviatorPreviousRoundInfoIdReq",
	"getHugeWinsInfoHandler":    "AviatorGetHugeWinsInfoIdReq",
	"getTopWinsInfoHandler":     "AviatorGetTopWinsInfoIdReq",
	"getTopRoundsInfoHandler":   "AviatorGetTopRoundsInfoIdReq",
	"changeProfileImageHandler": "AviatorChangeProfileImageIdReq",
	"roundFairnessHandler":      "AviatorRoundFairnessIdReq",
	"serverSeedHandler":         "ServerSeedHandler",
}

func FloatProbabilityCheck(percent float64) bool {
	if percent <= 0.0 {
		return false
	}
	if percent >= 100.0 {
		return true
	}
	return rand.Float64()*100 < percent
}

func GetClientErrRsp(initRsp *CommBody) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	pp.PutString("errorMessage", initRsp.Content.P.ErrorMessage)
	switch initRsp.Content.P.Code {
	case comm.ErrNoNotEnoughBalance:
		if initRsp.Content.P.BetID != 0 {
			pp.PutInt("betId", int32(initRsp.Content.P.BetID))
		}
		pp.PutDouble("bet", initRsp.Content.P.Bet)
	}
	if initRsp.Content.P.BetID != 0 {
		pp.PutInt("betId", int32(initRsp.Content.P.BetID))
	}
	if initRsp.Content.P.PlayerID != "" {
		pp.PutString("player_id", initRsp.Content.P.PlayerID)
	}
	if initRsp.Content.P.ProfileImage != "" {
		pp.PutString("profileImage", initRsp.Content.P.ProfileImage)
	}
	if initRsp.Content.P.Username != "" {
		pp.PutString("username", initRsp.Content.P.Username)
	}
	if initRsp.Content.P.OperatorKey != "" {
		pp.PutString("operatorKey", initRsp.Content.P.OperatorKey)
	}
	if len(initRsp.Content.P.Cashouts) != 0 {
		cashOuts := ut.NewSFSArray()
		for _, bet := range initRsp.Content.P.Cashouts {
			temp := ut.NewSFSObject()
			//temp.PutDouble("betAmount", bet.BetAmount)
			//temp.PutDouble("winAmount", bet.WinAmount)
			temp.PutString("player_id", bet.PlayerID)
			temp.PutInt("betId", int32(bet.BetID))
			//temp.PutBool("isMaxWinAutoCashOut", false) //todo
			cashOuts.Add(temp, ut.SFS_OBJECT, true)
		}
		pp.PutSFSArray("cashouts", cashOuts)
	}

	p.PutSFSObject("p", pp)
	p.PutString("c", initRsp.Content.C)

	so.AddCreatePAC(p, 1, 13)
	return so
}

func SetHashOne(client *redis.Client, key string, Score int64, value string) error {
	err := client.HSet(context.Background(), key, Score, value).Err()
	return err
}

func GetHashAll(client *redis.Client, key string) ([]string, error) {
	data, err := client.HGetAll(context.Background(), key).Result()

	var sorts []int
	var list []string
	for k, _ := range data {
		i, _ := strconv.Atoi(k)
		sorts = append(sorts, i)
	}
	sort.Ints(sorts)
	for _, v := range sorts {
		i := data[strconv.Itoa(v)]
		list = append(list, i)
	}

	return list, err
}

func GetHashOne(client *redis.Client, key string, Score string) (string, error) {
	data, err := client.HGet(context.Background(), key, Score).Result()
	return data, err
}
