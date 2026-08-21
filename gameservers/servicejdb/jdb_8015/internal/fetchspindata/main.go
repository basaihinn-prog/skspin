package main

import (
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	game := &jdbcomm.Game{
		ID:        "104",
		DBName:    "jdb_8015",
		MType:     "8015",
		GameType:  "8",
		GName:     "MoonlightTreasure_a89bf4a",
		UniqueKey: "1743573682996@6e0d11c9-9624-4750-b4f5-27bd8d1410f0 demo000010@XX",
		Bet:       40,
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
	sfsTop := jdbcomm.SFSObject{}
	sfsTop.Init()
	sfsTop.PutByte("c", 1)
	sfsTop.PutShort("a", 13)

	p := &jdbcomm.SFSObject{}
	p.Init()
	p.PutString("c", "h5.spin")
	p.PutInt("r", -1)

	pp := &jdbcomm.SFSObject{}
	pp.Init()
	pp.PutString("code", "h5.spin")
	pp.PutString("sn", "202504021357400")

	entity := &jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("betLine", "40")
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("lineBet", "1")
	entity.PutString("operation", "baseGame")
	entity.PutString("playerBet", "40")
	entity.PutString("wayGameBetColumn", "0")
	entity.PutString("waysBet", "1")

	sfsTop.PutSFSObject("p", p)
	p.PutSFSObject("p", pp)
	pp.PutSFSObject("entity", entity)

	return sfsTop
}
