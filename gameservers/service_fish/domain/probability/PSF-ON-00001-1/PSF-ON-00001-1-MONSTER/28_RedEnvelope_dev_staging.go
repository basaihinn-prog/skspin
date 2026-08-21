//go:build dev || staging
// +build dev staging

package PSF_ON_00001_1_MONSTER

func (r *redEnvelope) Hit(math *BsRedEnvelopeMath) (iconPays []int) {
	if MONSTER.isHit(math.HitWeight) {
		return r.pick(math)
	}
	return nil
}
