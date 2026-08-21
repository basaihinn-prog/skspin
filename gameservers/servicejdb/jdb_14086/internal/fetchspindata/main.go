package main

import (
	"fmt"
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	////decoding, _ := jdbcomm.Base64Decoding("gADhEgADAAFjAgEAAWEDAA0AAXASAAMAAWMIAAdoNS5zcGluAAFyBP////8AAXASAAEABmVudGl0eRIABgAFZGVub20IAAIxMAAMZXh0cmFCZXRUeXBlCAAKTm9FeHRyYUJldAALZ2FtZVN0YXRlSWQIAAEwAAlwbGF5ZXJCZXQIAAIyMAAOYnV5RmVhdHVyZVR5cGUIAA1CdXlGZWF0dXJlXzAxAApiZXRSZXF1ZXN0EgADAAdiZXRUeXBlCAAHV2F5R2FtZQAJYmV0Q29sdW1uBAAAAAYABndheUJldAQAAAAB")
	decoding, _ := jdbcomm.Base64Decoding("gAA1EgADAAFjAgEAAWEDAA0AAXASAAMAAWMIAA1HRU5fSEVBUlRCRUFUAAFyBP////8AAXASAAA=")
	fmt.Println(decoding)
	spinMap := map[int]jdbcomm.SpinFunc{
		90: norSpin, //90机率抓这个
		10: buySpin, //10机率抓这个
	}
	//enoughMap := map[int]string{}
	game := &jdbcomm.Game{
		ID:        "561",
		DBName:    "jdb_14086",
		MType:     "14086",
		GameType:  "14",
		GName:     "OpenSesameMega_152f421",
		UniqueKey: "1744180714073@1f6dd31b-1caa-4aeb-a06b-632b2a2b79ae demo003478@XX",
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
	betRequest.PutInt("betColumn", 6)
	betRequest.PutInt("wayBet", 1)

	entity := jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "20")
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

func buySpin() (jdbcomm.SFSObject, int) {
	betRequest := jdbcomm.SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "WayGame")
	betRequest.PutInt("betColumn", 6)
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
	return sfsTop, 1
}
