package dealData

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"serve/comm/db"
	"serve/comm/ut"
	"serve/service/pg_89/internal/gendata"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestSimulate(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_89")
	for i := 0; i < 6; i++ {
		num := i
		go dealNormal(num)
		go dealGame(num)
	}
	select {}
}

func TestCaclBucketId(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_89")
	coll := db.Collection("simulate")

	cur, _ := coll.Find(context.TODO(), db.D(), options.Find().SetProjection(db.D("times", 1, "hasgame", 1, "type", 1)))

	models := []mongo.WriteModel{}
	for cur.Next(context.TODO()) {
		var doc struct {
			ID      primitive.ObjectID `bson:"_id"`
			Times   float64            `bson:"times"`
			HasGame bool               `bson:"hasgame"`
			Type    db.BoundType       `bson:"type"`
		}
		assert.Nil(t, cur.Decode(&doc))

		bucketId := gendata.GBuckets.GetBucket(doc.Times, doc.HasGame, doc.Type)

		model := mongo.NewUpdateOneModel().SetFilter(db.ID(doc.ID)).SetUpdate(db.D("$set", db.D("bucketid", bucketId, "selected", true)))
		models = append(models, model)
	}

	coll.BulkWrite(context.TODO(), models, options.BulkWrite().SetOrdered(false))
}

func dealNormal(i int) {
	collName := fmt.Sprintf("pgSpinData_%d", i)
	coll := db.Collection(collName)
	newColl := db.Collection("simulate")
	var lastId primitive.ObjectID
	db.MiscGet(collName, &lastId)
	match := bson.M{}
	if lastId != primitive.NilObjectID {
		match["_id"] = bson.M{"$gte": lastId}
	}
	cursor, _ := coll.Find(context.TODO(), match)
	defer cursor.Close(context.TODO())

	pans := []bson.M{}
	for cursor.Next(context.TODO()) {
		var pan map[string]any
		err := cursor.Decode(&pan)
		if err != nil {
			fmt.Println("err:", err,
				"_id:", pan["_id"],
				"sid:", pan["sid"],
				"psid:", pan["psid"])
			break
		}
		_id := pan["_id"]
		pan["orignid"] = fmt.Sprintf("%s_%s", pan["_id"], collName)
		delete(pan, "_id")
		if pan["sid"].(string) == pan["psid"].(string) && len(pans) > 0 {
			if pans[0]["sid"].(string) == pans[0]["psid"].(string) { // 第一盘是否有效
				objectId := primitive.NewObjectID()
				hasGame := false
				for j := 0; j < len(pans); j++ {
					pans[j]["psid"] = objectId.Hex()
					pans[j]["sid"] = fmt.Sprintf("%s_%d_%d", objectId.Hex(), j+1, len(pans))
					if pans[j]["fs"] != nil {
						hasGame = true
					}
				}
				times := ut.GetFloat(pans[len(pans)-1]["aw"])
				data := gendata.SimulateData{
					Id:       objectId,
					DropPan:  pans,
					HasGame:  hasGame,
					Times:    times,
					BucketId: gendata.GBuckets.GetBucket(times, hasGame, db.BoundType(gendata.GameTypeNormal)),
					Type:     gendata.GameTypeNormal,
					Selected: true,
				}
				lo.Must(newColl.InsertOne(context.TODO(), data))
				fmt.Println("success", objectId.Hex())
			}
			db.MiscSet(collName, _id)
			pans = []bson.M{}
		}
		pans = append(pans, pan)
	}
}

func dealGame(i int) {
	collName := fmt.Sprintf("pgSpinDataGame_%d", i)
	coll := db.Collection(collName)
	newColl := db.Collection("simulate")
	var lastId primitive.ObjectID
	db.MiscGet(collName, &lastId)
	match := bson.M{}
	if lastId != primitive.NilObjectID {
		match["_id"] = bson.M{"$gte": lastId}
	}
	cursor, _ := coll.Find(context.TODO(), match)
	defer cursor.Close(context.TODO())

	pans := []bson.M{}
	for cursor.Next(context.TODO()) {
		var pan map[string]any
		err := cursor.Decode(&pan)
		if err != nil {
			fmt.Println("err:", err,
				"_id:", pan["_id"],
				"sid:", pan["sid"],
				"psid:", pan["psid"])
			break
		}
		_id := pan["_id"]
		pan["orignid"] = fmt.Sprintf("%s_%s", pan["_id"], collName)
		delete(pan, "_id")
		if pan["sid"].(string) == pan["psid"].(string) && len(pans) > 0 {
			if pans[0]["sid"].(string) == pans[0]["psid"].(string) { // 第一盘是否有效
				objectId := primitive.NewObjectID()
				hasGame := false
				for j := 0; j < len(pans); j++ {
					pans[j]["psid"] = objectId.Hex()
					pans[j]["sid"] = fmt.Sprintf("%s_%d_%d", objectId.Hex(), j+1, len(pans))
					if pans[j]["fs"] != nil {
						hasGame = true
					}
				}
				times := ut.GetFloat(pans[len(pans)-1]["aw"])
				data := gendata.SimulateData{
					Id:       objectId,
					DropPan:  pans,
					HasGame:  hasGame,
					Times:    times,
					BucketId: gendata.GBuckets.GetBucket(times, hasGame, db.BoundType(gendata.GameTypeGame)),
					Type:     gendata.GameTypeGame,
					Selected: true,
				}
				lo.Must(newColl.InsertOne(context.TODO(), data))
				fmt.Println("success", objectId.Hex())
			}
			db.MiscSet(collName, _id)
			pans = []bson.M{}
		}
		pans = append(pans, pan)
	}
}

type ClassData struct {
	Gte         int
	Lte         int
	Name        string
	Probability int
}

func TestRandomizedBucket(t *testing.T) {
	var BucketHeartBeat ClassData
	BucketHeartBeat.Gte = 1
	BucketHeartBeat.Lte = 6
	BucketHeartBeat.Name = "BucketHeartBeat"
	BucketHeartBeat.Probability = 40
	BucketHeartBeat.ClassDataFunc()

	var BucketStable ClassData
	BucketStable.Gte = 7
	BucketStable.Lte = 17
	BucketStable.Name = "BucketStable"
	BucketStable.Probability = 25
	BucketStable.ClassDataFunc()
}

func (c *ClassData) ClassDataFunc() {
	mongoaddr := "mongodb://myUser:123456@192.168.1.2:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_89")
	coll := db.Collection("simulate")
	var cancelGatherNum int
	{

		fmt.Println("----------------数据重置----------------")

		resetUpdateFilter := bson.M{
			"bucketid": bson.M{
				"$lte": 28,
			},
		}
		resetUpdate := bson.M{"$set": bson.M{c.Name: float32(1)}}

		_, err := coll.UpdateMany(context.TODO(), resetUpdateFilter, resetUpdate)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("----------------数据查询----------------")
		findFilter := bson.M{
			"$and": []bson.M{
				{"bucketid": bson.M{"$gte": c.Gte}},
				{"bucketid": bson.M{"$lte": c.Lte}},
				{"selected": true},
			},
		}
		projection := bson.M{"_id": true}

		cursor, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(projection))
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var ids []string
		for cursor.Next(context.TODO()) {
			var result struct {
				Id string `bson:"_id,omitempty"`
			}
			err = cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			ids = append(ids, result.Id)
		}

		fmt.Println("----------------数据处理----------------")
		fmt.Println("总条数", len(ids))

		var cancelGather []primitive.ObjectID
		cancelGatherNum = (len(ids) * c.Probability) / 100
		for i := 0; i < cancelGatherNum; i++ {
			GatherNum := len(ids)
			popNum := rand.Int64N(int64(GatherNum))

			var cancelResult string
			ids, cancelResult = popAtIndex(ids, popNum)
			hex, _ := primitive.ObjectIDFromHex(cancelResult)
			cancelGather = append(cancelGather, hex)
		}

		fmt.Println("概率条数", cancelGatherNum)
		fmt.Println("修改条数", len(cancelGather))
		fmt.Println("----------------数据处理----------------")

		fmt.Println("数据修改")

		updateFilter := bson.M{"_id": bson.M{"$in": cancelGather}}
		update := bson.M{"$set": bson.M{c.Name: float32(0)}}

		_, err = coll.UpdateMany(context.TODO(), updateFilter, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("处理完毕")
	}
	{
		OldCountFilter := db.D("selected", true, "bucketid", bson.D{{"$lte", 28}})
		OldCount, err := coll.CountDocuments(context.TODO(), OldCountFilter)
		if err != nil {
			log.Fatal(err)
		}

		OldMatchStage := bson.D{{"$match", bson.D{
			{"selected", true},
			{"bucketid", bson.D{{"$lte", 28}}},
		}}}
		OldGroupStage := bson.D{{"$group", bson.D{
			{"_id", nil},
			{"totalAmount", bson.D{{"$sum", "$times"}}},
		}}}
		OldCursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{OldMatchStage, OldGroupStage})
		if err != nil {
			log.Fatal(err)
		}

		var OldResults []struct {
			TotalAmount float64 `bson:"totalAmount"`
		}
		err = OldCursor.All(context.TODO(), &OldResults)
		if err != nil {
			log.Fatal(err)
		}
		var OldSum float64
		if len(OldResults) > 0 {
			OldSum = OldResults[0].TotalAmount
		}

		NewMatchStage := bson.D{{"$match", bson.D{
			{"selected", true},
			{c.Name, 1},
			{"bucketid", bson.D{{"$lte", 28}}},
		}}}
		NewGroupStage := bson.D{{"$group", bson.D{
			{"_id", nil},
			{"totalAmount", bson.D{{"$sum", "$times"}}},
		}}}
		NewCursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{NewMatchStage, NewGroupStage})
		if err != nil {
			log.Fatal(err)
		}

		var NewResults []struct {
			TotalAmount float64 `bson:"totalAmount"`
		}
		err = NewCursor.All(context.TODO(), &NewResults)
		if err != nil {
			log.Fatal(err)
		}
		var NewSum float64
		if len(NewResults) > 0 {
			NewSum = NewResults[0].TotalAmount
		}

		fmt.Println("OldCount", OldCount, "OldSum", OldSum, "NewSum", NewSum)
		Num := int(float64(OldCount) - NewSum/OldSum*float64(OldCount))
		fmt.Println("Num", Num)
		fmt.Println("cancelGatherNum", cancelGatherNum)
		fmt.Println("Num - cancelGatherNum", Num-cancelGatherNum)

		findFilter := bson.M{
			"$and": []bson.M{
				{"bucketid": 0},
				{c.Name: 1},
				{"selected": true},
			},
		}
		projection := bson.M{"_id": true}

		cursor, err := coll.Find(context.TODO(), findFilter, options.Find().SetProjection(projection))
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.TODO())

		var ids []string
		for cursor.Next(context.TODO()) {
			var result struct {
				Id string `bson:"_id,omitempty"`
			}
			err = cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			ids = append(ids, result.Id)
		}

		fmt.Println("原ids数量", len(ids))

		var cancelGather []primitive.ObjectID
		for i := 0; i < Num-cancelGatherNum; i++ {
			GatherNum := len(ids)
			popNum := rand.Int64N(int64(GatherNum))

			var cancelResult string
			ids, cancelResult = popAtIndex(ids, popNum)
			hex, _ := primitive.ObjectIDFromHex(cancelResult)
			cancelGather = append(cancelGather, hex)
		}

		fmt.Println("现ids数量", len(ids))

		updateFilter := bson.M{"_id": bson.M{"$in": cancelGather}}
		update := bson.M{"$set": bson.M{c.Name: float32(0)}}

		_, err = coll.UpdateMany(context.TODO(), updateFilter, update)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func popAtIndex(slice []string, index int64) ([]string, string) {
	poppedElement := slice[index]
	return append(slice[:index], slice[index+1:]...), poppedElement
}
