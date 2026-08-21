package api

import (
	"context"
	"game/comm/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// https://api.pg-demo.com/web-api/game-proxy/v2/BetSummary/Get?traceId=ZFLOQA20
// gid=39&dtf=1710864000000&dtt=1710950399999&atk=9A6074B5-B584-47C2-8508-A44FF6722BF4&pf=1&wk=0_C&btt=1
// {"dt":{"lut":1710921041616,"bs":{"gid":39,"bc":20,"btba":1380.00,"btwla":-600.00,"lbid":1770357278889876993}},"err":null}
// {"dt":{"lut":0,"bs":null},"err":null}

// type BsonA = bson.A
// type BsonD = bson.D
func getBetSummary(ps *PGParams, ret *M) (err error) {
	gid := ps.Get("gid")
	pipeline := mongo.Pipeline{
		bson.D{
			{"$match",
				bson.D{
					{"bt",
						bson.D{
							{"$gte", ps.GetInt("dtf")},
							{"$lte", ps.GetInt("dtt")},
						},
					},
					{"pid", ps.Pid},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$pid"},
					{"count", bson.D{{"$count", bson.D{}}}},
					{"betsum", bson.D{{"$sum", "$gtba"}}},
					{"npsum", bson.D{{"$sum", "$gtwla"}}},
					{"lbid", bson.D{{"$last", "$tid"}}},
					{"lut", bson.D{{"$last", "$bt"}}},
				},
			},
		},
	}

	coll := db.Collection2("pg_"+gid, "BetHistory")
	cur, err := coll.Aggregate(context.TODO(), pipeline)
	// _id 100308
	// count 4
	// betsum 4
	// npsum 196
	// lbid "1710928526745358572"
	// lut 1710928526747

	if err != nil {
		return
	}

	var docs []struct {
		Pid    int `bson:"_id"`
		Count  int
		Betsum float64
		Npsum  float64
		Lbid   string
		Lut    int64
	}
	err = cur.All(context.TODO(), &docs)
	if err != nil {
		return
	}

	if len(docs) == 0 {
		*ret = M{
			// {"dt":{"lut":0,"bs":null},"err":null}
			"lut": 0,
			"bs":  nil,
		}
		return
	}

	doc := docs[0]
	*ret = M{
		"lut": doc.Lut,
		"bs": M{
			"gid":   ps.GetInt("gid"),
			"bc":    doc.Count,
			"btba":  doc.Betsum,
			"btwla": doc.Npsum,
			"lbid":  doc.Lbid,
		},
	}

	return
}
