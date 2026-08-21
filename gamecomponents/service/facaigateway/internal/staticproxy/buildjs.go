package staticproxy

import (
	"bytes"
)

//func buildjs_nop_host_valid(content []byte) []byte {
//	//HOST2 := gamedata.Get().MyHost
//	content = bytes.ReplaceAll(content, []byte("{{HOST2}}"), []byte(HOST2))
//	return content
//}

func RemoveHTML(content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte("jdbdlbase."), []byte("jdbdl."))
	return content
}
