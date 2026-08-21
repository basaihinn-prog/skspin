package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"
	"game/duck/mongodb"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	RegMsgProc("/AdminInfo/GetOperatorList", "获取营运商列表", "AdminInfo", getOperatorList, getOperatorListParams{
		PageIndex: 1,
		PageSize:  10,
	})
}

type getOperatorListParams struct {
	PageIndex int64
	PageSize  int64
}

type getOperatorListResult struct {
	List     []*comm.Operator
	AllCount int64
}

func getOperatorList(ctx *Context, ps getOperatorListParams, ret *getOperatorListResult) (err error) {
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "参数错误"))
	}
	var list []*comm.Operator
	//err = CollAdminOperator.FindAll(bson.M{}, &list)
	//ret.List = list

	filter := mongodb.FindPageOpt{
		Page:     ps.PageIndex,
		PageSize: ps.PageSize,
		Sort:     bson.M{"_id": -1},
		Query:    bson.M{},
	}

	hasChildren := func(id int64) bool {
		for i := 0; i < len(list); i++ {
			_, ok = lo.Find(MenuList, func(item *comm.PerMenu) bool {
				return item.Pid == id
			})
			if ok {
				return true
			}
		}
		return false
	}

	count, err := CollAdminOperator.FindPage(filter, &list)
	for i := range list {
		for j := len(list[i].MenuIds) - 1; j >= 0; j-- {
			menu, ok := lo.Find(MenuList, func(item *comm.PerMenu) bool {
				return item.ID == list[i].MenuIds[j]
			})
			if !ok {
				continue
			}
			if menu.Pid == 0 && hasChildren(list[i].MenuIds[j]) {
				list[i].MenuIds = append(list[i].MenuIds[:j], list[i].MenuIds[j+1:]...)
				continue
			}
			_, ok = lo.Find(MenuList, func(item *comm.PerMenu) bool {
				return item.Pid == menu.ID
			})
			if ok {
				list[i].MenuIds = append(list[i].MenuIds[:j], list[i].MenuIds[j+1:]...)
			}
		}
	}
	ret.AllCount = count
	ret.List = list

	return err
}
