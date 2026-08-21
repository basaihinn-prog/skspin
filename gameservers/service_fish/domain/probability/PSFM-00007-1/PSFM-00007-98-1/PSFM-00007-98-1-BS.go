package PSFM_00007_98_1

import (
	"path/filepath"
	"runtime"
	PSFM_00007_1 "serve/service_fish/domain/probability/PSFM-00007-1"
	"serve/service_fish/models"
)

var _, bs, _, _ = runtime.Caller(0)

var RTP98BS = &rtp98bs{
	MathInfo: PSFM_00007_1.MathInfo{
		Docker:   PSFM_00007_1.DOCKER,
		Test:     filepath.Dir(bs),
		Folder:   PSFM_00007_1.FOLDER + PSFM_00007_1.RTP98,
		File:     "PSFM-00007-98-1-BS.json",
		FileType: models.BS,
	},
}

type rtp98bs struct {
	PSFM_00007_1.MathInfo
}

func init() {
	RTP98BS.Deserialization()
}
