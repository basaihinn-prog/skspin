package staticproxy

import (
	"bytes"
	"slices"
)

// func buildjs(w http.ResponseWriter, r *http.Request) {
// 	file := "/data/game/service/ppgateway/internal/staticproxy/build.js"
// 	http.ServeFile(w, r, file)
// 	w.Header().Add("Cache-Control", "no-cache")
// }

func buildjs_nop_host_valid(content []byte) []byte {
	const s = `(function(){
const prop = window['UHT_GAME_CONFIG_SRC']['sessionKeyV2'][0x0]['replace'](/[^0-9a-z-A-Z ]/g, '')

const o = {
	serverHostname: location.hostname,
	serverSeconds: Math.trunc ( Date.now() / 1000),
	userAgent: navigator.userAgent,
}

const str = JSON.stringify(o)
const data = new Array(234);
data.fill(32)
for(let i =0 ; i < str.length; ++i) {
	data[i] = str.charCodeAt(i)
}

Vars[prop] = data;
})();`

	i1 := bytes.Index(content, []byte("var Vars={"))
	if i1 == -1 {
		return content
	}

	i2 := bytes.Index(content[i1:], []byte("};"))
	if i2 == -1 {
		return content
	}

	pos := i1 + i2 + 2

	content = slices.Insert(content, pos, []byte(s)...)

	return content
}

func buildjs_nop_host_validB(content []byte) []byte {
	const s = `(function(){
const prop = window['UHT_GAME_CONFIG_SRC']['sessionKeyV2'][0x0]['replace'](/[^0-9a-z-A-Z ]/g, '')

const o = {
	serverHostname: location.hostname,
	serverSeconds: Math.trunc ( Date.now() / 1000),
	userAgent: navigator.userAgent,
}

const str = JSON.stringify(o)
const data = new Array(234);
data.fill(32)
for(let i =0 ; i < str.length; ++i) {
	data[i] = str.charCodeAt(i)
}

Vars[prop] = data;
})();`

	i1 := bytes.Index(content, []byte("var Vars={"))
	if i1 == -1 {
		return content
	}

	i2 := bytes.Index(content[i1:], []byte("};"))
	if i2 == -1 {
		return content
	}

	// 计算替换的位置
	start := i1
	end := i1 + i2 + 2 // 包括结束的 `};`

	// 用新代码替换旧的 Vars 定义
	content = append(content[:start], append([]byte(s), content[end:]...)...)
	return content
}

func buildjs_gtag_hook(content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte("globalTracking.QueuedTimers.length>0"), []byte(`globalTracking.QueuedTimers.length=0`))
	content = bytes.ReplaceAll(content, []byte("globalTracking.QueuedEvents.length>0"), []byte(`globalTracking.QueuedEvents.length=0`))
	return content
}

func buildjs_Number_hook(content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte("item.roundID=Number(dict[Keys.roundID])"), []byte(`item.roundID=dict[Keys.roundID]`))

	return content
}

func buildjs_Lang_hook(content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte(`ReplayConnection.watchURL=url+contextPath+"/replayGame.do"+watchQuery;`),
		[]byte(`ReplayConnection.watchURL = url + contextPath + "/replayGame.do" + watchQuery;
ReplayConnection.watchURL = ReplayConnection.watchURL + "&" + "lang=" + ServerOptions.language;`))
	content = bytes.ReplaceAll(content, []byte(`this.sharedLinkURL=url+contextPath+"/api/top/share/link"+query;`),
		[]byte(`this.sharedLinkURL = url + contextPath + "/api/top/share/link" + query;
this.sharedLinkURL = this.sharedLinkURL + "&" + "lang=" + ServerOptions.language;`))

	return content
}

func buildjs_post(content []byte) []byte {
	olds := "ResourceRequest.prototype.SendRequest=function(){"
	news := olds + `
function parseVars(data) {
	var variables = {};
	var pairs = String(data).split("&");
	for (var i = 0; i < pairs.length; ++i) {
		var pair = pairs[i].split("=");
		variables[pair[0]] = pair[1] || ""
	}
	return variables;
}
`
	content = bytes.Replace(content, []byte(olds), []byte(news), 1)

	olds = "req.onreadystatechange=function(){"
	news = olds + `
if (self.method=="POST" && this.readyState == 4 && this.status==200) {
	rawlog(self.method, self.url);
	rawlog("send >>>", parseVars(self.postData));
	rawlog("recv <<<", parseVars(this.responseText));
}
`
	content = bytes.Replace(content, []byte(olds), []byte(news), 1)
	return content
}
