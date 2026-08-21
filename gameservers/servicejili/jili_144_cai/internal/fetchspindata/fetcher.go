package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"serve/comm/db"
	"serve/comm/ut"
	"serve/servicejili/jili_144_cai/internal"
	"serve/servicejili/jili_144_cai/internal/gendata"
	"serve/servicejili/jili_144_cai/internal/message"
	"serve/servicejili/jili_144_cai/internal/models"
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
	gameReq      *message.Cai_GameReqData
	historyToken string
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

	launchUrl := lo.Must(jiliut.LaunchGame(username, fmt.Sprintf("%d", internal.GameNo), "en-US"))

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
		Path:   fmt.Sprintf("/%s/account/login", internal.GameShortName),
	}

	var loginDataReq = message.Cai_LoginDataReq{
		Token:          proto.String(token),
		OSType:         proto.String("Windows"),
		OSVersion:      proto.String(""),
		Browser:        proto.String("chrome"),
		Language:       proto.String("en"),
		BrowserVersion: proto.String("126.0.0.0"),
		Width:          proto.Int32(2560),
		Height:         proto.Int32(1440),
		Ratio:          proto.Float64(1.5),
		Machine:        proto.String(""),
		BrowserTag:     proto.Int32(31),
	}

	// http.Post(gameloginUrl.String(), "application/x-protobuf", jiliut.ProtoReader(&loginDataReq))

	var gameloginResData message.Cai_ResData
	lo.Must0(jiliut.PostProto(gameloginUrl.String(), &loginDataReq, &gameloginResData))

	// https://uat-wbslot-fd.jlfafafa1.com/csh/game/spin?token=
	f.spinUrl = &url.URL{
		Scheme: "https",
		Host:   ut.ReverseString(query.Get("gs")),
		Path:   fmt.Sprintf("/%s/game/spin", internal.GameShortName),
		RawQuery: url.Values{
			"token": {""},
		}.Encode(),
	}
	extra := int32(0)
	bet := float64(1)

	if ty == internal.GameTypeExtra {
		extra = 1
		bet = float64(2)
	} else if ty == internal.GameTypeExtraPlus {
		extra = 2
		bet = float64(10)
	}
	f.gameReq = &message.Cai_GameReqData{
		Token: proto.String(token),
		Aid:   proto.Int32(tokenST.Profile.Aid),
		Apiid: proto.Int32(tokenST.Profile.ApiId),
		// Encode: lo.Must(proto.Marshal(&spinReq)),
		Encode: jiliut.ProtoEncode(&message.Cai_SpinReq{
			Bet:   proto.Float64(bet),
			Extra: proto.Int32(extra),
		}),
	}
	return f
}

func (f *Fetcher) run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	for {
		time.Sleep(time.Second)
		var spinResData message.Cai_ResData
		lo.Must0(jiliut.PostProto(f.spinUrl.String(), f.gameReq, &spinResData))

		// spinResData.Data
		if len(spinResData.Data) == 0 {
			break
		}
		infoData := spinResData.Data[0]

		var spinAll message.Cai_SpinAllData
		lo.Must0(proto.Unmarshal(infoData.Encode, &spinAll))

		// spinResData.
		// fmt.Println(&spinResData)
		lo.Must0(len(spinAll.Info) == 1)
		spinAck := spinAll.Info[0]
		if spinAck.GetAckType() != 0 {
			slog.Error("pull spin data error", "acktype", spinAck.GetAckType(), "user", f.username, "type", f.ty)
			break
		}

		totalWin := spinAck.GetTotalWin()
		slog.Info("", "totalWin", totalWin)
		baseBet := float64(1)
		doc := &models.RawSpin{
			ID:       primitive.NewObjectID(),
			Times:    totalWin / baseBet,
			Type:     f.ty,
			HasGame:  len(spinAck.RoundQueue) > 1,
			Data:     spinAck,
			Selected: true,
		}
		doc.BucketId = gendata.GBuckets.GetBucket(doc.Times, doc.HasGame, db.BoundType(doc.Type))
		f.fetchHistory(*spinAck.ShowIndex, doc)

		f.coll.InsertOne(context.TODO(), doc)
		fmt.Println("插入成功。。。", f.ty)
	}
}
