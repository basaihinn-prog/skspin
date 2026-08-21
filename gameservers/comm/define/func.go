package define

import (
	"fmt"
	"strings"
)

// 千分位
func ThousandSeparator(s string) string {
	numStrs := strings.Split(s, ".")
	numStr := numStrs[0]
	length := len(numStr)
	if length < 4 {
		return s
	}
	cnt := (length - 1) / 3
	for i := 0; i < cnt; i++ {
		numStr = numStr[:length-(i+1)*3] + "," + numStr[length-(i+1)*3:]
	}
	if len(numStrs) == 1 {
		return numStr
	}
	return fmt.Sprintf("%s.%s", numStr, numStrs[1])
}
