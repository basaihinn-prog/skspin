package lazy

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"sync"
	"time"
)

type LoadFunc func(buf []byte) error
type LoadFuncWithName func(file string, buf []byte) error
type LoadData struct {
	callback LoadFuncWithName
	path     string
	isprefix bool
}

func (d *LoadData) IsMatch(file string) bool {
	if d.isprefix {
		return strings.HasPrefix(file, d.path)
	}
	return file == d.path
}

// type Watcher interface {
// 	Init() error
// 	Start(closeCh chan struct{}, dispatch func(string, []byte))
// 	Stop() error
// 	GetFileNames() ([]string, error)
// 	ReadFile(name string) ([]byte, error)
// 	WriteFile(name string, content []byte) error
// 	IsExists(name string) bool
// }

type ConfigManager struct {
	// Etcd     *etcd.Etcd
	md5Cache IMap[string, string]
	elements []*LoadData
	wg       sync.WaitGroup
	closeCh  chan struct{}
	Watcher  *FileWatcher
}

func NewConfigManager(w *FileWatcher) *ConfigManager {
	return &ConfigManager{
		md5Cache: NewSyncMap[string, string](),
		Watcher:  w,
		closeCh:  make(chan struct{}, 1),
	}
}

func (m *ConfigManager) Close() {
	close(m.closeCh)
	m.wg.Wait()
}

func (m *ConfigManager) Load(file string) {
	m.dispatch(file, nil)
}

func (m *ConfigManager) LoadAll() {

	mp := []*LoadData{}
	for _, v := range m.elements {
		if v.isprefix {
			mp = append(mp, v)
		} else {
			m.dispatch(v.path, nil)
		}
	}

	if len(mp) == 0 {
		return
	}

	names, err := m.Watcher.GetFileNames()
	if err != nil {
		slog.Info("读取文件列表出错", "error", err)
		return
	}

	for _, name := range names {
		for _, one := range mp {
			if one.IsMatch(name) {
				m.dispatch(name, nil)
			}
		}
	}
}

func (m *ConfigManager) Start() {
	err := m.Watcher.Init()
	if err != nil {
		// slog.Info("Watcher初始化错误", "error", err)
		log.Panicln("ConfigManager Start failed!", err.Error())
		return
	}

	m.Watcher.Start(m.closeCh, m.dispatch)
	err = m.Watcher.Stop()
	if err != nil {
		slog.Info("Watcher关闭错误", "error", err)
	}
}

func (m *ConfigManager) find(file string) *LoadData {
	for _, v := range m.elements {
		if v.IsMatch(file) {
			return v
		}
	}
	return nil
}

const errfile = "_config_last_err.json"

type errorInfo struct {
	File string
	Md5  string
	Err  string
	Time time.Time
}

// 如果配置读取错误，那么会将错误写到一个文件
// 提交者可以获取到这个错误
func (m *ConfigManager) dispatch(file string, value []byte) {
	ele := m.find(file)
	if ele == nil {
		return
	}

	err := m.dispatchErr(ele, file, value)

	if file == errfile {
		return
	}

	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	jsonbuf, _ := json.Marshal(errorInfo{
		File: file,
		Err:  errStr,
		Md5:  fmt.Sprintf("%x", md5.Sum(value)),
		Time: time.Now(),
	})
	m.Watcher.WriteFile(errfile, jsonbuf)
}

func (m *ConfigManager) GetErr(file, md5 string) string {
	buf, err := m.Watcher.ReadFile(errfile)
	if err != nil {
		slog.Info("读取错误", "error", err)
		return ""
	}
	resp := errorInfo{}
	json.Unmarshal(buf, &resp)

	if resp.File == file && time.Now().Sub(resp.Time) < 1*time.Second {
		return resp.Err
	}

	return ""
}

func (m *ConfigManager) dispatchErr(ele *LoadData, file string, value []byte) error {
	if value == nil {
		buf, err := m.Watcher.ReadFile(file)
		if err != nil {
			slog.Info("读取配置错误", "file", file, "error", err)
			return nil
		}
		if len(buf) == 0 {
			slog.Info("读取配置错误 长度==0", "file", file)
			return nil
		}
		value = buf
	}

	newMD5 := fmt.Sprintf("%x", md5.Sum(value))

	oldMD5, _ := m.md5Cache.Load(file)
	if oldMD5 == newMD5 {
		slog.Info("配置未改动，跳过加载", "file", file)
		return nil
	}

	err := ele.callback(file, value)
	if err != nil {
		slog.Error("加载配置错误", "file", file, "error", err)
		return err
	} else {
		m.md5Cache.Store(file, newMD5)
		slog.Info("加载配置成功", "file", file, "newMD5", newMD5)
		return nil
	}

}

func (m *ConfigManager) WatchAndLoad(file string, f LoadFunc) {
	m.Watch(file, f)
	m.Load(file)
}

func (m *ConfigManager) Watch(file string, f LoadFunc) {
	m.elements = append(m.elements, &LoadData{
		path:     file,
		isprefix: false,
		callback: func(file string, buf []byte) error {
			return f(buf)
		},
	})
}

func (m *ConfigManager) WatchPrefix(file string, f LoadFuncWithName) {
	m.elements = append(m.elements, &LoadData{
		callback: f,
		path:     file,
		isprefix: true,
	})
}
