
let throttleTimer: number | null
export const Throttle = (fn: Function, delay: number) :Function => {


    return (...args: unknown[]) => {

        if (throttleTimer) {
            return;
        }
        throttleTimer = setTimeout(() => {
            fn.apply(this, args);
            throttleTimer = null;
        }, delay);
    }
}



export const SetOperator = () => {

    let currentAppSecret  =  localStorage.getItem("app_secret")
    let currentAppId  = localStorage.getItem("app_id")

    let param = window.location.search.split("?")[1]
    if (!param){

        param = `app_secret=11c5d190-3add-4482-b6d4-ee990903f981&app_id=faketrans`
    }


    if (param){
        let params = param.split("&")


        let paramRecord:Record<string, string> = {}
        for (const i in params) {
            let paramItem = params[i].split("=")
            paramRecord[paramItem[0]] = paramItem[1]
        }

        if (currentAppSecret != paramRecord["app_secret"] || currentAppId != paramRecord["app_id"]) {

            localStorage.removeItem("userId")
        }


        localStorage.setItem("app_secret", paramRecord["app_secret"])
        localStorage.setItem("app_id", paramRecord["app_id"])

    }

    return param

}
