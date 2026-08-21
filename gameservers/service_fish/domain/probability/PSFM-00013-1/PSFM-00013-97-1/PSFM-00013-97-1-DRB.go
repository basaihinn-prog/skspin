package PSFM_00013_97_1

import (
	"path/filepath"
	"runtime"
	PSFM_00013_1 "serve/service_fish/domain/probability/PSFM-00013-1"
	"serve/service_fish/models"
)

var _, drb, _, _ = runtime.Caller(0)

var RTP97DRB = &rtp97drb{
	MathInfo: PSFM_00013_1.MathInfo{
		Docker:   PSFM_00013_1.DOCKER,
		Test:     filepath.Dir(drb),
		Folder:   PSFM_00013_1.FOLDER + PSFM_00013_1.RTP97,
		File:     "PSFM-00013-97-1-DRB.json",
		FileType: models.DRB,
	},
}

type rtp97drb struct {
	PSFM_00013_1.MathInfo
}

func init() {
	RTP97DRB.Deserialization()
}
