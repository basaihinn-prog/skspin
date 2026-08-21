package logger

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DailyWriter struct {
	namefmt string
	file    *os.File
	lock    sync.Mutex

	lastYearDay int
}

func NewDailyWriter(filenameFormat string) *DailyWriter {
	dir := filepath.Dir(filenameFormat)
	if dir != "" && dir != "." {
		os.MkdirAll(dir, os.ModePerm)
	}

	return &DailyWriter{
		namefmt: filenameFormat,
	}
}

func (dw *DailyWriter) Write(p []byte) (n int, err error) {
	err = dw.checkNewDay()
	if err != nil {
		return
	}

	return dw.file.Write(p)
}

func (dw *DailyWriter) checkNewDay() error {
	dw.lock.Lock()
	defer dw.lock.Unlock()

	now := time.Now()
	yd := now.YearDay()

	if yd == dw.lastYearDay {
		return nil
	}

	logfile, err := os.OpenFile(now.Format(dw.namefmt), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	dw.lastYearDay = yd

	if dw.file != nil {
		dw.file.Close()
	}

	dw.file = logfile

	return nil
}
