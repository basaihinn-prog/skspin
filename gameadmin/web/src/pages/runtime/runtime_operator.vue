<template>
    <div v-loading="pageLoading">
        <div class="searchView gameList">
            <el-form
                style="max-width: 100%"
            >
                <el-space wrap>

                    <operator_container  :is-init="true"
                                         :defaultGameEvent="defaultOperatorEvent"
                                         @select-operatorInfo="operatorListChange"/>
                    <gamelist_container
                            :hase-manufacturer="true"
                            :is-init="true"
                            :defaultGameEvent="defaultGameEvent"
                            @select-operator="selectGameList"/>
                </el-space>
            </el-form>
            <el-space wrap>
                <el-button type="primary" @click="initSearch">{{ $t('搜索') }}</el-button>
            </el-space>
        </div>
        <el-radio-group v-model="changeShow" @change="radioChange">
            <el-radio-button :label="$t('图表')" value="图表"/>
            <el-radio-button :label="$t('表格')" value="表格"/>
        </el-radio-group>
        <div v-show="changeShow === '图表' && echartShow">
            <el-row class="top-chart">
                <el-col :md="8" :span="24" class="top-chart-col">
                    <el-progress
                        :percentage="uiData.SevenDays.TotalReturnRate > 100 ? 100 : uiData.SevenDays.TotalReturnRate"
                        class="7days-chart" stroke-width="3"
                        indeterminate type="dashboard" color="steelblue" striped striped-flow>
                        <template #default="{ percentage }">
                            <el-row class="top-chart" :gutter="10">
                                <el-col :span="12">
                                    <el-statistic :title="$t('7日总流水')" :value="uiData.SevenDays.BetAmount"/>
                                </el-col>
                                <el-col :span="12">
                                    <el-statistic :title="$t('7日总产出')" :value="uiData.SevenDays.WinAmount"/>
                                </el-col>
                                <el-col :span="12">
                                    <el-statistic :title="$t('7日总回报率')" :value="uiData.SevenDays.TotalReturnRate"/>
                                </el-col>
                                <el-col :span="12">
                                    <el-statistic :title="$t('系统收益')" :value="uiData.SevenDays.SystemBenefits"/>
                                </el-col>
                            </el-row>

                        </template>
                    </el-progress>
                </el-col>
               <el-col :md="8" :span="24" class="top-chart-col">
                    <el-progress class="7days-chart" stroke-width="3"
                                 :percentage="uiData.MonthData.TotalReturnRate > 100 ? 100 : uiData.MonthData.TotalReturnRate"
                                 indeterminate type="dashboard" color="steelblue" striped striped-flow>
                        <template #default="{ percentage }">
                            <el-row class="top-chart" :gutter="10">
                                <el-col :span="12">
                                    <el-statistic :title="$t('本月总流水')" :value="uiData.MonthData.BetAmount"/>
                                </el-col>
                                <el-col :span="12">
                                    <el-statistic :title="$t('本月总产出')" :value="uiData.MonthData.WinAmount"/>
                                </el-col>
                                <el-col :span="12">
                                    <el-statistic :title="$t('本月总回报率')" :value="uiData.MonthData.TotalReturnRate"/>
                                </el-col>
                                <el-col :span="12">
                                    <el-statistic :title="$t('系统收益')" :value="uiData.MonthData.SystemBenefits"/>
                                </el-col>
                            </el-row>
                        </template>
                    </el-progress>
                </el-col>
               <el-col :md="8" :span="24" class="top-chart-col">
                    <el-progress stroke-width="3"
                                 :percentage="uiData.TotalData.TotalReturnRate > 100 ? 100 : uiData.TotalData.TotalReturnRate"
                                 class="7days-chart"
                                 indeterminate type="dashboard" color="steelblue" striped striped-flow>
                        <template #default="{ percentage }">
                            <el-row class="top-chart" :gutter="10">
                                <el-col :span="12">
                                    <el-statistic :title="$t('总流水')" :value="uiData.TotalData.BetAmount"/>
                                </el-col>
                                <el-col :span="12">
                                    <el-statistic :title="$t('总产出')" :value="uiData.TotalData.WinAmount"/>
                                </el-col>
                                <el-col :span="12">
                                    <el-statistic :title="$t('总回报率')" :value="uiData.TotalData.TotalReturnRate"/>
                                </el-col>
                                <el-col :span="12">
                                    <el-statistic :title="$t('系统收益')" :value="uiData.TotalData.SystemBenefits"/>
                                </el-col>
                            </el-row>

                        </template>
                    </el-progress>
                </el-col>
            </el-row>
            <el-row :gutter="20" class="bottom-chart">
               <el-col :md="8" :span="24">
                    <div class="bottom-chart-col">
                    <div id="DAUChart" ref="7days-chart" style="width: 100%; height: 400px;" class="DAU-chart"></div>
                    </div>
                </el-col>
               <el-col :md="8" :span="24">
                    <div class="bottom-chart-col">
                    <div id="totalPressurePointsChart" ref="total-pressure-points-chart"
                         style="width: 100%; height: 400px;" class="total-pressure-points-chart"></div>
                    </div>
                </el-col>
               <el-col :md="8" :span="24">
                       <div class="bottom-chart-col">
                    <div id="systemBenefitChart" ref="system-benefit-chart" style="width: 100%; height: 400px;"
                         class="system-benefit-chart"></div>
                       </div>
                </el-col>

                <el-col :md="8" :span="24">
                    <div class="bottom-chart-col">
                        <div id="responseRateChart" ref="response-rate-chart" style="width: 100%; height: 400px;"
                             class="response-rate-chart"></div>
                    </div>
                </el-col>
               <el-col :md="8" :span="24">
                       <div class="bottom-chart-col">
                    <div id="buyTotalBetChart" ref="buy-total-bet-chart" style="width: 100%; height: 400px;"
                         class="buy-total-bet-chart"></div>
                       </div>
                </el-col>
               <el-col :md="8" :span="24">
                       <div class="bottom-chart-col">
                    <div id="buyResponseRateChart" ref="buy-response-rate-chart" style="width: 100%; height: 400px;"
                         class="buy-response-rate-chart"></div>
                       </div>
                </el-col>
            </el-row>
        </div>
        <div v-show="changeShow === '表格'">
            <runtime-table :gameMsg="selectGameListVal?.Game"></runtime-table>

            <customtable :tableHeader="selectGameListVal?.Game?.startsWith('jili_') ?  JILITableHeader : tableHeader" :table-data="uiData.tableData" style="width: 100%" v-loading="loading">

                    <template #DAU="scope">{{ scope.scope.DauValue }}</template>
                    <template #BetValue="scope">{{ ut.toNumberWithComma(scope.scope.BetValue) }}</template>

                    <template #WinValue="scope">{{ ut.toNumberWithComma(scope.scope.WinValue) }}</template>

                    <template #SystemValue="scope">{{ ut.toNumberWithComma(scope.scope.SystemValue) }}</template>

                    <template #RateValue="scope">{{ percentFormatter(0, 0, scope.scope.RateValue) }}</template>

                    <template #BuyBet="scope">{{ ut.toNumberWithComma(scope.scope.BuyBet) }}</template>

                    <template #BuyWin="scope">{{ ut.toNumberWithComma(scope.scope.BuyWin) }}</template>
                    <template #BuySystemValue="scope">{{ ut.toNumberWithComma(scope.scope.BuySystemValue) }}</template>
                   <template #BuyRateValue="scope">{{ percentFormatter(0, 0, scope.scope.BuyRateValue) }}</template>

                      <template #ExtraBet="scope">{{ ut.toNumberWithComma(scope.scope.ExtraBet) }}</template>
                        <template #ExtraWin="scope">{{ ut.toNumberWithComma(scope.scope.ExtraWin) }}</template>
                        <template #ExtraSystemValue="scope">{{ ut.toNumberWithComma(scope.scope.ExtraSystemValue) }}</template>
                        <template #ExtraRateValue="scope">{{ percentFormatter(0, 0, scope.scope.ExtraRateValue) }}</template>

            </customtable>
        </div>
    </div>
</template>

<script lang="ts" setup>
import {onMounted, ref, reactive, onBeforeUnmount, nextTick, shallowRef, onUnmounted} from 'vue';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {AdminStatsRpc} from '@/api/stats/stats';
import ut, {initGolbal} from '@/lib/util'
import * as echarts from 'echarts/core';
import {
    TitleComponent,
    TooltipComponent,
    GridComponent,
    LegendComponent,
    MarkLineComponent,
} from 'echarts/components';
import {LineChart} from 'echarts/charts';
import {UniversalTransition} from 'echarts/features';
import {CanvasRenderer} from 'echarts/renderers';
import {useI18n} from 'vue-i18n';
import {excel} from "@/lib/excel";
import {markRaw} from 'vue'
import Gamelist_container from "@/components/gamelist_container.vue";
import {useStore} from "@/pinia";
import customtable from "@/components/customTable/tableComponent.vue"

const DAUChart = shallowRef(null);
const totalPressurePointsChart = shallowRef(null);
const systemBenefitChart = shallowRef(null);
const buyTotalBetChart = shallowRef(null);
const responseRateChart = shallowRef(null);
const buyResponseRateChart = shallowRef(null);
const {t} = useI18n()
const changeShow = ref('图表')
const echartShow = ref(false)
const piniaStore = useStore()
const {setRuntimeTable} = piniaStore

const pageLoading = ref(false)
import RuntimeTable from "@/pages/runtime/runtimeTable.vue";
import Operator_container from "@/components/operator_container.vue";
import tr from "element-plus/es/locale/lang/TR";
let uiData = reactive({
    tableData: [],
    allTableData: null,
    gameName: "",
    TotalExpense: 0,
    TotalIncome: 0,
    SevenDays: {
        TotalReturnRate: '0.00%', // 7日总回报率
        ExtraBetAmount: 0, // 购买小游戏流水
        ExtraWinAmount: 0, // 购买小游戏产出
        extra: '0.00%', // 7日额外模式回报率
        BuyBetAmount: 0, // 额外模式流水
        BuyWinAmount: 0, // 额外模式产出
        BetAmount: 0, // 7日总流水
        WinAmount: 0, // 7日总产出
        SystemBenefits: 0, // 系统受益
        percentage: 0, // 系统受益
    },
    MonthData: {
        TotalReturnRate: '0.00%', // 本月总回报率
        ExtraBetAmount: 0, // 购买小游戏流水
        ExtraWinAmount: 0, // 购买小游戏产出
        extra: '0.00%', // 本月额外模式回报率
        BuyBetAmount: 0, // 额外模式流水
        BuyWinAmount: 0, // 额外模式产出
        BetAmount: 0, // 本月总流水
        WinAmount: 0, // 本月总产出
        SystemBenefits: 0, // 系统受益
        percentage: 0, // 系统受益
    },
    TotalData: {
        TotalReturnRate: '0.00%', // 总回报率
        ExtraBetAmount: 0, // 购买小游戏流水
        ExtraWinAmount: 0, // 购买小游戏产出
        extra: '0.00%', // 总额外模式回报率
        BuyBetAmount: 0, // 额外模式流水
        BuyWinAmount: 0, // 额外模式产出
        BetAmount: 0, // 总流水
        WinAmount: 0, // 总产出
        SystemBenefits: 0, // 系统受益
        percentage: 0, // 系统受益
    },
    TotalMonthExpense: 0,
    TotalMonthIncome: 0,
    showColumn: {},
    chartDAU: [],
    chartTotalPressurePoints: [],
    chartSystemBenefit: [],
    chartBuyTotalBet: [],
    chartResponseRate: [],
    chartBuyResponseRate: [],
})
let init = ref(true)
let UserName = ref('')
const value = ref([])

const tableHeader = ref([
    {label: "日期", value: "CurrencyName", format:(row)=>ut.fmtDate(row.Time, 'YYYY-MM-DD')},
    {label: "DAU", value: "DAU", type:"custom",width:"90px"},
    {label: "总投注", value: "BetValue", type:"custom"},
    {label: "总赢分", value: "WinValue", type:"custom"},
    {label: "系统收益", value: "SystemValue", type:"custom"},
    {label: "回报率", value: "RateValue", type:"custom"},
    {label: "购买小游戏总投注", value: "BuyBet", type:"custom"},
    {label: "购买小游戏总赢分", value: "BuyWin", type:"custom"},
    {label: "购买小游戏系统收益", value: "BuySystemValue", type:"custom", width: "180px"},
    {label: "购买小游戏回报率", value: "BuyRateValue", type:"custom"},
])
const JILITableHeader = ref([
    {label: "日期", value: "CurrencyName", format:(row)=>ut.fmtDate(row.Time, 'YYYY-MM-DD')},
    {label: "DAU", value: "DAU", type:"custom",width:"90px"},
    {label: "总投注金额", value: "BetValue", type:"custom"},
    {label: "总赢得金额", value: "WinValue", type:"custom"},
    {label: "系统收益", value: "SystemValue", type:"custom"},
    {label: "回报率", value: "RateValue", type:"custom"},
    {label: "购买小游戏总投注", value: "BuyBet", type:"custom"},
    {label: "购买小游戏总赢分", value: "BuyWin", type:"custom"},
    {label: "购买小游戏系统收益", value: "BuySystemValue", type:"custom", width: "180px"},
    {label: "购买小游戏回报率", value: "BuyRateValue", type:"custom"},
    {label: "额外模式总押注", value: "CurrencyName", type:"custom"},
    {label: "额外模式总赢分", value: "CurrencyName", type:"custom"},
    {label: "额外模式系统收益", value: "CurrencyName", type:"custom"},
    {label: "额外模式回报率", value: "CurrencyName", type:"custom"},
])

const loading = ref(false)
onBeforeUnmount(async () => {
    [DAUChart, totalPressurePointsChart, systemBenefitChart, buyTotalBetChart, responseRateChart, buyResponseRateChart].forEach(chartRef => {
        if (chartRef.value) {
            chartRef.value.dispose();
        }
    });
});
onMounted(() => {



});

onUnmounted(()=>{

})

echarts.use([
    TitleComponent,
    TooltipComponent,
    GridComponent,
    LegendComponent,
    MarkLineComponent,
    LineChart,
    CanvasRenderer,
    UniversalTransition
]);

let initChart = (refData, id, name, legendData1, legendData2, data) => {
    let chartDom = document.getElementById(id);
    let myChart = echarts.init(chartDom);
    refData.value = markRaw(myChart)

    const option = {
        title: {
            text: name,
            left:"center",
            top: 10,
            textStyle:{
                fontSize: 12
            }
        },
        tooltip: {
            trigger: 'axis',
            // formatter: '{a} <br/>{b}: {c} ({d}%)'
            dataView: { readOnly: false },
        },
        legend: {
            left:"left",
            data: [legendData2, legendData1],
            padding:[50,0,0,20]
        },
        grid: {
            containLabel: true,
            top: 90,
            bottom: 10,
        },
        xAxis: {
            type: 'category',
            boundaryGap: false,
            data: [1, 2, 3, 4, 5, 6, 7]
        },
        yAxis: {
            type: 'value',
        },
        series: data,
    };

    // 使用配置项绘制图表
    refData.value.setOption(option);
}
let dataChartFun = (dataFun, name1, name2) => {
    return dataFun.map((subArray, index) => ({
        name: index > 0 ? name2 : name1,
        type: 'line',
        data: subArray
    }));
};
let dataChartFunPercent = (dataFun, name1, name2) => {

    return dataFun.map((subArray, index) => ({
        name: index > 0 ? name2 : name1,
        type: 'line',
        data: subArray,
        markLine: {
            silent: true,
            lineStyle: {
                color: '#333'
            },
            data: [
                {
                    yAxis: 1
                },
            ]
        }
    }));
};
const defaultOperatorEvent = ref({})
const defaultGameEvent = ref({})
const selectGameListVal = reactive({
    GameID: '',
    Type: null,
    AppID: ""
})

const selectGameList = (value, lotteryValue) => {

    selectGameListVal.GameID = value.gameData
    selectGameListVal.Type = lotteryValue.value
    InitData()
}


const operatorListChange = (value) => {

    if (value){

        selectGameListVal.AppID = value.AppID.toString()
    }else{
        selectGameListVal.AppID = ""
    }

    InitData()
}



const InitData = () => {
    if (init.value){
        if (selectGameListVal.AppID != "" && selectGameListVal.GameID != ""){
            init.value = false
            initSearch()
        }

    }
}

let initSearch = async () => {

    changeShow.value = "图表"
    echartShow.value = false
    if (selectGameListVal.AppID == ""){
        return tip.e(t("请选择商户"))
    }

    let [data, err] = await Client.Do(AdminStatsRpc.GetGameEarningsMethod,
        selectGameListVal
    )
    if (err) {
        uiData.allTableData = []
        uiData.tableData = []
        uiData.SevenDays = {
            TotalReturnRate: '0.00%', // 7日总回报率
            ExtraBetAmount: 0, // 购买小游戏流水
            ExtraWinAmount: 0, // 购买小游戏产出
            extra: '0.00%', // 7日额外模式回报率
            BuyBetAmount: 0, // 额外模式流水
            BuyWinAmount: 0, // 额外模式产出
            BetAmount: 0, // 7日总流水
            WinAmount: 0, // 7日总产出
            SystemBenefits: 0, // 系统受益
            percentage: 0, // 系统受益
        }
        uiData.MonthData = {
            TotalReturnRate: '0.00%', // 本月总回报率
            ExtraBetAmount: 0, // 购买小游戏流水
            ExtraWinAmount: 0, // 购买小游戏产出
            extra: '0.00%', // 本月额外模式回报率
            BuyBetAmount: 0, // 额外模式流水
            BuyWinAmount: 0, // 额外模式产出
            BetAmount: 0, // 本月总流水
            WinAmount: 0, // 本月总产出
            SystemBenefits: 0, // 系统受益
            percentage: 0, // 系统受益
        }
        uiData.TotalData = {
            TotalReturnRate: '0.00%', // 总回报率
            ExtraBetAmount: 0, // 购买小游戏流水
            ExtraWinAmount: 0, // 购买小游戏产出
            extra: '0.00%', // 总额外模式回报率
            BuyBetAmount: 0, // 额外模式流水
            BuyWinAmount: 0, // 额外模式产出
            BetAmount: 0, // 总流水
            WinAmount: 0, // 总产出
            SystemBenefits: 0, // 系统受益
            percentage: 0, // 系统受益
        }
        return tip.e(err)
    }
    echartShow.value = true
    uiData.SevenDays = returnData(data.SevenData)
    uiData.MonthData = returnData(data.MonthData)
    uiData.TotalData = returnData(data.TotalData)
    uiData.allTableData = data
    let LastWeek = JSON.parse(JSON.stringify(data.WeekData.LastWeek)).sort((a, b) => {
        return a.Time - b.Time
    })
    let NowWeek = JSON.parse(JSON.stringify(data.WeekData.NowWeek)).sort((a, b) => {
        return a.Time - b.Time
    })
    let weekDataList = LastWeek.concat(Object.values(NowWeek))
    let chartLastWeekList = LastWeek.map(t => {
        return {
            ...t,
            DauValue: t.DauValue,
            BetValue: t.BetValue / 10000,
            BuyBet: t.BuyBet / 10000,
            WinValue: t.WinValue / 10000,
            SystemValue: t.SystemValue / 10000,
            BuyRateValue: (Number(t.BuyRateValue) * 100).toFixed(2),
            BuySystemValue: t.BuySystemValue / 10000,
            BuyWin: t.BuyWin / 10000,
            RateValue: (t.RateValue*100).toPrecision(4),
        }
    })
    let chartNowWeekList = NowWeek.map(t => {
        return {
            ...t,
            DauValue: t.DauValue,
            BetValue: t.BetValue / 10000,
            BuyBet: t.BuyBet / 10000,
            WinValue: t.WinValue / 10000,
            SystemValue: t.SystemValue / 10000,
            BuyRateValue: (Number(t.BuyRateValue) * 100).toFixed(2),
            BuySystemValue: t.BuySystemValue / 10000,
            BuyWin: t.BuyWin / 10000,
            RateValue: (Number(t.RateValue) * 100).toFixed(2),
        }
    })
    await processChartData(chartLastWeekList, chartNowWeekList)
    await renderCharts()
    await radioChange()
}
// 渲染图表
const processChartData = async (last, now) => {
    let chartDAUfun = processTableData(last, now, 'DauValue')
    let chartTotalPressurePointsfun = processTableData(last, now, 'BetValue')
    let chartSystemBenefitfun = processTableData(last, now, 'SystemValue')
    let chartBuyTotalBetfun = processTableData(last, now, 'BuyBet')
    let chartResponseRatefun = processTableData(last, now, 'RateValue')
    let chartBuyResponseRatefun = processTableData(last, now, 'BuyRateValue')

    uiData.chartDAU = dataChartFun(chartDAUfun, t('上周DAU'), t('本周DAU'))
    uiData.chartTotalPressurePoints = dataChartFun(chartTotalPressurePointsfun, t('上周总投注'), t('本周总投注'))
    uiData.chartSystemBenefit = dataChartFun(chartSystemBenefitfun, t('上周系统收益'), t('本周系统收益'))
    uiData.chartBuyTotalBet = dataChartFun(chartBuyTotalBetfun, t('上周购买小游戏总投注'), t('本周购买小游戏总投注'))
    uiData.chartResponseRate = dataChartFunPercent(chartResponseRatefun, t('上周回报率'), t('本周回报率'))
    uiData.chartBuyResponseRate = dataChartFunPercent(chartBuyResponseRatefun, t('上周购买小游戏回报率'), t('本周购买小游戏回报率'))

};
// 渲染图表
const renderCharts = async () => {
    await nextTick(() => {
        initChart(DAUChart, 'DAUChart', t('DAU双周分析'), t('上周DAU'), t('本周DAU'), uiData.chartDAU)
        initChart(totalPressurePointsChart, 'totalPressurePointsChart', t('总投注双周分析'), t('上周总投注'), t('本周总投注'), uiData.chartTotalPressurePoints)
        initChart(systemBenefitChart, 'systemBenefitChart', t('系统收益双周分析'), t('上周系统收益'), t('本周系统收益'), uiData.chartSystemBenefit)
        initChart(buyTotalBetChart, 'buyTotalBetChart', t('购买小游戏总投注双周分析'), t('上周购买小游戏总投注'), t('本周购买小游戏总投注'), uiData.chartBuyTotalBet)
        initChart(responseRateChart, 'responseRateChart', t('回报率双周分析') + '(%)', t('上周回报率'), t('本周回报率'), uiData.chartResponseRate)
        initChart(buyResponseRateChart, 'buyResponseRateChart', t('购买小游戏回报率双周分析') + '(%)', t('上周购买小游戏回报率'), t('本周购买小游戏回报率'), uiData.chartBuyResponseRate)
    })
};


let returnData = (data) => {
    return {
        TotalReturnRate: ut.fmtPercent(data.WinAmount / data.BetAmount || 0), // 7日总回报率
        buyRateTotalReturnRate: ut.fmtPercent(data.BuyWinAmount / data.BuyBetAmount || 0), // 7日总回报率
        extra: ut.fmtPercent(data.ExtraWinAmount / data.ExtraBetAmount || 0), // 7日总回报率
        BetAmount: ut.toNumberWithComma(data.BetAmount), // 7日总流水
        WinAmount: ut.toNumberWithComma(data.WinAmount), // 7日总产出
        BuyWinAmount: ut.toNumberWithComma(data.BuyWinAmount), // 小游戏产出
        BuyBetAmount: ut.toNumberWithComma(data.BuyBetAmount), // 小游戏流水
        ExtraBetAmount: ut.toNumberWithComma(data.ExtraBetAmount), // 购买小游戏流水
        ExtraWinAmount: ut.toNumberWithComma(data.ExtraWinAmount), // 购买小游戏产出
        SystemBenefits: ut.toNumberWithComma(data.BetAmount - data.WinAmount), // 系统收益
        percentage: Number((data.WinAmount / data.BetAmount).toPrecision(4)) * 100, // 环
    }
}
let processTableData = (chartLastWeekList, chartNowWeekList, field) => {

    let extractedData = chartLastWeekList.map((t) => t[field]);
    let extractedData1 = chartNowWeekList.map((t) => t[field]);
    let resultArrays = [];
    resultArrays.push(extractedData, extractedData1);

    return resultArrays;
}

let radioChange = async () => {
    if (!uiData.allTableData) {
        return
    }
    let LastWeek = JSON.parse(JSON.stringify(uiData.allTableData.WeekData.LastWeek))
    let NowWeek = JSON.parse(JSON.stringify(uiData.allTableData.WeekData.NowWeek))
    uiData.tableData = NowWeek.concat(Object.values(LastWeek))
    uiData.SevenDays = returnData(uiData.allTableData.SevenData)
    uiData.MonthData = returnData(uiData.allTableData.MonthData)
    uiData.TotalData = returnData(uiData.allTableData.TotalData)
    setRuntimeTable(uiData)
}
let xiazai = () => {
    let fmtGold = (v) => ut.fmtGold(v)
    let fmtPercent = (v) => ut.fmtPercent(v)
    let LastWeek = JSON.parse(JSON.stringify(uiData.allTableData.WeekData.LastWeek))
    let NowWeek = JSON.parse(JSON.stringify(uiData.allTableData.WeekData.NowWeek))
    uiData.tableData = NowWeek.concat(Object.values(LastWeek))
    excel.dump(uiData.tableData, UserName, [
        {key: "Time", name: t("日期")},
        {key: "DauValue", name: t("DAU")},
        {key: "BetValue", name: t("总押分"), fmt: fmtGold},
        {key: "SystemValue", name: t("系统受益"), fmt: fmtGold},
        {key: "RateValue", name: t("回报率"), fmt: fmtPercent},
        {key: "BuyBet", name: t("购买小游戏总押注"), fmt: fmtGold},
        {key: "BuyRateValue", name: t("购买小游戏回报率"), fmt: fmtPercent},
    ])
}
</script>
<style lang="scss" scoped>
.gameList{
  margin-bottom: 15px;
}
.el-statistic {
    --el-statistic-content-font-size: 20px;
}

.top-chart {
    .top-chart-col {
        display: flex;
        justify-content: center;
        min-height: 270px;

        :deep(.el-progress-circle) {
            width: calc(100% / 3 - 100px);
            min-width: 300px;
            margin: 0 auto;
        }

        :deep(.el-progress-circle__track) {
            stroke: darkgray;
        }

        :deep(.el-progress__text) {
            font-size: 14px !important;
            width: 90%;
            transform: translate(5%, -50%);

            .top-chart {
                .el-col {
                    margin-top: 1rem;
                }
            }
        }
    }
}

.bottom-chart {

    .bottom-chart-col {
        display: flex;
        justify-content: center;
        margin-top: 5rem;
        border: 1px solid #c0ffc9;
        border-radius: 8px;
        background: #ffffff;
    }

}
:deep(.el-descriptions__body) {
    margin-top: 1rem;
    padding: 12px 0 0 12px;

    .el-descriptions__label:not(.is-bordered-label) {
        color: rgba(18, 31, 62, 0.74);
    }

    .el-descriptions__content {
        font-weight: bold;
    }
}

</style>
