package room

import (
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log/slog"
	"serve/comm/redisx"
	"serve/comm/ut"
	"serve/serviceSpribeGame/aviator/comm"
	"strconv"
)

func ClientPreviousRoundInfo(c *websocket.Conn, messageType websocket.MessageType, data *ut.SFSObject, r *Room) {
	//content := `{"id":13,"targetController":1,"content":{"c":"updateCurrentBets","p":{"betsCount":1,"code":200}}}`
	content := &ContentCache{
		C: "previousRoundInfoResponse",
	}
	var err error
	content.P.Bets, err = GetSortedBets(redisx.GetClient(), r.Name+"_previous")
	if err != nil {
		slog.Error("查询上一手失败 GetSortedBets Err", "err", err)
	}
	content.P.RoundsInfo, err = GetPreviousRoundInfo(r.Name + "_previous:" + "roundInfo")
	content.P = PCache{
		Bets:       content.P.Bets,
		RoundsInfo: content.P.RoundsInfo,
		Code:       comm.Succ,
	}
	rsp := CommBodyCache{
		Id:               13,
		TargetController: 1,
		Content:          content,
	}
	//marshal, err := json.Marshal(rsp)
	//if err != nil {
	//	slog.Error("上一手反序列化失败 Marshal Err", "err", err)
	//	rsp := GetErrMsg("previousRoundInfoResponse", comm.ServerErr, "Marshal Err")
	//	rspMarshal, _ := json.Marshal(rsp)
	//	c.WriteMessage(messageType, rspMarshal)
	//}
	marshal, _ := GetClientPreviousRoundInfoRsp(&rsp).ToBinary()
	c.WriteMessage(websocket.BinaryMessage, marshal)
}
func GetClientPreviousRoundInfoRsp(initRsp *CommBodyCache) *ut.SFSObject {
	so := ut.NewSFSObject()
	p := ut.NewSFSObject()

	pp := ut.NewSFSObject()
	roundInfo := ut.NewSFSObject()
	roundInfo.PutDouble("multiplier", initRsp.Content.P.RoundsInfo.Multiplier)
	roundInfo.PutLong("roundStartDate", initRsp.Content.P.RoundsInfo.RoundStartDate)
	roundInfo.PutLong("roundEndDate", initRsp.Content.P.RoundsInfo.RoundEndDate)
	roundInfo.PutInt("roundId", int32(initRsp.Content.P.RoundsInfo.RoundID))
	pp.PutSFSObject("roundInfo", roundInfo)
	pp.PutInt("code", int32(initRsp.Content.P.Code))
	bets := ut.NewSFSArray()
	for _, bet := range initRsp.Content.P.Bets {
		temp := ut.NewSFSObject()
		temp.PutDouble("bet", bet.Bet)
		temp.PutDouble("winAmount", bet.WinAmount)
		temp.PutString("currency", bet.Currency)
		temp.PutBool("win", bet.Win)
		atoi, _ := strconv.Atoi(bet.RoundBetId)
		temp.PutInt("roundBetId", int32(atoi))
		temp.PutDouble("payout", bet.Payout)
		temp.PutBool("isFreeBet", false)
		temp.PutString("profileImage", bet.ProfileImage)
		temp.PutString("username", bet.Username)
		bets.Add(temp, ut.SFS_OBJECT, true)
	}
	pp.PutSFSArray("bets", bets)

	p.PutSFSObject("p", pp)
	p.PutString("c", "previousRoundInfoResponse")

	so.AddCreatePAC(p, 1, 13)
	return so
}

type ContentCache struct {
	Api string `json:"api,omitempty"`
	Cl  string `json:"cl,omitempty"`
	Ct  int    `json:"ct,omitempty"`
	Ms  int    `json:"ms,omitempty"`
	Tk  string `json:"tk,omitempty"`

	Zn string `json:"zn,omitempty"`
	Un string `json:"un,omitempty"`
	Pw string `json:"pw,omitempty"`
	P  PCache `json:"p,omitempty"` //需要调查有没有同名但是字段不同的

	Pi int             `json:"pi,omitempty"`
	Rl [][]interface{} `json:"rl,omitempty"`
	Rs int             `json:"rs,omitempty"`

	C string `json:"c,omitempty"`
	R int    `json:"r,omitempty"`
}

type PCache struct { //猜测是通用字段，内部包含很多字段
	Token            string            `json:"token,omitempty"`
	Currency         string            `json:"currency,omitempty"`
	Lang             string            `json:"lang,omitempty"`
	SessionToken     string            `json:"sessionToken,omitempty"`
	Platform         *Platform         `json:"platform,omitempty"`
	Version          string            `json:"version,omitempty"`
	AdditionalParams *AdditionalParams `json:"additionalParams,omitempty"`

	RoundsInfo         *RoundInfo     `json:"roundInfo,omitempty"`
	Code               int            `json:"code,omitempty"`
	ActiveBets         []interface{}  `json:"activeBets,omitempty"`
	OnlinePlayers      int            `json:"onlinePlayers,omitempty"`
	ChatSettings       *ChatSettings  `json:"chatSettings,omitempty"`
	ChatHistory        []*ChatHistory `json:"chatHistory,omitempty"`
	User               *User          `json:"user,omitempty"`
	Config             *Config        `json:"config,omitempty"`
	RoundID            int64          `json:"roundID,omitempty"`
	StageID            int            `json:"stageId,omitempty"`
	ActiveFreeBetsInfo []interface{}  `json:"activeFreeBetsInfo,omitempty"`

	Bets      []*Bets     `json:"bets,omitempty"`
	BetsCount int         `json:"betsCount,omitempty"`
	CashOuts  []*CashOuts `json:"cashOuts,omitempty"`

	Bet        float64 `json:"bet,omitempty"`
	ClientSeed string  `json:"clientSeed,omitempty"`
	BetID      int     `json:"betId,omitempty"`
	FreeBet    bool    `json:"freeBet,omitempty"`

	PlayerID     string `json:"player_id,omitempty"`
	ProfileImage string `json:"profileImage,omitempty"`
	Username     string `json:"username,omitempty"`

	Cashouts    []*Cashouts `json:"cashouts,omitempty"`
	Multiplier  float64     `json:"multiplier,omitempty"`
	OperatorKey string      `json:"operatorKey,omitempty"`

	NewBalance float64 `json:"newBalance,omitempty"`

	NewStateID int `json:"newStateId,omitempty"`
	TimeLeft   int `json:"timeLeft,omitempty"`

	X      float64 `json:"x,omitempty"`
	CrashX float64 `json:"crashX,omitempty"`
	Crash  bool    `json:"crash,omitempty"`

	MaxMultiplier float64 `json:"maxMultiplier,omitempty"`
	RoundId       int64   `json:"roundId,omitempty"`

	LastBetId int64 `json:"lastBetId,omitempty"`
}

type CommBodyCache struct {
	Id               int           `json:"id"`
	Channel          string        `json:"channel,omitempty"`
	TargetController int           `json:"targetController"`
	Content          *ContentCache `json:"content"`
}
