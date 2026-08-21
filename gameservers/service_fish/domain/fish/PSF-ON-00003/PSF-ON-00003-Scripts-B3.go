package PSF_ON_00003

import (
	"serve/service_fish/domain/file"
)

var ScriptB3 = &scriptB3{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00003",
			File:   "PSF-ON-00003-Scripts-B3.json",
		},
	},
}

type scriptB3 struct {
	file.ScriptInfo
}

func init() {
	ScriptB3.ScriptInfo.Deserialization()
}
