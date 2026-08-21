- fuck debugger
window._Function = window.Function; window.Function = function (...args) {
    if (args[0].startsWith("debugger")) { return console.log; }; return window._Function(...args); 
}


debug throw Error(Cw);