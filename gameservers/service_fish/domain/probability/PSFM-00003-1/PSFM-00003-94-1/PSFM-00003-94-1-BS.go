package PSFM_00003_94_1

import (
	"path/filepath"
	"runtime"
	PSFM_00003_1 "serve/service_fish/domain/probability/PSFM-00003-1"
	"serve/service_fish/models"
)

var _, bs, _, _ = runtime.Caller(0)

var RTP94BS = &rtp94bs{
	MathInfo: PSFM_00003_1.MathInfo{
		Docker:                PSFM_00003_1.DOCKER,
		Test:                  filepath.Dir(bs),
		Folder:                PSFM_00003_1.FOLDER + PSFM_00003_1.RTP94,
		File:                  "PSFM-00003-94-1-BS.json",
		PSF_ON_00002_1_BsMath: nil,
		FileType:              models.BS,
	},
}

type rtp94bs struct {
	PSFM_00003_1.MathInfo
}

func init() {
	RTP94BS.Deserialization()
}
