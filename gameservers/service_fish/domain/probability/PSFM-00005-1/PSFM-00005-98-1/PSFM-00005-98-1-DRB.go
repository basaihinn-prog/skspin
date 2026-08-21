package PSFM_00005_98_1

import (
	"path/filepath"
	"runtime"
	PSFM_00005_1 "serve/service_fish/domain/probability/PSFM-00005-1"
	"serve/service_fish/models"
)

var _, drb, _, _ = runtime.Caller(0)

var RTP98DRB = &rtp98drb{
	MathInfo: PSFM_00005_1.MathInfo{
		Docker:   PSFM_00005_1.DOCKER,
		Test:     filepath.Dir(drb),
		Folder:   PSFM_00005_1.FOLDER + PSFM_00005_1.RTP98,
		File:     "PSFM-00005-98-1-DRB.json",
		FileType: models.DRB,
	},
}

type rtp98drb struct {
	PSFM_00005_1.MathInfo
}

func init() {
	RTP98DRB.Deserialization()
}
