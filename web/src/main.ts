import './assets/main.css'

import { createApp } from 'vue'
// 如果您正在使用CDN引入，请删除下面一行。
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import pinia from "@/stores"
import ElementPlus, {ElNotification} from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/display.css'
import {moteAxiosComm} from "@/util/moteRequestComm";

import "@/assets/main.css"
import {SetOperator} from "@/util/util";

const app = createApp(App)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

const localLang = localStorage.getItem("systemLanguage")
if (!localLang){
    localStorage.setItem("systemLanguage", "en")
}



SetOperator()

app.use(ElementPlus)

app.use(pinia)
app.use(router)

app.mount('#app')
