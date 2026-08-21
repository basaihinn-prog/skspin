package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
)

const (
	logdir = "log"
	Debug  = false
	EnvKey = "launch103167f931427968546c9386d875e12d"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("eg. ", os.Args[0], "ping", "www.baidu.com")
		return
	}

	if os.Getenv(EnvKey) == "" {
		// fmt.Println("我是父进程!", os.Getpid())
		cmd := &exec.Cmd{
			Path: lo.Must(exec.LookPath(os.Args[0])),
			Args: os.Args,      //注意,此处是包含程序名的
			Env:  os.Environ(), //父进程中的所有环境变量
			// Stdout: os.Stdout,
			// Stderr: os.Stderr,
		}
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", EnvKey, "1"))
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}

		// fmt.Println("拉起子进程!", os.Getpid())
		err := cmd.Start()
		if err != nil {
			slog.Error("fork fail!", "path", cmd.Path, "args", os.Args)
		}
		// fmt.Println("父进程 Exit!", os.Getpid())
		os.Exit(0)
		return
	}
	// fmt.Println("我是子进程!", os.Getpid())
	launch()
	// signalProc()
}

func signalProc() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Block until a signal is received.
	for s := range c {
		slog.Error("Got signal", "signal", s)
		// if s == syscall.SIGINT || s == syscall.SIGHUP {
		// 	continue
		// }
		os.Exit(0)
	}
}

func launch() {
	// [./launch echo 1 2 3]
	// slog.Default().()
	exepath := os.Args[1]

	exepath, err := exec.LookPath(exepath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))

	if Debug {
		c.AddFunc("* * * * *", func() { onDaily(path.Base(exepath)) })
		// onDaily(path.Base(exepath))
	} else {
		c.AddFunc("@daily", func() { onDaily(path.Base(exepath)) })
	}
	c.Start()

	var log = newLog(exepath)
	cmd := exec.Command(exepath, os.Args[2:]...)
	cmd.Stdout = log
	cmd.Stderr = log
	err = cmd.Run()
	slog.Error("process exit!", "exepath", exepath, "error", err, "args", os.Args)
}

func onDaily(exe string) {
	now := time.Now()

	var early time.Time
	if Debug {
		early = now.Add(-3 * time.Minute)
	} else {
		early = now.Add(-7 * 24 * time.Hour)
	}

	name := joinPath(logdir, exe, early)
	dir := path.Dir(name)
	earlyname := path.Base(name)

	// filepath.WalkDir()
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if !file.Type().IsRegular() {
			return
		}
		if file.Name() < earlyname {
			delname := path.Join(dir, file.Name())
			err := os.Remove(delname)
			slog.Info("remove log", "name", delname, "error", err)
		}
	}
}

type Log struct {
	exe  string
	file *os.File
}

func newLog(exepath string) *Log {
	exe := path.Base(exepath)
	log := &Log{
		exe: exe,
	}

	log.changefile()

	return log
}

func joinPath(logdir string, exe string, t time.Time) string {
	datestr := t.Format("20060102")
	if Debug {
		datestr = t.Format("200601021504")
	}
	return path.Join(logdir, exe, datestr)
}

func (log *Log) changefile() {
	name := joinPath(logdir, log.exe, time.Now())

	oldfile := log.file
	if oldfile != nil {
		if oldfile.Name() == name {
			// slog.Info("same logfile", "name", name)
			return
		}
		slog.Info("diff logfile", "name", name, "oldname", oldfile.Name())
		oldfile.Close()
		log.file = nil
	} else {
		slog.Info("init logfile", "name", name)
	}

	lo.Must0(os.MkdirAll(path.Dir(name), 0777))

	slog.Info("create logfile", "name", name)
	newfile := lo.Must(os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666))
	log.file = newfile

	linkname := log.exe + ".log"
	os.Remove(linkname)
	os.Symlink(name, linkname)
}

func (log *Log) Write(p []byte) (n int, err error) {
	log.changefile()
	file := log.file
	return file.Write(p)
}
