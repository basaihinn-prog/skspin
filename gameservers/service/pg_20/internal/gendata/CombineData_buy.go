package gendata

import (
	"slices"
	"sync"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CombineDataBuy struct {
	rawSeries []primitive.ObjectID
	series    []primitive.ObjectID
	idx       int
	mtx       sync.Mutex
}

func (data *CombineDataBuy) next() (id primitive.ObjectID) {
	data.mtx.Lock()
	defer data.mtx.Unlock()

	if data.idx >= len(data.series) {
		data.shuffle()
		data.idx = 0
	}
	id = data.series[data.idx]
	data.idx++
	return
}

func (data *CombineDataBuy) shuffle() {
	data.series = lo.Shuffle(slices.Clone(data.rawSeries))
}
