package PSF_ON_00003_1_RTP

import (
	"serve/fish_comm/rng"
	PSF_ON_00003_1 "serve/service_fish/domain/probability/PSF-ON-00003-1"
	"strconv"
)

const (
	threw = 0
	bite  = 1
)

var Service = &service{}

type service struct {
}

func (s *service) rngWinLoss(state int, drbMath *PSF_ON_00003_1.DrbMath) int {
	switch state {
	case threw:
		return bite
	case bite:
		return threw
	default:
		options := make([]rng.Option, 0, 2)

		options = append(options, rng.Option{drbMath.NetWinGroup.InitialWeight[0], threw})
		options = append(options, rng.Option{drbMath.NetWinGroup.InitialWeight[1], bite})

		return s.rng(options).(int)
	}
}

func (s *service) rngWinLossLevel(initial int, drbMath *PSF_ON_00003_1.DrbMath) int {
	switch initial {
	case bite:
		return s.getWinOrLoss(drbMath.NetWinGroup.NetLoss, drbMath.NetWinGroup.NetLossWeight)

	case threw:
		fallthrough
	default:
		return s.getWinOrLoss(drbMath.NetWinGroup.NetWin, drbMath.NetWinGroup.NetWinWeight)
	}
}

func (s *service) getBullet(initial, netWinOrLoss int, drbMath *PSF_ON_00003_1.DrbMath) (rtpId string, bullets uint64) {
	switch initial {
	case threw:
		rtpId, bullets = s.getWinRTP(netWinOrLoss, &drbMath.NetWinGroup.NetWinRTP)
	case bite:
		rtpId, bullets = s.getLossRTP(netWinOrLoss, &drbMath.NetWinGroup.NetLossRTP)
	default:
		return "", 0
	}

	return rtpId, bullets
}

func (s *service) getWinRTP(netWinOrLoss int, winMath *struct{ PSF_ON_00003_1.DrbWinMath }) (rtpId string, bullets uint64) {
	var winRTP int

	switch netWinOrLoss {
	case winMath.NetWin1.Win[0]:
		winRTP = s.getWinOrLoss(winMath.NetWin1.RTPGroup, winMath.NetWin1.Weight)
		rtpId = strconv.Itoa(winRTP)
		bullets = s.getWinRtpBullet(rtpId, &winMath.NetWin1)

	case winMath.NetWin2.Win[0]:
		winRTP = s.getWinOrLoss(winMath.NetWin2.RTPGroup, winMath.NetWin2.Weight)
		rtpId = strconv.Itoa(winRTP)
		bullets = s.getWinRtpBullet(rtpId, &winMath.NetWin2)

	case winMath.NetWin3.Win[0]:
		winRTP = s.getWinOrLoss(winMath.NetWin3.RTPGroup, winMath.NetWin3.Weight)
		rtpId = strconv.Itoa(winRTP)
		bullets = s.getWinRtpBullet(rtpId, &winMath.NetWin3)

	case winMath.NetWin4.Win[0]:
		winRTP = s.getWinOrLoss(winMath.NetWin4.RTPGroup, winMath.NetWin4.Weight)
		rtpId = strconv.Itoa(winRTP)
		bullets = s.getWinRtpBullet(rtpId, &winMath.NetWin4)

	case winMath.NetWin5.Win[0]:
		winRTP = s.getWinOrLoss(winMath.NetWin5.RTPGroup, winMath.NetWin5.Weight)
		rtpId = strconv.Itoa(winRTP)
		bullets = s.getWinRtpBullet(rtpId, &winMath.NetWin5)

	default:
		return "", 0
	}

	return rtpId, bullets
}

func (s *service) getLossRTP(netWinOrLoss int, lossMath *struct{ PSF_ON_00003_1.DrbLossMath }) (rtpId string, bullets uint64) {
	var lossRTP int

	switch netWinOrLoss {
	case lossMath.NetLoss1.Loss[0]:
		lossRTP = s.getWinOrLoss(lossMath.NetLoss1.RTPGroup, lossMath.NetLoss1.Weight)
		rtpId = strconv.Itoa(lossRTP)
		bullets = s.getLossRtpBullet(rtpId, &lossMath.NetLoss1)

	case lossMath.NetLoss2.Loss[0]:
		lossRTP = s.getWinOrLoss(lossMath.NetLoss2.RTPGroup, lossMath.NetLoss2.Weight)
		rtpId = strconv.Itoa(lossRTP)
		bullets = s.getLossRtpBullet(rtpId, &lossMath.NetLoss2)

	case lossMath.NetLoss3.Loss[0]:
		lossRTP = s.getWinOrLoss(lossMath.NetLoss3.RTPGroup, lossMath.NetLoss3.Weight)
		rtpId = strconv.Itoa(lossRTP)
		bullets = s.getLossRtpBullet(rtpId, &lossMath.NetLoss3)

	case lossMath.NetLoss4.Loss[0]:
		lossRTP = s.getWinOrLoss(lossMath.NetLoss4.RTPGroup, lossMath.NetLoss4.Weight)
		rtpId = strconv.Itoa(lossRTP)
		bullets = s.getLossRtpBullet(rtpId, &lossMath.NetLoss4)

	case lossMath.NetLoss5.Loss[0]:
		lossRTP = s.getWinOrLoss(lossMath.NetLoss5.RTPGroup, lossMath.NetLoss5.Weight)
		rtpId = strconv.Itoa(lossRTP)
		bullets = s.getLossRtpBullet(rtpId, &lossMath.NetLoss5)

	default:
		return "", 0
	}

	return rtpId, bullets
}

func (s *service) getWinRtpBullet(rtpId string, rtp *struct{ PSF_ON_00003_1.DrbNetWin }) uint64 {
	switch rtpId {
	case rtp.RTP1.ID:
		return uint64(rtp.RTP1.Bullet[0])
	case rtp.RTP2.ID:
		return uint64(rtp.RTP2.Bullet[0])
	case rtp.RTP3.ID:
		return uint64(rtp.RTP3.Bullet[0])
	case rtp.RTP4.ID:
		return uint64(rtp.RTP4.Bullet[0])
	case rtp.RTP5.ID:
		return uint64(rtp.RTP5.Bullet[0])
	default:
		return 0
	}
}

func (s *service) getLossRtpBullet(rtpId string, rtp *struct{ PSF_ON_00003_1.DrbNetLoss }) uint64 {
	switch rtpId {
	case rtp.RTP1.ID:
		return uint64(rtp.RTP1.Bullet[0])
	case rtp.RTP2.ID:
		return uint64(rtp.RTP2.Bullet[0])
	case rtp.RTP3.ID:
		return uint64(rtp.RTP3.Bullet[0])
	case rtp.RTP4.ID:
		return uint64(rtp.RTP4.Bullet[0])
	case rtp.RTP5.ID:
		return uint64(rtp.RTP5.Bullet[0])
	default:
		return 0
	}
}

func (s *service) getWinOrLoss(item, weight []int) int {
	options := make([]rng.Option, 0, len(weight))

	for i := 0; i < len(weight); i++ {
		options = append(options, rng.Option{
			Weight: weight[i],
			Item:   item[i],
		})
	}

	return s.rng(options).(int)
}

func (s *service) rng(options []rng.Option) interface{} {
	return rng.New(options).Item
}
