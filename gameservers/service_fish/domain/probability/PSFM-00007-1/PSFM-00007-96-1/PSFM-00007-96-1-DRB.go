package PSFM_00007_96_1

import (
	"path/filepath"
	"runtime"
	PSFM_00007_1 "serve/service_fish/domain/probability/PSFM-00007-1"
	"serve/service_fish/models"
)

var _, drb, _, _ = runtime.Caller(0)

var RTP96DRB = &rtp96drb{
	MathInfo: PSFM_00007_1.MathInfo{
		Docker:   PSFM_00007_1.DOCKER,
		Test:     filepath.Dir(drb),
		Folder:   PSFM_00007_1.FOLDER + PSFM_00007_1.RTP96,
		File:     "PSFM-00007-96-1-DRB.json",
		FileType: models.DRB,
	},
}

type rtp96drb struct {
	PSFM_00007_1.MathInfo
}

func init() {
	RTP96DRB.Deserialization()
}
