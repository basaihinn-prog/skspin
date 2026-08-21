package RKF_H5_00001

import "serve/service_fish/domain/file"

var ScriptD2 = &scriptD2{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/RKF-H5-00001",
			File:   "RKF-H5-00001-Scripts-D2.json",
		},
	},
}

type scriptD2 struct {
	file.ScriptInfo
}

func init() {
	ScriptD2.ScriptInfo.Deserialization()
}
