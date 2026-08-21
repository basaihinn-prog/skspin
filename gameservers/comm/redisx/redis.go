package redisx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/mq"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// Redis 客户端
	redisClient *redis.Client
	serviceName string
)

const (
	Bucket1 = 0
	Bucket2 = 1
	Bucket3 = 2
)

// 初始化 Redis 客户端
func Init(redisURL string) {
	opt, err := redis.ParseURL(redisURL)
	redisClient = redis.NewClient(opt)
	// 测试 Redis 连接
	ctx := context.Background()
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		slog.Error("Redis 连接失败:", err)
		return
	}
	slog.Info("Redis 连接成功:", pong)
}

func GetClient() *redis.Client {
	return redisClient
}

// AppStore 结构体
type AppStore struct {
	RewardPercent      int       `json:"reward_percent"`
	NoAwardPercent     int       `json:"no_award_percent"`
	AppId              string    `json:"appId"`
	GamePatten         int       `json:"gamePatten"`   //是哪个桶子
	MaxWinPoints       float64   `json:"MaxWinPoints"` //最大应分 bat * time < MaxWinPoints
	MaxMultiple        float64   `json:"maxMultiple"`
	GameId             string    `json:"gameId"`
	Cs                 []float64 `json:"cs"`
	StopLoss           int64     `json:"StopLoss"`            //止盈止损开关
	ShowNameAndTimeOff int64     `json:"ShowNameAndTimeOff" ` //名字和时间显示
	ShowExitBtnOff     int64     `json:"ShowExitBtnOff" `     //退出按钮显示
	FreeGameOff        int64     `json:"FreeGameOff"`         //是否开启免费游戏（是否允许游客玩）

	//PG独有
	DefaultCs       float64 `json:"DefaultCs"`       //默认CS
	DefaultBetLevel int64   `json:"DefaultBetLevel"` //默认ML

	//新增
	IsProtection                int `json:"IsProtection"`                //是否开启 0 关 1开
	ProtectionRotateCount       int `json:"ProtectionRotateCount"`       //保护次数
	ProtectionRewardPercentLess int `json:"ProtectionRewardPercentLess"` //保护内获取比率
	BuyMinAwardPercent          int `json:"BuyMinAwardPercent"`          //千分制

	RTP                   int           `json:"RTP"`                   //rtp
	MoniterConfig         MoniterConfig `json:"MoniterConfig"`         //游戏监控
	PersonalMoniterConfig MoniterConfig `json:"PersonalMoniterConfig"` //个人监控
	// 利润率
	ProfitMargin float64 `json:"ProfitMargin"`
	// 坠毁概率
	CrashRate float64 `json:"CrashRate"`
	// 刻度
	Scale         float64 `json:"Scale"`
	OnlineUpNum   float64 `json:"OnlineUpNum"`
	OnlineDownNum float64 `json:"OnlineDownNum"`
}

type MoniterConfig struct {
	IsMoniter                  int64                  `json:"IsMoniter"`                  //游戏监控
	MoniterNewbieNum           int64                  `json:"MoniterNewbieNum"`           //监控新手数量
	MoniterRTPErrorValue       int64                  `json:"MoniterRTPErrorValue"`       //监控RTP错误值
	MoniterNumCycle            int64                  `json:"MoniterNumCycle"`            //监控周期
	MoniterAddRTPRangeValue    []MoniterRtpRangeValue `json:"MoniterAddRTPRangeValue"`    //监控增加RTP范围值
	MoniterReduceRTPRangeValue []MoniterRtpRangeValue `json:"MoniterReduceRTPRangeValue"` //监控减少RTP范围值
}

type MoniterRtpRangeValue struct {
	RangeMinValue  int64   `json:"RangeMinValue"`  //最小值
	RangeMaxValue  int64   `json:"RangeMaxValue"`  //最大值
	NewbieValue    float64 `json:"NewbieValue"`    //新手概率
	NotNewbieValue float64 `json:"NotNewbieValue"` //老手概率
}

// DocPlayer 结构体
type DocPlayer struct {
	Pid int64
}

// PersonRtpSettings 结构体
type PersonRtpSettings struct {
	RewardPercent      int   `json:"reward_percent"`
	NoAwardPercent     int   `json:"no_award_percent"`
	TargetRTP          int   `json:"target_rtp"`
	RelieveRTP         int   `json:"relieve_rtp"`
	PlayerId           int64 `json:"player_id"`
	BuyMinAwardPercent int   `json:"buy_min_award_percent"`
}

// Redis 相关常量
const (
	AppIdCacheKey   = "appIdCache"   // 用于存储 appIdCache 的 Redis 键
	ReceiveCountKey = "receiveCount" // 用于存储 receiveCount 的 Redis 键

	AppStoreKey    = "appStore_"
	PlayerCsKey    = "playerCs_"
	PlayerNextData = "pnd_%d_%s"

	ALLGameBatchSize = 1000000
)

//// 删除 Redis 中的 appIdCache
//func deleteAppIdCache(ctx context.Context, key string) error {
//	return redisClient.Del(ctx, AppIdCacheKey, key).Err()
//}

// 加载 receiveCount
func loadReceiveCount() (int, error) {
	count, err := redisClient.Get(context.Background(), ReceiveCountKey).Int()
	if err == redis.Nil {
		return 0, nil // 如果不存在，返回 0
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 增加 receiveCount
func incrementReceiveCount() error {
	return redisClient.Incr(context.Background(), ReceiveCountKey).Err()
}

// 生成键
func MakeAppKey(appId string) string {
	return fmt.Sprintf("%v:%s", appId, serviceName)
}

// 生成键
func MakePlayerKey(plrId int64) string {
	return fmt.Sprintf("%d:%s", plrId, serviceName)
}

// 生成键
func MakePlayerNextDataKey(plrId int64) string {
	return fmt.Sprintf(PlayerNextData, plrId, serviceName)
}

// 注册订阅
func RegSubscribe(gameId, redisURL string) {
	fmt.Println("register subscribe gameId:", gameId)
	serviceName = gameId
	Init(redisURL)
	mq.Subscribe(fmt.Sprintf("/store/setInfo_%s", gameId), SubscribeInfo)
	mq.Subscribe(fmt.Sprintf("/player/setPlayerSettings_%s", gameId), SubscribePlayerSettings)
}

// 处理玩家设置订阅
func SubscribePlayerSettings(settings *PersonRtpSettings) {
	fmt.Println("has changed player id:", settings.PlayerId)
	fmt.Println("has changed player:", settings)
	if settings.PlayerId == 0 {
		return
	}
	filter := bson.D{{"_id", settings.PlayerId}}
	var plr DocPlayer
	cur := db.Collection("players")
	err := cur.FindOne(context.TODO(), filter).Decode(&plr)
	if err == mongo.ErrNoDocuments {
		return
	}
	_, err = cur.UpdateOne(context.Background(), bson.M{"_id": plr.Pid}, bson.M{"$set": bson.M{"reward_percent": settings.RewardPercent, "no_award_percent": settings.NoAwardPercent, "target_rtp": settings.TargetRTP, "relieve_rtp": settings.RelieveRTP, "BuyMinAwardPercent": settings.BuyMinAwardPercent}})
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

// 处理商户信息订阅
func SubscribeInfo(store *AppStore) {
	err := incrementReceiveCount()
	if err != nil {
		slog.Error("Increment receiveCount fail:", err)
		return
	}
	receiveCount, err := loadReceiveCount()
	if err != nil {
		slog.Error("Load receiveCount fail:", err)
		return
	}
	fmt.Println("receiveCount:", receiveCount)
	key := MakeAppKey(store.AppId)
	fmt.Println("has changed key:", key)
	fmt.Printf("%v", store)

	// 将 AppStore 存储到 Redis
	err = storeAppIdCache(key, store)
	if err != nil {
		slog.Error("Store AppStore to Redis fail:", err)
		return
	}
}

// 加载 Redis 中的 appIdCache
func LoadAppIdCache(appId string) (*AppStore, error) {
	StoreKey := MakeAppKey(appId)
	// 检查键是否存在
	exists, err := redisClient.Exists(context.Background(), StoreKey).Result()
	if err != nil {
		return nil, err
	}
	var store *AppStore
	//不存在就拉取
	if exists == 0 {
		store, err = GetAppInfo(appId, serviceName)
		if err != nil {
			return nil, err
		}
		//存进去
		err = storeAppIdCache(StoreKey, store)
		if err != nil {
			return nil, err
		}
		return store, nil
	}
	jsonStore, err := redisClient.Get(context.Background(), StoreKey).Result()
	if err == redis.Nil {
		return nil, errors.New("appIdCache not found")
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(jsonStore), &store)
	if err != nil {
		return nil, err
	}
	return store, nil
}

// 存储 Redis 中的 appIdCache
func storeAppIdCache(key string, v *AppStore) error {
	jsonStore, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return redisClient.SetEX(context.Background(), key, jsonStore, time.Second*60*60*24).Err()
}

func SetPlayerCs(PID int64, Cs []float64) {
	key := MakePlayerKey(PID)
	jsonData, err := json.Marshal(Cs)
	if err != nil {
		slog.Error("Failed to marshal float64 array to JSON: %v", err)
	}

	// 将 JSON 字符串存储到 Redis 中 ttl 60 day
	err = redisClient.Set(context.Background(), key, jsonData, time.Second*60*60*24*30*2).Err()
	if err != nil {
		slog.Error("Failed to set float64 array in Redis: %v", err)
	}
}

// 获取玩家 Cs
func GetPlayerCs(AppID string, PID int64, IsEnd bool) ([]float64, error) {
	appStore, err := LoadAppIdCache(AppID)
	if err != nil {
		return nil, err
	}

	playerKey := MakePlayerKey(PID)

	exists, err := redisClient.Exists(context.Background(), playerKey).Result()
	if err != nil {
		slog.Error("Error checking playerKey existence: %v", err)
		return nil, err
	}
	//不存在就设置为商户的
	if exists == 0 {
		SetPlayerCs(PID, appStore.Cs)
		return appStore.Cs, nil
	}
	cs, err := redisClient.Get(context.Background(), playerKey).Result()
	if err == redis.Nil {
		return nil, errors.New("not found")
	}
	var csSlice []float64
	err = json.Unmarshal([]byte(cs), &csSlice)
	if err != nil {
		return nil, err
	}

	//最后一局了，不同的话就让他报错 然后覆盖
	if IsEnd && !reflect.DeepEqual(appStore.Cs, csSlice) {
		SetPlayerCs(PID, appStore.Cs)
		return nil, define.PGNotEnoughCashErr
	}
	return csSlice, nil
}

// 玩家下次获取的区间
type NextMulti struct {
	MinMulti float64 `json:"MinMult"`
	MaxMulti float64 `json:"MaxMult"`
}

// 获取玩家 Cs
func GetPlayerNextData(PID int64) *NextMulti {

	pndKey := MakePlayerNextDataKey(PID)

	exists, err := redisClient.Exists(context.Background(), pndKey).Result()
	if err != nil {
		return nil
	}
	//不存在
	if exists == 0 {
		return nil
	}
	result, err := redisClient.Get(context.Background(), pndKey).Result()
	if err == redis.Nil {
		return nil
	}
	var nextMulti NextMulti
	err = json.Unmarshal([]byte(result), &nextMulti)
	if err != nil {
		return nil
	}
	defer redisClient.Del(context.Background(), pndKey)
	return &nextMulti
}

func GetAppInfo(appId, gameId string) (appStore *AppStore, err error) {

	err = mq.Invoke("/gamecenter/AppStore/getAppInfo", map[string]any{
		"app_id":  appId,
		"game_id": gameId,
	}, &appStore)
	if err != nil {
		slog.Error("GetAppInfo::error", "err", err.Error())
		return
	}

	return
}

func LoadNoAwardPercent(appId string) int {
	app, err := LoadAppIdCache(appId)
	if err != nil {
		slog.Error(err.Error())
		return 0
	}
	// slog.Error("loadNoAwardPercent", app.NoAwardPercent)
	return app.NoAwardPercent
}

func LoadAwardPercent(appId string) int {
	app, err := LoadAppIdCache(appId)
	if err != nil {
		slog.Error(err.Error())
		return 0
	}
	// slog.Error("loadNoAwardPercent", app.NoAwardPercent)
	return app.RewardPercent
}
