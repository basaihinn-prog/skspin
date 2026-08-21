package RKF_H5_00001

import "serve/service_fish/domain/file"

var ScriptC3 = &scriptC3{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/RKF-H5-00001",
			File:   "RKF-H5-00001-Scripts-C3.json",
		},
	},
}

type scriptC3 struct {
	file.ScriptInfo
}

func init() {
	ScriptC3.ScriptInfo.Deserialization()
}
