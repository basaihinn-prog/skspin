package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HistoryRecord struct {
	Type           string
	RoundIndex     string
	NetValue       string
	Bet            string
	Win            string
	CreateTime     int64
	PreMoney       string
	PostMoney      string
	ItemName       json.RawMessage
	CurrencyNo     int
	AlterId        int
	BetType        int
	JpContribution string
}

type SingleRoundLogSummary struct {
	LogIndex string
	Title    []string
	Desc     []string
	AlterId  int
}

type LogPlateInfo struct {
	Type         int
	Direction    int
	Array        []int
	Plate        string
	SymbolLength string
	ColorNo      int
	List         []*struct {
		S, L, C int
		W       string
		D       string
		Award   []int
	}
	TopStr        json.RawMessage
	Mult          string
	MysterySymbol interface{}
}

type HistoryDoc struct {
	ID                      primitive.ObjectID `bson:"_id"`
	Pid                     int64
	HistoryRecord           *HistoryRecord
	SingleRoundLogSummaries []*SingleRoundLogSummary
	LogPlateInfos           map[string][]*LogPlateInfo
	Tid                     string             `json:"tid" bson:"tid"`
	OId                     primitive.ObjectID `bson:"oid"`
}
