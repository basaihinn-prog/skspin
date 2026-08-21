package dealData

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"serve/comm/db"
	"serve/comm/ut"
	"serve/service/pg_50/internal/gendata"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestSimulate(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_50")

	collNames, _ := db.Client().Database("pg_50").ListCollectionNames(context.TODO(), bson.M{})

	for i := 0; i < len(collNames); i++ {
		if strings.HasPrefix(collNames[i], "pgSpinData_") {
			name := collNames[i]
			go dealNormal(name)
		}
	}
	// go dealGame()
	select {}
}

func TestCaclBucketId(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_50")
	coll := db.Collection("simulate")
	// 查询返回这三个字段，times，hasgame，type
	// options.FindOptions是mongoDB中配置查询的方法
	cur, _ := coll.Find(context.TODO(), db.D(), options.Find().SetProjection(db.D("times", 1, "hasgame", 1, "type", 1)))

	// mongo.WriteModel---mongoDB写入操作
	models := []mongo.WriteModel{}
	// 遍历查询的结果
	// cur.Next(context.TODO()) 是循环的条件部分，表示在游标 cur 上调用 Next 方法，如果能够移动到下一个文档则返回 true，否则返回 false。
	// cur.Decode 将文档解码为 Go 中的数据结构
	for cur.Next(context.TODO()) {
		var doc struct {
			ID      primitive.ObjectID `bson:"_id"`
			Times   float64            `bson:"times"`
			HasGame bool               `bson:"hasgame"`
			Type    db.BoundType       `bson:"type"`
		}
		assert.Nil(t, cur.Decode(&doc))
		// 这一步是干嘛用的
		bucketId := gendata.GBuckets.GetBucket(doc.Times, doc.HasGame, doc.Type)

		model := mongo.NewUpdateOneModel().SetFilter(db.ID(doc.ID)).SetUpdate(db.D("$set", db.D("bucketid", bucketId, "selected", true)))
		models = append(models, model)
	}

	_, err := coll.BulkWrite(context.TODO(), models, options.BulkWrite().SetOrdered(false))
	if err != nil {
		fmt.Println(err)
	}
}

func dealNormal(collName string) {
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
				times := ut.Round6(ut.GetFloat(pans[len(pans)-1]["aw"]) / 0.9)
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

func TestAddLimitCount(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "pg_50")

	coll := db.Collection("simulate")
	match := bson.M{
		"bucketid": 0,
	}

	for i := 0; i < 1000; i++ {
		cursor, err := coll.Aggregate(context.TODO(), bson.A{
			bson.D{{"$match", match}},
			bson.D{
				{"$sample", bson.D{{"size", 10}}},
			},
		})
		// id, _ := primitive.ObjectIDFromHex("6604d1cce04f4526f342ac8a")
		// cursor, _ := coll.Find(context.TODO(), db.ID(id))
		if err != nil {
			fmt.Println(err)
			return
		}

		// cursor.Next()
		var docs []*gendata.SimulateData
		err = cursor.All(context.TODO(), &docs)
		if err != nil {
			fmt.Println(err)
			return
		}
		var bsonDocs []interface{}
		for j := range docs {
			docs[j].Id = primitive.NewObjectID()
			bsonDocs = append(bsonDocs, docs[j])
		}
		_, err = coll.InsertMany(context.TODO(), bsonDocs)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(i)
	}
}
