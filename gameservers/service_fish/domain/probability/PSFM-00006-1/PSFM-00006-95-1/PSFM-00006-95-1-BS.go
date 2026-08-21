package PSFM_00006_95_1

import (
	"path/filepath"
	"runtime"
	PSFM_00006_1 "serve/service_fish/domain/probability/PSFM-00006-1"
	"serve/service_fish/models"
)

var _, bs, _, _ = runtime.Caller(0)

var RTP95BS = &rtp95bs{
	MathInfo: PSFM_00006_1.MathInfo{
		Docker:   PSFM_00006_1.DOCKER,
		Test:     filepath.Dir(bs),
		Folder:   PSFM_00006_1.FOLDER + PSFM_00006_1.RTP95,
		File:     "PSFM-00006-95-1-BS.json",
		FileType: models.BS,
	},
}

type rtp95bs struct {
	PSFM_00006_1.MathInfo
}

func init() {
	RTP95BS.Deserialization()
}
