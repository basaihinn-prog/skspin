package PSFM_00004_97_1

import (
	"path/filepath"
	"runtime"
	PSFM_00004_1 "serve/service_fish/domain/probability/PSFM-00004-1"
	"serve/service_fish/models"
)

var _, bs, _, _ = runtime.Caller(0)

var RTP97BS = &rtp97bs{
	MathInfo: PSFM_00004_1.MathInfo{
		Docker:   PSFM_00004_1.DOCKER,
		Test:     filepath.Dir(bs),
		Folder:   PSFM_00004_1.FOLDER + PSFM_00004_1.RTP97,
		File:     "PSFM-00004-97-1-BS.json",
		FileType: models.BS,
	},
}

type rtp97bs struct {
	PSFM_00004_1.MathInfo
}

func init() {
	RTP97BS.Deserialization()
}
