package main

import (
	"fmt"
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	spinMap := map[int]jdbcomm.SpinFunc{
		80: norSpin,
		10: buySpin01,
		11: buySpin02,
	}
	//enoughMap := map[int]string{}
	game := &jdbcomm.Game{
		ID:        "825",
		DBName:    "jdb_14093",
		MType:     "14093",
		GameType:  "14",
		GName:     "MagicAceWildLock_049b158",
		UniqueKey: "1744681582484@1e8b80c2-7a57-44f1-a906-196009cb0bb6 demo002689@XX",
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
	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "WayGame")
	betRequest.PutInt("betColumn", 5)
	betRequest.PutInt("wayBet", 3)

	entity := jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "NoExtraBet")
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
	fmt.Println("current is normal")
	return sfsTop, 0
}

func buySpin01() (jdbcomm.SFSObject, int) {
	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "WayGame")
	betRequest.PutInt("betColumn", 5)
	betRequest.PutInt("wayBet", 1)

	entity := jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "20")
	entity.PutString("buyFeatureType", "BuyFeature_02")
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
	fmt.Println("current is Buy")
	return sfsTop, 1
}

func buySpin02() (jdbcomm.SFSObject, int) {
	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "WayGame")
	betRequest.PutInt("betColumn", 5)
	betRequest.PutInt("wayBet", 1)

	entity := jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "20")
	entity.PutString("buyFeatureType", "BuyFeature_01")
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
	fmt.Println("current is Buy")
	return sfsTop, 2
}
