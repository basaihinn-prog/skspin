package rtp

func psf_on_00004_AvgPayCheck(fishId int32, triggerIconId, inputIconPay, avgPay, bet int, haveAvgPay bool) (iconPay int) {
	switch {
	case fishId == 100 && triggerIconId == 100:
		iconPay = avgPay * bet
	case fishId == 101 && triggerIconId == 101:
		iconPay = avgPay * bet
	default:
		if haveAvgPay && inputIconPay > 0 {
			iconPay = avgPay * bet
		} else {
			iconPay = inputIconPay * bet
		}
	}

	return iconPay
}
