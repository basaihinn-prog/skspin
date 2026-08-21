
(function(){
    // const _Function = window.Function;
    // window.Function = function (...args) {
    //     if (args[0].startsWith("debugger")) { return console.log; }; return _Function(...args); 
    // }
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
        // console.log("XMLHttpRequest", method, url,  ...args);
        open.call(this, method, geturl(url),  ...args);
        // if (url.startsWith('/shared/')) {
        //     this.setRequestHeader('rp-href', location.href);
        // }
    };

    const fetch = window.fetch;
    window.fetch = function(resource, options) {
        return fetch(geturl(resource), options);
    };
})();