package PSFM_00002_97_1

import (
	"fmt"
	"testing"
)

/*
func TestRtp98bs_Deserialization(t *testing.T) {
	//fmt.Println(RTP98BS.BsMath.Icons)
	fmt.Println("\n", RTP98BS.PSF_ON_00003_BsMath.Icons.Slot.RTP9)
}

func TestRtp98drb_Deserialization(t *testing.T) {
	fmt.Println("\n", RTP98DRB.PSF_ON_00003_DrbMath.RTPGroup.RTP9)
}
*/

func TestRtp98fs_Deserialization(t *testing.T) {
	fmt.Println("\n", "length HighRTPGroupWeight: ", len(RTP97DRB.PSF_ON_00002_DrbMath.HighRTPGroupWeight),
		"length LowRTPGroupWeight: ", len(RTP97DRB.PSF_ON_00002_DrbMath.LowRTPGroupWeight))
}
