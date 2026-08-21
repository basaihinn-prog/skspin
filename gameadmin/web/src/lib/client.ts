import { Callback, Hook, HttpConn,  } from "./conn";
import {useI18n} from "vue-i18n";


export class Client {
	private static http: HttpConn = new HttpConn('/api')

	// 另一种写法
	// let [data, err] = Client.Do(AuthApi.Login, {})
	public static async Do<TReq, TResp> (
		fn: (client, req: TReq) => Promise<[TResp, any]>,
		req: TReq
	): Promise<[TResp, any]> {
		return fn(this, req)
	}

	// 回调式写法:
	// Client.send("server/method", {}, (data, err) => console.log(data, err))
	public static send(method: string, args: any, cb: Callback): void;
	// await 写法：
	// let [data, err] = await Client.send("server/method", {})
	public static send(path: string, args: any): Promise<[any, any]>;
	public static send(path, args, cb?) {
		if (cb) {
			this.sendCB(path, args, cb)
		} else {
			return new Promise((resolve, reject) => {
				this.sendCB(path, args, (msg, err) => resolve([msg, err]))
			})
		}
	}

    private static sendCB(method: string, args: any, cb?: Callback) {
		this.http.send(method, args, cb)
	}

	// 设置登录信息
	public static setToken (token: string | null) {
		this.http.setToken(token)
	}

	// 设置调用语言
	public static setLanguage (language: string) {
		this.http.setLanguage(language)
	}

	public static setHook (hook: Hook) {
		this.http.hook = hook
	}
}
import {tip} from '@/lib/tip';
/**
 * 处理并验证数据，确保数据的完整性和有效性。
 * @param {Array|Object} data - 需要验证的数据。
 * @param {string} name - 需要检查的属性名。
 * @param {boolean} [showError=true] - 是否显示错误信息。
 * @param {Function} [errorHandler] - 自定义错误处理函数。
 * @returns {Array|Object|null} 处理后的数据或 null（如果数据无效）。
 */
export default function ReturnErrorHandle(data,name,showError = true, errorHandler = null) {
    const defaultErrorHandler = (message) => {
        if (showError) {
            console.error(message); // 使用 console.error 替代 tip.e
        }
        return null;
    };
    const handleError = errorHandler || defaultErrorHandler;
    // 检查 data 是否为 null
    if (data === null) {
        return handleError('Missing data');
    }
    const validateItem = (item) => {
        if (item[name] === null || item[name] === undefined || item[name] === '') {
            item[name] = 'Wrong data';
        }
    };
    // 检查 data 是否为数组
    if (Array.isArray(data)) {
        // 检查数组是否为空
        if (data.length === 0) {
            return handleError('Data is empty');
        }
        // 遍历数组检查是否包含空字段
        // 遍历数组检查是否包含空字段
        data.forEach(validateItem);
    }

    // 如果 data 不是数组，也不是 null，则直接返回 data
    return data;
}
