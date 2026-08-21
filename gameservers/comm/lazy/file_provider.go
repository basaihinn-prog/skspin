package lazy

import (
	"os"

	"github.com/samber/lo"

	"gopkg.in/yaml.v3"
)

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type FilePortProvider struct {
	config IMap[string, string]
}

func NewFilePortProvider() *FilePortProvider {
	fd := &FilePortProvider{}
	fd.config = NewSyncMap[string, string]()

	bytes := lo.Must(os.ReadFile("config/grpc_route.yaml"))
	lo.Must0(fd.loadConfig(bytes))
	return fd
}

func (am *FilePortProvider) Get(key string) (addr string, ok bool) {
	return am.config.Load(key)
}

func (f *FilePortProvider) loadConfig(buf []byte) error {
	tmp := map[string]string{}
	err := yaml.Unmarshal(buf, &tmp)
	if err != nil {
		return err
	}

	f.config.Clear()
	for k, v := range tmp {
		f.config.Store(k, v)
	}

	return nil
}
