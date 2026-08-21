package main

import (
	"log"
	"time"
)

type ResponseAction101 struct {
	Status string `json:"status"`
	Data   struct {
		Ots string `json:"ots"`
	} `json:"data"`
}
type UserInfo struct {
	UID        string `json:"uid"`
	UserName   string `json:"userName"`
	Lvl        int    `json:"lvl"`
	UserStatus int    `json:"userStatus"`
	Currency   string `json:"currency"`
}

type ResponseAction6 struct {
	Status string     `json:"status"`
	Data   []UserInfo `json:"data"`
}

type ResponseAction19 struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	IsShowAutoPlay bool     `json:"isShowAutoPlay"`
	Result4        Result4  `json:"result4"`
	Result6        Result6  `json:"result6"`
	Result10       Result10 `json:"result10"`
}

type Result4 struct {
	Currency         string `json:"currency"`
	IsDemoAccount    bool   `json:"isDemoAccount"`
	IsApiAccount     bool   `json:"isApiAccount"`
	IsShowJackpot    bool   `json:"isShowJackpot"`
	IsShowCurrency   bool   `json:"isShowCurrency"`
	IsShowDollarSign bool   `json:"isShowDollarSign"`
	DecimalPoint     int    `json:"decimalPoint"`
	GameGroup        []int  `json:"gameGroup"`
	FunctionList     []int  `json:"functionList"`
}

type Result6 struct {
	UID        string `json:"uid"`
	UserName   string `json:"userName"`
	Lvl        int    `json:"lvl"`
	UserStatus int    `json:"userStatus"`
	Currency   string `json:"currency"`
}

type Result10 struct {
	Status               string   `json:"status"`
	SessionID            []string `json:"sessionID"`
	Zone                 string   `json:"zone"`
	GsInfo               string   `json:"gsInfo"`
	GameType             int32    `json:"gameType"`
	MachineType          int32    `json:"machineType"`
	IsRecovery           bool     `json:"isRecovery"`
	S0                   string   `json:"s0"`
	S1                   string   `json:"s1"`
	S2                   string   `json:"s2"`
	S3                   string   `json:"s3"`
	S4                   string   `json:"s4"`
	GameUid              string   `json:"gameUid"`
	GamePass             string   `json:"gamePass"`
	UseSSL               bool     `json:"useSSL"`
	StreamingUrl         struct{} `json:"streamingUrl"`
	AchievementServerUrl string   `json:"achievementServerUrl"`
	ChatServerUrl        string   `json:"chatServerUrl"`
	IsWSBinary           bool     `json:"isWSBinary"`
}

func main() {
	game := &hacksawcomm.Game{
		ID:        "94",
		DBName:    "hacksaw_1620",
		MType:     "14027",
		GameType:  "14",
		GName:     "LuckySeven_f11946f",
		UserName:  "demo003470@XX",
		UniqueKey: "1741146017009@eb98bf56-ff5f-4dc8-896a-1920cb612498 demo001652@XX",
		Base:      "gADYEgADAAFjAgEAAWEDAA0AAXASAAMAAWMIAAdoNS5zcGluAAFyBP////8AAXASAAEABmVudGl0eRIABgAFZGVub20IAAIxMAAMZXh0cmFCZXRUeXBlCAAKTm9FeHRyYUJldAALZ2FtZVN0YXRlSWQIAAEwAAlwbGF5ZXJCZXQIAAI1MAAOYnV5RmVhdHVyZVR5cGUIAARudWxsAApiZXRSZXF1ZXN0EgADAAdiZXRUeXBlCAAHV2F5R2FtZQAJYmV0Q29sdW1uBAAAAAUABndheUJldAQAAAAB",
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
			game.Run(hacksawcomm.CreateSpin)
		}()
	}
}
