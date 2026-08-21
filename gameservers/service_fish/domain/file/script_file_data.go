package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var scriptInfo = &ScriptInfo{}

type ScriptInfo struct {
	FileInfo
	Model map[string]interface{}
}

func (s *ScriptInfo) Deserialization() bool {
	var file *os.File
	var err error

	defer file.Close()
	wd, _ := os.Getwd()
	fmt.Println(filepath.Join(wd+"\\service_fish\\"+s.Folder, s.File))
	if file, err = os.Open(filepath.Join(s.Docker, s.File)); err != nil {
		if file, err = os.Open(filepath.Join(wd+"\\service_fish\\"+s.Folder, s.File)); err != nil {
			if file, err = os.Open(s.File); err != nil {
				panic(err.Error())
			}
		}
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&s.Model); err != nil {
		panic(err.Error())
	}

	return true
}

func (s *ScriptInfo) Reload() bool {
	for k := range s.Model {
		delete(s.Model, k)
	}

	return true
}

func (s *ScriptInfo) Data() map[string]interface{} {
	// TODO Deserialization have nil interface
	if len(s.Model) == 0 {
		s.Deserialization()
	}
	return s.Model["Scripts"].(map[string]interface{})
}

func (s *ScriptInfo) Size() int {
	return len(s.Model["Scripts"].(map[string]interface{}))
}

func (s *ScriptInfo) ResetTime() int {
	data := s.Model["Transform scene"].(map[string]interface{})
	return int(data["Time"].(float64))
}
