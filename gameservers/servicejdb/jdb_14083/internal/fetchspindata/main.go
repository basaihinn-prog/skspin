package main

import (
	"fmt"
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	spinMap := map[int]jdbcomm.SpinFunc{
		100: norSpin,
	}
	//enoughMap := map[int]string{}
	game := &jdbcomm.Game{
		ID:        "158",
		DBName:    "jdb_14083",
		MType:     "14083",
		GameType:  "14",
		GName:     "CooCooFarm_0237e87",
		UniqueKey: "1744438161221@65f7f0d0-a3a7-49cd-b7a0-09b7fe3bfa23 demo000951@XX",
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
	entity.PutString("playerBet", "10")
	entity.PutString("buyFeatureType", "null")

	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "LineGame")
	betRequest.PutInt("betLine", 10)
	betRequest.PutInt("lineBet", 1)

	sfsTop.PutSFSObject("p", &p)
	p.PutSFSObject("p", &pp)
	pp.PutSFSObject("entity", &entity)
	entity.PutSFSObject("betRequest", &betRequest)

	fmt.Println("current is normal")
	return sfsTop, 0
}
