package PSFM_00013_96_1

import (
	"path/filepath"
	"runtime"
	PSFM_00013_1 "serve/service_fish/domain/probability/PSFM-00013-1"
	"serve/service_fish/models"
)

var _, fs, _, _ = runtime.Caller(0)

var RTP96FS = &rtp96fs{
	MathInfo: PSFM_00013_1.MathInfo{
		Docker:   PSFM_00013_1.DOCKER,
		Test:     filepath.Dir(fs),
		Folder:   PSFM_00013_1.FOLDER + PSFM_00013_1.RTP96,
		File:     "PSFM-00013-96-1-FS.json",
		FileType: models.FS,
	},
}

type rtp96fs struct {
	PSFM_00013_1.MathInfo
}

func init() {
	RTP96FS.Deserialization()
}
