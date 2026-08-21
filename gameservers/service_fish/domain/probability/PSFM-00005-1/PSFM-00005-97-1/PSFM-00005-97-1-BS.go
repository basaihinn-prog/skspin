package PSFM_00005_97_1

import (
	"path/filepath"
	"runtime"
	PSFM_00005_1 "serve/service_fish/domain/probability/PSFM-00005-1"
	"serve/service_fish/models"
)

var _, bs, _, _ = runtime.Caller(0)

var RTP97BS = &rtp97bs{
	MathInfo: PSFM_00005_1.MathInfo{
		Docker:   PSFM_00005_1.DOCKER,
		Test:     filepath.Dir(bs),
		Folder:   PSFM_00005_1.FOLDER + PSFM_00005_1.RTP97,
		File:     "PSFM-00005-97-1-BS.json",
		FileType: models.BS,
	},
}

type rtp97bs struct {
	PSFM_00005_1.MathInfo
}

func init() {
	RTP97BS.Deserialization()
}
