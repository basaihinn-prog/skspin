package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"serve/servicejili/jiliut"

	"github.com/samber/lo"
)

func main() {
	var (
		bundle  string
		outfile string
	)
	flag.StringVar(&bundle, "bundle", "", "/data/game/bin/cache/jilid.rslotszs001.com/gem3/src/chunks/bundle.032c8.js  或 https://uat-wbgame.jlfafafa3.com/astarte/3.6/web-mobile/assets/other/index.946a0.js")
	flag.StringVar(&outfile, "outfile", "tmp.proto", "tmp.proto")
	flag.Parse()

	if bundle == "" || outfile == "" {
		return
	}

	/*
		fmt.Println("输入游戏简称：")
		var shortName string
		shortName = "kk2"
		_, err := fmt.Scan(&shortName)
		if err != nil {
			fmt.Println("输入简称错误：", err)
			return
		}
		fmt.Println("输入游戏Id：")
		var gameId string
		gameId = "jili_16_kk2"
		_, err = fmt.Scan(&gameId)
		if err != nil {
			fmt.Println("输入简称错误：", err)
			return
		}

		var bundle string
		dir := fmt.Sprintf("/data/game/bin/cache/jilid.rslotszs001.com/%s/src/chunks", shortName)
		files, err := os.ReadDir(dir)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "bundle") {
				bundle = fmt.Sprintf("%s/%s", dir, file.Name())
				break
			}
		}
		if bundle == "" {
			fmt.Println("路径错误")
			return
		}
	*/

	// raw_content, _ := os.ReadFile("bundle.543e0.js")

	var content []byte
	if strings.HasPrefix(bundle, "http") {
		content = lo.Must(jiliut.GetBody(bundle))
	} else {
		content = lo.Must(os.ReadFile(bundle))
	}

	// prolog := regexp.MustCompile(`,\w=new \w\.(\w+)Proto\.(\w+);`)
	// , r = new p.snProto.lifeResp; e.pos < o;
	prolog := regexp.MustCompile(`,\w=new \w\.(\w+)Proto\.([^;]+);\w\.pos<\w;`)

	var out strings.Builder
	out.WriteString(`syntax = "proto3";
package message;
option go_package = "./";
import "google/protobuf/timestamp.proto";

`)

	for i := 0; ; i++ {
		loc := prolog.FindSubmatchIndex(content)
		if loc == nil {
			break
		}

		classname := string(content[loc[2]:loc[3]]) + "_" + strings.ReplaceAll(string(content[loc[4]:loc[5]]), ".", "_")

		slog.Info("", "classname", classname)
		fmt.Fprintf(&out, "message %s {\n", classname)

		content = content[loc[1]:]

		endpos := bytes.Index(content, []byte("return "))

		if !(strings.HasPrefix(classname, "mission_") ||
			strings.HasPrefix(classname, "sn_") ||
			false) {
			procOne(&out, content[:endpos])
		}
		out.WriteString("}\n\n")

		content = content[endpos:]
	}
	// os.MkdirAll(fmt.Sprintf("../%s/internal/message/", gameId), 0777)
	// outfile := fmt.Sprintf("../%s/internal/message/message.proto", gameId)
	// var err error
	os.WriteFile(outfile, []byte(out.String()), 0644)
	// fmt.Println(err)

	// os.Chdir(fmt.Sprintf("../%s/internal/message", gameId))
	/*
		cmd := exec.Command("protoc", "-I=.", "--go_out=.", "message.proto")
		err = cmd.Run()
		if err != nil {
			fmt.Println("shell err", err)
		}
	*/
}

var (
	caseExp  = regexp.MustCompile(`case (\d+):`)
	caseExp1 = regexp.MustCompile(`\w\.(\w+)=\w\.(\w+)\(\)`)
	caseExp2 = regexp.MustCompile(`\w\.(\w+)\&\&\w\.\w+.length\|\|\(\w\.\w+=\[\]\),\w\.\w+\.push\(\w\.(\w+)Proto\.(\w+)\.decode\(\w,\w\.uint32\(\)\)\);`)

	// n.award.push(p.customProto.Plate.AwardReel.decode(e, e.uint32()));
	caseExp2_1 = regexp.MustCompile(`\w\.(\w+)\&\&\w\.\w+.length\|\|\(\w\.\w+=\[\]\),\w\.\w+\.push\(\w\.(\w+)Proto\.(\w+\.\w+)\.decode\(\w,\w\.uint32\(\)\)\);`)
	caseExp3   = regexp.MustCompile(`if\(\w\.(\w+)\&\&\w\.\w+\.length\|\|\(\w\.\w+=\[\]\),2==\(7\&\w\)\)for\(v?a?r? ?\w=\w\.uint32\(\)\+\w\.pos;\w\.pos\<\w;\)\w\.\w+\.push\(\w\.(\w+)\(\)\);else \w\.\w+\.push\(\w\.\w+\(\)\);`)

	caseExp4 = regexp.MustCompile(`\w\.(\w+)=\w\.(\w+)Proto\.(\w+)\.decode\(\w\,\w\.uint32\(\)\);`)

	caseExp5 = regexp.MustCompile(`\w\.(\w+)\&\&\w\.\w+\.length\|\|\(\w\.\w+=\[\]\),\w\.\w+\.push\(\w\.string\(\)\);`)
	caseExp6 = regexp.MustCompile(`\w\.(\w+)=\w\.google\.protobuf\.Timestamp\.decode\(\w,\w\.uint32\(\)\);`)
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
	if vv := caseExp2_1.FindSubmatch(row); vv != nil {
		field = vv[1]
		// typename = append([]byte("repeated "),vv[2]..., '_', vv[2]...)
		typename = []byte(fmt.Sprintf("repeated %s_%s", vv[2], bytes.ReplaceAll(vv[3], []byte{'.'}, []byte{'_'})))
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
	if vv := caseExp6.FindSubmatch(row); vv != nil {
		field = vv[1]
		typename = []byte("google.protobuf.Timestamp")
		return
	}

	// panic("sorry")
	log.Panicf("parseCase panic, row=%s", row)
	return
}
