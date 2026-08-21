package PSFM_00006_97_1

import (
	"path/filepath"
	"runtime"
	PSFM_00006_1 "serve/service_fish/domain/probability/PSFM-00006-1"
	"serve/service_fish/models"
)

var _, fs, _, _ = runtime.Caller(0)

var RTP97FS = &rtp97fs{
	MathInfo: PSFM_00006_1.MathInfo{
		Docker:   PSFM_00006_1.DOCKER,
		Test:     filepath.Dir(fs),
		Folder:   PSFM_00006_1.FOLDER + PSFM_00006_1.RTP97,
		File:     "PSFM-00006-97-1-FS.json",
		FileType: models.FS,
	},
}

type rtp97fs struct {
	PSFM_00006_1.MathInfo
}

func init() {
	RTP97FS.Deserialization()
}
