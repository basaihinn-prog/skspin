package main

import (
	"fmt"
	"log"
	"serve/servicejdb/jdbcomm"
	"time"
)

func main() {

	decoding, _ := jdbcomm.Base64Decoding("gAD6EgADAAFjAgEAAWEDAA0AAXASAAMAAWMIAAhiYXIuc3BpbgABcgT/////AAFwEgADAARjb2RlCAAIYmFyLnNwaW4AAnNuBQAAuC0kQWEmAAZlbnRpdHkSAAQABWRlbm9tBAAAAAoACXBsYXllckJldAUAAAAAAAAAKAADYmV0DQAIAAAAAAAAAAUAAAAAAAAABQAAAAAAAAAFAAAAAAAAAAUAAAAAAAAABQAAAAAAAAAFAAAAAAAAAAUAAAAAAAAABQANZXh0ZW5kUmVxdWVzdBIAAgARZXh0ZW5kUmVxdWVzdFR5cGUIAARFUjAyAARkYXRhDAABAAAAAA==")
	// decoding, _ := jdbcomm.Base64Decoding("gAA1EgADAAFjAgEAAWEDAA0AAXASAAMAAWMIAA1HRU5fSEVBUlRCRUFUAAFyBP////8AAXASAAA=")
	fmt.Println(decoding)
	game := &jdbcomm.Game{
		ID:        "86",
		DBName:    "jdb_9007",
		MType:     "9007",
		GameType:  "9",
		GName:     "SuperSuperFruit_e3c5b00",
		UniqueKey: "1743564748046@facd703e-5dad-49c2-ac60-f70d9d2950e7 demo002803@XX",
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
	pp.PutString("sn", "202504011409590")

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
