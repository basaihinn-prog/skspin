package facaicomm

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"reflect"
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

func (sd *SimulateData) Deal2(multiple float64, balance int64, roundId string) Variables {
	pan := sd.DropPan[0]
	fmt.Println(pan)

	pan.MKMulFloat("totalWinBase", multiple)
	pan.MKMulFloat("totalWin", multiple)
	pan.MKMulFloat("totalWinFree", multiple)
	pan["timestamp"] = time.Now().UnixNano()
	pan.SetInt("userPoint", int(balance))
	//pan.SetInt("winOdds", 0)
	//pan.SetInt("userPls", 1)
	pan.MKMulFloat("normalBroadcastAmount", multiple)
	onceResults, _ := convertMongoData(pan["onceResults"])
	for i := range onceResults {
		onceResults[i].MKMulFloat("winCoins", multiple)
		fallenResults, _ := convertMongoData(onceResults[i]["fallenResults"])
		for i2 := range fallenResults {
			fallenResults[i2].MKMulFloat("beforeMultipleWinCoins", multiple)
			fallenResults[i2].MKMulFloat("fallenWinCoins", multiple)
			normalDetail, _ := convertMongoData(fallenResults[i2]["normalDetail"])
			for i3 := range normalDetail {
				normalDetail[i3].MKMulFloat("winCoins", multiple)

			}
		}
	}
	fmt.Println(pan)
	return pan
}
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

func truncateFloat(num float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(num*shift) / shift
}

func convertToSlice(data any) ([]any, error) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return nil, fmt.Errorf("input is not a slice or array")
	}

	result := make([]any, val.Len())
	for i := 0; i < val.Len(); i++ {
		result[i] = val.Index(i).Interface()
	}
	return result, nil
}
