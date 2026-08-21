package room

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"log/slog"
	"math"
	"math/rand/v2"
	"os"
	"runtime/debug"
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/lazy"
	"serve/comm/redisx"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"serve/serviceSpribeGame/aviator/gamedata"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

const (
	TIme     = 5
	Duration = TIme * 1000
)

const (
	ROOMSTATE_WAIT_BET = iota
	ROOMSTATE_FLY
	ROOMSTATE_SETTLE
)
const (
	PLAYERSTATE_OB = iota
	PLAYERSTATE_BET
	PLAYERSTATE_CASHOUT
)

var DefaultXArr = []float64{1.00}

type Client struct {
	Conn        *websocket.Conn `json:"-"`
	PlayerState int
	PlayerId    string
	Uid         string
	PlayerName  string
	PlayerIcon  string
	AppId       string
	Multiplier  float64         `json:"multiplier,omitempty"`
	BetMap      map[int]float64 `json:"betMap,omitempty"` //投注窗口列表 0是betId == 1 1是betId == 2
	Seed        string          // 当前轮次的种子投注传入
}

type Room struct {
	Name    string
	Ch      chan *comm.MSG `json:"-"`
	Players map[string]*Client
	State   int
	XArr    []float64
	XIndex  int
	RoundId int64
	Timer   int
	Bets    map[string]*Bets // 当前轮次游戏的投注池子 - 记录当前所有的投注信息
	OldBets []*Bets          // 当前轮次游戏的投注池子 - 记录当前所有的投注信息
	// 上一轮的
	Currency          string
	RoundStartDate    int64
	RoundEndDate      int64
	OnlinePlayerTimer int
	RoundHistories    []*RoundsInfo
	CashRWMU          sync.RWMutex
	CashOutArr        []*CashOut
	Robots            []*Robot
	RobotsBets        map[string]*Bets
	BetRobots         map[string]*Robot
	RobotBetIndex     int //投注阶段加载本次投注机器人的节点
	AdminRobots       *Robot
	//BetList           []*Bets
	StateCh    chan int `json:"-"`
	ServerSeed string   // 当前轮次的服务器种子
	MinBet     float64
}

type RedisRoom struct {
	Name    string
	Players map[string]*Client
	State   int
	XArr    []float64
	XIndex  int
	RoundId int64
	Timer   int
	Bets    map[string]*Bets // 当前轮次游戏的投注池子 - 记录当前所有的投注信息
	OldBets []*Bets          // 当前轮次游戏的投注池子 - 记录当前所有的投注信息
	// 上一轮的
	Currency          string
	RoundStartDate    int64
	RoundEndDate      int64
	OnlinePlayerTimer int
	RoundHistories    []*RoundsInfo
	CashOutArr        []*CashOut
}

type CashOut struct {
	BetId       int
	Bet         float64 //投注
	PlayerId    string  //用户id
	Multiplying float64 //倍率
}

type BetsInfo struct {
	RoomName string
	Bets     map[string]map[int][]*Bets
}

func NewBetsInfo(name string) *BetsInfo {
	return &BetsInfo{
		RoomName: name,
		Bets:     make(map[string]map[int][]*Bets),
	}
}

func NewRoom(name string) *Room {
	var room *Room

	roomKey := fmt.Sprintf("backup:%s", name)
	roomInfo, err := GetRoomInfoFromRedis(roomKey)
	if roomInfo != nil && err == nil {
		room = roomInfo
		room.Ch = make(chan *comm.MSG, 1024)
	} else {
		room = &Room{
			Name:       name,
			Ch:         make(chan *comm.MSG, 1024),
			Players:    make(map[string]*Client),
			Bets:       make(map[string]*Bets),
			RobotsBets: make(map[string]*Bets),
			StateCh:    make(chan int, 1),
		}
		flag := gamedata.Settings.RobotFlag
		flag = true
		if flag {
			for i := 0; i < 300; i++ {
				sf := ut.NewSnowflake()
				id := strconv.Itoa(int(sf.NextID()))
				robot := &Robot{
					Id:         id,
					PlayerIcon: fmt.Sprintf(`av-%v.png`, rand.IntN(72)+1),
					PlayerName: ut.GenerateRandomString2(5),
					RoomId:     name,
				}
				room.Robots = append(room.Robots, robot)
			}
		}
	}
	if err != nil {
		slog.Error("NewRoom::GetRoomInfoFromRedis Err", "err", err)
	}

	return room
}

func (r *Room) AddMsg(m *comm.MSG) {
	r.Ch <- m
}

func (r *Room) Run() {
	// 自定义日志格式
	logger := log.New(os.Stdout, "", 0) // 禁用默认的时间戳和前缀
	// 自定义时间格式
	logger.SetFlags(log.Lmsgprefix)                                      // 只保留消息前缀
	logger.SetPrefix(time.Now().Format("2006-01-02 15:04:05.000") + " ") // 精确到毫秒

	//房间主逻辑
	go func() {
		defer func() {
			slog.Info(r.Name + "该房间消息已损坏")
		}()
		for {
			select {
			case m := <-r.Ch: //房间的管道收到每秒的信号就开始消费
				//logger.Println(m)
				r.Msg(m) //消费的逻辑
			}
		}
	}()
}
func (r *Room) Msg(m *comm.MSG) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "AppID ", r.Name, `Msg处理消息出错`)
			os.Stdout.Write(debug.Stack())
		}
	}()

	switch m.Typ {
	case comm.NETMSG:
	case comm.TIMER_1s: //每秒房间要干的事
		r.MissionMaker()
	case comm.TIMER_100: //每秒房间要干的事
		r.MissionMaker100()
	case comm.NETCLOSE:
		//断线判断需不需要还钱
		session, err := comm.GetSession(m.Conn)
		if err != nil {
			slog.Error("RemovePlayer Err", "AppID ", r.Name, "GetSession err", err)
			return
		}
		plr := session.Plr
		if r.State == ROOMSTATE_WAIT_BET {
			for i := 0; i < 2; i++ {
				historyBetKey := fmt.Sprintf("%s_%v", session.PlayerId, i+1)
				if _, ok := r.Bets[historyBetKey]; ok {
					//huanqian
					costGold := ut.Money2Gold(r.Bets[historyBetKey].Bet)
					//调用游戏中心去修改玩家钱包信息，这里是还原下注金额
					balance, err := slotsmongo.ModifyGold(&slotsmongo.ModifyGoldPs{
						Pid:     plr.PID,
						Change:  costGold,
						RoundID: strconv.Itoa(int(r.RoundId)),
						Reason:  slotsmongo.ReasonRefund,
					})
					if err != nil {
						slog.Error("newUpgrader Err", "AppID ", r.Name, "slotsmongo.ModifyGold err", err)
						return
					}
					sf := ut.NewSnowflake()
					id := sf.NextID()
					betId := strconv.Itoa(int(int32(id % math.MaxInt32)))
					nowTime := time.Now().Unix()
					betLog := &slotsmongo.AddBetLogParamsAviator{
						ID:               strconv.FormatInt(id, 10),
						Pid:              plr.PID,
						AppID:            plr.AppID,
						Uid:              plr.Uid,
						UserName:         strconv.FormatInt(plr.PID, 10),
						CurrencyKey:      plr.CurrencyKey,
						RoomId:           plr.AppID,
						Bet:              ut.Money2Gold(r.Bets[historyBetKey].Bet),
						Balance:          balance,
						CashOutDate:      nowTime,
						RoundBetId:       betId,
						RoundId:          r.RoundId,
						GameId:           comm.GameID,
						LogType:          comm.LogType_ReMove,
						FinishType:       comm.FinishType_0,
						ManufacturerName: "spribe",
						BetId:            i + 1,
						GameType:         define.GameType_Mini,
						Completed:        false,
					}
					slog.Info("ClientCancelBet时,即将插入取消投注记录", "入参: belog", betLog)
					slotsmongo.AddBetLogAviator(betLog) //待改为update
				}
			}
		}
		//移除玩家
		RemovePlayer(m.Conn, r)

	default:
		r.ActionHandler(m.Typ, m.Conn, m.Msg)

	}
}

/*
	case "AviatorGameInfoIdReq":
		c.WriteMessage(messageType, []byte(`{"id":1,"targetController":0,"content":{"pi":0,"rl":[[1,"game_state","default",false,false,0,20,[]]],"rs":0,"un":"67590\u0026\u0026release","zn":"aviator_core_inst4"}}`))
*/

func (r *Room) NotifyAll(h func(c *websocket.Conn, r *Room)) {
	if h != nil {
		for _, player := range r.Players {
			h(player.Conn, r)
		}
	}
}

func (r *Room) NotifyAll2(data []byte) {
	for _, player := range r.Players {
		err := player.Conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			slog.Error("player.Conn.WriteMessage", "Err", err.Error())
		}
	}
}

func (r *Room) MissionMaker() {
	switch r.State {
	case ROOMSTATE_WAIT_BET: //等待阶段 其实是下注阶段，(虽然别的阶段也能下注，但都是前端在下一回合的这个阶段发请求)。这有这里是给当次回合下注的。这个阶段持续?秒
		//等待阶段分两种，等待开始、等待结束 开始时"newStateId": 1, 结束时"newStateId": 2,
		//结束后秒接第一次的飞机节点
		if 0 < r.Timer && r.Timer < TIme {
			//根据redis的排行榜信息返回
			r.Timer++
			r.NotifyAll(UpdateCurrentBets)
			return
		}
		r.NotifyAll(ChangeState)
		if r.Timer == TIme {
			//遍历所有投注的人，全部同步到clickhouse中
			r.State = ROOMSTATE_FLY
			r.StateCh <- ROOMSTATE_FLY
			// 每次转阶段之后向redis同步房间状态
			err := RoomBackup2Redis(r)
			if err != nil {
				slog.Error("MissionMaker::RoomBackup2Redis Err", "AppID ", r.Name, "err", err.Error())
			}
			r.Timer = 0
			//提前配置房间内下次使用的节点数组
			info, err := redisx.LoadAppIdCache(r.Name)
			if err != nil {
				slog.Error("MissionMaker::LoadAppIdCache Err", "AppID ", r.Name, "err", err.Error())
			}
			if FloatProbabilityCheck(info.CrashRate) {
				r.XArr = DefaultXArr
			} else {
				if info.ProfitMargin != 0 {
					r.XArr = comm.CalcPoints(math.Round(100-info.ProfitMargin) / 100)
				} else {
					r.XArr = DefaultXArr
					break
				}
				bet, err := GetBet(r.Name)
				if err != nil {
					slog.Error("MissionMaker::GetBet Err", "AppID ", r.Name, "err", err.Error())
				}
				if bet != nil {
					winPoints := info.MaxWinPoints / bet.Bet
					if winPoints < 1.00 {
						winPoints = 1.00
					}
					if r.XArr[len(r.XArr)-1] > info.MaxMultiple || r.XArr[len(r.XArr)-1] > info.MaxWinPoints {
						r.XArr = DefaultXArr
					}
				}
			}
			//r.XArr = comm.CalcPointsDev(0.97)
		} else {
			coll := db.Collection2(lazy.ServiceName, "RoundIdCount")
			update := db.D("$inc", db.D("RoundCount", 1))
			opts := options.FindOneAndUpdate().SetUpsert(true).SetProjection(db.D("RoundCount", 1)).SetReturnDocument(options.After)
			var doc RoundCount
			err := coll.FindOneAndUpdate(context.TODO(), db.ID(r.Name), update, opts).Decode(&doc)
			if err != nil {
				slog.Error("MissionMaker::FindOneAndUpdate Err", "AppID ", r.Name, "err", err.Error())
			}
			r.RoundId = doc.RoundCount
			r.ServerSeed = slotsmongo.GenerateRandStr(40)
			r.Timer++
			r.NotifyAll(UpdateCurrentBetsInit)
			// 删除redis房间数据
			DeleteKeyAsync(redisx.GetClient(), `backup:`+r.Name)
			DeleteKeyAsync(redisx.GetClient(), r.Name)
			r.RoundStartDate = time.Now().UnixMilli()
			//每次转阶段之后向redis同步房间状态
			err = RoomBackup2Redis(r)
			if err != nil {
				slog.Error("MissionMaker::RoomBackup2Redis Err", "AppID ", r.Name, "err", err.Error())
			}
			//准备当前轮次投注的机器人
			{
				go func() {
					defer func() {
						if err := recover(); err != nil {
							log.Println(err, "AppID ", r.Name, ``+`准备当前轮次投注的机器人`)
							os.Stdout.Write(debug.Stack())
						}
					}()
					info, err := redisx.LoadAppIdCache(r.Name)
					if err != nil {
						slog.Error("ClientGameInfo Err", "AppID ", r.Name, "LoadAppIdCache err", err)
					}
					if info.OnlineUpNum < 1 {
						slog.Error("ClientGameInfo Err", "AppID ", r.Name, "OnlineUpNum", info.OnlineUpNum)
					}
					if info.OnlineDownNum < 1 {
						slog.Error("ClientGameInfo Err", "AppID ", r.Name, "OnlineDownNum", info.OnlineUpNum)
					}
					if info.OnlineUpNum > info.OnlineDownNum {
						slog.Error("ClientGameInfo Err", "AppID ", r.Name, "err", "后台配置中房间投注下限大于了投注上限")
					}
					r.RobotsBets = make(map[string]*Bets)
					r.BetRobots = make(map[string]*Robot)
					for i := 0; i < len(r.Robots); i++ {
						randomNumber := 35 + rand.IntN(40-35+1)
						if rand.IntN(100) < randomNumber {
							multiplier := r.Robots[i].GetMultiplier(1)
							if multiplier == 0.0 {
								continue
							}
							historyBetKey := fmt.Sprintf("%s_%v", r.Robots[i].Id, 1)
							bet := &Bets{
								Bet:          float64(GetRandInt(int(info.OnlineUpNum)*10, int(info.OnlineDownNum)*10)) / 10,
								BetID:        1,
								Currency:     r.Currency,
								PlayerID:     r.Robots[i].Id,
								ProfileImage: r.Robots[i].PlayerIcon,
								Username:     r.Robots[i].PlayerName,
								RoundId:      r.RoundId,
								Multiplier:   multiplier,
							}
							r.RobotsBets[historyBetKey] = bet
							if rand.IntN(100) < 40 {
								multiplier = r.Robots[i].GetMultiplier(2)
								if multiplier == 0.0 {
									continue
								}
								historyBetKey = fmt.Sprintf("%s_%v", r.Robots[i].Id, 2)
								bet2 := &Bets{
									Bet:          float64(GetRandInt(int(info.OnlineUpNum)*10, int(info.OnlineDownNum)*10)) / 10,
									BetID:        2,
									Currency:     r.Currency,
									PlayerID:     r.Robots[i].Id,
									ProfileImage: r.Robots[i].PlayerIcon,
									Username:     r.Robots[i].PlayerName,
									RoundId:      r.RoundId,
									Multiplier:   multiplier,
								}
								r.RobotsBets[historyBetKey] = bet2
							}
						}
					}
					slog.Info(fmt.Sprintf("房间%v要投注的机器人准备完毕", r.Name))

					slog.Info(fmt.Sprintf("房间%v机器人准备投注", r.Name))
					for key, bets := range r.RobotsBets {
						select {
						case c := <-r.StateCh:
							if c == ROOMSTATE_FLY {
								break
							}
						default: // 如果没有信号，执行默认逻辑
							// 机器人投注需要生成seed
							RobotBet2(bets, r)
							seed := slotsmongo.GenerateRandStr(20)
							key := strings.Split(key, "_")
							robots := &Robot{
								BetInfo:    []*Bets{bets},
								BetNum:     1,
								Id:         key[0],
								PlayerIcon: bets.ProfileImage,
								PlayerName: bets.Username,
								RoomId:     r.Name,
								Seed:       seed,
							}
							r.CashRWMU.Lock()
							r.BetRobots[key[0]] = robots
							r.CashRWMU.Unlock()
							//slog.Info(fmt.Sprintf("房间%v插入了一个机器人", r.Name))
						}
					}
					slog.Info(fmt.Sprintf("机器人投注完成 房间：%v 一共%v个", r.Name, len(r.RobotsBets)))
				}()
			}

		}
	case ROOMSTATE_SETTLE: //结算阶段 就是扣钱阶段，这个阶段中会有一批玩家还属于PLAYERSTATE_BET状态，遍历这些玩家扣钱，并且重置他们的状态。然后想办法拖延5秒，再变换房间状态
		if r.Timer == 0 {
			// 每次转阶段之后向redis同步房间状态
			err := RoomBackup2Redis(r)
			if err != nil {
				slog.Error("MissionMaker::RoomBackup2Redis Err", "AppID ", r.Name, "err", err.Error())
			}
			if len(r.RoundHistories) < 60 {
				r.RoundHistories = append(r.RoundHistories, &RoundsInfo{
					MaxMultiplier: r.XArr[len(r.XArr)-1],
					RoundID:       r.RoundId,
				})
			} else {
				temp := make([]*RoundsInfo, len(r.RoundHistories[1:]))
				copy(temp, r.RoundHistories[1:])
				temp = append(temp, &RoundsInfo{
					MaxMultiplier: r.XArr[len(r.XArr)-1],
					RoundID:       r.RoundId,
				})
				r.RoundHistories = temp
			}
			r.CashOutFuncEliminate()
			//这是飞机爆炸后等待阶段的0秒
			//给r.Bets下的所有人update投注记录的状态
			{
				go func() {
					defer func() {
						if err := recover(); err != nil {
							log.Println(err, "AppID ", r.Name, ``+`给r.Bets下的所有人update投注记录的状态`)
							os.Stdout.Write(debug.Stack())
						}
					}()
					for _, plrInfo := range r.Bets {
						if plrInfo.HistoryFlag == true {
							// todo update这些人的RoundMaxMultiplier字段 - 如果未来飞机会在xArr数组len-1之前爆炸
							//slotsmongo.UpdateBetLogAviator()
						} else {
							sf := ut.NewSnowflake()
							id := sf.NextID()
							betId := strconv.Itoa(int(int32(id % math.MaxInt32)))
							plr := r.Players[plrInfo.PlayerID]
							playerId, err := strconv.ParseInt(plr.PlayerId, 10, 64)
							if err != nil {
								slog.Error("ClientLogin Err", "strconv.ParseInt err", err)
							}
							nowTime := time.Now().UnixMilli()
							balance, err := slotsmongo.GetBalance(playerId)
							if err != nil {
								slog.Error("MissionMaker::slotsmongo.GetBalance Err", "AppID ", r.Name, "GetBalance err", err)
								continue
							}
							betLog := &slotsmongo.AddBetLogParamsAviator{
								ID:                 strconv.FormatInt(id, 10),
								Pid:                playerId,
								AppID:              plr.AppId,
								Uid:                plr.Uid,
								UserName:           strconv.FormatInt(playerId, 10),
								CurrencyKey:        plrInfo.Currency,
								RoomId:             plr.AppId,
								Bet:                ut.Money2Gold(plrInfo.Bet),
								Win:                0,
								Balance:            balance,
								CashOutDate:        nowTime,
								Payout:             r.XArr[len(r.XArr)-1],
								Profit:             ut.Money2Gold(-plrInfo.Bet),
								RoundBetId:         betId,
								RoundId:            r.RoundId,
								GameId:             comm.GameID,
								LogType:            comm.LogType_Bet,
								FinishType:         comm.FinishType_1,
								RoundMaxMultiplier: r.XArr[len(r.XArr)-1],
								ManufacturerName:   "spribe",
								BetId:              plrInfo.BetID,
								GameType:           define.GameType_Mini,
								Completed:          true,
							}
							slotsmongo.AddBetLogAviator(betLog)
						}
					}
				}()
			}
			// 构建验证结构体插入到mongo中
			{
				go func() {
					defer func() {
						if err := recover(); err != nil {
							log.Println(err, "AppID ", r.Name, ``+`构建验证结构体插入到mongo中`)
							os.Stdout.Write(debug.Stack())
						}
					}()
					fairnessPlay := []PlayerSeed{}
					fairnessData := &FairnessDataBson{
						RoomID:   r.Name,
						Number:   126.83721838688118,
						CreateAt: time.Now().UTC(),
					}
					fairnessData.ServerSeed = r.ServerSeed
					fairnessData.RoundID = r.RoundId
					fairnessData.RoundStartDate = r.RoundStartDate
					fairnessData.Result = r.XArr[len(r.XArr)-1]
					longServerSeed := r.ServerSeed
					for _, player := range r.Players {
						if player.Seed != "" {
							fairnessPlay = append(fairnessPlay, PlayerSeed{
								ProfileImage: player.PlayerIcon,
								Seed:         player.Seed,
								Username:     player.PlayerName,
							})
						}
						longServerSeed += player.Seed
						if len(fairnessPlay) == 3 {
							break
						}
					}
					if len(fairnessPlay) < 3 {
						r.CashRWMU.RLock()
						for _, robot := range r.BetRobots {
							fairnessPlay = append(fairnessPlay, PlayerSeed{
								ProfileImage: robot.PlayerIcon,
								Seed:         robot.Seed,
								Username:     robot.PlayerName,
							})
							longServerSeed += robot.Seed
							if len(fairnessPlay) == 3 {
								break
							}
						}
						r.CashRWMU.RUnlock()
					}
					fairnessData.PlayerSeeds = fairnessPlay
					fairnessData.SeedSHA256, fairnessData.PartSeedHexNumber, fairnessData.PartSeedDecimalNumber, _ = slotsmongo.GenerateSha512(longServerSeed)
					coll := db.Collection2(lazy.ServiceName, "fairnessData")
					coll.InsertOne(context.TODO(), fairnessData)
				}()
			}
			// 添加数据到上一手的key里
			{
				go func() {
					defer func() {
						if err := recover(); err != nil {
							log.Println(err, "AppID ", r.Name, " "+`添加数据到上一手的key里`)
							os.Stdout.Write(debug.Stack())
						}
					}()
					slog.Info("即将插入上一手数据", "房间", r.Name)
					roundInfo := &RoundInfo{
						RoundID:        r.RoundId,
						RoundStartDate: r.RoundStartDate,
						RoundEndDate:   r.RoundEndDate,
						Multiplier:     r.XArr[len(r.XArr)-1],
					}
					err := AddPreviousRoundInfo(r.Name, roundInfo)
					if err != nil {
						slog.Error("MissionMaker::AddPreviousRoundInfo", "AppID ", r.Name, "err", err.Error())
					}
					err = AddPreviousBetInfo(redisx.GetClient(), r.Name)
					if err != nil {
						slog.Error("添加上一手失败MissionMaker::AddPreviousBetInfo Err", "Err", err.Error(), "房间", r.Name)
					}
					slog.Info("插入上一手数据处理完成", "房间", r.Name)
				}()
			}
		}
		if 0 <= r.Timer && r.Timer < 5 {
			r.Timer++
			return
		}
		if r.Timer == 5 {
			r.State = ROOMSTATE_WAIT_BET
			r.Timer = 0
			r.RoundEndDate = time.Now().UnixMilli()
			// 每次转阶段之后向redis同步房间状态
			err := RoomBackup2Redis(r)
			if err != nil {
				slog.Error("MissionMaker::RoomBackup2Redis Err", "AppID ", r.Name, "err", err.Error())
			}
			r.Bets = make(map[string]*Bets)
			roundInfo := TopRound{
				EndDate:        r.RoundEndDate,
				MaxMultiplier:  r.XArr[len(r.XArr)-1],
				RoundId:        r.RoundId,
				RoundStartDate: r.RoundStartDate / 1000,
				ServerSeed:     `cbdcf0fddd8f93e5feb11fcfa7abe485c980189f14d5c93bfb3938e13093031a`,
				Zone:           "aviator_core",
			}
			for _, date := range DateList {
				keyTopRounds := fmt.Sprintf("%s_%s_%s", r.Name, "topround", date)
				err := StoreTopRounds(redisx.GetClient(), keyTopRounds, roundInfo)
				if err != nil {
					slog.Error("StoreTopRounds", "AppID ", r.Name, "err", err.Error())
				}
			}
		}
	}
	if r.OnlinePlayerTimer == 24 {
		r.NotifyAll(OnlinePlayers)
		r.OnlinePlayerTimer = 0
	} else {
		r.OnlinePlayerTimer++
	}
}

func (r *Room) MissionMaker100() {
	switch r.State {
	case ROOMSTATE_FLY: //起飞阶段 每秒返回飞机价格节点直到结束(准备好的节点数组到头) 进入结算阶段
		if r.CashOutArr != nil {
			r.CashOutFuncDetection()
		}
		r.NotifyAll(SendCoordinate)
		r.NotifyAll(UpdateCurrentCashOuts)
		//处理这一节点需要提现的机器人
		{
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println(err, "AppID ", r.Name, " "+`处理这一节点需要提现的机器人`)
						os.Stdout.Write(debug.Stack())
					}
				}()
				r.CashRWMU.Lock()
				for k, val := range r.RobotsBets {
					if val.Multiplier <= r.XArr[r.XIndex] && val.Effective {
						RobotCashOut2(val, r)
						delete(r.RobotsBets, k)
					}
				}
				r.CashRWMU.Unlock()
			}()
		}
		//判断是否是最后一次飞行，如果是最后一次飞行，还要发送两条消息，并且重置该房间的飞机节点数组和数组下标以及变换房间状态
		{
			if r.XIndex == len(r.XArr)-1 {
				r.State = ROOMSTATE_SETTLE
				r.NotifyAll(RoundChartInfo)
				r.NotifyAll(SendCrashCoordinate)
				r.NotifyAll(NewBalance)
				r.XIndex = 0
				for _, client := range r.Players {
					client.PlayerState = PLAYERSTATE_OB
					client.BetMap = make(map[int]float64)
					client.Multiplier = 0
				}
			} else { //飞机节点数组下标+1
				r.XIndex++
			}
		}
	}
}

func (r *Room) CashOutFuncAdd(BetId int, Bet, AutoCashOut float64, PlayerId string) {
	r.CashRWMU.Lock()
	defer r.CashRWMU.Unlock()
	if AutoCashOut != float64(0) {
		r.CashOutArr = append(r.CashOutArr, &CashOut{
			BetId:       BetId,
			Bet:         Bet,
			PlayerId:    PlayerId,
			Multiplying: AutoCashOut,
		})
	}
}

func (r *Room) CashOutFuncCancel(PlayerId string) {
	r.CashRWMU.Lock()
	defer r.CashRWMU.Unlock()
	var tmp []*CashOut
	for _, v := range r.CashOutArr {
		if PlayerId != v.PlayerId {
			tmp = append(tmp, v)
		}
	}
	if len(tmp) > 0 {
		r.CashOutArr = tmp
	}
}

func (r *Room) CashOutFuncDetection() {
	Multiplying := r.XArr[r.XIndex]
	r.CashRWMU.Lock()
	defer r.CashRWMU.Unlock()
	var tmp []*CashOut
	for _, v := range r.CashOutArr {
		if v.Multiplying <= Multiplying {
			val := *v
			go func(value CashOut) {
				defer func() {
					if err := recover(); err != nil {
						log.Println(err, "AppID ", r.Name, ` `+`CashOutFuncDetection`)
						os.Stdout.Write(debug.Stack())
					}
				}()
				//reqs, _ := json.Marshal(CommBody{
				//	Channel:          "AviatorCashOutIdReq",
				//	TargetController: 1,
				//	Content: &Content{
				//		C: "cashOut",
				//		R: -1,
				//		P: P{
				//			BetID:       value.BetId,
				//			Multiplying: v.Multiplying,
				//		},
				//	},
				//})
				pp := ut.NewSFSObject()
				p := ut.NewSFSObject()
				p.PutDouble("betId", float64(value.BetId))
				p.PutDouble("multiplying", v.Multiplying)
				pp.PutSFSObject("p", p)
				ClientCashOut(r.Players[v.PlayerId].Conn, websocket.BinaryMessage, pp, r)
			}(val)
		} else {
			tmp = append(tmp, v)
		}
	}
	r.CashOutArr = tmp
}

func (r *Room) CashOutFuncEliminate() {
	r.CashRWMU.Lock()
	defer r.CashRWMU.Unlock()
	r.CashOutArr = nil
}

func (r *Room) ActionHandler(action uint8, c *websocket.Conn, p *ut.SFSObject) {
	messageType := websocket.BinaryMessage

	switch action {
	case comm.AviatorLoginIdReq:
		ClientLogin(c, messageType, p, r)
	case comm.AviatorGameInfoIdReq:
		ClientGameInfo(c, messageType, p, r)
	case comm.AviatorCurrentBetsInfoIdReq:
		ClientCurrentBetsInfo(c, messageType, p, r)
	case comm.AviatorBetIdReq:
		ClientBet(c, messageType, p, r)
	case comm.AviatorCancelBetIdReq:
		ClientCancelBet(c, messageType, p, r)
	case comm.AviatorCashOutIdReq:
		ClientCashOut(c, messageType, p, r)
	case comm.AviatorBetHistoryIdReq:
		ClientBetHistory(c, messageType, p, r)
	case comm.AviatorGameStatePingIdReq:
		so := ut.NewSFSObject()
		p := ut.NewSFSObject()
		pp := ut.NewSFSObject()
		p.PutSFSObject("p", pp)
		p.PutString("c", "PING_RESPONSE")
		so.AddCreatePAC(p, 1, 13)
		marshal, _ := so.ToBinary()
		//c.WriteMessage(messageType, []byte(`{"id":13,"targetController":1,"content":{"c":"PING_RESPONSE"}}`))
		c.WriteMessage(messageType, marshal)
	case comm.AviatorAddChatMessageIdReq:
		ClientAddChatMessage(c, messageType, p, r)
	case comm.AviatorLikeIdReq:
		ClientLikeChat(c, messageType, p, r)
	case comm.ClientSearchGifs:
		ClientSearchGifs(c, messageType, p, r)
	case comm.AviatorPreviousRoundInfoIdReq:
		ClientPreviousRoundInfo(c, messageType, p, r)
	case comm.AviatorGetHugeWinsInfoIdReq:
		ClientGetHugeWinsInfo(c, messageType, p, r)
	case comm.AviatorGetTopWinsInfoIdReq:
		ClientGetTopWinsInfo(c, messageType, p, r)
	case comm.AviatorGetTopRoundsInfoIdReq:
		ClientGetTopRoundsInfo(c, messageType, p, r)
	case comm.AviatorChangeProfileImageIdReq:
		ClientChangeProfileImage(c, messageType, p, r)
	case comm.AviatorRoundFairnessIdReq:
		ClientRoundFairnessData(c, messageType, p, r)
	case comm.ServerSeedHandler:
		ClientServerSeedHandler(c, messageType, p, r)

	default:
		slog.Error("这消息没地方走呢？")
	}
}
