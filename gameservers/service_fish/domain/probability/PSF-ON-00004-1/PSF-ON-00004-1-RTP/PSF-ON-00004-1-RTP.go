package PSF_ON_00004_1_RTP

import (
	"reflect"
	PSF_ON_00004_1 "serve/service_fish/domain/probability/PSF-ON-00004-1"
)

func (s *service) RngFraction(drbMathI interface{}) (denominator, molecular, multiplier int) {
	drbMath := reflect.ValueOf(drbMathI).Interface().(*PSF_ON_00004_1.DrbMath)

	denominator = s.rngDenominator(drbMath)
	molecular = s.rngMolecular(drbMath)
	multiplier = s.getMultiplier(molecular, drbMath)

	molecular = denominator * molecular / multiplier
	return denominator, molecular, multiplier
}

func (s *service) RngRtp() (rtpIdOrder int) {
	return s.rngRtpId()
}

func (s *service) GetFirstMultiplier(subgameId int, drbMathI interface{}) int {
	drbMath := reflect.ValueOf(drbMathI).Interface().(*PSF_ON_00004_1.DrbMath)
	return s.firstMultiplier(subgameId, drbMath)
}
