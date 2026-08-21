package main

import (
	"fmt"
	"time"

	"serve/comm/db"
	"serve/servicejili/jili_137_ge/internal"
)

func main() {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	for i := 0; i < 5; i++ {
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("%d_123456_%d", internal.GameNo, id), coll, internal.GameTypeNormal).run()
				fmt.Println("拉取Normal数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)

		// 小游戏数据拉够了
		// go func(id int) {
		// 	var num int
		// 	for {
		// 		newFetcher(fmt.Sprintf("%d_buyccc_123456_%d", internal.GameNo, id), coll, gendata.GameTypeBuy).run()
		// 		fmt.Println("拉取Buy数据失败了，重试中。。。")
		// 		time.Sleep(5 * time.Second)
		// 		num++
		// 	}
		// }(i)

		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("%d_extra_123456_%d", internal.GameNo, id), coll, internal.GameTypeExtra).run()
				fmt.Println("拉取Extra数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
	}
	select {}
}

//
//func fetch(username string, coll *mongo.Collection) {
//
//	// var jili = jili.JILI()
//	launchUrl := lo.Must(jiliut.LaunchGame(username, fmt.Sprintf("%d", internal.GameNo), "en-US"))
//
//	ul := lo.Must(url.Parse(launchUrl))
//
//	query := ul.Query()
//	// fmt.Println(query)
//
//	be := query.Get("be")
//
//	loginUrl := url.URL{
//		Scheme: "https",
//		Host:   ut.ReverseString(be),
//		Path:   "/sso-login.api",
//		RawQuery: url.Values{
//			"key":  {query.Get("ssoKey")},
//			"lang": {"en-US"},
//		}.Encode(),
//	}
//
//	// https://uat-wbgame.jlfafafa3.com/csh/?ssoKey=e062d6e682e4c67294de904238b60ff607ce433f&lang=en-US&gameID=2&gs=moc.1afafaflj.df-tolsbw-tau&domain_gs=1afafaflj&domain_platform=moc.1afafaflj.df-tolsbw-tau&be=moc.2afafaflj.ipabewbw-tau&apiId=604&iu=true&legalLang=true
//
//	// https://uat-wbwebapi.jlfafafa2.com/sso-login.api?key=a5d5741c927fd5c6b6e9c4fd15116adaf5713bed&lang=en-US
//
//	var tokenST struct {
//		Token   string
//		Profile struct {
//			Aid   int32
//			ApiId int32
//		}
//	}
//	jiliut.PostJson(loginUrl.String(), nil, &tokenST)
//
//	// lo.Must0(json.Unmarshal(respBody, &tokenST))
//
//	token := tokenST.Token
//	// fmt.Println(token)
//
//	// https://uat-wbslot-fd.jlfafafa1.com/csh/account/login?
//	gameloginUrl := url.URL{
//		Scheme: "https",
//		Host:   ut.ReverseString(query.Get("gs")),
//		Path:   fmt.Sprintf("/%s/account/login", internal.GameShortName),
//	}
//
//	var loginDataReq = message.Rs2_LoginDataReq{
//		Token:          proto.String(token),
//		OSType:         proto.String("Windows"),
//		OSVersion:      proto.String(""),
//		Browser:        proto.String("chrome"),
//		Language:       proto.String("en"),
//		BrowserVersion: proto.String("125.0.0.0"),
//		Width:          proto.Int32(2560),
//		Height:         proto.Int32(1440),
//		Ratio:          proto.Float64(1.5),
//		Machine:        proto.String(""),
//		BrowserTag:     proto.Int32(31),
//	}
//
//	// http.Post(gameloginUrl.String(), "application/x-protobuf", jiliut.ProtoReader(&loginDataReq))
//
//	var gameloginResData message.Rs2_ResData
//	lo.Must0(jiliut.PostProto(gameloginUrl.String(), &loginDataReq, &gameloginResData))
//
//	// https://uat-wbslot-fd.jlfafafa1.com/csh/game/spin?token=
//	spinUrl := url.URL{
//		Scheme: "https",
//		Host:   ut.ReverseString(query.Get("gs")),
//		Path:   fmt.Sprintf("/%s/game/spin", internal.GameShortName),
//		RawQuery: url.Values{
//			"token": {""},
//		}.Encode(),
//	}
//
//	// var spinReq = message.Csh_SpinReq{
//	// 	Bet: proto.Float64(1),
//	// }
//	var gameReq = message.Rs2_GameReqData{
//		Token: proto.String(token),
//		Aid:   proto.Int32(tokenST.Profile.Aid),
//		Apiid: proto.Int32(tokenST.Profile.ApiId),
//		// Encode: lo.Must(proto.Marshal(&spinReq)),
//		Encode: jiliut.ProtoEncode(&message.Rs2_SpinReq{
//			Bet: proto.Float64(internal.BaseBet),
//		}),
//	}
//
//	for ; ; time.Sleep(time.Second) {
//		var spinResData message.Rs2_ResData
//		lo.Must0(jiliut.PostProto(spinUrl.String(), &gameReq, &spinResData))
//
//		if len(spinResData.Data) <= 0 {
//			fmt.Println("拉取数据错误")
//			continue
//		}
//		// spinResData.Data
//		infoData := spinResData.Data[0]
//
//		var spinAll message.Rs2_SpinAllData
//		lo.Must0(proto.Unmarshal(infoData.Encode, &spinAll))
//		// spinResData.
//		fmt.Println(&spinResData)
//		lo.Must0(len(spinAll.Info) == 1)
//		spinAck := spinAll.Info[0]
//		if spinAck.AckType != nil {
//			fmt.Println("拉取数据错误")
//			continue
//		}
//		totalWin := spinAck.GetTotalWin()
//		slog.Info("", "totalWin", totalWin)
//
//		hasGame := false
//		for _, round := range spinAck.RoundQueue {
//			if round.GetFreeRemainRound() > 0 {
//				hasGame = true
//				break
//			}
//		}
//		doc := &models.RawSpin{
//			ID:       primitive.NewObjectID(),
//			Times:    totalWin / float64(internal.BaseBet),
//			Type:     gendata.GameTypeNormal,
//			HasGame:  hasGame,
//			Data:     spinAck,
//			Selected: true,
//		}
//		coll.InsertOne(context.TODO(), doc)
//	}
//}
//
//func fetchBuy(username string, coll *mongo.Collection) {
//
//	// var jili = jili.JILI()
//	launchUrl := lo.Must(jiliut.LaunchGame(username, "102", "en-US"))
//
//	ul := lo.Must(url.Parse(launchUrl))
//
//	query := ul.Query()
//	// fmt.Println(query)
//
//	be := query.Get("be")
//
//	loginUrl := url.URL{
//		Scheme: "https",
//		Host:   ut.ReverseString(be),
//		Path:   "/sso-login.api",
//		RawQuery: url.Values{
//			"key":  {query.Get("ssoKey")},
//			"lang": {"en-US"},
//		}.Encode(),
//	}
//
//	// https://uat-wbgame.jlfafafa3.com/csh/?ssoKey=e062d6e682e4c67294de904238b60ff607ce433f&lang=en-US&gameID=2&gs=moc.1afafaflj.df-tolsbw-tau&domain_gs=1afafaflj&domain_platform=moc.1afafaflj.df-tolsbw-tau&be=moc.2afafaflj.ipabewbw-tau&apiId=604&iu=true&legalLang=true
//
//	// https://uat-wbwebapi.jlfafafa2.com/sso-login.api?key=a5d5741c927fd5c6b6e9c4fd15116adaf5713bed&lang=en-US
//
//	var tokenST struct {
//		Token   string
//		Profile struct {
//			Aid   int32
//			ApiId int32
//		}
//	}
//	jiliut.PostJson(loginUrl.String(), nil, &tokenST)
//
//	// lo.Must0(json.Unmarshal(respBody, &tokenST))
//
//	token := tokenST.Token
//	// fmt.Println(token)
//
//	// https://uat-wbslot-fd.jlfafafa1.com/csh/account/login?
//	gameloginUrl := url.URL{
//		Scheme: "https",
//		Host:   ut.ReverseString(query.Get("gs")),
//		Path:   fmt.Sprintf("/%s/account/login", internal.GameShortName),
//	}
//
//	var loginDataReq = message.Rs2_LoginDataReq{
//		Token:          proto.String(token),
//		OSType:         proto.String("Windows"),
//		OSVersion:      proto.String(""),
//		Browser:        proto.String("chrome"),
//		Language:       proto.String("en"),
//		BrowserVersion: proto.String("125.0.0.0"),
//		Width:          proto.Int32(2560),
//		Height:         proto.Int32(1440),
//		Ratio:          proto.Float64(1.5),
//		Machine:        proto.String(""),
//		BrowserTag:     proto.Int32(31),
//	}
//
//	// http.Post(gameloginUrl.String(), "application/x-protobuf", jiliut.ProtoReader(&loginDataReq))
//
//	var gameloginResData message.Rs2_ResData
//	lo.Must0(jiliut.PostProto(gameloginUrl.String(), &loginDataReq, &gameloginResData))
//
//	// https://uat-wbslot-fd.jlfafafa1.com/csh/game/spin?token=
//	spinUrl := url.URL{
//		Scheme: "https",
//		Host:   ut.ReverseString(query.Get("gs")),
//		Path:   fmt.Sprintf("/%s/game/spin", internal.GameShortName),
//		RawQuery: url.Values{
//			"token": {""},
//		}.Encode(),
//	}
//
//	// var spinReq = message.Csh_SpinReq{
//	// 	Bet: proto.Float64(1),
//	// }
//	var gameReq = message.Rs2_GameReqData{
//		Token: proto.String(token),
//		Aid:   proto.Int32(tokenST.Profile.Aid),
//		Apiid: proto.Int32(tokenST.Profile.ApiId),
//		// Encode: lo.Must(proto.Marshal(&spinReq)),
//		Encode: jiliut.ProtoEncode(&message.Rs2_SpinReq{
//			Bet:     proto.Float64(internal.BaseBet),
//			MallBet: proto.Float64(internal.BaseBet * internal.BuyMul),
//			MallID:  proto.Int32(50),
//		}),
//	}
//
//	for ; ; time.Sleep(time.Second) {
//		var spinResData message.Rs2_ResData
//		lo.Must0(jiliut.PostProto(spinUrl.String(), &gameReq, &spinResData))
//
//		if len(spinResData.Data) <= 0 {
//			continue
//		}
//		// spinResData.Data
//		infoData := spinResData.Data[0]
//
//		var spinAll message.Rs2_SpinAllData
//		lo.Must0(proto.Unmarshal(infoData.Encode, &spinAll))
//
//		// spinResData.
//		fmt.Println(&spinResData)
//		lo.Must0(len(spinAll.Info) == 1)
//		spinAck := spinAll.Info[0]
//		totalWin := spinAck.GetTotalWin()
//		slog.Info("", "totalWin", totalWin)
//
//		doc := &models.RawSpin{
//			ID:       primitive.NewObjectID(),
//			Times:    totalWin / float64(internal.BaseBet),
//			HasGame:  true,
//			Type:     gendata.GameTypeGame,
//			Data:     spinAck,
//			Selected: true,
//		}
//		coll.InsertOne(context.TODO(), doc)
//	}
//}
