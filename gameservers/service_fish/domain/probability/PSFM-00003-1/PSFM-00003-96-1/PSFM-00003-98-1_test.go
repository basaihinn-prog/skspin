package PSFM_00003_96_1

import (
	"fmt"
	"testing"
)

func TestRtp96bs_Reload(t *testing.T) {

	RTP96BS.Reload()

	fmt.Println(RTP96BS.PSF_ON_00002_1_BsMath)
}

func TestRtp96drb_Reload(t *testing.T) {
	RTP96DRB.Reload()

	fmt.Println(RTP96DRB.PSF_ON_00002_1_DrbMath)
}

func TestRtp96fs_Reload(t *testing.T) {
	RTP96FS.Reload()

	fmt.Println(RTP96FS.PSF_ON_00002_1_FsMath)
}
