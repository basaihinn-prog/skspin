package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"serve/servicejili/jili_110_ali/internal"
	"serve/servicejili/jili_110_ali/internal/gendata"
	"serve/servicejili/jili_110_ali/internal/message"
	"serve/servicejili/jili_110_ali/internal/models"
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
	ty       int

	launchUrl    *url.URL
	spinUrl      *url.URL
	gameReq      *jiliUtMessage.Server_Request
	historyToken string
	token        string
}

func newFetcher(username string, coll *mongo.Collection, ty int) *Fetcher {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	f := &Fetcher{
		username: username,
		coll:     coll,
		ty:       ty,
	}

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
			Bet:     proto.Float64(internal.BaseBet),
			Extra:   lo.Ternary(ty == internal.GameTypeExtra, proto.Uint32(1), nil),
			MallBet: lo.Ternary(ty == internal.GameTypeGame, proto.Float64(48), nil),
			//todo mallID
			//MallID:  lo.Ternary(ty == gendata.GameTypeGame, proto.Int32(50), nil),
			// Cheat:          proto.Uint32(0),
			// CurrencyNumber: proto.Int32(0),
			// Vip:            proto.Uint32(0),
			// IsFree:         proto.Bool(false),
			// Isturbo:        proto.Bool(false),
			// Istournament:   proto.Bool(false),
		}),
	}

	return f
}

func (f *Fetcher) run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			debug.PrintStack()
		}
	}()
	for ; ; time.Sleep(time.Second) {
		fmt.Println("拉取数据。。。")
		var response jiliUtMessage.Server_GaiaResponse
		lo.Must0(jiliut.PostProtoWithHeaders(f.spinUrl.String(), f.gameReq, &response, map[string]string{"token": f.token}))

		// spinResData.Data
		if len(response.Data) == 0 {
			break
		}

		var spinResponse jiliUtMessage.Server_SpinResponse
		lo.Must0(proto.Unmarshal(response.Data, &spinResponse))

		var aliData message.Ali_SpinAck
		lo.Must0(proto.Unmarshal(spinResponse.Data, &aliData))
		spinResponse.Data = []byte{}

		if aliData.GetAckType() != 0 {
			slog.Error("pull spin data error", "acktype", aliData.GetAckType(), "user", f.username, "type:", f.ty)
			break
		}

		totalWin := aliData.GetTotalWin()
		slog.Info("", "totalWin", totalWin)

		doc := &models.RawSpin{
			ID:    primitive.NewObjectID(),
			Times: totalWin / internal.BaseBet,
			//Type:     lo.Ternary(f.buy, gendata.GameTypeGame, gendata.GameTypeNormal),
			Type:     f.ty,
			HasGame:  len(aliData.GetRoundQueue()) > 1,
			ZeusData: &aliData,
			UtData:   &spinResponse,
			Selected: true,
		}
		doc.BucketId = gendata.GBuckets.GetBucket(doc.Times, doc.HasGame, db.BoundType(doc.Type))
		f.fetchHisotry(strconv.Itoa(int(spinResponse.GetRoundIndexV2())), doc)

		f.coll.InsertOne(context.TODO(), doc)
		fmt.Println("插入成功。。。")
	}
}
