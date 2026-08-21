package main

import (
	"fmt"
	"time"

	"serve/servicejili/jili_40_ols3/internal"
	"serve/servicejili/jiliut"

	"serve/comm/db"
)

func main() {
	mongoaddr := jiliut.GetFetchMongoAddr()
	db.DialToMongo(mongoaddr, internal.GameID)
	coll := db.Collection("rawSpinData")

	for i := 0; i < 5; i++ {
		go func(id int) {
			var num int
			for {
				newFetcher(fmt.Sprintf("%d_40123456_%d_%d", internal.GameNo, id, num), coll).run()
				fmt.Println("normal 拉取数据失败了，重试中。。。")
				time.Sleep(5 * time.Second)
				num++
			}
		}(i)
	}
	select {}
}

// func fetch(username string, coll *mongo.Collection) {

// 	// var jili = jili.JILI()
// 	launchUrl := lo.Must(jiliut.LaunchGame(username, "45", "en-US"))

// 	ul := lo.Must(url.Parse(launchUrl))

// 	query := ul.Query()
// 	// fmt.Println(query)

// 	be := query.Get("be")

// 	loginUrl := url.URL{
// 		Scheme: "https",
// 		Host:   ut.ReverseString(be),
// 		Path:   "/sso-login.api",
// 		RawQuery: url.Values{
// 			"key":  {query.Get("ssoKey")},
// 			"lang": {"en-US"},
// 		}.Encode(),
// 	}

// 	// https://uat-wbgame.jlfafafa3.com/csh/?ssoKey=e062d6e682e4c67294de904238b60ff607ce433f&lang=en-US&gameID=2&gs=moc.1afafaflj.df-tolsbw-tau&domain_gs=1afafaflj&domain_platform=moc.1afafaflj.df-tolsbw-tau&be=moc.2afafaflj.ipabewbw-tau&apiId=604&iu=true&legalLang=true

// 	// https://uat-wbwebapi.jlfafafa2.com/sso-login.api?key=a5d5741c927fd5c6b6e9c4fd15116adaf5713bed&lang=en-US

// 	var tokenST struct {
// 		Token   string
// 		Profile struct {
// 			Aid   int32
// 			ApiId int32
// 		}
// 	}
// 	jiliut.PostJson(loginUrl.String(), nil, &tokenST)

// 	// lo.Must0(json.Unmarshal(respBody, &tokenST))

// 	token := tokenST.Token
// 	// fmt.Println(token)

// 	// https://uat-wbslot-fd.jlfafafa1.com/csh/account/login?
// 	gameloginUrl := url.URL{
// 		Scheme: "https",
// 		Host:   ut.ReverseString(query.Get("gs")),
// 		Path:   fmt.Sprintf("/%s/account/login", internal.GameShortName),
// 	}

// 	var loginDataReq = message.Ols3_LoginDataReq{
// 		Token:          proto.String(token),
// 		OSType:         proto.String("Windows"),
// 		OSVersion:      proto.String(""),
// 		Browser:        proto.String("chrome"),
// 		Language:       proto.String("en"),
// 		BrowserVersion: proto.String("125.0.0.0"),
// 		Width:          proto.Int32(2560),
// 		Height:         proto.Int32(1440),
// 		Ratio:          proto.Float64(1.5),
// 		Machine:        proto.String(""),
// 		BrowserTag:     proto.Int32(31),
// 	}

// 	// http.Post(gameloginUrl.String(), "application/x-protobuf", jiliut.ProtoReader(&loginDataReq))

// 	var gameloginResData message.Ols3_ResData
// 	lo.Must0(jiliut.PostProto(gameloginUrl.String(), &loginDataReq, &gameloginResData))

// 	// https://uat-wbslot-fd.jlfafafa1.com/csh/game/spin?token=
// 	spinUrl := url.URL{
// 		Scheme: "https",
// 		Host:   ut.ReverseString(query.Get("gs")),
// 		Path:   fmt.Sprintf("/%s/game/spin", internal.GameShortName),
// 		RawQuery: url.Values{
// 			"token": {""},
// 		}.Encode(),
// 	}

// 	// var spinReq = message.Csh_SpinReq{
// 	// 	Bet: proto.Float64(1),
// 	// }
// 	var gameReq = message.Ols3_GameReqData{
// 		Token: proto.String(token),
// 		Aid:   proto.Int32(tokenST.Profile.Aid),
// 		Apiid: proto.Int32(tokenST.Profile.ApiId),
// 		// Encode: lo.Must(proto.Marshal(&spinReq)),
// 		Encode: jiliut.ProtoEncode(&message.Ols3_SpinReq{
// 			Bet: proto.Float64(internal.BaseBet),
// 		}),
// 	}

// 	for ; ; time.Sleep(time.Second) {
// 		var spinResData message.Ols3_ResData
// 		lo.Must0(jiliut.PostProto(spinUrl.String(), &gameReq, &spinResData))

// 		// spinResData.Data
// 		if len(spinResData.Data) <= 0 {
// 			fmt.Println("获取数据为空")
// 			continue
// 		}
// 		infoData := spinResData.Data[0]

// 		var spinAll message.Ols3_SpinAllData
// 		lo.Must0(proto.Unmarshal(infoData.Encode, &spinAll))

// 		// spinResData.
// 		fmt.Println(&spinResData)
// 		lo.Must0(len(spinAll.Info) == 1)
// 		spinAck := spinAll.Info[0]
// 		totalWin := spinAck.GetTotalWin()
// 		slog.Info("", "totalWin", totalWin)
// 		hasGame := false

// 		doc := &models.RawSpin{
// 			ID:       primitive.NewObjectID(),
// 			Times:    totalWin / float64(internal.BaseBet),
// 			Type:     gendata.GameTypeNormal,
// 			HasGame:  hasGame,
// 			Data:     spinAck,
// 			Selected: true,
// 		}
// 		coll.InsertOne(context.TODO(), doc)
// 	}
// }
