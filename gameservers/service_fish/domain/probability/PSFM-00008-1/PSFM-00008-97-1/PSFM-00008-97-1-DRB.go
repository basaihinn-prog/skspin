package PSFM_00008_97_1

import (
	"path/filepath"
	"runtime"
	PSFM_00008_1 "serve/service_fish/domain/probability/PSFM-00008-1"
	"serve/service_fish/models"
)

var _, drb, _, _ = runtime.Caller(0)

var RTP97DRB = &rtp97drb{
	MathInfo: PSFM_00008_1.MathInfo{
		Docker:   PSFM_00008_1.DOCKER,
		Test:     filepath.Dir(drb),
		Folder:   PSFM_00008_1.FOLDER + PSFM_00008_1.RTP97,
		File:     "PSFM-00008-97-1-DRB.json",
		FileType: models.DRB,
	},
}

type rtp97drb struct {
	PSFM_00008_1.MathInfo
}

func init() {
	RTP97DRB.Deserialization()
}
