package gendata

import (
	"log/slog"
	"serve/comm/db"
	"serve/comm/redisx"
	"serve/comm/ut"
	"slices"
	"sync"

	"github.com/samber/lo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CombineData struct {
	// raw data
	rtp     float64
	combine Combine
	buckets map[int][][]primitive.ObjectID

	combineUseCounts []int

	series map[int][]primitive.ObjectID
	idx    map[int]int
	mtx    sync.Mutex
}

func (data *CombineData) next(gameType db.BoundType, WitchBucket int) (id primitive.ObjectID) {
	data.mtx.Lock()
	defer data.mtx.Unlock()
	if data.idx == nil {
		data.idx = make(map[int]int)
	}
	if data.idx[WitchBucket] >= len(data.series[WitchBucket]) {
		data.shuffle(gameType, WitchBucket)
		data.idx[WitchBucket] = 0
	}

	id = data.series[WitchBucket][data.idx[WitchBucket]]
	data.idx[WitchBucket]++
	return
}

func (data *CombineData) shuffle(gameType db.BoundType, bucketId int) {
	if bucketId == BucketGov {
		var buckets = make([][]primitive.ObjectID, len(data.buckets[bucketId]))
		for j := 0; j < len(buckets); j++ {
			buckets[j] = slices.Clone(data.buckets[bucketId][j])
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
		for _, c := range data.combine {
			if c.Type != gameType {
				continue
			}
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
		data.combineUseCounts = make([]int, len(oldcountmap))
		for i := 0; i < len(oldcountmap); i++ {
			data.combineUseCounts[i] = oldcountmap[i] - newcountmap[i]
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
	batchSize := redisx.ALLGameBatchSize
	var buckets = make([][]primitive.ObjectID, len(data.buckets[bucketId]))
	for j := 0; j < len(buckets); j++ {
		buckets[j] = slices.Clone(data.buckets[bucketId][j])
	}
	// 将所有 buckets 中的 ID 拉平到一个切片中
	allIDs := lo.Flatten(data.buckets[bucketId])
	perBatch := len(allIDs) / batchSize
	remainder := len(allIDs) % batchSize
	if remainder > 0 {
		perBatch += 1
	}
	// 创建一个切片来保存每批次的数据
	pending := make([][]primitive.ObjectID, perBatch)
	//start := 0
	minFunc := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	// 分配每个批次的数据，并处理剩余部分
	for i := 0; i < perBatch; i++ {
		for j := 0; j < len(buckets); j++ {
			if buckets[j] != nil {
				bucketLen := len(buckets[j])
				batchLen := bucketLen / perBatch // 每个批次的基本长度
				extra := 0
				// 判断是否需要在当前批次分配额外的一个元素
				if i < remainder {
					extra = 1
				}
				start := i*batchLen + minFunc(i, bucketLen%perBatch)
				end := start + batchLen + extra
				if end > bucketLen { // 防止越界
					end = bucketLen
				}
				if start < bucketLen { // 确保 start 不越界
					pending[i] = append(pending[i], buckets[j][start:end]...)
				}
			}
		}
	}

	// 打乱每个批次内的数据
	for i := range pending {
		lo.Shuffle(pending[i])
	}

	if data.series == nil {
		data.series = make(map[int][]primitive.ObjectID)
	}
	// 将乱序后的 ID 保存到 data.series
	data.series[bucketId] = lo.Flatten(pending)

	return
}
