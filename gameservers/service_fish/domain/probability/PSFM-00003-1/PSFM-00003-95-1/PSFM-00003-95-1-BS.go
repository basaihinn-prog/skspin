package PSFM_00003_95_1

import (
	"path/filepath"
	"runtime"
	PSFM_00003_1 "serve/service_fish/domain/probability/PSFM-00003-1"
	"serve/service_fish/models"
)

var _, bs, _, _ = runtime.Caller(0)

var RTP95BS = &rtp95bs{
	MathInfo: PSFM_00003_1.MathInfo{
		Docker:                PSFM_00003_1.DOCKER,
		Test:                  filepath.Dir(bs),
		Folder:                PSFM_00003_1.FOLDER + PSFM_00003_1.RTP95,
		File:                  "PSFM-00003-95-1-BS.json",
		PSF_ON_00002_1_BsMath: nil,
		FileType:              models.BS,
	},
}

type rtp95bs struct {
	PSFM_00003_1.MathInfo
}

func init() {
	RTP95BS.Deserialization()
}
