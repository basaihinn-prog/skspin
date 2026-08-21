package PSFM_00013_96_1

import (
	"path/filepath"
	"runtime"
	PSFM_00013_1 "serve/service_fish/domain/probability/PSFM-00013-1"
	"serve/service_fish/models"
)

var _, bs, _, _ = runtime.Caller(0)

var RTP96BS = &rtp96bs{
	MathInfo: PSFM_00013_1.MathInfo{
		Docker:   PSFM_00013_1.DOCKER,
		Test:     filepath.Dir(bs),
		Folder:   PSFM_00013_1.FOLDER + PSFM_00013_1.RTP96,
		File:     "PSFM-00013-96-1-BS.json",
		FileType: models.BS,
	},
}

type rtp96bs struct {
	PSFM_00013_1.MathInfo
}

func init() {
	RTP96BS.Deserialization()
}
