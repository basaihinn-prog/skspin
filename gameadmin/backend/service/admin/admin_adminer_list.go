package main

import (
	"errors"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"game/duck/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

// http://192.168.1.14:8080/api/AdminInfo/AdminerList
func init() {
	RegMsgProc("/AdminInfo/AdminerList", "获取后台管理员列表", "AdminInfo", adminerList, adminerListPs{
		PageIndex: 1,
		PageSize:  10,
		AppID:     "",
		UserName:  "",
	})
}

// type httpcheckPs struct {
// }

type adminerListPs struct {
	PageIndex int64
	PageSize  int64
	AppID     string
	UserName  string
}

type adminerListRet struct {
	Count int64
	List  []*comm.User
}

func adminerList(ctx *Context, ps adminerListPs, ret *adminerListRet) (err error) {
	query := bson.M{}
	var user comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}

	_, ok := IsAdminUser(ctx)

	id, err := GetOperatopAppID(ctx)
	if err != nil {
		return err
	}

	id = append(id, user.AppID)
	query["AppID"] = bson.M{"$in": id}
	if ps.AppID != "" {
		query["AppID"] = ps.AppID
	}

	// 如果当前是管理员的情况下 显示除了admin主账号外的用户
	if ok {

		query["$or"] = []bson.M{{"OperatorAdmin": true}, {"AppID": "admin"}}
		query["_id"] = bson.M{"$ne": 1}
	}

	if ps.UserName != "" {
		query["Username"] = bson.M{
			"$regex": ps.UserName,
		}
	}

	filter := mongodb.FindPageOpt{
		Page:     ps.PageIndex,
		PageSize: ps.PageSize,
		Sort:     db.D("CreateAt", -1, "_id", -1, "Status", 1), //可以为空
		Query:    query,
	}
	count, err := CollAdminUser.FindPage(filter, &ret.List)
	if err != nil {
		return
	}
	ret.Count = count

	//if len(ret.List) == 0 {
	//	return
	//}

	//var groups = map[int64]bool{}
	//for _, v := range ret.List {
	//	// groups = append(groups, v.GroupId)
	//	groups[v.GroupId] = true
	//}
	//
	//// fmt.Print(groups)
	//
	//gids := maps.Keys(groups)
	//
	//var auths []*struct {
	//	ID   int64  `protobuf:"varint,1,opt,name=ID,proto3" bson:"_id"`
	//	Name string `protobuf:"bytes,3,opt,name=Name,proto3" ` //权限组名称
	//}
	//CollAdminAuth.FindAllOpt(bson.M{
	//	"_id": bson.M{"$in": gids},
	//}, &auths, options.Find().SetProjection(bson.M{"_id": 1, "Name": 1}))
	//
	//fmt.Println(auths)
	//
	//m := make(map[int64]string, len(auths))
	//for _, a := range auths {
	//	m[a.ID] = a.Name
	//}

	//for _, v := range ret.List {
	//	v.GroupName = m[v.GroupId]
	//}

	return
}
