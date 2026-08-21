package internal

import (
	"log/slog"
	"math/rand/v2"
	"net/http"
	"serve/comm/robot"
	"serve/comm/ut"
	"serve/servicepp/ppcomm"
	"strings"
	"sync"
	"time"

	"github.com/samber/lo"
)

type Fetcher struct {
	game, user string

	mgckey string
	// location *url.URL
	ul         string
	httpClient *http.Client
	lastResp   ppcomm.Variables
	c          float64
	l          int
	ty         int
	bonusTy    int
}

type FetchDataParam struct {
	Game        string
	ArrC        map[float64]int
	GBuckets    *ppcomm.Buckets
	Ty          int
	bonusTy     int
	IsEnough    bool //拉取够了设置为true
	L           int
	InsertC     float64
	ApiVersion  string
	HuiDuSerial string
}

var robotMap map[string]bool
var robotLock sync.Mutex

func FetchData(ps *FetchDataParam) {
	if ps == nil {
		return
	}
	robotMap = make(map[string]bool)
	for ; ; time.Sleep(30 * time.Second) {
		//数量够了就停止
		if robot.CheckGameEnd(ps.Game) || CheckGameEnough(ps.Game, ps.Ty) {
			break
		}
		robots := robot.GetRobot(ps.Game)

		if len(robots) <= 0 {
			slog.Info("get robot fail, robot busy, GameId:%s, type:%d\n", ps.Game, ps.Ty)
		}
		//fmt.Println(len(robotMap))
		for i := range robots {
			if ps.IsEnough {
				if rand.Int()%100 < 70 {
					continue
				}
			}
			robotLock.Lock()
			if _, ok := robotMap[robots[i].ID.Hex()]; ok {
				robotLock.Unlock()
				continue
			}
			robotMap[robots[i].ID.Hex()] = true
			robotLock.Unlock()

			fetcher := NewFetcher(ps.Game, ps.ApiVersion, ps.HuiDuSerial, ps.bonusTy)
			if fetcher == nil {
				robotLock.Lock()
				delete(robotMap, robots[i].ID.Hex())
				robotLock.Unlock()
				continue
			}
			ppR := NewPPRobot(robots[i], fetcher, ps)
			go ppR.Run()
		}
	}
	slog.Info("over: FetchData")
}

func (f *Fetcher) Do(ps ppcomm.Variables) ppcomm.Variables {
	// form, _ := url.ParseQuery("action=doSpin&symbol=vs20olympx&c=0.5&l=20&bl=0&index=3&counter=5&repeat=0&mgckey=123")
	lo.Must0(ps.Str("action") != "")

	if ps.Str("action") == "doSpin" {
		ps.SetFloat("c", f.c)
		ps.SetInt("l", f.l)
		ps.Set("bl", "0")
	}
	if ps.Str("action") == "doBonus" {
		ps.SetInt("ind", f.bonusTy)
	}
	if f.game == "vs40bigjuan" && ps.Str("status") == "1" && f.ty == ppcomm.GameTypeGame {
		ps.SetInt("pur", 0)
		ps.Delete("pur")
	} else if f.ty == ppcomm.GameTypeGame && f.game != "vs40bigjuan" {
		ps.SetInt("pur", 0)
	} else if f.ty == ppcomm.GameTypeSuperGame1 {
		ps.SetInt("pur", 1)
	} else if f.ty == ppcomm.GameTypeSuperGame2 {
		ps.SetInt("pur", 2)
	}
	ps.Set("symbol", f.game)
	ps.SetInt("index", f.lastResp.Int("index")+1)
	ps.SetInt("counter", f.lastResp.Int("counter")+1)
	ps.Set("mgckey", f.mgckey)
	ps.Set("repeat", "0")

	req, _ := http.NewRequest(http.MethodPost, f.ul, strings.NewReader(ps.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	// http.DefaultClient.Do(req)
	body, _ := lo.Must2(ut.DoHttpReq(f.httpClient, req))

	vars := ppcomm.ParseVariables(string(body))
	f.lastResp = vars

	//slog.Info("robot_game365_center.DO", "ps", ps.Encode(),
	//	"user", f.user,
	//	"na", vars["na"],
	//	"resp", mux.TruncMsg(body, 256),
	//)
	return vars
}

func (f *Fetcher) NextAction() string {
	return f.lastResp.Str("na")
}

func (f *Fetcher) LastResp() ppcomm.Variables {
	return f.lastResp
}
