
<template>
    <div >
        <div class="searchView">
            <el-form
                :model="uiData"
                style="max-width: 100%"
            >
                <el-space wrap>

                    <operator_container :defaultOperatorEvent ="defaultOperatorEvent" @select-operatorInfo="operatorListChange" :hase-all="true"></operator_container>
                        <gamelist_container
                            :hase-manufacturer="true"
                            :hase-all="true"
                            :defaultGameEvent="defaultGameEvent"
                            :isInit="true"
                            @select-operator="selectGameList"/>

                    <el-form-item :label="$t('日期')">
                        <el-date-picker
                            v-model="uiData.Date"
                            type="date"
                            :placeholder="$t('日期')" value-format="YYYYMMDD"
                        />
                    </el-form-item>

                    <el-form-item>
                        <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
                    </el-form-item>
                </el-space>
            </el-form>

        </div>


        <div class="page_table_context">
        <customTable
            table-name="appAnalysis_list"
            :table-header="tableHeader"
            :table-data="uiData.tableData" v-loading="loading" :page="uiData.PageIndex" :page-size="uiData.PageSize" :count="uiData.Count"
            @refresh-table="queryList"
            @page-change="pageChange">

            <template #Game="scope">

                <div style="display: flex;align-items: center;">
                    <el-tag size="small" type="primary">{{ scope.scope.GameManufacturerName }}</el-tag>

                    <div style="margin: 0 10px">
                    <el-avatar :src="scope.scope.GameIcon" size="small" />
                    </div>
                    {{ scope.scope.Game }}
                </div>
            </template>
            <template #BetAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope.BetAmount) }}

            </template>

            <template #WinAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope.WinAmount) }}

            </template>

            <template #shouyi="scope">
                <el-text :type="scope.scope.shouyi >= 0 ? 'success' : 'danger'">
                    {{ ut.toNumberWithCommaNormal(ut.fmtGold(scope.scope.shouyi)) }}
                </el-text>
            </template>

            <template #huibao="scope">


                {{ scope.scope.BetAmount <= 0 ? '∞' : percentFormatter(0,0, scope.scope.huibao) }}


            </template>
        </customTable>
        </div>
    </div>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive } from 'vue';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { AdminAnalysis, SearchType } from '@/api/gamepb/analysis';
import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import {useStore} from "@/pinia";
import ut from "@/lib/util";
const { t } = useI18n()
import Currency_container from "@/components/currency_container.vue";
import {useI18n} from "vue-i18n";
import Gamelist_container from "@/components/gamelist_container.vue";
import {excel} from "@/lib/excel";
const defaultOperatorEvent = ref({})
const defaultGameEvent = ref({})
const defaultCurrencyEvent = ref({})



const store = useStore();
let loading = ref(false)
const sortOrder = ref('ascending');


let tableHeader = [
    {label: "日期", value: "Date", format:(row)=> ut.fmtDate(row.Date, "YYYY-MM-DD")},
    {label: "游戏名称", value: "Game", width:"250px", type: "custom"},
    {label: "商户AppID", value: "AppID"},
    {label: "币种", value: "CurrencyName"},
    {label: "投注金额", value: "BetAmount", type: "custom", sortable: true},
    {label: "总赢分", value: "WinAmount", type: "custom", sortable: true},
    {label: "收益", value: "shouyi", type: "custom", sortable: true},
    {label: "回报率", value: "huibao", type: "custom", sortable:true},
]

let uiData = reactive({
    tableData: [],
    tableDataSlice: [],
    Count: 0,
    PageIndex: 1,
    PageSize: 20,
    Manufacturer: "",
    Date: '',
    Game: '',
    Operator: 0,
    CurrencyCode: "",
    NeedLast: false,
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
    shouyi: {
        top:false,
        down:false
    },
})
const queryList = async () => {

    loading.value=true
    uiData.EndDate = uiData.Date
    let [data, err] = await Client.Do(AdminAnalysis.GameAnalysisList, uiData)
    loading.value=false
    uiData.tableData = []
    if (err) {
        return tip.e(err)
    }
    uiData.Count = data.All
    if (data.List){
        uiData.tableData = data.List.map(t=>{
            return {
                ...t,
                shouyi:t.BetAmount - t.WinAmount,
                huibao:isNaN(t.WinAmount/ t.BetAmount) ? 0 : (t.WinAmount/ t.BetAmount),
            }
        }) || []
    }


    const startIndex = (uiData.PageIndex - 1) * uiData.PageSize;
    const endIndex = startIndex + uiData.PageSize;
    uiData.tableDataSlice = uiData.tableData.slice(startIndex, endIndex);
}

onMounted(() => {

});
const CurrencyCode = ref("")
const operatorListChange = (value) =>{

    if (value){

        uiData.Operator = value.Id
        CurrencyCode.value = value.CurrencyKey
    } else{
        uiData.Operator = null
        CurrencyCode.value = ""
    }
}
const selectGameList = (value) =>{

    if (value.gameData){

    }
    uiData.Game = value.gameData

    if (value.manufacturer || value.manufacturer == null){

        uiData.Manufacturer = value.manufacturer
    }


}
const currencyListChange = (value) =>{
    uiData.CurrencyCode = value.CurrencyCode
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

const pageChange = (page) => {
    uiData.PageSize = page.dataSize
    uiData.PageIndex = page.currentPage
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
        shouyi: {
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
const generatorExcel = () => {

    const excelTableHeader = [
        {label: "游戏厂商", value: "GameManufacturerName"},
        {label: "游戏图标", value: "GameIcon"},
    ]

    excel.DataGeneratorExcel([...tableHeader, ...excelTableHeader], uiData.tableDataSlice, `产品数据`)
}
</script>
