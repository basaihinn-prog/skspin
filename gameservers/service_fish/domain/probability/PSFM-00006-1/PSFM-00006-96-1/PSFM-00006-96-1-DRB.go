package PSFM_00006_96_1

import (
	"path/filepath"
	"runtime"
	PSFM_00006_1 "serve/service_fish/domain/probability/PSFM-00006-1"
	"serve/service_fish/models"
)

var _, drb, _, _ = runtime.Caller(0)

var RTP96DRB = &rtp96drb{
	MathInfo: PSFM_00006_1.MathInfo{
		Docker:   PSFM_00006_1.DOCKER,
		Test:     filepath.Dir(drb),
		Folder:   PSFM_00006_1.FOLDER + PSFM_00006_1.RTP96,
		File:     "PSFM-00006-96-1-DRB.json",
		FileType: models.DRB,
	},
}

type rtp96drb struct {
	PSFM_00006_1.MathInfo
}

func init() {
	RTP96DRB.Deserialization()
}
