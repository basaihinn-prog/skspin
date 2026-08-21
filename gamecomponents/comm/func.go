package comm

import (
	"context"
	"fmt"
	"game/comm/db"
	"game/comm/mq"
	"game/duck/lazy"
	"game/duck/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand/v2"
	"os"
)

const statdir = "statistics"

func OpenStatFile(file string, flag int) (f *os.File, err error) {
	os.Mkdir(statdir, 0755)
	filename := fmt.Sprintf("%v/%v_%v.stat", statdir, lazy.ServiceName, file)
	f, err = os.OpenFile(filename, flag|os.O_RDWR, 0644)
	if err != nil {
		logger.Err(err)
	}
	return
}

func DIV[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](num1, num2 T) T {
	if num2 == 0 {
		return 0
	}
	return num1 / num2
}

func RandPassword() string {
	baseStr := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*+-_="

	length := 8
	bytes := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		bytes = append(bytes, baseStr[rand.Int()%len(baseStr)])
	}
	return string(bytes)
}

type GetOperatorInfoRes struct {
	//商户币种        AdminOperator_copy1 / CurrencyKey
	CurrencyKey                   string         `json:"CurrencyKey" bson:"CurrencyKey"`
	CurrencyManufactureVisibleOff map[string]int `json:"CurrencyManufactureVisibleOff" bson:"CurrencyVisibleOff"`
}

func GetOperatorInfo(appId string) GetOperatorInfoRes {
	resInfo := GetOperatorInfoRes{}
	err := mq.Invoke("/AdminInfo/Interior/operatorInfo", map[string]any{
		"AppID": appId,
	}, &resInfo)

	if err != nil {
		return resInfo
	}
	return resInfo
}

type ExParams struct {
	ShowNameAndTimeOff int64  `bson:"ShowNameAndTimeOff" ` //名字和时间显示
	ShowExitBtnOff     int64  `bson:"ShowExitBtnOff" `     //退出按钮显示
	StopLoss           int64  `bson:"StopLoss"`            //止盈止损开关
	BackPackOff        int64  `bson:"BackPackOff"`         //背包显示开关
	CarouselOff        int64  `bson:"CarouselOff"`         //轮播开关
	ExitBtnOff         int64  `bson:"ExitBtnOff"`          //退出按钮显示
	ExitLink           string `bson:"ExitLink"`            //退出按钮链接
	OpenScreenOff      int64  `bson:"OpenScreenOff"`       //开屏弹窗开关
	SidebarOff         int64  `bson:"SidebarOff"`          //赢钱多
	OfficialVerify     int64  `bson:"OfficialVerify"`      //官方验证
}

func GetEXParams(pid int64, gameId string) *ExParams {
	var doc struct {
		AppID string `bson:"AppID"` // 所属产品
	}
	// var appid string
	coll := db.Collection2("game", "Players")
	coll.FindOne(context.TODO(), db.ID(pid), options.FindOne().SetProjection(db.D("AppID", 1))).Decode(&doc)

	var tmp *ExParams

	coll = db.Collection2("GameAdmin", "GameConfig")
	coll.FindOne(context.TODO(), bson.M{"AppID": doc.AppID, "GameId": gameId}, options.FindOne().SetProjection(db.D("ShowNameAndTimeOff", 1, "ShowExitBtnOff", 1, "StopLoss", 1, "BackPackOff", 1, "OpenScreenOff", 1, "SidebarOff", 1, "CarouselOff", 1, "ExitBtnOff", 1, "ExitLink", 1, "OfficialVerify", 1))).Decode(&tmp)

	return tmp
}

func GetCurrentItem(pid int64) *lazy.CurrencyItem {
	var item *lazy.CurrencyItem
	var doc struct {
		AppID string `bson:"AppID"` // 所属产品
	}
	// var appid string
	coll := db.Collection2("game", "Players")
	coll.FindOne(context.TODO(), db.ID(pid), options.FindOne().SetProjection(db.D("AppID", 1))).Decode(&doc)

	type currencyItem struct {
		CurrencyKey        string         `bson:"CurrencyKey"`
		CurrencyVisibleOff map[string]int `json:"CurrencyManufactureVisibleOff" bson:"CurrencyVisibleOff"`
	}
	var tmp *currencyItem

	coll = db.Collection2("GameAdmin", "AdminOperator")
	coll.FindOne(context.TODO(), bson.M{"AppID": doc.AppID}, options.FindOne().SetProjection(db.D("CurrencyKey", 1, "CurrencyVisibleOff", 1))).Decode(&tmp)
	if tmp != nil {
		item = lazy.GetCurrencyItem(tmp.CurrencyKey)
		item.CurrencyManufactureVisibleOff = tmp.CurrencyVisibleOff
	} else {
		item = lazy.GetCurrencyItem("")
	}
	return item
}

var currencyMap = map[string]string{
	"AED": "AE",
	"ALL": "ALL",
	"AMD": "AMD",
	"AOA": "AOA",
	"ARS": "ARS",
	"AZN": "AZN",
	"BAM": "BAM",
	"BGN": "BGN",
	"BHD": "BHD",
	"BIF": "BIF",
	"BMD": "BMD",
	"BND": "BN",
	"BOB": "BOB",
	"BRL": "BR",
	"BDT": "BT",
	"BYN": "BYN",
	"BZD": "BZ",
	"CHF": "CH",
	"CHW": "CH",
	"COP": "COP",
	"CZK": "CZK",
	"DKK": "DKK",
	"SEK": "DKK",
	"EGP": "EGP",
	"ETB": "ETB",
	"CHE": "EU",
	"EUR": "EU",
	"FKP": "GB",
	"GBP": "GB",
	"GIP": "GB",
	"SHP": "GB",
	"SYP": "GB",
	"GEL": "GEL",
	"GHS": "GHS",
	"GTQ": "GTQ",
	"HNL": "HNL",
	"HUF": "HUF",
	"ILS": "ILS",
	"IQD": "IQD",
	"IRR": "IRR",
	"ISK": "ISK",
	"JOD": "JOD",
	"JPY": "JP",
	"KES": "KES",
	"KGS": "KGS",
	"KMF": "KRWK",
	"KRW": "KRWK",
	"KWD": "KWD",
	"KZT": "KZT",
	"LAK": "LAK",
	"LBP": "LBP",
	"LKR": "LK",
	"MAD": "MAD",
	"MDL": "MDL",
	"MKD": "MKD",
	"MMK": "MM2",
	"MNT": "MNT",
	"MOP": "MOP",
	"MUR": "MUR",
	"MVR": "MVR",
	"MWK": "MWK",
	"MZN": "MZN",
	"NGN": "NGN",
	"CAD": "NIO",
	"NIO": "NIO",
	"NOK": "NO",
	"PEN": "PEN",
	"PLN": "PLN",
	"CUP": "PP",
	"PYG": "PYG",
	"CNY": "RB",
	"MYR": "RM",
	"RON": "RON",
	"IDR": "RPO",
	"INR": "RS",
	"RUB": "RUB",
	"THB": "TB",
	"TJS": "TJS",
	"TMT": "TMT",
	"TND": "TND",
	"SOS": "UGX",
	"USD": "US",
}

func GetJDBCurrentItem(pid int64) *lazy.CurrencyItem {
	var item lazy.CurrencyItem
	var doc struct {
		AppID string `bson:"AppID"` // 所属产品
	}
	// var appid string
	coll := db.Collection2("game", "Players")
	coll.FindOne(context.TODO(), db.ID(pid), options.FindOne().SetProjection(db.D("AppID", 1))).Decode(&doc)

	type currencyItem struct {
		CurrencyKey        string         `bson:"CurrencyKey"`
		CurrencyVisibleOff map[string]int `json:"CurrencyManufactureVisibleOff" bson:"CurrencyVisibleOff"`
	}
	var tmp *currencyItem

	coll = db.Collection2("GameAdmin", "AdminOperator")
	coll.FindOne(context.TODO(), bson.M{"AppID": doc.AppID}, options.FindOne().SetProjection(db.D("CurrencyKey", 1, "CurrencyVisibleOff", 1))).Decode(&tmp)
	if tmp != nil {
		item = *lazy.GetCurrencyItem(tmp.CurrencyKey)
		item.CurrencyManufactureVisibleOff = tmp.CurrencyVisibleOff
	} else {
		item = *lazy.GetCurrencyItem("")
	}
	item.Key = currencyMap[item.Key]
	//item.Key = "ZFS"
	return &item
}
