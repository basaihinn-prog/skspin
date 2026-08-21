package PSFM_00002_95_1

import (
	"path/filepath"
	"runtime"
	PSFM_00002_1 "serve/service_fish/domain/probability/PSFM-00002-1"
	"serve/service_fish/models"
)

var _, fs, _, _ = runtime.Caller(0)

var RTP95FS = &rtp95fs{
	MathInfo: PSFM_00002_1.MathInfo{
		Docker:   PSFM_00002_1.DOCKER,
		Test:     filepath.Dir(fs),
		Folder:   PSFM_00002_1.FOLDER + PSFM_00002_1.RTP95,
		File:     "PSFM-00002-95-1-FS.json",
		FileType: models.FS,
	},
}

type rtp95fs struct {
	PSFM_00002_1.MathInfo
}

func init() {
	RTP95FS.Deserialization()
}
