package bonuslottery

import "serve/fish_comm/flux"

func rkf_h5_00001_RedEnvelope(action *flux.Action, depositMultiple uint64) {
	psf_on_00001_RedEnvelope(action, depositMultiple)
}

func rkf_h5_00001_Slot(action *flux.Action, depositMultiple uint64) {
	psf_on_00001_Slot(action, depositMultiple)
}
