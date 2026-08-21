<template>
    <div class="container">
        <SidebarContainer></SidebarContainer>
        <div class="rightView">
            <TopContainer></TopContainer>
            <div class="page_history_list" id="page_history_list">
                <div @click="router.push(tabMenuList[0].Url)"
                     class="page_history_item"
                     :class="router.currentRoute.value.path == tabMenuList[0].Url ? 'page_history_item_active' : ''"
                     :style="{zIndex: router.currentRoute.value.path == tabMenuList[0].Url ? 1 : 99}"
                >


                    <el-image :src="HouseIcon" style="width:20px;height: 20px;margin-right: 10px"></el-image>

                    {{ $t('游戏报表') }}



                    <div class="el-menu-icon ele-tab-corner-left"></div>
                    <div class="el-menu-icon ele-tab-corner-right"></div>

                </div>


                <template v-for="(tag, index) in tabMenuList" :key="index">
                    <div @click="router.push(tag.Url)"
                         v-if="index != 0"
                         class="page_history_item"
                         :id="tag.Url"
                         :class="router.currentRoute.value.path == tag.Url ? 'page_history_item_active' : ''"
                         :style="{zIndex: router.currentRoute.value.path == tag.Url ? 1 : 999 - index}"
                    >

                        <el-icon style="height: auto;vertical-align: middle;margin-top: -3px;margin-right: 3px">
                            <component :is="tag.Icon"></component>
                        </el-icon>

                        {{ $t(tag.Title) }}
                        <el-icon v-if="index != 0" style="vertical-align: -0.15em;display: inline-block;font-size: 12px;margin-left: 8px"
                                 @click.stop="handleTabRemove(tag.Title)">
                            <Close/>
                        </el-icon>

                        <div class="el-menu-icon ele-tab-corner-left"></div>
                        <div class="el-menu-icon ele-tab-corner-right"></div>

                    </div>
                </template>
            </div>


            <div class="router-view">
                <router-view v-slot="{ Component,route }" :key="router.currentRoute.value.name" ref="routerRef" v-if="isRouterAlive" class="page_content">

                        <transition name="transitionName">
                                <component :is="Component" :key="router.currentRoute.value.name" />

                        </transition>

                </router-view>
            </div>
        </div>
    </div>
</template>

<script setup lang='ts'>
import {computed, getCurrentInstance, nextTick, onMounted, ref, watch, watchEffect} from 'vue'
import SidebarContainer from '@/components/sidebar-container.vue'
import TopContainer from '@/components/top-container.vue'
import router from '@/router'
import {useRouter,useRoute} from 'vue-router';
import {GetAdminInfo, useStore} from '@/pinia';
import {useI18n} from 'vue-i18n';

import HouseIcon from "@/assets/login/house.png"

const store = useStore()
const routers = useRouter();
const route = useRoute()
const { t,locale } = useI18n();
const { proxy } = getCurrentInstance();
const { tabMenuList, viewModel,language,activeTabs } = storeToRefs(store)
const piniaStore = useStore()
const { initAdminInfo,setActiveTabs,setElMenuActive } = piniaStore
import { storeToRefs } from 'pinia'
import {House} from "@element-plus/icons-vue";
import {ElMessageBox} from "element-plus";
import {RTPConstructor} from "@/lib/RTP_config";
import ut from "@/lib/util";
import {Client} from "@/lib/client";
let defaultTab = { ID: 0, Icon: "Home", Pid: 0, Sort: 0, Title: '游戏报表', Url: "/dashboard"};
const translateTitle = (title: string) => t(title.trim())
const isTabClosable = (tab) => tab.Title !== t('游戏报表')
const routerRef = ref(null)
const isRouterAlive = ref(true)
let computedActiveTabs = computed(() =>  activeTabs.value || '/')




watch(()=>route.path, (newData)=>{
    // let screenNum = document.getElementById(newData).offsetLeft / window.innerWidth
    // if(document.getElementById(newData).offsetLeft >  window.innerWidth){
    //
    // }

    // document.getElementById("page_history_list").scrollTo({
    //     left: (window.innerWidth * screenNum) + window.innerWidth
    // })
})





onMounted(async () => {
    let getStorgelang = localStorage.getItem('lang')
    let getStorgeactiveTabs = localStorage.getItem('activeTabs')
    if (getStorgelang || getStorgeactiveTabs) {
        localStorage.clear()
        location.reload()
    }
    document.title = t('游戏管理后台')
    await initAdminInfo()
    if (route.path !== '/dashboard') {
        // tabMenuList.value = [defaultTab]
        setElMenuActive(
            {
                title:computedActiveTabs.value,
                value:route.path
            }
        )
    } else {
        setActiveTabs(t('游戏报表'))
        setElMenuActive(
            {
                title:'',
                value:''
            }
        )
    }
    tabMenuList.value = tabMenuList.value.map(tab => ({
        ...tab,
        Title: translateTitle(tab.Title)// 确保所有标题经过国际化处理
    }));
    locale.value = language.value || 'zh'
    if (piniaStore.AdminInfo.GroupId == 3){
        setInterval(async ()=>{
            let info = await GetAdminInfo()

            if(info.GroupId == 3 && info.Businesses.Balance < 0){



                ElMessageBox.confirm(



                    `<div style="display: flex;align-items: center;justify-content: center;flex-wrap: wrap">
                        <div style="width: 100%; text-align: center">${t('商户余额不足，请联系客服及时充值！')}</div>
                        <div style="width: 100%; text-align: center">${t('如不及时充值将会关闭游戏服务')}</div>
                        <div style="width:100%;color: red;text-align: center">${t('当前余额:{Num}', {Num: ut.toNumberWithCommaNormal(info.Businesses.Balance.toFixed(2))})}</div>
                    </div>`,
                    t('商户余额提示'),
                    {
                        dangerouslyUseHTMLString: true,
                        cancelButtonText: t('关闭'),
                    }
                )
                    .then(async () => {


                    })
            }
        },piniaStore.AdminInfo.Businesses.ArrearsThresholdInterval)
    }



});

const handleTabClick = (tabPaneContext,name) =>{

    let tabClick = tabMenuList.value.find((list) => translateTitle(list.Title) === tabPaneContext.paneName);
    if (tabClick && route.path !== tabClick.Url) {
        setActiveTabs(translateTitle(tabClick.Title))
        let datas = {
            title:translateTitle(tabClick.Title),
            value:tabClick.Url
        }
        setElMenuActive(datas)
        router.push(tabClick.Url);
        // window.location.reload()
    }
}

const handleTabRemove = (name) =>{
    if (name === t('游戏报表')) return

    const tabIndex = tabMenuList.value.findIndex((list) => translateTitle(list.Title) === name);
    if (tabIndex !== -1){
        tabMenuList.value.splice(tabIndex, 1);
    }
    // 返回上一个标签页
    if (tabMenuList.value.length) {
        const lastTab = tabMenuList.value[tabMenuList.value.length - 1]
        setActiveTabs(translateTitle(lastTab.Title))
        router.push(lastTab.Url);
    } else {
        setActiveTabs(t('游戏报表'));
        router.push('/dashboard');
    }


}

const reloadRouterView = () => {
    // router
    isRouterAlive.value = false;
    nextTick(() => {
        isRouterAlive.value = true;
    })

}
</script>

<style lang="scss">
.container {
    width: 100%;
    overflow-y: hidden;
    flex-direction: row;
    display: flex;
    box-sizing: border-box;
}

.admin-container {
    padding: 15px;
}

.rightView {
    transition: 0.3s all;
    width: 100%;
    overflow: auto;
    z-index: 1;
    height: 100%;
    .el-tabs{
        padding: 0 10px;
    }
}

.router-view {
    width: 100%;
    min-height: calc(100vh - 86px);
    background: var(--el-bg-color-page);
    padding: 15px;

}

.transitionName-enter-from,
.transitionName-leave-to {
    transform: translateX(20px);
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

.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.userDetails {
    .el-tab-pane {
        width: var(--tabPaneWidth);
        display: flex;
        justify-content: center;
        flex-wrap: wrap;

        .el-form {
            width: 100%;
        }
    }

}
.el-descriptions {
    .el-descriptions__body{
        box-shadow: 0 2px 6px 0 rgba(18, 31, 62, 0.1);
    }
}
.elTable {
    width: 100%;
    box-shadow: 0 2px 6px 0 rgba(18, 31, 62, 0.1);
    .cell {
        padding: 0 10px;
    }

}

.expand {
    width: 100%;
    div {
        padding-bottom: 10px;

        h4 {
            padding-bottom: 6px;
        }
        p{
            padding-bottom: 6px;
        }
    }
}

.tableImg {
    width: 200px;
    height: 200px;
    margin: 0 auto;
    text-align: center;
    display: flex;
}

.noneA {
    text-decoration: none;
}

.copy {
    margin-left: 5px;
    cursor: pointer;
}
:deep(.el-tabs__header){
    margin: 0.5rem;
    background: #e3e7ee;
}

:deep(.el-tabs__item){
    width: auto;
}
.page_history_list{
    height: auto;
    background-color: var(--el-bg-color);
    border: 1px solid var(--el-border-color-light);
    box-shadow: 0 1px 1px var(--el-box-shadow-light);
    overflow-x: auto;
    white-space: nowrap;
    overflow-y: hidden;
    display: flex;
    align-items: flex-end;
}



.page_history_list::-webkit-scrollbar-thumb:hover {
    width: 4px!important;
    background: var(--el-color-primary)
}
.page_history_item:first-child{
    margin-left: 0.8rem;
    display: flex;
    align-items: flex-end;
}
.page_history_item:hover{
    background-color: #f4f4f5;
    border-radius:10px 10px 0 0!important;
    z-index: 0!important;
}
.page_history_item_active{
    background-color: #e8f1ff!important;
    font-weight: 600;
    color: var(--el-color-primary);
    border-radius:10px 10px 0 0!important;

}

.page_history_item{
    width: auto;
    height: auto;
    display: inline-block;
    cursor: pointer;
    padding: 8px 12px;
    border-radius: 10px;
    position: relative;
    font-size: $size16;
    z-index: 999999;
    background: #ffffff;
    margin-left: 0;
}
.page_content{
    width: 100%;
    height: calc(100vh - 120px);
    padding: 8px 15px;
    overflow-y: auto;
}

.page_history_item_active .el-menu-icon{
    border-bottom: 0.8rem solid #e8f1ff;
    display: block;
}
.page_history_item:hover .el-menu-icon{
    display: block;
}
.el-menu-icon{
    position: absolute;
    background: rgba(0,0,0,0);
    width: 0;
    height: 0;
    bottom: 0;
    z-index: 99999;
    border-bottom: 0.8rem solid #f4f4f5;
    display: none;
}



.ele-tab-corner-left{
    left: -0.5rem;
    right: 0;
    border-right: 0 solid transparent;
    border-left: 1rem solid transparent;
    display: none;
}
.ele-tab-corner-right{
    right: -0.5rem;
    border-left: 0 solid transparent;
    border-right: 1rem solid transparent;
}
</style>
