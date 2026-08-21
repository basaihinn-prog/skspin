package RKF_H5_00001

import "serve/service_fish/domain/file"

var ScriptC1 = &scriptC1{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/RKF-H5-00001",
			File:   "RKF-H5-00001-Scripts-C1.json",
		},
	},
}

type scriptC1 struct {
	file.ScriptInfo
}

func init() {
	ScriptC1.ScriptInfo.Deserialization()
}
