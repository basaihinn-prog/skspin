package logger

import (
	"path/filepath"
	"runtime"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

type partLevel struct {
}

func PartLevel() *partLevel {
	return &partLevel{}
}

var lv2str = []string{
	LevelDebug: "DBG",
	LevelInfo:  "INF",
	LevelWarn:  "WAN",
	LevelError: "ERR",
	LevelFatal: "FTL",
	LevelPanic: "PNC",
}

func (p *partLevel) Output(args *Args) {
	args.Buffer.WriteString(lv2str[args.Level])
}

////
type partCaller struct {
	skip      int
	shortFile bool
}

func PartCaller(skip int, shortFile bool) *partCaller {
	return &partCaller{skip: skip, shortFile: shortFile}
}

func (p *partCaller) Output(args *Args) {

	_, file, line, ok := runtime.Caller(p.skip)
	if !ok {
		file = "???"
		line = 0
	}

	if p.shortFile {
		file = filepath.Base(file)
	}

	args.Buffer.WriteString(file)
	args.Buffer.Write(':')
	args.Buffer.WriteInt(line)
}

////

type partDateTime struct {
	layout string
}

func PartDateTime(layout string) *partDateTime {
	return &partDateTime{layout: layout}
}

func (p *partDateTime) Output(args *Args) {
	args.Buffer.buf = args.Time.AppendFormat(args.Buffer.buf, p.layout)
}

////

type partMessage struct {
}

func PartMessage() *partMessage {
	return &partMessage{}
}

func (p *partMessage) Output(args *Args) {
	args.Buffer.WriteString(args.Msg)
}

////
