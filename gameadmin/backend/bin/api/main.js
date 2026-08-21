var _curApiBtn = null
var _curApi = null

function POST(url, data, cb) {
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			let obj =  JSON.parse(this.responseText)
			cb(obj)
		}
	};
	xhttp.open("POST", url, true);

	let token = document.getElementById("token").value
	xhttp.setRequestHeader("Authorization", token)

	xhttp.setRequestHeader("pid", document.getElementById("pid").value)
	xhttp.setRequestHeader("AppID", "faketrans")
	xhttp.setRequestHeader("AppSecret", "b6337af9-a91a-4085-b1f2-466923470735")
	// xhttp.setRequestHeader("appid", document.getElementById("appid").value)

	xhttp.send(JSON.stringify(data));
}

function onExcute() {
	if (_curApi == null) {
		return
	}


	let params = document.getElementById("params")
	_curApi.Params = params.value
	commExcute(_curApi.Path, _curApi.Params)
}


function commExcute(path, params) {
	console.log("excute", path)

	// let argsObj = {}
	// argsObj.seq = new Date().getTime()


	// argsObj.params = JSON.parse(params)
	let argsObj = JSON.parse(params)

	console.log("请求参数:", argsObj)
	POST(path, argsObj, (resp)=>{
		console.log("返回:", resp)

		let result = document.getElementById("result")
		result.value = JSON.stringify(resp)
		if (!resp.error) {
			result.style.color = "white"
		}else {
			result.style.color = "red"
		}
	})

}

function wrapCommExcute(path, params) {
	return () => {
		commExcute(path, params)
	}
}


function addClickHandler(p, i) {
	p.onclick = (e) => {
		if(_curApiBtn != null) {
			_curApiBtn.style.backgroundColor = ""
		}
		_curApiBtn = e.currentTarget
		_curApi = i

		_curApiBtn.style.backgroundColor = "CornflowerBlue"

		let params = document.getElementById("params")
		params.value = i.Params

		document.getElementById("excute").innerText = "执行     " + _curApi.Path
	}
}

function renderApis(apis) {
	let left = document.getElementById("left")

	let lastKind = apis[0].Kind
	let loginBtn = null
	for (let i of apis) {
		if (lastKind != i.Kind) {
			let line = document.createElement("hr")
			// line.className = "item"
			// line.innerHTML = i.Kind
			left.append(line)
			// line.style.height = "1ch"
		}
		lastKind = i.Kind


		let p = document.createElement("button")
		p.className = "item"
		p.innerHTML = `<span class="kind">${i.Kind} ${i.Class} ${i.OnlyDev?"👀":""}</span>
<span class="path">${i.Path}</span>
<span class="desc">${i.Desc}</span>`


		addClickHandler(p, i)

		left.append(p)

		if (i.Path == "/login") {
			loginBtn = p
		}
		// left.append(document.createElement("hr"))
	}

	if (loginBtn) {
		//loginBtn.onclick()
		loginBtn.click()
	}
}

function loadApiList() {
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			let obj =  JSON.parse(this.responseText)
			renderApis(obj.data)
		}
	};
	xhttp.open("POST", "/list_api", true);
	xhttp.send("{}");
}

function onParamsChange() {
	if (_curApi != null) {
		_curApi.Params = this.value
	}
}

function main() {
	console.log("欢迎来到测试页!!")

	document.getElementById("excute").onclick = onExcute
	// document.getElementById("excutelogin").onclick = wrapCommExcute("/login", "{}")
	// document.getElementById("excutereset").onclick = wrapCommExcute("/admin/reset_player", "{}")
	// document.getElementById("excuteunlock").onclick = wrapCommExcute("/admin/unlock_events", "null")
	// document.getElementById("excuteinit").onclick = wrapCommExcute("/admin/init_items", "null")

	let params = document.getElementById("params")
	params.onchange = onParamsChange

	loadApiList()
}
