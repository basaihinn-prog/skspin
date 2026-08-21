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
	"serve/servicejili/jili_67_mp/internal"
	"serve/servicejili/jili_67_mp/internal/gendata"
	"serve/servicejili/jili_67_mp/internal/models"
	"serve/servicejili/jiliut"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Fetcher struct {
	username string
	coll     *mongo.Collection
	buy      bool

	launchUrl    *url.URL
	spinUrl      *url.URL
	gameReq      map[string]interface{}
	historyToken string
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

	launchUrl := lo.Must(jiliut.LaunchGame(username, "67", "en-US"))

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

	var loginDataReq = map[string]interface{}{
		"BrowserTag":     31,
		"BrowserVersion": "126.0.0.0",
		"Height":         1080,
		"OSVersion":      "",
		"Ratio":          1,
		"Width":          1920,
		"browser":        "chrome",
		"language":       "en",
		"machine":        "",
		"os":             "Windows",
		"source":         0,
		"token":          token,
	}

	// http.Post(gameloginUrl.String(), "application/x-protobuf", jiliut.ProtoReader(&loginDataReq))

	var gameloginResData map[string]interface{}
	lo.Must0(jiliut.PostJson(gameloginUrl.String(), &loginDataReq, &gameloginResData))

	// https://uat-wbslot-fd.jlfafafa1.com/csh/game/spin?token=
	f.spinUrl = &url.URL{
		Scheme: "https",
		Host:   ut.ReverseString(query.Get("gs")),
		Path:   fmt.Sprintf("/%s/game/spin", internal.GameShortName),
		RawQuery: url.Values{
			"token": {token},
			"aid":   {fmt.Sprintf("%d", tokenST.Profile.Aid)},
			"apild": {fmt.Sprintf("%d", tokenST.Profile.ApiId)},
			"delay": {"0"},
		}.Encode(),
	}
	f.gameReq = map[string]interface{}{
		"bet":      1,
		"currency": 0,
		"free":     0,
	}
	return f
}

func (f *Fetcher) run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	for ; ; time.Sleep(time.Second) {
		fmt.Println("拉取数据。。。")
		var spinResData map[string]interface{}
		lo.Must0(jiliut.PostJson(f.spinUrl.String(), f.gameReq, &spinResData))

		// spinResData.Data
		if spinResData == nil {
			break
		}

		info := spinResData["info"].(map[string]interface{})
		spinAck := info["AckType"].(float64)
		if spinAck != 0 {
			slog.Error("pull spin data error", "acktype", spinAck, "user", f.username, "buy", f.buy)
			break
		}
		plateInfos := info["PlateInfo"].([]interface{})
		totalWin := info["TotalWin"].(float64)
		showIndex := info["ShowIndex"].(string)

		slog.Info("", "totalWin", totalWin)

		doc := &models.RawSpin{
			ID:       primitive.NewObjectID(),
			Times:    totalWin / 1.0, // bet is 1.0 for csh
			Type:     gendata.GameTypeNormal,
			HasGame:  len(plateInfos) > 1,
			Data:     spinResData,
			Selected: true,
		}
		doc.BucketId = gendata.GBuckets.GetBucket(doc.Times, doc.HasGame, db.BoundType(doc.Type))
		f.fetchHisotry(showIndex, doc)

		f.coll.InsertOne(context.TODO(), doc)
		fmt.Println("插入成功。。。")
	}
}
