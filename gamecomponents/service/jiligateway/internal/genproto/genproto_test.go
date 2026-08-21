package main

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	row := "if (r.Row && r.Row.length || (r.Row = []), 2 == (7 & o)) for (var a = e.uint32() + e.pos; e.pos < a;)r.Row.push(e.int32()); else r.Row.push(e.int32());"
	row2 := " if (r.rtags && r.rtags.length || (r.rtags = []), 2 == (7 & i)) for (o = e.uint32() + e.pos; e.pos < o;)r.rtags.push(e.int32()); else r.rtags.push(e.int32()); "

	// caseExp3 := regexp.MustCompile(`else r\.\w+\.push\(e\.\w+\(\)\);`)

	caseExp3 := regexp.MustCompile(`if \(r\.(\w+) \&\& r\.\w+\.length \|\| \(r\.\w+ = \[\]\), 2 == \(7 \& \w\)\) for \((var |)\w = e\.uint32\(\) \+ e\.pos; e\.pos \< \w;\)r\.\w+\.push\(e\.(\w+)\(\)\); else r\.\w+\.push\(e\.\w+\(\)\);`)
	assert.True(t, caseExp3.MatchString(row))
	assert.True(t, caseExp3.MatchString(row2))

}

func TestYYY(t *testing.T) {
	s := []byte("case 1: r.currencyNumber = e.int32(); break; case 2: r.currencyName = e.string(); break; case 3: r.currencySymbol = e.string(); break; case 4: r.coin = e.double(); break;")

	loc := caseExp.FindSubmatchIndex([]byte(s))
	os.Stdout.Write(s[loc[4]:loc[5]])
}
