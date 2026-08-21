package main

import (
	"context"
	"errors"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/adminpb"
	"reflect"
	"strings"
	"unicode"
)

type AdminGroupRpc struct{}

// 获取权限组列表
func (AdminGroupRpc) GroupList(ctx context.Context, req *adminpb.PageListReq) (*adminpb.GroupListResp, error) {
	return nil, errors.New("Method Unimplemented")
	/*
		//uid := GetSession(ctx).Pid
		var list []*adminpb.DBAuth
		filter := mongodb.FindPageOpt{
			Page:     req.PageIndex,
			PageSize: req.PageSize,
			Sort:     bson.M{"_id": -1}, //可以为空
			Query:    bson.M{},
		}
		count, err := CollAdminAuth.FindPage(filter, &list)
		if err != nil {
			return nil, err
		}
		return &adminpb.GroupListResp{Count: count, List: list}, nil
	*/
}

// 添加权限组
func (AdminGroupRpc) AddGroup(ctx context.Context, req *adminpb.AddGroupReq) (*pb.Empty, error) {
	return nil, errors.New("Method Unimplemented")
	/*
		ss := GetSession(ctx)

		if req.Name == "" {
			return nil, ss.Err("请输入权限组名称")
		}

		DBAuth := adminpb.DBAuth{
			ID:         int64(ut2.RandomInt(1000000, 9999000)),
			Name:       req.Name,
			Remark:     req.Remark,
			CreateTime: mongodb.NewTimeStamp(time.Now()),
		}
		if err := CollAdminAuth.InsertOne(&DBAuth); err != nil {
			return nil, err
		}
		return nil, nil
	*/
}

// 根据权限ID查看权限详情
func (AdminGroupRpc) GetPowerDetailsById(ctx context.Context, req *adminpb.GroupIdReq) (*adminpb.PowerListResp, error) {
	return nil, errors.New("Method Unimplemented")
	/*
		ss := GetSession(ctx)

		if req.Gid == 0 {
			return nil, ss.Err("请传入权限ID")
		}

		group := adminpb.DBAuth{}
		if err := CollAdminAuth.FindId(req.Gid, &group); err != nil {
			return nil, err
		}

		menuList := make([]*adminpb.DBMenu, 0)

		if err := CollAdminMenu.FindAll(bson.M{}, &menuList); err != nil {
			return nil, err
		}

		roleMenuList := strings.Split(group.MenuIds, ",")
		powerList := make([]int64, 0)
		if !Empty(roleMenuList) {
			for _, v := range roleMenuList {
				iv, _ := strconv.Atoi(v)
				powerList = append(powerList, int64(iv))
			}
		}

		return &adminpb.PowerListResp{MenuList: menuList, Gid: req.Gid, CheckedArray: powerList}, nil
	*/
}

// 修改权限组菜单
func (AdminGroupRpc) EditPowerDetails(ctx context.Context, req *adminpb.EditPowerReq) (*pb.Empty, error) {
	return nil, errors.New("Method Unimplemented")
	/*
		ss := GetSession(ctx)

		if req.Gid == 0 {
			return nil, ss.Err("请传入权限ID")
		}

		var menus = ""
		for _, v := range req.PowerList {
			menus += fmt.Sprintf("%v", v) + ","
		}
		updateInfo := make(map[string]interface{}, 0)
		updateInfo["MenuIds"] = Rtrim(menus, ",")
		if err := CollAdminAuth.UpdateId(req.Gid, bson.M{"$set": updateInfo}); err != nil {
			return nil, err
		}
		return nil, nil
	*/
}

func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}

func Rtrim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimRightFunc(str, unicode.IsSpace)
	}
	return strings.TrimRight(str, characterMask[0])
}
