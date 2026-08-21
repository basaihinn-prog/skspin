package gendata

import "fmt"

type Bound struct {
	ID    int
	Need  int
	Group int

	Min, Max float64

	//需要额外扣除玩家个人奖池5Bet
	PoolCost int
}

func (b *Bound) name() string {
	s := fmt.Sprintf("(%v, %v]", b.Min, b.Max)

	return s
}

type Buckets struct {
	bounds []*Bound
}

var GBuckets = NewBuckets()

// （20，50】需要额外扣除玩家个人奖池3Bet
// （50，100】需要额外扣除玩家个人奖池5Bet
// （100，9999999】需要额外扣除玩家个人奖池10et

func NewBuckets() *Buckets {
	bounds := []*Bound{
		{
			Need:  1e4,
			Group: 0,
			Min:   -1,
			Max:   0,
		},
		{
			Need:  1e4,
			Group: 1,
			Min:   0,
			Max:   3,
		},
		{
			Need:  5e3,
			Group: 1,
			Min:   3,
			Max:   5,
		},
		{
			Need:  5e3,
			Group: 2,
			Min:   5,
			Max:   8,
		},
		{
			Need:  5e3,
			Group: 2,
			Min:   8,
			Max:   10,
		},
		{
			Need:  5e3,
			Group: 3,
			Min:   10,
			Max:   15,
		},
		{
			Need:     5e3,
			Group:    3,
			Min:      15,
			Max:      20,
			PoolCost: 0,
		},
		{
			Need:     5e3,
			Group:    4,
			Min:      20,
			Max:      50,
			PoolCost: 3,
		},
		{
			Need:     5e3,
			Group:    4,
			Min:      50,
			Max:      100,
			PoolCost: 5,
		},
		{
			Need:     5e3,
			Group:    5,
			Min:      100,
			Max:      9999999,
			PoolCost: 10,
		},
	}

	b := Buckets{
		bounds: bounds,
	}

	for i, v := range b.bounds {
		v.ID = i
	}

	return &b
}

func (b *Buckets) GetBucket(multi float64) int {
	for _, b := range b.bounds {
		if b.Min < multi && multi <= b.Max {
			return b.ID
		}
	}

	return -1

}

func (b *Buckets) GetBound(i int) *Bound {
	return b.bounds[i]
}
