package RKF_H5_00001

import "serve/service_fish/domain/file"

var ScriptB3 = &scriptB3{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/RKF-H5-00001",
			File:   "RKF-H5-00001-Scripts-B3.json",
		},
	},
}

type scriptB3 struct {
	file.ScriptInfo
}

func init() {
	ScriptB3.ScriptInfo.Deserialization()
}
