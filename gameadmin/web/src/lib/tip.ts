import { ElLoading, ElMessage, ElMessageBox, } from 'element-plus'

export class tip {
	static i(text: string, title?: string) {
		ElMessage({ message: text, type: "info", })
	}

	static e(text: string, title?: string) {
		ElMessage({ message: text, type: "error" })
	}

    static w(text: string, title?: string) {
        ElMessage({ message: text, type: "warning" })
    }

	static s(text: string, t: number = 3) {
		ElMessage({ message: text, type: "success", duration: t * 1000 })
	}

	static async ask(text: string): Promise<"ok" | "cancel"> {

		return new Promise((resolve, reject) => {
			ElMessageBox.confirm(text).then(() => {
				resolve("ok")
			}).catch(() => {
				resolve("cancel")
			})
		})
	}

	static async input(title: string, value: any, fn?: () => boolean): Promise<[data: any, isCancel: boolean]> {
		return new Promise((resolve, reject) => {
			ElMessageBox.prompt("", title, {
				confirmButtonText: '确定',
				cancelButtonText: '取消',
				inputValidator: fn,
				inputValue: value,
				// inputErrorMessage: 'Invalid Email',
			})
				.then(({ value }) => {
					resolve([value, true])
				})
				.catch(() => {
					resolve([null, false])
				})
		})
	}

	static loading(text?) {
		return ElLoading.service({
			lock: true,
			text: text ? text : 'Loading',
			background: 'rgba(0, 0, 0, 0.7)',
		})
	}
}
