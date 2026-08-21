package PSF_ON_00001_1_MONSTER

var Envelope = &envelope{}

type envelope struct{}

func (e *envelope) Pick(math *BsRedEnvelopeMath) (iconPays []int) {
	return RedEnvelope.pick(math)
}
