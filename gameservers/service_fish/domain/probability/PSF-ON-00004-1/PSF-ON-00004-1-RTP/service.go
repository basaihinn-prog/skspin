package PSF_ON_00004_1_RTP

import (
	"serve/fish_comm/rng"
	PSF_ON_00004_1 "serve/service_fish/domain/probability/PSF-ON-00004-1"
)

var Service = &service{}

type service struct {
}

func (s *service) rngDenominator(rtpDrb *PSF_ON_00004_1.DrbMath) int {
	denominatorOption := make([]rng.Option, 0, 10)

	denominatorOption = s.appendDenominator(denominatorOption, &rtpDrb.BetGroup.Bet1)
	denominatorOption = s.appendDenominator(denominatorOption, &rtpDrb.BetGroup.Bet2)
	denominatorOption = s.appendDenominator(denominatorOption, &rtpDrb.BetGroup.Bet3)
	denominatorOption = s.appendDenominator(denominatorOption, &rtpDrb.BetGroup.Bet4)
	denominatorOption = s.appendDenominator(denominatorOption, &rtpDrb.BetGroup.Bet5)
	denominatorOption = s.appendDenominator(denominatorOption, &rtpDrb.BetGroup.Bet6)
	denominatorOption = s.appendDenominator(denominatorOption, &rtpDrb.BetGroup.Bet7)
	denominatorOption = s.appendDenominator(denominatorOption, &rtpDrb.BetGroup.Bet8)
	denominatorOption = s.appendDenominator(denominatorOption, &rtpDrb.BetGroup.Bet9)
	denominatorOption = s.appendDenominator(denominatorOption, &rtpDrb.BetGroup.Bet10)

	return s.rng(denominatorOption).(int)
}

func (s *service) rngMolecular(rtpDrb *PSF_ON_00004_1.DrbMath) (win int) {
	molecularOption := make([]rng.Option, 0, 10)

	molecularOption = s.appendMolecular(molecularOption, &rtpDrb.WinGroup.Win1)
	molecularOption = s.appendMolecular(molecularOption, &rtpDrb.WinGroup.Win2)
	molecularOption = s.appendMolecular(molecularOption, &rtpDrb.WinGroup.Win3)
	molecularOption = s.appendMolecular(molecularOption, &rtpDrb.WinGroup.Win4)
	molecularOption = s.appendMolecular(molecularOption, &rtpDrb.WinGroup.Win5)
	molecularOption = s.appendMolecular(molecularOption, &rtpDrb.WinGroup.Win6)
	molecularOption = s.appendMolecular(molecularOption, &rtpDrb.WinGroup.Win7)
	molecularOption = s.appendMolecular(molecularOption, &rtpDrb.WinGroup.Win8)
	molecularOption = s.appendMolecular(molecularOption, &rtpDrb.WinGroup.Win9)
	molecularOption = s.appendMolecular(molecularOption, &rtpDrb.WinGroup.Win10)

	return s.rng(molecularOption).(int)
}

func (s *service) appendDenominator(option []rng.Option, bet *struct{ PSF_ON_00004_1.DrbBetGroupMath }) []rng.Option {
	option = append(option, rng.Option{
		Weight: bet.Weight,
		Item:   bet.Turnover,
	})

	return option
}

func (s *service) appendMolecular(option []rng.Option, win *struct{ PSF_ON_00004_1.DrbWinGroupMath }) []rng.Option {
	option = append(option, rng.Option{
		Weight: win.Weight,
		Item:   win.Win,
	})

	return option
}

func (s *service) getMultiplier(win int, rtpDrb *PSF_ON_00004_1.DrbMath) (multiplier int) {
	switch win {
	case rtpDrb.WinGroup.Win1.Win:
		return rtpDrb.WinGroup.Win1.Multiplier
	case rtpDrb.WinGroup.Win2.Win:
		return rtpDrb.WinGroup.Win2.Multiplier
	case rtpDrb.WinGroup.Win3.Win:
		return rtpDrb.WinGroup.Win3.Multiplier
	case rtpDrb.WinGroup.Win4.Win:
		return rtpDrb.WinGroup.Win4.Multiplier
	case rtpDrb.WinGroup.Win5.Win:
		return rtpDrb.WinGroup.Win5.Multiplier
	case rtpDrb.WinGroup.Win6.Win:
		return rtpDrb.WinGroup.Win6.Multiplier
	case rtpDrb.WinGroup.Win7.Win:
		return rtpDrb.WinGroup.Win7.Multiplier
	case rtpDrb.WinGroup.Win8.Win:
		return rtpDrb.WinGroup.Win8.Multiplier
	case rtpDrb.WinGroup.Win9.Win:
		return rtpDrb.WinGroup.Win9.Multiplier
	case rtpDrb.WinGroup.Win10.Win:
		return rtpDrb.WinGroup.Win10.Multiplier
	default:
		return 0
	}
}

func (s service) rngRtpId() (rtpIdOrder int) {
	options := make([]rng.Option, 0, 5)
	options = append(options, rng.Option{Weight: 1, Item: 1})
	options = append(options, rng.Option{Weight: 1, Item: 2})
	options = append(options, rng.Option{Weight: 1, Item: 3})
	options = append(options, rng.Option{Weight: 1, Item: 4})
	options = append(options, rng.Option{Weight: 1, Item: 5})

	return s.rng(options).(int)
}

func (s service) rng(options []rng.Option) interface{} {
	return rng.New(options).Item
}

func (s *service) firstMultiplier(subgameId int, rtpDrb *PSF_ON_00004_1.DrbMath) int {
	switch subgameId {
	case 0:
		return rtpDrb.FirstMultiplier1[0]
	case 1:
		return rtpDrb.FirstMultiplier2[0]
	case 2:
		return rtpDrb.FirstMultiplier3[0]
	default:
		return 0
	}
}
