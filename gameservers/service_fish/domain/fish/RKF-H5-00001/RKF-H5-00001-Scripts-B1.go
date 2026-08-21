package RKF_H5_00001

import "serve/service_fish/domain/file"

var ScriptB1 = &scriptB1{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/RKF-H5-00001",
			File:   "RKF-H5-00001-Scripts-B1.json",
		},
	},
}

type scriptB1 struct {
	file.ScriptInfo
}

func init() {
	ScriptB1.ScriptInfo.Deserialization()
}
