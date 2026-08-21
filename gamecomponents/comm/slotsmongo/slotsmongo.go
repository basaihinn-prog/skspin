package slotsmongo

import (
	"game/comm/db"

	"go.mongodb.org/mongo-driver/mongo"
)

func SlotsCollection(name string) *mongo.Collection {
	return db.Collection2("slots", name)
}

func GameCollection(name string) *mongo.Collection {
	return db.Collection2("game", name)
}

func ReportsCollection(name string) *mongo.Collection {
	return db.Collection2("reports", name)
}
