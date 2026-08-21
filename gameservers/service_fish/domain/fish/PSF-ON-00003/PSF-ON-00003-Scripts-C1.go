package PSF_ON_00003

import (
	"serve/service_fish/domain/file"
)

var ScriptC1 = &scriptC1{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00003",
			File:   "PSF-ON-00003-Scripts-C1.json",
		},
	},
}

type scriptC1 struct {
	file.ScriptInfo
}

func init() {
	ScriptC1.ScriptInfo.Deserialization()
}
