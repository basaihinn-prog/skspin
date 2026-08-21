<template>


  <!--  顶部 PC -->
    <div class="pageHeader hidden-sm-and-down ">
        <el-row>
            <el-col :span="20" :offset="2">
                <div class="pageHeader--container">

                    <div>

                    </div>

                    <a href="/" class="logo_container">
                        <el-image :src="Logo" style="height: 45px">
                            <template #placeholder>
                                <div></div>
                            </template>
                        </el-image>
                    </a>
                    <div class="pageHeader--operation">

                        <div class="header-icon-item" @click="openAboutUsDialog">
                            <el-icon color="#ffffff" style="margin-right: 8px">
                                <InfoFilled />
                            </el-icon>
                            {{ systemText.aboutUs }}
                        </div>

                        <div class="header-icon-item" @click="openLinkDialog" style="display: none">
                            <el-icon color="#ffffff" style="margin-right: 8px">
                                <Link/>
                            </el-icon>
                            {{ systemText.link }}
                        </div>

<!--                        // TODO-->
                        <div class="header-icon-item" @click="openConnectDialog" v-if="envName != 'VIP环境'">
                            <el-icon color="#ffffff" style="margin-right: 8px">
                                <ChatSquare/>
                            </el-icon>
                             {{ systemText.contactInformation }}
                        </div>
                        <div class="header-icon-item"
                             style="display: flex; align-items: center;justify-items: center;position: relative;z-index: 999"
                             @click="openLanguageDialog">
                            <el-image :src="LanguageIcon" fit="contain" align="middle" style="width: 16px;margin-right: 8px;margin-top: 10px">
                                <template #placeholder>
                                    <div></div>
                                </template>
                            </el-image>
                            {{ systemText.language }}
                        </div>
                    </div>
                </div>
            </el-col>
        </el-row>
    </div>

  <!--  顶部 移动 -->
    <div class="pageHeader hidden-md-and-up" id="searchContentId">
        <el-row :gutter="24">
            <el-col :span="22" :offset="1" v-if="searchStatus == 'init'">
                <div class="pageHeader--container">

                    <div class="pageHeader--container--left" >
                        <div class="pageHeader--container--panelIcon" @click="openPanel" v-if="envName != 'VIP环境'">
                            <el-icon size="25" color="#ffffff"><Expand /></el-icon>
                        </div>
                        <el-image :src="Logo" style="width: 80px;height: 25px" fit="fill">
                            <template #error>
                                <div class="image-slot">
                                    <el-image class="game_img" :src="failImg" fit="fill">
                                        <template #placeholder>
                                            <div></div>
                                        </template>
                                    </el-image>
                                </div>
                            </template>
                            <template #placeholder>
                                <div class="image-slot">
                                    <el-image class="game_img game_img_load" :src="loadding" fit="fill">
                                        <template #placeholder>
                                            <div></div>
                                        </template>
                                    </el-image>
                                </div>
                            </template>
                        </el-image>

                    </div>



                    <div class="pageHeader--operation">
                        <span class="header-icon-item" @click="RTPVisible = true" v-if="appId == 'faketrans'" >
                            <SetRTP :RTPDialogVisible="RTPVisible" @closeRTPDialogVisible="closeRTP"></SetRTP>

                        </span>
                        <span class="header-icon-item" @click="changeSearchStatus">
                            <el-image :src="SearchIcon" fit="fill" style="width: 25px" :style="{opacity:RTPVisible ? '0.3': null}" v-if="!isAboutPage">
                                         <template #placeholder>
                                            <div></div>
                                        </template>
                            </el-image>
                        </span>
                        <span class="header-icon-item" @click="openFilterDraw" :style="{opacity:RTPVisible ? '0.3': null}" v-if="!isAboutPage">

                            <el-image :src="FilterIcon" fit="fill" style="width: 25px">
                                         <template #placeholder>
                                            <div></div>
                                        </template> 
                            </el-image>
                        </span>
                        <span class="header-icon-item" @click="openLanguageDialog" :style="{opacity:RTPVisible ? '0.3': null}">
                            <el-image :src="LanguageIcon" fit="fill" style="width: 25px">
                                         <template #placeholder>
                                            <div></div>
                                        </template>
                            </el-image>
                        </span>
                    </div>
                </div>
            </el-col>
            <el-col :span="22" :offset="1" v-if="searchStatus == 'search'">
                <div class="pageHeader--container--center">
                    <div class="search--content">

                        <el-input :placeholder="systemText.searchPlaceHolder" v-model="searchName" class="borderless"
                                  autofocus
                                  id="searchHeaderId"
                                  :input-style="[{borderColor: 'rgba(255,255,255)'}]" @input="searchGame">
                            <template #prefix>
                                <el-icon size="23" color="#000000">
                                    <Search/>
                                </el-icon>
                            </template>
                        </el-input>

                    </div>
                    <div style="height: 50px;display: flex;align-items: center;margin-left: 20px"
                         @click="changeSearchStatus">
                        <el-icon size="23" color="#ffffff" style="font-weight: bold">
                            <CloseBold/>
                        </el-icon>
                    </div>

                </div>

            </el-col>
        </el-row>
    </div>


  <!--  筛选的抽屉  -->
    <el-drawer
            size="50%"
            v-model="drawerStatus"
            title="I am the title"
            direction="ltr"
            custom-class="selectGamePanel"
            :show-close="false"
    >
        <template #header>
            <div style="display: flex;justify-content: right;align-items: center">

                <el-icon size="30" @click="openPanel">
                    <Close/>
                </el-icon>
            </div>


            <div></div>
            <div></div>
        </template>
        <div class="pageHeader--operation drawList">
            <div class="header-icon-item" @click="openAboutUsDialog(true)" >
                <el-icon color="#ffffff" style="margin-right: 8px">
                    <InfoFilled />
                </el-icon>
                {{ systemText.aboutUs }}
            </div>
            <div class="header-icon-item" @click="openLinkDialog(true)" style="display: none">
                <el-icon color="#ffffff" size="22" style="margin-right: 12px">
                    <Link/>
                </el-icon>
                {{ systemText.link }}
            </div>
            <div class="header-icon-item"  @click="openConnectDialog(true)"  v-if="envName != 'VIP环境'">
                <el-icon color="#ffffff" size="22" style="margin-right: 12px">
                    <ChatDotSquare />
                </el-icon>
                {{ systemText.contactInformation }}
            </div>

        </div>

    </el-drawer>

</template>

<script setup lang="ts">

import {
    ChatSquare,
    Setting,
    Link,
    Picture,
    Expand,
    Close,
    Filter,
    Search,
    CloseBold,
    ChatDotSquare
} from "@element-plus/icons-vue";
import LanguageIcon from "@/assets/Language.png"
import SearchIcon from "@/assets/search.png"
import Logo from "@/assets/logo.png"
import FilterIcon from "@/assets/filter.png"
import type {LanguageConfigInterface} from "@/interface/languageConfigInterface";

import {ElIcon} from "element-plus";
import type {Ref} from "vue";
import {nextTick, onMounted, ref} from "vue";
import {Throttle} from "@/util/util";
import failImg from "@/assets/failImg.png";
import loadding from "@/assets/loadding.webp";
import {useLanguageStore} from "@/stores/store/langageStore";
import SetRTP from "@/components/setRTP.vue";
import {storeToRefs} from "pinia";

const emit = defineEmits(["openLanguageConfigDialog", "openLinkConfigDialog", "openAboutUsConfigDialog", "openConnectConfigDialog", "openFilterConfigDraw", "searchGameByName", "manuPosition", "changeFilterView"])
const drawerStatus: Ref<boolean> = ref<boolean>(false)
const searchStatus: Ref<string> = ref<string>("init")

const appId = localStorage.getItem("app_id")

const searchName: Ref<string> = ref<string>("")

const envName = import.meta.env.VITE_NAME

const languageStore = useLanguageStore()

const { systemText }:{systemText:LanguageConfigInterface} = storeToRefs(languageStore)

const RTPVisible = ref(false)
const closeRTP = () => {
    RTPVisible.value = false
    event.stopPropagation()
}

const isAboutPage = ref(false);

onMounted(() => {
  isAboutPage.value = window.location.pathname === '/about-us'; 
});

const changeSearchStatus = () => {

    if (RTPVisible.value){
        RTPVisible.value = false
        return
    }

    searchStatus.value = searchStatus.value == "init" ? "search" : "init"
    searchName.value = ""
    nextTick(()=>{
        document.getElementById("searchHeaderId")?.focus()
    })

    emit("searchGameByName", "")
    emit("changeFilterView", searchStatus.value)
}
const openLinkDialog = (isCloseDraw?: boolean) => {
    if (isCloseDraw){
        drawerStatus.value = false
    }
    emit("openLinkConfigDialog")
}

const openAboutUsDialog = (isCloseDraw?: boolean) => {
    if (isCloseDraw){
        drawerStatus.value = false
    }
    window.location.href = '/about-us'; 

    emit("openAboutUsConfigDialog")
}

const openFilterDraw = () => {

    if (RTPVisible.value){
        RTPVisible.value = false
        return
    }

    emit("openFilterConfigDraw")
}
const openConnectDialog = (isCloseDraw?: boolean) => {
    if (isCloseDraw){
        drawerStatus.value = false
    }

    emit("openConnectConfigDialog")
}
const openLanguageDialog = () => {

    if (RTPVisible.value){
        RTPVisible.value = false
        return
    }
    emit("openLanguageConfigDialog")
}

const openPanel = () => {
    drawerStatus.value = !drawerStatus.value
}
const searchGame = () => {
    Throttle(() => {
        emit("searchGameByName", searchName.value.trim())
    }, 800)()
}
const clearSearName = () => {
    searchName.value = ""
    emit("searchGameByName", "")
}

</script>

<style scoped>

.pageHeader--container {
    width: 100%;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;

}

.pageHeader--operation {
    display: flex;
    align-items: center;
    position: relative;
    z-index: 999;
}
.pageHeader--operation .header-icon-item{

    font-size: 14px;
}

.pageHeader--icon {
    width: auto;
    height: 60px;
}
.header-icon-item {
    width: auto;
    height: auto;
    display: flex;
    align-items: center;
    color: #ffffff;
    margin-right: 30px;
    cursor: pointer;
    font-size: 14px;
    position: relative;
    z-index: 999;

}
.header-icon-item:hover{
    color: #4d93fb;
}


.pageHeader {
    width: 100%;
    height: 64px;
    background-color: rgba(0, 0, 0, 0.7);
    position: relative;
    z-index: 999;
}

.logo_container{
    position: absolute;
    width: 100%;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    left: 0;
    top: 0;
    z-index: 9;
}
.hidden-md-and-up .header-icon-item {
    width: auto;
    height: 64px;
    margin-left: 22px;
    display: flex;
    justify-items: center;
    margin-right: 0;
    font-size: 12px
}

.pageHeader--container--left {
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.pageHeader--container--panelIcon {
    width: 45px;
    height: 64px;
    display: flex;
    align-items: center;
    font-size: 10px;
}

.pageHeader--container--center {
    width: 100%;

    height: 64px;
    display: flex;
    justify-content: center;
    align-items: center;
}

.search--content {
    width: 80%;
    height: 50px;
    display: flex;
    align-items: center;
    background-color: #ffffff;
    padding: 5px 8px;
    border-radius: 70px;
}
.hidden-md-and-up.pageHeader{
    position: fixed;
    top: 0;
    left: 0;
    background-color: rgba(0, 0, 0, 1);
}
.drawList{
    width: 100%;
    display: block;


}
.drawList .header-icon-item{
    width: 100%;
    margin: auto;
    padding: 20px 0;


}
.placeholder-slot{
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0);
}

.image-slot{
    max-width: 228px;
    position: absolute;
    width: 100%;
    height: 100%;
    display: flex;
    background: rgba(0,0,0,0);
    align-items: center;
    justify-items: center;
}
</style>
