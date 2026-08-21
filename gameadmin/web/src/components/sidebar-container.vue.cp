<template>
    <div
        v-if="shouldShowMenuModel"
        class="menuModel"
        @click="hideSidebar">
    </div>

    <transition name="transitionName">
        <div class="menu_container"
             :class="{ 'show-sidebar': store.sidebarShow }"
             v-if="(shouldShowMenuModel && (store.viewModel === 2)) || store.viewModel === 1"
        >
            <div class="logo-container">

                    <div class="wh-full flex flex-center" :key="+ sidebarLeftPosition != 0 ">
                        <div class="logo-context">
                            <el-image fit="fill" :src="logo" class="logo-image" />

                        </div>

                        <span class="operatorName">{{ $t(title) }}</span>

                    </div>

            </div>

            <el-menu
                ref="menuRef"
                :unique-opened="true"
                :router="menuRouter"
                :default-active="defaultActive"
                :collapse="false"
                :collapse-transition="false"
                :background-color="variables['menu-background']"
                :text-color="variables['menu-text']"
                :active-text-color="variables['menu-active-text']"

        >

                <template v-for="menu in AdminInfo.MenuList" :key="menu.Url">
                    <el-sub-menu v-if="menu.Children" :index="menu.Url" >
                        <template #title>
                            <el-icon>
                                <component :is="menu.Icon" />
                            </el-icon>
                            <span class="menu_text font_size16"> {{  $t(menu.Title).replace(/^\s+|\s+$/g, "") }}</span>
                        </template>
                        <el-menu-item-group>
                            <template v-for="(menu1,menu1Index) in menu.Children" :key="menu1.Url">
                                <el-sub-menu class="menu2" v-if="menu1.Children" :index="menu1.Url">
                                    <template #title>
                                        <el-icon>
                                            <component :is="menu1.Icon" />
                                        </el-icon>
                                        <span style="color:#7c8a96;" class="font_size16">{{ $t(menu1.Title.replace(/^\s+|\s+$/g, "")) }}</span></template>
                                    <el-menu-item v-for="menu2 in menu1.Children" :index="menu2.Url" :key="menu2.Url" @click="tabRouterClick(menu2)">
                                        <template #title>
                                            <el-icon>
                                                <component :is="menu2.Icon" />
                                            </el-icon>
                                            <span class="font_size16">{{ $t(menu2.Title).replace(/^\s+|\s+$/g, "") }}</span>
                                            <router-link :to="menu2.Url" ></router-link>
                                        </template>
                                    </el-menu-item>
                                </el-sub-menu>
                                <el-menu-item v-else :index="menu1.Url" class="menu_text"  @click="tabRouterClick(menu1)">
                                    <el-icon>
                                        <component :is="menu1.Icon" />
                                    </el-icon>
                                    <span >{{ $t(menu1.Title).replace(/^\s+|\s+$/g, "") }}</span>
                                    <template #title>
                                        <router-link :to="menu1.Url" ></router-link>
                                    </template>
                                </el-menu-item>
                            </template>
                        </el-menu-item-group>
                    </el-sub-menu>
                    <el-menu-item v-else :index="menu.Url"  @click="tabRouterClick(menu)">
                        <el-icon>
                            <component :is="menu.Icon" />
                        </el-icon>
                        <router-link :to="menu.Url" class="menu_text"><span class="font_size16">{{ $t(menu.Title).replace(/^\s+|\s+$/g, "") }}</span></router-link>
                    </el-menu-item>
                </template>
            </el-menu>

        </div>
    </transition>
</template>

<script setup lang='ts'>
import logo from "@/assets/style/logo.png"
import {nextTick, onMounted, ref, reactive, watch, computed, toRaw, watchEffect} from 'vue';

import variables from "@/assets/style/variables.module.scss";
import { useRouter } from 'vue-router';
import {useStore} from '@/pinia/index';

const store = useStore()
let router = useRouter();
import { useI18n } from 'vue-i18n';
import {storeToRefs} from "pinia";
import {Client} from "@/lib/client";
import {AdminGroup} from "@/api/adminpb/group";
import {tip} from "@/lib/tip";
import ut from "@/lib/util";
const { t } = useI18n();
const piniaStore = useStore()
const { setMenuList } = piniaStore
const shouldShowMenuModel = computed(() => store.viewModel === 2 && store.sidebarShow);
const sidebarLeftPosition = computed(() => store.viewModel === 1 ? '0' : '-210px');
const { AdminInfo,elMenuActive, tabMenuList } = storeToRefs(store)
const { setActiveTabs,addTab,setElMenuActive } = piniaStore
const defaultActive = computed(() => elMenuActive.value.value || '/')
const menuRouter = ref(true)

const title = ref("")
    if (AdminInfo.value.GroupId <= 1){
        title.value = t("总控")
    }
    else if (AdminInfo.value.GroupId == 2){
        title.value = t("线路商")
    }else{

        title.value = t("商户")
    }



const hideSidebar = () => {
    store.setSidebarShow(false);
};
let tabRouterClick = (menu:any) =>{

    // 检查是否已经存在相同的标签页
    let nTitle = t(menu.Title).replace(/^\s+|\s+$/g, "")
    menu.Title = nTitle
    setActiveTabs(nTitle)
    setElMenuActive({
        title:nTitle,
        value:menu.Url
    })
    addTab(menu);
}
function internationalizeData(data) {
    // 假设您已经有了语言翻译的资源，例如一个翻译函数 translate(key, language)
    // 递归遍历数据并翻译标题
    for (const item of data) {
        item.Title = t(item.Title).replace(/^\s+|\s+$/g, "");
        if (item.Children) {
            internationalizeData(item.Children);
        }
    }
    return data;
}
onMounted( async () => {
    if (!AdminInfo.value.MenuList) {
        await store.initAdminInfo();
    }
    setMenuList(internationalizeData(AdminInfo.value.MenuList))
});

</script>
<style scoped lang='scss'>
.logo-container {
    width: 100%;
    height: $navbar-height;
    background-color: $sidebar-logo-background;
    display: flex;
    align-items: center;
    justify-content: center;
    --left:auto
}

.logo-context{
    width: 30%;
    height: $navbar-height;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-left: 20px;
}
.logo-title {
    flex-shrink: 0; /* 防止容器在空间不足时缩小 */
    margin-left: 10px;
    font-size: 14px;
    font-weight: bold;
    color: white;
}
.logo-image{
    width: 100%;
    height: auto;
    margin: 0 auto;
}
.menu_container{
    width: 240px;
    height: 100vh;
    position: fixed;
    left: var(--left);
    overflow: hidden;
    z-index: 4;
}
.el-menu:not(.el-menu--collapse) {
    width: 240px;
}
.el-menu {
    height: calc(100vh - $navbar-height);
    position: fixed;
    left: var(--left);
    overflow-y: auto;
    overflow-x: hidden;
    z-index: 4;
    font-size: var(--size16)!important;
}

.show-sidebar {
    left: 0;
}
.el-sub-menu .el-icon, .el-icon:not(.is-active){
    font-size: var(--size16);
    color: #7c8a96;
}

.el-menu-item:not(.is-active){
    color: #7c8a96;
}

.el-sub-menu .menu_text{
    font-size: $size16;
}
.el-sub-menu .menu_text:not(.is-active),
.menu_text:not(.is-active){
    font-weight: 540;
    color: #7c8a96;

}
.menuModel {
    width: 100vw;
    height: 100vh;
    position: fixed;
    left: 0;
    background: rgba($color: #000000, $alpha: 0.5);
    z-index: 4;
}




.transitionName-enter-from,
.transitionName-leave-to {
    transform: translateX(-20px);
    opacity: 0;
}

.transitionName-enter-to,
.transitionName-leave-from {
    opacity: 1;
}

.transitionName-enter-active {
    transition: all 0.7s ease;
}

.transitionName-leave-active {
    transition: all 0.3s cubic-bezier(1, 0.6, 0.6, 1);
}
.menu2 .el-sub-menu__title{
    color: #7c8a96;
}
.flex-center{
    align-items: center;
}
.operatorName{
    color: #7c8a96;
    margin-left: 20px;
    font-size: 20px;
    font-weight: bolder;
}
</style>
<style>
.el-sub-menu__title span{
    width: 80%;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    display: block;
}
</style>
