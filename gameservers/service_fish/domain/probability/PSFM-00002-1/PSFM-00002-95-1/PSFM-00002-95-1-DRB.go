package PSFM_00002_95_1

import (
	"path/filepath"
	"runtime"
	PSFM_00002_1 "serve/service_fish/domain/probability/PSFM-00002-1"
	"serve/service_fish/models"
)

var _, drb, _, _ = runtime.Caller(0)

var RTP95DRB = &rtp95drb{
	MathInfo: PSFM_00002_1.MathInfo{
		Docker:   PSFM_00002_1.DOCKER,
		Test:     filepath.Dir(drb),
		Folder:   PSFM_00002_1.FOLDER + PSFM_00002_1.RTP95,
		File:     "PSFM-00002-95-1-DRB.json",
		FileType: models.DRB,
	},
}

type rtp95drb struct {
	PSFM_00002_1.MathInfo
}

func init() {
	RTP95DRB.Deserialization()
}
