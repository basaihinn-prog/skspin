package PSFM_00006_95_1

import (
	"path/filepath"
	"runtime"
	PSFM_00006_1 "serve/service_fish/domain/probability/PSFM-00006-1"
	"serve/service_fish/models"
)

var _, drb, _, _ = runtime.Caller(0)

var RTP95DRB = &rtp95drb{
	MathInfo: PSFM_00006_1.MathInfo{
		Docker:   PSFM_00006_1.DOCKER,
		Test:     filepath.Dir(drb),
		Folder:   PSFM_00006_1.FOLDER + PSFM_00006_1.RTP95,
		File:     "PSFM-00006-95-1-DRB.json",
		FileType: models.DRB,
	},
}

type rtp95drb struct {
	PSFM_00006_1.MathInfo
}

func init() {
	RTP95DRB.Deserialization()
}
