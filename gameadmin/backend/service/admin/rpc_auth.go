package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/duck/lang"
	"game/duck/lazy"
	"game/duck/mongodb"
	"game/pb/_gen/pb/adminpb"
	"log/slog"
	"time"

	"github.com/samber/lo"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

const MinPlayerID = 1000

type AdminAuthRpc struct {
}

// 登录
func (AdminAuthRpc) Login(ctx context.Context, req *adminpb.AdminLoginReq) (*adminpb.AdminLoginResp, error) {
	ss := GetSession(ctx)

	if req.Username == "" {
		return nil, ss.Err("账号不存在")
	}

	if req.Password == "" {
		return nil, ss.Err("密码不存在")
	}

	userInfo := adminpb.DBAdminer{}
	query := bson.M{
		"Username": req.Username,
		"Password": req.Password,
	}
	err := CollAdminUser.FindOne(query, &userInfo)
	if err != nil {
		return nil, ss.Err("账号不存在或密码错误")
	}

	if userInfo.IsDelete == adminpb.IsDelete_YesDelete {
		return nil, ss.Err("账号不存在")
	}

	if userInfo.Status == adminpb.AdminStatus_Frozen {
		return nil, ss.Err("账号被冻结")
	}

	if userInfo.IsOpenGoogle {
		if code, err := NewGoogleAuth().VerifyCode(userInfo.GoogleCode, req.GooleCode); !code || err != nil {
			return nil, ss.Err("验证失败")
		}
	}

	token := uuid.NewString()
	//tokenExp := time.Now().AddDate(1, 0, 0)

	ip := ss.getString("ip", "")
	if IsBlockLoc(ip) {
		return nil, errors.New(lang.GetLang(ss.getString("lang", "en"), "IP Limit"))
	}
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": req.Username}, &user)
	// 用户审核状态判断
	var operator comm.Operator_V2
	if user.AppID != "admin" {
		err = CollAdminOperator.FindOne(bson.M{"AppID": user.AppID}, &operator)
		if err != nil {
			return nil, err
		}
		if operator.OperatorType == 2 && operator.ReviewStatus == 0 {
			return nil, errors.New("商户未审核")
		}
	}
	// 白名单IP地址验证
	if setting.OpenWhite {
		if err != nil {
			return nil, err
		}
		if user.AppID != "admin" {
			if len(operator.WhiteIps) > 0 && !lo.Contains(operator.WhiteIps, ip) {
				slog.Error("login failed!",
					"username", req.Username,
					"ip", ip,
					"whiteips", operator.WhiteIps,
				)
				return nil, errors.New(lang.GetLang(ss.getString("lang", "en"), "IP Limit"))
			}
			err = CollAdminOperator.Update(bson.M{"AppID": user.AppID}, bson.M{"$set": bson.M{"TokenExpireAt": mongodb.NewTimeStamp(time.Now())}})
			if err != nil {
				fmt.Println(err)
			}
		} else {
			if !lo.Contains(setting.WhiteIps, ip) {
				slog.Error("admin login failed!",
					"username", req.Username,
					"ip", ip,
					"whiteips", setting.WhiteIps,
				)
				return nil, errors.New(lang.GetLang(ss.getString("lang", "en"), "IP Limit"))
			}
			err = CollAdminOperator.Update(bson.M{"AppID": user.AppID}, bson.M{"$set": bson.M{"TokenExpireAt": mongodb.NewTimeStamp(time.Now())}})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	err = CollAdminOperator.Update(bson.M{"AppID": user.AppID}, bson.M{"$set": bson.M{"TokenExpireAt": mongodb.NewTimeStamp(time.Now())}})

	if lazy.CommCfg().IsDev {
		token = req.Username
	}

	tokenExp := time.Now().Add(time.Hour)
	err = CollAdminUser.UpdateId(userInfo.ID, bson.M{
		"$set": bson.M{
			"Token":         token,
			"TokenExpireAt": tokenExp,
		},
	})

	if err != nil {
		return nil, err
	}

	tokenMap.Delete(userInfo.Token)
	tokenMap.Store(token, &Token{
		Token:    token,
		ExpireAt: tokenExp,
		Pid:      userInfo.ID,
	})

	return &adminpb.AdminLoginResp{
		Token:    token,
		ExpireAt: mongodb.NewTimeStamp(tokenExp),
	}, nil
}
