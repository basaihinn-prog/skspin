package PSF_ON_00005

import (
	"serve/service_fish/domain/file"
)

var ScriptB2 = &scriptB2{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00005",
			File:   "PSF-ON-00005-Scripts-B2.json",
		},
	},
}

type scriptB2 struct {
	file.ScriptInfo
}

func init() {
	ScriptB2.ScriptInfo.Deserialization()
}
