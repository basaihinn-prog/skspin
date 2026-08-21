package gendata

import (
	"log/slog"
	"serve/comm/db"
	"serve/comm/ut"
	"slices"
	"sync"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CombineData struct {
	// raw data
	combine          Combine
	combineUseCounts []int
	mtx              sync.Mutex

	//修改
	series  map[int][]db.Playitem
	buckets map[int][][]db.Playitem
	idx     map[int]int
}

func (data *CombineData) next(WitchBucket int) (now primitive.ObjectID, nextIds []primitive.ObjectID) {
	data.mtx.Lock()
	defer data.mtx.Unlock()

	if data.idx == nil {
		data.idx = make(map[int]int)
	}

	if data.idx[WitchBucket] >= len(data.series[WitchBucket]) {
		data.shuffle(WitchBucket)
		data.idx[WitchBucket] = 0
	}
	da := data.series[WitchBucket][data.idx[WitchBucket]]
	//fmt.Print("doc是: ")
	//fmt.Println(da.ID.Hex())
	//fmt.Println(fmt.Sprintf("我的下标是: %v", data.idx[WitchBucket]))
	data.idx[WitchBucket]++
	var i = 0
	hasBigReward := false
	for ; i < 9; i++ {
		if data.idx[WitchBucket]+i >= len(data.series[WitchBucket]) {
			break
		}
		//if i == 2 {
		//	data.series[WitchBucket][i+data.idx[WitchBucket]].ID, _ = primitive.ObjectIDFromHex("6679c5f25e1bab741e9bb753")
		//	data.series[WitchBucket][i+data.idx[WitchBucket]].Data.GameType = 2
		//}
		//fmt.Print("nextIds添加: ")
		//fmt.Println(data.series[WitchBucket][i+data.idx[WitchBucket]].ID.Hex())
		//fmt.Println(fmt.Sprintf("我的下标是: %v", i+data.idx[WitchBucket]))
		nextIds = append(nextIds, data.series[WitchBucket][i+data.idx[WitchBucket]].ID)
		if data.series[WitchBucket][i+data.idx[WitchBucket]].Data.GameType != 1 {
			hasBigReward = true
			break
		}
	}
	if hasBigReward {
		i++
		data.idx[WitchBucket] += i
	} else {
		nextIds = []primitive.ObjectID{}
	}
	//fmt.Println(fmt.Sprintf("本次取桶完毕, 现在ind是: %v", data.idx[WitchBucket]))

	now = da.ID
	return
}

func (data *CombineData) shuffle(bucketId int) {
	var buckets = make([][]db.Playitem, len(data.buckets[bucketId]))
	for i := 0; i < len(buckets); i++ {
		buckets[i] = slices.Clone(data.buckets[bucketId][i])
	}

	oldcountmap := lo.Map(buckets, func(arr []db.Playitem, _ int) int {
		return len(arr)
	})

	var pop = func(ids []int, count int) []db.Playitem {
		ans := make([]db.Playitem, 0, count)

		for i := 0; i < count; i++ {
			id, ok := ut.SampleByWeightsPred(ids, func(item int) (weight int) {
				return len(buckets[item])
			})

			if ok {
				ans = append(ans, ut.PopRand(&buckets[id]))
			} else {
				slog.Error("cannot found valid item", "ids", ids, "i", i, "count", count)
			}
		}

		return ans
	}

	var pending [][]db.Playitem
	for _, c := range data.combine {
		for i := 0; i < c.Count; i++ {
			// 随机基本单元
			var idSeries []db.Playitem
			for _, m := range c.Meta {
				// lo.Must0(len(m.BucketIds) != 0)
				if len(m.BucketIds) != 0 {
					picks := pop(m.BucketIds, m.Count)
					idSeries = append(idSeries, picks...)
				} else {
					// 待填充
					picks := make([]db.Playitem, m.Count)
					idSeries = append(idSeries, picks...)
				}
			}

			lo.Shuffle(idSeries)
			pending = append(pending, idSeries)
		}
		slog.Info("gen success", "combineId", c.ID, "count", c.Count, "name", c.Name)
	}

	freeIds := lo.Range(len(buckets))

	for _, a := range pending {
		for i, id := range a {
			if id.ID.IsZero() {
				// 用剩下的数据填充
				ids := pop(freeIds, 1)
				if len(ids) <= 0 {
					continue
				}
				a[i] = ids[0]
			}
		}
		a = slices.DeleteFunc(a, func(id db.Playitem) bool {
			return id.ID.IsZero()
		})
	}

	newcountmap := lo.Map(buckets, func(arr []db.Playitem, _ int) int {
		return len(arr)
	})

	data.combineUseCounts = make([]int, len(oldcountmap))
	for i := 0; i < len(oldcountmap); i++ {
		data.combineUseCounts[i] = oldcountmap[i] - newcountmap[i]
	}

	freeBucket := lo.Flatten(buckets)

	// 将组合外的数据和组合洗牌
	freePending := lo.Map(freeBucket, func(item db.Playitem, _ int) []db.Playitem {
		return []db.Playitem{item}
	})

	pending = append(pending, freePending...)
	lo.Shuffle(pending)
	if data.series == nil {
		data.series = make(map[int][]db.Playitem)
	}

	data.series[bucketId] = lo.Flatten(pending)

	return
}
