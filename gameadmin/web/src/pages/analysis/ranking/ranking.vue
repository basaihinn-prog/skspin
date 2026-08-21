<template>
    <div>
      <div class="searchView gameList">
          <el-space wrap>
              <el-form
                      style="max-width: 100%"
                      @keyup.enter="initSearch"
              >
                  <el-space wrap>
                      <operator_container :defaultOperatorEvent="defaultOperatorEvent"
                                          @select-operatorInfo="operatorListChange" :is-init="true" :hase-all="true"/>

                      <el-form-item :label="$t('游戏类型')">

                          <el-select v-model="uiData.GameType" style="width: 150px" clearable>
                              <el-option v-for='(value, key) in GameClass' :label="value"
                                         :value="parseInt(key as unknown as string)"/>
                          </el-select>

                      </el-form-item>
                      <gamelist_container :defaultGameEvent="defaultGameEvent"
                                          @select-operator="selectGameList" :haseManufacturer="true" :is-init="true" :hase-all="true"></gamelist_container>
                      <el-form-item :label="$t('时间')">
                          <el-date-picker
                              v-model="timeRange"
                              type="daterange"
                              value-format="x"
                              format="YYYY-MM-DD"
                              :range-separator="$t('至')"
                              :shortcuts="shortcuts"
                              :clearable="false"
                              :start-placeholder="$t('开始时间')"
                              :end-placeholder="$t('结束时间')"
                              :disabled-date="option"
                          />

                      </el-form-item>
                      <el-form-item>


                          <el-button type="primary" @click="initSearch">{{ $t('搜索') }}</el-button>
                      </el-form-item>
                  </el-space>
              </el-form>
          </el-space>
      </div>

        <div class="page_table_context">
            <customTable
                v-loading="loading"
                table-name="gameUserRetainedList"
                :table-header="tableHeader"
                :table-data="tableData"
                :page-size="uiData.PageSize"
                :page="uiData.Page"
                :count="count"
                @refresh-table="initSearch"
                @pageChange="PageChange"
            >

                <template #handleTools>
                    <el-radio-group v-model="uiData.RankingListType" @change="initSearch">
                        <el-radio-button :label="$t('赢钱榜')" :value="0"/>
                        <el-radio-button :label="$t('输钱榜')" :value="1"/>
                        <el-radio-button :label="$t('投注金额')" :value="2"/>
                        <el-radio-button :label="$t('胜率榜')" :value="3"/>
                    </el-radio-group>
                </template>
                <template #TotalWinLoss="scope">

                    <el-text :type="scope.scope.TotalWinLoss >= 0 ? 'success' : 'danger'">
                        {{ ut.toNumberWithComma(scope.scope.TotalWinLoss) }}
                    </el-text>
                </template>
                <template #HistoryWinLoss="scope">

                    <el-text :type="scope.scope.HistoryWinLoss >= 0 ? 'success' : 'danger'">
                        {{ ut.toNumberWithComma(scope.scope.HistoryWinLoss) }}
                    </el-text>
                </template>
                <template #Balance="scope">


                    <el-text type="info">
                        {{ ut.toNumberWithComma(scope.scope.Balance) }}
                    </el-text>
                </template>
                <template #TotalBet="scope">

                    <el-text type="info">
                        {{ ut.toNumberWithComma(scope.scope.TotalBet) }}
                    </el-text>
                </template>
                <template #WinRate="scope">

                    <el-text type="info">
                        {{ `${scope.scope.WinRate.toFixed(2)}%` }}
                    </el-text>
                </template>
                <template #LoginAt="scope">

                    {{ ut.fmtSelectedUTCDateFormat(new Date(scope.scope.LoginAt).getTime()) }}

                </template>
                <template #CreateAt="scope">

                    {{ ut.fmtSelectedUTCDateFormat(new Date(scope.scope.CreateAt).getTime()) }}

                </template>
            </customTable>
        </div>






    </div>
</template>

<script setup lang="ts">

import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import {onMounted, Reactive, reactive, ref, watch} from "vue";
import {useI18n} from "vue-i18n";
import {Client} from "@/lib/client";
import {GameUserRetained, GameUserRetainedReq, GetRankingListParam} from "@/api/adminpb/gameUserRetained";
import ut from "@/lib/util";
import {useStore} from "@/pinia";
import Gamelist_container from "@/components/gamelist_container.vue";
import {GameType} from "@/api/gamepb/customer";
import {GameClass} from "@/enum";
const {t} = useI18n()
const store = useStore()

const ReportId = ref("")
const Operator = ref(0)
const count = ref(0)
const loading = ref(false)
const dialog = ref(false)
const timeRange = ref([new Date().getTime(), new Date().getTime()])

watch(dialog, (newData)=>{

    if (!newData){
        ReportId.value = ""
    }
})



const option = (value) => {


    return value > new Date()
}

const formatTime = (time) => {
    let timeStr = ut.fmtDate(new Date(time).getTime()/1000)
    return ut.fmtDate(timeStr, "YYYY-MM-DD")
}

// 留存数据列表头
const tableUpHeader = [
    {label: "排名", value: "GameID", type: "index",width: "60px"},
    {label: "商户AppID", value: "AppID"},
    {label: "商户币种", value: "CurrencyName"},
    {label: "唯一标识", value: "Pid"},
]
const tableDownHeader = [
    {label: "累计输赢", value: "TotalWinLoss", type: "custom"},
    {label: "历史输赢", value: "HistoryWinLoss", type: "custom"},
    {label: "余额", value: "Balance", type: "custom"},
    {label: "投注次数", value: "SpinCount"},
    {label: "玩家注册时间", value: "CreateAt", type: "custom", width: "180px"},
    {label: "最后登录时间", value: "LoginAt", type: "custom", width: "180px"},
]




const tableHeader = ref([])

// 留存数据列表
const tableData = ref([])

// 时间查询
const shortcuts = [
    {
        text: t('过去7天'),
        value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            return [Date.parse(start.toString()), Date.parse(end.toString())]
        },
    },
    {
        text: t('今天'),
        value: () => {
            const today = new Date();
            const year = today.getFullYear();
            const month = today.getMonth();
            const date = today.getDate();

            const startTime = new Date(year, month, date, 0, 0, 0, 0o00).getTime();
            const endTime = new Date(year, month, date, 23, 59, 59, 0o00).getTime();

            return [startTime, endTime];
        },
    },
]

const uiData:Reactive<GetRankingListParam> = reactive<GetRankingListParam>({
    OperatorId: null,
    GameID: null,
    Manufacturer: "",
    RankingListType: 0,
    StartTime: 0,
    EndTime: 0,
    PageSize:20,
    Page:1,
    GameType: 0
})

onMounted(()=>{
    initSearch()
})

const defaultGameEvent = ref()
const selectGameList = (value) => {

    if (value){

        uiData.Manufacturer = value.manufacturer


        uiData.GameID = value.gameData
    }
}
// GameID
// Manufacturer

const defaultOperatorEvent = ref()
const operatorListChange = (value) => {
    if (value){

        uiData.OperatorId = value.Id == "ALL" ? null : value.Id
    }else{

        uiData.OperatorId = null
    }
}


const initSearch = async () => {


    if (uiData.RankingListType<3){

        tableHeader.value = [
            ...tableUpHeader,
            {label: "投注金额", value: "TotalBet", type:"custom"},
            ...tableDownHeader
        ]
    }else{
        tableHeader.value = [
            ...tableUpHeader,
            {label: "胜率", value: "WinRate", type:"custom"},
            ...tableDownHeader
        ]
    }

    const requestData = {
        ...uiData,
    }


    const startTime = new Date(timeRange.value[0]);
    startTime.setHours(0)
    startTime.setMinutes(0)
    startTime.setSeconds(0)
    startTime.setMilliseconds(0)

    const endTime = new Date(timeRange.value[1]);
    endTime.setHours(23)
    endTime.setMinutes(59)
    endTime.setSeconds(59)
    endTime.setMilliseconds(0)

    requestData.StartTime = ut.fmtUTCDate(startTime.getTime())
    requestData.EndTime = ut.fmtUTCDate(endTime.getTime())
    if(requestData.Manufacturer){

        requestData.Manufacturer = requestData.Manufacturer.toLowerCase()
    }

    tableData.value = []
    loading.value = true
    let [data, err] = await Client.Do(GameUserRetained.GetRankingList, requestData);
    loading.value = false



    if (data.Count) {
        tableData.value = data.List
    }
    count.value = data.Count
}

const PageChange = (page) => {
    uiData.Page = page.currentPage
    uiData.PageSize = page.dataSize
    initSearch()
}

const checkInfo = (row) => {
    ReportId.value = row.ID
    dialog.value = true
}


</script>

<style scoped lang="scss">

</style>
