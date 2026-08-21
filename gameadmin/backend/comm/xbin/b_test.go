package xbin

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Resp struct {
	X int8
}

type Resp1 struct {
	X int8
}

type HitLine struct {
	// Line   *Line
	G       int
	Counts  []int
	Rate    int
	Multi   int
	Color   int
	Formula string `xbin:"-"`
}

type PlayResp struct {
	// Grade     int
	// Di        int
	// TotalG    uint
	// Times     float64
	// UniqueStr string
	// Bytes     []int8
	// Pan       [][]int8

	HitLine *HitLine
	R       Resp
	Arr     []*Resp
}

func TestXxx(t *testing.T) {

	v := &PlayResp{
		// Grade:     -123,
		// UniqueStr: "aklsdj",
		// Di:        123123,
		// Times:     0.002,
		// TotalG:    8888,
		// Bytes:     []int8{1, 2, 3, 4, 5},
		// Pan:       [][]int8{{1, 3, 5, 6, 8}, {2, 2, 2, 2}},

		HitLine: &HitLine{
			G:       123,
			Counts:  []int{1, 2, 3, 4},
			Rate:    4,
			Multi:   10,
			Color:   9,
			Formula: "hello formula",
		},
		R: Resp{X: 123},
		Arr: []*Resp{
			{X: 25},
			{X: 33},
		},
	}
	{
		s, _ := json.Marshal(v)
		// fmt.Println(len(bytesW), len(s), string(s))

		var bb *PlayResp
		json.Unmarshal(s, &bb)

		kk, _ := json.Marshal(bb)
		fmt.Println("kk:", string(kk))
	}

	bytesW, _ := Marshal(v)

	var o1 = PlayResp{}
	test1(bytesW, &o1)
	fmt.Println("1", tojson(o1))

	var o2 *PlayResp
	Unmarshal(bytesW, &o2)
	fmt.Println("2", tojson(o2))

	o3 := PlayResp{}
	Unmarshal(bytesW, &o3)
	fmt.Println("3", tojson(o3))
}

func test1(b []byte, to any) {
	test2(b, &to)
}

func test2(b []byte, to any) {

	Unmarshal(b, &to)
}

func tojson(a any) string {
	s, _ := json.Marshal(a)
	return string(s)
}
