<template>
    <div>
        <customTable :tableHeader="languageTableConfig" :tableData="tableData">
            <template #handleTools>


                <uploadExcel :fileName="fileName">
                    <el-button type="primary" plain>{{ $t('上传') }}</el-button>
                </uploadExcel>

            </template>

            <template #gameName="scope">
                <div style="display: flex;align-items: center;">
                    <el-tag size="small" type="primary">{{ scope.scope.GameManufacturerName }}</el-tag>
                    <div style="margin: 0 10px">
                        <el-avatar :src="scope.scope.GameIcon" size="small" />
                    </div>
                    {{ scope.scope.gameName }}
                </div>
            </template>

            <template #uploadFile="scope">

                <UploadIcon @uploadFile="uploadFileRequest($event, scope.scope)" :upload-file-keys="0">
                    <el-button size="small" type="primary" plain>{{ $t('上传') }}</el-button>
                </UploadIcon>

            </template>

        </customTable>

    </div>
</template>
<script setup lang="ts">

import {onMounted, ref} from "vue";
import CustomTable from "@/components/customTable/tableComponent.vue";
import {Client} from "@/lib/client";
import UploadExcel from "@/pages/systemManagement/config/component/component/uploadExcel.vue";
import UploadIcon from "./component/uploadImage.vue";
import Gamelist_container from "@/components/gamelist_container.vue";
import {Language} from "@/api/adminpb/language";
import {tip} from "@/lib/tip";
import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";
import UploadComponent from "@/components/customTable/uploadComponent.vue";
import {storeToRefs} from "pinia";
import axios from "axios";
import {Upload} from "@/api/comm";

const {t} = useI18n()
const store = useStore()


const {GameList} = storeToRefs(store)
const {fileName} = defineProps(["fileName"])

const languageTableConfig = ref([

    {label:"游戏名称", value:"gameName", type:"custom", width:"300px"},
    {label:"游戏类型", value:"gameType"},
    {label:"Id", value:"id"},
    // {label:"厂商", value:"manufacturer"},
    {label:"厂商名称", value:"manufacturerName"},
    {label:"运行状态", value:"runStatus"},
    {label:"排序", value:"sort"},
    {label:"新游戏排序", value:"newSort"},
    {label:"上传游戏图标", value:"uploadFile",type:"custom", width:"120px"},

])
const tableData = ref([])
const filterStatus = ref(false)
const dialogVisible = ref(false)
const GameId = ref("")
const Manufacturer = ref("")
const GameMap = ref({})


onMounted(() => {
    console.log("----------------------GameList---------------------")
    console.log(GameList.value)
    console.log("----------------------GameList---------------------")
    for (const i in GameList.value) {


            let item = GameList.value[i]
            GameMap.value[item.ID] = item

    }


    initTips()
})

const openDialog = () => {

}

const defaultGameEvent = ref({})
const selectGameList = (value) => {

    GameId.value = value.gameData || null
    Manufacturer.value = value.manufacturer || null


}





const initTips = async () => {
    const [response, err] = await Client.Do(Language.SelectConfig, {fileName: fileName} as any)
    if (err) return tip.e(t(err))

    let GameNameConfig = {
        "SPRIBE-1":"Mines",
        "SPRIBE-2":"Aviator",
        "SPRIBE-3":"Goal",
        "SPRIBE-4":"HiLo",
        "SPRIBE-5":"HotLine",
        "SPRIBE-6":"Keno",
        "SPRIBE-7":"Keno80",
        "SPRIBE-8":"MiniRoulette",
        "SPRIBE-9":"Plinko",
        "SPRIBE-10":"Dice",
        "SPRIBE-11":"Limbo"

    }

    tableData.value = response.Context.map(item => {


        let gameConfig = GameMap.value[`${item.manufacturerName.toLowerCase()}_${item.id}`]


        let GameId = gameConfig && gameConfig.ID || ""
        let gameName = ""
        if (gameConfig && gameConfig.Name){
            gameName = gameConfig.Name
        }else{
            gameName = GameNameConfig[`${item.manufacturerName.toUpperCase()}-${item.id}`]
        }
        let GameManufacturerName = item.manufacturerName

        let GameIcon = ""
        if (gameConfig && gameConfig.IconUrl){
            GameIcon = gameConfig.IconUrl
        }else{
            GameIcon = `https://admin.ppply.site/BHdownload/${item.manufacturerName.toUpperCase()}-${item.id}.webp`
        }
        return {
            ...item,
            GameId,
            gameName,
            GameManufacturerName,
            GameIcon,

        }
    })
    console.log(tableData.value)



}

const uploadFileRequest = async (uploadMap,game)=> {

    let formData = new FormData();
    formData.append('file', uploadMap.file);
    formData.append('fileName', uploadMap.file);
    formData.append('manufacturer', game.GameManufacturerName);
    formData.append('game_id', game.id);
    formData.append('id', game.GameId);



    const [response, err] = await Client.Do(Upload.UploadFile, formData )

    if (err){
        return tip.e(t("上传失败"))
    }


}

</script>
<style scoped lang="scss">

</style>
