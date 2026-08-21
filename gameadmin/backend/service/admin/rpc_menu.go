package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"game/comm"
	"game/duck/mongodb"
	"game/duck/ut2"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/adminpb"
	"github.com/samber/lo"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminMenuRpc struct{}

// 获取菜单列表
func (AdminMenuRpc) MenuList(ctx context.Context, req *pb.Empty) (*adminpb.MenuListResp, error) {
	//uid := common.GetSession(ctx).Pid

	menuList := make([]*adminpb.DBMenu, 0)
	if _, err := CollAdminMenu.FindPage(mongodb.FindPageOpt{
		Page:     1,
		PageSize: 1000,
		Sort:     bson.M{"Sort": 1}, //可以为空
		Query:    bson.M{},
	}, &menuList); err != nil {
		return nil, err
	}
	return &adminpb.MenuListResp{MenuList: menuList}, nil
}

// 添加菜单
func (AdminMenuRpc) AddMenu(ctx context.Context, req *adminpb.DBMenu) (*pb.Empty, error) {
	ss := GetSession(ctx)

	var user comm.User
	err := CollAdminUser.FindOne(bson.M{"_id": ss.Pid}, &user)
	if err != nil || user.AppID != "admin" {
		return nil, ss.Err("权限不足")
	}

	if req.Title == "" {
		return nil, ss.Err("请输入菜单名")
	}

	dbmenu := adminpb.DBMenu{}
	err = CollAdminMenu.FindOne(bson.M{"Title": req.Title}, &dbmenu)
	if err == nil {
		return nil, ss.Err("该菜单已存在")
	}

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	req.ID = int64(ut2.RandomInt(1000, 9999))

	if err = CollAdminMenu.InsertOne(&req); err != nil {
		return nil, err
	}

	//admin权限组添加
	CollAdminPermission.UpdateId(1, bson.M{"$push": bson.M{
		"MenuIds": req.ID,
	}})
	//admin权限组添加
	CollAdminPermission.UpdateId(setting.DefaultPermissionId, bson.M{"$push": bson.M{
		"MenuIds": req.ID,
	}})
	UpdateMenList()
	return nil, nil
}

// 编辑菜单
func (AdminMenuRpc) UpdateMenu(ctx context.Context, req *adminpb.DBMenu) (*pb.Empty, error) {
	ss := GetSession(ctx)

	var user comm.User
	err := CollAdminUser.FindOne(bson.M{"_id": ss.Pid}, &user)
	if err != nil || user.AppID != "admin" {
		return nil, ss.Err("权限不足")
	}

	if req.ID == 0 {
		return nil, ss.Err("请传入菜单id")
	}

	if req.Title == "" {
		return nil, ss.Err("请输入菜单名")
	}

	updateInfo := make(map[string]interface{}, 0)
	updateInfo["Pid"] = req.Pid
	updateInfo["Title"] = req.Title
	updateInfo["Url"] = req.Url
	updateInfo["Sort"] = req.Sort
	updateInfo["Icon"] = req.Icon
	if err := CollAdminMenu.UpdateId(req.ID, bson.M{"$set": updateInfo}); err != nil {
		return nil, err
	}
	UpdateMenList()
	return nil, nil
}

// 删除菜单
func (AdminMenuRpc) DelMenu(ctx context.Context, req *adminpb.DelMenuReq) (*pb.Empty, error) {
	ss := GetSession(ctx)
	var user comm.User
	err := CollAdminUser.FindOne(bson.M{"_id": ss.Pid}, &user)
	if err != nil || user.AppID != "admin" {
		return nil, ss.Err("权限不足")
	}
	if req.ID == 0 {
		return nil, ss.Err("请传入菜单id")
	}

	if err := CollAdminMenu.DeleteId(req.ID); err != nil {
		return nil, err
	}
	CollAdminPermission.Update(
		bson.M{
			"$in": bson.M{"MenuIds": req.ID},
		},
		bson.M{
			"$pull": bson.M{"MenuIds": req.ID},
		})
	UpdateMenList()
	return nil, nil
}

// Md5 md5()
func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func NormalMenu(ids []int64) []int64 {
	cnt := 0
	for i := 0; i < 5; i++ {
		cnt = 0
		for j := range ids {
			menu, ok := lo.Find(MenuList, func(item *comm.PerMenu) bool {
				return item.ID == ids[j]
			})
			if !ok {
				continue
			}
			if menu.Pid != 0 {
				if slices.Contains(ids, menu.Pid) {
					continue
				}
				cnt++
				ids = append(ids, menu.Pid)
			}
		}
		if cnt == 0 {
			break
		}
	}
	return ids
}
