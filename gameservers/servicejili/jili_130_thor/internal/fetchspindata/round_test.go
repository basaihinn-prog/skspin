package main

import (
	"fmt"
	"testing"
)

func TestRoundId(t *testing.T) {
	// arr := []int32{294074, -1371455870, 294074680, -1361302324}

	// 11873938AEDC2CCC
	// 47CBA

	// 118750A666B00C7C

	// 66B00C7C
	// 118750A6

	id := (294080678 << 32) | 1722813564

	id = (294082667 << 32) | (-1261995476)&0xffffffff

	fmt.Print(id)
	// 1263066896118320252
}
