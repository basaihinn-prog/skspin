package main

import (
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	game := &jdbcomm.Game{
		ID:        "94",
		DBName:    "jdb_14027",
		MType:     "14027",
		GameType:  "14",
		GName:     "LuckySeven_f11946f",
		UserName:  "demo003470@XX",
		UniqueKey: "1741146017009@eb98bf56-ff5f-4dc8-896a-1920cb612498 demo001652@XX",
		Spin:      spin,
	}

	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("********* 程序panic，错误信息: %+v *********", r)
					log.Println("********* 3秒后将自动重启... *********")
					time.Sleep(time.Second * 3) // 休息3秒再重启，不然疯狂panic太快了
				}
			}()
			game.Run()
		}()
	}
}

func spin() jdbcomm.SFSObject {
	sfsTop := jdbcomm.SFSObject{}
	sfsTop.Init()
	sfsTop.PutByte("c", 1)
	sfsTop.PutShort("a", 13)

	so := jdbcomm.SFSObject{}
	so.Init()
	so.PutString("c", "h5.spin")
	so.PutInt("r", -1)
	entity := jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "9")
	entity.PutString("buyFeatureType", "null")
	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "LineGame")
	betRequest.PutInt("betLine", 9)
	betRequest.PutInt("lineBet", 1)
	entity.PutSFSObject("betRequest", &betRequest)
	p := jdbcomm.SFSObject{}
	p.Init()
	p.PutSFSObject("entity", &entity)
	so.PutSFSObject("p", &p)
	sfsTop.PutSFSObject("p", &so)
	return sfsTop
}
