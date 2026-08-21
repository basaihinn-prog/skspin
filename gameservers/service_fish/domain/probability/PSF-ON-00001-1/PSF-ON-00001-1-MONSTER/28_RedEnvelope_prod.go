//go:build prod
// +build prod

package PSF_ON_00001_1_MONSTER

func (r *redEnvelope) Hit(math *BsRedEnvelopeMath) (iconPays []int) {
	if MONSTER.isHit(math.HitWeight) {
		return r.pick(math)
	}
	return nil
}
