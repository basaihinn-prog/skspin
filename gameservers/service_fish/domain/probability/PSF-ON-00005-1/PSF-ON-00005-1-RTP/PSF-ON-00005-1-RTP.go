package PSF_ON_00005_1_RTP

import (
	"reflect"
	PSF_ON_00005_1 "serve/service_fish/domain/probability/PSF-ON-00005-1"
)

func (s *service) RngRtp(state int, drbMathI interface{}) (rtpState int, rtpId string, rtpBullets uint64) {
	drbMath := reflect.ValueOf(drbMathI).Interface().(*PSF_ON_00005_1.DrbMath)

	rtpGroupId := -1
	rtpState, rtpGroupId = s.highLowState(state, drbMath)

	switch rtpGroupId {
	case 2:
		rtpId = drbMath.RTPGroup.RTP2.ID

		rtpBullets = s.rngBullets(
			drbMath.RTPGroup.RTP2.MinBullet,
			drbMath.RTPGroup.RTP2.MaxBullet,
		)

	case 3:
		rtpId = drbMath.RTPGroup.RTP3.ID

		rtpBullets = s.rngBullets(
			drbMath.RTPGroup.RTP3.MinBullet,
			drbMath.RTPGroup.RTP3.MaxBullet,
		)

	case 4:
		rtpId = drbMath.RTPGroup.RTP4.ID

		rtpBullets = s.rngBullets(
			drbMath.RTPGroup.RTP4.MinBullet,
			drbMath.RTPGroup.RTP4.MaxBullet,
		)

	case 5:
		rtpId = drbMath.RTPGroup.RTP5.ID

		rtpBullets = s.rngBullets(
			drbMath.RTPGroup.RTP5.MinBullet,
			drbMath.RTPGroup.RTP5.MaxBullet,
		)

	case 6:
		rtpId = drbMath.RTPGroup.RTP6.ID

		rtpBullets = s.rngBullets(
			drbMath.RTPGroup.RTP6.MinBullet,
			drbMath.RTPGroup.RTP6.MaxBullet,
		)

	case 7:
		rtpId = drbMath.RTPGroup.RTP7.ID

		rtpBullets = s.rngBullets(
			drbMath.RTPGroup.RTP7.MinBullet,
			drbMath.RTPGroup.RTP7.MaxBullet,
		)

	case 8:
		rtpId = drbMath.RTPGroup.RTP8.ID

		rtpBullets = s.rngBullets(
			drbMath.RTPGroup.RTP8.MinBullet,
			drbMath.RTPGroup.RTP8.MaxBullet,
		)

	case 9:
		rtpId = drbMath.RTPGroup.RTP9.ID

		rtpBullets = s.rngBullets(
			drbMath.RTPGroup.RTP9.MinBullet,
			drbMath.RTPGroup.RTP9.MaxBullet,
		)

	case 1:
		fallthrough
	default:
		rtpId = drbMath.RTPGroup.RTP1.ID

		rtpBullets = s.rngBullets(
			drbMath.RTPGroup.RTP1.MinBullet,
			drbMath.RTPGroup.RTP1.MaxBullet,
		)
	}

	return rtpState, rtpId, rtpBullets
}
