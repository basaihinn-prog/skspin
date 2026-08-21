package jdbcomm

import (
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"slices"
	"sync"
)

type CombineDataBuy struct {
	RawSeries []primitive.ObjectID
	Series    []primitive.ObjectID
	Idx       int
	Mtx       sync.Mutex
}

func (data *CombineDataBuy) Next() (id primitive.ObjectID) {
	data.Mtx.Lock()
	defer data.Mtx.Unlock()

	if data.Idx >= len(data.Series) {
		data.Shuffle()
		data.Idx = 0
	}
	id = data.Series[data.Idx]
	data.Idx++
	return
}

func (data *CombineDataBuy) Shuffle() {
	data.Series = lo.Shuffle(slices.Clone(data.RawSeries))
}
