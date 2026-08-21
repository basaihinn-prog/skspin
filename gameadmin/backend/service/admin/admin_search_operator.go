package main

import (
	"errors"
	"game/comm"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func init() {
	RegMsgProc("/AdminInfo/SearchOperator", "查询商户", "AdminInfo", getSearchOperator, getSearchOperatorParams{})
	RegMsgProc("/AdminInfo/GetOperatorById", "通过appid获取商户", "AdminInfo", getOperatorById, getOperatorParams{})

}

type getSearchOperatorParams struct {

	//RealAppID        string `json:"AppID" bson:"AppSecret"`           //商户名称
	Status           int64  `json:"Status" bson:"Status"`                     //状态
	AppID            string `json:"AppID" bson:"AppID"`                       //商户名称
	BelongingCountry string `json:"BelongingCountry" bson:"BelongingCountry"` //商户类型
	OperatorType     int    `json:"OperatorType" bson:"OperatorType"`         //商户类型
	CreatedStartTime int64  `json:"CreatedStartTime"`
	CreatedEndTime   int64  `json:"CreatedEndTime"`
	PageIndex        int64  `json:"PageIndex"`
	PageSize         int64  `json:"PageSize"`
	AppSecret        string `json:"AppSecret"`
	ReviewStatus     int64  `json:"ReviewStatus"`
}

type getOperatorParams struct {
	AppID string `json:"AppID" bson:"AppID"` //商户id
}

type searchOperatorResultV2 struct {
	List     []*comm.Operator_V2
	AllCount int64
}

type getOperatorByIdResult struct {
	OperatorInfo *comm.Operator_V2
}

func getSearchOperator(ctx *Context, ps getSearchOperatorParams, ret *searchOperatorResultV2) (err error) {
	user, _ := IsAdminUser(ctx)
	var list []*comm.Operator_V2
	query := bson.M{}
	if ps.CreatedStartTime != 0 {
		startTime := mongodb.NewTimeStamp(time.UnixMilli(ps.CreatedStartTime))
		endTime := mongodb.NewTimeStamp(time.UnixMilli(ps.CreatedEndTime))
		query = bson.M{
			//"Name":user.AppID,
			"CreateTime": bson.M{
				"$gte": startTime, // 大于或等于开始时间
				"$lte": endTime,   // 小于或等于结束时间
			},
		}
	}
	var ok = false
	if ps.AppID != "" { //模糊查询
		//query["AppID"] = bson.M{"$regex": "/" + ps.AppID + "/", "$options": 'i'}
		query["AppID"] = primitive.Regex{
			Pattern: ps.AppID,
			Options: "i",
		}
		ok = true
	}

	if ps.AppSecret != "" { //模糊查询
		//query["AppID"] = bson.M{"$regex": "/" + ps.AppID + "/", "$options": 'i'}
		query["AppSecret"] = primitive.Regex{
			Pattern: ps.AppSecret,
			Options: "i",
		}
		ok = true
	}
	if user.GroupId == 2 {
		query["Name"] = user.AppID
		ok = true
	}
	if ps.OperatorType != 0 {
		query["OperatorType"] = ps.OperatorType
		ok = true
	}

	if ps.Status != -1 {
		query["Status"] = ps.Status
		ok = true
	}
	if ps.ReviewStatus == -1 {
		query["ReviewStatus"] = 0
		ok = true
	}

	//  二期根据所属国家查询
	if ps.BelongingCountry != "" && ps.BelongingCountry != "ALL" {
		query["BelongingCountry"] = ps.BelongingCountry
		ok = true
	}

	if user.GroupId == 3 {
		query["AppID"] = user.AppID
		ok = true
	}
	if ok == false {
		switch user.GroupId {
		case 0:
			query["Name"] = "admin"
		case 1:
			query["Name"] = "admin"
		case 2:
			query["AppID"] = user.AppID
		case 3:
			query["AppID"] = user.AppID
		default:
		}
	}

	pageFileter := mongodb.FindPageOpt{
		Page:     ps.PageIndex,
		PageSize: ps.PageSize,
		Sort:     bson.M{"CreateTime": -1},
		Query:    query,
	}
	ret.AllCount, err = CollAdminOperator.FindPage(pageFileter, &list)
	if err != nil {
		return err
	}
	for i, op := range list {
		if op.OperatorType == 1 {

			count, err := CollAdminOperator.CountDocuments(bson.M{"Name": op.AppID})
			if err != nil {
				return err
			}
			if count > 0 {
				list[i].HasChildren = true
			} else {
				list[i].HasChildren = false
			}
		}
	}

	ret.List, err = GetOperatorProfit(ctx, list)
	if err != nil {
		return err
	}

	return err
}

func getOperatorById(ctx *Context, ps getOperatorParams, ret *getOperatorByIdResult) (err error) {
	user, _ := IsAdminUser(ctx)

	if ps.AppID == "" { //模糊查询
		return errors.New("AppID is empty")
	}
	if user.GroupId == 3 && ps.AppID != user.AppID {
		return errors.New("Insufficient operator permissions")
	}

	err = CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, &ret.OperatorInfo)

	if ret.OperatorInfo.CurrencyManufactureVisibleOff != nil {

		for _, v := range ret.OperatorInfo.CurrencyManufactureVisibleOff {

			ret.OperatorInfo.CurrencyVisibleOff = v

			if v == 0 {
				break
			}

		}

	}

	if err != nil {
		return err
	}
	if user.GroupId == 2 {
		if ret.OperatorInfo.Name != user.AppID && ret.OperatorInfo.AppID != user.AppID {
			return errors.New("Insufficient operator permissions")
		}
	}

	return err
}
