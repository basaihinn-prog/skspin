package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/mq"
	"game/comm/mux"
	"game/duck/mongodb"
	"game/service/admin/channel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/CreatePlayerRTP", "新建玩家RTP设置", "AdminInfo", CreatePlayerinfo, CreatePlayerRTP{})
	mux.RegHttpWithSample("/AdminInfo/SetDemoPlayerRTP", "试玩站设置超高RTP", "AdminInfo", SetDemoPlayerRTP, CreatePlayerRTP{})
	RegMsgProc("/AdminInfo/PlayerRTPControlList", "玩家RTP控制列表", "AdminInfo", PlayerRTPControlList, getrtplistreq{})
	//主动取消控制
	RegMsgProc("/AdminInfo/CancelPlayerRTPControl", "主动取消玩家RTP控制", "AdminInfo", CancelPlayerRTP, CancelRTP{})

	RegMsgProc("/AdminInfo/CancelPlayerRTPControlByPId", "主动取消玩家RTP控制", "AdminInfo", CancelPlayerRTPByPId, CancelRTPByPId{})

	// 内部调用
	mux.RegRpc("/AdminInfo/Interior/CreatePlayerRTP", "内部调用MQ-新建玩家RTP控制", "AdminInfo", createPlayerRTP, CreatePlayerRTP{})

}

var (
	minScore = int64(1)
	maxScore = int64(1000000)
	minMult  = 30
	maxMult  = 10000
)

type CancelRTP struct {
	IDList string `json:"Ids"`
}
type CancelRTPByPId struct {
	Pid int64 `json:"Pid"`
}

type CreatePlayerRTP struct {
	DemoUserId        string // 试玩站用户id
	GameID            string // 游戏列表
	Pid               int64  // 玩家唯一标识
	AppID             string // 所属商户
	AppIDandPlayerID  string
	ContrllRTP        float64 //控制RTP
	AutoRemoveRTP     float64 //自动解除RTP
	BuyRTP            int     //购买游戏RTP
	PersonWinMaxMult  int     //个人最高盈利倍数
	PersonWinMaxScore int64   //个人最高盈利分值
}

type PlayerRTPControlModel struct {
	ID                 int64  `bson:"_id"`
	AppID              string //运营商  关联运营商
	GameID             string //游戏编号 关联游戏
	CreateAt           time.Time
	Uid                int64   //用户账号
	GameName           string  //游戏名称
	GameRTP            int64   //游戏配置RTP 关联游戏配置
	PlayerRTP          int64   //玩家历史RTP RTP控制之前的
	ControlRTP         float64 //账号RTP
	RewardPercent      int64   //账号RTP
	NoAwardPercent     int64   //账号RTP
	AutoRemoveRTP      float64 //自动解除RTP
	AutoRewardPercent  int64   //自动解除RTP
	AutoNoAwardPercent int64   //自动解除RTP
	Status             int     //状态
}
type rtplistres struct {
	ID                   int64     `bson:"_id" json:"Id"`
	AppID                string    `bson:"AppID" json:"OperatorName"`          //运营商  关联运营商
	GameID               string    `bson:"GameID"`                             //游戏编号 关联游戏
	CreateAt             time.Time `bson:"CreateAt" json:"ControlTime" `       //控制时间
	Uid                  string    `bson:"Uid" json:"Uid"`                     //用户账号
	Pid                  int64     `bson:"Pid" json:"Pid"`                     //用户唯一标识
	GameName             string    `bson:"GameName" json:"GameName"`           //游戏名称
	GameIcon             string    `bson:"-" json:"GameIcon"`                  //游戏图标			// 需求优化  查询列表添加游戏图标  与游戏厂商
	GameManufacturerName string    `bson:"-" json:"GameManufacturerName"`      //游戏厂商			// 需求优化  查询列表添加游戏图标  与游戏厂商
	GameRTP              int64     `bson:"RTP" json:"GameRTP"`                 //游戏配置RTP 关联游戏配置
	PlayerRTP            int64     `bson:"PlayerRTP" json:"PlayerHistoryRTP"`  //玩家历史RTP RTP控制之前的
	ControlRTP           int64     `bson:"ControlRTP" json:"ControlRTP"`       //账号RTP
	PresentRTP           int64     `bson:"PresentRTP" json:"ControllingRTP"`   //当前RTP
	AutoRemoveRTP        int64     `bson:"AutoRemoveRTP" json:"AutoRemoveRTP"` //自动解除RTP
	BuyRTP               int       `bson:"BuyRTP" json:"BuyRTP"`
	BuyType              int       `bson:"BuyType" json:"BuyType"`
}

type playerRTPListRes struct {
	List        []*rtplistres
	RTPCofigOff int64 //是否开启rtp控制
	Count       int64
}
type PlayerRtpControlReq struct {
	OperateId string //运营商id
	Uid       string //用户账号
	CreateAt  time.Time
	GameID    string //游戏id
}

// 玩家RTP控制列表
type PlayerRtpControlRes struct {
	AppID              string //运营商  关联运营商
	GameID             string //游戏编号 关联游戏
	CreateAt           time.Time
	Uid                string //用户账号
	GameName           string //游戏名称
	GameRTP            int64  //游戏配置RTP 关联游戏配置
	PlayerRTP          int64  //玩家历史RTP RTP控制之前的
	ControlRTP         int64  //账号RTP
	RewardPercent      int    //账号RTP
	NoAwardPercent     int    //账号RTP
	AutoRemoveRTP      int64  //自动解除RTP
	AutoRewardPercent  int    //自动解除RTP
	AutoNoAwardPercent int    //自动解除RTP
}

type PersonRtpSettings struct {
	RewardPercent  int   `json:"reward_percent"`
	NoAwardPercent int   `json:"no_award_percent"`
	PlayerId       int64 `json:"player_id"`
}
type PlayerRtpSettings struct {
	RewardPercent      int   `json:"reward_percent"`
	NoAwardPercent     int   `json:"no_award_percent"`
	TargetRTP          int   `json:"target_rtp"`
	RelieveRTP         int   `json:"relieve_rtp"`
	PlayerId           int64 `json:"player_id"`
	BuyMinAwardPercent int   `json:"buy_min_award_percent"`
}

type getrtplistreq struct {
	OperatorId       int    `json:"OperatorId"`       // 运营商ID
	Pid              int64  `json:"UId"`              // 玩家唯一标识
	ControlTimeStart string `json:"ControlTimeStart"` // 控制时间开始
	ControlTimeEnd   string `json:"ControlTimeEnd"`   // 控制时间结束
	GameId           string `json:"GameId"`           // 游戏ID
	Page             int64  `json:"Page"`             // 页码
	PageSize         int64  `json:"PageSize"`         // 每页大小
	Manufacturer     string `json:"Manufacturer"`     //游戏厂商
}

func PlayerRTPControlList(ctx *Context, ps getrtplistreq, ret *playerRTPListRes) (err error) {
	var appidlist []string
	if ps.OperatorId > 0 {
		var operator comm.Operator
		err = CollAdminOperator.FindOne(bson.M{"_id": ps.OperatorId}, &operator)
		if err != nil {
			return err
		}
		appidlist = append(appidlist, operator.AppID)
	} else {
		appidlist, err = GetOperatopAppID(ctx)
		if err != nil {
			return err
		}
	}

	query := bson.M{}

	if ps.GameId != "ALL" && ps.GameId != "" {
		query["GameID"] = ps.GameId
	} else {
		return errors.New("请选择游戏")
	}

	query["FromPlant"] = "admin"
	query["Status"] = 1
	query["AppID"] = bson.M{
		"$in": appidlist,
	}
	//startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if ps.Pid != 0 {
		query["Pid"] = ps.Pid
	}

	rDate := ps.ControlTimeStart + ps.ControlTimeEnd

	switch rDate {
	case "":
	case ps.ControlTimeStart:
		t, _ := time.Parse("2006-01-02 15:04:05", ps.ControlTimeStart)
		mtime := mongodb.NewTimeStamp(t)
		query["CreateAt"] = bson.M{"$gte": mtime}
	case ps.ControlTimeEnd:
		t, _ := time.Parse("2006-01-02 15:04:05", ps.ControlTimeEnd)
		mtime := mongodb.NewTimeStamp(t)
		query["CreateAt"] = bson.M{"$lte": mtime}
	default:
		if len(rDate) > 0 {
			t, _ := time.Parse("2006-01-02 15:04:05", ps.ControlTimeStart)
			mstarttime := mongodb.NewTimeStamp(t)
			t, _ = time.Parse("2006-01-02 15:04:05", ps.ControlTimeEnd)
			mendtime := mongodb.NewTimeStamp(t)

			query["CreateAt"] = bson.M{
				"$gte": mstarttime,
				"$lte": mendtime,
			}
		}
	}
	var find = mongodb.FindPageOpt{
		Page:     ps.Page,
		PageSize: ps.PageSize,
		Sort:     bson.M{"CreateAt": -1},
		Query:    query,
	}

	ret.Count, err = CollPlayerRTPControl.FindPage(find, &ret.List)
	if err != nil {
		return
	}

	var pl struct {
		BetAmount int64 `bson:"betamount"`
		WinAmount int64 `bson:"winamount"`
	}
	var msg struct {
		Pid    int64
		GameId string
		AppId  string
	}
	path := "/gamecenter/player/getPlayerInGame"
	gameNameMap := map[string]string{}
	gameMap := map[string]*comm.Game2{} // hml   用过ID映射game信息
	_, games := GetGameListFormatName(ctx.Lang, nil)

	for _, game := range games {
		gameNameMap[game.ID] = game.Name
		gameMap[game.ID] = game
	}
	appidlist = appidlist[:0]
	gameidlist := []string{}
	for _, r := range ret.List {
		appidlist = append(appidlist, r.AppID)
		gameidlist = append(gameidlist, r.GameID)
	}
	filter := bson.M{
		"AppID": bson.M{
			"$in": appidlist,
		},
		"GameId": bson.M{
			"$in": gameidlist,
		}}
	opconfig := []*comm.GameConfig{}
	findOp := options.Find().SetSort(bson.D{{"AppID", -1}})
	err = CollGameConfig.FindAllOpt(filter, &opconfig, findOp)
	if err != nil {
		return
	}
	mapOpConfig := map[string][]*comm.GameConfig{}
	for _, config := range opconfig {
		mapOpConfig[config.AppID] = append(mapOpConfig[config.AppID], config)
	}

	//商户游戏配置
	coll := NewOtherDB("game").Collection("Players").Coll()
	for i, play := range ret.List {
		msg.AppId = play.AppID
		msg.GameId = play.GameID
		msg.Pid = play.Pid
		err = mq.Invoke(path, msg, &pl)

		var pRTP int64
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				pRTP = 0
			} else {
				//pRTP = 0
				return err
			}
		} else {
			pRTP = int64(float64(pl.WinAmount) / float64(pl.BetAmount) * 100)
			if pl.BetAmount == 0 {
				pRTP = 0
			}
		}
		ret.List[i].PresentRTP = pRTP
		sel := bson.M{"_id": 0, "Uid": 1}
		single := coll.FindOne(context.TODO(), bson.M{"_id": play.Pid}, options.FindOne().SetProjection(sel))
		err = single.Decode(&ret.List[i])
		if err != nil {
			return
		}
		sel = bson.M{"_id": 0, "RTP": 1}
		_, ok := gameNameMap[play.GameID]
		if ok == false {
			gameNameMap[play.GameID] = "/"
		}
		ret.List[i].GameName = gameNameMap[play.GameID]
		ret.List[i].GameIcon = gameMap[play.GameID].IconUrl                      // 图标
		ret.List[i].GameManufacturerName = gameMap[play.GameID].ManufacturerName // 厂商
		ret.List[i].BuyType = gameMap[play.GameID].BuyType

		//operator gameRTP config
		if _, ok = mapOpConfig[play.AppID]; ok == true {
			for _, opc := range mapOpConfig[play.AppID] {
				if play.GameID == opc.GameId {
					ret.List[i].GameRTP = int64(opc.RTP)
				}
			}
		}

	}
	if ret.Count == 0 {
		ret.RTPCofigOff = 0
	}
	err = nil
	return
}

func CreatePlayerinfo(ctx *Context, ps CreatePlayerRTP, ret *comm.Empty) (err error) {
	// 1.发布地址   /player/setPlayerSettings_%s
	// 发布结构体
	//
	//	type PersonRtpSettings struct {
	//		RewardPercent  int   `json:"reward_percent"`
	//		NoAwardPercent int   `json:"no_award_percent"`
	//		PlayerId       int64 `json:"player_id"`
	//	}

	// 验证RTP可以设置的最大值
	err, b := comm.VerifyMaxRTP(ps.ContrllRTP)
	if b == false {
		return err
	}
	if ok, err := VerifyAuth(ctx, ps); !ok || err != nil {
		return err
	}
	err = createPlayerRTP(ps, ret)

	return
}

func SetDemoPlayerRTP(r *http.Request, ps CreatePlayerRTP, ret *comm.Empty) (err error) {

	ps.AppID = r.Header.Get("appid")

	if ps.AppID == "" {

		return errors.New("商户ID错误")
	}
	if ps.DemoUserId == "" {
		return errors.New("玩家不存在")
	}

	// -----------------------查询玩家信息--------------------
	var playerInfo player
	playerColl := db.Collection2("game", "Players")
	cour := playerColl.FindOne(context.TODO(), bson.M{"Uid": ps.DemoUserId})
	err = cour.Decode(&playerInfo)
	if err != nil {
		return err
	}

	// --------------------只有faketranss商户才能设置--------------------

	if (ps.ContrllRTP > 120 && playerInfo.AppID != "faketrans") ||
		(ps.ContrllRTP > 300 && playerInfo.AppID == "faketrans") {
		return errors.New("超出RTP最大值")
	}

	var gamePlayer map[string]string

	count := 0
	for gamePlayer == nil {
		time.Sleep(5 * time.Second)
		cur := db.Collection2(ps.GameID, "players").FindOne(context.TODO(), bson.M{"Uid": ps.DemoUserId})
		err := cur.Decode(&gamePlayer)
		if err != nil {
			fmt.Println("未找到")
			count++
		}

		if count > 5 {
			gamePlayer = map[string]string{}
			return errors.New("设置RTP失败")
		}
	}

	// ================ RTP控制表写入数据 ========================
	id, err := NextID(CollPlayerRTPControl, 1)

	var pl struct {
		BetAmount int64 `bson:"betamount"`
		WinAmount int64 `bson:"winamount"`
	}
	var msg struct {
		Pid    int64
		GameId string
		AppId  string
	}
	path := "/gamecenter/player/getPlayerInGame"
	msg.AppId = ps.AppID
	msg.GameId = ps.GameID
	msg.Pid = playerInfo.Id
	err = mq.Invoke(path, msg, &pl)

	var pRTP int64
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			pRTP = 0
		} else {
			//pRTP = 0
			return err
		}
	} else {
		pRTP = int64(float64(pl.WinAmount) / float64(pl.BetAmount) * 100)
		if pl.BetAmount == 0 {
			pRTP = 0
		}
	}

	rw, nw := GetGameNwandAw(ps.GameID, float64(ps.ContrllRTP))
	arw, anw := GetGameNwandAw(ps.GameID, float64(ps.AutoRemoveRTP))

	rTPControl := comm.PlayerRTPControlModel{
		ID:                 id,
		AppID:              ps.AppID,
		GameID:             ps.GameID,
		CreateAt:           time.Now(),
		Pid:                playerInfo.Id,
		GameName:           "",
		PlayerRTP:          pRTP,
		ControlRTP:         ps.ContrllRTP,
		RewardPercent:      rw,
		NoAwardPercent:     nw,
		AutoRemoveRTP:      ps.AutoRemoveRTP,
		AutoRewardPercent:  arw,
		AutoNoAwardPercent: anw,
		Status:             1,
		FromPlant:          "demo",
	}

	updata_rTPControl := PlayerRTPControlModel{}
	err = CollPlayerRTPControl.FindOne(bson.M{"GameID": ps.GameID, "Pid": playerInfo.Id, "FromPlant": "demo"}, &updata_rTPControl)
	if err != mongo.ErrNoDocuments {
		settingInfo := bson.M{
			"CreateAt":           rTPControl.CreateAt,
			"PlayerRTP":          rTPControl.PlayerRTP,
			"ControlRTP":         rTPControl.ControlRTP,
			"RewardPercent":      rTPControl.RewardPercent,
			"NoAwardPercent":     rTPControl.NoAwardPercent,
			"AutoRemoveRTP":      rTPControl.AutoRemoveRTP,
			"AutoRewardPercent":  rTPControl.AutoRewardPercent,
			"AutoNoAwardPercent": rTPControl.AutoNoAwardPercent,
			"Status":             rTPControl.Status,
		}
		set := bson.M{"$set": settingInfo}
		err = CollPlayerRTPControl.UpsertOne(bson.M{"GameID": ps.GameID, "Pid": playerInfo.Id, "FromPlant": "demo"}, set)
		if err != nil {
			return err
		}
	} else {

		err = CollPlayerRTPControl.InsertOne(rTPControl)
		if err != nil {
			return err
		}
	}

	// ================ 发出消息 ==========================

	var p = PlayerRtpSettings{
		RewardPercent:  int(rw),
		NoAwardPercent: int(nw),
		TargetRTP:      int(ps.ContrllRTP),
		RelieveRTP:     0,
		PlayerId:       playerInfo.Id,
	}

	sub := "/player/setPlayerSettings_" + ps.GameID
	_ = mq.PublishMsg(sub, p)

	time.Sleep(800 * time.Millisecond)

	return err
}

func createPlayerRTP(ps CreatePlayerRTP, ret *comm.Empty) (err error) {
	// 1.发布地址   /player/setPlayerSettings_%s
	// 发布结构体
	//
	//	type PersonRtpSettings struct {
	//		RewardPercent  int   `json:"reward_percent"`
	//		NoAwardPercent int   `json:"no_award_percent"`
	//		PlayerId       int64 `json:"player_id"`
	//	}

	if ps.AppIDandPlayerID == "" {
		return errors.New("请选择玩家")
	}
	if ps.GameID == "" {
		return errors.New("请选择游戏")
	}
	if ps.PersonWinMaxMult < minMult || ps.PersonWinMaxMult > maxMult {
		return errors.New("请输入正确的倍数")
	}
	if ps.PersonWinMaxScore < minScore || ps.PersonWinMaxScore > maxScore {
		return errors.New("请输入正确的最高钱数")
	}
	game_Filter := bson.M{}
	gameid := strings.Split(ps.GameID, ",")
	if gameid[0] != "ALL" {
		game_Filter["_id"] = bson.M{
			"$in": gameid,
		}
	}
	gamelist := []*comm.Game2{}
	err = NewOtherDB("game").Collection("Games").FindAll(game_Filter, &gamelist)
	if err != nil {
		return
	}

	appidAndPid := strings.Split(ps.AppIDandPlayerID, ",")
	appidlist := []string{}
	pidlist := []string{}
	for _, v := range appidAndPid {
		temp := strings.Split(v, ":")
		appidlist = append(appidlist, temp[0])
		pidlist = append(pidlist, temp[1])
	}
	var pl struct {
		BetAmount int64 `bson:"betamount"`
		WinAmount int64 `bson:"winamount"`
	}
	var msg struct {
		Pid    int64
		GameId string
		AppId  string
	}
	path := "/gamecenter/player/getPlayerInGame"
	go func() {
		for _, g := range gamelist {

			rw, nw := GetGameNwandAw(ps.GameID, float64(ps.ContrllRTP))
			arw, anw := GetGameNwandAw(ps.GameID, float64(ps.AutoRemoveRTP))

			for _, v := range appidAndPid {
				temp := strings.Split(v, ":")
				if temp[0] == "" {
					return
				}
				appid := temp[0]
				pid, err := strconv.ParseInt(temp[1], 10, 64)
				//pRTP, _ := getBetByWinRTP(appid, pid)
				msg.AppId = appid
				msg.GameId = g.ID
				msg.Pid = pid
				err = mq.Invoke(path, msg, &pl)
				var pRTP int64
				if err != nil {
					if err.Error() == "mongo: no documents in result" {
						pRTP = 0
					} else {
						//pRTP = 0
						return
					}
				} else {
					pRTP = int64(float64(pl.WinAmount) / float64(pl.BetAmount) * 100)
					if pl.BetAmount == 0 {
						pRTP = 0
					}
				}
				var result comm.Operator_V2

				err = CollAdminOperator.FindOne(bson.M{"AppID": appid}, &result)
				if err != nil {
					return
				}

				id, err := NextID(CollPlayerRTPControl, 1)
				rTPControl := comm.PlayerRTPControlModel{
					ID:                 id,
					AppID:              appid,
					GameID:             g.ID,
					CreateAt:           time.Now(),
					Pid:                pid,
					GameName:           g.Name,
					PlayerRTP:          pRTP,
					ControlRTP:         ps.ContrllRTP,
					RewardPercent:      rw,
					NoAwardPercent:     nw,
					AutoRemoveRTP:      ps.AutoRemoveRTP,
					AutoRewardPercent:  arw,
					AutoNoAwardPercent: anw,
					Status:             1,
					FromPlant:          "admin",
					PersonWinMaxMult:   ps.PersonWinMaxMult,
					PersonWinMaxScore:  ps.PersonWinMaxScore,
				}
				if result.BuyRTPOff != 0 {
					rTPControl.BuyRTP = ps.BuyRTP
					rTPControl.BuyMinAwardPercent = MapBuyGameRTP[ps.BuyRTP]
				}
				updata_rTPControl := comm.PlayerRTPControlModel{}
				err = CollPlayerRTPControl.FindOne(bson.M{"GameID": g.ID, "Pid": pid, "FromPlant": "admin"}, &updata_rTPControl)
				if err != mongo.ErrNoDocuments {
					settingInfo := bson.M{
						"CreateAt":           rTPControl.CreateAt,
						"PlayerRTP":          rTPControl.PlayerRTP,
						"ControlRTP":         rTPControl.ControlRTP,
						"RewardPercent":      rTPControl.RewardPercent,
						"NoAwardPercent":     rTPControl.NoAwardPercent,
						"AutoRemoveRTP":      rTPControl.AutoRemoveRTP,
						"AutoRewardPercent":  rTPControl.AutoRewardPercent,
						"AutoNoAwardPercent": rTPControl.AutoNoAwardPercent,
						"Status":             rTPControl.Status,
						"FromPlant":          "admin",
						"PersonWinMaxMult":   rTPControl.PersonWinMaxMult,
						"PersonWinMaxScore":  rTPControl.PersonWinMaxScore,
					}
					if result.BuyRTPOff != 0 {
						settingInfo["BuyRTP"] = rTPControl.BuyRTP
						settingInfo["BuyMinAwardPercent"] = rTPControl.BuyMinAwardPercent
					}
					set := bson.M{"$set": settingInfo}
					err = CollPlayerRTPControl.UpsertOne(bson.M{"GameID": g.ID, "Pid": pid}, set)
					if err != nil {
						return
					}
				} else {
					err = CollPlayerRTPControl.InsertOne(rTPControl)
					if err != nil {
						return
					}
				}
				// BuyMinAwardPercent值不准确当为-1不接收方接收到消息应该忽略BuyMinAwardPercent的值使用旧值或者数据库中的值
				var p = PlayerRtpSettings{
					RewardPercent:  int(rw),
					NoAwardPercent: int(nw),
					TargetRTP:      int(ps.ContrllRTP),
					RelieveRTP:     int(ps.AutoRemoveRTP),
					PlayerId:       pid,
				}
				if result.BuyRTPOff != 0 {
					p.BuyMinAwardPercent = -1
				}

				sub := "/player/setPlayerSettings_" + g.ID
				_ = mq.PublishMsg(sub, p)
			}
		}

	}()
	err = nil
	return
}

func CancelPlayerRTP(ctx *Context, ps CancelRTP, ret *comm.Empty) (err error) {
	// 1.发布地址   /player/setPlayerSettings_%s
	// 发布结构体
	//
	//	type PersonRtpSettings struct {
	//		RewardPercent  int   `json:"reward_percent"`
	//		NoAwardPercent int   `json:"no_award_percent"`
	//		PlayerId       int64 `json:"player_id"`
	//	}

	sIdList := strings.Split(ps.IDList, ",")
	iIdList := []int64{}
	for _, s := range sIdList {
		temp, _ := strconv.ParseInt(s, 10, 64)
		iIdList = append(iIdList, temp)
	}
	filter := bson.M{"_id": bson.M{"$in": iIdList}}
	set := bson.M{"$set": bson.M{"Status": 0}}
	err = CollPlayerRTPControl.Update(filter, set)
	if err != nil {
		return
	}

	go func() {
		var playerRpt []*comm.PlayerRTPControlModel
		err = CollPlayerRTPControl.FindAll(bson.M{"_id": bson.M{"$in": iIdList}}, &playerRpt)
		if err != nil {
			slog.Error("主动取消RTP错误：playRTPControl\"$in\": ", iIdList)
		}
		for _, play := range playerRpt {
			channel.TaskQueue <- channel.RtpSettings{
				GameId:         play.GameID,
				RewardPercent:  0,
				NoAwardPercent: 0,
				PlayerId:       play.Pid,
			}
			//sub := "/player/setPlayerSettings_" + play.GameID
			//_ = mq.PublishMsg(sub, oper)
		}
	}()

	return
}

func CancelPlayerRTPByPId(ctx *Context, ps CancelRTPByPId, ret *comm.Empty) (err error) {
	// 1.发布地址   /player/setPlayerSettings_%s
	// 发布结构体
	//
	//	type PersonRtpSettings struct {
	//		RewardPercent  int   `json:"reward_percent"`
	//		NoAwardPercent int   `json:"no_award_percent"`
	//		PlayerId       int64 `json:"player_id"`
	//	}
	//pid, err := strconv.Atoi(ps.Pid)
	filter := bson.M{"Pid": ps.Pid}
	set := bson.M{"$set": bson.M{"Status": 0}}
	err = CollPlayerRTPControl.Update(filter, set)
	if err != nil {
		return
	}

	go func() {
		var playerRpt []*comm.PlayerRTPControlModel
		err = CollPlayerRTPControl.FindAll(bson.M{"Pid": ps.Pid}, &playerRpt)
		if err != nil {
			slog.Error("主动取消RTP错误：playRTPControl\"Pid = \": ", ps.Pid)
		}
		for _, play := range playerRpt {
			channel.TaskQueue <- channel.RtpSettings{
				GameId:         play.GameID,
				RewardPercent:  0,
				NoAwardPercent: 0,
				PlayerId:       play.Pid,
			}
			//sub := "/player/setPlayerSettings_" + play.GameID
			//_ = mq.PublishMsg(sub, oper)
		}
	}()

	return
}

type playerbetSum struct {
	Bet int64 `bson:"Bet"` // 下注
	Win int64 `bson:"Win"` // 单次消除输赢
}
