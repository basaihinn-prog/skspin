package PSF_ON_00005

import (
	"serve/service_fish/domain/file"
)

var ScriptC2 = &scriptC2{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00005",
			File:   "PSF-ON-00005-Scripts-C2.json",
		},
	},
}

type scriptC2 struct {
	file.ScriptInfo
}

func init() {
	ScriptC2.ScriptInfo.Deserialization()
}
