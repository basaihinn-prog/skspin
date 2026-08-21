package staticproxy

import (
	"bytes"
	"game/service/jdbgateway/internal/gamedata"
)

func buildjs_nop_host_valid(content []byte) []byte {
	HOST2 := gamedata.Get().MyHost
	content = bytes.ReplaceAll(
		content,
		[]byte(`i.prototype.nextAutoSpin=function(t){`),
		[]byte(`i.prototype.nextAutoSpin = function(t) { 
				if(this.voData.gameState=="BASE_SHOW_MONEY")this.voData.gameState="BASE_IDLE";`),
	)

	content = bytes.ReplaceAll(content, []byte("this.data.curExtraBet=e"), []byte(`if(e.indexOf("0_")!=-1){e=e.split("0_")[1]}this.data.curExtraBet=e`))

	content = bytes.ReplaceAll(content, []byte("{{HOST2}}"), []byte(HOST2))
	return content
}
func buildjs_nop_host_valid2(content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte("jdbdlbase."), []byte("jdbdl."))
	return content
}
