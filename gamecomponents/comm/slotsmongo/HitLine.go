package slotsmongo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type HitLine struct {
	// Line   *Line
	G       int
	Counts  []int
	Rate    int
	Multi   int
	Color   int
	Formula string `xbin:"-"`
	// PanG    int    `xbin:"-"`
}

func (hl *HitLine) String() string {
	return fmt.Sprintf("%+v", *hl)
}

// 生成方程式, 后台复现使用
func (hl *HitLine) GenFormula(di, mul int) {
	// if hl.Formula != "" {
	// 	return
	// }
	// di := float64(di_) / 10000
	hl.G = hl.G * mul
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(mul))
	sb.WriteString("*(")

	counts := lo.Map(hl.Counts, func(item int, _ int) string {
		return strconv.Itoa(item)
	})

	sb.WriteString(strings.Join(counts, "*"))

	sb.WriteString(")*")
	sb.WriteString(strconv.Itoa(hl.Rate))
	sb.WriteString("*")

	fdi := float64(di) / 10000
	sb.WriteString(fmt.Sprintf("%.2f", fdi))

	sb.WriteString("=")

	// ans := float64(hl.G) / 10000
	ans := float64(hl.G) / 10000
	sb.WriteString(fmt.Sprintf("%.2f", ans))
	hl.Formula = sb.String()
	// (1*2*1)*6*0.05=0.60
}
