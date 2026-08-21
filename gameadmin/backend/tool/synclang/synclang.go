package main

import (
	"context"
	"flag"
	"game/comm/db"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func dialToMongo(addr string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOps := options.Client().ApplyURI(addr)
	client := lo.Must(mongo.Connect(ctx, clientOps))

	return client
}
func main() {
	var (
		from, to string
	)
	flag.StringVar(&from, "from", "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin&directConnection=true", "from mongourl")
	flag.StringVar(&to, "to", "mongodb://127.0.0.1:27017/", "to mongourl")
	flag.Parse()

	// from := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin&directConnection=true"
	// to := "mongodb://myUserAdmin:doudou123456@127.0.0.1:27017/?authSource=admin&directConnection=true"
	// to := "mongodb://127.0.0.1:27017/"
	client_from := dialToMongo(from)
	client_to := dialToMongo(to)

	coll_from := client_from.Database("game").Collection("Lang")
	coll_to := client_to.Database("game").Collection("Lang")

	cur := lo.Must(coll_from.Find(context.TODO(), db.D()))

	var models []mongo.WriteModel
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		// cur.Decode()
		// coll_to.InsertOne(context.TODO(), cur.Current)

		md := mongo.NewReplaceOneModel().
			SetFilter(db.ID(cur.Current.Lookup("_id"))).
			SetReplacement(cur.Current).
			SetUpsert(true)

		models = append(models, md)
	}

	lo.Must(coll_to.BulkWrite(context.TODO(), models, options.BulkWrite().SetOrdered(false)))

	// bson.RawValue
	// cur.All()
	// cur.All(context.TODO(), )
}
