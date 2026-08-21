package comm

import (
	"encoding/json"
	"log/slog"
	"math"
	"math/rand/v2"
	"serve/comm/db"
	"serve/comm/jwtutil"
	"serve/comm/ut"
	"strconv"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

const (
	GameID = "spribe_01"
)

const (
	NETOPEN   = 0
	NETMSG    = 1
	NETCLOSE  = 2
	TIMER_1s  = 3
	TIMER_100 = 4

	AviatorLoginIdReq = iota
	AviatorGameInfoIdReq
	AviatorCurrentBetsInfoIdReq
	AviatorBetIdReq
	AviatorCancelBetIdReq
	AviatorCashOutIdReq
	AviatorBetHistoryIdReq
	AviatorGameStatePingIdReq
	AviatorAddChatMessageIdReq
	AviatorLikeIdReq
	ClientSearchGifs
	AviatorPreviousRoundInfoIdReq
	AviatorGetHugeWinsInfoIdReq
	AviatorGetTopWinsInfoIdReq
	AviatorGetTopRoundsInfoIdReq
	AviatorChangeProfileImageIdReq
	AviatorRoundFairnessIdReq
	ServerSeedHandler
)

const (
	LogType_Bet    = 0 //投注
	LogType_ReMove = 1 //取消投注
	FinishType_0   = 0 // 该次游戏未结束
	FinishType_1   = 1 // 该次游戏结束
)

type MSG struct {
	Conn *websocket.Conn
	Msg  *ut.SFSObject
	Typ  uint8
}

func calcMax(rtp float64) float64 {
	// 公式计算 y = 97 * 1024 / (1024 - x)
	ran := 1048576 //1024*1024
	x := rand.IntN(ran)
	max := rtp * float64(ran) / float64(ran-x)

	// 当 x 取值在 0~20 时，如果结果小于 1.00，则取 1.00
	if (x >= 0 && x <= 20) || max < 1.00 {
		max = 1.00
	}

	// 保留两位小数
	max = math.Round(max*100) / 100
	return max
}

func CalcPoints(rtp float64) (arr []float64) {
	//base := 1.08
	max := calcMax(rtp)
	highY := 0.0
	k := 1.0
	numT := 0.0
	for {
		numT += 0.1
		highY = 1.08
		res := truncateFloat(math.Pow(highY, numT*k), 2)
		if res > max {
			arr = append(arr, max)
			break
		}
		arr = append(arr, res)
	}
	return
}

func CalcPointsDev(rtp float64) (arr []float64) {
	//base := 1.08
	max := calcMax(rtp)
	const bbb = 6
	if max > bbb {
		max = bbb
	}
	if max < 1 {
		max = 1
	}
	num := float64(1)
	arr = append(arr, num)
	for {
		if num > 2.5 {
			num = truncateFloat(num+0.03, 2)
		} else {
			num = truncateFloat(num+0.01, 2)
		}
		if num > max {
			break
		}
		arr = append(arr, num)
	}
	arr[len(arr)-1] = max
	return
}

func CalcPointsTes(rtp float64) (arr []float64) {
	//base := 1.08
	max := calcMax(rtp)
	highY := 0.0
	k := 1.0
	numT := 0.0
	for {
		numT += 0.1
		highY = 1.08
		res := truncateFloat(math.Pow(highY, numT*k), 2)
		if res > max {
			arr = append(arr, max)
			break
		}
		arr = append(arr, res)
	}
	return
}

type WsSessionBody struct {
	Token    string        `json:"token"`
	PlayerId string        `json:"playerId"`
	Plr      *db.DocPlayer `json:"plr"`
}

func GetSession(c *websocket.Conn) (WsSessionBody, error) {
	session := WsSessionBody{}
	marshal, err := json.Marshal(c.Session())
	if err != nil {
		slog.Error("GetPlayerBySession", "json.Marshal err", err)
		return WsSessionBody{}, err
	}
	err = json.Unmarshal(marshal, &session)
	if err != nil {
		slog.Error("GetPlayerBySession", "json.Unmarshal err", err)
		return WsSessionBody{}, err
	}
	return session, nil
}

func GetPlr(c *websocket.Conn) (*db.DocPlayer, error) {
	var plr *db.DocPlayer
	session, err := GetSession(c)
	if err != nil {
		slog.Error("GetPlr Err", "GetSession err", err)
	}
	playerIdStr := session.PlayerId
	playerId, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		slog.Error("ClientLogin Err", "strconv.ParseInt err", err)
	}
	plr, err = db.GetDocPlayer(playerId)
	if err != nil {
		slog.Error("ClientLogin Err", "db.GetDocPlayer err", err)
	}
	return plr, nil
}
func truncateFloat(num float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(num*shift) / shift
}

func TokenGetPlr(tk string) (*db.DocPlayer, error) {
	var plr *db.DocPlayer
	playerId, _, err := jwtutil.ParseTokenData(tk)
	if err != nil {
		slog.Error("ClientLogin Err", "jwt.ParseTokenData", err)
	}
	plr, err = db.GetDocPlayer(playerId)
	if err != nil {
		slog.Error("spin_stat_err", "GetPlr", err)
		return nil, err
	}
	return plr, nil
}
