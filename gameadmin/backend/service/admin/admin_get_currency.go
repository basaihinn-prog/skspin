package main

import (
	"game/comm"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/GetCurrency", "获取币种", "AdminInfo", getCurrency, getLandParams{})
}

type getLandParams struct{}

type getLandResult struct {
	List     []*comm.CurrencyType
	Allcount int `json:"Allcount"`
}

func getCurrency(ctx *Context, ps getGameCofigParams, ret *getLandResult) (err error) {
	var list []*comm.CurrencyType
	err = DB.Collection("CurrencyType").FindAll(bson.M{}, &list)
	ret.Allcount = len(list)
	ret.List = list
	if ctx.Lang != "zh" {
		for i, currencyType := range ret.List {
			ret.List[i].CurrencyName = currencyType.CurrencyCode
		}
	}

	if err != nil {
		return err
	}
	return
}
