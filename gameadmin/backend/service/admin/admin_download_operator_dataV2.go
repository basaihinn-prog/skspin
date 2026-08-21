package main

import (
	"fmt"
	"game/comm"
	"game/duck/mongodb"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func init() {
	//弃用
	//RegMsgProc("/AdminInfo/DownloadOperatorDataV2", "下载运营商数据", "AdminInfo", downloadOperatorDataV2, downloadOperatorParamsV2{})
}

type downloadOperatorParamsV2 struct {
	Status           int64  `json:"Status" bson:"Status"`             //状态
	AppID            string `json:"Name" bson:"Name"`                 //商户名称
	OperatorType     int    `json:"OperatorType" bson:"OperatorType"` //商户类型
	CreatedStartTime int64  `json:"CreatedStartTime"`
	CreatedEndTime   int64  `json:"CreatedEndTime"`
	PageIndex        int64  `json:"PageIndex"`
	PageSize         int64  `json:"PageSize"`
}

type downloadOperatorResultV2 struct {
	List     *excelize.File
	AllCount int64
}

func downloadOperatorDataV2(ctx *Context, ps downloadOperatorParamsV2, ret *downloadOperatorResultV2) (err error) {
	user, _ := IsAdminUser(ctx)
	query := bson.M{}
	if ps.CreatedStartTime != 0 {
		startTime := mongodb.NewTimeStamp(time.UnixMilli(ps.CreatedStartTime))
		endTime := mongodb.NewTimeStamp(time.UnixMilli(ps.CreatedEndTime))
		query = bson.M{
			//"Name":user.AppID,
			"CreateTime": bson.M{
				"$gte": startTime, // 大于或等于开始时间
				"$lte": endTime,   // 小于或等于结束时间
			},
		}
	}

	if ps.AppID != "" { //模糊查询
		//query["AppID"] = bson.M{"$regex": "/" + ps.AppID + "/", "$options": 'i'}
		query["AppID"] = primitive.Regex{
			Pattern: ps.AppID,
			Options: "i",
		}
	}
	if ps.Status != -1 {
		query["Status"] = ps.Status
	}
	if ps.OperatorType != 0 {
		query["OperatorType"] = ps.OperatorType
	}

	switch user.GroupId {
	case 1:
	case 2:
		query["Name"] = user.AppID
	case 3:
		query["AppID"] = user.AppID
	default:
	}

	//var list *[]*comm.Operator_V2
	list := make([]*comm.Operator_V2, 0)
	//err = GetOperatorInfo(query, &list)
	f := down(list)
	//if err != nil {
	//	return err
	//}
	ret.List = f
	ret.AllCount = int64(len(list))

	return err
}
func down(list []*comm.Operator_V2) *excelize.File {
	f := excelize.NewFile()
	//生成表头
	headers := []string{
		"Id",
		"Name",
		"MenuIds",
		"ExcluedGameIds",
		"PermissionId",
		"AppID",
		"AppSecret",
		"UserName",
		"PlatformPay",
		"CooperationType",
		"PresentRate",
		"NextRate",
		"OperatorType",
		"Advance",
		"Status",
		"CreateTime",
		"TokenExpireAt",
		"CurrencyKey",
		"Contact",
		"WalletMode",
		"Surname",
		"Lang",
		"ServiceIp",
		"WhiteIps",
		"Address",
		"UserWhite",
		"LoginOff",
		"FreeOff",
		"DormancyOff",
		"RestoreOff",
		"ManualFullScreenOff",
		"NewGameDefaul",
		"MassageOff",
		"MassageIp",
		"RTPOff",
		"StopLoss",
		"MaxMultipleOff",
		"LineMerchant",
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i))) // A1, B1, C1
		err := f.SetCellValue("Sheet1", cell, header)
		if err != nil {
			return f
		}
	}
	return f

	// 写入结构体数据
	//for i, person := range list {
	//	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), person.Id)                   // A2, A3, A4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), person.Name)                 // B2, B3, B4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), person.MenuIds)              // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), person.ExcluedGameIds)       // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), person.PermissionId)         // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), person.AppID)                // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("G%d", i+2), person.AppSecret)            // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("H%d", i+2), person.UserName)             // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("I%d", i+2), person.PlatformPay)          // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("J%d", i+2), person.CooperationType)      // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("K%d", i+2), person.PresentRate)          // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("L%d", i+2), person.NextRate)             // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("M%d", i+2), person.OperatorType)         // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("N%d", i+2), person.Advance)              // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("O%d", i+2), person.Status)               // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("P%d", i+2), person.CreateTime)           // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("Q%d", i+2), person.TokenExpireAt)        // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("R%d", i+2), person.CurrencyKey)          // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("S%d", i+2), person.Contact)              // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("T%d", i+2), person.WalletMode)           // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("U%d", i+2), person.Surname)              // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("V%d", i+2), person.Lang)                 // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("W%d", i+2), person.ServiceIp)            // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("X%d", i+2), person.WhiteIps)             // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("Y%d", i+2), person.Address)              // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("Z%d", i+2), person.UserWhite)            // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AA%d", i+2), person.LoginOff)            // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AB%d", i+2), person.FreeOff)             // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AC%d", i+2), person.DormancyOff)         // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AD%d", i+2), person.RestoreOff)          // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AE%d", i+2), person.ManualFullScreenOff) // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AF%d", i+2), person.NewGameDefaulOff)       // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AG%d", i+2), person.MassageOff)          // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AH%d", i+2), person.MassageIp)           // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AI%d", i+2), person.RTPOff)              // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AJ%d", i+2), person.StopLoss)            // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AK%d", i+2), person.MaxMultipleOff)      // C2, C3, C4
	//	f.SetCellValue("Sheet1", fmt.Sprintf("AL%d", i+2), person.LineMerchant)        // C2, C3, C4
	//}
}
