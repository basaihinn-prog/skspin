package main

import (
	"context"
	"encoding/json"
	"fmt"
	"game/comm/db"
	"game/duck/lang"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

type LangMap map[string]string

type NewLang struct {
	ZH         string `bson:"_id"`
	Permission int32  `bson:"permission"`
	LangMap    `json:",inline" bson:"inline"`
}

func TestJsonInline(t *testing.T) {
	mongoaddr := "mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin"
	db.DialToMongo(mongoaddr, "game")

	coll := db.Collection("Lang")
	var newLang NewLang
	err := coll.FindOne(context.TODO(), bson.M{}).Decode(&newLang)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(newLang)

	data, _ := json.Marshal(newLang)
	fmt.Println(string(data))
}

func TestMarsh(t *testing.T) {
	lang := &lang.Lang{
		ZH:         "测试",
		Permission: 0,
		LangMap: map[string]string{
			"en": "test",
			"th": "测试泰文",
		},
	}
	data, _ := json.Marshal(lang)
	fmt.Println(string(data))
	fmt.Println(data)
}

func TestUnMarsh(t *testing.T) {
	data := "{\"EN\":\"test\",\"Permission\":0,\"TH\":\"测试泰文\",\"ZH\":\"测试\"}"
	var lang *lang.Lang
	json.Unmarshal([]byte(data), &lang)
	fmt.Println(lang)
}

func TestRemoveBetLogDownload(t *testing.T) {
	DB.Connect("mongodb://myUserAdmin:doudou123456@156.241.5.141:27017/?authSource=admin")
	removeBetLogDownload()
}
