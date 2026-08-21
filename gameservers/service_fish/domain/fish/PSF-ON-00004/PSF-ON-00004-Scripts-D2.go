package PSF_ON_00004

import (
	"serve/service_fish/domain/file"
)

var ScriptD2 = &scriptD2{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00004",
			File:   "PSF-ON-00004-Scripts-D2.json",
		},
	},
}

type scriptD2 struct {
	file.ScriptInfo
}

func init() {
	ScriptD2.ScriptInfo.Deserialization()
}
