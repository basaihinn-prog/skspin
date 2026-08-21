- fuck debugger
nodbg=1
window._Function = window.Function; window.Function = function (...args) {
    console.log("Function", ...args);
    if (args[0].startsWith("debugger")) { return console.log; }; return window._Function(...args); 
}

window.location = new Proxy(location_dta, {
    get(target, p) {
        let v = target[p];
        console.log('get', p, '==', v);

        var s = new Error().stack;
        if (!s.includes('google')) {
            debugger;
        }
        return v;
    }
})


new Proxy(G9, {
    get(target, p) {
        let v = target[p];
        console.log('get', p, '==', v);
        debugger;
        return v;
    }
})

C = new Proxy(C, {
    get(target, p) {
        let v = target[p];

        if (p == "isSecureContext") {
            console.log('get', p, '==', v);
            debugger;
        }
        return v;
    }
})


cc.log = function(...args) {
    console.log("cc.log", ...args)
}

var b = document.querySelector('.npve_continue_button_port')
b.click()


ssh -D 1080  -Nf   root@rp-hk-dev
zy@myubuntu:~/Desktop/chrom-mybuild$ ./chrome --proxy-server="socks5://127.0.0.1:1080"

debug throw Error(Cw);

var spinbtn = cc.director.getRunningScene().getChildByName('Canvas').getChildByName('spin_button_holder').getChildByName('spin_button_controller')

document.querySelector('.start-button-inner').click()


更新http 步骤

- 更新pgggateway

- Caddyfile
  50002 添加 http://m-pghttp.rpgamestest.com 

- game_config.yaml 修改 **or**, 添加 **api_origin**
    pglaunchurl: "http://**m-pghttp**.rpgamestest.com/39/index.html?ot=abcd1234abcd123432531532111kkafa&btt=1&l=en&ops=000102030404e44708090b853992e98c&or=https%3A%2F%2Fstatic-pg.rpgamestest.com&api_origin=https%3A%2F%2Fapi.rpgamestest.com"

- pggateway_config.yaml 添加 **m-pghttp.**
  reverseproxy:
    "static-pg.": "https://static.pg-demo.com"
    "m-pg.": "https://m.pgsoft-games.com"
    "m-pghttp.": "https://m.pgsoft-games.com"  


----
- 更新pgggateway

- Caddyfile
    禁用自动升级https auto_https disable_redirects  
    50002 添加 http://m-pg.rpgamestest.com 

- game_config.yaml 修改 **or**, 添加 **api_origin**
    pglaunchurl: "http://m-pg.rpgamestest.com/39/index.html?ot=abcd1234abcd123432531532111kkafa&btt=1&l=en&ops=000102030404e44708090b853992e98c&or=https%3A%2F%2Fstatic-pg.rpgamestest.com&api_origin=https%3A%2F%2Fapi.rpgamestest.com"



C[ep(0x633)+ep(0xdc3)+eC(0x1b6a)+ep(0xaa1)+eC(0xa28)]


This request has been blocked; the content must be served over HTTPS.
