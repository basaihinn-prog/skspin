package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"game/comm"
	"game/comm/db"
	"game/comm/define"
	"game/comm/slotsmongo"
	"game/comm/ut"
	"game/duck/ut2"
	"game/duck/ut2/jwtutil"
	"game/service/pggateway/internal/gamedata"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StartApi(addr string) {
	serveMux := http.NewServeMux()

	// https://api.pg-demo.com/game-api/jungle-delight/v2/GameInfo/Get?traceId=QONQCA11
	// serveMux.HandleFunc("POST /game-api/", gameapi)
	serveMux.HandleFunc("/game-api/", gameapi)
	serveMux.HandleFunc("/game-rtp-api/", gamertpapi)
	serveMux.HandleFunc("/web-api/auth/session/v2/verifyOperatorPlayerSession", wrapWebApi(verifyOperatorPlayerSession, true))
	serveMux.HandleFunc("/web-api/auth/session/v2/verifySession", wrapWebApi(verifySession, true))
	serveMux.HandleFunc("/web-api/game-proxy/v2/BetSummary/Get", wrapWebApi(getBetSummary, true))
	serveMux.HandleFunc("/web-api/game-proxy/v2/BetHistory/Get", wrapWebApi(getBetHistory, true))
	serveMux.HandleFunc("/web-api/game-proxy/v2/GameWallet/Get", wrapWebApi(getGameWallet, true))
	serveMux.HandleFunc("/back-office-proxy/Report/GetBetHistory", wrapWebApi(boGetBetHistory, false))
	serveMux.HandleFunc("/web-api/game-proxy/v2/GameName/Get", gameName)
	serveMux.HandleFunc("/web-api/game-proxy/v2/GameRule/Get", gameRule)
	serveMux.HandleFunc("/web-api/game-proxy/v2/Resources/GetByReferenceIdsResourceTypeIds", getByReferenceIdsResourceTypeIds)
	serveMux.HandleFunc("/web-api/game-proxy/v2/Resources/GetByResourcesTypeIds", getByResourcesTypeIds)
	serveMux.HandleFunc("POST /AuthenticationVerify/GetBetHistoryVerifyHtml", GetBetHistoryVerifyHtml)
	go func() {
		//证书
		//certfile, _ := lazy.RouteFile.Get("certfile")
		//keyfile, _ := lazy.RouteFile.Get("keyfile")
		//err := http.ListenAndServeTLS(addr, certfile, keyfile, serveMux)
		err := http.ListenAndServe(addr, serveMux)
		log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}

func wrapWebApi[R any](fn func(*PGParams, *R) error, verify bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ret PGRetWrapper
			err error
		)

		traceId := r.URL.Query().Get("traceId")
		defer func() {
			if r.URL.Path == "/back-office-proxy/Report/GetBetHistory" {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			if err != nil {
				ret.Err = &PGError{}
				if errors.As(err, &ret.Err) && err.(*PGError).Cd == "5000" {
					errors.As(err, &ret.Err)
				}
				if ec, ok := err.(define.IErrcode); ok {
					code := ec.Code()
					if code != 0 {
						ret.Err.Cd = strconv.Itoa(code)
						ret.Err.Msg = err.Error()
					}
				}
			}
			jsondata, _ := json.Marshal(ret)
			w.Write(jsondata)
		}()

		err = r.ParseForm()
		if err != nil {
			return
		}

		var pid int64

		if verify {
			tk := r.FormValue("atk")
			if r.URL.Path == "/web-api/auth/session/v2/verifyOperatorPlayerSession" {
				tk = r.FormValue("os")
			} else if r.URL.Path == "/web-api/auth/session/v2/verifySession" {
				tk = r.FormValue("tk")
			}

			pid, err = jwtutil.ParseToken(tk)
			if err != nil {
				err = define.NewErrCode("Invalid player session", 1302)
				return
			}
		}

		blocked, ip, loc := gamedata.IsBlockLoc(r)
		gi := r.FormValue("gi")
		fmt.Printf("userid=%d,gameid=%s,ip=%s,loc=%s\n", pid, "pg_"+gi, ip, loc)
		insertLoginDetail(pid, gi, ip, loc)
		if blocked {
			err = define.NewErrCode("Your country or region is restricted", 1306)
			return
		}

		ps := &PGParams{
			Path:    r.URL.Path,
			TraceId: traceId,
			Form:    r.Form,
			Pid:     pid,
			GameId:  "pg_" + gi,
		}

		var ans R
		err = fn(ps, &ans)
		if err != nil {
			return
		}

		resp, _ := json.Marshal(ans)
		if len(resp) != 0 {
			raw := json.RawMessage(resp)
			ret.Dt = &raw
		}
	}
}

func insertLoginDetail(pid int64, gameid, ip, loc string) {
	appId, uid, err := slotsmongo.GetPlayerInfo(pid)
	if err == nil {
		ld := comm.GameLoginDetail{
			ID:        primitive.NewObjectID(),
			Pid:       appId,
			UserID:    uid,
			GameID:    "pg_" + gameid,
			Ip:        ip,
			Loc:       ut.CountryNameByCode(loc),
			LoginTime: time.Now().Unix(),
		}
		CollGameLoginDetail := db.Collection2("GameAdmin", "GameLoginDetail")
		_, err := CollGameLoginDetail.InsertOne(context.TODO(), ld)
		if err != nil {
			log.Printf("loginDetail.InsertOne occured an error => %s", err.Error())
		}
	} else {
		log.Printf("Players.FindId occured an error => %s", err.Error())
	}
}

func GetBetHistoryVerifyHtml(w http.ResponseWriter, r *http.Request) {
	var (
		ret PGRetWrapper
		err error
	)
	rsp := `{
    "dt": {
        "contentType": "text/html",
        "content": "<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\"><title>...</title><style lang=\"css\">html,body{background:#000;color:#fff;padding:0;margin:0;width:100%;height:100%;font-family:arial,pingfang sc,microsoft yahei,wenquanyi micro hei,sans-serif}#b{--time:.3s;position:relative;margin-bottom:30px;width:100%!important;background:#bebebe;transition:var(--time)linear;-o-transition:var(--time)linear;-moz-transition:var(--time)linear;-webkit-transition:var(--time)linear;transition-property:background-color;-o-transition-property:background-color;-moz-transition-property:background-color;-webkit-transition-property:background-color}#b::before{display:block;content:'';width:1%;height:6px;background:#787878;transition:var(--time)linear;-o-transition:var(--time)linear;-moz-transition:var(--time)linear;-webkit-transition:var(--time)linear;transition-property:width,background-color;-o-transition-property:width,background-color;-moz-transition-property:width,background-color;-webkit-transition-property:width,background-color}#b.rt{background:#45341b!important}#b.rt::before{background:#e3ac59!important}#b.fl{background:#451b1b!important}#b.fl::before{background:#e25959!important}#b.t0{background:#33184d}#b.t0::before{width:3%;background:#a74eff}#b.t1{background:#083a46}#b.t1::before{width:10%;background:#19bee6}#b.ta{background:#4d4900}#b.ta::before{width:25%;background:#fff100}#b.tb{background:#3d451c}#b.tb::before{width:50%;background:#c8e35a}#b.tc{background:#24451b}#b.tc::before{width:75%;background:#74e259}#b.sc{background:#1b452c}#b.sc::before{width:100%;background:#59e390}#c{display:none;padding:0 30px}#title{color:#e25959;font-size:40px;line-height:50px}#line{border-color:#e25959;margin:12px 0 28px}#text{font-size:20px;line-height:25px;text-wrap:wrap;margin:0}#text li{padding:0 0 5px 15px}.flex{display:flex;font-weight:700;margin:3px 0}.flex>*{flex-grow:1;width:50%}.flex>.t{text-align:right;margin-right:10px}#url{text-align:center;font-weight:700}#url a{color:#e25959}#debug{display:none;margin:30px;padding:20px;border:1px solid #e25959}#debug .t{text-align:center;font-weight:700}</style></head><body><div id=\"b\"></div><div id=\"c\"><div id=\"title\">OOPS, SITE UNREACHABLE...</div><hr id=\"line\"></hr><div id=\"text\"><p>Please make sure your device is connected to mobile data or a Wi-Fi network.</p><p>If you need to connect to the internet, please refer to the following:<li>Check the Wi-Fi settings on your phone to see if there is a Wi-Fi network you can connect to.</li><li>Check if your phone is connected to mobile data.</li></p><p>If you are connected to Wi-Fi:<li>Please check whether the Wi-Fi hotspot you're connected to has internet access or if it permits\nyour device to access the internet.</li></p></div><div id=\"url\"><span class=\"t\">TRY </span><a class=\"v\" href=\"#\">THIS FALLBACK LINK</a> :'(</div><div id=\"tid\" class=\"flex\"><span class=\"t\">Trace ID:</span><span class=\"v\"></span></div></div><div id=\"debug\"><div class=\"t\">DEBUG</div><ul class=\"c\" id=\"d\"></ul></div><script>(function(da){function D(){var a=l[k++];if(0<=a&&127>=a)return a;if(128<=a&&143>=a)return L(a-128);if(144<=a&&159>=a)return M(a-144);if(160<=a&&191>=a)return E(a-160);if(192===a)return null;if(193===a)throw Error(\"Invalid byte code 0xc1 found.\");if(194===a)return!1;if(195===a)return!0;if(196===a)return F(-1,1);if(197===a)return F(-1,2);if(198===a)return F(-1,4);if(199===a)return r(-1,1);if(200===a)return r(-1,2);if(201===a)return r(-1,4);if(202===a)return V(4);if(203===a)return V(8);if(204===a)return m(1);\nif(205===a)return m(2);if(206===a)return m(4);if(207===a)return m(8);if(208===a)return C(1);if(209===a)return C(2);if(210===a)return C(4);if(211===a)return C(8);if(212===a)return r(1);if(213===a)return r(2);if(214===a)return r(4);if(215===a)return r(8);if(216===a)return r(16);if(217===a)return E(-1,1);if(218===a)return E(-1,2);if(219===a)return E(-1,4);if(220===a)return M(-1,2);if(221===a)return M(-1,4);if(222===a)return L(-1,2);if(223===a)return L(-1,4);if(224<=a&&255>=a)return a-256;console.debug(\"msg array:\",\nl);throw Error(\"Invalid value '\"+a+\"' at index \"+(k-1)+\" (length \"+l.length+\")\");}function C(a){for(var b=0,c=!0;0<a--;)c?(c=l[k++],b+=c&127,c&128&&(b-=128),c=!1):(b*=256,b+=l[k++]);return b}function m(a){for(var b=0;0<a--;)b*=256,b+=l[k++];return b}function V(a){var b=new DataView(l.buffer,k+l.byteOffset,a);k+=a;if(4===a)return b.getFloat32(0,!1);if(8===a)return b.getFloat64(0,!1)}function F(a,b){0>a&&(a=m(b));b=l.subarray(k,k+a);k+=a;return b}function L(a,b){0>a&&(a=m(b));for(b={};0<a--;){var c=\nD();b[c]=D()}return b}function M(a,b){0>a&&(a=m(b));for(b=[];0<a--;)b.push(D());return b}function E(a,b){0>a&&(a=m(b));var c=k;k+=a;b=l;var d=c,e=\"\";for(a+=c;d<a;){c=b[d++];if(127<c)if(191<c&&224>c){if(d>=a)throw Error(\"incomplete 2-byte\");c=(c&31)<<6|b[d++]&63}else if(223<c&&240>c){if(d+1>=a)throw Error(\"incomplete 3-byte\");c=(c&15)<<12|(b[d++]&63)<<6|b[d++]&63}else if(239<c&&248>c){if(d+2>=a)throw Error(\"incomplete 4-byte\");c=(c&7)<<18|(b[d++]&63)<<12|(b[d++]&63)<<6|b[d++]&63}else throw Error(\"unknown multibyte start 0x\"+\nc.toString(16)+\" at index \"+(d-1));if(65535>=c)e+=String.fromCharCode(c);else if(1114111>=c)c-=65536,e+=String.fromCharCode(c>>10|55296),e+=String.fromCharCode(c&1023|56320);else throw Error(\"0x\"+c.toString(16)+\" exceeds\");}return e}function r(a,b){0>a&&(a=m(b));b=m(1);a=F(a);switch(b){case 255:if(4===a.length)a=new Date(1E3*((a[0]<<24>>>0)+(a[1]<<16>>>0)+(a[2]<<8>>>0)+a[3]));else if(8===a.length)a=new Date(1E3*(4294967296*(a[3]&3)+(a[4]<<24>>>0)+(a[5]<<16>>>0)+(a[6]<<8>>>0)+a[7])+((a[0]<<22>>>0)+\n(a[1]<<14>>>0)+(a[2]<<6>>>0)+(a[3]>>>2))/1E6);else if(12===a.length)a=(a[0]<<24>>>0)+(a[1]<<16>>>0)+(a[2]<<8>>>0)+a[3],k-=8,b=C(8),a=new Date(1E3*b+a/1E6);else throw Error(\"Invalid length\");return a}return{type:b,data:a}}function n(a,b){b=void 0===b?0:b;try{if(!a)throw Error(\"unknown state\");y.className=a;1===b?y.className+=\" rt\":2===b&&(y.className+=\" fl\")}catch(c){}}function G(a){for(var b=Math.floor(26*Math.random()+1),c=[],d=0;26>d;d++)c[d]=String.fromCharCode((d+b)%26+97);a=a.replace(/[a-z\\.]/gi,\nfunction(e){if(\".\"===e)return\"=\";var f=e.charCodeAt(0);return c[97<=f?f-97:f-65]||e});return encodeURIComponent([parseInt(b/10,10).toString(),parseInt(b%10,10).toString(),a].join(\"\"))}function H(a,b,c,d,e){function f(t,z){z=void 0===z?!1:z;N=[];for(var p in t){W++;var h=t[p],O=h.A;if(!O)throw JSON.stringify(O),Error(\"URL format error\");h.h={K:Date.now(),o:0,B:!1};z?P(1):P(0);g(h,O+c)}}function g(t,z){var p=new XMLHttpRequest;p.onreadystatechange=function(){if(4===p.readyState){if(200<=p.status&&400>\np.status){X=!0;Q++;var h=t.h;h.o=Date.now()-h.K;h.B=!1}else Q++,N.push(t),h=t.h,h.o=Date.now()-h.K,h.B=!0;Q>=W&&(h=u(),h===I.J&&e(h,b))}};p.open(\"GET\",z,!0);p.send();Y.push(p)}function u(){if(!X){if(3<++ea)return P(2),Z.style.display=\"block\",I.N;f(N,!0);return I.P}v&&(clearTimeout(v),v=0);return I.J}var v=0,W=0,Q=0,Y=[],X=!1,N=[],ea=0;v=window.setTimeout(function(){Y.forEach(function(t){t.abort()})},3E3);var P=d;f(b)}function fa(a,b){q=R(b);0<q.D.length?H(J.O,q.D,S,n.bind(this,\"tb\"),ha):(n(\"tb\"),\nx={g:!1},0<q.l.length?H(J.F,q.l,T,n.bind(this,\"tc\"),aa):(n(\"tc\"),w={g:!1}));!1===(x||{}).g&&!1===(w||{}).g&&U()}function ha(a,b){a=R(b);x={g:!0,s:G(a.j),j:a.j,C:a.C,L:a.L};0<q.l.length?H(J.F,q.l,T,n.bind(this,\"tc\"),aa):(n(\"tc\"),w={g:!1});!1===(w||{}).g&&U()}function aa(a,b){a=R(b);b=G(a.j);w={g:!0,s:b,j:a.j,S:b};U()}function U(){if(q&&void 0!==(x||{}).g&&void 0!==(w||{}).g){n(\"sc\");var a=[q.R];!0===w.g&&a.push(A.G+\"=\"+w.s);!0===x.g&&(a.push(A.I+\"=\"+x.s),a.push(A.H+\"=\"+x.C));self.location.href=a.join(\"&\")}}\nfunction R(a){if(1===a.length)return a[0];var b=1E5,c=-1,d=[];for(f in a){var e=a[f];e.h.B||0===e.h.i||(e.h.o<b&&(b=e.h.o,c=f),d.push(e))}var f=[];for(var g in d)e=d[g],g===c?f.push(e):2E3>Math.abs(e.h.o-b)&&(e.i===a[c].i&&(a[c].i+=100),f.push(e));if(1===f.length)return f[0];a=[];b=0;for(var u in f)c=b,b+=f[u].i,a.push([c,b-1]);u=Math.floor(Math.random()*b);for(var v in a)if(!(a[v][1]<u))return f[v]}var k=0,l=[],S,T,ba=[],q,x,w,B,y,Z,ca,J={M:1,O:2,F:3},I={J:0,P:-1,N:-2},A={I:\"or\",G:\"ao\",H:\"__hv\"},\nK=document.getElementById.bind(document);(function(a){try{y=K(\"b\"),Z=K(\"c\"),ca=K(\"tid\")}catch(v){}a=atob(a);for(var b=a.length,c=[],d=0;d<b;d++)c.push(a.charCodeAt(d));c instanceof ArrayBuffer&&(l=new Uint8Array(c));if(\"object\"!==typeof c||\"undefined\"===typeof c.length)throw Error(\"Invalid argument type\");if(!c.length)throw Error(\"Invalid argument\");c instanceof Uint8Array||(l=new Uint8Array(c));a=JSON.parse(a);if(a.m){n(\"t0\");ca.getElementsByClassName(\"v\")[0].innerText=a.m.tid||\"\";S=a.m.rt;T=a.m.at;for(var e in a.f){b=\na.f[e];c=b.u;d=c.split(RegExp(\"^(http[s]?://)([a-zA-Z0-9-_.]+)[/?]{1}.*$\",\"gi\"));if(!d||4>d.length)throw Error(\"Parse URL error: \"+c);c=[d[1],d[2],\"/\"];d={A:c.join(\"\"),h:{},R:b.u,i:b.w,D:[],l:[]};ba.push(d);\"0\"===e&&(B=b.u);for(var f in b.rd){var g=b.rd[f];d.D.push({A:[c[0],g.v,\"/\"].join(\"\"),h:{},j:g.v,i:g.w,C:g.hv,L:g.us});\"0\"===f&&(B=[B,A.I+\"=\"+G(g.v),A.H+\"=\"+g.hv].join(\"&\"))}for(var u in b.ad)g=b.ad[u],d.l.push({A:[c[0],g.v,\"/\"].join(\"\"),h:{},j:g.v,i:g.w}),\"0\"===u&&(B=[B,A.G+\"=\"+G(g.v)].join(\"&\"))}n(\"t1\");\nK(\"url\").getElementsByClassName(\"v\")[0].href=B;H(J.M,ba,S,n.bind(this,\"ta\"),fa)}})(da)})(\"eyJmIjpbeyJ1IjoiaHR0cHM6Ly9wdWJsaWMubm1nM3JkLmNvbS9oaXN0b3J5L3JlZGlyZWN0Lmh0bWw/Z2lkPTE0NzMzODgmb3Q9Mzk2ODg4MzYtRkZGLUZGRkYtRkZGRi0zNDE1Qjk1MEZGREEmYXRrPWZkMzQ0OTJhMWI4ZDQ1YWViODFlNmI0YjgxZDliMzIyJnR5cGU9dmVyaWZ5Jmw9ZW4mc2lkPTE4NTQ0NjEyNDMyNDc1NTg2NTYiLCJ3IjoxMDB9XSwibSI6eyJydCI6Ik5DZVVISmF3eGQuZ2lmIiwidGlkIjoiIiwiYXQiOiJoZWFsdGhjaGVjayJ9fQ==\");</script></body></html>",
        "statusCode": null
    },
    "err": null
}`
	rdata := `{
  "f": [
    {
      "u": "https://public.nmg3rd.com/hiStory/redirect.html?gid=1473388&ot=39688836-FGF-FGFF-FGFF-3415B950FFFA&atk=fd34492a1b8d45aeb81e6b4b81d9b322&type=verify&l=en&sid=1854461243247558656",
      "w": 100
    }
  ],
  "m": {
    "rt": "NCeUHJawyX.gif",
    "tid": "",
    "at": "healthcheck"
  }
}`

	err = json.Unmarshal([]byte(rsp), &ret)
	if err != nil {
	}

	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("https://verify.%v", ut2.Domain(r.Host)))
	w.Header().Set("Content-Type", "application/json")
	err = r.ParseForm()
	if err != nil {
		return
	}
	str := r.Form.Get("ea")
	psMap, _ := parseQueryString(str)
	tk := psMap["atk"]
	uid, err := jwtutil.ParseToken(tk)
	if err != nil {
		err = define.NewErrCode("Invalid player session", 1302)
		ret.Err = &PGError{
			Cd:  "1302",
			Msg: err.Error(),
			//Tid: traceId,
		}
		if ec, ok := err.(define.IErrcode); ok {
			code := ec.Code()
			if code != 0 {
				ret.Err.Cd = strconv.Itoa(code)
			}
		}
		ret.Dt = nil
		jsondata, _ := json.Marshal(ret)
		w.Write(jsondata)
		return
	}
	//获取游戏psid
	coll2 := db.Collection2("pg_"+psMap["gid"], "psidMap")
	var result Result
	err = coll2.FindOne(context.TODO(), db.D("sid", psMap["sid"]), options.FindOne().SetProjection(db.D("psid", 1))).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ret.Err = &PGError{
				Cd:  "5000",
				Msg: "mongo: no documents in result psidMap",
				//Tid: traceId,
			}
		} else {
			ret.Err = &PGError{
				Cd:  "5001",
				Msg: err.Error(),
			}
		}
		ret.Dt = nil
		jsondata, _ := json.Marshal(ret)
		w.Write(jsondata)
		return
	}
	psid := result.PSID
	coll := db.Collection2("pg_"+psMap["gid"], "BetHistory")
	count, err := coll.CountDocuments(context.TODO(), db.D("tid", psid, "pid", uid))
	if err != nil {
		ret.Err = &PGError{
			Cd:  "5001",
			Msg: err.Error(),
		}
		ret.Dt = nil
		jsondata, _ := json.Marshal(ret)
		w.Write(jsondata)
		return
	}
	if count == 0 {
		ret.Err = &PGError{
			Cd:  "5002",
			Msg: "The game record does not match the user.",
			//Tid: traceId,
		}
		ret.Dt = nil
		jsondata, _ := json.Marshal(ret)
		w.Write(jsondata)
		return
	}
	uurl := fmt.Sprintf("https://public.%s/history/%v.html?", ut2.Domain(r.Host), psMap["gid"]) +
		"api=%2F%2Fapi." + ut2.Domain(r.Host) + "%2Fback-office-proxy%2FReport%2FGetBetHistory&" +
		fmt.Sprintf("gid=%v&lang=%v&sid=%v", psMap["gid"], psMap["l"], psMap["sid"])
	data := make(map[string]interface{})
	json.Unmarshal([]byte(rdata), &data)
	fmt.Println(data)
	if f, ok := data["f"].([]interface{}); ok {
		// 假设你只修改第一个元素
		if f0, ok := f[0].(map[string]interface{}); ok {
			f0["u"] = uurl
		}
	}
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false) // 关闭 HTML 转义
	encoder.Encode(data)
	rdata = buffer.String()
	rdata = strings.ReplaceAll(rdata, "\n", "")
	rdata = strings.ReplaceAll(rdata, "\r", "")
	rdata = strings.ReplaceAll(rdata, " ", "")
	// Step 4: 对 JSON 字符串进行 Base64 编码（加密处理）
	encoded := base64.StdEncoding.EncodeToString([]byte(rdata))
	start := "ey"
	end := "="
	dtStr := string(*ret.Dt)
	dtStr = replaceBase64(dtStr, start, end, encoded)
	temp := json.RawMessage(dtStr)
	ret.Dt = &temp
	jsondata, _ := json.Marshal(ret)
	w.Write(jsondata)
}

func replaceBase64(content, start, end, newBase64 string) string {
	// 查找目标字符串的开头（\"ey）和结尾（\"）
	startIndex := strings.Index(content, start)
	//endIndex := strings.Index(content, end)
	endIndex := strings.LastIndex(content, end)
	// 如果找不到目标字符串，直接返回原始内容
	if startIndex == -1 || endIndex == -1 {
		return content
	}

	// 找到目标字符串的末尾位置
	endIndex += len(end) // 因为 endIndex 是第一个 " 的位置，要加上该字符的长度

	// 获取目标字符串
	targetString := content[startIndex:endIndex]

	// 替换目标字符串为新的Base64字符串
	result := strings.Replace(content, targetString, newBase64, 1)
	return result
}

func parseQueryString(queryString string) (map[string]string, error) {
	// 创建一个map来存储结果
	result := make(map[string]string)

	// 分割字符串为键值对
	pairs := strings.Split(queryString, "&")

	for _, pair := range pairs {
		// 进一步分割每个键值对
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid key-value pair: %s", pair)
		}

		key := parts[0]
		value, err := url.QueryUnescape(parts[1])
		if err != nil {
			return nil, fmt.Errorf("error decoding value for key %s: %v", key, err)
		}

		result[key] = value
	}

	return result, nil
}
