package discovery

import (
	"errors"
	"game/duck/logger"
	"game/duck/ut2"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type FilePortProvider struct {
	config ut2.IMap[string, string]
}

func NewFilePortProvider() *FilePortProvider {
	fd := &FilePortProvider{}
	fd.config = ut2.NewSyncMap[string, string]()

	bytes, err := os.ReadFile("config/grpc_route.yaml")
	if err != nil {
		logger.Err(err)
	}
	err = fd.loadConfig(bytes)
	if err != nil {
		logger.Err(err)
	}
	return fd
}

func (am *FilePortProvider) Close() {
}

func (am *FilePortProvider) GetPort(service string) (int, error) {
	arr, _ := am.GetAddr(service)
	if len(arr) != 0 {
		addr := arr[0]
		idx := strings.IndexByte(addr, ':')
		if idx != -1 {
			return strconv.Atoi(addr[idx+1:])
		}
	}

	return 0, errors.New("no config for" + service)
}

func (am *FilePortProvider) GetAddr(service string) ([]string, error) {

	addr, ok := am.config.Load(service)
	if ok {
		return []string{addr}, nil
	}
	return nil, nil
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
		// switch dest := v.(type) {
		// case string:
		f.config.Store(k, v)
		// case []any:
		// 	for _, one := range dest {
		// 		f.config.Store(k, one.(string))
		// 	}
		// default:
		// 	logger.Info("类型错误", dest)
		// }
	}

	// out, _ := yaml.Marshal(tmp)
	// os.WriteFile("config/grpc_route.yaml", out, 0644)

	return nil
}
