package gamedata

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Setting struct {
	RobotFlag bool `json:"RobotFlag"`
}

var Settings *Setting

func Load() {
	// 读取 YAML 文件
	data, err := ioutil.ReadFile("config/minigateway_config.yaml")
	if err != nil {
		log.Fatalf("无法读取文件: %v", err)
	}
	// 解析 YAML 数据
	Settings = &Setting{}
	err = yaml.Unmarshal(data, Settings)
	if err != nil {
		log.Fatalf("无法解析 YAML: %v", err)
	}
}
