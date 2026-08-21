package jilicomm

import "serve/comm/db"

type JILIPlayer struct {
	db.DocPlayer `bson:"inline"`
}
