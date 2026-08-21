package main

import (
	"serve/comm/db"
	"serve/comm/lazy"
	"serve/comm/mq"
	"serve/servicepp/ppcomm"
	"serve/servicepp/robot_dispatch_center/internal"

	"github.com/samber/lo"
)

func main() {
	lazy.Init("robots_center")
	mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
	db.DialToMongo(mgoaddr, lazy.ServiceName)

	addr := lo.Must(lazy.RouteFile.Get("proxy.mq"))
	mqconn := lo.Must(mq.ConnectServerMust(addr))
	ppcomm.RegistRpcToMQ(mqconn)
	//go internal.RobotsDispatch()
	go internal.TasksDispatch()
	//go func() {
	//	for ; ; time.Sleep(5 * time.Minute) {
	//
	//		gamesDb := []string{"pp_vs40bigjuan"}
	//		//gamesDb := []string{"pp_vs10bhallbnza2", "pp_vs10txbigbass", "pp_vs12bbbxmas", "pp_vs10bbsplxmas", "pp_vs12bbb"}
	//		for _, gamename := range gamesDb {
	//			coll := db.Collection2(gamename, "simulate")
	//			//cur, _ := coll.Find(context.TODO(), db.D())
	//
	//			filter := bson.M{"droppan": bson.M{"$size": 1}, "type": 1}
	//			cursor, err := coll.Find(context.TODO(), filter)
	//			if err != nil {
	//				log.Fatal(err)
	//			}
	//			defer cursor.Close(context.TODO())
	//
	//			var simulateData []ppcomm.SimulateData
	//			if err := cursor.All(context.TODO(), &simulateData); err != nil {
	//				log.Fatal(err)
	//			}
	//			fmt.Printf("game:%v flash done all data count:%v \n", gamename, len(simulateData))
	//			if len(simulateData) == 0 {
	//				continue
	//			}
	//			////处理所有数据
	//			//for i := range simulateData {
	//			//
	//			//}
	//			//批量更新
	//			var updates []mongo.WriteModel
	//			for _, data := range simulateData {
	//
	//				update := mongo.NewDeleteOneModel()
	//				update.SetFilter(bson.M{"_id": data.Id})
	//				//update.SetUpdate(bson.M{"$set": bson.M{"droppan": data.DropPan}})
	//				updates = append(updates, update)
	//			}
	//
	//			_, err = coll.BulkWrite(context.TODO(), updates)
	//			if err != nil {
	//				log.Fatal(err)
	//			}
	//			fmt.Printf("game:%v flash done del data count: %v \n", gamename, len(updates))
	//		}
	//	}
	//}()
	lazy.Serve()
}
