package main

import (
	"game/comm"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/UpdatePermission", "修改权限组", "AdminInfo", updatePermission, updatePermissionParams{
		PermissionId: 1,
	})
}

type updatePermissionParams struct {
	PermissionId int64
	Name         string
	Remark       string
	MenuIds      []int64
}

func updatePermission(ctx *Context, ps updatePermissionParams, ret *comm.Empty) (err error) {
	var user comm.User
	var parentUser comm.User
	err = CollAdminUser.FindOne(bson.M{"Username": ctx.Username}, &user)
	if err != nil {
		return err
	}
	err = CollAdminUser.FindOne(bson.M{"AppID": user.AppID}, &parentUser)
	if err != nil {
		return err
	}

	permissionSet := bson.M{
		"$set": bson.M{
			"Name":    ps.Name,
			"Remark":  ps.Remark,
			"MenuIds": ps.MenuIds,
		},
	}

	err = CollAdminPermission.UpdateOne(bson.M{"_id": ps.PermissionId, "OperatorId": parentUser.Id}, permissionSet)

	return err
}
