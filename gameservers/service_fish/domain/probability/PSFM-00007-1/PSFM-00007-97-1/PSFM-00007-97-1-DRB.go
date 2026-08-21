package PSFM_00007_97_1

import (
	"path/filepath"
	"runtime"
	PSFM_00007_1 "serve/service_fish/domain/probability/PSFM-00007-1"
	"serve/service_fish/models"
)

var _, drb, _, _ = runtime.Caller(0)

var RTP97DRB = &rtp97drb{
	MathInfo: PSFM_00007_1.MathInfo{
		Docker:   PSFM_00007_1.DOCKER,
		Test:     filepath.Dir(drb),
		Folder:   PSFM_00007_1.FOLDER + PSFM_00007_1.RTP97,
		File:     "PSFM-00007-97-1-DRB.json",
		FileType: models.DRB,
	},
}

type rtp97drb struct {
	PSFM_00007_1.MathInfo
}

func init() {
	RTP97DRB.Deserialization()
}
