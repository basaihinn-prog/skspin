package PSFM_00003_97_1

import (
	"fmt"
	"testing"
)

func TestRtp97bs_Reload(t *testing.T) {
	RTP97BS.deserialization()
	fmt.Println(RTP97BS.PSF_ON_00001_2_BsMath)
}

func TestRtp97drb_Reload(t *testing.T) {
	RTP97DRB.deserialization()
	fmt.Println(RTP97DRB.PSF_ON_00001_2_DrbMath)
}

func TestRtp97fs_Reload(t *testing.T) {
	RTP97FS.deserialization()
	fmt.Println(RTP97FS.PSF_ON_00001_2_FsMath)
}
