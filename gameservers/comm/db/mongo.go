package db

import (
	"bytes"
	"context"
	"log/slog"
	"path"
	"sync"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	mongoDB     *mongo.Database
	collections = map[string]*mongo.Collection{}
	mtx         sync.Mutex
	client      *mongo.Client
)

func GetMongoDbName() string {
	return mongoDB.Name()
}

func UpsertOpt() *options.UpdateOptions {
	return options.Update().SetUpsert(true)
}

func Client() *mongo.Client {
	return client
}

func SetupClient(c *mongo.Client) {
	client = c
}

func DialToMongo(addr, db string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOps := options.Client().ApplyURI(addr)
	client = lo.Must(mongo.Connect(ctx, clientOps))

	lo.Must0(client.Ping(context.TODO(), readpref.Primary()))

	// connstr, _ := connstring.Parse(addr)

	slog.Info("Mongo Dail SUCCESS!", "addr", addr, "db", db)
	mongoDB = client.Database(db)
}

func D(ps ...any) bson.D {
	n := len(ps)
	lo.Must0(n%2 == 0)

	d := make(bson.D, n/2)
	for i := 0; i < n/2; i++ {
		k, v := ps[2*i], ps[2*i+1]
		d[i].Key = k.(string)
		d[i].Value = v
	}

	return d
}

func AppendE(d *bson.D, k string, v any) {
	*d = append(*d, bson.E{Key: k, Value: v})
}

func ID(id any) bson.D {
	return bson.D{
		{Key: "_id", Value: id},
	}
}

func Collection(name string) *mongo.Collection {
	mtx.Lock()
	defer mtx.Unlock()

	c := collections[name]
	if c == nil {
		c = mongoDB.Collection(name)
		collections[name] = c
	}

	return c
}

func Collection2(dbname, collname string) *mongo.Collection {
	mtx.Lock()
	defer mtx.Unlock()

	fname := path.Join(dbname, collname)

	c := collections[fname]
	if c == nil {
		c = client.Database(dbname).Collection(collname)
		collections[fname] = c
	}

	return c

}

func Json2Bson(jsonbuf []byte) (raw bson.Raw, err error) {
	vr, err := bsonrw.NewExtJSONValueReader(bytes.NewReader(jsonbuf), true)
	if err != nil {
		return
	}

	decoder, err := bson.NewDecoder(vr)
	if err != nil {
		return
	}
	// var obj bson.Raw
	err = decoder.Decode(&raw)

	return
}

func Json2BsonD(jsonbuf []byte) (d bson.D, err error) {
	vr, err := bsonrw.NewExtJSONValueReader(bytes.NewReader(jsonbuf), true)
	if err != nil {
		return
	}

	decoder, err := bson.NewDecoder(vr)
	if err != nil {
		return
	}

	// var obj bson.Raw
	err = decoder.Decode(&d)
	return
}

// CloseCursor 关闭游标，防止泄露
func CloseCursor(cur *mongo.Cursor) {
	// 确保游标关闭，防止资源泄漏
	err := cur.Close(context.TODO())
	if err != nil {
	}
}
