package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"serve/servicejili/jili_17_sweetheart/internal"
	"serve/servicejili/jili_17_sweetheart/internal/gendata"
	tmp "serve/servicejili/jili_17_sweetheart/internal/message"
	"serve/servicejili/jili_17_sweetheart/internal/models"

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
	gameReq      *tmp.Server_Request
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

	launchUrl := lo.Must(jiliut.LaunchGame(username, "17", "en-US"))

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

	var loginDataReq = tmp.Server_InfoReq{
		Os:       proto.String("Windows"),
		Language: proto.String("en-US"),
		Browser: []*tmp.Server_Browser{
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

	// http.Post(gameloginUrl.String(), "application/x-protobuf", jiliut.ProtoReader(&loginDataReq))
	var gameloginResData tmp.Server_InfoReq
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

	f.gameReq = &tmp.Server_Request{
		Ack: proto.Int32(0),
		Req: jiliut.ProtoEncode(&tmp.Server_SpinReq{
			Bet:            proto.Int64(1000),
			CurrencyNumber: proto.Int32(0),
			MallBet:        lo.Ternary(buy, proto.Int64(29500), nil),
		}),
	}

	return f
}

func (f *Fetcher) run() bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			debug.PrintStack()
		}
	}()
	for ; ; time.Sleep(time.Second) {
		fmt.Println("拉取数据。。。")
		var response tmp.Server_Response
		lo.Must0(jiliut.PostProtoWithHeaders(f.spinUrl.String(), f.gameReq, &response, map[string]string{"token": f.token}))

		// spinResData.Data
		if len(response.Response) == 0 {
			break
		}

		var spinResponse tmp.Server_SpinResponse
		lo.Must0(proto.Unmarshal(response.Response, &spinResponse))

		var sweetheartData tmp.Custom_AllPlate
		lo.Must0(proto.Unmarshal(spinResponse.AllPlate, &sweetheartData))
		spinResponse.AllPlate = []byte{}

		if len(sweetheartData.Plate) <= 0 {
			slog.Error("pull spin data error", " plate length is 0", " err")
			break
		}

		var totalWin int64

		for _, plate := range sweetheartData.Plate {
			totalWin += plate.GetWin()
		}
		totalWin = totalWin / 1000
		slog.Info("", "totalWin", totalWin)

		doc := &models.RawSpin{
			ID:             primitive.NewObjectID(),
			Times:          float64(totalWin), // bet is 1.0 for csh
			Type:           lo.Ternary(f.buy, internal.GameTypeGame, internal.GameTypeNormal),
			HasGame:        len(sweetheartData.GetPlate()) > 1,
			SweetheartData: &sweetheartData,
			UtData:         &spinResponse,
			Selected:       true,
		}
		doc.BucketId = gendata.GBuckets.GetBucket(doc.Times, doc.HasGame, db.BoundType(doc.Type))
		f.fetchHisotry(strconv.Itoa(int(spinResponse.GetRoundIndex())), doc)

		f.coll.InsertOne(context.TODO(), doc)
		fmt.Println("插入成功。。。")

		ty := lo.Ternary(f.buy, internal.GameTypeGame, internal.GameTypeNormal)
		if rand.Int()%100 > 70 {
			cnt, _ := f.coll.CountDocuments(context.TODO(), db.D("type", ty))
			if ty == internal.GameTypeNormal {
				if cnt >= 320000 {
					wg.Done()
					return true
				}
			} else if ty == internal.GameTypeGame {
				if cnt >= 60000 {
					wg.Done()
					return true
				}
			}
		}
	}
	return false
}
