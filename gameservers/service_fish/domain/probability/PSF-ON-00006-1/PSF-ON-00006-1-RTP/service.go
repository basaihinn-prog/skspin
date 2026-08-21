package PSF_ON_00006_1_RTP

import (
	"math/rand"
	PSF_ON_00006_1 "serve/service_fish/domain/probability/PSF-ON-00006-1"
	"time"

	"serve/fish_comm/rng"
)

const (
	low  = 0
	high = 1
)

var Service = &service{}

type service struct {
}

func (s *service) highLowState(state int, rtpDrb *PSF_ON_00006_1.DrbMath) (rtpState, rtpGroupId int) {
	switch state {
	case high:
		rtpGroupId = s.rngRtpGroupId(rtpDrb.LowRTPGroupWeight, rtpDrb.RTPGroupID)
		rtpState = low

	case low:
		rtpGroupId = s.rngRtpGroupId(rtpDrb.HighRTPGroupWeight, rtpDrb.RTPGroupID)
		rtpState = high

	default:
		if s.rngHighLowGroup() == high {
			rtpGroupId = s.rngRtpGroupId(rtpDrb.LowRTPGroupWeight, rtpDrb.RTPGroupID)
			rtpState = low
		} else {
			rtpGroupId = s.rngRtpGroupId(rtpDrb.HighRTPGroupWeight, rtpDrb.RTPGroupID)
			rtpState = high
		}
	}

	return rtpState, rtpGroupId
}

func (s *service) rngRtpGroupId(weights, group []int) int {
	rtpGroupId := make([]rng.Option, 0, len(weights))

	for i := 0; i < len(weights); i++ {
		rtpGroupId = append(rtpGroupId, rng.Option{weights[i], group[i]})
	}

	return s.rng(rtpGroupId)
}

func (s *service) rngBullets(min, max int) uint64 {
	rand.Seed(time.Now().UnixNano())
	return uint64(rand.Intn(max-min+1) + min)
}

func (s *service) rngHighLowGroup() int {
	highLowOptions := []rng.Option{
		{1, high},
		{1, low},
	}
	return s.rng(highLowOptions)
}

func (r service) rng(options []rng.Option) int {
	return rng.New(options).Item.(int)
}
