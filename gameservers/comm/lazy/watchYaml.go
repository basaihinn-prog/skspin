package lazy

import (
	"encoding/json"
	"log"
	"path"
	"reflect"

	"serve/comm/ut"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

func wrapYaml(f any) LoadFuncWithName {
	fn := reflect.ValueOf(f)
	t := reflect.TypeOf(f)
	paramsType := t.In(0)

	lo.Must0(t.NumIn() == 1)
	lo.Must0(t.NumOut() == 1)
	u := reflect.TypeOf((*error)(nil)).Elem()
	lo.Must0(t.Out(0).Implements(u))
	return func(file string, buf []byte) error {
		paramsValue, err := ut.Bindwrap(paramsType, func(i interface{}) error {
			extname := path.Ext(file)
			switch extname {
			case ".json":
				return json.Unmarshal(buf, i)
			case ".yml", ".yaml":
				return yaml.Unmarshal(buf, i)
			default:
				// log.Panic("WatchUnmarshal Ext invalid!", "file", file)
				panic(extname)
			}
		})
		if err != nil {
			return err
		}

		out := fn.Call([]reflect.Value{paramsValue})
		if !out[0].IsNil() {
			err = out[0].Interface().(error)
			return err
		}

		return nil
	}
}

// func(set *Setting) error
func (m *ConfigManager) WatchUnmarshal(file string, f any) {
	switch path.Ext(file) {
	case ".json", ".yml", ".yaml":
	default:
		log.Panic("WatchUnmarshal Ext invalid!", "file", file)
	}
	m.elements = append(m.elements, &LoadData{
		path:     file,
		isprefix: false,
		callback: wrapYaml(f),
	})
}
