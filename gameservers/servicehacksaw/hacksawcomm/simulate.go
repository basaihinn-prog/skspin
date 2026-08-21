package hacksawcomm

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SimulateData struct {
	Id                   primitive.ObjectID `bson:"_id"`
	DropPan              []Variables        `bson:"droppan"` //自行解析的数据
	HasGame              bool               `bson:"hasgame"`
	Times                float64            `bson:"times"`
	BucketId             int                `bson:"bucketid"`
	Type                 int                `bson:"type"`
	Selected             bool               `bson:"selected"`
	RoundID              int                `bson:"RoundID"`     //数据ID
	TurnIndex            int                `bson:"TurnIndex"`   //轮次，他跟数据ID是属于一组的数据
	FreeFlag             int                `bson:"FreeFlag"`    //旋转类型   0 就是普通旋转, 1就是买免费  2  超级购买   1010代表10次购买   1012代表12次  2010  也是代表10次购买
	GroupFlag            int                `bson:"GroupFlag"`   //普通下注，双倍下注，部分游戏有双倍下注
	QueryString          Variables          `bson:"QueryString"` //数据值
	BucketHeartBeat      int                `bson:"BucketHeartBeat"`
	BucketWave           int                `bson:"BucketWave"`
	BucketGov            int                `bson:"BucketGov"`
	BucketMix            int                `bson:"BucketMix"`
	BucketStable         int                `bson:"BucketStable"`
	BucketHighAward      int                `bson:"BucketHighAward"`
	BucketSuperHighAward int                `bson:"BucketSuperHighAward"`
}

//	func (sd *SimulateData) Deal(c float64, line float64, balance float64) Variables {
//		bet := c * line
//		//slog.Info("SimulateData, ", "ObjectID", sd.Id.Hex())
//
//		pan := sd.DropPan[0]
//
//		rid := strings.Split(pan.Str("gid"), "_")[0]
//		pan.Set("rid", rid)
//		pan.SetFloat("c", c)
//		pan.MKMulFloat("tw", bet)
//		pan.MKMulFloat("tmb_win", bet)
//		pan.MKMulFloat("w", bet)
//		pan.MKMulFloat("wp", bet)
//		pan.MKMulFloat("fsres", bet)
//		pan.MKMulFloat("fswin", bet)
//		pan.MKMulFloat("fsres_total", bet)
//		pan.MKMulFloat("fswin_total", bet)
//		pan.MKMulFloat("tmb_res", bet)
//		pan.MKMulFloat("rs_iw", bet)
//		if pan.Get("rs_iw") != "" {
//			pan.Set("rs_iw", strings.ReplaceAll(pan.Get("rs_iw"), ",", "")) //去除逗号
//		}
//		pan.MKMulFloat("rs_win", bet)
//		if pan.Get("rs_win") != "" {
//			pan.Set("rs_win", strings.ReplaceAll(pan.Get("rs_win"), ",", "")) //去除逗号
//		}
//		pan.MKMulFloat("mo_tw", bet)
//		pan.MKMulFloat("apwa", bet)
//		pan.MKMulFloat("rw", bet)
//		psym := pan.Get("psym")
//		split := strings.Split(psym, "~")
//		if len(split) == 3 {
//			ret, _ := strconv.ParseFloat(split[1], 64)
//			p := message.NewPrinter(language.English)
//			split[1] = p.Sprintf("%.2f", ret*bet)
//			psym = strings.Join(split, "~")
//			pan.Set("psym", psym)
//		}
//		gsfA := pan.Get("gsf_a")
//		if gsfA != "" {
//			split = strings.Split(gsfA, ";")
//			for i, cell := range split {
//				cellsplit := strings.Split(cell, "~")
//				if len(cellsplit) != 2 {
//					break
//				}
//				ret, _ := strconv.ParseFloat(cellsplit[1], 64)
//				p := message.NewPrinter(language.English)
//				cellsplit[1] = p.Sprintf("%.2f", ret*bet)
//				split[i] = strings.Join(cellsplit, "~")
//			}
//			gsfA = strings.Join(split, ";")
//			pan.Set("gsf_a", gsfA)
//		}
//
//		// 消除
//		for k := range pan {
//			// 正则表达式，使用捕获组
//			r := regexp.MustCompile(`^(l\d+)`)
//			// 查找匹配
//			match := r.FindStringSubmatch(k)
//			// 判断是否匹配并获取结果
//			if len(match) > 0 {
//				lstrs := strings.Split(pan[k], "~")
//				lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
//				score, err := strconv.ParseFloat(lstrs[1], 64)
//				if err != nil {
//					panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, k))
//				}
//				p := message.NewPrinter(language.English)
//				lstrs[1] = p.Sprintf("%.2f", score*bet)
//				pan[k] = strings.Join(lstrs, "~")
//			}
//		}
//		//for i := 0; i < pan.Int("l"); i++ {
//		//	l := fmt.Sprintf("l%d", i)
//		//	if lstr, ok := pan[l]; ok {
//		//		lstrs := strings.Split(lstr, "~")
//		//		lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
//		//		score, err := strconv.ParseFloat(lstrs[1], 64)
//		//		if err != nil {
//		//			panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, l))
//		//		}
//		//		p := message.NewPrinter(language.English)
//		//		lstrs[1] = p.Sprintf("%.2f", score*bet)
//		//		pan[l] = strings.Join(lstrs, "~")
//		//	}
//		//}
//
//		pan.SetCurrency("balance", balance)
//		pan.SetCurrency("balance_cash", balance)
//
//		return pan
//	}
func parseMongoArray(raw interface{}) []map[string]interface{} {
	// 尝试断言为 bson.A，即 []interface{}
	array, ok := raw.(bson.A)
	if !ok {
		fmt.Println("类型断言失败，raw 不是 bson.A")
		return nil
	}

	var result []map[string]interface{}
	for i, item := range array {
		doc, ok := item.(bson.M) // bson.M == map[string]interface{}
		if !ok {
			fmt.Printf("第 %d 个元素不是 bson.M\n", i)
			continue
		}
		result = append(result, doc)
	}
	return result
}

func (sd *SimulateData) Deal2(multiple float64, balance int64, roundId string) Variables {
	pan := sd.DropPan[0]
	round := pan["round"].(Variables)
	accountBalance := pan["accountBalance"].(Variables)
	accountBalance["balance"] = balance
	round["roundId"] = roundId
	events, _ := convertMongoData(round["events"])
	for i := range events {
		events[i].MKMulFloat("wa", multiple)
		events[i].MKMulFloat("awa", multiple)
		c := convertToMapSafely(events[i]["c"])
		actions, _ := convertMongoData(c["actions"])
		for i2 := range actions {
			data := convertToMapSafely(actions[i2]["data"])
			data.MKMulFloat("winAmount", multiple)
			data.MKMulFloat("baseWinAmount", multiple)
			//winMultipliers, _ := convertMongoData(actions[i2]["winMultipliers"])
			//for i3 := range winMultipliers {
			//	winMultipliers[i3]
			//}
		}
	}

	//spinResult.MKMulFloat("totalWin", multiply)
	//pan.SetCurrency("balance", balance)
	pan["serverTime"] = time.Now().UTC().Format(time.RFC3339)
	return pan
}
