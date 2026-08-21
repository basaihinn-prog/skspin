package hacksawcomm

import (
	"encoding/json"
	"serve/comm/jwtutil"
)

// Bet 结构体表示 bets 数组中的单个对象
type Bet struct {
	BetAmount string `json:"betAmount"`
	BuyBonus  string `json:"buyBonus"`
}

// ContinueInstructions 结构体表示 continueInstructions 对象
type ContinueInstructions struct {
	Action string `json:"action"`
}

// SessionData 结构体表示整个 JSON 数据
type SessionData struct {
	Seq                  int                   `json:"seq"`
	SessionUUID          string                `json:"sessionUuid"`
	Bets                 []Bet                 `json:"bets,omitempty"` // omitempty 表示字段为空时不序列化
	OfferID              string                `json:"offerId"`        // 使用指针类型表示可以为 null
	PromotionID          string                `json:"promotionId"`    // 使用指针类型表示可以为 null
	Autoplay             bool                  `json:"autoplay"`
	RoundID              string                `json:"roundId,omitempty"`
	ContinueInstructions *ContinueInstructions `json:"continueInstructions,omitempty"`
}

func ParseHackSawReq(payload []byte) (int64, *SessionData, error) {
	var req SessionData
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return 0, nil, err
	}
	pid, err := jwtutil.ParseToken(req.SessionUUID)
	if err != nil {
		return 0, nil, err
	}
	return pid, &req, nil
}
