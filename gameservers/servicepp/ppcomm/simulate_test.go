package ppcomm

import (
	"context"
	"fmt"
	"os"
	"serve/comm/db"
	"testing"
)

func TestSimulate(t *testing.T) {
	mongoaddr := "mongodb://myAdmin:myAdminPassword1@47.237.3.219:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pp_vs20olympx")

	coll := db.Collection("simulate")

	cur, _ := coll.Find(context.TODO(), db.D())
	ridMap := map[string]string{}
	for cur.Next(context.TODO()) {
		var doc SimulateData
		err := cur.Decode(&doc)
		if err != nil {
			panic(err)
		}
		rid := doc.DropPan[0].Str("rid")
		for i := 0; i < len(doc.DropPan); i++ {
			doc.DropPan[i].MKMulFloat("tw", 1)
			if rid != doc.DropPan[i].Str("rid") {
				fmt.Printf("数据处理错误, id:%s\n", doc.Id.Hex())
				os.Exit(1)
			}
		}
		if id, ok := ridMap[rid]; ok {
			fmt.Printf("数据处理错误, Oid:%s, id:%s\n", id, doc.Id.Hex())
			os.Exit(1)
		}
		ridMap[rid] = doc.Id.Hex()
	}
	fmt.Println("====很好，处理没有问题！！！")
}
