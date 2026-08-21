package PSFM_00004_96_1

import (
	"path/filepath"
	"runtime"
	PSFM_00004_1 "serve/service_fish/domain/probability/PSFM-00004-1"
	"serve/service_fish/models"
)

var _, fs, _, _ = runtime.Caller(0)

var RTP96FS = &rtp96fs{
	MathInfo: PSFM_00004_1.MathInfo{
		Docker:   PSFM_00004_1.DOCKER,
		Test:     filepath.Dir(fs),
		Folder:   PSFM_00004_1.FOLDER + PSFM_00004_1.RTP96,
		File:     "PSFM-00004-96-1-FS.json",
		FileType: models.FS,
	},
}

type rtp96fs struct {
	PSFM_00004_1.MathInfo
}

func init() {
	RTP96FS.Deserialization()
}
