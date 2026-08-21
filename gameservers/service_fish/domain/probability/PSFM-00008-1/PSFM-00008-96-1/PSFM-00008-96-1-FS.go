package PSFM_00008_96_1

import (
	"path/filepath"
	"runtime"
	PSFM_00008_1 "serve/service_fish/domain/probability/PSFM-00008-1"
	"serve/service_fish/models"
)

var _, fs, _, _ = runtime.Caller(0)

var RTP96FS = &rtp96fs{
	MathInfo: PSFM_00008_1.MathInfo{
		Docker:   PSFM_00008_1.DOCKER,
		Test:     filepath.Dir(fs),
		Folder:   PSFM_00008_1.FOLDER + PSFM_00008_1.RTP96,
		File:     "PSFM-00008-96-1-FS.json",
		FileType: models.FS,
	},
}

type rtp96fs struct {
	PSFM_00008_1.MathInfo
}

func init() {
	RTP96FS.Deserialization()
}
