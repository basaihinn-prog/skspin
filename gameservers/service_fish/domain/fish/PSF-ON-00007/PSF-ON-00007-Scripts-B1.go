package PSF_ON_00007

import (
	"serve/service_fish/domain/file"
)

var ScriptB1 = &scriptB1{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00007",
			File:   "PSF-ON-00007-Scripts-B1.json",
		},
	},
}

type scriptB1 struct {
	file.ScriptInfo
}

func init() {
	ScriptB1.ScriptInfo.Deserialization()
}
