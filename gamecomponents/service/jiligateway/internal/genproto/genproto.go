package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

func main() {
	raw_content, _ := os.ReadFile("bundle.543e0.js")
	content := raw_content

	exp := regexp.MustCompile(`r = new \w.(\w+)Proto\.(\w+);`)

	var out strings.Builder
	out.WriteString(`syntax = "proto3";
package message;
option go_package = "./";

`)
	// out.WriteByte('\n')

	for i := 0; ; i++ {
		loc := exp.FindSubmatchIndex(content)
		if loc == nil {
			break
		}

		classname := string(content[loc[2]:loc[3]]) + "_" + string(content[loc[4]:loc[5]])

		// fmt.Printf("%d, type %s\n", i, classname)
		slog.Info("", "classname", classname)
		fmt.Fprintf(&out, "message %s {\n", classname)
		// out.WriteString()
		// out.WriteByte('\n')

		content = content[loc[1]:]

		endpos := bytes.Index(content, []byte("return r"))
		procOne(&out, content[:endpos])
		out.WriteString("}\n\n")

		content = content[endpos:]
	}

	os.WriteFile("message.proto", []byte(out.String()), 0644)
}

var (
	// case 1: r.AwardSymbol = e.int32(); break;
	// caseExp = regexp.MustCompile(`case (\d+): r\.(\w+) = e\.(\w+)\(\);\s+break;`)

	// case 3: r.ColumnSymbol && r.ColumnSymbol.length || (r.ColumnSymbol = []), r.ColumnSymbol.push(p.cshProto.Column.decode(e, e.uint32())); break;

	// if (r.Row && r.Row.length || (r.Row = []), 2 == (7 & o)) for (var a = e.uint32() + e.pos; e.pos < a;)r.Row.push(e.int32()); else r.Row.push(e.int32());
	// if (r.Plate && r.Plate.length || (r.Plate = []), 2 == (7 & o)) for (var a = e.uint32() + e.pos; e.pos < a;)r.Plate.push(e.int32()); else r.Plate.push(e.int32());

	// " if (r.bet && r.bet.length || (r.bet = []), 2 == (7 & i)) for (var o = e.uint32() + e.pos; e.pos < o;)r.bet.push(e.double()); else r.bet.push(e.double()); "
	// " if (r.rtags && r.rtags.length || (r.rtags = []), 2 == (7 & i)) for (o = e.uint32() + e.pos; e.pos < o;)r.rtags.push(e.int32()); else r.rtags.push(e.int32()); "

	// r.Prefer = p.cshProto.PreferRoundShow.decode(e, e.uint32());

	// r.Name && r.Name.length || (r.Name = []), r.Name.push(e.string());
	// " r.List && r.List.length || (r.List = []), r.List.push(m.buffProto.WBuffOut.decode(e, e.uint32())); "

	caseExp  = regexp.MustCompile(`case (\d+):`)
	caseExp1 = regexp.MustCompile(`r\.(\w+) = e\.(\w+)\(\)`)
	caseExp2 = regexp.MustCompile(`r\.(\w+) \&\& r\.\w+.length \|\| \(r\.\w+ = \[\]\), r\.\w+\.push\(\w\.(\w+)Proto\.(\w+)\.decode\(e, e\.uint32\(\)\)\);`)
	caseExp3 = regexp.MustCompile(`if \(r\.(\w+) \&\& r\.\w+\.length \|\| \(r\.\w+ = \[\]\), 2 == \(7 \& \w\)\) for \(v?a?r? ?\w = e\.uint32\(\) \+ e\.pos; e\.pos \< \w;\)r\.\w+\.push\(e\.(\w+)\(\)\); else r\.\w+\.push\(e\.\w+\(\)\);`)

	caseExp4 = regexp.MustCompile(`r\.(\w+) = \w\.(\w+)Proto\.(\w+)\.decode\(e\, e\.uint32\(\)\);`)

	caseExp5 = regexp.MustCompile(`r\.(\w+) \&\& r\.\w+\.length \|\| \(r\.\w+ = \[\]\), r\.\w+\.push\(e\.string\(\)\);`)
)

func procOne(w io.Writer, data []byte) {
	// os.Stdout.Write(data)
	for {
		loc := caseExp.FindSubmatchIndex(data)
		if loc == nil {
			break
		}

		// fmt.Println(rets)
		index := data[loc[2]:loc[3]]
		var field, typename []byte

		// data = data[loc[1]:]

		endpos := bytes.Index(data, []byte("break;"))

		// row := data[loc[4]:loc[5]]
		row := data[loc[1]:endpos]

		field, typename = parseCase(row)

		// field := data[loc[4]:loc[5]]
		// typename := data[loc[6]:loc[7]]

		slog.Info("procOne", "index", string(index), "field", string(field), "typename", string(typename))

		optionalPrefix := "optional "
		if strings.HasPrefix(string(typename), "repeated") {
			optionalPrefix = ""
		}
		fmt.Fprintf(w, "    %s%s %s = %s;\n", optionalPrefix, typename, field, index)

		data = data[endpos+len("break;"):]
	}
}

func parseCase(row []byte) (field, typename []byte) {
	if vv := caseExp1.FindSubmatch(row); vv != nil {
		field = vv[1]
		typename = vv[2]
		return
	}
	if vv := caseExp2.FindSubmatch(row); vv != nil {
		field = vv[1]
		// typename = append([]byte("repeated "),vv[2]..., '_', vv[2]...)
		typename = []byte(fmt.Sprintf("repeated %s_%s", vv[2], vv[3]))
		return
	}
	if vv := caseExp3.FindSubmatch(row); vv != nil {
		field = vv[1]
		typename = append([]byte("repeated "), vv[2]...)
		return
	}
	if vv := caseExp4.FindSubmatch(row); vv != nil {
		field = vv[1]
		// typename = vv[2]
		typename = []byte(fmt.Sprintf("repeated %s_%s", vv[2], vv[3]))
		return
	}

	if vv := caseExp5.FindSubmatch(row); vv != nil {
		field = vv[1]
		typename = []byte("string")
		return
	}

	// panic("sorry")
	log.Panicf("parseCase panic, row=%s", row)
	return
}
