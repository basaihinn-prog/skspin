package main

import (
	"encoding/json"
	"errors"
	"game/comm/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

func init() {
	RegMsgProc("/AdminInfo/AddRickRule", "添加风控管理", "AdminInfo", addRickRule, RickRuleReq{})
	RegMsgProc("/AdminInfo/GetRickRule", "查询风控管理", "AdminInfo", getRickRule, RickRuleReq{})
	mux.RegRpc("/AdminInfo/Interior/PushRickRuleTransferOut", "内部接口-发出预警风控信息", "AdminInfo", pushRickRuleTransferOut, RickRuleReq{})
	mux.RegRpc("/AdminInfo/Interior/PushRickRuleReturnRate", "内部接口-发出转账风控信息", "AdminInfo", pushRickRuleReturnRate, RickRuleReq{})
}

const TRANSFER_OUT_MONEY = 1000

type RickRuleReq struct {
	AppID    string
	RickRule string
	Type     string
}

type RickRule struct {
	Condition []*Condition `json:"condition"`
	Execute   Execute      `json:"execute"`
}

type RickRuleModel struct {
	AppID            string           `bson:"AppID"`
	ReturnRate       [][]RickRuleItem `bson:"ReturnRate"`
	OriginReturnRate string           `bson:"OriginReturnRate"`
	Transfer         RickRuleItem     `bson:"Transfer"`
	OriginTransfer   string           `bson:"OriginTransfer"`
}

type RickRuleItem struct {
	Condition Condition `bson:"condition"`
	Execute   Execute   `bson:"execute"`
}
type Condition struct {
	Filed    string  `json:"filed"`
	Contrast string  `json:"contrast"`
	Value    float64 `json:"value"`
}

type Execute struct {
	Value string `json:"value"`
}

func getRickRule(ctx *Context, ps RickRuleReq, ret *RickRuleModel) (err error) {
	user, _ := IsAdminUser(ctx)

	findBson := bson.M{}
	if ps.AppID != "" {
		findBson["AppID"] = ps.AppID
	} else {
		findBson["AppID"] = user.AppID
	}

	coll := NewOtherDB("GameAdmin").Collection("RickRule")
	err = coll.FindOne(findBson, &ret)
	if err != nil {
		ret = nil
		return nil
	}
	ret = &RickRuleModel{
		AppID:      ps.AppID,
		ReturnRate: [][]RickRuleItem{},
		Transfer:   RickRuleItem{},
	}
	return err
}
func addRickRule(ctx *Context, ps RickRuleReq, ret *getPermissionListResult) (err error) {
	if ps.RickRule == "" {
		return errors.New("设置风控失败")
	}

	var RickRuleList []*RickRule
	err = json.Unmarshal([]byte(ps.RickRule), &RickRuleList)
	if err != nil {
		return err
	}

	var rickRuleModel RickRuleModel
	rickRuleModel.ReturnRate = nil
	// 默认值要求变更为-1
	rickRuleModel.Transfer = RickRuleItem{
		Condition: Condition{Value: -1},
	}
	for _, rule := range RickRuleList {

		var returnRateAlarm []RickRuleItem
		for _, condition := range rule.Condition {
			if condition.Filed == "TransferOutMoney" {
				rickRuleModel.Transfer = RickRuleItem{
					Condition: *condition,
					Execute:   rule.Execute,
				}
			} else {

				returnRateAlarm = append(returnRateAlarm, RickRuleItem{
					Condition: *condition,
					Execute:   rule.Execute,
				})
			}
		}

		// 如果为转账预警时   returnRateAlarm值可能为nil
		if returnRateAlarm != nil {

			rickRuleModel.ReturnRate = append(rickRuleModel.ReturnRate, returnRateAlarm)
		}

	}

	var exitRickRuleModel RickRuleModel
	coll := NewOtherDB("GameAdmin").Collection("RickRule")
	flag := false
	err = coll.FindOne(bson.M{"AppID": ps.AppID}, &exitRickRuleModel)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			flag = true
		} else {
			return err
		}
	}
	rickRuleModel.AppID = ps.AppID

	if strings.ToLower(ps.Type) == "returnrate" {
		rickRuleModel.OriginReturnRate = ps.RickRule
	}

	if strings.ToLower(ps.Type) == "transfer" {
		rickRuleModel.OriginTransfer = ps.RickRule
	}

	if flag {
		err = coll.InsertOne(rickRuleModel)

	} else {

		setBson := bson.M{}

		if strings.ToLower(ps.Type) == "returnrate" {
			setBson["ReturnRate"] = rickRuleModel.ReturnRate
			setBson["OriginReturnRate"] = ps.RickRule
		}
		if strings.ToLower(ps.Type) == "transfer" {
			setBson["Transfer"] = rickRuleModel.Transfer
			setBson["OriginTransfer"] = ps.RickRule
		}

		err = coll.Update(bson.M{"AppID": ps.AppID}, bson.M{"$set": setBson})

	}
	return err
}

// 转账预警线
type transferOut struct {
	TransferOutMoney float64
}

// 回报预警线
type cond struct {
	List []*SingleCond
}
type SingleCond struct {
	Single float64
	Total  float64
	RTP    float64
}

func pushRickRuleTransferOut(ps RickRuleReq, ret *transferOut) (err error) {

	var rickRuleModel RickRuleModel

	if ps.AppID == "" {
		ret.TransferOutMoney = TRANSFER_OUT_MONEY
		return nil
	}

	coll := NewOtherDB("GameAdmin").Collection("RickRule")
	err = coll.FindOne(bson.M{"AppID": ps.AppID}, &rickRuleModel)
	if err != nil {

		if err == mongo.ErrNoDocuments {

			ret.TransferOutMoney = TRANSFER_OUT_MONEY
		} else {
			ret.TransferOutMoney = 0
		}

		return nil
	}

	ret.TransferOutMoney = rickRuleModel.Transfer.Condition.Value

	return nil
}

func pushRickRuleReturnRate(ps RickRuleReq, ret *cond) (err error) {

	defaultSingleCond := SingleCond{
		Single: 100000,
		Total:  1000000,
		RTP:    106,
	}
	if ps.AppID == "" {
		ret.List = append(ret.List, &defaultSingleCond)
		return nil
	}

	var rickRuleModel RickRuleModel

	coll := NewOtherDB("GameAdmin").Collection("RickRule")
	err = coll.FindOne(bson.M{"AppID": ps.AppID}, &rickRuleModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			ret.List = append(ret.List, &defaultSingleCond)
		} else {
			ret.List = []*SingleCond{}
		}
		return nil
	}
	ret.List = []*SingleCond{}

	if rickRuleModel.ReturnRate == nil {
		return nil
	}

	for _, items := range rickRuleModel.ReturnRate {
		var singleCond SingleCond
		for _, childItem := range items {

			switch childItem.Condition.Filed {
			case "Single":
				singleCond.Single = childItem.Condition.Value
				break
			case "Total":
				singleCond.Total = childItem.Condition.Value
				break
			case "RTP":
				singleCond.RTP = childItem.Condition.Value
				break
			}

		}
		ret.List = append(ret.List, &singleCond)
	}

	return nil
}
