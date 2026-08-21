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

// function ff(...ps) {
//     console.log("%cencode  {ps}", 'color:red;')
// }

// ff('hello', 123)

//-------------------------

globalThis.myinject = function (o, name) {
    // const o = proto[name]
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
for(const protokey of Object.getOwnPropertyNames(Tt)) {
    const proto = Tt[protokey]
    for(const name of Object.getOwnPropertyNames(proto)) {
        const o = proto[name]
        globalThis.myinject(o, name)    
    }
}      