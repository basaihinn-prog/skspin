package PSFM_00003_93_1

import (
	"path/filepath"
	"runtime"
	PSFM_00003_1 "serve/service_fish/domain/probability/PSFM-00003-1"
	"serve/service_fish/models"
)

var _, drb, _, _ = runtime.Caller(0)

var RTP93DRB = &rtp93drb{
	MathInfo: PSFM_00003_1.MathInfo{
		Docker:   PSFM_00003_1.DOCKER,
		Test:     filepath.Dir(drb),
		Folder:   PSFM_00003_1.FOLDER + PSFM_00003_1.RTP93,
		File:     "PSFM-00003-93-1-DRB.json",
		FileType: models.DRB,
	},
}

type rtp93drb struct {
	PSFM_00003_1.MathInfo
}

func init() {
	RTP93DRB.Deserialization()
}
