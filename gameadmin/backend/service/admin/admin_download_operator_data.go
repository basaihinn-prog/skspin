package main

import (
	"errors"
	"fmt"
	"game/comm"
	"game/duck/lang"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/DownloadOperatorData", "下载运营商数据", "AdminInfo", downloadOperatorData, downloadOperatorParams{
		OperatorId: 0,
	})
}

type downloadOperatorParams struct {
	OperatorId int64
}

type OperatorData struct {
	AppId         string
	AppSecret     string
	Address       string
	WhiteIps      []string
	AdminUserName string
	AdminPassword string
	GoogleCode    string
	Qrcode        string
	ApiUrl        string
}

type downloadOperatorResults struct {
	List []*OperatorData
}

func downloadOperatorData(ctx *Context, ps downloadOperatorParams, ret *downloadOperatorResults) (err error) {
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	var operators []*comm.Operator
	if ps.OperatorId > 0 { //选择了运营商
		var operator *comm.Operator
		err = CollAdminOperator.FindId(ps.OperatorId, &operator)
		if err != nil {
			return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
		}
		operators = append(operators, operator)
	} else { //没有选择运营商
		err = CollAdminOperator.FindAll(bson.M{}, &operators)
	}
	appIds := make([]string, 0, len(operators))
	appNames := make([]string, 0, len(operators))
	for i := range operators {
		if operators[i].AppID == "admin" {
			continue
		}
		appIds = append(appIds, operators[i].AppID)
		appNames = append(appNames, fmt.Sprintf("admin_%s", operators[i].AppID))
	}

	var Users []*comm.User
	err = CollAdminUser.FindAll(
		bson.M{
			"AppID":    bson.M{"$in": appIds},
			"Username": bson.M{"$in": appNames},
		}, &Users)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}

	findUser := func(appId string) *comm.User {
		for i := range Users {
			if Users[i].AppID == appId {
				return Users[i]
			}
		}
		return nil
	}
	ret.List = make([]*OperatorData, 0, len(operators))
	for i := range operators {
		user := findUser(operators[i].AppID)
		if user == nil {
			ret.List = append(ret.List, &OperatorData{
				AppId:     operators[i].AppID,
				AppSecret: operators[i].AppSecret,
				Address:   operators[i].Address,
				WhiteIps:  operators[i].WhiteIps,
				//AdminUserName: user.UserName,
				//AdminPassword: user.PassWord,
				//GoogleCode:    user.GoogleCode,
				//Qrcode:        user.Qrcode,
				ApiUrl: setting.ApiUrl,
			})
			continue
		}
		ret.List = append(ret.List, &OperatorData{
			AppId:         operators[i].AppID,
			AppSecret:     operators[i].AppSecret,
			Address:       operators[i].Address,
			WhiteIps:      operators[i].WhiteIps,
			AdminUserName: user.UserName,
			AdminPassword: user.PassWord,
			GoogleCode:    user.GoogleCode,
			Qrcode:        user.Qrcode,
			ApiUrl:        setting.ApiUrl,
		})
	}
	return err
}
