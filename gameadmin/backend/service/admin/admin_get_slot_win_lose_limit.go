package main

import (
	"errors"
	"game/comm"
	"game/duck/lang"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	RegMsgProc("/AdminInfo/SlotWinLoseLimitList", "玩家slots净输赢记录(列表)", "AdminInfo", slotWinLoseLimitList, slotWinLoseLimitListParams{
		OperatorId:      0,
		StartTime:       0,
		EndTime:         0,
		EnsureName:      "",
		EnsureStatus:    1,
		EnsureStartTime: 0,
		EnsureEndTime:   0,
		pageSetting: pageSetting{
			PageSize:   10,
			PageNumber: 1,
		},
	})

	RegMsgProc("/AdminInfo/doEnsureSlotWinLoseLimit", "审核玩家slots净输赢的记录", "AdminInfo", doEnsureSlotWinLoseLimit, doEnsureSlotWinLoseLimitParams{
		ID:     "",
		Remark: "",
	})
}

type slotWinLoseLimitListParams struct {
	OperatorId      int64  //运营商
	StartTime       int64  //记录时间
	EndTime         int64  //记录时间
	Pid             int64  //玩家id
	EnsureName      string //审核人
	EnsureStatus    int64  //0:未审核,1已审核
	EnsureStartTime int64  //审核时间
	EnsureEndTime   int64  //审核时间
	pageSetting
}

type slotWinLoseLimitListResults struct {
	List *[]comm.SlotWinLoseLimit
	pageSetting
}

func slotWinLoseLimitList(ctx *Context, ps slotWinLoseLimitListParams, ret *slotWinLoseLimitListResults) (err error) {

	ids, _ := GetOperatopAppID(ctx)
	if ps.PageNumber == 0 {
		ps.PageNumber = 1
		ret.PageNumber = 1
	} else {
		ret.PageNumber = ps.PageNumber
	}
	if ps.PageSize == 0 {
		ps.PageSize = 50
		ret.PageSize = 50
	} else {
		ret.PageSize = ps.PageSize
	}
	filter := bson.M{}
	err = filter4WinLoseLimitList(filter, ctx.Lang, ps, ids)
	if err != nil {
		return
	}
	findOptions := options.Find()
	skip := (ps.PageNumber - 1) * ps.PageSize
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.M{"Time": -1})
	findOptions.SetLimit(int64(ps.PageSize))

	coll := DB.Collection("WinLoseLimit")
	tmp := make([]comm.SlotWinLoseLimit, ps.PageSize)
	err = coll.FindAllOpt(filter, &tmp, findOptions)
	ret.List = &tmp
	count, _ := coll.CountDocuments(filter)
	ret.Count = int(count)
	return err
}

type doEnsureSlotWinLoseLimitParams struct {
	ID string
	// EnsureStatus int    //0:未审核,1已审核
	Remark string //审核备注
}

func doEnsureSlotWinLoseLimit(ctx *Context, ps doEnsureSlotWinLoseLimitParams, ret *comm.Empty) (err error) {
	id, err := primitive.ObjectIDFromHex(ps.ID)
	if err != nil {
		return
	}
	// 已审核不能再次审核
	ensureFilter := bson.M{"_id": id, "Remark": bson.M{"$exists": false}, "EnsureStatus": bson.M{"$exists": false}}
	coll := DB.Collection("WinLoseLimit")
	tmp := &comm.SlotWinLoseLimit{}
	err = coll.FindOne(ensureFilter, tmp)
	if err != nil {
		return
	}
	filter := bson.M{"_id": id}
	update :=
		bson.M{"$set": bson.M{"EnsureStatus": 1,
			"Remark":       ps.Remark,
			"EnsureOpPid":  ctx.PID,
			"EnsureOpName": ctx.Username,
			"EnsureTime":   time.Now()},
		}

	err = coll.UpdateOne(filter, update)
	if err != nil {
		return
	}
	return
}

func filter4WinLoseLimitList(match bson.M, Lang string, ps slotWinLoseLimitListParams, AppIDs []string) error {
	if ps.EnsureStatus == 0 {
		match["EnsureStatus"] = bson.M{"$ne": 1}
	} else if ps.EnsureStatus == 1 {
		match["EnsureStatus"] = 1
	} else {
		return errors.New(lang.GetLang(Lang, "参数错误"))
	}
	if ps.OperatorId != 0 {
		coll := DB.Collection("AdminOperator")
		op := &comm.Operator{}
		coll.FindId(ps.OperatorId, op)
		AppIDs = append(AppIDs, op.Name)
	}

	match["AppID"] = bson.M{
		"$in": AppIDs,
	}
	if ps.StartTime != 0 || ps.EndTime != 0 {
		timeTemp := bson.M{}
		match["Time"] = timeTemp
		if ps.StartTime != 0 {
			timeTemp["$gte"] = time.Unix(ps.StartTime, 0)
		}
		if ps.EndTime != 0 {
			timeTemp["$lte"] = time.Unix(ps.EndTime, 0)
		}
	}
	if ps.EnsureStartTime != 0 || ps.EnsureEndTime != 0 {
		timeTemp := bson.M{}
		match["EnsureTime"] = timeTemp
		if ps.EnsureStartTime != 0 {
			timeTemp["$gte"] = time.Unix(ps.EnsureStartTime, 0)
		}
		if ps.EnsureEndTime != 0 {
			timeTemp["$lte"] = time.Unix(ps.EnsureEndTime, 0)
		}
	}
	if len(ps.EnsureName) != 0 {
		match["EnsureOpName"] = ps.EnsureName
	}
	if ps.Pid != 0 {
		match["ID"] = ps.Pid
	}
	return nil
}
