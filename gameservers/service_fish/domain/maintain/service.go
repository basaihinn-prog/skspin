package maintain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var Service = &service{
	Docker: "/data/sys",
	test:   "../maintain",
	folder: "domain/maintain",
	File:   "maintain.json",
}

type service struct {
	Docker   string
	test     string
	folder   string
	File     string
	maintain *Maintain
}

func init() {
	Service.deserialization()
}

func (s *service) deserialization() bool {
	var file *os.File
	var err error

	defer file.Close()

	wd, _ := os.Getwd()
	fmt.Println(filepath.Join(wd+"\\service_fish\\"+s.folder, s.File))
	if file, err = os.Open(filepath.Join(s.Docker, s.File)); err != nil {
		if file, err = os.Open(filepath.Join(wd+"\\service_fish\\"+s.folder, s.File)); err != nil {
			if file, err = os.Open(s.File); err != nil {
				panic(err.Error())
			}
		}
	}

	jsonParser := json.NewDecoder(file)

	if err = jsonParser.Decode(&s.maintain); err != nil {
		panic(err.Error())
	}
	return true
}

func (s *service) Reload() bool {
	s.maintain = nil
	return s.deserialization()
}

func (s *service) Shutdown() bool {
	return s.maintain.ShutDown
}

func (s *service) Message() string {
	return s.maintain.Message
}

func (s *service) DatabaseReload() bool {
	return s.maintain.DatabaseReload
}
