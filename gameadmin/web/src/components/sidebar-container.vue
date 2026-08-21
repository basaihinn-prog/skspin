<template>

    <div class="menu_container" >


        <div class="login-container">
            <el-image :src="logo" style="width: 60px;display: flex;align-items: center;"
                      :style="{marginLeft: store.sidebarShow ? 0 : '20px'}" fit="fill"/>
            <span v-if="!store.sidebarShow" style="width: 140px;padding-left: 10px">{{ $t(title) }}</span>
        </div>

        <el-menu
            class="menu_list"
            :default-active="defaultActive"
            :collapse="store.sidebarShow"
            :collapse-transition="true"
            :router="true"
            :unique-opened="true"
        >

            <template v-for="(menu, index) in AdminInfo.MenuList">
                <el-sub-menu v-if="menu.Children" :index="menu.Url" :show-timeout="800">
                    <template #title>
                        <el-icon>
                            <component :is="menu.Icon"/>
                        </el-icon>
                        <span>{{ $t(menu.Title.replace(/^\s+|\s+$/g, "")) }}</span>
                    </template>


                    <el-menu-item-group>
                        <template v-for="(menu1,menu1Index) in menu.Children" :key="menu1.Url">


                            <el-sub-menu class="menu2" v-if="menu1.Children" :index="menu1.Url">
                                <template #title>
                                    <el-icon>
                                        <component :is="menu1.Icon"/>
                                    </el-icon>
                                    <span>{{ $t(menu1.Title.replace(/^\s+|\s+$/g, "")) }}</span>
                                </template>


                                <el-menu-item-group>

                                    <template v-for="(menu2,menu2Index) in menu1.Children" :key="menu1.Url">
                                        <el-menu-item @click="tabRouterClick(menu2)" :index="menu2.Url">
                                            <template #title>
                                                <el-icon>
                                                    <component :is="menu2.Icon"/>
                                                </el-icon>
                                                <span>{{ $t(menu2.Title.replace(/^\s+|\s+$/g, "")) }}</span>
                                            </template>
                                        </el-menu-item>
                                    </template>
                                </el-menu-item-group>

                            </el-sub-menu>


                            <el-menu-item v-else @click="tabRouterClick(menu1)" :index="menu1.Url">
                                <template #title>
                                    <el-icon>
                                        <component :is="menu1.Icon"/>
                                    </el-icon>
                                    <span>{{ $t(menu1.Title.replace(/^\s+|\s+$/g, ""))}}</span>
                                </template>
                            </el-menu-item>


                        </template>
                    </el-menu-item-group>


                </el-sub-menu>


                <el-menu-item v-else @click="tabRouterClick(menu)" :index="menu.Url">
                    <el-icon>
                        <component :is="menu.Icon"/>
                    </el-icon>
                    <template #title>

                        <span>{{ $t(menu.Title.replace(/^\s+|\s+$/g, "")) }}</span>
                    </template>
                </el-menu-item>


            </template>
        </el-menu>

    </div>
</template>

<script setup lang='ts'>
import logo from "@/assets/style/logo.png"
import {nextTick, onMounted, ref, reactive, watch, computed, toRaw, watchEffect} from 'vue';

import variables from "@/assets/style/variables.module.scss";
import {useRouter} from 'vue-router';
import {useStore} from '@/pinia/index';

const store = useStore()
let router = useRouter();
import {useI18n} from 'vue-i18n';
import {storeToRefs} from "pinia";
import {Client} from "@/lib/client";
import {AdminGroup} from "@/api/adminpb/group";
import {tip} from "@/lib/tip";
import ut from "@/lib/util";
import {debug} from "util";

const {t} = useI18n();
const piniaStore = useStore()
const {setMenuList} = piniaStore
const shouldShowMenuModel = computed(() => store.viewModel === 2 && store.sidebarShow);
const sidebarLeftPosition = computed(() => store.viewModel === 1 ? '0' : '-210px');
const {AdminInfo, elMenuActive, tabMenuList} = storeToRefs(store)
const {setActiveTabs, addTab, setElMenuActive} = piniaStore
const defaultActive = computed(() => elMenuActive.value.value || '/dashboard')
const menuRouter = ref(true)

const title = ref("")


const hideSidebar = () => {
    store.setSidebarShow(false);
};

let tabRouterClick = (menu: any) => {
    // 检查是否已经存在相同的标签页
    let nTitle = t(menu.Title).replace(/^\s+|\s+$/g, "")
    menu.Title = nTitle
    setActiveTabs(nTitle)
    setElMenuActive({
        title: nTitle,
        value: menu.Url
    })
    addTab(menu);

    router.push(menu.Url)
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

onMounted(async () => {
    if (!AdminInfo.value.MenuList) {
        await store.initAdminInfo();
    }
    setMenuList(internationalizeData(AdminInfo.value.MenuList))

    if (AdminInfo.value.GroupId <= 1) {
        title.value = "总控"
    } else if (AdminInfo.value.GroupId == 2) {
        title.value = "线路商"

    } else {
        title.value = "商户"
    }
});



</script>
<style scoped lang='scss'>
.login-container {
    width: 100%;
    height: auto;
    display: flex;
    align-items: center;
    margin: 10px auto;
    color: #ffffff;
    font-weight: bolder;
}

.mobileMenu {
  background: #001529;
}

.menu_container {
  background: #001529;
}

.el-menu-vertical-demo {

  height: calc(100dvh - 50px);
  overflow: auto;
}
</style>

<style>
.el-menu{
    background: #001529;
}
.el-sub-menu__title, .el-menu-item{
    color: rgba(255,255,255,0.8);
    font-size: 14px;
}
.el-sub-menu__title{
    border-radius: 8px;
}
.el-sub-menu__title:hover{
    background: rgba(255,255,255,0.2);
}
.el-menu-item-group .el-menu-item,.el-menu .el-menu-item{

    width: 96%;
    margin: 5px auto;
    border-radius: 8px;
}
.el-menu-item-group .el-menu-item:hover,.el-menu .el-menu-item:hover{

    background: rgba(255,255,255,0.2);
    color: rgba(255,255,255, 0.8);

}

.menu_list{
    height: calc(100dvh - 40px);
    overflow-y: auto;
    font-size: 14px;
}



.el-menu-item.is-active{
    background: var(--el-color-primary)!important;
    color: rgba(255,255,255, 0.8)!important;
}
</style>
