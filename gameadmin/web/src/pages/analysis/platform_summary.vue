
<template>
    <div >
        <AnalysisPlatform :uiData='uiData' :hidden-time="true" @query-list="btnQuery" :hase-all="true"></AnalysisPlatform>

        <div class="page_table_context">
            <customTable
                v-loading="loading"
                table-name="platform_summary_list"
                :table-header="tableHeader"
                :table-data="uiData.tableDataSlice"
                :page="uiData.PageIndex"
                :page-size="uiData.PageSize"
                :count="uiData.Count"
                @refresh-table="btnQuery"
                @page-change="changePage"
            >

                <template #BetAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope.BetAmount) }}
                </template>
                <template #WinAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope.WinAmount) }}
                </template>
                <template #playersWinLose="scope">
                    <el-text :type="scope.scope.playersWinLose >= 0 ? 'success' : 'danger'" size="small">

                        {{ ut.toNumberWithComma(scope.scope.playersWinLose) }}
                    </el-text>
                </template>
                <template #companyWinsLoses="scope">
                    <el-text :type="scope.scope.companyWinsLoses >= 0 ? 'success' : 'danger'" size="small">
                        {{ ut.toNumberWithComma(scope.scope.companyWinsLoses) }}
                    </el-text>
                </template>
                <template #companyWinRate="scope">
                    <el-text :type="(scope.scope.BetAmount-scope.scope.WinAmount)/scope.scope.BetAmount >= 0 ? 'success' : 'danger'" size="small">
                        {{  percentFormatter(0,0, scope.scope.companyWinRate) }}
                    </el-text>
                </template>
            </customTable>
        </div>

    </div>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive } from 'vue';
import { useStore } from '@/pinia/index';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { useI18n } from 'vue-i18n';
import { AdminAnalysis, SearchType } from '@/api/gamepb/analysis';
import AnalysisPlatform from '@/components/analysis_platform.vue';
import customTable from '@/components/customTable/tableComponent.vue';
import moment from "moment";
const { t } = useI18n()
const store = useStore()
import ut from '@/lib/util'
import {excel} from "@/lib/excel";
let loading = ref(false)
const tableHeader = [
    {label: "序号", value: "GameID", type: "index",width: "60px"},
    {label: "商户AppID", value: "AppID",},
    {label: "玩家数量", value: "LoginPlrCount"},
    {label: "新玩家数量", value: "RegistPlrCount"},
    {label: "投注次数", value: "SpinCount"},
    {sortable:true,width: "140px",label: "投注金额", value: "BetAmount", type: "custom"},
    {sortable:true,width: "140px",label: "总赢分", value: "WinAmount", type: "custom"},
    {sortable:true,width: "140px",label: "玩家输赢", value: "playersWinLose", type: "custom"},
    {sortable:true,width: "140px",label: "公司输赢", value: "companyWinsLoses", type: "custom"},
    {sortable:true,width: "140px",label: "公司输赢占比", value: "companyWinRate", type: "custom"},
]
let uiData = reactive({
    tableData: [],
    tableDataSlice: [],
    isSortTop: false,
    isSortDown: false,
    Count: 0,
    PageIndex: 1,
    PageSize: 20,
    SearchType: SearchType.Day,
    Pid: null,
    times: "",
})
let sortBool = reactive({
    BetAmount: {
        top:false,
        down:false
    },
    WinAmount: {
        top:false,
        down:false
    },
    playersWinLose: {
        top:false,
        down:false
    },
    companyWinsLoses: {
        top:false,
        down:false
    },
})
const btnQuery = () => {
    queryList()
}
const queryList = async () => {
    console.log(uiData);
    loading.value=true
    let [data, err] = await Client.Do(AdminAnalysis.Player, {
        PageIndex: uiData.PageIndex,
        PageSize: uiData.PageSize,
        Date: !uiData.times?'':moment(uiData.times[0]).format('YYYYMMDD'),
        EndDate: !uiData.times?'':moment(uiData.times[1]).format('YYYYMMDD'),
        Operator: uiData.Pid,
        Type: "ALL",
    })
    uiData.tableDataSlice = []
    if (err) {
        return tip.e(err)
    }
    loading.value=false
    uiData.Count = data.All
    uiData.tableData = data.List?.map(t=>{
        return {
            ...t,
            playersWinLose:(t.WinAmount-t.BetAmount),
            companyWinsLoses:(t.BetAmount-t.WinAmount),
            companyWinRate: isNaN((t.BetAmount-t.WinAmount)/t.BetAmount) ? 0 : (t.BetAmount-t.WinAmount)/t.BetAmount,
        }
    }) || []

    const startIndex = (uiData.PageIndex - 1) * uiData.PageSize;
    const endIndex = startIndex + uiData.PageSize;
    uiData.tableDataSlice = uiData.tableData
}

onMounted(() => {
    queryList()
});
const changePage = (page) => {
    uiData.PageIndex = page.currentPage
    uiData.PageSize = page.dataSize
    queryList()
}
const columnSort = (name, sorts) => {
    const isAscending = sorts === 'ascending';
    const isDescending = !isAscending;
    sortBool = {
        BetAmount: {
            top:false,
            down:false
        },
        WinAmount: {
            top:false,
            down:false
        },
        playersWinLose: {
            top:false,
            down:false
        },
        companyWinsLoses: {
            top:false,
            down:false
        },
    }
    if (isAscending) {
        sortBool[name].top = !sortBool[name].top;
        sortBool[name].down = false;
    } else {
        sortBool[name].down = !sortBool[name].down
        sortBool[name].top = false;
    }

    if ((isAscending && sortBool[name].top) || (isDescending && sortBool[name].down)) {
        uiData.tableData = uiData.tableData.sort((a, b) => (isAscending ? a[name] - b[name] : b[name] - a[name]));

    } else {
        queryList();
    }
};
</script>
