package hacksawcomm

import (
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"serve/comm/ut"
	"slices"
	"sync"
)

type CombineData struct {
	// raw data
	Combine Combine
	Buckets map[int][][]primitive.ObjectID

	CombineUseCounts []int

	series map[int][]primitive.ObjectID
	idx    map[int]int
	mtx    sync.Mutex
}

func (data *CombineData) Next(WitchBucket int) (id primitive.ObjectID) {
	data.mtx.Lock()
	defer data.mtx.Unlock()
	if data.idx == nil {
		data.idx = make(map[int]int)
	}
	if data.idx[WitchBucket] >= len(data.series[WitchBucket]) {
		data.Shuffle(WitchBucket)
		data.idx[WitchBucket] = 0
	}

	id = data.series[WitchBucket][data.idx[WitchBucket]]
	data.idx[WitchBucket]++
	return
}

func (data *CombineData) Shuffle(bucketId int) {
	var buckets = make([][]primitive.ObjectID, len(data.Buckets[bucketId]))
	for j := 0; j < len(buckets); j++ {
		buckets[j] = slices.Clone(data.Buckets[bucketId][j])
	}

	oldcountmap := lo.Map(buckets, func(arr []primitive.ObjectID, _ int) int {
		return len(arr)
	})

	var pop = func(ids []int, count int) []primitive.ObjectID {
		ans := make([]primitive.ObjectID, 0, count)

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

	var pending [][]primitive.ObjectID
	for _, c := range data.Combine {
		for i := 0; i < c.Count; i++ {
			// 随机基本单元
			var idSeries []primitive.ObjectID
			for _, m := range c.Meta {
				// lo.Must0(len(m.BucketIds) != 0)
				if len(m.BucketIds) != 0 {
					picks := pop(m.BucketIds, m.Count)
					idSeries = append(idSeries, picks...)
				} else {
					// 待填充
					picks := make([]primitive.ObjectID, m.Count)
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
			if id.IsZero() {
				// 用剩下的数据填充
				ids := pop(freeIds, 1)
				if len(ids) <= 0 {
					continue
				}
				a[i] = ids[0]
			}
		}
		a = slices.DeleteFunc(a, func(id primitive.ObjectID) bool {
			return id.IsZero()
		})
	}

	newcountmap := lo.Map(buckets, func(arr []primitive.ObjectID, _ int) int {
		return len(arr)
	})

	data.CombineUseCounts = make([]int, len(oldcountmap))
	for i := 0; i < len(oldcountmap); i++ {
		data.CombineUseCounts[i] = oldcountmap[i] - newcountmap[i]
	}

	freeBucket := lo.Flatten(buckets)

	// 将组合外的数据和组合洗牌
	freePending := lo.Map(freeBucket, func(item primitive.ObjectID, _ int) []primitive.ObjectID {
		return []primitive.ObjectID{item}
	})

	pending = append(pending, freePending...)
	lo.Shuffle(pending)
	if data.series == nil {
		data.series = make(map[int][]primitive.ObjectID)
	}
	data.series[bucketId] = lo.Flatten(pending)

	return
}
