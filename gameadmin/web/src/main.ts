import { createApp } from 'vue';
import '@/assets/index.scss'; // 全局样式
import 'element-plus/dist/index.css'; // 按需引入
import '@/assets/style/main.css';

import ElementPlus from 'element-plus';
import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus'; // 按需引入

import i18n from "@/language/i18n";
import { createPinia } from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';

import App from '@/App.vue';
import router from './router/index';
import { Client } from './lib/client';
import { initGolbal } from './lib/util';
import pagination from '@/components/pagination.vue';
import { useStore } from "@/pinia";
import {AdminGameCenter} from "@/api/gamepb/admin";


// 创建 Pinia 实例并使用持久化插件
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
// 初始化应用
let app = createApp(App)
// 初始化全局功能
initGolbal(app)
Client.setHook((resp) => {
  if (resp.error == "用户登录已超时") {
      setTimeout(function () {
          router.replace({path: "/login"})
      }, 1000)

    return true
  }
  return false
})

if (process.env.NODE_ENV === 'production') {
    if (window) {
        window.console.log = function () {};
    }
}

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
// 从 localStorage 读取数据并更新 Pinia store
try {
    const myStore = useStore();
    const storedData = JSON.parse(localStorage.getItem('game_store'));
    if (storedData) {


        myStore.$state = storedData;
    }
} catch (error) {
    console.error("Error loading store data:", error);
    ElMessage.error("Failed to load saved data.");
}



// 使用插件和组件
app.use(pinia)
    .use(i18n)
    .use(router)
    .use(ElementPlus)
    .component('Pagination', pagination)
    .mount('#app');
