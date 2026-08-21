package PSF_ON_00004

import (
	"serve/service_fish/domain/file"
)

var ScriptB2 = &scriptB2{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00004",
			File:   "PSF-ON-00004-Scripts-B2.json",
		},
	},
}

type scriptB2 struct {
	file.ScriptInfo
}

func init() {
	ScriptB2.ScriptInfo.Deserialization()
}
