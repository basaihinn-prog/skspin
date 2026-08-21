<template>
    <el-dialog
            v-model="props.modelValue"
            :title="$t('游戏留存详情')"
            :width="store.viewModel === 2 ? '85%' : '960px'"
            @open="checkInfo"
            @close="emits('update:modelValue')"
            align-center
    >

        <div class="searchView gameList">
            <el-form
                :inline="true"
                @keyup.enter="initInfoData"
            >
                <el-space wrap>
                    <gamelist_container
                        :haseManufacturer="true"
                            :is-init="true"
                            :hase-all="true"
                            :defaultGameEvent="defaultGameEvent"
                            @select-operator="selectGameList"
                  />
                </el-space>
            </el-form>
        </div>
        <customTable
                v-loading="loading"
                table-name="gameUserRetainedInfo"
                height="400"
                :table-header="childTableHeader"
                :table-Data="childTableData"
                :page-size="uiData.PageSize"
                :page="uiData.Page"
                :count="count"
                @refresh-table="initInfoData"
                @pageChange="pageChange"
        >

            <template #GameID="scope">

                <div style="display: flex;align-items: center;">
                    <el-tag size="small" type="primary">{{ scope.scope.GameGameManufacturer }}</el-tag>

                    <div style="margin: 0 10px">
                        <el-avatar :src="scope.scope.GameIcon" size="small" />
                    </div>
                    {{ scope.scope.GameName }}
                </div>
            </template>
        </customTable>

    </el-dialog>
</template>

<script setup lang="ts">

import {useStore} from "@/pinia";
import {onMounted, Reactive, reactive, Ref, ref} from "vue";
import {useI18n} from "vue-i18n";
import ut from "@/lib/util";

const {t} = useI18n()
const store = useStore()

import customTable from "@/components/customTable/tableComponent.vue";
import {Client} from "@/lib/client";
import {GameUserRetained, GameUserRetainedReq} from "@/api/adminpb/gameUserRetained";
import Gamelist_container from "@/components/gamelist_container.vue";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {tip} from "@/lib/tip";

const props = defineProps(["modelValue", "ReportId"])
const emits = defineEmits(["update:modelValue"])

const loading = ref(false)
const defaultGameEvent = ref({})

const selectGameList = (value) => {

   uiData.value.GameId = value.gameData || null
   uiData.value.Manufacturer = value.manufacturer || null

    initInfoData()
}


const formatTime = (time) => {
    let timeStr = ut.fmtDate(new Date(time).getTime()/1000)
    return ut.fmtDate(timeStr, "YYYY-MM-DD")


}
// 留存详情列表头
const childTableHeader = ref([
    {label: "日期", value: "Date", format: (row) => formatTime(row.Date),width: "150px"},

    {label: "游戏名称", value: "GameID", type: "custom", width:"300px"},
    {label: "投注用户数", value: "RetentionPlayerCount", width: "160px"},
    {label: "次留", value: "RetentionPlayer1d", format:(row) => `${row.RetentionPlayer1d.toFixed(2)}%`, width: "100px"},
    {label: "三日留存", value: "RetentionPlayer3d", format:(row) => `${row.RetentionPlayer3d.toFixed(2)}%`, width: "100px"},
    {label: "七日留存", value: "RetentionPlayer7d", format:(row) => `${row.RetentionPlayer7d.toFixed(2)}%`, width: "100px"},
    {label: "十四日留存", value: "RetentionPlayer14d", format:(row) => `${row.RetentionPlayer14d.toFixed(2)}%`, width: "100px"},
    {label: "三十日留存", value: "RetentionPlayer30d", format:(row) => `${row.RetentionPlayer30d.toFixed(2)}%`, width: "100px"},
])


// 留存详情列表
const childTableData = ref([])
const count = ref(0)


const uiData: Ref<GameUserRetainedReq> = ref<GameUserRetainedReq>(<GameUserRetainedReq>{
    ReportId: "",
    Operator: 0,
    Manufacturer: "",
    GameId: "",
    Page: 1,
    PageSize: 20,
})


let initGameList = ref(null)
let GameManufacturerGroup = reactive({})
const initInfoData = async () => {

    const [infoData, err] = await Client.Do(GameUserRetained.GameUserRetainedInfo, uiData.value)

    childTableData.value = []

    if (infoData.All) {


        let data = [...infoData.List]

        debugger
        if (uiData.value.Manufacturer){
            data = data.filter(item=>GameManufacturerGroup[uiData.value.Manufacturer].indexOf(item.GameID) > -1)
        }
        if (uiData.value.GameId){
            data = data.filter(item=> uiData.value.GameId == item.GameID)
        }

        count.value = data.length

        let startIndex = (uiData.value.Page - 1) * uiData.value.PageSize
        let endIndex = startIndex + uiData.value.PageSize

        childTableData.value = data.slice(startIndex, endIndex)


        childTableData.value.forEach(item=>{
            const currentGame = initGameList.value.find(gitem=> gitem.ID == item.GameID)

            if (currentGame){

                item.GameName = currentGame.Name
                item.GameGameManufacturer = currentGame.ManufacturerName
                item.GameIcon = currentGame.IconUrl
            }


        })

    }

}

let initSearch = async () => {
    loading.value = true
    let [datas, errs] = await Client.Do(AdminGameCenter.GameList, {})
    loading.value = false
    if (errs) {
        return tip.e(errs)
    }
    initGameList.value = datas.List
    initGameList.value.forEach(item=>{

        if (!GameManufacturerGroup[item.ManufacturerName]){
            GameManufacturerGroup[item.ManufacturerName] = []
        }
        GameManufacturerGroup[item.ManufacturerName].push(item.ID)
    })

}

const pageChange = (page) => {
    uiData.value.Page = page.currentPage
    uiData.value.PageSize = page.dataSize
    initInfoData()
}

const checkInfo = async () => {
    uiData.value.ReportId = props.ReportId
    await initSearch()
    initInfoData()
}
</script>

<style scoped lang="scss">
.el-form--inline .el-form-item {
    margin-right: 0;
}
</style>
