package bonuslottery

import (
	"serve/fish_comm/flux"
)

func psf_on_00002_RedEnvelope(action *flux.Action, depositMultiple uint64) {
	psf_on_00001_RedEnvelope(action, depositMultiple)
}

func psf_on_00002_Slot(action *flux.Action, depositMultiple uint64) {
	psf_on_00001_Slot(action, depositMultiple)
}
