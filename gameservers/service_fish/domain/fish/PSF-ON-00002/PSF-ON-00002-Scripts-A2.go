package PSF_ON_00002

import (
	"serve/service_fish/domain/file"
)

var ScriptA2 = &scriptA2{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00002",
			File:   "PSF-ON-00002-Scripts-A2.json",
		},
	},
}

type scriptA2 struct {
	file.ScriptInfo
}

func init() {
	ScriptA2.ScriptInfo.Deserialization()
}
