package cfgmgr

import (
	"game/duck/ut2/fileutil"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
)

type FileWatcher struct {
	configPath string
	watcher    *fsnotify.Watcher
}

func NewFileWatcher(path string) *FileWatcher {
	c := &FileWatcher{configPath: path}
	return c
}

func (m *FileWatcher) GetFileNames() (ret []string, err error) {

	d, err := os.ReadDir(m.configPath)
	if err != nil {
		return
	}

	d = slices.DeleteFunc(d, func(item fs.DirEntry) bool {
		return item.IsDir()
	})

	ret = lo.Map(d, func(item fs.DirEntry, _ int) string {
		return item.Name()
	})

	// filepath.Walk(m.configPath, func(path string, info fs.FileInfo, err error) error {
	// 	if info.IsDir() && info.Name() == ".git" {
	// 		return
	// 	}
	// 	if info == nil || info.IsDir() {
	// 		return nil
	// 	}

	// 	ret = append(ret, info.Name())

	// 	return nil
	// })

	return
}

func (m *FileWatcher) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(m.configPath, name))
}

func (m *FileWatcher) WriteFile(name string, content []byte) error {
	return os.WriteFile(filepath.Join(m.configPath, name), content, 0644)
}

func (m *FileWatcher) IsExists(name string) bool {
	return fileutil.FileExists(filepath.Join(m.configPath, name))
}

func (m *FileWatcher) Init() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	m.watcher = watcher

	return watcher.Add(m.configPath)
}

func (m *FileWatcher) Stop() error {
	return m.watcher.Close()
}

func (m *FileWatcher) Start(closeCh chan struct{}, dispatch func(string, []byte)) {

	for {
		select {
		case event := <-m.watcher.Events:
			if (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Create == fsnotify.Create) {
				dispatch(filepath.Base(event.Name), nil)
			}
		case <-m.watcher.Errors:

		case <-closeCh:
			return
		}
	}
}
