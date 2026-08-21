package main

import (
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	spinMap := map[int]jdbcomm.SpinFunc{
		0:   norSpin,
		100: buySpin,
	}

	game := &jdbcomm.Game{
		ID:        "141",
		DBName:    "jdb_14065",
		MType:     "14065",
		GameType:  "14",
		GName:     "BlossomOfWealth_3fb1c56",
		UniqueKey: "1743213787517@c682c890-eb29-426b-ac85-1baa3f8d1765 demo000410@XX",
		WeightMap: spinMap,
		// Spin: buySpin,
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

func norSpin() (jdbcomm.SFSObject, int) {
	sfsTop := jdbcomm.SFSObject{}
	sfsTop.Init()
	sfsTop.PutByte("c", 1)
	sfsTop.PutShort("a", 13)

	p := jdbcomm.SFSObject{}
	p.Init()
	p.PutString("c", "h5.spin")
	p.PutInt("r", -1)

	pp := jdbcomm.SFSObject{}
	pp.Init()

	entity := jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("playerBet", "50")
	entity.PutString("gameStateId", "0")

	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "WayGame")
	betRequest.PutInt("betColumn", 5)
	betRequest.PutInt("wayBet", 1)

	sfsTop.PutSFSObject("p", &p)
	p.PutSFSObject("p", &pp)
	pp.PutSFSObject("entity", &entity)
	entity.PutSFSObject("betRequest", &betRequest)

	return sfsTop, 0
}

func buySpin() (jdbcomm.SFSObject, int) {
	sfsTop := jdbcomm.SFSObject{}
	sfsTop.Init()
	sfsTop.PutByte("c", 1)
	sfsTop.PutShort("a", 13)

	p := jdbcomm.SFSObject{}
	p.Init()
	p.PutString("c", "h5.spin")
	p.PutInt("r", -1)

	pp := jdbcomm.SFSObject{}
	pp.Init()

	entity := jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("playerBet", "50")
	entity.PutString("gameStateId", "0")
	entity.PutString("buyFeatureType", "BuyFeature_01")

	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "WayGame")
	betRequest.PutInt("betColumn", 5)
	betRequest.PutInt("wayBet", 1)

	sfsTop.PutSFSObject("p", &p)
	p.PutSFSObject("p", &pp)
	pp.PutSFSObject("entity", &entity)
	entity.PutSFSObject("betRequest", &betRequest)

	return sfsTop, 1
}
