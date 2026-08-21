import axios from "axios";
import type {AxiosInstance, AxiosResponse} from "axios"


export const axiosComm:AxiosInstance = axios.create({
    baseURL: window.location.protocol + "//" + window.location.host
})
export interface ResponseModel<T = any> {
    data: T;
}

axiosComm.interceptors.request.use((requestConfig) => {
    requestConfig.url =  requestConfig.url + `?nowTime=${new Date().getTime()}`

    return requestConfig
})


axiosComm.interceptors.response.use((response: AxiosResponse<ResponseModel>): AxiosResponse['data'] => {

    return response
})
