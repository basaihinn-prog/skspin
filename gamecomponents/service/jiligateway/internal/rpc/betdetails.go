package rpc

import (
	"cmp"
	"context"
	"fmt"
	"game/comm/db"
	"game/comm/mux"
	"game/duck/ut2/jwtutil"
	"game/service/jiligateway/internal/gamedata"
	"go.mongodb.org/mongo-driver/bson"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	mux.RegRpc("/jiligateway/getBetDetailsUrl", "gameinfo", "game-api", getBetDetailsUrl, getBetDetailsUrlPs{
		Gid:   "jili_2_csh",
		BetID: "6683766ad8add15a97539202",
		Lang:  "en",
	})
}

type getBetDetailsUrlPs struct {
	Gid   string
	BetID string
	Lang  string
	Token string
}
type getBetDetailsUrlRet struct {
	Url string
}

// https://uat-history.kafa010.com/ingame/gamehistory/csh/17199133299370523741?token=eyJQIjoxMDA3OTQsIkUiOjE3MTk5NTYxMzgsIlMiOjEwMDQsIkQiOiJqaWxpXzJfY3NoIn0.OyVYQZ4xru5vAmFZXzPr2btfmZRnyCoUEoyMOQnDWt8&game=2&layout=nomenu

// https://uat-www.jlfafafa3.com/ingame/gamehistory/csh/1719827657170767002?token=9d949cb1f5519ce045b396c780f803950f03a3c9&game=2&ri=1719827657170767002&li=1719827657170767002&layout=nomenu
func getBetDetailsUrl(ps getBetDetailsUrlPs, ret *getBetDetailsUrlRet) (err error) {
	tmpl := gamedata.Get().BoBetDetailsTempUrl

	u, err := url.Parse(tmpl)
	if err != nil {
		return
	}

	gid, _ := splitJILI(ps.Gid)
	gameNo, err := strconv.Atoi(gid)
	if err != nil {
		return
	}

	// gid := strings.TrimPrefix(ps.Gid, "pg_")

	// u.Path = path.Join("/history", gid+".html")

	q := u.Query()
	q.Set("game", gid)
	// q.Set("gid", gid)
	// q.Set("lang", "en-US")

	// u.RawQuery = q.Encode()

	// ret.Url = u.String()

	coll := db.Collection2(ps.Gid, "BetHistory")

	id := ps.BetID

	var doc HistoryDoc
	proj := db.D(
		"pid", 1,
		"historyrecord.roundindex", 1,
	)
	err = coll.FindOne(context.Background(), bson.D{{Key: "tid", Value: id}}, options.FindOne().SetProjection(proj)).Decode(&doc)
	if err != nil {
		return
	}

	token, err := jwtutil.NewTokenWithData(doc.Pid, time.Now().Add(12*time.Hour), ps.Gid)
	q.Set("token", token)

	u.Path = fmt.Sprintf("/%s/ingame/gamehistory/%s/%s", getlang4(ps.Lang), gamedata.GameMap()[gameNo].Id, doc.HistoryRecord.RoundIndex)

	u.RawQuery = q.Encode()

	ret.Url = u.String()
	return
}

type HistoryRecord struct {
	RoundIndex string
}

type HistoryDoc struct {
	ID            primitive.ObjectID `bson:"_id"`
	Pid           int64
	HistoryRecord *HistoryRecord
}

func splitJILI(gameID string) (id, name string) {
	arr := strings.SplitN(gameID, "_", 3)
	lo.Must0(len(arr) == 3)

	id, name = arr[1], arr[2]
	return
}

var lang4map = map[string]string{
	"en":  "en-US",
	"da":  "da-DK",
	"de":  "de-DE",
	"es":  "es-AR",
	"fi":  "fi-FI",
	"fr":  "fr-FR",
	"idr": "id-ID",
	"it":  "it-IT",
	"ja":  "ja-JP",
	"ko":  "ko-KR",
	"nl":  "nl-NL",
	"no":  "nn-NO",
	"pl":  "pl-PL",
	"pt":  "pt-PT",
	"ro":  "ro-RO",
	"ru":  "ru-RU",
	"sv":  "sv-SE",
	"th":  "th-TH",
	"tr":  "tr-TR",
	"vi":  "vi-VN",
	"zh":  "zh-CN",
	// "my": "",
}

func getlang4(lang2 string) string {
	// switch lang2 {
	// case "en":
	// 	return "en-US"
	// case "da":
	// 	return "da-DK"
	// }

	ret := lang4map[lang2]

	return cmp.Or(ret, "en-US")
}
