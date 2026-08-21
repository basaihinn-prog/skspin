package PSFM_00002_98_1

import (
	"path/filepath"
	"runtime"
	PSFM_00002_1 "serve/service_fish/domain/probability/PSFM-00002-1"
	"serve/service_fish/models"
)

var _, fs, _, _ = runtime.Caller(0)

var RTP98FS = &rtp98fs{
	MathInfo: PSFM_00002_1.MathInfo{
		Docker:   PSFM_00002_1.DOCKER,
		Test:     filepath.Dir(fs),
		Folder:   PSFM_00002_1.FOLDER + PSFM_00002_1.RTP98,
		File:     "PSFM-00002-98-1-FS.json",
		FileType: models.FS,
	},
}

type rtp98fs struct {
	PSFM_00002_1.MathInfo
}

func init() {
	RTP98FS.Deserialization()
}
