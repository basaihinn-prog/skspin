package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"serve/fish_comm/rng"
)

var fishInfo = &FishInfo{}

type FishInfo struct {
	FileInfo
	IsFishPath bool
	Model      map[string]interface{}
}

func (f *FishInfo) Deserialization() bool {
	var file *os.File
	var err error

	defer file.Close()
	wd, _ := os.Getwd()
	fmt.Println(filepath.Join(wd+"\\service_fish\\"+f.Folder, f.File))
	if file, err = os.Open(filepath.Join(f.Docker, f.File)); err != nil {
		if file, err = os.Open(filepath.Join(wd+"\\service_fish\\"+f.Folder, f.File)); err != nil {
			if file, err = os.Open(f.File); err != nil {
				panic(err.Error())
			}
		}
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&f.Model); err != nil {
		panic(err.Error())
	}

	if f.IsFishPath {
		for k, _ := range f.Model {
			fmt.Println(fmt.Sprintf("%s", k))
		}
	}

	return true
}

func (f *FishInfo) Reload() bool {
	for k := range f.Model {
		delete(f.Model, k)
	}

	return f.Deserialization()
}

func (f *FishInfo) Rng(options []rng.Option) interface{} {
	return rng.New(options).Item
}
