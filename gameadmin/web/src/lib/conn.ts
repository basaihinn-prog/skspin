import ut from "./util";
import router from "@/router";
import {tip} from "@/lib/tip";
import Handsontable from "handsontable";
import languages = Handsontable.languages;
import {ElMessage} from "element-plus";


// 登录超时多语言 MAP
let loginTImeOut = {
    "zh": "登录超时.",
    "en": "Login timeout.",
    "it": "Per effettuare il log on in eccesso.",
    "es": "El login ha expirado. Por favor.",
    "th": "หมดเวลาเข้าสู่ระบบแล้วโปรดรีเฟรชส่วนติดต่อ",
    "idr": "Log masuk sudah habis.",
}

var timeOut = false


export interface Resp {
    id: number | string
    data: any
    error?: string
}
export type Callback = (data: any, err?:string) => void
// return true 表示中断后面的逻辑
export type Hook = (resp : Resp) => boolean

class Conn {
	protected _url :string  = ""
	protected _token :string  = ""
	protected _language:string  = ""

    hook : Hook = null

    constructor(url : string) {
        this._url = url
    }

    public send(method: string, data: any, cb?: (msg: any, err: string) => void) {

    }

    public setLanguage(language: string) {
        this._language = language
    }

    public setToken(token : string | null) {
        this._token = token
    }

    public isOk() : boolean {
        return true
    }
}



export class HttpConn extends Conn {

    public send(method: string, data: any | FormData, cb?: (msg: any, err: string) => void) {
        let xhr = new XMLHttpRequest()
        const languagesLocal = localStorage.getItem("language")



		xhr.onload = () => {

            if (xhr.status == 200) {
                let resp = JSON.parse(xhr.responseText)

                if (resp.error == "用户登录已超时" && !timeOut){
                    timeOut = true

                    ElMessage({ message: loginTImeOut[languagesLocal || 'en'], type: "error" , duration: 1000})
                    setTimeout(function () {
                        router.replace("/login")
                    }, 1000)
                    return;
                }

                if (this.hook && this.hook(resp)) {
                   return
                }

                if (resp.error) {


                    cb(null, resp.error)
                } else {
                    // console.log("Recv", resp)
                    cb(resp.data, null)
                }
            } else if (xhr.status != 0) {
                router.replace("/login")
                cb(null, xhr.responseText ? xhr.responseText : ("status:"+xhr.status))
            }
		}

		xhr.timeout = 60000;
		xhr.ontimeout = function (e) {
            console.log("on timeout",e)
            cb(null, xhr.statusText ? xhr.statusText : "timeout")
		}

		xhr.onerror = function (e) {
            console.log("on error", e)
            cb(null, xhr.statusText ? xhr.statusText :"error")
		}

        let url = ut.UrlJoin(this._url,  method)
		xhr.open("POST", url, true)
        xhr.setRequestHeader("Lang", this._language || JSON.parse(localStorage.getItem('game_store'))?.language || 'zh')

        if (this._token) {
            xhr.setRequestHeader('Authorization', this._token);
        }

        const isLoginApi = url.split("/").includes("AdminAuth") && url.split("/").includes("Login")
        if(isLoginApi){
            timeOut = false
        }
        if(!isLoginApi && timeOut){
            return
        }




        if (data instanceof FormData) {
            xhr.send(data)
        } else {
            let jsondata = JSON.stringify(data)
            xhr.setRequestHeader("Content-Type", "application/json")
            xhr.send(jsondata)
        }
	}

}
