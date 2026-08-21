package PSF_ON_00002

import (
	"serve/service_fish/domain/file"
)

var ScriptD1 = &scriptD1{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00002",
			File:   "PSF-ON-00002-Scripts-D1.json",
		},
	},
}

type scriptD1 struct {
	file.ScriptInfo
}

func init() {
	ScriptD1.ScriptInfo.Deserialization()
}
