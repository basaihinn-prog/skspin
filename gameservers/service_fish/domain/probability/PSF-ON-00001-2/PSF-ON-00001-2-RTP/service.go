package PSF_ON_00001_2_RTP

import (
	"math/rand"
	"serve/fish_comm/rng"
	PSF_ON_00001_2 "serve/service_fish/domain/probability/PSF-ON-00001-2"
	"time"
)

const (
	high = 1
	low  = 0
)

var Service = &service{}

type service struct {
}

func (s *service) highLowState(state int, rtpDrb *PSF_ON_00001_2.DrbMath) (rtpState, rtpGroupId int) {

	switch state {
	case high:
		rtpGroupId = s.rngLowRtpGroupId(rtpDrb)
		rtpState = low

	case low:
		rtpGroupId = s.rngHighRtpGroupId(rtpDrb)
		rtpState = high

	default:
		if s.rngHighLowGroup(rtpDrb.InitLowGroupWeight) == high {
			rtpGroupId = s.rngLowRtpGroupId(rtpDrb)
			rtpState = low
		} else {
			rtpGroupId = s.rngHighRtpGroupId(rtpDrb)
			rtpState = high
		}
	}
	return rtpState, rtpGroupId
}

func (s *service) rngHighRtpGroupId(rtpDrb *PSF_ON_00001_2.DrbMath) int {
	highRtpGroupId := make([]rng.Option, 0, len(rtpDrb.HighRTPGroupWeight))

	for i := 0; i < len(rtpDrb.HighRTPGroupWeight); i++ {
		highRtpGroupId = append(highRtpGroupId, rng.Option{
			rtpDrb.HighRTPGroupWeight[i],
			rtpDrb.RTPGroupID[i],
		})
	}

	return s.rng(highRtpGroupId)
}

func (s *service) rngLowRtpGroupId(rtpDrb *PSF_ON_00001_2.DrbMath) int {
	lowRtpGroupId := make([]rng.Option, 0, len(rtpDrb.LowRTPGroupWeight))

	for i := 0; i < len(rtpDrb.LowRTPGroupWeight); i++ {
		lowRtpGroupId = append(lowRtpGroupId, rng.Option{
			rtpDrb.LowRTPGroupWeight[i],
			rtpDrb.RTPGroupID[i],
		})
	}

	return s.rng(lowRtpGroupId)
}

func (s *service) rngBullets(min, max int) uint64 {
	rand.Seed(time.Now().UnixNano())
	return uint64(rand.Intn(max-min+1) + min)
}

func (s *service) rngHighLowGroup(initGroup []int) int {
	highLowOptions := []rng.Option{
		{initGroup[0], high},
		{initGroup[1], low},
	}
	return s.rng(highLowOptions)
}

func (r service) rng(options []rng.Option) int {
	return rng.New(options).Item.(int)
}
