package PSF_ON_00001

import (
	"serve/service_fish/domain/file"
)

var ScriptA1 = &scriptA1{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00001",
			File:   "PSF-ON-00001-Scripts-A1.json",
		},
	},
}

type scriptA1 struct {
	file.ScriptInfo
}

func init() {
	ScriptA1.ScriptInfo.Deserialization()
}
