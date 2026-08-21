<template>

    <HeaderComponent
        @openLanguageConfigDialog="openLanguageDialog"
        @openConnectConfigDialog="openConnectDialog"
        @openAboutUsConfigDialog="openAboutUsDialog"
        @openFilterConfigDraw="openDraw"
        @searchGameByName="searchByName"
    />
    <!--  内容  -->
    <el-row :gutter="24" justify="center">

        <!--    游戏列表部分    -->
        <el-col :md="18" :sm="22" :xs="22" >

            <AboutUs />

            <div class="bottom_text" id="bottom_text">
                <el-link :href="hrefText">{{ hrefText }}</el-link>
            </div>

            <div style="height: 90px" class="hidden-md-and-up"></div>

        </el-col>
    </el-row>




    <!--  系统语言  -->
    <Dialog v-model="languageDialogVisible" :title="systemText.language" width="60%" max-width="600px">
        <div class="dialog-list">

            <el-col class="language-item" :lg="12" :sm="6" :xs="12" :span="6" v-for="item in LanguageList"
                    :key="item.code" @click="changeSystemLanguage(item.code)">
                <el-image class="language-item--img" size="large" :src="item.icon ? `${HOST}/${item.icon}` : ''"
                          :class="{'language-item-active': language == item.code}" fit="fill"/>
                <div class="language-item--label"
                     :style="{'color': language == item.code ? 'rgb(249, 190, 0)' : ''}">{{ item.label }}
                </div>
            </el-col>

        </div>
    </Dialog>

    <!--  联系方式  -->
    <Dialog v-model="connectDialogVisible" :title="systemText.contactInformation" width="60%" max-width="600px">

        <div class="connect-list dialog-list">
            <div class="hidden-sm-and-down connect-item" v-for="(item, index) in aboutInfo.connect" :key="index">
                <el-link :href='item.href' target="_blank" :underline="false"
                         style="width: 100%;margin-bottom: 24px;margin-top: 24px;position: relative">

                    <div @click.stop>
                        <el-image :src="`${HOST}public/${item.image}`" fit="fill"></el-image>
                        <div class="language-item--label">{{ item.label }}</div>
                    </div>
                </el-link>
            </div>


            <div class="hidden-md-and-up connect-item" v-for="(item, index) in aboutInfo.connect" :key="index">
                <el-link :href='item.href' target="_blank" :underline="false"
                         style="width: 100%;margin-bottom: 24px;margin-top: 24px;position: relative">

                    <div @click.stop>
                        <el-image :src="`${HOST}public/${item.image}`" fit="fill"></el-image>
                        <div class="language-item--label">{{ item.label }}</div>
                    </div>
                </el-link>
            </div>
        </div>

    </Dialog>



    <!--  回到顶部   移动端 -->
    <div class="to_top hidden-md-and-up" v-show="isShow" @click="scrollTop">
        <el-icon color="#ffffff" size="24">
            <ArrowUpBold/>
        </el-icon>
    </div>

    <!--  回到顶部  PC端 -->
    <div class="to_top hidden-sm-and-down" v-show="isShow" @click="scrollTop">
        <el-icon color="#ffffff">
            <ArrowUpBold/>
        </el-icon>
    </div>
</template>


<script setup lang="ts">
import AboutUs from '@/views/component/aboutComponent.vue';
import MenuComponent from "@/views/component/menuComponent.vue";
import GameList from "@/views/component/gameList.vue";
import type {Ref} from "vue";
import {nextTick, onMounted, onUnmounted, ref} from "vue";
import type {GameInterface} from "@/interface/gameInterface";
import Dialog from "@/components/dialog.vue";
import {useLanguageStore} from "@/stores/store/langageStore";

import type {LanguageConfigInterface} from "@/interface/languageConfigInterface";
import {ArrowUpBold, Close} from "@element-plus/icons-vue";
import HeaderComponent from "@/views/component/headerComponent.vue";
import type {FilterInterface} from "@/interface/filterInterface";
import {register} from 'swiper/element/bundle';
import customSwiper from "./component/customswiper.vue";
import {GetGameList} from "@/api/Game";
import {ScrollWindow} from "@/util/windowUtil";
import {useGameStore} from "@/stores/store/game";
import type {LanguageInterface} from "@/interface/languageInterface";
import type {AboutInterface} from "@/interface/aboutInterface";
import {storeToRefs} from "pinia";
import type {manufacturerInterface} from "@/interface/manufacturer";
import type {ThemeInterface} from "@/interface/themeInterface";
import {Throttle} from "@/util/util";


const HOST: string = window.location.protocol + "//" + window.location.host + "/" + import.meta.env.VITE_SOURCE_BASE

const hrefText = window.location.href
const languageStore = useLanguageStore()
const gameStore = useGameStore()

const { LanguageList, systemText, aboutInfo, language }:{
    LanguageList: Ref<LanguageInterface[]>
    systemText: Ref<LanguageConfigInterface>
    aboutInfo: Ref<AboutInterface>
    language: string
} = storeToRefs(languageStore)

const { Manufacturer,ThemeList }: {
    Manufacturer:Ref<manufacturerInterface[]>,
    ThemeList:Ref<ThemeInterface[]>
} = storeToRefs(gameStore)

const loading: Ref<boolean> = ref<boolean>(true)

const menuRef: Ref<InstanceType<typeof MenuComponent>> = ref<InstanceType<typeof MenuComponent>>(null)

let sortValue:Ref<number | null> = ref(null)
let sortOptions:Ref<Array<Record<string, any>>> = ref<Array<Record<string, any>>>([])

let searchViewStatus = ref("")
const gameList: Ref<GameInterface[]> = ref<GameInterface[]>([])
const AllGameList: Ref<Record<string, GameInterface[]> | GameInterface[]> = ref<Record<string, GameInterface[]> | GameInterface[]>({})



// 游戏筛选条件
// 通过厂商筛选
const currentAllFilterConf: Ref<string[]> = ref<string[]>([])

// 用来展示的字段
const currentAllFilter: Ref<FilterInterface[]> = ref<FilterInterface[]>([])

const filterGameType: Ref<string> = ref<string>("ALL")
// 语言列表

// 语言选择面板是否显示
const languageDialogVisible: Ref<boolean> = ref<boolean>(false)
// 联系方式是否显示
const connectDialogVisible: Ref<boolean> = ref<boolean>(false)


// 抽屉开关
const drawerStatus: Ref<boolean> = ref<boolean>(false)

const searchName: Ref<string> = ref<string>("")
// 是否回到顶部
const isShow: Ref<boolean> = ref<boolean>(false)
const scrollHeight = ref(0)


// 获取游戏
const getGame = async () => {



    loading.value = true

    AllGameList.value = await GetGameList(
        systemText.value.gameList as Record<string, string>,
        searchName.value,
        currentAllFilter.value,
        filterGameType.value,
        Number(sortValue.value ? sortValue.value : 1)
    )


    setTimeout(()=>{

        generatorMap()
    },500)



    loading.value = false
}


// 获取语言配置 初始化游戏
const init = () => {

    systemText.value.sortList = systemText.value.sortList as Array<Record<string, any>>

    sortValue.value = systemText.value.sortList[0].value
    sortOptions.value =  systemText.value.sortList

    getGame()
    loading.value = false
}

const generatorMap = () => {



    if (sortValue.value == 2){
        return
    }


    let map = {}
    for (let index in Manufacturer.value) {



        let currentDom = `${Manufacturer.value[index].label}-container`


        let currentDomTop = document.getElementById(currentDom)?.offsetTop
        let currentDomHeight = document.getElementById(currentDom)?.offsetHeight

        const currentTitleName = Manufacturer.value[index].label


        let currentDomBottom = 99999
        if (currentDomTop && currentDomHeight){
            currentDomBottom = parseFloat(currentDomTop.toString()) + parseFloat(currentDomHeight.toString())
        }

        map[currentTitleName] = [
            currentDomTop,
            currentDomBottom
        ]

    }

    console.log(map)
    localStorage.setItem("scrollMap", JSON.stringify(map))
}

// 切换分类
const toManu = ({toScroll}: { toScroll: number }) => {

    ScrollWindow(toScroll)
}

const GameTop = () => {

    const offsetTop: any = document.getElementById('classificationContent')?.offsetTop
    ScrollWindow(offsetTop - 100)

}
const checkMenu = async ({named}: { named: string }) => {

    filterGameType.value = named



    currentAllFilter.value = []
    searchName.value = ""
    currentAllFilterConf.value = []
    await getGame()
    GameTop()
}

// 打开语言面板
const openLanguageDialog = () => {
    languageDialogVisible.value = !languageDialogVisible.value
}

const openAboutUsDialog = () => {
    aboutDialogVisible.value = !aboutDialogVisible.value
}

// 打开联系面板
const openConnectDialog = () => {
    connectDialogVisible.value = !connectDialogVisible.value
}
// 打开筛选面板
const openDraw = () => {
    drawerStatus.value = !drawerStatus.value
}

// 添加筛选的厂商
const manufacturerFilter = (appendFilter: FilterInterface, filterType: string) => {

    searchName.value = ""

    appendFilter["type"] = filterType
    let filterList = [...currentAllFilter.value]
    if (filterType == "theme"){

        const typeIndex = filterList.findIndex(item => item.type == appendFilter.type)
        const idIndex =  filterList.findIndex(item => item.id == appendFilter.id && item.type == appendFilter.type)
        if (typeIndex == -1){
            filterList.push(appendFilter)
        }else{
            filterList[typeIndex] = appendFilter
            // if (idIndex == -1){
            //     filterList.push(appendFilter)
            //
            // }else{
            //
            // }


        }

    }

    if (filterType == "manufacturer"){

        const index = filterList.findIndex(item => item.id == appendFilter.id && item.type == appendFilter.type)
        if (index == -1){
            filterList.push(appendFilter)
        }else{
            filterList.splice(index, 1)
        }

    }
    currentAllFilter.value = filterList


    currentAllFilterConf.value = currentAllFilter.value.map(item=>`${item.type}-${item.id}`)


    filterGameType.value = "ALL"

    systemText.value.sortList = systemText.value.sortList as Array<Record<string, any>>

    sortValue.value = Number(systemText.value.sortList[0].value)

    getGame()
}

// 清除所有筛选条件
const clearFilter = () => {

    currentAllFilter.value = []
    RenderData()

}

// 语言切换
const changeSystemLanguage = (code: string = "en") => {

    localStorage.setItem("systemLanguage", code)
    window.location.reload()
}



const filterHandle = async (index: number) => {


    currentAllFilter.value.splice(index, 1)
    RenderData()
}

const RenderData = () => {
    let filterConfData: string[] = []

    currentAllFilter.value.forEach(item => {
        filterConfData.push(`${item.type}-${item.id}`)
    })


    currentAllFilterConf.value = filterConfData

    getGame()

    GameTop()

    filterGameType.value = "ALL"

}



const searchByName = async (selectStr: string) => {
    searchName.value = selectStr
    await getGame()
}



const changeView = async (handleType: string) => {



    searchViewStatus.value = handleType

    filterGameType.value = "ALL"
    searchName.value = ""
    await getGame()

}


const scrollTop = () => {
    ScrollWindow(0)
}

const handleScroll = () => {

    isShow.value = window.pageYOffset >= 380;


}

onMounted(() => {

// 初始化
    setTimeout(() => {
        init()
    }, 400)


    register()

    window.addEventListener('scroll', handleScroll);
    window.addEventListener('resize', ()=>{
        nextTick(()=>{
            Throttle(()=>{
                generatorMap()
            }, 500)()

        })

    });



});

onUnmounted(() => {
    // 移除监听事件
    window.removeEventListener('scroll', null);
    window.removeEventListener('resize', ()=>{
        nextTick(()=>{
            Throttle(()=>{
                generatorMap()
            }, 500)()
        })
    });


});



</script>

<style scoped>

.dialog-list {
    display: flex;
    flex-wrap: wrap;
}


.language-item {
    width: 20%;
    height: auto;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-wrap: wrap;
    margin-bottom: 32px;

}

.language-item--label {
    width: 100%;
    text-align: center;
    margin-top: 10px;
    color: #ffffff;
}

.language-item--img {
    width: 60%;
    height: auto;
    border: 5px solid rgba(0, 0, 0, 0);
    border-radius: 10px;
}

.connect-list {
    justify-content: center;
    align-items: center;
}

.connect-item {
    position: relative;
    display: flex;
    -webkit-box-pack: center;
    justify-content: center;
    -webkit-box-align: center;
    align-items: center;
    flex-direction: column;
    width: 20%;
    margin: 0 1.66% 16px;
    cursor: pointer;
}

.linkItem:active .language-item--label {
    color: rgb(249, 190, 0);
}

.drawerContent {
    width: 94%;
    height: 70vh;
    overflow-y: auto;
    margin: 0 auto;
}

.clear--filter {
    width: 80%;
    height: auto;
    display: flex;
    justify-content: center;
}

.selectGameItem {

    float: left;
    color: rgba(255, 255, 255, 0.5);
    padding: 7px 20px;
    border-radius: 25px;
    border: 1px solid rgba(130, 136, 151, 0.6);
    cursor: pointer;
    margin-right: 16px;
    margin-bottom: 8px;
}

.selectGameItemActive {
    background-color: #4d93fb;
    color: rgb(255, 255, 255);
}

.clear--filter-button {
    width: 95%;
    padding: 18px 25px 17px;
    color: rgb(255, 255, 255);
    border-radius: 30px;
    border: 1px solid rgba(130, 136, 151, 0.6);
    background-color: rgb(20, 22, 32);
    cursor: pointer;
    opacity: 0.4;
    text-align: center;
}

.to_top {
    width: 62px;
    height: 62px;
    position: fixed;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    right: 175px;
    bottom: 148px;
    z-index: 999;
    background-color: #4d93fb;

}

.language-item-active {
    border: 5px solid #ffffff;

}

.all-title {
    font-size: 18px !important;
    height: auto;
    display: block;

}
.manufacturer-title{
    margin-top: -15px;
}

.game-type-title {
    width: 100%;
    font-size: 15px;
    font-weight: 600;
    color: #ffffff;
    margin-bottom: 30px;
    position: relative;
}

.filter-list-container {
    margin-bottom: 30px;
}

.filter-item {
    font-size: 16px;
    color: rgb(255, 255, 255);
    padding: 0 16px;
    border-radius: 25px;
    border: 1px solid rgba(130, 136, 151, 0.6);
    cursor: pointer;
    margin-right: 16px;
    margin-bottom: 8px;
    float: left;
    height: 50px;
    line-height: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.manufacturer-name {
    color: rgba(255, 255, 255, .5);
}

.hidden-md-and-up.to_top {
    width: 60px;
    height: 60px;
    right: 36px;
    bottom: 94px;
    z-index: 999;
    background-color: #4d93fb;
}

.hidden-md-and-up.connect-item {
    width: 50% !important;
}

.gameFilterTitle {
    width: 100%;
    height: auto;
    margin-top: 10px;
    margin-bottom: 20px;
    font-size: 16px;
    color: #ffffff;
}
</style>

<style>


.bottom_text {
    font-size: 13px;
    height: 122px;
    line-height: 122px;
    text-align: center;
    border-top: 1px solid rgba(130, 136, 151, 0.6);
    color: rgba(255, 255, 255, 0.5);
}

.class-game-border {
    font-size: 24px;
    position: relative;
}

.class-game-border:before {
    content: "";
    background: #4d93fb;
    width: 56px;
    height: 2px;
    display: block;
    position: absolute;
    bottom: 0;
}

.sortSelect {
    width: 140px;
    float: right;
}
</style>

<style>

.sortSelect .el-select__wrapper {

    border-color: #ffffff;
    background: #242937;
}


.sortSelect .is-focused {
    box-shadow: 0 0 0 1px #4d93fb inset !important;

}

.sortSelect .el-popper__arrow {
    background: #242937;
}


.el-drawer {
    background-color: rgb(20, 22, 32) !important;
}

.swiperContainer .el-carousel__indicator {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: #4d93fb;

}

.placeholder-slot {
    width: 100%;
    height: 100%;
    opacity: 0;
}

@keyframes clll {
    from {
        width: 10px;
        border-radius: 50%
    }
    to {
        width: 44px;
        border-radius: 4px;
    }
}

.dialog--header:before {
    content: "";
    width: 20%;
    height: 4px;
    background: #4d93fb;
    display: inline-block;
    position: absolute;
    bottom: -10px;
    left: 40%;
}

.mySwiper .el-image {
    height: 25vw;
}

.swiper-slide {
    background-position: center;
    height: 25vw;
    background-size: cover;

}


.mySwiper .el-image {
    height: 25vw;
}

.swiper-slide {
    background-position: center;
    height: 25vw;
    background-size: cover;

}


.hidden-md-and-up .mySwiper .el-image {
    height: 40vw !important;
}

.hidden-md-and-up .swiper-slide {
    background-position: center;
    height: 40vw !important;
    background-size: cover;

}


.swiper-slide:not(.swiper-slide-active) {
    filter: blur(5px);
    background-color: rgb(255, 255, 255, 0.1);
    box-shadow: 0 25px 45px rgba(0, 0, 0, 0.1);
}
.swiper-slide .el-image{
    width: 100%;
}
</style>
