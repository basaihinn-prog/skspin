package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"serve/servicejili/jili_49_fullhouse/internal"
	"serve/servicejili/jili_49_fullhouse/internal/gendata"
	"serve/servicejili/jili_49_fullhouse/internal/message"
	"serve/servicejili/jili_49_fullhouse/internal/models"
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
			Bet:           proto.Float64(1),
			MallBet:       lo.Ternary(buy, proto.Float64(40.5), nil),
		}),
	}

	fmt.Println(username, " new inoke 4")
	return f
}

func (f *Fetcher) run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	for ; ; time.Sleep(time.Second) {
		fmt.Println(f.username, " run inoke 1")
		var response jiliUtMessage.Server_GaiaResponse
		lo.Must0(jiliut.PostProtoWithHeaders(f.spinUrl.String(), f.gameReq, &response, map[string]string{"token": f.token}))
		fmt.Println(f.username, " run inoke 2")

		// spinResData.Data
		if len(response.Data) == 0 {
			break
		}

		var spinResponse jiliUtMessage.Server_SpinResponse
		lo.Must0(proto.Unmarshal(response.Data, &spinResponse))

		fmt.Println(f.username, " run inoke 3")

		var zeusData message.Custom_AllPlate
		lo.Must0(proto.Unmarshal(spinResponse.Data, &zeusData))
		spinResponse.Data = []byte{}

		fmt.Println(f.username, " run inoke 4")

		if len(zeusData.Plate) <= 0 {
			slog.Error("pull spin data error", " plate length is 0", " err")
			break
		}

		fmt.Println(f.username, " run inoke 5")

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

		fmt.Println(f.username, " run inoke 6")

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
		fmt.Println(f.username, " run inoke 7")
		f.fetchHisotry(strconv.Itoa(int(spinResponse.GetRoundIndexV2())), doc)
		fmt.Println(f.username, " run inoke 8")

		f.coll.InsertOne(context.TODO(), doc)
		fmt.Println(f.username, " run inoke 9")
		fmt.Println("插入成功。。。")
	}
}
