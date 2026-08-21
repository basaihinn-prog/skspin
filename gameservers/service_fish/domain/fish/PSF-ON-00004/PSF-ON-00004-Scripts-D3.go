package PSF_ON_00004

import (
	"serve/service_fish/domain/file"
)

var ScriptD3 = &scriptD3{
	ScriptInfo: file.ScriptInfo{
		Model: make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00004",
			File:   "PSF-ON-00004-Scripts-D3.json",
		},
	},
}

type scriptD3 struct {
	file.ScriptInfo
}

func init() {
	ScriptD3.ScriptInfo.Deserialization()
}
