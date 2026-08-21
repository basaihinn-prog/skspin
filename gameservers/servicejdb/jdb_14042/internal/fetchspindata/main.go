package main

import (
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	spinMap := map[int]jdbcomm.SpinFunc{
		0:   norSpin, //90机率抓这个
		100: buySpin, //10机率抓这个
	}

	game := &jdbcomm.Game{
		ID:        "112",
		DBName:    "jdb_14042",
		MType:     "14042",
		GameType:  "14",
		GName:     "TreasureBowl_b023fe0",
		UniqueKey: "1744167335420@2eaefc69-f85d-46f6-85d0-95892a501a52 demo003998@XX",
		// Spin:      spin,
		WeightMap: spinMap,
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
	entity.PutString("playerBet", "25")
	entity.PutString("gameStateId", "0")

	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "QuantityGame")
	betRequest.PutInt("quantityBet", 25)

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
	entity.PutString("playerBet", "5")
	entity.PutString("gameStateId", "0")
	entity.PutString("buyFeatureType", "BuyFeature_01")

	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "QuantityGame")
	betRequest.PutInt("quantityBet", 5)

	sfsTop.PutSFSObject("p", &p)
	p.PutSFSObject("p", &pp)
	pp.PutSFSObject("entity", &entity)
	entity.PutSFSObject("betRequest", &betRequest)

	return sfsTop, 1
}
