package main

import (
	"context"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/duck/lang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	RegMsgProc("/AdminInfo/UpdateOperatorStatus", "修改运营商", "AdminInfo", updateOperatorStatus, UpdateOperatorStatusParams{
		OperatorId: 0,
		Status:     0,
	})
}

type UpdateOperatorStatusParams struct {
	OperatorId int64
	Status     int
}

func updateOperatorStatus(ctx *Context, ps UpdateOperatorStatusParams, ret *comm.Empty) (err error) {
	_, ok := IsAdminUser(ctx)
	if !ok {
		return errors.New(lang.GetLang(ctx.Lang, "权限不足"))
	}
	err = CollAdminOperator.UpdateId(ps.OperatorId, bson.M{
		"$set": bson.M{
			"Status": ps.Status,
		},
	})

	return err
}

// 自动同步下个月的费率脚本
func SyncPlantRate() {
	fmt.Println("______SyncPlantRate  begin_____________")
	var results []comm.Operator_V2
	err := CollAdminOperator.FindAll(bson.M{"NextRate": bson.M{"$gt": 0}}, &results)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	models := []mongo.WriteModel{}

	if len(results) > 0 {

		for _, doc := range results {

			model := mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": doc.Id}).SetUpdate(
				bson.M{
					"$set": bson.M{
						"PlatformPay": float32(doc.NextRate),
					},
				},
			)
			models = append(models, model)

		}
		_, err = db.Collection2("GameAdmin", "AdminOperator").BulkWrite(context.TODO(), models, options.BulkWrite().SetOrdered(true))
		if err != nil {
			fmt.Println("批量处理错误" + err.Error())
			return
		}
	}
	fmt.Println("______SyncPlantRate  End_____________")
}
