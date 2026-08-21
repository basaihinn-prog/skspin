package staticproxy

import "bytes"

func hookhtml(content []byte) []byte {
	content = bytes.Replace(content, []byte("</script>"), []byte(`</script><script>window.gtag=console.log;window.dataLayer=[];</script>`), 1)
	content = bytes.Replace(content, []byte("<div class=\"loader-circle\"></div>"), []byte(""), 3)
	return content
}

func hookfetch(content []byte) []byte {
	s := `<body>
<script>
(function(){
    function geturl(u) {
        const url = new URL(u);
        if (url.pathname.startsWith('/shared/') && (url.pathname.endsWith('/index.json') || url.pathname.endsWith('/index.js'))) {
            const gid = location.pathname.split('/')[1];
            url.pathname = url.pathname.replace('/index.', '/index-'+gid+'.')
            u = url.toString();
        }
        
        return u;
    }

    const open = XMLHttpRequest.prototype.open;
    XMLHttpRequest.prototype.open = function(method, url, ...args) {
        open.call(this, method, geturl(url),  ...args);
    };

    const fetch = window.fetch;
    window.fetch = function(resource, options) {
        return fetch(geturl(resource), options);
    };
})();
</script>
	`

	return bytes.Replace(content, []byte("<body>"), []byte(s), 1)
}

func hookSkipInsecurePmt(content []byte) []byte {
	oldS := "<head>"
	newS := oldS + `
<script>
(function(){
let intervalID = 0
const t0 = new Date();
const cb = ()=>{
    const node = document.querySelector('.npve_continue_button_port')||document.querySelector('.npve_continue_button_land');
	if(node){
        node.click(); 
		console.log("clearInterval", intervalID);
		clearInterval(intervalID);
        return;
	}

    const now = new Date();
    if (now - t0 > 30000) {
		console.log("clearInterval", intervalID);
		clearInterval(intervalID);
    }
}
intervalID = window.setInterval(cb, 100);
console.log("intervalID", intervalID);
})();
</script>`

	return bytes.Replace(content, []byte(oldS), []byte(newS), 1)
}

// {scope: '/'}
func hookNopSW(content []byte) []byte {
	oldS := "<head>"
	newS := oldS + `
<script>
(function(){
const register = navigator.serviceWorker.register;
navigator.serviceWorker.register = function(scriptURL, options) {
console.log(scriptURL, options);
//options.scope = '/nonononononofucksw';
//scriptURL = "/nonono/404.js"
// debugger;
// return register.call(this, scriptURL, options);
return new Promise(resolve => setTimeout(resolve, 100))
}
})();
</script>`

	return bytes.Replace(content, []byte(oldS), []byte(newS), 1)
}
