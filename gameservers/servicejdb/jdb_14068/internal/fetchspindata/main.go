package main

import (
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	game := &jdbcomm.Game{
		ID:        "144",
		DBName:    "jdb_14068",
		MType:     "14068",
		GameType:  "14",
		GName:     "ProsperityTiger_f9015dd",
		UniqueKey: "1743128689407@c7c45051-335e-4c4e-9b75-cf67860046da demo001685@XX",
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
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", "55")
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
