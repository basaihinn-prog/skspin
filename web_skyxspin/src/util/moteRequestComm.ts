import axios from "axios";
import type { AxiosInstance, AxiosResponse } from "axios";

export interface ResponseModel<T = any> {
    data: T;
}

/**
 * Creates an Axios instance with a dynamic base URL.
 * @param baseURL - The base URL for API requests.
 * @returns AxiosInstance
 */
const createAxiosInstance = (baseURL: string): AxiosInstance => {
    const instance = axios.create({ baseURL });

    instance.interceptors.request.use((requestConfig) => {
        const app_secret = localStorage.getItem("app_secret");
        const app_id = localStorage.getItem("app_id");
        const systemLanguage = localStorage.getItem("systemLanguage");

        requestConfig.headers.set("AppID", app_id);
        requestConfig.headers.set("AppSecret", app_secret);
        requestConfig.headers.set("Lang", systemLanguage);

        return requestConfig;
    });

    instance.interceptors.response.use((response: AxiosResponse<ResponseModel>): AxiosResponse['data'] => {
        return response.data;
    });

    return instance;
};

// Base URLs for different APIs
const BASE_URL_1 = window.location.protocol + "//" + window.location.host + "/api/"
const BASE_URL_2 = window.location.protocol + "//" + window.location.host.replace("game", "external") + "/api/"; 

// Create Axios instances
export const moteAxiosComm = createAxiosInstance(BASE_URL_1);
export const externalMoteAxiosComm = createAxiosInstance(BASE_URL_2);
