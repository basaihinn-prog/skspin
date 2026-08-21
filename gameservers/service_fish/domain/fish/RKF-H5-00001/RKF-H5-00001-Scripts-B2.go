package RKF_H5_00001

import "serve/service_fish/domain/file"

var ScriptB2 = &scriptB2{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/RKF-H5-00001",
			File:   "RKF-H5-00001-Scripts-B2.json",
		},
	},
}

type scriptB2 struct {
	file.ScriptInfo
}

func init() {
	ScriptB2.ScriptInfo.Deserialization()
}
