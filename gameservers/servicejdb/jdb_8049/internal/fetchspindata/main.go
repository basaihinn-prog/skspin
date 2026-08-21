package main

import (
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	game := &jdbcomm.Game{
		ID:        "102",
		DBName:    "jdb_8049",
		MType:     "8049",
		GameType:  "8",
		GName:     "FlirtingScholarTangII_3babd4a",
		UserName:  "demo003470@XX",
		UniqueKey: "1743124795317@983873d8-00a8-4fa7-b5b8-9a5a35633e4a demo001716@XX",
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
	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "LineGame")
	betRequest.PutInt("betLine", 50)
	betRequest.PutInt("lineBet", 1)

	entity := jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "ExtraBet_01")
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "60")
	entity.PutString("buyFeatureType", "null")
	entity.PutSFSObject("betRequest", &betRequest)

	pp := jdbcomm.SFSObject{}
	pp.Init()
	pp.PutSFSObject("entity", &entity)

	p := jdbcomm.SFSObject{}
	p.Init()
	p.PutString("c", "h5.spin")
	p.PutInt("r", -1)
	p.PutSFSObject("p", &pp)

	sfsTop := jdbcomm.SFSObject{}
	sfsTop.Init()
	sfsTop.PutByte("c", 1)
	sfsTop.PutShort("a", 13)
	sfsTop.PutSFSObject("p", &p)

	return sfsTop
}
