// ee.cshProto.ExchangeReq
globalThis.myinject = function (proto) {
    for(const name of Object.getOwnPropertyNames(proto)) {
        const o = proto[name]
        o._raw_encode =  o._raw_encode || o.encode
        o.encode = (...ps)=>{
            console.log(">> encode before", name, ...ps)
            const r = o._raw_encode(...ps)
            console.log(">> encode after", name, r)
            return r
        }

        o._raw_decode = o._raw_decode || o.decode
        o.decode = (...ps)=>{
            console.log("<< decode before", name, ...ps)
            const r = o._raw_decode(...ps)
            console.log("<< decode after", name, r)
            return r
        }
    }
}

// file:///Users/a123/Desktop/jilitmp/uat-wbgame.jlfafafa3.com/aa/assets/game/index.2d197.js

// m.Login = function(t) {

    // globalThis.myinject = function (o, name) {
    //     // const o = proto[name]
    //     o._raw_encode =  o._raw_encode || o.encode
    //     o.encode = (...ps)=>{
    //         console.log(">> encode before", name, ...ps)
    //         const r = o._raw_encode(...ps)
    //         console.log(">> encode after", name, r)
    //         return r
    //     }

    //     o._raw_decode = o._raw_decode || o.decode
    //     o.decode = (...ps)=>{
    //         console.log("<< decode before", name, ...ps)
    //         const r = o._raw_decode(...ps)
    //         console.log("<< decode after", name, r)
    //         return r
    //     }
    // }
    // for(const protokey of Object.getOwnPropertyNames(Wt)) {
    //     const proto = Wt[protokey]
    //     for(const name of Object.getOwnPropertyNames(proto)) {
    //         const o = proto[name]
    //         globalThis.myinject(o, name)    
    //     }
    // }
    
//     var e = this;
//     dt.Token = t,
//     this.m_isLogin = !1,
//     this.m_token = t;
//     for (var a, i = location.search.split(/(?:\?|&)/g), n = 0; n < i.length; n++) {
//         var o = i[n].split("ct=")[1];
//         o && (a = o)
//     }
//     var s = new Wt.aaProto.LoginDataReq;
//     s.Token = t,
//     s.OSType = L.os,
//     s.OSVersion = L.osVersion,
//     s.Browser = String(L.browserType),
//     s.Language = ut.CurrLang,
//     s.BrowserVersion = String(L.browserVersion),
//     s.Width = window.screen.width,
//     s.Height = window.screen.height,
//     s.Ratio = window.devicePixelRatio,
//     s.Machine = lt.GetDeviceName(),
//     s.BrowserTag = lt.GetBrowserTag(),
//     a && Number(a) && (s.Cheat = Number(a));
//     var r = new W;
//     this.m_isProto ? r.reqData = Wt.aaProto.LoginDataReq.encode(s).finish() : r.reqData = s,
//     At.GetInstance().SendCommand(r, (function(t) {
//         e.OnRecvAck(t)
//     }
//     ), !1)
// }

// --------

// file:///Users/a123/Desktop/jilitmp/uat-wbgame.jlfafafa3.com/aa/assets/other/index.18ae2.js

// n.Post = function(t) {
//     var e = this
//       , n = t.content.reqData instanceof Uint8Array ? t.content.reqData : JSON.stringify(t.content.reqData)
//       , o = t.content.URL + "?";
//     if (t.content.info)
//         for (var a in t.content.info)
//             "constructor" != a && (o += a + "=" + t.content.info[a] + "&");
//     t.hasToken && (o += "token=" + this.m_token);
//     var l = Date.now()
//       , m = new XMLHttpRequest;
//     (this.m_isProto || t.isProto) && (m.responseType = "arraybuffer"),
//     m.open("POST", this.m_ip + o, !0),
//     m.setRequestHeader("Content-Type", this.m_isProto || t.isProto ? "application/x-protobuf" : "application/json"),
//     m.timeout = 1e4,
//     m.onload = function() {
//         if (4 == m.readyState) {
//             if (m.status >= 200 && m.status < 300 || 304 == m.status) {
//                 console.log("POST return", e.m_ip + o)
//                 if (e.m_isProto || t.isProto) {
//                     if (m.response)
//                         return t.callback(m.response),
//                         void e.ExecuteNext()
//                 } else if (m.responseText) {
//                     var n = JSON.parse(m.responseText);
//                     return n.token && (e.m_token = n.token),
//                     t.callback(n),
//                     void e.ExecuteNext()
//                 }
//                 t.callback(!1),
//                 e.ExecuteNext()
//             } else
//                 "/account/login" === t.content.URL ? e.m_loginFailCallback && e.m_loginFailCallback() : t.content.Type === u.ErrorHandleType.Handle && s.getInstance().ShowMessageBox(i.StringKey.MSGBOX_CONNECT_FAIL, i.StringKey.MSGBOX_TITLE_SYSTEM_INFO, r.CloseGameTab, r.GetErrorCode(i.ErrorStatus.Game, i.ErrorStatus.Unknow, m.status)),
//                 t.callback(!1),
//                 e.ExecuteNext();
//             c.GetInstance().SetSpeedData(r.minus(Date.now(), l))
//         }
//     }
//     ,
//     m.ontimeout = function() {
//         c.GetInstance().SetSpeedData(r.minus(Date.now(), l)),
//         e.ErrorHandle(t.content.Type)
//     }
//     ,
//     m.onerror = function() {
//         c.GetInstance().SetSpeedData(r.minus(Date.now(), l)),
//         e.ErrorHandle(t.content.Type)
//     }
//     ,
//     console.log("POST", this.m_ip + o),
//     m.send(n)
// }




// function ff(...ps) {
//     console.log("%cencode  {ps}", 'color:red;')
// }

// ff('hello', 123)