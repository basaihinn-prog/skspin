package main

import (
	"game/duck/ut2"
	"game/pb/_gen/pb/adminpb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Token struct {
	Pid      int64
	UserName string
	Token    string
	ExpireAt time.Time
}

var tokenMap = ut2.NewSyncMap[string, *Token]()

func getToken(token string) *Token {

	tk, ok := tokenMap.Load(token)
	if !ok || tk.UserName == "" {
		one := adminpb.DBAdminer{}
		err := CollAdminUser.FindOne(bson.M{"Token": token}, &one)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil
			}
			return nil
		}
		// 禁用管理平台账号，不返回token
		if one.Status == 1 {
			return nil
		}

		if time.Now().After(one.TokenExpireAt.AsTime()) {
			return nil
		}

		tk = &Token{
			Token:    one.Token,
			ExpireAt: one.TokenExpireAt.AsTime(),
			Pid:      one.ID,
			UserName: one.Username,
		}
		tokenMap.Store(tk.Token, tk)
	}

	return tk
}
