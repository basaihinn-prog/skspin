package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"serve/servicejili/jili_76_nightclub/internal"
	"serve/servicejili/jili_76_nightclub/internal/gendata"
	"serve/servicejili/jili_76_nightclub/internal/message"
	"serve/servicejili/jili_76_nightclub/internal/models"
	"serve/servicejili/jiliut/jiliUtMessage"

	"serve/comm/db"
	"serve/comm/ut"

	"serve/servicejili/jiliut"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/proto"
)

type Fetcher struct {
	username string
	coll     *mongo.Collection
	buy      bool
	ty       int

	launchUrl    *url.URL
	spinUrl      *url.URL
	gameReq      *jiliUtMessage.Server_Request
	historyToken string
	token        string
}

func newFetcher(username string, coll *mongo.Collection, buy bool) *Fetcher {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	f := &Fetcher{
		username: username,
		coll:     coll,
		buy:      buy,
	}
	f.ty = lo.Ternary(buy, gendata.GameTypeGame, gendata.GameTypeNormal)
	fmt.Println(username, " new inoke 1")

	launchUrl := lo.Must(jiliut.LaunchGame(username, fmt.Sprintf("%v", internal.GameNo), "en-US"))

	ul := lo.Must(url.Parse(launchUrl))
	f.launchUrl = ul

	query := ul.Query()
	// f.launchQuery = query
	// fmt.Println(query)

	be := query.Get("be")

	loginUrl := url.URL{
		Scheme: "https",
		Host:   ut.ReverseString(be),
		Path:   "/sso-login.api",
		RawQuery: url.Values{
			"key":  {query.Get("ssoKey")},
			"lang": {"en-US"},
		}.Encode(),
	}
	var tokenST struct {
		Token   string
		Profile struct {
			Aid   int32
			ApiId int32
		}
	}
	jiliut.PostJson(loginUrl.String(), nil, &tokenST)

	fmt.Println(username, " new inoke 2")

	token := tokenST.Token
	f.token = token

	historyTokenUrl := url.URL{
		Scheme: "https",
		Host:   "history-api" + ul.Host[strings.IndexByte(ul.Host, '.'):],
		Path:   "/token",
	}

	var historyTokenRet struct {
		Code    int
		Message string
		Data    struct {
			Token string
		}
	}
	jiliut.PostJson(historyTokenUrl.String(), map[string]string{
		"Token": token,
	}, &historyTokenRet)

	lo.Must0(historyTokenRet.Data.Token != "")
	f.historyToken = historyTokenRet.Data.Token

	gameloginUrl := url.URL{
		Scheme: "https",
		Host:   ut.ReverseString(query.Get("gs")),
		// Path:   fmt.Sprintf("/%s/req?D=1", internal.GameShortName),
		Path:     fmt.Sprintf("/%s/req", internal.GameShortName),
		RawQuery: "D=1&",
	}

	var loginDataReq = jiliUtMessage.Server_InfoReq{
		Os:       proto.String("Windows"),
		Language: proto.String("en-US"),
		Browser: []*jiliUtMessage.Server_Browser{
			{
				Type:     proto.String("chrome"),
				Version:  proto.String("126.0.0.0"),
				Language: proto.String("en-US"),
				Width:    proto.Uint32(1920),
				Height:   proto.Uint32(1080),
				Ratio:    nil,
			},
		},
		Version: proto.String(""),
		Model:   proto.String(""),
	}

	fmt.Println(username, " new inoke 3")

	// http.Post(gameloginUrl.String(), "application/x-protobuf", jiliut.ProtoReader(&loginDataReq))
	var gameloginResData jiliUtMessage.Server_GameInfoAck
	err := jiliut.PostProtoWithHeaders(gameloginUrl.String(), &loginDataReq, &gameloginResData, map[string]string{"token": token})
	if err != nil {
		panic(err)
	}
	// https://uat-wbslot-fd.jlfafafa1.com/csh/game/spin?token=
	f.spinUrl = &url.URL{
		Scheme: "https",
		Host:   ut.ReverseString(query.Get("gs")),
		// Path:   fmt.Sprintf("/%s/req?D=0", internal.GameShortName),
		Path:     fmt.Sprintf("/%s/req", internal.GameShortName),
		RawQuery: "D=0&",
	}

	f.gameReq = &jiliUtMessage.Server_Request{
		Ack: proto.Int32(0),
		Req: jiliut.ProtoEncode(&jiliUtMessage.Server_SpinReq{
			BuyCustomType: proto.Int32(0),
			Bet:           proto.Float64(internal.BaseBet),
			MallBet:       lo.Ternary(buy, proto.Float64(39), nil),
		}),
	}

	fmt.Println(username, " new inoke 4")
	return f
}

func (f *Fetcher) run() bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	for ; ; time.Sleep(time.Second) {
		var response jiliUtMessage.Server_GaiaResponse
		lo.Must0(jiliut.PostProtoWithHeaders(f.spinUrl.String(), f.gameReq, &response, map[string]string{"token": f.token}))

		// spinResData.Data
		if len(response.Data) == 0 {
			break
		}

		var spinResponse jiliUtMessage.Server_SpinResponse
		lo.Must0(proto.Unmarshal(response.Data, &spinResponse))

		var zeusData message.Custom_AllPlate
		lo.Must0(proto.Unmarshal(spinResponse.Data, &zeusData))
		spinResponse.Data = []byte{}

		if len(zeusData.Plate) <= 0 {
			slog.Error("pull spin data error", " plate length is 0", " err")
			break
		}

		//todo
		// if zeusData.GetAckType() != 0 {
		// 	slog.Error("pull spin data error", "acktype", zeusData.GetAckType(), "user", f.username, "buy", f.buy)
		// 	break
		// }

		//todo
		// totalWin := zeusData.GetWin()
		var totalWin float64
		for _, plate := range zeusData.Plate {
			totalWin += plate.GetWin()
		}
		slog.Info(f.username, " totalWin", totalWin)

		doc := &models.RawSpin{
			ID:    primitive.NewObjectID(),
			Times: totalWin / 1.0, // bet is 1.0 for csh
			Type:  lo.Ternary(f.buy, gendata.GameTypeGame, gendata.GameTypeNormal),
			//HasGame:  len(zeusData.GetRoundQueue()) > 1,
			HasGame:  len(zeusData.Plate) > 1,
			ZeusData: &zeusData,
			UtData:   &spinResponse,
			Selected: true,
		}
		doc.BucketId = gendata.GBuckets.GetBucket(doc.Times, doc.HasGame, db.BoundType(doc.Type))
		f.fetchHisotry(strconv.Itoa(int(spinResponse.GetRoundIndexV2())), doc)
		f.coll.InsertOne(context.TODO(), doc)
		fmt.Println("插入成功。。。")

		if rand.Int()%100 > 70 {
			cnt, _ := f.coll.CountDocuments(context.TODO(), db.D("type", f.ty))
			if f.ty == gendata.GameTypeNormal {
				if cnt >= 320000 {
					wg.Done()
					return true
				}
			} else if f.ty == gendata.GameTypeGame {
				if cnt >= 50000 {
					wg.Done()
					return true
				}
			}
		}
	}
	return false
}
