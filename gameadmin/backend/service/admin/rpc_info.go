package main

import (
	"context"
	"errors"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/adminpb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type AdminInfoRpc struct{}

func (AdminInfoRpc) GetAdminInfo(ctx context.Context, req *pb.Empty) (*adminpb.AdminInfoResp, error) {
	return nil, errors.New("Method Unimplemented")
	/*ss := GetSession(ctx)

	userInfo := adminpb.DBAdminer{}
	if err := CollAdminUser.FindId(ss.Pid, &userInfo); err != nil {
		return nil, err
	}

	query := bson.M{}
	if userInfo.GroupId != 1 {
		group := adminpb.DBAuth{}
		err := CollAdminAuth.FindId(userInfo.GroupId, &group)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				return nil, err
			}
		}
		//获取菜单列表
		menus := strings.Split(group.MenuIds, ",")

		menusInt := make([]int, 0)
		for _, v := range menus {
			menusID, _ := strconv.Atoi(v)
			menusInt = append(menusInt, menusID)
		}

		query["_id"] = bson.M{"$in": menusInt}
	}

	// 这个是超级权限 直接拿所有的菜单id
	menuList := make([]*adminpb.DBMenu, 0)
	_, err := CollAdminMenu.FindPage(mongodb.FindPageOpt{
		Page:     1,
		PageSize: 1000,
		Sort:     bson.M{"Sort": 1}, //可以为空
		Query:    query,
	}, &menuList)
	if err != nil {
		return nil, err
	}
	return &adminpb.AdminInfoResp{
		ID:       ss.Pid,
		Username: userInfo.Username,
		Gid:      userInfo.GroupId,
		MenuList: menuList,
	}, nil*/
}

// 获取管理员列表
func (AdminInfoRpc) AdminerList(ctx context.Context, req *adminpb.PageListReq) (*adminpb.AdminListResp, error) {
	return nil, errors.New("Method Unimplemented")
	//uid := common.GetSession(ctx).Pid
	/*var list []*adminpb.DBAdminer
	filter := mongodb.FindPageOpt{
		Page:     req.PageIndex,
		PageSize: req.PageSize,
		Sort:     bson.M{"_id": -1}, //可以为空
		Query:    bson.M{"IsDelete": adminpb.IsDelete_NoDelete},
	}
	count, err := CollAdminUser.FindPage(filter, &list)
	if err != nil {
		return nil, err
	}
	return &adminpb.AdminListResp{Count: count, List: list}, nil*/
}

// 添加管理员
func (AdminInfoRpc) AddAdminer(ctx context.Context, req *adminpb.AddAdminReq) (*pb.Empty, error) {
	return nil, errors.New("Method Unimplemented")
	/*
		ss := GetSession(ctx)
		if req.Username == "" {
			return nil, ss.Err("账号不存在")
		}

		userInfo := adminpb.DBAdminer{}
		err := CollAdminUser.FindOne(bson.M{"Username": req.Username}, &userInfo)
		if err == nil {
			return nil, ss.Err("不可重复注册")
		}

		if err != nil {
			if err != mongo.ErrNoDocuments {
				return nil, err
			}
		}

		secret := NewGoogleAuth().GetSecret()
		qrcode := NewGoogleAuth().GetQrcodeUrl(req.Username, secret)

		DBAdminer := adminpb.DBAdminer{
			ID:           int64(ut2.RandomInt(10000, 99999)),
			Username:     req.Username,
			Password:     "123456",
			GoogleCode:   secret,
			Qrcode:       qrcode,
			IsOpenGoogle: false,
			Status:       adminpb.AdminStatus_Normal,
			IsDelete:     adminpb.IsDelete_NoDelete,
			CreateAt:     mongodb.NewTimeStamp(time.Now()),
			LoginAt:      mongodb.NewTimeStamp(time.Now()),
			GroupId:      req.GroupId,
		}
		if err = CollAdminUser.InsertOne(&DBAdminer); err != nil {
			return nil, err
		}
		return nil,
	*/
}

// 解冻（冻结管理员）
func (AdminInfoRpc) AdminerStatus(ctx context.Context, req *adminpb.AdminerStatusReq) (*pb.Empty, error) {
	return nil, errors.New("Method Unimplemented")
	/*
		updateInfo := make(map[string]interface{}, 0)
		updateInfo["Status"] = req.Status
		if err := CollAdminUser.UpdateId(req.Pid, bson.M{"$set": updateInfo}); err != nil {
			return nil, err
		}
		return nil, nil
	*/
}

// 删除管理员
func (AdminInfoRpc) DelAdminer(ctx context.Context, req *adminpb.AdminPid) (*pb.Empty, error) {
	return nil, errors.New("Method Unimplemented")
	/*
		ss := GetSession(ctx)
		if req.Pid == 3846 {
			return nil, ss.Err("不能删除admin")
		}

		if req.Pid == ss.Pid {
			return nil, ss.Err("不能删除自己")
		}

		updateInfo := make(map[string]interface{}, 0)
		updateInfo["IsDelete"] = adminpb.IsDelete_YesDelete
		if err := CollAdminUser.UpdateId(req.Pid, bson.M{"$set": updateInfo}); err != nil {
			return nil, err
		}
		return nil, nil
	*/
}

// 打开（关闭）谷歌验证码
func (AdminInfoRpc) UpdateOpenGoole(ctx context.Context, req *adminpb.UpdateOpenGooleReq) (*pb.Empty, error) {
	return nil, errors.New("Method Unimplemented")
	/*
		updateInfo := bson.M{}
		updateInfo["IsOpenGoole"] = req.IsOpenGoole

		one := adminpb.DBAdminer{}
		err := CollAdminUser.FindId(req.Pid, &one)
		if err != nil {
			return nil, err
		}

		if one.GoogleCode == "" {
			secret := NewGoogleAuth().GetSecret()
			qrcode := NewGoogleAuth().GetQrcodeUrl("admin", secret)
			updateInfo["GoogleCode"] = secret
			updateInfo["Qrcode"] = qrcode
		}

		if err := CollAdminUser.UpdateId(req.Pid, bson.M{"$set": updateInfo}); err != nil {
			return nil, err
		}

		return &pb.Empty{}, nil
	*/
}

// 修改密码
func (AdminInfoRpc) UpdatePasswd(ctx context.Context, req *adminpb.UpdatePasswdReq) (*pb.Empty, error) {
	ss := GetSession(ctx)

	if req.OldPasswd == "" {
		return nil, ss.Err("请输入旧密码")
	}

	if req.NewPasswd == "" || req.ConfirmPasswd != req.NewPasswd {
		return nil, ss.Err("密码不存在或两次输入不一致")
	}

	userInfo := adminpb.DBAdminer{}
	if err := CollAdminUser.FindId(ss.Pid, &userInfo); err != nil {
		return nil, err
	}

	if userInfo.Password != req.OldPasswd {
		return nil, ss.Err("用户不存在或密码错误")
	}

	updateInfo := make(map[string]interface{}, 0)
	updateInfo["Password"] = req.NewPasswd
	if err := CollAdminUser.UpdateId(ss.Pid, bson.M{"$set": updateInfo}); err != nil {
		return nil, err
	}

	return nil, nil
}

// 退出登录
func (AdminInfoRpc) LoginOut(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	ss := GetSession(ctx)

	userInfo := adminpb.DBAdminer{}
	if err := CollAdminUser.FindId(ss.Pid, &userInfo); err != nil {
		return nil, err
	}

	tokenExp := time.Now().AddDate(0, 0, 0)

	tokenMap.Delete(userInfo.Token)
	tokenMap.Store("", &Token{
		Token:    "",
		ExpireAt: tokenExp,
		Pid:      userInfo.ID,
	})

	return nil, nil
}
