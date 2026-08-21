import axios from "axios";
import type {AxiosInstance, AxiosResponse} from "axios"
export const moteAxiosComm:AxiosInstance = axios.create({
    baseURL: window.location.protocol + "//" + window.location.host + "/api/"
    // baseURL: "http://localhost:11100/"
})
export interface ResponseModel<T = any> {
    data: T;
}

moteAxiosComm.interceptors.request.use((requestConfig) => {

    // const AppSecret = localStorage.getItem("app_secret")
    // const AppID = localStorage.getItem("app_id")
    //
    // requestConfig.headers.set("AppID", AppID)
    // requestConfig.headers.set("AppSecret", AppSecret)


    const app_secret = localStorage.getItem("app_secret")
    const app_id = localStorage.getItem("app_id")
    const systemLanguage = localStorage.getItem("systemLanguage")



    requestConfig.headers.set("AppID", app_id)
    requestConfig.headers.set("AppSecret", app_secret)
    requestConfig.headers.set("Lang", systemLanguage)


    return requestConfig
})


moteAxiosComm.interceptors.response.use((response: AxiosResponse<ResponseModel>): AxiosResponse['data'] => {

    return response.data
})
