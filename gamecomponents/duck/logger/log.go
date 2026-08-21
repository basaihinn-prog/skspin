package logger

import "time"

var DefaultLogger = New()

// 写到控制台
func New() *Logger {
	l := NewLogger()
	l.Use(
		PartLevel(),
		PartDateTime(time.RFC3339),
		PartCaller(4, true),
		PartMessage(),
	)
	return l
}

// 按日期写入文件
func SetDaily(folder, prefix string) {
	DefaultLogger = NewLogger(
		OptOutput(
			NewDailyWriter(folder + "/" + prefix + "-2006-01-02.log"),
		),
	)

	DefaultLogger.Use(
		PartLevel(),
		PartDateTime(time.RFC3339),
		PartCaller(4, true),
		PartMessage(),
	)
}

func Debug(i ...any)                 { DefaultLogger.Debug(i...) }
func Debugf(format string, i ...any) { DefaultLogger.Debugf(format, i...) }
func Info(i ...any)                  { DefaultLogger.Info(i...) }
func Infof(format string, i ...any)  { DefaultLogger.Infof(format, i...) }
func Warn(i ...any)                  { DefaultLogger.Warn(i...) }
func Warnf(format string, i ...any)  { DefaultLogger.Warnf(format, i...) }
func Err(i ...any)                   { DefaultLogger.Err(i...) }
func Errf(format string, i ...any)   { DefaultLogger.Errf(format, i...) }
func Fatal(i ...any)                 { DefaultLogger.Fatal(i...) }
func Fatalf(format string, i ...any) { DefaultLogger.Fatalf(format, i...) }
func Panic(i ...any)                 { DefaultLogger.Panic(i...) }
func Panicf(format string, i ...any) { DefaultLogger.Panicf(format, i...) }
