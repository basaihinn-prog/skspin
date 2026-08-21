package PSF_ON_00005_1_MONSTER

import "serve/fish_comm/rng"

const (
	SumHitWeight = 100000
)

var MONSTER = &monster{}

type monster struct {
}

func (m *monster) isHit(hitWeight []int) bool {
	m.checkHitWeight(hitWeight)

	options := make([]rng.Option, 0, 1)
	options = append(options, rng.Option{hitWeight[0], true})
	options = append(options, rng.Option{hitWeight[1], false})

	return m.rng(options).(bool)
}

func (m *monster) checkHitWeight(hitWeight []int) {
	if hitWeight[0]+hitWeight[1] != SumHitWeight {
		panic("hitWeight[0] + hitWeight[1] error")
	}
}

func (m *monster) rng(options []rng.Option) interface{} {
	return rng.New(options).Item
}
