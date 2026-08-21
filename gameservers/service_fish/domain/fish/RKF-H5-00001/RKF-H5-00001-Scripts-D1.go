package RKF_H5_00001

import "serve/service_fish/domain/file"

var ScriptD1 = &scriptD1{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/RKF-H5-00001",
			File:   "RKF-H5-00001-Scripts-D1.json",
		},
	},
}

type scriptD1 struct {
	file.ScriptInfo
}

func init() {
	ScriptD1.ScriptInfo.Deserialization()
}
