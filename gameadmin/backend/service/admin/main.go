package main

import (
	"context"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/mq"
	"game/comm/mux"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/lazy"
	"game/duck/logger"
	"game/duck/mongodb"
	"game/duck/rpc1"
	"game/pb/_gen/pb/adminpb"
	"game/service/admin/channel"
	"log"
	"log/slog"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/encoding"
)

func main() {
	now := time.Now()
	fmt.Println("服务器时区：", strings.Split(now.String(), " ")[2])
	lazy.Init("admin")
	clickHouseAddr := lo.Must(lazy.RouteFile.Get("clickhouse"))

	LoadCfg()
	ConnectMongoByConfig()
	err := db.DialToClickHouse(clickHouseAddr)
	if err != nil {
		log.Fatalf("Failed to conn ClickHouse: %v", err)
	}
	initEnsureIndex()

	DownBetHistory()

	// 添加加载  游戏相关配置
	initGameConfig()

	must, err := mq.ConnectServerMust()
	if err != nil {
		return
	}
	mux.RegistRpcToMQ(must)
	initBase1()
	UpdateMenList()

	encoding.RegisterCodec(rpc1.JsonBytesCodec{})

	adminpb.RegisterAdminAuthServer(lazy.GrpcServer, &AdminAuthRpc{})
	adminpb.RegisterAdminConfigFileServer(lazy.GrpcServer, &AdminConfigFile{})
	adminpb.RegisterAdminInfoServer(lazy.GrpcServer, &AdminInfoRpc{})
	adminpb.RegisterAdminGroupServer(lazy.GrpcServer, &AdminGroupRpc{})
	adminpb.RegisterAdminMenuServer(lazy.GrpcServer, &AdminMenuRpc{})
	// 订阅消息
	RegSub()
	go channel.PushPlayersetPlayerSettings()

	//go addPlayerRetentionReport()

	DataDeal()
	//addFakeOPerator()
	go func() {
		err := operatorSyncGames()
		if err != nil {
			fmt.Println("启动同步商户脚本错误：" + err.Error())
			return
		}
	}()
	go resetGameManufacturer()
	lazy.ServeFn(httpServe)
}

func httpServe() {

	// mx := http.NewServeMux()
	mx := mux.DefaultRpcMux.ServeMux
	mx.Handle("/", NewHttpHandler(&checker{}, "admin"))
	mx.Handle("/mq/", http.StripPrefix("/mq", http.HandlerFunc(adminMQHandler)))
	mx.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	mx.HandleFunc("/AdminInfo/BHdownload", doDownloadBetHistoryExcle)
	mx.HandleFunc("/AdminInfo/AddLangToDB", uploadFiles)
	mx.HandleFunc("/comm/saveUpload", saveUploadFile)
	mx.HandleFunc("/AdminInfo/uploadGame", uploadGameFile)
	mx.HandleFunc("/AdminInfo/uploadConfig", uploadConfig)
	// mx.HandleFunc("/AdminInfo/GetGameEarningsData/DownLoad", DownGameEarningsData)

	fs := http.FileServer(http.Dir("./BHdownload"))
	http.Handle("/BHdownload/", http.StripPrefix("/BHdownload/", fs))

	// RegistHandlers()

	port, err := lazy.PortProvider.GetPort("admin.http")
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf(":%v", port)
	// err = http.ListenAndServe(fmt.Sprintf(":%v", port), mx)
	// if err != nil {
	// 	panic(err)
	// }
	mux.DefaultRpcMux.StartHttpServer(addr)

}

type checker struct{}

func (c *checker) GetPidAndUname(token string) (int64, string) {
	tk := getToken(token)
	if tk == nil {
		return 0, ""
	}
	return tk.Pid, tk.UserName
}

func (c *checker) IsAllowAndShoudLogin(path *Path) (allow bool, shouldLogin bool) {
	// if path.Service == "/list_api" {
	// 	allow = true
	// 	shouldLogin = false
	// 	return
	// }

	// 这里有包名 baozang.XXRPC
	service := path.Service
	arr := strings.Split(path.Service, ".")
	if len(arr) == 2 {
		service = arr[1]
	}

	allow = strings.HasPrefix(service, "Admin")

	shouldLogin = service != "AdminAuth"
	return
}

type DataReport struct {
	Bet int64 `bson:"betamount"`
	Win int64 `bson:"winamount"`
}

func DataDeal() {

	//msg := YesterdayTotal()
	//slog.Info("YesterDayTotal", "data", msg)
	//err := slotsmongo.SendTelegramChat(msg)
	//if err != nil {
	//	slog.Warn("SendTelegramChat", "error", ut.ErrString(err))
	//}

	//支持秒级
	//c := cron.New(cron.WithSeconds())

	c := cron.New()

	//report.Betlog  5分钟一次 删除5个月前数据
	//betlog 5分钟一次 3个月之前的
	//游戏库中 BetHistory 3个月的数据
	c.AddFunc("0 5 * * *", func() {
		//slotsmongo.DelExpiredTTLLogColls()
		//
		//coll := db.Collection2("reports", "BetLog")
		//// 通用的下注历史
		//before := time.Now().AddDate(0, -6, 0)
		//res, _ := coll.DeleteMany(context.TODO(), bson.M{
		//	"InsertTime": bson.M{"$lte": before},
		//})
		//slog.Info("delete bet log", "count", res.DeletedCount)

		before := time.Now().AddDate(0, -4, 0)
		collNames, _ := db.Client().Database("betlog").ListCollectionNames(context.TODO(), bson.M{})
		for i := range collNames {
			collN := db.Collection2("betlog", collNames[i])
			res, _ := collN.DeleteMany(context.TODO(), bson.M{
				"InsertTime": bson.M{"$lte": before},
			})
			slog.Info("delete bet log", "count", res.DeletedCount)
		}
		coll := db.Collection2("reports", "PlayerRetentionReport")
		// 通用的下注历史
		before = time.Now().AddDate(0, -4, 0)
		res, _ := coll.DeleteMany(context.TODO(), bson.M{
			"date": bson.M{"$lte": before},
		})
		slog.Info("delete player retentionReport log", "count", res.DeletedCount)
		coll = db.Collection2("reports", "PlayerDailyRetention")
		// 通用的下注历史
		before = time.Now().AddDate(0, -4, 0)
		res, _ = coll.DeleteMany(context.TODO(), bson.M{
			"date": bson.M{"$lte": before},
		})
		slog.Info("delete player DailyRetention log", "count", res.DeletedCount)
		before = time.Now().AddDate(0, -4, 0)
		databases, _ := db.Client().ListDatabaseNames(context.TODO(), bson.D{})
		beforeObjId := primitive.NewObjectIDFromTimestamp(before)
		for i := range databases {
			coll := db.Collection2(databases[i], "BetHistory")
			// if strings.HasPrefix(databases[i], "pg_") {
			res, err := coll.DeleteMany(context.TODO(), bson.M{
				"_id": bson.M{"$lte": beforeObjId},
			})
			if err == nil {
				slog.Info("delete bet log", "gameid", databases[i], "count", res.DeletedCount)
			}
			// }
		}
	})
	//
	c.AddFunc("@hourly", func() {
		if setting.Alert == nil || setting.Alert.Hour {
			msg := HourReports()
			slog.Info("HourTotal", "data", msg)
			err := slotsmongo.SendTelegramChat(msg)
			if err != nil {
				slog.Warn("SendTelegramChat", "error", ut.ErrString(err))
			}
		}
	})
	updateGameBet(true)
	c.AddFunc("0 5 1 1,4,7,10 *", func() {
		updateGameBet(false)
	})
	c.AddFunc("30 9 * * *", func() {
		removeBetLogDownload()  // 删除下注历史CSV文件(保留7天)
		removeGameLoginDetail() // 删除IP访问数据(保留7天)
	})

	c.AddFunc("0 12 L * ?", func() {
		//c.AddFunc("*/5 * * * * ?", func() {
		SyncPlantRate()
	})
	//
	c.AddFunc("0 1 * * *", func() {
		err := addOperatorMonthlyGGR()
		if err != nil {
			slog.Error("update reports OperatorMonthlyGGR err:  ", err)
		}
		err = addBetDailyGGR()
		if err != nil {
			slog.Error("update reports OperatorDailyReport err:  ", err)
		}
	})
	// 数据留存每小时更新一次
	c.AddFunc("0 * * * *", func() {
		err := addPlayerRetentionReport()
		if err != nil {
			slog.Error("update reports PlayerRetentionReport err:  ", err)
		}
	})
	// 数据每天0点更新一次
	c.AddFunc("0 0 * * *", func() {
		BalanceAlert(24)
	})
	// 数据每12小时更新一次
	c.AddFunc("0 0 */12 * *", func() {
		BalanceAlert(12)
	})
	// 数据每5小时更新一次
	c.AddFunc("0 0 */5 * *", func() {
		BalanceAlert(5)
	})
	// 数据每3小时更新一次
	c.AddFunc("0 0 */3 * *", func() {
		BalanceAlert(3)
	})
	// 数据每小时更新一次
	c.AddFunc("0 * * * *", func() {
		BalanceAlert(1)
	})
	c.Start()
}

// 插入betDailyReport脚本
func initOperatorMonthlyGGR() (err error) {

	slog.Info("---------------------addOperatorMonthlyGGR begin---------------------")

	oper := []*comm.Operator_V2{}
	CollAdminOperator.FindAll(bson.M{}, &oper)
	mapRale := make(map[string]float64)
	for _, op := range oper {
		mapRale[op.AppID] = op.PlatformPay
	}
	now := time.Now().Format("20060102")

	betDateReport := []*comm.BetDailyGGR{}

	NewOtherDB("reports").Collection("BetDailyReport").FindAll(bson.M{"date": bson.M{"$lt": now}}, &betDateReport)
	coll := NewOtherDB("reports").Collection("OperatorMonthlyGGR").Coll()
	for i, ggr := range betDateReport {
		rate := mapRale[ggr.AppID]
		win := betDateReport[i].WinAmount
		bet := betDateReport[i].BetAmount

		GGR := float64(bet-win) * (rate * 0.01)
		filter := bson.M{
			"appid": ggr.AppID,
			"date":  ggr.Date[:6],
		}

		id := fmt.Sprintf("%s:%s", ggr.Date[:6], ggr.AppID)

		update := bson.M{
			"$inc": bson.M{
				"betamount": ggr.BetAmount,
				"winamount": ggr.WinAmount,
				"GGR":       GGR,
			},
			"$setOnInsert": bson.M{ // 在插入时设置默认值
				"_id":       id,
				"date":      ggr.Date[:6],
				"plantrate": rate,
			},
		}

		opts := options.Update().SetUpsert(true)

		_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			return err
		}

	}
	slog.Info("---------------------addOperatorMonthlyGGR End---------------------")
	return

}

// 月operator数据汇总，没有今天的，对账单用
func addOperatorMonthlyGGR() (err error) {

	slog.Info("---------------------addOperatorMonthlyGGR begin---------------------")

	oper := []*comm.Operator_V2{}
	CollAdminOperator.FindAll(bson.M{}, &oper)
	mapRale := make(map[string]float64)
	mapParent := make(map[string]string)
	mapCooperation := make(map[string]int)
	mapTurnoverRale := make(map[string]float64)
	for _, op := range oper {
		mapRale[op.AppID] = op.PlatformPay
		mapParent[op.AppID] = op.Name
		mapCooperation[op.AppID] = op.CooperationType
		mapTurnoverRale[op.AppID] = op.TurnoverPay
	}
	now := time.Now().AddDate(0, 0, -1).Format("20060102")

	betDateReport := []*comm.BetDailyGGR{}

	NewOtherDB("reports").Collection("BetDailyReport").FindAll(bson.M{"date": now}, &betDateReport)
	coll := NewOtherDB("reports").Collection("OperatorMonthlyGGR").Coll()
	for i, ggr := range betDateReport {
		rate := mapRale[ggr.AppID]
		if mapCooperation[ggr.AppID] == 2 {
			rate = mapTurnoverRale[ggr.AppID]
		}
		win := betDateReport[i].WinAmount
		bet := betDateReport[i].BetAmount

		GGR := float64(0.0)
		InCome := float64(0.0)
		if mapParent[ggr.AppID] != "admin" {
			if mapCooperation[ggr.AppID] == 1 {
				parentRate := mapRale[mapParent[ggr.AppID]]
				InCome = float64(bet-win) * ((rate - parentRate) * 0.01)
				GGR = float64(bet-win) * (rate * 0.01)
			} else {
				parentRate := mapTurnoverRale[mapParent[ggr.AppID]]
				InCome = float64(bet) * ((rate - parentRate) * 0.01)
				GGR = float64(bet) * (rate * 0.01)
			}
		}

		filter := bson.M{
			"appid": ggr.AppID,
			"date":  ggr.Date[:6],
		}

		id := fmt.Sprintf("%s:%s", ggr.Date[:6], ggr.AppID)

		update := bson.M{
			"$inc": bson.M{
				"betamount": ggr.BetAmount,
				"winamount": ggr.WinAmount,
				"GGR":       GGR,
				"InCome":    InCome,
			},
			"$set": bson.M{
				"plantrate": rate,
			},
			"$setOnInsert": bson.M{ // 在插入时设置默认值
				"_id":  id,
				"date": ggr.Date[:6],
				//"plantrate": rate,
			},
		}

		opts := options.Update().SetUpsert(true)

		_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			return err
		}

	}
	//query := bson.M{
	//	"date": now,
	//}
	//operDateReport := []*comm.OperatorDailyGGR{}
	//NewOtherDB("reports").Collection("OperatorDailyReport").FindAll(query, &operDateReport)
	//coll := NewOtherDB("reports").Collection("OperatorMonthlyGGR").Coll()
	//for i, daily := range operDateReport {
	//	rate := mapRale[daily.AppID]
	//	win := operDateReport[i].WinAmount
	//	bet := operDateReport[i].BetAmount
	//	operatorBalance := operDateReport[i].OperatorBalance
	//
	//	GGR := float64(bet-win) * (rate * 0.01)
	//	filter := bson.M{
	//		"appid": daily.AppID, := fmt.Sprintf("%s:%s", daily.Date[:6], daily.AppID)
	//	//	update := bson.M{
	//	//		"$inc": bson.M{
	//	//			"betamount": daily.BetAmount,
	//	//			"winamount": daily.WinAmount,
	//	//			"GGR":       GGR,
	//	//		},
	//	//		"$set": bson.M{
	//	//			"plantrate":       rate,
	//	//			"OperatorBalance": operatorBalance,
	//	//		},
	//	//		"$setOnInsert": bson.M{ // 在插入时设置默认值
	//	//			"_id":  id,
	//	//			"date": daily.Date[:6],
	//	//			//"plantrate": rate,
	//	//		},
	//	//	}
	//	//
	//	//	opts := options.Update().SetUpsert(true)
	//		"date":  daily.Date[:6],
	//	}
	//	id
	//
	//	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	//	if err != nil {
	//		return err
	//	}
	//}
	slog.Info("---------------------addOperatorMonthlyGGR End---------------------")
	return

}

type result struct {
	Pids []int64 `bson:"pids"`
}

// 定义一个函数来从slice中去除重复元素
func RemoveDuplicates(slice []*result) []int64 {
	// 使用map来记录元素是否已经见过
	m := make(map[int64]bool)
	res := make([]int64, 0)

	// 遍历原始slice，检查元素是否已经存在于map中
	for _, v := range slice {
		for _, v1 := range v.Pids {
			if _, ok := m[v1]; !ok {
				m[v1] = true
				res = append(res, v1)
			}
		}
	}

	return res
}

// 添加玩家留存数据
func addPlayerRetentionReport() (err error) {

	slog.Info("---------------------addPlayerRetentionReport begin---------------------")

	oper := []*comm.Operator_V2{}
	CollAdminOperator.FindAll(bson.M{"OperatorType": 2}, &oper)
	query := bson.M{}
	nowStr := time.Now().Format("20060102")
	now1d := time.Now().Truncate(24 * time.Hour)
	now30d := time.Now().AddDate(0, 0, -31).Truncate(24 * time.Hour)
	query["date"] = bson.M{
		"$gte": now30d,
		"$lte": now1d,
	}
	retentionReport := []*comm.RetentionReport{}
	coll := db.Collection2("reports", "PlayerRetentionReport")
	err = NewOtherDB("reports").Collection("PlayerRetentionReport").FindAll(query, &retentionReport)
	if err != nil {
		return err
	}
	collDaily := NewOtherDB("reports").Collection("PlayerDailyRetention")
	for _, r := range retentionReport {
		query1 := bson.M{}
		query3 := bson.M{}
		query7 := bson.M{}
		query14 := bson.M{}
		query30 := bson.M{}
		dateTime := r.Date.AsTime()
		dateStr := dateTime.Format("20060102")
		query = bson.M{}
		dateStr1 := dateTime.AddDate(0, 0, 1).Format("20060102")
		dateStr3 := dateTime.AddDate(0, 0, 3).Format("20060102")
		dateStr7 := dateTime.AddDate(0, 0, 7).Format("20060102")
		dateStr14 := dateTime.AddDate(0, 0, 14).Format("20060102")
		dateStr30 := dateTime.AddDate(0, 0, 30).Format("20060102")
		id := ut.JoinStr(':', dateStr, r.AppID)
		id1 := ut.JoinStr(':', dateStr1, r.AppID)
		id3 := ut.JoinStr(':', dateStr3, r.AppID)
		id7 := ut.JoinStr(':', dateStr7, r.AppID)
		id14 := ut.JoinStr(':', dateStr14, r.AppID)
		id30 := ut.JoinStr(':', dateStr30, r.AppID)
		query["_id"] = primitive.Regex{
			Pattern: "^" + id + ":",
		}
		query1["_id"] = primitive.Regex{
			Pattern: "^" + id1 + ":",
		}
		query3["_id"] = primitive.Regex{
			Pattern: "^" + id3 + ":",
		}
		query7["_id"] = primitive.Regex{
			Pattern: "^" + id7 + ":",
		}
		query14["_id"] = primitive.Regex{
			Pattern: "^" + id14 + ":",
		}
		query30["_id"] = primitive.Regex{
			Pattern: "^" + id30 + ":",
		}
		if r.GameID != "" {
			query["gameId"] = r.GameID
			query1["gameId"] = r.GameID
			query3["gameId"] = r.GameID
			query7["gameId"] = r.GameID
			query14["gameId"] = r.GameID
			query30["gameId"] = r.GameID
		}

		results := []*result{}
		pidList := []int64{}
		collDaily.FindAll(query, &results)
		for _, p := range results {
			pidList = append(pidList, p.Pids...)
		}
		m := make(map[interface{}]bool)
		count := 0

		// 遍历pidList，检查元素是否已经存在于map中
		for _, v := range pidList {
			if _, ok := m[v]; !ok {
				m[v] = true
				count++
			}
		}
		retentionPlayer1d := float64(0)  // 次日留存率
		retentionPlayer3d := float64(0)  // 3日留存率
		retentionPlayer7d := float64(0)  // 7日留存率
		retentionPlayer14d := float64(0) // 14日留存率
		retentionPlayer30d := float64(0) // 30日留存率

		if count != 0 {
			results = []*result{}
			pidList = []int64{}
			collDaily.FindAll(query1, &results)
			pidList = RemoveDuplicates(results)
			count1 := 0

			// 遍历pidList，检查元素是否已经存在于map中
			for _, v := range pidList {
				if _, ok := m[v]; ok {
					count1++
				}
			}

			results = []*result{}
			pidList = []int64{}
			collDaily.FindAll(query3, &results)
			pidList = RemoveDuplicates(results)
			count3 := 0

			// 遍历pidList，检查元素是否已经存在于map中
			for _, v := range pidList {
				if _, ok := m[v]; ok {
					count3++
				}
			}

			results = []*result{}
			pidList = []int64{}
			collDaily.FindAll(query7, &results)
			pidList = RemoveDuplicates(results)
			count7 := 0

			// 遍历pidList，检查元素是否已经存在于map中
			for _, v := range pidList {
				if _, ok := m[v]; ok {
					count7++
				}
			}

			results = []*result{}
			pidList = []int64{}
			collDaily.FindAll(query14, &results)
			pidList = RemoveDuplicates(results)
			count14 := 0

			// 遍历pidList，检查元素是否已经存在于map中
			for _, v := range pidList {
				if _, ok := m[v]; ok {
					count14++
				}
			}

			results = []*result{}
			pidList = []int64{}
			collDaily.FindAll(query30, &results)
			pidList = RemoveDuplicates(results)
			count30 := 0

			// 遍历pidList，检查元素是否已经存在于map中
			for _, v := range pidList {
				if _, ok := m[v]; ok {
					count30++
				}
			}
			retentionPlayer1d = math.Round(float64(count1)/float64(count)*10000) / 100
			retentionPlayer3d = math.Round(float64(count3)/float64(count)*10000) / 100
			retentionPlayer7d = math.Round(float64(count7)/float64(count)*10000) / 100
			retentionPlayer14d = math.Round(float64(count14)/float64(count)*10000) / 100
			retentionPlayer30d = math.Round(float64(count30)/float64(count)*10000) / 100
		}
		updateRetention := db.D(
			"$set", db.D(
				"RetentionPlayer1d", retentionPlayer1d,
				"RetentionPlayer3d", retentionPlayer3d,
				"RetentionPlayer7d", retentionPlayer7d,
				"RetentionPlayer14d", retentionPlayer14d,
				"RetentionPlayer30d", retentionPlayer30d,
			))

		_, err = coll.UpdateByID(context.TODO(), r.ID, updateRetention)
		if err != nil {
			slog.Info("updatePayerRetention", "err", ut.ErrString(err))
		}
	}
	for _, op := range oper {
		id := ut.JoinStr(':', nowStr, op.AppID)
		coll.InsertOne(context.TODO(), &comm.PlayerRetentionReport{ID: id, AppID: op.AppID, Date: mongodb.NowTimeStamp()})
	}
	slog.Info("---------------------addPlayerRetentionReport End---------------------")
	return

}

// 商户每日 每月 数据汇总
func addBetDailyGGR() (err error) {
	slog.Info("---------------------addBetDailyGGR begin---------------------")

	oper := []*comm.Operator_V2{}
	CollAdminOperator.FindAll(bson.M{}, &oper)
	mapRale := make(map[string]float64)
	mapCooperstion := make(map[string]int)
	mapTurnoverRale := make(map[string]float64)
	for _, op := range oper {
		mapRale[op.AppID] = op.PlatformPay
		mapCooperstion[op.AppID] = op.CooperationType
		mapTurnoverRale[op.AppID] = op.TurnoverPay
	}
	yesterday := time.Now().UTC().AddDate(0, 0, -1)

	betDateReport := []*comm.BetDailyGGR{}
	NewOtherDB("reports").Collection("BetDailyReport").FindAll(bson.M{"date": yesterday.Format("20060102")}, &betDateReport)
	//NewOtherDB("reports").Collection("BetDailyReport").FindAll(bson.M{}, &betDateReport)
	for i, daily := range betDateReport {
		rate := mapRale[daily.AppID]
		win := betDateReport[i].WinAmount
		bet := betDateReport[i].BetAmount

		betDateReport[i].PresentRate = rate
		betDateReport[i].GGR = float64(bet-win) * (rate * 0.01)

	}
	err = NewOtherDB("reports").Collection("BetDailyGGR").InsertMany(lo.ToAnySlice(betDateReport))
	if err != nil {
		return err
	}
	fmt.Println("success insert")

	//OperatorDailyGGR 每日汇总要用,将昨天的数据插入
	operDateReport := []*comm.OperatorDailyGGR{}
	query := bson.M{
		"date": yesterday.Format("20060102"),
	}
	NewOtherDB("reports").Collection("OperatorDailyReport").FindAll(query, &operDateReport)

	if err != nil {
		return
	}

	for i, daily := range operDateReport {
		rate := mapRale[daily.AppID]
		win := operDateReport[i].WinAmount
		bet := operDateReport[i].BetAmount
		GGR := float64(bet-win) * (rate * 0.01)
		if mapCooperstion[daily.AppID] == 2 {
			rate = mapTurnoverRale[daily.AppID]
			GGR = float64(bet) * (rate * 0.01)
		}
		operDateReport[i].PresentRate = rate
		operDateReport[i].GGR = GGR
	}
	err = NewOtherDB("reports").Collection("OperatorDailyGGR").InsertMany(lo.ToAnySlice(operDateReport))
	if err != nil {
		return err
	}
	slog.Info("success insert")
	return
}

func HourReports() string {
	coll := db.Collection2("reports", "BetDailyReport")
	getData := func(startT, endT string) DataReport {
		res := make([]DataReport, 0, 10)
		match := bson.M{
			"date": bson.M{
				"$gte": startT,
				"$lte": endT,
			},
		}

		cursor, err := coll.Aggregate(context.TODO(), []bson.M{
			{
				"$match": match,
			}, {
				"$group": bson.M{
					"_id":       "1",
					"winamount": bson.M{"$sum": "$winamount"},
					"betamount": bson.M{"$sum": "$betamount"},
				},
			},
		})
		if err != nil {
			slog.Warn("db cursor", "error", ut.ErrString(err))
			return DataReport{}
		}
		err = cursor.All(context.TODO(), &res)
		if err != nil {
			slog.Warn("db unmarshal", "error", ut.ErrString(err))
			return DataReport{}
		}
		if len(res) == 0 {
			return DataReport{}
		}
		return res[0]
	}

	addReport := bson.M{}
	str := strings.Builder{}
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:00:00")
	addReport["time"] = nowStr
	str.WriteString(fmt.Sprintf("时间：%s\n", nowStr))
	if now.Hour() == 0 {
		now = now.Add(-time.Hour)
	}
	date := now.Format("20060102")

	data := getData(date, date)
	addReport["bet"] = data.Bet
	addReport["win"] = data.Win
	addReport["rate"] = fmt.Sprintf("%.2f%%", float64(data.Win)/float64(data.Bet)*100)
	str.WriteString("总打码量：")
	str.WriteString(addThousandSeparator(fmt.Sprintf("%.1f\n", ut.Gold2Money(data.Bet))))
	str.WriteString("系统收益：")
	str.WriteString(addThousandSeparator(fmt.Sprintf("%.1f\n", ut.Gold2Money(data.Bet-data.Win))))
	str.WriteString(fmt.Sprintf("当日回报率：%.3f%%\n", float64(data.Win)/float64(data.Bet)*100))

	//七日
	nowB7 := now.AddDate(0, 0, -6).Format("20060102")
	data = getData(nowB7, date)
	addReport["bet7"] = data.Bet
	addReport["win7"] = data.Win
	addReport["rate7"] = fmt.Sprintf("%.2f%%", float64(data.Win)/float64(data.Bet)*100)
	str.WriteString(fmt.Sprintf("七日回报率：%.4f%%\n", float64(data.Win)/float64(data.Bet)*100))

	//月
	monthFirst := fmt.Sprintf("%s00", now.Format("200601"))
	monthEnd := fmt.Sprintf("%s31", now.Format("200601"))
	data = getData(monthFirst, monthEnd)
	addReport["betM"] = data.Bet
	addReport["winM"] = data.Win
	addReport["rateM"] = fmt.Sprintf("%.2f%%", float64(data.Win)/float64(data.Bet)*100)
	str.WriteString(fmt.Sprintf("当月回报率：%.4f%%", float64(data.Win)/float64(data.Bet)*100))

	coll = db.Collection2("reports", "hourReports")
	coll.InsertOne(context.TODO(), addReport)
	return str.String()
}

type BetInfo struct {
	Game string `bson:"_id"`
	Bet  int64  `bson:"betamount"`
}

func updateGameBet(checkTime bool) {
	var lastTime time.Time
	db.MiscGet2("GameAdmin", "GameBet", &lastTime)
	if lastTime.AddDate(0, 3, 0).After(time.Now()) && checkTime {
		return
	}
	now := time.Now()
	coll := db.Collection2("reports", "BetDailyReport")
	match := bson.M{
		"date": bson.M{
			"$gte": now.AddDate(0, 0, -6).Format("20060102"),
			"$lte": now.Format("20060102"),
		},
	}
	cursor, err := coll.Aggregate(context.TODO(), []bson.M{
		{
			"$match": match,
		}, {
			"$group": bson.M{
				"_id":       "$game",
				"betamount": bson.M{"$sum": "$betamount"},
			},
		},
	})
	if err != nil {
		return
	}
	res := make([]BetInfo, 0, 10)
	err = cursor.All(context.TODO(), &res)
	if err != nil {
		slog.Warn("db unmarshal", "error", ut.ErrString(err))
		return
	}
	if len(res) == 0 {
		return
	}
	getBet := func(id string) int64 {
		for i := range res {
			if res[i].Game == id {
				return res[i].Bet
			}
		}
		return 0
	}
	coll2 := db.Collection2("game", "Games")
	cur, _ := coll2.Find(context.TODO(), db.D(), options.Find().SetProjection(db.D("_id", 1)))
	models := []mongo.WriteModel{}
	for cur.Next(context.TODO()) {
		var doc struct {
			ID string `bson:"_id"`
		}
		err = cur.Decode(&doc)
		if err != nil {
			return
		}
		model := mongo.NewUpdateOneModel().SetFilter(db.ID(doc.ID)).SetUpdate(db.D("$set", db.D("Bet", getBet(doc.ID))))
		models = append(models, model)
	}
	coll2.BulkWrite(context.TODO(), models, options.BulkWrite().SetOrdered(false))
	db.MiscSet2("GameAdmin", "GameBet", now)
}

func YesterdayTotal() string {
	coll := db.Collection2("reports", "BetDailyReport")
	getData := func(startT, endT string) DataReport {
		res := make([]DataReport, 0, 10)
		match := bson.M{
			"date": bson.M{
				"$gte": startT,
				"$lte": endT,
			},
		}

		cursor, err := coll.Aggregate(context.TODO(), []bson.M{
			{
				"$match": match,
			}, {
				"$group": bson.M{
					"_id":       "1",
					"winamount": bson.M{"$sum": "$winamount"},
					"betamount": bson.M{"$sum": "$betamount"},
				},
			},
		})
		if err != nil {
			slog.Warn("db cursor", "error", ut.ErrString(err))
			return DataReport{}
		}
		err = cursor.All(context.TODO(), &res)
		if err != nil {
			slog.Warn("db unmarshal", "error", ut.ErrString(err))
			return DataReport{}
		}
		if len(res) == 0 {
			return DataReport{}
		}
		return res[0]
	}
	yesterday := time.Now().AddDate(0, 0, -1)

	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("时间：%s\n", yesterday.Format("2006-01-02")))
	// 昨天
	date := yesterday.Format("20060102")
	data := getData(date, date)
	str.WriteString("总打码量：")
	str.WriteString(addThousandSeparator(fmt.Sprintf("%.1f\n", ut.Gold2Money(data.Bet))))
	str.WriteString("系统总收益：")
	str.WriteString(addThousandSeparator(fmt.Sprintf("%.1f\n", ut.Gold2Money(data.Bet-data.Win))))
	str.WriteString(fmt.Sprintf("当日回报率：%.4f%%\n", float64(data.Win)/float64(data.Bet)*100))

	// 七日
	nowB7 := yesterday.AddDate(0, 0, -6).Format("20060102")
	data = getData(nowB7, date)
	str.WriteString(fmt.Sprintf("七日回报率：%.4f%%\n", float64(data.Win)/float64(data.Bet)*100))

	// 当月
	monthFirst := fmt.Sprintf("%s00", yesterday.Format("200601"))
	monthEnd := fmt.Sprintf("%s31", yesterday.Format("200601"))
	data = getData(monthFirst, monthEnd)
	str.WriteString(fmt.Sprintf("当月回报率：%.4f%%", float64(data.Win)/float64(data.Bet)*100))

	return str.String()
}

func addThousandSeparator(s string) string {
	numStrs := strings.Split(s, ".")
	numStr := numStrs[0]
	length := len(numStr)
	if length < 4 {
		return s
	}
	cnt := (length - 1) / 3
	for i := 0; i < cnt; i++ {
		numStr = numStr[:length-(i+1)*3] + "," + numStr[length-(i+1)*3:]
	}
	if len(numStrs) == 1 {
		return numStr
	}
	return fmt.Sprintf("%s.%s", numStr, numStrs[1])
}
func initEnsureIndex() {
	coll := db.Collection2("GameAdmin", "SlotsPoolHistory")
	indexmodels := []mongo.IndexModel{
		{
			Keys: db.D("AnimUserPid", 1),
		},
	}
	coll.Indexes().CreateMany(context.TODO(), indexmodels)
}

func initGameConfig() {

	err, gameList := GetGameListFormatName("", bson.M{})
	if err != nil {
		logger.Fatalf("加载游戏配置失败：获取游戏列表错误(%s)", err.Error())
		return
	}

	var GameBetList []*GameBet
	var GameNowList []*GameNow

	for _, game := range gameList {
		GameBetList = append(GameBetList, &GameBet{
			GameID:          game.ID,
			Bet:             game.BetList,
			GameName:        game.Name,
			DefaultBet:      game.DefaultBet,
			DefaultBetLevel: game.DefaultBetLevel,
		})
		GameNowList = append(GameNowList, &GameNow{
			GameID:        game.ID,
			RewardPercent: game.RewardPercent,
		})
	}

	err = loadGameBet(GameBetList)
	if err != nil {
		logger.Fatalf("加载游戏配置失败：存储GameBet(%s)", err.Error())

	}

	err = loadGameNow(GameNowList)
	if err != nil {
		logger.Fatalf("加载游戏配置失败：存储GameNow(%s)", err.Error())

	}

}

func BalanceAlert(timeInterval int) {
	slog.Info("--------------------BalanceAlert start--------------------", "timeInterval", timeInterval)
	oper := []*comm.Operator_V2{}
	CollAdminOperator.FindAll(bson.M{"BalanceAlertTimeInterval": timeInterval}, &oper)
	for _, op := range oper {
		if op.Balance < op.BalanceAlert {
			ret := &comm.BalanceAlertPs{
				AppID:           op.AppID,
				CooperationType: op.CooperationType,
				WalletType:      op.WalletMode,
				MerchantBalance: op.Balance,
			}
			//发送邮件
			_ = mq.PublishMsg("/alerter/OperatorBalanceAlert", ret)
		}
	}
	slog.Info("--------------------BalanceAlert end--------------------", "timeInterval", timeInterval)
}
