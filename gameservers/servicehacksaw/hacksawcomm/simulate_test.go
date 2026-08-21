package hacksawcomm

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"serve/comm/db"
	"serve/comm/ut"
	"testing"
)

func TestSimulate(t *testing.T) {

	mongoaddr := "mongodb://myAdmin:myAdminPassword1@18.61.185.51:27017/?authSource=admin"

	gamesDb := []string{"hacksaw_1620"}
	for _, gamename := range gamesDb {
		db.DialToMongo(mongoaddr, gamename)

		coll := db.Collection2(gamename, "simulate")

		//cur, _ := coll.Find(context.TODO(), db.D())

		filter := bson.M{"selected": true}
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var simulateData []SimulateData
		if err := cursor.All(context.TODO(), &simulateData); err != nil {
			log.Fatal(err)
		}
		if len(simulateData) == 0 {
			fmt.Printf("game:%v is jump\n", gamename)
			continue
		}
		for i := range simulateData {
			spinResult := simulateData[i].DropPan["spinResult"].(Variables)
			totalWin := spinResult.Float("totalWin")
			if totalWin != 0 {
				boardDisplayResult := spinResult["boardDisplayResult"].(Variables)
				originC := boardDisplayResult.Float("displayBet")
				simulateData[i].Times = ut.Round6(totalWin / originC)
			}
		}
		//批量更新
		var updates []mongo.WriteModel
		for _, data := range simulateData {
			update := mongo.NewUpdateOneModel()
			update.SetFilter(bson.M{"_id": data.Id})
			update.SetUpdate(bson.M{"$set": bson.M{"droppan": data.DropPan, "times": data.Times}})
			updates = append(updates, update)
		}

		_, err = coll.BulkWrite(context.TODO(), updates)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game:%v flash done\n", gamename)
	}

}
