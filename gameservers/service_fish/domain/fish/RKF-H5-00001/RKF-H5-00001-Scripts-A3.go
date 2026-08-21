package RKF_H5_00001

import "serve/service_fish/domain/file"

var ScriptA3 = &scriptA3{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/RKF-H5-00001",
			File:   "RKF-H5-00001-Scripts-A3.json",
		},
	},
}

type scriptA3 struct {
	file.ScriptInfo
}

func init() {
	ScriptA3.ScriptInfo.Deserialization()
}
