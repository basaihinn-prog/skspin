<template>


  <!--  移动端显示  -->
    <div class="hidden-md-and-up" :class="{'sticky-class': isSticky}">


        <el-row :gutter="20" v-if="propSearchName == '' ">
            <el-col :span="18" :offset="!isSticky ? 0 : 1"
                    v-if="currentGameTypeId != 3 && props.sortValue != 2 && currentAllFilter.length == 0">

                    <span class="manuList">
                        <template v-for="item in manufacturer" :key="item.id">


                            <span class="manu_classification_item"

                                  :class="item.label == manuSelectCheck ? 'manu_classification_item_active' : ''"
                                  @click="toManuClassification(item.label, item.id)">{{ item.label }}</span>
                        </template>
                    </span>
            </el-col>
        </el-row>

    </div>

    <div class="menu--container hidden-md-and-up">
        <div class="menu--content" v-if="filterViewStatus == 'classification' && propSearchName == '' ">
            <div v-for="item in gameType" :key="item.id" class="menu--item"
                 v-if="isPhone()"
                 @touchstart.native.prevent="checkMenu(item.id, item.label)"
                 style="font-size: 14px">

                <el-image v-show="item.id !== currentGameTypeId" :src="`${menuMap[item.icon]}`" class="game_type_icon">
                    <template #placeholder>
                        <div></div>
                    </template>
                </el-image>
                    <div v-show="item.id !== currentGameTypeId" style="width: 100%;line-height: normal;text-align: center">{{ item.label }}</div>

                    <el-image v-show="item.id == currentGameTypeId" :src="`${menuMap[ 'select' +item.icon]}`" class="game_type_icon">
                        <template #placeholder>
                            <div></div>
                        </template>
                    </el-image>
                    <div v-show="item.id == currentGameTypeId" style="width: 100%;line-height: normal;text-align: center">{{ item.label }}</div>


            </div>
            <div v-for="item in gameType" :key="item.id" class="menu--item"
                 v-else
                 @click="checkMenu(item.id, item.label)"
                 style="font-size: 14px">

                <el-image v-show="item.id !== currentGameTypeId" :src="`${menuMap[item.icon]}`" class="game_type_icon">
                             <template #placeholder>
                                            <div></div>
                                        </template>
                </el-image>
                    <div v-show="item.id !== currentGameTypeId" style="width: 100%;line-height: normal;text-align: center">{{ item.label }}</div>

                <el-image v-show="item.id == currentGameTypeId" :src="`${menuMap[ 'select' +item.icon]}`" class="game_type_icon">
                             <template #placeholder>
                                            <div></div>
                                        </template>
                </el-image>
                    <div v-show="item.id == currentGameTypeId" style="width: 100%;line-height: normal;text-align: center">{{ item.label }}</div>


            </div>
        </div>

    </div>

  <!--  PC端显示  -->
    <div class="menu--container hidden-sm-and-down" :class="{'sticky': isSticky}">
        <div class="menu--content" v-if="filterViewStatus == 'classification'">
            <div v-for="item in gameType" :key="item.id" class="menu--item"
                 :class="item.id === currentGameTypeId ? 'checked--menu--item' : ''"
                 @click.stop="checkMenu(item.id, item.label)">
                {{ item.label }}
            </div>
        </div>
        <div class="menu--content" v-if="filterViewStatus == 'search'">
            <div class="search--container">
                <div class="search--content">

                    <el-input :placeholder="systemText.searchPlaceHolder" v-model="GamesearchName"
                              :autofocus="true"
                              id="searchNameId"
                              class="borderless"
                              :input-style="[{borderColor: 'rgba(255,255,255)'}]"
                              @input="searchGame">
                        <template #prefix>
                            <el-icon size="23" color="#000000">
                                <Search/>
                            </el-icon>
                        </template>
                    </el-input>

                </div>
                <div style="height: 50px;display: flex;align-items: center;margin-left: 20px" @click="clearSearName">
                    <el-icon size="23" color="#ffffff" style="font-weight: bold">
                        <CloseBold/>
                    </el-icon>
                </div>
            </div>
        </div>

        <div class="filter_game_container">
            <span class="filter_game_item" @click="openDialog" v-if="appId == 'faketrans'">
                <SetRTP :RTPDialogVisible="RTPVisible" @closeRTPDialogVisible="closeRTP"></SetRTP>

            </span>
            <span class="filter_game_item" @click="changeFilter">
                <el-icon size="24">
                    <Search/>
                </el-icon>
            </span>
            <span class="filter_game_item" @click="openDrawerPanel">
                <el-icon size="24">
                    <Filter/>
                </el-icon>
            </span>

        </div>


        <div class="manuContent"
             v-if="!propSearchName && currentAllFilter.length == 0 && currentGameTypeId != 3 && props.sortValue != 2">

            <el-row :gutter="20">

                <!--    banner轮播图部分    -->
                <el-col :span="18" :offset="!isSticky ? 0 : 3">


                    <span class="manuList">
                        <template v-for="item in manufacturer" :key="item.id">

                                <span class="manu_classification_item"

                                  :class="item.label == manuSelectCheck ? 'manu_classification_item_active' : ''"
                                  @click="toManuClassification(item.label, item.id)">{{ item.label }} </span>
                        </template>
                    </span>
                </el-col>

            </el-row>
        </div>

    </div>

</template>

<script setup lang="ts">
import type {Ref} from "vue"
import {ref, onMounted, onUnmounted, watch, nextTick} from "vue";
import {axiosComm} from "@/util/requestComm";
import type {GameTypeInterface} from "@/interface/menuInterface";
import {Throttle} from "@/util/util";
import {CloseBold, Filter, Search} from "@element-plus/icons-vue";
import type {manufacturerInterface} from "@/interface/manufacturer";

import all from "@/assets/all.png"
import mini from "@/assets/mini.png"
import slot from "@/assets/slot.png"
import selectAll from "@/assets/selectall.png"
import selectMini from "@/assets/selectmini.png"
import selectSlot from "@/assets/selectslot.png"
import SetRTP from "@/components/setRTP.vue";
import {useLanguageStore} from "@/stores/store/langageStore";
import {storeToRefs} from "pinia";
import {LanguageConfigInterface} from "@/interface/languageConfigInterface";
import {useGameStore} from "@/stores/store/game";

const props = defineProps({
    sortValue: {
        type: Number,
        default: -1
    },
    currentAllFilter: {
        type: Array<any>,
        default: []
    },
    languageRef: {
        type: Object,
        default: []
    },
    propSearchName: {
        type: String,
        default: ""
    }
})

const appId = localStorage.getItem("app_id")

// 注册事件
const emit = defineEmits(["checkMenu", "openDrawer", "searchGameByName", "changeFilterView", "toManu", "toGameList"])
const languageStore = useLanguageStore()
const gameStore = useGameStore()
const { systemText }:{systemText:LanguageConfigInterface} = storeToRefs(languageStore)
const { Manufacturer }:{manufacturer: manufacturerInterface} = storeToRefs(gameStore)

const manufacturer: Ref<manufacturerInterface[]> = Manufacturer.value

const RTPVisible: Ref<boolean> = ref<boolean>(false)

// 当前的游戏分类
const currentGameTypeId: Ref<number> = ref<number>(0)

const GamesearchName: Ref<string> = ref<string>("")

// 查询视图状态
const filterViewStatus: Ref<string> = ref<string>("classification")

// 游戏分类
const gameType: Ref<GameTypeInterface[]> = ref<GameTypeInterface[]>([])

const HOST: string =  window.location.protocol + "//" + window.location.host + import.meta.env.VITE_SOURCE_BASE
axiosComm.get("/mock/gameType.json").then((res) => {
    currentGameTypeId.value = res.data[0].id
    gameType.value = res.data

})
const isPhone = () => {
    const userAgentInfo = navigator.userAgent;
    const mobileAgents = ["Android", "iPhone", "SymbianOS", "Windows Phone", "iPad", "iPod"];
    const mobileFlag = mobileAgents.some((mobileAgent) => {
        return userAgentInfo.indexOf(mobileAgent) > 0;
    });

    return mobileFlag;
}


const menuMap = ref({
    "all.png" : all,
    "mini.png" : mini,
    "slot.png" : slot,
    "selectall.png" : selectAll,
    "selectmini.png" : selectMini,
    "selectslot.png" : selectSlot,
})

const isSticky = ref(false);
const manuSelectCheck = ref("PG");
const stickyOffset = ref(80);

const openDialog = () => {
    RTPVisible.value = true
}
const closeRTP = () => {
    RTPVisible.value = false
    event.stopPropagation()

}
function handleScroll(manuSelectCheckAttr?: number) {


    if (window.pageYOffset >= stickyOffset.value) {
        isSticky.value = true;
    } else {
        isSticky.value = false;
    }


    Throttle(()=>{


        let currentWindowTop = window.pageYOffset

        const scrollManuMapData = JSON.parse(localStorage.getItem("scrollMap"))

        if (scrollManuMapData[manuSelectCheck.value] && currentWindowTop <= scrollManuMapData[manuSelectCheck.value][0]) {
            manuSelectCheck.value = "PG"
        }

        for (let key in scrollManuMapData) {
            if (scrollManuMapData[key][0] != 0 && scrollManuMapData[key][1] != 0) {


                if ((scrollManuMapData[key][0] - 250) < currentWindowTop && (scrollManuMapData[key][1] - 250) > currentWindowTop) {
                    manuSelectCheck.value = key
                    break
                }
            }
        }


    }, 400)()
}

onMounted(() => {
    // 获取吸顶元素的偏移位置
    const element = document.querySelector('.menu--container');

    if (element) {
        stickyOffset.value = window.innerHeight * 0.7
    }

    // 监听滚动事件
    window.addEventListener('scroll', ()=>{
        handleScroll()
    });
});

onUnmounted(() => {
    // 移除监听事件
    window.removeEventListener('scroll', ()=>{

        handleScroll()
    });
});


const checkMenu = (id: number, named: string) => {

    currentGameTypeId.value = id
    handleScroll()



    emit("checkMenu", {named})
}

// 打开抽屉面板
const openDrawerPanel = () => {
    emit("openDrawer")
}


const changeFilter = () => {
    filterViewStatus.value = filterViewStatus.value === "search" ? "classification" : "search"
    GamesearchName.value = ""
    currentGameTypeId.value = 1

    emit("changeFilterView", filterViewStatus.value)

}

const searchGame = (value) => {

    Throttle(() => {
        emit("searchGameByName", GamesearchName.value)
    }, 300)()

}

const searchFocus = () => {
    emit("toGameList")
}
const clearSearName = () => {
    filterViewStatus.value = filterViewStatus.value === "search" ? "classification" : "search"
    GamesearchName.value = ""
    emit("searchGameByName", "")
    emit("changeFilterView", filterViewStatus.value)
}
const toManuClassification = (manufacturerName: string, toScroll) => {



    const scrollManuMap = JSON.parse(localStorage.getItem("scrollMap"))


    let scrollNum = scrollManuMap[manufacturerName]

    if (scrollNum) {
        toScroll = scrollNum[0]
    } else {
        toScroll = null
    }


    //
    // let init = 80
    // if (!isPhone()){
    //     init = 210
    // }



    emit("toManu", {toScroll: toScroll - 185})
}


defineExpose({
    currentGameTypeId,
    handleScroll,
    manuSelectCheck
})
</script>

<style scoped>
.menu--container {
    width: 100%;
    height: 80px;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 40px;
    color: rgb(255, 255, 255);
    cursor: pointer;
    font-size: 18px;
    position: relative;

}

.menu--content {
    width: 100%;
    height: 80px;
    display: flex;
    justify-content: center;
    align-items: center;
    color: rgb(255, 255, 255);
    cursor: pointer;
    font-size: 18px;
    position: relative;
}

.menu--item {
    width: auto;
    height: auto;
    margin: 0 40px;
    line-height: 31px;
    position: relative;

}

.checked--menu--item:before {
    width: 100%;
    height: 4px;
    background: url("@/assets/daohang-hengxian.png") no-repeat 100%;
    position: absolute;
    bottom: -2px;
    left: 0;
    content: '';
    background-size: 100% 2px;
}


.sticky {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 1000; /* 确保吸顶效果在最上方 */
    backdrop-filter: blur(25px);
    box-shadow: rgba(0, 0, 0, 1) 0 2px 16px 0;
    background-color: rgba(0, 0, 0, 0.7);
}

.sticky-class {
    position: fixed;
    top: 64px;
    left: 0;
    width: 100%;
    z-index: 200; /* 确保吸顶效果在最上方 */
    backdrop-filter: blur(25px);
    box-shadow: rgba(0, 0, 0, 1) 0 2px 16px 0;
    background-color: rgba(0, 0, 0, 0.7);
}

.filter_game_container {
    width: auto;
    height: 80px;
    position: absolute;
    right: 0;
    top: 0;
    display: flex;
    justify-content: center;
    align-items: center;

}

.filter_game_item {
    width: 60px;
    height: 60px;
    margin-right: 15px;
    border-radius: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    background-color: rgb(0, 0, 0);
    cursor: pointer;

}

.search--container {
    width: 65%;
    height: 80px;
    display: flex;
    justify-content: center;
    align-items: center;

}

.search--content {
    width: 70%;
    height: 50px;
    display: flex;
    align-items: center;
    background-color: #ffffff;
    padding: 5px 8px;
    border-radius: 70px;
}

.hidden-md-and-up.menu--container {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 84px;
    z-index: 999;
    background-color: rgba(0, 0, 0, 0.85);
    backdrop-filter: blur(20px);
    box-shadow: rgba(255, 255, 255, 0.3) 0 -0.5px 0 0;
    margin-bottom: 0;
}

.hidden-md-and-up .menu--item {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    line-height: normal;
    height: auto;
}

.game_type_icon {
    width: 36px;
    height: 36px;
    display: block;
}

.manuContent {
    width: 100%;
    height: 60px;
    position: absolute;
    right: 0;
    bottom: -60px;
}

.manuList {
    width: auto;
    align-items: center;
    display: inline-flex;
    height: 48px;
    background: #13161F;
    border-radius: 40px 40px 40px 40px;
    margin-top: 5px;
    color: #ffffff;

}





.sticky .manuContent {
    backdrop-filter: blur(45px);
}

.sticky .manuList {
    justify-content: flex-end;
}

.manu_classification_item {
    padding: 5px 20px;
    border-radius: 50px;
    text-align: center;
    min-width: 87px;
}

.manu_classification_item_active {
    background: #4d93fb;
}
.menu--item:active{
    background: rgba(0, 0, 0, 0.85);
}


@media (max-width: 640px) {
    .manuList {
        overflow: auto;
        width: 92vw;
    }

    .manu_classification_item{
        padding: 5px 12px;
        width: auto;
        min-width: unset;
    }
}


</style>

<style>
.borderless .el-input__wrapper {
    border: none !important;
    box-shadow: none !important;
}

</style>
