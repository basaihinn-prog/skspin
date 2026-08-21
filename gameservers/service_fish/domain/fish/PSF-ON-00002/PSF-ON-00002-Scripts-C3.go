package PSF_ON_00002

import (
	"serve/service_fish/domain/file"
)

var ScriptC3 = &scriptC3{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00002",
			File:   "PSF-ON-00002-Scripts-C3.json",
		},
	},
}

type scriptC3 struct {
	file.ScriptInfo
}

func init() {
	ScriptC3.ScriptInfo.Deserialization()
}
