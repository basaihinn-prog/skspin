package main

import (
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {
	game := &jdbcomm.Game{
		ID:        "92",
		DBName:    "jdb_8022",
		MType:     "8022",
		GameType:  "8",
		GName:     "Mahjong_bf2cf1e",
		UniqueKey: "1743573009701@5b5b6e29-a7b7-4c5c-a305-6438b2ef9c0c demo005267@XX",
		Bet:       25,
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
	pp.PutString("sn", "202504021351160")

	entity := &jdbcomm.SFSObject{}
	entity.Init()
	entity.PutString("betLine", "25")
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("lineBet", "1")
	entity.PutString("operation", "baseGame")
	entity.PutString("playerBet", "25")
	entity.PutString("wayGameBetColumn", "0")
	entity.PutString("waysBet", "1")

	sfsTop.PutSFSObject("p", p)
	p.PutSFSObject("p", pp)
	pp.PutSFSObject("entity", entity)

	return sfsTop
}
