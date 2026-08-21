<template>
    <div >
        <div class="searchView ">
            <el-form
                :model="uiData"
                style="max-width: 100%"
            >
                <el-space wrap >
                    <operator_container :defaultOperatorEvent ="defaultOperatorEvent" @select-operatorInfo="operatorListChange" :hase-all="true"></operator_container>
                    <el-form-item :label="$t('日期查询')">
                        <el-date-picker type="month" :placeholder="$t('请选择')" v-model="uiData.times" size="default"/>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
                    </el-form-item>
                </el-space>
            </el-form>
        </div>

        <div class="page_table_context">
            <customTable
                v-loading="loading"
                table-name="monthly_summary_list"
                :table-header="tableHeader"
                :table-data="uiData.tableDataSlice"
                :page="uiData.PageIndex"
                :count="uiData.Count"
                :page-size="uiData.PageSize"
                @page-change="pageChange"
                @refresh-table="queryList"

            >


                <template #BetAmount="scope">

                    {{ ut.toNumberWithComma(scope.scope.BetAmount) }}
                </template>
                <template #WinAmount="scope">

                    {{ ut.toNumberWithComma(scope.scope.WinAmount) }}
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
                    {{ percentFormatter(0, 0, scope.scope.companyWinRate)}}
                    </el-text>
                </template>
                <template #GGR="scope">
                    <el-text :type="scope.scope.GGR >= 0 ? 'success' : 'danger'">
                        {{ ut.toNumberWithComma(scope.scope.GGR) }}
                    </el-text>

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
import { onMounted, ref, reactive } from 'vue';
import { useStore } from '@/pinia/index';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { useI18n } from 'vue-i18n';
import { AdminAnalysis, SearchType } from '@/api/gamepb/analysis';
import AnalysisPlatform from '@/components/analysis_platform.vue';
import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import ut from '@/lib/util'
import moment from 'moment';
import {AdminInfo} from "@/api/adminpb/info";
import {excel} from "@/lib/excel";
const defaultOperatorEvent = ref({})


const { t } = useI18n()
const store = useStore()


const Rate = ref(0)
let loading = ref(false)
let OperatorConfig = ref(null)
let operatorBalance = ref({})
let tableHeader = [
    { label: "投注月份", value: "Date",format:(row)=> ut.fmtDate(row.Date, "YYYY-MM")},
    { label: "商户AppID", value: "AppID"},

    { label: "投注次数", value: "SpinCount", sortable:true},
    { label: "投注金额", value: "BetAmount", type: "custom", sortable:true},
    { label: "总赢分", value: "WinAmount", type: "custom", sortable:true},
    { label: "玩家输赢", value: "playersWinLose", type: "custom", sortable:true},
    { label: "公司输赢", value: "companyWinsLoses", type: "custom", sortable:true},
    { label: "公司输赢占比", value: "companyWinRate", type: "custom", sortable:true, width:"160px"},
    { label: "合作模式", value: "CooperationType"},
    { label: "平台费", value: "GGR", type: "custom", sortable:true, width:"180px"},
    { label: "币种", value: "CurrencyName"},
    { label: "余额", value: "OperatorBalance", type: "custom", width:"150px"},

]
let uiData = reactive({
    tableData: [],
    tableDataSlice: [],
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
const operatorListChange = (value) =>{

    if (value){

        uiData.Pid = value.Id
        Rate.value = value.PresentRate
    }else{

        uiData.Pid = ""
        Rate.value = 0
    }
}

uiData.times = new Date()
const queryList = async () => {

    // if (!uiData.Pid && store.AdminInfo.GroupId != 3){
    //     return tip.e("请选择商户")
    // }

    loading.value=true
    let range = moment(uiData.times).format('YYYYMM')
    uiData.tableData =[]
        let timesStart = range +'00'
    let timesEnd = range +'31'
    let [data, err] = await Client.Do(AdminAnalysis.Player, {
        PageIndex: uiData.PageIndex,
        PageSize: uiData.PageSize,
        Operator: uiData.Pid ? uiData.Pid : null,
        Date: !uiData.times?'':timesStart,
        EndDate: !uiData.times?'':timesEnd,
        Type: 'Month',
    })
    loading.value=false
    if (err) {
        return tip.e(err)
    }

    uiData.Count = data.All
    uiData.tableData = data.List?.map(m=>{

        return {
            ...m,
            playersWinLose:(m.WinAmount-m.BetAmount),
            companyWinsLoses:(m.BetAmount-m.WinAmount),
            companyWinRate:isNaN((m.BetAmount-m.WinAmount)/m.BetAmount) ? 0 : (m.BetAmount-m.WinAmount)/m.BetAmount,
            CooperationType:["", t("收益分成"), t("流水分成")][OperatorConfig.value[m.AppID].CooperationType]
        }
    }) || []
    const startIndex = (uiData.PageIndex - 1) * uiData.PageSize;
    const endIndex = startIndex + uiData.PageSize;
    uiData.tableDataSlice = uiData.tableData.slice(startIndex, endIndex);
}


const changePage = () => {
    const startIndex = (uiData.PageIndex - 1) * uiData.PageSize;
    const endIndex = startIndex + uiData.PageSize;
    uiData.tableDataSlice = uiData.tableData.slice(startIndex, endIndex);
}
const handleCurrentChange = async (page: number) => {
    uiData.PageIndex = page
    changePage()
}
const handleSizeChange = (size: number) => {
    uiData.PageSize = size
    changePage()
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
        changePage();
    } else {
        queryList();
    }
};


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
    // operatorBalance.value['data']

}
getOperatorData()
const pageChange = (page) => {
    uiData.PageSize = page.dataSize
    uiData.PageIndex = page.currentPage
    queryList()
}

</script>
