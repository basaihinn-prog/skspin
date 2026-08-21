package main

import (
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	game := &jdbcomm.Game{
		ID:        "98",
		DBName:    "jdb_14035",
		MType:     "14035",
		GameType:  "14",
		GName:     "DragonsWorld_1b7d877",
		UserName:  "demo003470@XX",
		UniqueKey: "1743130604272@4f091222-e1c5-4632-81d6-dedc159aaa3e demo004166@XX",
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
	betRequest.PutString("betType", "WayGame")
	betRequest.PutInt("betColumn", 5)
	betRequest.PutInt("wayBet", 1)

	entity := jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "ExtraBet_01")
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "45")
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
