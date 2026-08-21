package main

import (
	"fmt"
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	spinMap := map[int]jdbcomm.SpinFunc{
		90: norSpin,
		10: buySpin,
	}
	//enoughMap := map[int]string{}
	game := &jdbcomm.Game{
		ID:        "566",
		DBName:    "jdb_14087",
		MType:     "14087",
		GameType:  "14",
		GName:     "PopPopCandy_b0bd53a",
		UniqueKey: "1744680024812@7bbcb101-1af2-45fb-b84c-84a7d41f2da5 demo001526@XX",
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
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "20")
	entity.PutString("buyFeatureType", "null")

	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "QuantityGame")
	betRequest.PutInt("quantityBet", 1)

	sfsTop.PutSFSObject("p", &p)
	p.PutSFSObject("p", &pp)
	pp.PutSFSObject("entity", &entity)
	entity.PutSFSObject("betRequest", &betRequest)

	fmt.Println("current is normal")
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
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "20")
	entity.PutString("buyFeatureType", "BuyFeature_01")

	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "QuantityGame")
	betRequest.PutInt("quantityBet", 1)

	sfsTop.PutSFSObject("p", &p)
	p.PutSFSObject("p", &pp)
	pp.PutSFSObject("entity", &entity)
	entity.PutSFSObject("betRequest", &betRequest)

	fmt.Println("current is normal")
	return sfsTop, 1
}
