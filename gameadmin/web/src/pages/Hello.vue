<template>
    <div>
        <div class="searchView gameList">
            <el-form
                label-position="top"
                label-width="100px"
                style="max-width: 100%"
            >
                <el-space wrap>
                    <gamelist_container :defaultGameEvent="defaultGameEvent"
                                        @select-operator="selectGameList"></gamelist_container>


                    <operator_container :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operator="operatorListChange"></operator_container>

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
                <el-col :span="8" class="top-chart-col">
                    <el-progress
                        :percentage="uiData.SevenDays.percentage > 100 ? uiData.SevenDays.percentage - 100 : uiData.SevenDays.percentage"
                        class="7days-chart" stroke-width="3"
                        indeterminate type="dashboard" color="steelblue" striped striped-flow>
                        <template #default="{ percentage }">
                            <el-row class="top-chart" :gutter="10">
                                <el-col :span="12">
                                    <el-statistic :title="$t('7日总流水')" :value="uiData.SevenDays.BetAmount"/>
                                </el-col>
                                <el-col :span="12">
                                    <el-statistic :title="$t('7日总产出')" :value="uiData.SevenDays.BetAmount"/>
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
                <el-col :span="8" class="top-chart-col">
                    <el-progress class="7days-chart" stroke-width="3"
                                 :percentage="uiData.MonthData.percentage>100?uiData.MonthData.percentage-100:uiData.MonthData.percentage"
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
<!--                <el-col :span="8">-->

<!--                    <div style="border: 1px solid #bcbcbc;background: white">-->
<!--                        <el-table :data="tableData" border style="width: 100%">-->
<!--                            <el-table-column prop="notifyType" label="公告类型" width="100"/>-->
<!--                            <el-table-column prop="notifyTitle" label="公告标题" width="100"/>-->
<!--                            <el-table-column prop="notifyContent" label="公告内容" />-->

<!--                        </el-table>-->
<!--                    </div>-->

<!--                </el-col>-->
            </el-row>
            <el-row class="bottom-chart">
                <el-col :span="8" class="bottom-chart-col">
                    <div id="DAUChart" ref="7days-chart" style="width: 800px; height: 400px;" class="DAU-chart"></div>
                </el-col>
                <el-col :span="8" class="bottom-chart-col">
                    <div id="totalPressurePointsChart" ref="total-pressure-points-chart"
                         style="width: 800px; height: 400px;" class="total-pressure-points-chart"></div>
                </el-col>
                <el-col :span="8" class="bottom-chart-col">
                    <div id="systemBenefitChart" ref="system-benefit-chart" style="width: 800px; height: 400px;"
                         class="system-benefit-chart"></div>
                </el-col>
            </el-row>
        </div>
        <div v-show="changeShow === '表格'">
            <runtime-table :gameMsg="selectGameListVal?.Game"></runtime-table>

            <div class="table-icon" @click="xiazai">
                <el-icon>
                    <Download/>
                </el-icon>
            </div>
            <el-table :data="uiData.tableData" style="width: 100%">
                <el-table-column :label="$t('日期')" prop="Time"></el-table-column>
                <el-table-column :label="$t('DAU')">
                    <template #default="scope">{{ scope.row.DauValue }}</template>
                </el-table-column>
                <el-table-column :label="$t('总押分')" prop="BetValue">
                    <template #default="scope">{{ ut.toNumberWithComma(scope.row.BetValue) }}</template>
                </el-table-column>
                <el-table-column :label="$t('总赢得金额')" prop="WinValue">
                    <template #default="scope">{{ ut.toNumberWithComma(scope.row.WinValue) }}</template>
                </el-table-column>
                <el-table-column :label="$t('系统收益')">
                    <template #default="scope">{{ ut.toNumberWithComma(scope.row.SystemValue) }}</template>
                </el-table-column>
                <el-table-column :label="$t('回报率')" prop="RateValue">
                    <template #default="scope">{{ percentFormatter(0, 0, scope.row.RateValue) }}</template>
                </el-table-column>
                <el-table-column :label="$t('购买小游戏总押注')" prop="BuyBet">
                    <template #default="scope">{{ ut.toNumberWithComma(scope.row.BuyBet) }}</template>
                </el-table-column>
                <el-table-column :label="$t('购买小游戏总赢分')" prop="BuyWin">
                    <template #default="scope">{{ ut.toNumberWithComma(scope.row.BuyWin) }}</template>
                </el-table-column>
                <el-table-column :label="$t('购买小游戏系统收益')" prop="BuySystemValue">
                    <template #default="scope">{{ ut.toNumberWithComma(scope.row.BuySystemValue) }}</template>
                </el-table-column>
                <el-table-column :label="$t('购买小游戏回报率')" prop="BuyRateValue">
                    <template #default="scope">{{ percentFormatter(0, 0, scope.row.BuyRateValue) }}</template>
                </el-table-column>
                <template v-if="selectGameListVal?.Game?.startsWith('jili_')">
                    <el-table-column :label="$t('额外模式总押注')" prop="BuyRateValue">
                        <template #default="scope">{{ ut.toNumberWithComma(scope.row.ExtraBet) }}</template>
                    </el-table-column>
                    <el-table-column :label="$t('额外模式总赢分')" prop="BuyRateValue">
                        <template #default="scope">{{ ut.toNumberWithComma(scope.row.ExtraWin) }}</template>
                    </el-table-column>
                    <el-table-column :label="$t('额外模式系统收益')" prop="BuyRateValue">
                        <template #default="scope">{{ ut.toNumberWithComma(scope.row.ExtraSystemValue) }}</template>
                    </el-table-column>
                    <el-table-column :label="$t('额外模式回报率')" prop="BuyRateValue">
                        <template #default="scope">{{ percentFormatter(0, 0, scope.row.ExtraRateValue) }}</template>
                    </el-table-column>
                </template>
            </el-table>
        </div>
    </div>
</template>

<script lang="ts" setup>
import {onMounted, ref, reactive, onBeforeUnmount, nextTick, shallowRef} from 'vue';
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

const DAUChart = shallowRef(null);
const totalPressurePointsChart = shallowRef(null);
const systemBenefitChart = shallowRef(null);
const buyTotalBetChart = shallowRef(null);
const responseRateChart = shallowRef(null);
const buyResponseRateChart = shallowRef(null);
const {t} = useI18n()

const changeShow = ref(t('图表'))
const echartShow = ref(false)
const piniaStore = useStore()
const {setRuntimeTable, language} = piniaStore
import RuntimeTable from "@/pages/runtime/runtimeTable.vue";
import {Notify} from "@/api/adminpb/notify";
import Operator_container from "@/components/operator_container.vue";
import {debug} from "util";
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
let UserName = ref('')
const value = ref([])
const loading = ref(false)

const tableData = ref([])



onBeforeUnmount(async () => {
    [DAUChart, totalPressurePointsChart, systemBenefitChart, buyTotalBetChart, responseRateChart, buyResponseRateChart].forEach(chartRef => {
        if (chartRef.value) {
            chartRef.value.dispose();
        }
    });

});


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
            subtext: name,
        },
        tooltip: {
            trigger: 'axis',
            // formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: {
            data: [legendData2, legendData1]
        },
        grid: {
            containLabel: true
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
    Type: null
})
const selectOperatorVal = reactive({
    OperatorId: 0
})
const selectGameList = (value, lotteryValue) => {

    selectGameListVal.GameID = value.value
    selectGameListVal.Type = lotteryValue.value
}
const operatorListChange = (value) => {

    selectOperatorVal.OperatorId = value.value
}
let initSearch = async () => {
    echartShow.value = false
    let [data, err] = await Client.Do(AdminStatsRpc.GetGameEarningsMethod,
        {...selectGameListVal,...selectOperatorVal }
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
            BuyRateValue: ut.fmtPercentFour(t.BuyRateValue, 4),
            BuySystemValue: t.BuySystemValue / 10000,
            BuyWin: t.BuyWin / 10000,
            RateValue: (t.RateValue).toPrecision(4),
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
            BuyRateValue: ut.fmtPercentFour(t.BuyRateValue, 4),
            BuySystemValue: t.BuySystemValue / 10000,
            BuyWin: t.BuyWin / 10000,
            RateValue: (t.RateValue).toPrecision(4),
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
    uiData.chartTotalPressurePoints = dataChartFun(chartTotalPressurePointsfun, t('上周总押分'), t('本周总押分'))
    uiData.chartSystemBenefit = dataChartFun(chartSystemBenefitfun, t('上周系统收益'), t('本周系统收益'))
    uiData.chartBuyTotalBet = dataChartFun(chartBuyTotalBetfun, t('上周购买小游戏总押注'), t('本周购买小游戏总押注'))
    uiData.chartResponseRate = dataChartFunPercent(chartResponseRatefun, t('上周回报率'), t('本周回报率'))
    uiData.chartBuyResponseRate = dataChartFunPercent(chartBuyResponseRatefun, t('上周购买小游戏回报率'), t('本周购买小游戏回报率'))
};
// 渲染图表
const renderCharts = async () => {
    await nextTick(() => {
        initChart(DAUChart, 'DAUChart', t('DAU双周分析'), t('上周DAU'), t('本周DAU'), uiData.chartDAU)
        initChart(totalPressurePointsChart, 'totalPressurePointsChart', t('总押分双周分析'), t('上周总押分'), t('本周总押分'), uiData.chartTotalPressurePoints)
        initChart(systemBenefitChart, 'systemBenefitChart', t('系统收益双周分析'), t('上周系统收益'), t('本周系统收益'), uiData.chartSystemBenefit)
        initChart(buyTotalBetChart, 'buyTotalBetChart', t('购买小游戏总押注双周分析'), t('上周购买小游戏总押注'), t('本周购买小游戏总押注'), uiData.chartBuyTotalBet)
        initChart(responseRateChart, 'responseRateChart', t('回报率双周分析'), t('上周回报率'), t('本周回报率'), uiData.chartResponseRate)
        initChart(buyResponseRateChart, 'buyResponseRateChart', t('购买小游戏回报率双周分析'), t('上周购买小游戏回报率'), t('本周购买小游戏回报率'), uiData.chartBuyResponseRate)
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


// 图标，表格  tab切换
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


// 公告类型
const notifyType = {
    "1": "系统",
    "2": "维护",
    "3": "优惠",
    "4": "推荐"
}
// 状态
const notifyStatus = {
    "1": "启用",
    "2": "停用",
}

const getNotify = async () => {
    // let response = await Client.Do(Notify.GetEffectiveNotify, {Page: 1, PageSize:10000} as any)
    //
    // response[0].Notifications.forEach(item => {
    //     item.notifyType = notifyType[item.notifyType]
    //     item.notifyContent = item[language]
    //     tableData.value.push(item)
    // })
}
getNotify()

</script>
<style lang="scss" scoped>
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
