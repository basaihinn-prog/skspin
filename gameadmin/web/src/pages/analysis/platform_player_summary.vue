
<template>
    <div >
        <AnalysisPlatform :uiData='uiData' @query-list="queryList" hase-all></AnalysisPlatform>
        <!-- 数据 -->

        <div class="page_table_context">
            <customTable
                v-loading="loading"
                table-name="player_summary_list"
                :table-data="tableDataSlice"
                :table-header="tableHeader"
                :count="uiData.Count"
                :page="uiData.PageIndex"
                :page-size="uiData.PageSize"
                @refresh-table="queryList"
                @page-change="pageChange"
            >

                <template #WinAmount="scope">

                        {{ ut.toNumberWithComma(scope.scope.WinAmount) }}

                </template>
                <template #BetAmount="scope">

                        {{ ut.toNumberWithComma(scope.scope.BetAmount) }}

                </template>
                <template #playersWinLose="scope">

                    <el-text :type="scope.scope.playersWinLose >= 0 ? 'success' : 'danger'">
                        {{ ut.toNumberWithComma(scope.scope.playersWinLose) }}
                    </el-text>
                </template>
                <template #companyWinsLoses="scope">

                    <el-text :type="scope.scope.companyWinsLoses >= 0 ? 'success' : 'danger'">
                        {{ ut.toNumberWithComma(scope.scope.companyWinsLoses) }}
                    </el-text>
                </template>
                <template #companyWinRate="scope">

                    <el-text :type="scope.scope.companyWinRate >= 0 ? 'success' : 'danger'">
                        {{ percentFormatter(0,0, (scope.scope.companyWinRate)) }}
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
import customTable from "@/components/customTable/tableComponent.vue";
import moment from "moment";
const { t } = useI18n()
const store = useStore()
import ut from '@/lib/util'
let loading = ref(false)



let tableHeader = [
    { label: "投注日期", value: "Date",format:(row)=> ut.fmtDate(row.Date, "YYYY-MM-DD")},
    { label: "商户AppID", value: "AppID"},
    { label: "唯一标识", value: "Pid"},
    { label: "投注次数", value: "SpinCount"},
    { label: "投注金额", value: "BetAmount", type: "custom", sortable:true},
    { label: "总赢分", value: "WinAmount", type: "custom", sortable:true},
    { label: "玩家输赢", value: "playersWinLose", type: "custom", sortable:true},
    { label: "公司输赢", value: "companyWinsLoses", type: "custom", sortable:true},
    { label: "公司输赢占比", value: "companyWinRate", type: "custom", sortable:true},

]
const tableData = ref([])
const tableDataSlice = ref([])
let uiData = reactive({
    tableData: [],
    tableDataSlice: [],
    Count: 0,
    PageIndex: 1,
    PageSize: 20,
    SearchType: SearchType.Day,
    times: "",
    Pid: null,
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
const queryList = async () => {

    loading.value=true
    let [data, err] = await Client.Do(AdminAnalysis.PlayerReport, {
        PageIndex: uiData.PageIndex,
        PageSize: uiData.PageSize,
        Date: !uiData.times?'':moment(uiData.times[0]).format('YYYYMMDD'),
        EndDate: !uiData.times?'':moment(uiData.times[1]).format('YYYYMMDD'),
        Operator: parseInt(uiData.Pid),
    })
    loading.value=false
    tableDataSlice.value = []
    if (err) {return tip.e(err)}
    uiData.Count = data.All
    tableData.value = data.List?.map(t=>{
        return {
            ...t,
            playersWinLose:(t.WinAmount-t.BetAmount),
            companyWinsLoses:(t.BetAmount-t.WinAmount),
            companyWinRate:isNaN((t.BetAmount-t.WinAmount)/t.BetAmount) ? 0 : (t.BetAmount-t.WinAmount)/t.BetAmount,
        }
    }) || []

    tableDataSlice.value = tableData.value
}

const handleCurrentChange = async (page: number) => {
    uiData.PageIndex = page
    changePage()
}
const handleSizeChange = (size: number) => {
    uiData.PageSize = size
    changePage()
}
const changePage = () => {
    const startIndex = (uiData.PageIndex - 1) * uiData.PageSize;
    const endIndex = startIndex + uiData.PageSize;
    tableDataSlice.value = tableData.value.slice(startIndex, endIndex);
}

const pageChange = (page) => {
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
        tableData.value = tableData.value.sort((a, b) => (isAscending ? a[name] - b[name] : b[name] - a[name]));
        changePage();
    } else {
        queryList();
    }
};
</script>
