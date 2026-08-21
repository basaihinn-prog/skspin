package internal

import (
	"context"
	"fmt"
	"log"
	"serve/comm/db"
	"serve/servicepp/ppcomm"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 救所有等于drop pan =2 的 数据
func TestSimulate(t *testing.T) {

	mongoaddr := "mongodb://myAdmin:myAdminPassword1@54.251.234.111:27017/?authSource=admin"

	//gamesDb := []string{"pp_vs10fangfree", "pp_vs10jokerhot", "pp_vs10noodles", "pp_vs20drgbless", "vs20mustanggld2", "pp_vs40wildrun", "pp_vs50juicyfr", "pp_vs5jjwild", "pp_vs5joker", "pp_vswaysfreezet", "pp_vswayssevenc"}
	//gamesDb := []string{"pp_vs10fangfree", "pp_vs10jokerhot", "pp_vs10noodles", "pp_vs20drgbless", "vs20mustanggld2", "pp_vs40wildrun", "pp_vs50juicyfr", "pp_vs5jjwild", "pp_vs5joker", "pp_vswaysfreezet", "pp_vswayssevenc"}
	gamesDb := []string{"pp_vs10fangfree", "pp_vs10jokerhot", "pp_vs10noodles", "pp_vs20drgbless", "vs20mustanggld2", "pp_vs40wildrun", "pp_vs50juicyfr", "pp_vs5jjwild", "pp_vs5joker", "pp_vswaysfreezet", "pp_vswayssevenc"}
	for _, gamename := range gamesDb {
		db.DialToMongo(mongoaddr, gamename)

		coll := db.Collection2(gamename, "simulate")

		//cur, _ := coll.Find(context.TODO(), db.D())

		filter := bson.M{"droppan": bson.M{"$size": 2}}
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var simulateData []ppcomm.SimulateData
		if err := cursor.All(context.TODO(), &simulateData); err != nil {
			log.Fatal(err)
		}
		if len(simulateData) == 0 {
			fmt.Printf("game:%v is jump\n", gamename)
			continue
		}
		for i := range simulateData {
			if len(simulateData[i].DropPan) >= 2 {
				simulateData[i].DropPan = simulateData[i].DropPan[1:]
			}
		}
		//批量更新
		var updates []mongo.WriteModel
		for _, data := range simulateData {
			update := mongo.NewUpdateOneModel()
			update.SetFilter(bson.M{"_id": data.Id})
			update.SetUpdate(bson.M{"$set": bson.M{"droppan": data.DropPan}})
			updates = append(updates, update)
		}

		_, err = coll.BulkWrite(context.TODO(), updates)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done\n", gamename)
	}

}

// 救freeGame 起始不扣钱的
func TestSimulate2(t *testing.T) {
	mongoaddr := "mongodb://myAdmin:myAdminPassword1@54.251.234.111:27017/?authSource=admin"

	//gamesDb := []string{"pp_vs10fangfree", "pp_vs10jokerhot", "pp_vs10noodles", "pp_vs20drgbless", "vs20mustanggld2", "pp_vs40wildrun", "pp_vs50juicyfr", "pp_vs5jjwild", "pp_vs5joker", "pp_vswaysfreezet", "pp_vswayssevenc"}
	gamesDb := []string{"pp_vs12bbb"}
	//gamesDb := []string{"pp_vs20mustanggld2", "pp_vs40wildrun", "pp_vswaysfreezet", "pp_vswayssevenc"}
	for _, gamename := range gamesDb {
		db.DialToMongo(mongoaddr, gamename)

		coll := db.Collection2(gamename, "simulate")

		//cur, _ := coll.Find(context.TODO(), db.D())

		filter := bson.M{"hasgame": true, "type": 0}
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var simulateData []ppcomm.SimulateData
		if err := cursor.All(context.TODO(), &simulateData); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done all data count:%v \n", gamename, len(simulateData))

		if len(simulateData) == 0 {
			fmt.Printf("game:%v is jump\n", gamename)
			continue
		}
		//处理所有数据
		for i := range simulateData {
			//if len(simulateData[i].DropPan) >= 2 ||
			//	simulateData[i].DropPan[0].Int("fs") == 0 && simulateData[i].DropPan[1].Str("fs") != simulateData[i].DropPan[0].Str("fs") && simulateData[i].DropPan[1].Int("fs") > 0 && simulateData[i].DropPan[1].Int("fsmax") > 0 {
			//fmt.Println(simulateData[i].DropPan[len(simulateData[i].DropPan)-1].Int("fs_total"))
			if simulateData[i].DropPan[len(simulateData[i].DropPan)-1].Int("fs_total") != len(simulateData[i].DropPan)-1 {
				//切断
				simulateData[i].DropPan = simulateData[i].DropPan[1:]
				//重新计算 gid
				params := strings.Split(simulateData[i].DropPan[0].Str("gid"), "_")
				for j := 0; j < len(simulateData[i].DropPan); j++ {
					simulateData[i].DropPan[j]["gid"] = fmt.Sprintf("%s_%d_%d", params[0], j+1, len(simulateData[i].DropPan))
				}
			}
		}
		//批量更新
		var updates []mongo.WriteModel
		for _, data := range simulateData {
			update := mongo.NewUpdateOneModel()
			update.SetFilter(bson.M{"_id": data.Id})
			update.SetUpdate(bson.M{"$set": bson.M{"droppan": data.DropPan}})
			updates = append(updates, update)
		}

		_, err = coll.BulkWrite(context.TODO(), updates)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done fix data count: %v \n", gamename, len(updates)) //15887
	}

}

// 救freeGame 粘连spin + free game 用 st: rect 认为是一次正常的连转进free game
func TestSimulate3(t *testing.T) {
	mongoaddr := "mongodb://myAdmin:myAdminPassword1@54.251.234.111:27017/?authSource=admin"

	gamesDb := []string{"pp_vs10txbigbass", "pp_vs20drgbless"}
	for _, gamename := range gamesDb {
		db.DialToMongo(mongoaddr, gamename)

		coll := db.Collection2(gamename, "simulate")

		//cur, _ := coll.Find(context.TODO(), db.D())

		filter := bson.M{"hasgame": true, "type": 0}
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var simulateData []ppcomm.SimulateData
		if err := cursor.All(context.TODO(), &simulateData); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done all data count:%v \n", gamename, len(simulateData))

		if len(simulateData) == 0 {
			fmt.Printf("game:%v is jump\n", gamename)
			continue
		}
		//处理所有数据
		for i := range simulateData {
			fsInd := -1
			for j := range simulateData[i].DropPan {
				if simulateData[i].DropPan[j].Int("fs") == 1 {
					fsInd = j
				}
			}
			//不是连转进的fg
			if fsInd != -1 && fsInd != 0 && simulateData[i].DropPan[fsInd-1].Str("st") != "rect" {
				//切断
				simulateData[i].DropPan = simulateData[i].DropPan[fsInd:]
				//重新计算 gid
				params := strings.Split(simulateData[i].DropPan[0].Str("gid"), "_")
				for j := 0; j < len(simulateData[i].DropPan); j++ {
					simulateData[i].DropPan[j]["gid"] = fmt.Sprintf("%s_%d_%d", params[0], j+1, len(simulateData[i].DropPan))
				}
			}
		}
		//批量更新
		var updates []mongo.WriteModel
		for _, data := range simulateData {
			update := mongo.NewUpdateOneModel()
			update.SetFilter(bson.M{"_id": data.Id})
			update.SetUpdate(bson.M{"$set": bson.M{"droppan": data.DropPan}})
			updates = append(updates, update)
		}

		_, err = coll.BulkWrite(context.TODO(), updates)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done fix data count: %v \n", gamename, len(updates)) //15887
	}

}

// 救freeGame 粘连spin + free game 用 rs_t: ""
func TestSimulate4(t *testing.T) {
	mongoaddr := "mongodb://myAdmin:myAdminPassword1@54.251.234.111:27017/?authSource=admin"

	//gamesDb := []string{"pp_vswaysfirewmw"}
	gamesDb := []string{"pp_vs20clreacts", "pp_vswaysfirewmw", "pp_vswaysmfreya", "pp_vs20procountx", "pp_vs20aztecgates"}
	for _, gamename := range gamesDb {
		db.DialToMongo(mongoaddr, gamename)

		coll := db.Collection2(gamename, "simulate")

		//cur, _ := coll.Find(context.TODO(), db.D())

		filter := bson.M{"hasgame": true, "type": 0}
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var simulateData []ppcomm.SimulateData
		if err := cursor.All(context.TODO(), &simulateData); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done all data count:%v \n", gamename, len(simulateData))

		if len(simulateData) == 0 {
			fmt.Printf("game:%v is jump\n", gamename)
			continue
		}
		//处理所有数据
		for i := range simulateData {
			fsInd := -1
			for j := range simulateData[i].DropPan {
				if simulateData[i].DropPan[j].Int("fs") == 1 {
					fsInd = j
				}
			}
			//不是连转进的fg
			if fsInd != -1 && fsInd != 0 && simulateData[i].DropPan[fsInd-1].Str("rs_t") == "" {
				//切断
				simulateData[i].DropPan = simulateData[i].DropPan[fsInd:]
				//重新计算 gid
				params := strings.Split(simulateData[i].DropPan[0].Str("gid"), "_")
				for j := 0; j < len(simulateData[i].DropPan); j++ {
					simulateData[i].DropPan[j]["gid"] = fmt.Sprintf("%s_%d_%d", params[0], j+1, len(simulateData[i].DropPan))
				}
			}
		}
		//批量更新
		var updates []mongo.WriteModel
		for _, data := range simulateData {
			update := mongo.NewUpdateOneModel()
			update.SetFilter(bson.M{"_id": data.Id})
			update.SetUpdate(bson.M{"$set": bson.M{"droppan": data.DropPan}})
			updates = append(updates, update)
		}

		_, err = coll.BulkWrite(context.TODO(), updates)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done fix data count: %v \n", gamename, len(updates)) //15887
	}

}

// 大鲈鱼 所有 drop pan = 2的
func TestSimulate5(t *testing.T) {
	mongoaddr := "mongodb://myAdmin:myAdminPassword1@54.251.234.111:27017/?authSource=admin"

	gamesDb := []string{"pp_vs40bigjuan"}
	//gamesDb := []string{"pp_vs10bhallbnza2", "pp_vs10txbigbass", "pp_vs12bbbxmas", "pp_vs10bbsplxmas", "pp_vs12bbb"}
	for _, gamename := range gamesDb {
		db.DialToMongo(mongoaddr, gamename)

		coll := db.Collection2(gamename, "simulate")

		//cur, _ := coll.Find(context.TODO(), db.D())

		filter := bson.M{"droppan": bson.M{"$size": 1}, "type": 1}
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var simulateData []ppcomm.SimulateData
		if err := cursor.All(context.TODO(), &simulateData); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done all data count:%v \n", gamename, len(simulateData))
		if len(simulateData) == 0 {
			continue
		}
		////处理所有数据
		//for i := range simulateData {
		//
		//}
		//批量更新
		var updates []mongo.WriteModel
		for _, data := range simulateData {

			update := mongo.NewDeleteOneModel()
			update.SetFilter(bson.M{"_id": data.Id})
			//update.SetUpdate(bson.M{"$set": bson.M{"droppan": data.DropPan}})
			updates = append(updates, update)
		}

		_, err = coll.BulkWrite(context.TODO(), updates)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done del data count: %v \n", gamename, len(updates)) //15887
	}

}

// 大鲈鱼 删除所有 drop pan = 1的 并且 rs_c= 1
func TestSimulate6(t *testing.T) {
	mongoaddr := "mongodb://myAdmin:myAdminPassword1@54.251.234.111:27017/?authSource=admin"

	//gamesDb := []string{"pp_vs10fangfree"}
	//gamesDb := []string{"pp_vs10bhallbnza2", "pp_vs10txbigbass", "pp_vs12bbbxmas", "pp_vs10bbsplxmas", "vs12bbb"}
	gamesDb := []string{"pp_vs10bhallbnza2", "pp_vs10txbigbass", "pp_vs12bbbxmas", "pp_vs10bbsplxmas", "pp_vs12bbb"}
	for _, gamename := range gamesDb {
		db.DialToMongo(mongoaddr, gamename)

		coll := db.Collection2(gamename, "simulate")

		//cur, _ := coll.Find(context.TODO(), db.D())

		filter := bson.M{"droppan": bson.M{"$size": 1}}
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var simulateData []ppcomm.SimulateData
		var simulateData2 []ppcomm.SimulateData
		if err := cursor.All(context.TODO(), &simulateData); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done all data count:%v \n", gamename, len(simulateData))
		////处理所有数据
		for i := range simulateData {
			if simulateData[i].DropPan[0].Int("rs_c") == 1 {
				simulateData2 = append(simulateData2, simulateData[i])
			}
		}
		if len(simulateData2) == 0 {
			continue
		}
		//批量更新
		var updates []mongo.WriteModel
		for _, data := range simulateData2 {
			update := mongo.NewDeleteOneModel()
			update.SetFilter(bson.M{"_id": data.Id})
			//update.SetUpdate(bson.M{"$set": bson.M{"droppan": data.DropPan}})
			updates = append(updates, update)
		}

		_, err = coll.BulkWrite(context.TODO(), updates)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done del data count: %v \n", gamename, len(updates)) //15887
	}

}

// 大鲈鱼 删除所有 drop pan = 1的 并且 rs_c > 0
func TestSimulate7(t *testing.T) {
	mongoaddr := "mongodb://myAdmin:myAdminPassword1@54.251.234.111:27017/?authSource=admin"

	//gamesDb := []string{"pp_vs10fangfree"}
	gamesDb := []string{"pp_vs10bhallbnza2", "pp_vs10txbigbass", "pp_vs12bbbxmas", "pp_vs10bbsplxmas", "pp_vs12bbb"}
	//gamesDb := []string{"pp_vs12bbbxmas", "pp_vs10bbsplxmas", "vs12bbb"}
	for _, gamename := range gamesDb {
		db.DialToMongo(mongoaddr, gamename)

		coll := db.Collection2(gamename, "simulate")

		//cur, _ := coll.Find(context.TODO(), db.D())

		filter := bson.M{"droppan": bson.M{"$size": 1}}
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var simulateData []ppcomm.SimulateData
		var simulateData2 []ppcomm.SimulateData
		if err := cursor.All(context.TODO(), &simulateData); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done all data count:%v \n", gamename, len(simulateData))
		////处理所有数据
		for i := range simulateData {
			if simulateData[i].DropPan[0].Int("rs_c") > 1 {
				simulateData2 = append(simulateData2, simulateData[i])
			}
		}
		if len(simulateData2) == 0 {
			continue
		}
		//批量更新
		var updates []mongo.WriteModel
		for _, data := range simulateData2 {
			update := mongo.NewDeleteOneModel()
			update.SetFilter(bson.M{"_id": data.Id})
			//update.SetUpdate(bson.M{"$set": bson.M{"droppan": data.DropPan}})
			updates = append(updates, update)
		}

		_, err = coll.BulkWrite(context.TODO(), updates)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done del data count: %v \n", gamename, len(updates)) //15887
	}

}

// pp_vs25bkofkngdm 删除所有 type
func TestSimulate8(t *testing.T) {
	mongoaddr := "mongodb://myAdmin:myAdminPassword1@54.251.234.111:27017/?authSource=admin"

	//gamesDb := []string{"pp_vs10fangfree"}
	//gamesDb := []string{"pp_vs10bhallbnza2", "pp_vs10txbigbass", "pp_vs12bbbxmas", "pp_vs10bbsplxmas", "vs12bbb"}
	gamesDb := []string{"pp_vs25bkofkngdm"}
	for _, gamename := range gamesDb {
		db.DialToMongo(mongoaddr, gamename)

		coll := db.Collection2(gamename, "simulate")
		findOptions := options.Find().SetLimit(60000)
		filter := bson.D{{"bucketid", 0}}
		cursor, err := coll.Find(context.TODO(), filter, findOptions)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var idsToDelete []primitive.ObjectID
		for cursor.Next(context.TODO()) {
			var result ppcomm.SimulateData
			if err := cursor.Decode(&result); err != nil {
				log.Fatal(err)
			}
			idsToDelete = append(idsToDelete, result.Id)
		}

		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}

		// 使用 $in 操作符构建删除过滤器
		deleteFilter := bson.D{{"_id", bson.D{{"$in", idsToDelete}}}}

		// 执行 DeleteMany 操作
		deleteResult, err := coll.DeleteMany(context.TODO(), deleteFilter)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted %v documents\n", deleteResult.DeletedCount)
	}

}
