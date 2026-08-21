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
		ID:        "781",
		DBName:    "jdb_14091",
		MType:     "14091",
		GameType:  "14",
		GName:     "PiggyBank_56b6a82",
		UniqueKey: "1744679938640@5ea74abf-d6e6-4961-816d-eac7dd5d85bc demo001393@XX",
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
	entity.PutString("denom", "1000")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "1")
	entity.PutString("buyFeatureType", "null")

	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "LineGame")
	betRequest.PutInt("betLine", 1)
	betRequest.PutInt("lineBet", 1)

	sfsTop.PutSFSObject("p", &p)
	p.PutSFSObject("p", &pp)
	pp.PutSFSObject("entity", &entity)
	entity.PutSFSObject("betRequest", &betRequest)

	fmt.Println("current is normal")
	return sfsTop, 0
}

func buySpin() (jdbcomm.SFSObject, int) {
	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "LineGame")
	betRequest.PutInt("betLine", 1)
	betRequest.PutInt("lineBet", 1)

	entity := jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("denom", "1000")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "1")
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
	return sfsTop, 1
}
