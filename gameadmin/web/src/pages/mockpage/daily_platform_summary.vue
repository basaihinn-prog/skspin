
<template>
    <div >
        <div class="searchView">
            <AnalysisPlatform :uiData='uiData' @query-list="queryList" :hase-all="true"></AnalysisPlatform>
        </div>
        <!-- 数据 -->
        <div class="page_table_context">
            <customTable
                v-loading="loading"
                table-name="daily_summary_list"
                :table-header="tableHeader"
                :table-data="uiData.tableDataSlice"
                :page="uiData.PageIndex"
                :page-size="uiData.PageSize"
                :count="uiData.Count"
                @refresh-table="queryList"
                @page-change="changePage"
            >

                <template #handleTools="scope">
                    <el-space>
                        <upload-excel @uploadFile="resolveExcelData" file-name="每日平台汇总"/>

                        <el-button type="primary" plain @click="generatorExcel">生成Excel</el-button>

                    </el-space>
                </template>

                <template #BetAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope.BetAmount * 1e4) }}
                </template>

                <template #WinAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope.WinAmount * 1e4) }}
                </template>

                <template #playersWinLose="scope">
                    <el-text :type="scope.scope.playersWinLose < 0 ? 'danger' : 'success'" size="small">{{ ut.toNumberWithComma(scope.scope.playersWinLose * 1e4) }}</el-text>

                </template>

                <template #companyWinsLoses="scope">

                    <el-text :type="scope.scope.companyWinsLoses < 0 ? 'danger' : 'success'" size="small">{{ ut.toNumberWithComma(scope.scope.companyWinsLoses * 1e4) }}</el-text>

                </template>
                <template #companyProportion="scope">
                    <el-text :type="scope.scope.companyProportion < 0 ? 'danger' : 'success'" size="small">{{ percentFormatter(0,0, scope.scope.companyProportion )}}</el-text>
                </template>
                <template #GGR="scope">
                    <el-text :type="scope.scope.GGR < 0 ? 'danger' : 'success'" size="small">{{ ut.toNumberWithComma(scope.scope.GGR * 1e4) }}</el-text>

                </template>
                <template #OperatorBalance="scope">
                    <div style="display: flex;align-items: center;justify-content: center">
                        <template v-if="scope.scope.OperatorBalance>0">
                            <el-text v-if="operatorBalance[scope.scope.AppID] < scope.scope.OperatorBalance" type='success' size="small"><span style="float: left" v-html="scope.scope.CurrencySymbol"></span>&nbsp;{{ ut.toNumberWithComma(scope.scope.OperatorBalance * 1e4) }}</el-text>
                            <el-text v-if="operatorBalance[scope.scope.AppID] >= scope.scope.OperatorBalance" type='primary' size="small"><span style="float: left" v-html="scope.scope.CurrencySymbol"></span>&nbsp;{{ ut.toNumberWithComma(scope.scope.OperatorBalance * 1e4) }}</el-text>

                        </template>
                        <template v-else>
                            <el-text v-if="scope.scope.OperatorBalance <= 0" type='danger' size="small"><span style="float: left" v-html="scope.scope.CurrencySymbol"></span>&nbsp;{{ ut.toNumberWithComma(scope.scope.OperatorBalance * 1e4) }}</el-text>

                        </template>
                    </div>

                </template>
            </customTable>
        </div>




    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, nextTick} from 'vue';
import { useStore } from '@/pinia/index';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { useI18n } from 'vue-i18n';
import { AdminAnalysis, SearchType } from '@/api/gamepb/analysis';
import AnalysisPlatform from '@/components/analysis_platform.vue';
import customTable from "@/components/customTable/tableComponent.vue"
import moment from "moment";
const { t } = useI18n()
const store = useStore()
import ut from "@/lib/util";
import {AdminInfo} from "@/api/adminpb/info";
import UploadExcel from "./component/uploadExcel.vue";
import {excel} from "@/lib/excel";
let loading = ref(false)
let uploadStatus = ref(false)
let operatorBalance = ref({})
let OperatorConfig = ref(null)
let tableHeaderMap = ref(null)
let tableHeader = [
    { label: "投注日期", value: "Date", format:(row)=> ut.fmtDate(row.Date, "YYYY-MM-DD")},
    { label: "商户AppID", value: "AppID"},
    { label: "玩家数量", value: "LoginPlrCount"},
    { label: "新玩家数量", value: "RegistPlrCount"},
    { label: "投注次数", value: "SpinCount"},
    { label: "投注金额", value: "BetAmount", type: "custom", sortable:true},
    { label: "总赢分", value: "WinAmount", type: "custom", sortable:true},
    { label: "玩家输赢", value: "playersWinLose", type: "custom", sortable:true},
    { label: "公司输赢", value: "companyWinsLoses", type: "custom", sortable:true},
    { label: "公司输赢占比", value: "companyProportion", type: "custom", sortable:true, width: "150px"},
    { label: "平台费", value: "GGR", type: "custom", sortable:true},
    { label: "币种", value: "CurrencyName"},
    { label: "合作模式", value: "CooperationType"},
    { label: "余额", value: "OperatorBalance", type: "custom", width: "150px"},
]
let tableData = ref([])
let uiData = reactive({
    tableData: [],
    tableDataSlice: [],
    Count: 0,
    PageIndex: 1,
    PageSize: 20,
    SearchType: SearchType.Day,
    Pid: null,
    times: "",
    rate: 0
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



uiData.times = [new Date(),new Date()]

const queryList = async () => {
    // if (!uiData.Pid && store.AdminInfo.GroupId != 3){
    //     return tip.e(t("商户不能为空"))
    // }
    loading.value=true
    let [data, err] = await Client.Do(AdminAnalysis.Player, {
        PageIndex: uiData.PageIndex,
        PageSize: uiData.PageSize,
        Operator: uiData.Pid,
        Date: !uiData.times?'':moment(uiData.times[0]).format('YYYYMMDD'),
        EndDate: !uiData.times?'':moment(uiData.times[1]).format('YYYYMMDD'),
        Type: '',
    })
    loading.value=false
    if (err) {
        return tip.e(err)
    }
    uiData.tableData = []
    uiData.Count = data.All


    uiData.tableData = data.List?.map(m=>{

        return {
            ...m,
            playersWinLose:(m.WinAmount-m.BetAmount),
            companyWinsLoses:(m.BetAmount-m.WinAmount),
            companyProportion: isNaN((m.BetAmount-m.WinAmount)/m.BetAmount) ? 0 : (m.BetAmount-m.WinAmount)/m.BetAmount,
            CooperationType:["", t("收益分成"), t("流水分成")][OperatorConfig.value[m.AppID].CooperationType]
        }
    }) || []
    const startIndex = (uiData.PageIndex - 1) * uiData.PageSize;
    const endIndex = startIndex + uiData.PageSize;
    uiData.tableDataSlice = uiData.tableData

}
const getOperatorData = async () => {
    let [data, err] = await Client.Do(AdminInfo.GetOperatorList, {
        PageIndex: 1,
        PageSize: 10000,
        OperatorType: 2,
        Status: -1
    })
    if (err) {
        return tip.e(err)
    }
    OperatorConfig.value = {}
    for (const i in data.List) {
        operatorBalance.value[data.List[i].AppID] = data.List[i].BalanceThreshold
        OperatorConfig.value[data.List[i].AppID] = data.List[i]
    }


}

getOperatorData()

const changePage = (page) => {
    uiData.PageIndex = page.currentPage
    uiData.PageSize = page.dataSize
    if (uploadStatus.value){
        uploadStatus.value = false
        return
    }
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
    }
};


const generatorExcel = (data) => {
    let generatorHeader = [...tableHeader].filter(item => item.label != '操作')


    excel.DataGeneratorExcel(generatorHeader, uiData.tableDataSlice, `每日平台汇总`)
}


const resolveExcelData = (data) => {
    uploadStatus.value = true
    if (!tableHeaderMap.value) {
        tableHeaderMap.value = {}
        for (const i in tableHeader) {

            tableHeaderMap.value[tableHeader[i].label] = tableHeader[i].value
        }
    }
    uiData.Count = data.data.length
    let renderData = data.data.slice(0, uiData.PageSize)
    uiData.PageIndex = 1
    nextTick(()=>{

        uiData.tableData = excel.excelResolveData(tableHeaderMap.value, renderData)
        uiData.tableDataSlice = excel.excelResolveData(tableHeaderMap.value, renderData);
        uploadStatus.value = false
    })

}


</script>
<style scoped>
.searchView{
    padding-top: 10px;
    padding-left: 10px;
    padding-bottom: 10px;
}
</style>
