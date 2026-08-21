package ppcomm

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type SimulateData struct {
	Id              primitive.ObjectID `bson:"_id"`
	DropPan         []Variables        `bson:"droppan"` //自行解析的数据
	HasGame         bool               `bson:"hasgame"`
	Times           float64            `bson:"times"`
	BucketId        int                `bson:"bucketid"`
	Type            int                `bson:"type"`
	Selected        bool               `bson:"selected"`
	RoundID         int                `bson:"RoundID"`     //数据ID
	TurnIndex       int                `bson:"TurnIndex"`   //轮次，他跟数据ID是属于一组的数据
	FreeFlag        int                `bson:"FreeFlag"`    //旋转类型   0 就是普通旋转, 1就是买免费  2  超级购买   1010代表10次购买   1012代表12次  2010  也是代表10次购买
	GroupFlag       int                `bson:"GroupFlag"`   //普通下注，双倍下注，部分游戏有双倍下注
	QueryString     Variables          `bson:"QueryString"` //数据值
	BucketMix       int                `bson:"BucketMix"`
	BucketStable    int                `bson:"BucketStable"`
	BucketHeartBeat int                `bson:"BucketHeartBeat"`
}

func (sd *SimulateData) Deal(c float64, line float64, balance float64) Variables {
	bet := c * line
	//slog.Info("SimulateData, ", "ObjectID", sd.Id.Hex())

	pan := sd.DropPan[0]

	rid := strings.Split(pan.Str("gid"), "_")[0]
	pan.Set("rid", rid)
	pan.SetFloat("c", c)
	pan.MKMulFloat("tw", bet)
	pan.MKMulFloat("tmb_win", bet)
	pan.MKMulFloat("w", bet)
	pan.MKMulFloat("wp", bet)
	pan.MKMulFloat("fsres", bet)
	pan.MKMulFloat("fswin", bet)
	pan.MKMulFloat("fsres_total", bet)
	pan.MKMulFloat("fswin_total", bet)
	pan.MKMulFloat("tmb_res", bet)
	pan.MKMulFloat("rs_iw", bet)
	if pan.Get("rs_iw") != "" {
		pan.Set("rs_iw", strings.ReplaceAll(pan.Get("rs_iw"), ",", "")) //去除逗号
	}
	pan.MKMulFloat("rs_win", bet)
	if pan.Get("rs_win") != "" {
		pan.Set("rs_win", strings.ReplaceAll(pan.Get("rs_win"), ",", "")) //去除逗号
	}
	pan.MKMulFloat("mo_tw", bet)
	pan.MKMulFloat("apwa", bet)
	pan.MKMulFloat("rw", bet)
	psym := pan.Get("psym")
	split := strings.Split(psym, "~")
	if len(split) == 3 {
		ret, _ := strconv.ParseFloat(split[1], 64)
		p := message.NewPrinter(language.English)
		split[1] = p.Sprintf("%.2f", ret*bet)
		psym = strings.Join(split, "~")
		pan.Set("psym", psym)
	}
	gsfA := pan.Get("gsf_a")
	if gsfA != "" {
		split = strings.Split(gsfA, ";")
		for i, cell := range split {
			cellsplit := strings.Split(cell, "~")
			if len(cellsplit) != 2 {
				break
			}
			ret, _ := strconv.ParseFloat(cellsplit[1], 64)
			p := message.NewPrinter(language.English)
			cellsplit[1] = p.Sprintf("%.2f", ret*bet)
			split[i] = strings.Join(cellsplit, "~")
		}
		gsfA = strings.Join(split, ";")
		pan.Set("gsf_a", gsfA)
	}

	// 消除
	for k := range pan {
		// 正则表达式，使用捕获组
		r := regexp.MustCompile(`^(l\d+)`)
		// 查找匹配
		match := r.FindStringSubmatch(k)
		// 判断是否匹配并获取结果
		if len(match) > 0 {
			lstrs := strings.Split(pan[k], "~")
			lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
			score, err := strconv.ParseFloat(lstrs[1], 64)
			if err != nil {
				panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, k))
			}
			p := message.NewPrinter(language.English)
			lstrs[1] = p.Sprintf("%.2f", score*bet)
			pan[k] = strings.Join(lstrs, "~")
		}
	}
	//for i := 0; i < pan.Int("l"); i++ {
	//	l := fmt.Sprintf("l%d", i)
	//	if lstr, ok := pan[l]; ok {
	//		lstrs := strings.Split(lstr, "~")
	//		lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
	//		score, err := strconv.ParseFloat(lstrs[1], 64)
	//		if err != nil {
	//			panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, l))
	//		}
	//		p := message.NewPrinter(language.English)
	//		lstrs[1] = p.Sprintf("%.2f", score*bet)
	//		pan[l] = strings.Join(lstrs, "~")
	//	}
	//}

	pan.SetCurrency("balance", balance)
	pan.SetCurrency("balance_cash", balance)

	return pan
}

func (sd *SimulateData) Deal2(c float64, line float64, balance float64) Variables {
	bet := c * line
	//slog.Info("SimulateData, ", "ObjectID", sd.Id.Hex())
	pan := sd.DropPan[0]
	originC := pan.Float("c")
	multiply := c / originC // 倍数，所有与钱相关的都在入库的时候乘以了100转换为整数
	bet = multiply
	rid := strings.Split(pan.Str("gid"), "_")[0]
	pan.Set("rid", rid)
	pan.SetFloat("c", c)
	pan.MKMulFloat("tw", bet)
	pan.MKMulFloat("tmb_win", bet)
	pan.MKMulFloat("w", bet)
	pan.MKMulFloat("wp", bet)
	pan.MKMulFloat("fsres", bet)
	pan.MKMulFloat("fswin", bet)
	pan.MKMulFloat("fsres_total", bet)
	pan.MKMulFloat("fswin_total", bet)
	pan.MKMulFloat("tmb_res", bet)
	pan.MKMulFloat("pw", bet)
	pan.MKMulFloat("rs_iw", bet)
	if pan.Get("rs_iw") != "" {
		pan.Set("rs_iw", strings.ReplaceAll(pan.Get("rs_iw"), ",", "")) //去除逗号
	}
	pan.MKMulFloat("rs_win", bet)
	if pan.Get("rs_win") != "" {
		pan.Set("rs_win", strings.ReplaceAll(pan.Get("rs_win"), ",", "")) //去除逗号
	}
	pan.MKMulFloat("mo_tw", bet)
	//针对pp_vswayschilheat的处理
	pan.MKMulmo_twInG("mo_tw", bet)
	//pan.MKMulFloat("apwa", bet)
	if pan.Get("apwa") != "" {
		splitAPWA := strings.Split(pan.Get("apwa"), ",")
		for k, v := range splitAPWA {
			ret, _ := strconv.ParseFloat(v, 64)
			p := message.NewPrinter(language.English)
			splitAPWA[k] = p.Sprintf("%.2f", ret*bet)
		}
		pan.Set("apwa", strings.Join(splitAPWA, ","))
	}
	pan.MKMulFloat("rw", bet)
	psym := pan.Get("psym")
	split := strings.Split(psym, "~")
	if len(split) == 3 {
		ret, _ := strconv.ParseFloat(split[1], 64)
		p := message.NewPrinter(language.English)
		split[1] = p.Sprintf("%.2f", ret*bet)
		psym = strings.Join(split, "~")
		pan.Set("psym", psym)
	}
	gsfA := pan.Get("gsf_a")
	if gsfA != "" {
		split = strings.Split(gsfA, ";")
		for i, cell := range split {
			cellsplit := strings.Split(cell, "~")
			if len(cellsplit) != 2 {
				break
			}
			ret, _ := strconv.ParseFloat(cellsplit[1], 64)
			p := message.NewPrinter(language.English)
			cellsplit[1] = p.Sprintf("%.2f", ret*bet)
			split[i] = strings.Join(cellsplit, "~")
		}
		gsfA = strings.Join(split, ";")
		pan.Set("gsf_a", gsfA)
	}
	//ReplaceEmbedScoreField 兴许可以使用
	//ReplaceScoreField2(pan, bet, "wlc_v", '~')
	if _, ok := pan["wlc_v"]; ok {
		wlc_vStr := pan.Get("wlc_v")
		wlc_v := strings.Split(wlc_vStr, ";")
		for i := range wlc_v {
			temp := strings.Split(wlc_v[i], "~")
			if len(temp) > 2 {
				temp[1] = strings.ReplaceAll(temp[1], ",", "")
				ret, _ := strconv.ParseFloat(temp[1], 64)
				p := message.NewPrinter(language.English)
				temp[1] = p.Sprintf("%.2f", ret*bet)
				wlc_v[i] = strings.Join(temp, "~")
				fmt.Println(wlc_v[i])
			}
		}
		wlc_vStr = strings.Join(wlc_v, ";")
		pan.Set("wlc_v", wlc_vStr)
	}
	// 消除
	for k := range pan {
		// 正则表达式，使用捕获组
		r := regexp.MustCompile(`^(l\d+)`)
		// 查找匹配
		match := r.FindStringSubmatch(k)
		// 判断是否匹配并获取结果
		if len(match) > 0 {
			lstrs := strings.Split(pan[k], "~")
			lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
			score, err := strconv.ParseFloat(lstrs[1], 64)
			if err != nil {
				panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, k))
			}
			p := message.NewPrinter(language.English)
			lstrs[1] = p.Sprintf("%.2f", score*bet)
			pan[k] = strings.Join(lstrs, "~")
		}

	}
	//for i := 0; i < pan.Int("l"); i++ {
	//	l := fmt.Sprintf("l%d", i)
	//	if lstr, ok := pan[l]; ok {
	//		lstrs := strings.Split(lstr, "~")
	//		lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
	//		score, err := strconv.ParseFloat(lstrs[1], 64)
	//		if err != nil {
	//			panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, l))
	//		}
	//		p := message.NewPrinter(language.English)
	//		lstrs[1] = p.Sprintf("%.2f", score*bet)
	//		pan[l] = strings.Join(lstrs, "~")
	//	}
	//}

	pan.SetCurrency("balance", balance)
	pan.SetCurrency("balance_cash", balance)

	return pan
}

func (sd *SimulateData) DealSP(c float64, line float64, balance float64, originC float64) Variables {
	bet := c * line
	//slog.Info("SimulateData, ", "ObjectID", sd.Id.Hex())
	pan := sd.DropPan[0]
	//originC := pan.Float("c")
	multiply := c / originC // 倍数，所有与钱相关的都在入库的时候乘以了100转换为整数
	bet = multiply
	rid := strings.Split(pan.Str("gid"), "_")[0]
	pan.Set("rid", rid)
	pan.SetFloat("c", c)
	pan.MKMulFloat("tw", bet)
	pan.MKMulFloat("tmb_win", bet)
	pan.MKMulFloat("w", bet)
	pan.MKMulFloat("wp", bet)
	pan.MKMulFloat("fsres", bet)
	pan.MKMulFloat("fswin", bet)
	pan.MKMulFloat("fsres_total", bet)
	pan.MKMulFloat("fswin_total", bet)
	pan.MKMulFloat("tmb_res", bet)
	pan.MKMulFloat("pw", bet)
	pan.MKMulFloat("rs_iw", bet)
	if pan.Get("rs_iw") != "" {
		pan.Set("rs_iw", strings.ReplaceAll(pan.Get("rs_iw"), ",", "")) //去除逗号
	}
	pan.MKMulFloat("rs_win", bet)
	if pan.Get("rs_win") != "" {
		pan.Set("rs_win", strings.ReplaceAll(pan.Get("rs_win"), ",", "")) //去除逗号
	}
	pan.MKMulFloat("mo_tw", bet)
	//针对pp_vswayschilheat的处理
	pan.MKMulmo_twInG("mo_tw", bet)
	//pan.MKMulFloat("apwa", bet)
	if pan.Get("apwa") != "" {
		splitAPWA := strings.Split(pan.Get("apwa"), ",")
		for k, v := range splitAPWA {
			ret, _ := strconv.ParseFloat(v, 64)
			p := message.NewPrinter(language.English)
			splitAPWA[k] = p.Sprintf("%.2f", ret*bet)
		}
		pan.Set("apwa", strings.Join(splitAPWA, ","))
	}
	pan.MKMulFloat("rw", bet)
	psym := pan.Get("psym")
	split := strings.Split(psym, "~")
	if len(split) == 3 {
		ret, _ := strconv.ParseFloat(split[1], 64)
		p := message.NewPrinter(language.English)
		split[1] = p.Sprintf("%.2f", ret*bet)
		psym = strings.Join(split, "~")
		pan.Set("psym", psym)
	}
	gsfA := pan.Get("gsf_a")
	if gsfA != "" {
		split = strings.Split(gsfA, ";")
		for i, cell := range split {
			cellsplit := strings.Split(cell, "~")
			if len(cellsplit) != 2 {
				break
			}
			ret, _ := strconv.ParseFloat(cellsplit[1], 64)
			p := message.NewPrinter(language.English)
			cellsplit[1] = p.Sprintf("%.2f", ret*bet)
			split[i] = strings.Join(cellsplit, "~")
		}
		gsfA = strings.Join(split, ";")
		pan.Set("gsf_a", gsfA)
	}
	//ReplaceEmbedScoreField 兴许可以使用
	//ReplaceScoreField2(pan, bet, "wlc_v", '~')
	if _, ok := pan["wlc_v"]; ok {
		wlc_vStr := pan.Get("wlc_v")
		wlc_v := strings.Split(wlc_vStr, ";")
		for i := range wlc_v {
			temp := strings.Split(wlc_v[i], "~")
			if len(temp) > 2 {
				temp[1] = strings.ReplaceAll(temp[1], ",", "")
				ret, _ := strconv.ParseFloat(temp[1], 64)
				p := message.NewPrinter(language.English)
				temp[1] = p.Sprintf("%.2f", ret*bet)
				wlc_v[i] = strings.Join(temp, "~")
				fmt.Println(wlc_v[i])
			}
		}
		wlc_vStr = strings.Join(wlc_v, ";")
		pan.Set("wlc_v", wlc_vStr)
	}
	// 消除
	for k := range pan {
		// 正则表达式，使用捕获组
		r := regexp.MustCompile(`^(l\d+)`)
		// 查找匹配
		match := r.FindStringSubmatch(k)
		// 判断是否匹配并获取结果
		if len(match) > 0 {
			lstrs := strings.Split(pan[k], "~")
			lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
			score, err := strconv.ParseFloat(lstrs[1], 64)
			if err != nil {
				panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, k))
			}
			p := message.NewPrinter(language.English)
			lstrs[1] = p.Sprintf("%.2f", score*bet)
			pan[k] = strings.Join(lstrs, "~")
		}

	}
	//for i := 0; i < pan.Int("l"); i++ {
	//	l := fmt.Sprintf("l%d", i)
	//	if lstr, ok := pan[l]; ok {
	//		lstrs := strings.Split(lstr, "~")
	//		lstrs[1] = strings.ReplaceAll(lstrs[1], ",", "")
	//		score, err := strconv.ParseFloat(lstrs[1], 64)
	//		if err != nil {
	//			panic(fmt.Sprintf("deal score error, rid:%s, line:%s", rid, l))
	//		}
	//		p := message.NewPrinter(language.English)
	//		lstrs[1] = p.Sprintf("%.2f", score*bet)
	//		pan[l] = strings.Join(lstrs, "~")
	//	}
	//}

	pan.SetCurrency("balance", balance)
	pan.SetCurrency("balance_cash", balance)

	return pan
}

// pp的数据进行解析成返回前端的格式
func (sd *SimulateData) ParsingData(bet int64, balance float64, multiply int, enterBalance float64) Variables {
	pan := sd.QueryString

	//处理嵌入字段的分数*倍数
	ReplaceEmbedScoreField(pan, multiply)

	// 添加余额
	AppendBalance(pan, bet, balance, multiply, enterBalance)

	// 替换其他字段
	ReplaceAllScoreField(pan, multiply)

	pan.SetCurrency("balance", balance)
	pan.SetCurrency("balance_cash", balance)
	pan.SetCurrency("balance_bonus", 0.00)

	pan.SetInt("stime", int(time.Now().UnixMilli()))

	return pan
}

func GetBet(isBuyFree bool, coin float64, lines int, bl string, hasDouble int, FreeFlag int, TurnIndex int) int64 {
	if isBuyFree {
		bet := coin * float64(lines) // 0.01 * 20 = 0.20
		// 待观察
		bet *= 100
		return int64(bet)
	}

	if FreeFlag == 0 && TurnIndex == 0 {
		if hasDouble == 1 && bl == "1" { // 加倍模式
			return int64(coin * float64(lines) * 1.25)
		}
		return int64(coin * float64(lines))
	}

	return 0
}

func parseToMap(query string) (Variables, error) {
	// 解析查询字符串
	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

	// 转换为 map[string]string
	result := make(map[string]string)
	for key, val := range values {
		// url.ParseQuery 会返回 []string，即使每个键只有一个值
		if len(val) > 0 {
			result[key] = val[0]
		}
	}

	return result, nil
}
