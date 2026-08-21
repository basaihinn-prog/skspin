package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Args struct {
	Level  Level
	Msg    string
	Time   time.Time
	Buffer buffer
}

type IPart interface {
	Output(args *Args)
}

type Logger struct {
	parts  []IPart
	output io.Writer
	pool   sync.Pool
	Level  Level
}

type OptionFunc func(l *Logger)

func OptOutput(o io.Writer) OptionFunc {
	return func(l *Logger) {
		l.output = o
	}
}

func OptLevel(lv Level) OptionFunc {
	return func(l *Logger) {
		l.Level = lv
	}
}

func NewLogger(options ...OptionFunc) *Logger {
	l := &Logger{}
	l.parts = []IPart{}
	l.output = os.Stdout

	for _, o := range options {
		o(l)
	}

	l.pool.New = func() any {
		return &Args{}
	}

	return l
}

func (l *Logger) Write(p []byte) (n int, err error) {
	return l.output.Write(p)
}

func (l *Logger) Use(parts ...IPart) *Logger {
	l.parts = parts
	return l
}

func (l *Logger) Log(level Level, msg string) {

	if l.Level > level {
		return
	}

	args := l.pool.Get().(*Args)
	args.Level = level
	args.Msg = msg
	args.Time = time.Now()
	args.Buffer.Reset()

	for i, part := range l.parts {
		part.Output(args)
		if i != len(l.parts)-1 {
			args.Buffer.Write(' ')
		}
	}

	bf := args.Buffer.buf
	if len(bf) == 0 || bf[len(bf)-1] != '\n' {
		args.Buffer.Write('\n')
	}

	l.output.Write(args.Buffer.buf)

	l.pool.Put(args)
}

func (l *Logger) Debug(i ...any)                 { l.Log(LevelDebug, fmt.Sprintln(i...)) }
func (l *Logger) Debugf(format string, i ...any) { l.Log(LevelDebug, fmt.Sprintf(format, i...)) }
func (l *Logger) Info(i ...any)                  { l.Log(LevelInfo, fmt.Sprintln(i...)) }
func (l *Logger) Infof(format string, i ...any)  { l.Log(LevelInfo, fmt.Sprintf(format, i...)) }
func (l *Logger) Warn(i ...any)                  { l.Log(LevelWarn, fmt.Sprintln(i...)) }
func (l *Logger) Warnf(format string, i ...any)  { l.Log(LevelWarn, fmt.Sprintf(format, i...)) }
func (l *Logger) Err(i ...any)                   { l.Log(LevelError, fmt.Sprintln(i...)) }
func (l *Logger) Errf(format string, i ...any)   { l.Log(LevelError, fmt.Sprintf(format, i...)) }

func (l *Logger) Fatal(i ...any) {
	l.Log(LevelFatal, fmt.Sprintln(i...))
	os.Exit(1)
}
func (l *Logger) Fatalf(format string, i ...any) {
	l.Log(LevelFatal, fmt.Sprintf(format, i...))
	os.Exit(1)
}
func (l *Logger) Panic(i ...any) {
	s := fmt.Sprintln(i...)
	l.Log(LevelPanic, s)
	panic(s)
}
func (l *Logger) Panicf(format string, i ...any) {
	s := fmt.Sprintf(format, i...)
	l.Log(LevelPanic, s)
	panic(s)
}

type ILogger interface {
	Debug(i ...any)
	Debugf(format string, i ...any)
	Info(i ...any)
	Infof(format string, i ...any)
	Warn(i ...any)
	Warnf(format string, i ...any)
	Err(i ...any)
	Errf(format string, i ...any)
	Fatal(i ...any)
	Fatalf(format string, i ...any)
	Panic(i ...any)
	Panicf(format string, i ...any)
}
