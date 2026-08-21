package rpc

import (
	"fmt"
	"strings"
	"testing"

	"serve/comm/ut"
)

// var winReg = regexp.MustCompile(`^(\w+ game win: )([0-9.+-]+)$`)

func TestXxx(t *testing.T) {
	desc := "Free game win: 1.36"
	desc = "Main game win: 1.6"
	rets := winReg.FindStringSubmatch(desc)
	fmt.Println(rets)
	if len(rets) == 3 {
		winstr := rets[2]
		ut.FloatStrPtrMul(&winstr, 2)

		ret := rets[1] + winstr
		fmt.Println(ret)
	}
}

func TestAA(t *testing.T) {
	s := "/history/csh/get-log-plate-info/1718692563274366002/xxxxxx"

	ret := strings.Split(s, "/")
	// strings.Sp
	fmt.Println(ret)
}
