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
                </el-space>
            </el-form>
            <el-space wrap>
                <el-button type="primary" @click="initSearch">{{ $t('搜索') }}</el-button>
            </el-space>
            <template v-if="selectGameListVal.Game === 'Hilo' && changeShow === $t('表格')">
                <el-form
                    label-position="top"
                    label-width="100px"
                    style="max-width: 100%"
                >
                    <el-form-item :label="$t('Hilo')">
                        <el-space wrap>
                            <el-input-number :controls="false" v-model.trim="HiloData" clearable :placeholder="$t('请输入')"
                                             size="default"/>
                            <el-button circle plain @click="setHiloData()" style="margin-bottom:0;margin-left: .5rem">
                                <el-icon><Select/></el-icon>
                            </el-button>
                            <el-popover placement="right" :width="600" trigger="click">
                                <template #reference>
                                    <div class="table-icon" @click="setHistory" style="margin: 0">
                                        <el-icon ><Tickets /></el-icon>
                                    </div>
                                </template>
                                <el-table :data="setHistoryList" style="width: 100%" stripe  height="480" highlight-current-row>
                                    <el-table-column prop="InsertTime" label="时间" width="150" :formatter="dateFormater" />
                                    <el-table-column prop="Username" label="操作人" width="80" />
                                    <el-table-column prop="OldGold" label="调整前数据" >
                                        <template #default="scope">
                                            {{ut.toNumberWithComma(scope.row.Before)}}
                                        </template>
                                    </el-table-column>
                                    <el-table-column prop="Change" label="变化值" >
                                        <template #default="scope">
                                            <el-tag effect="plain">
                                                {{scope.row.Change > 0 ?'+ ' + ut.toNumberWithComma(scope.row.Change) : ut.toNumberWithComma(scope.row.Change)}}
                                            </el-tag>
                                        </template>
                                    </el-table-column>
                                    <el-table-column prop="NewGold" label="调整后数据" >
                                        <template #default="scope">
                                            {{ut.toNumberWithComma(scope.row.After)}}
                                        </template>
                                    </el-table-column>
                                </el-table>
                            </el-popover>
                        </el-space>

                    </el-form-item>
                </el-form>
            </template>
        </div>

        <el-radio-group v-model="changeShow" @change="radioChange">
            <el-radio-button :label="$t('图表')"/>
            <el-radio-button :label="$t('表格')"/>
        </el-radio-group>
        <div v-show="changeShow === $t('图表') && echartShow">
            <el-row class="top-chart">
                <el-col :span="8" class="top-chart-col">
                    <el-progress
                        :percentage="uiData.SevenDays.percentage>100?uiData.SevenDays.percentage-100:uiData.SevenDays.percentage"
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
                <el-col :span="8" class="top-chart-col">
                    <el-progress stroke-width="3"
                                 :percentage="uiData.TotalData.percentage>100?uiData.TotalData.percentage-100:uiData.TotalData.percentage"
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
                <el-col :span="8" class="bottom-chart-col">
                    <div id="buyTotalBetChart" ref="buy-total-bet-chart" style="width: 800px; height: 400px;"
                         class="buy-total-bet-chart"></div>
                </el-col>
                <el-col :span="8" class="bottom-chart-col">
                    <div id="responseRateChart" ref="response-rate-chart" style="width: 800px; height: 400px;"
                         class="response-rate-chart"></div>
                </el-col>
                <el-col :span="8" class="bottom-chart-col">
                    <div id="buyResponseRateChart" ref="buy-response-rate-chart" style="width: 800px; height: 400px;"
                         class="buy-response-rate-chart"></div>
                </el-col>
            </el-row>
        </div>
        <div v-show="changeShow === $t('表格')">
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
                <el-table-column :label="$t('总旋转次数')" prop="SpinCount"></el-table-column>
                <el-table-column :label="$t('总押分')" prop="BetValue">
                    <template #default="scope">{{ ut.toNumberWithComma(scope.row.BetValue) }}</template>
                </el-table-column>
                <el-table-column :label="$t('总赢分')" prop="WinValue">
                    <template #default="scope">{{ ut.toNumberWithComma(scope.row.WinValue) }}</template>
                </el-table-column>
                <el-table-column :label="$t('系统收益')">
                    <template #default="scope">{{ ut.toNumberWithComma(scope.row.SystemValue) }}</template>
                </el-table-column>

                <el-table-column :label="$t('个人奖池爆奖')" prop="BigReward">
                    <template #default="scope">{{ ut.toNumberWithComma(scope.row.BigReward) }}</template>
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
import type {FormInstance, FormRules} from 'element-plus'
import {GetAdminInfo, useStore} from '@/pinia/index';
import {useI18n} from 'vue-i18n';
import {excel} from "@/lib/excel";
import {markRaw} from 'vue'
import Gamelist_container from "@/components/gamelist_container.vue";

const DAUChart = ref(null);
const totalPressurePointsChart = ref(null);
const systemBenefitChart = ref(null);
const buyTotalBetChart = ref(null);
const responseRateChart = ref(null);
const buyResponseRateChart = ref(null);
const {t} = useI18n()
const changeShow = ref(t('图表'))
const echartShow = ref(false)
const piniaStore = useStore()
const {setRuntimeTable} = piniaStore
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
        BuyBetAmount: 0, // 额外模式流水
        BuyWinAmount: 0, // 额外模式产出
        extra: '0.00%', // 7日额外模式回报率
        buyRateTotalReturnRate: '0.00%', // 7日总回报率
        BetAmount: 0, // 7日总流水
        WinAmount: 0, // 7日总产出
        SystemBenefits: 0, // 系统受益
        percentage: 0, // 系统受益
    },
    MonthData: {
        buyRateTotalReturnRate: '0.00%', // 本月总回报率
        ExtraBetAmount: 0, // 购买小游戏流水
        ExtraWinAmount: 0, // 购买小游戏产出
        extra: '0.00%', // 本月额外模式回报率
        TotalReturnRate: '0.00%', // 本月总回报率
        BetAmount: 0, // 本月总流水
        WinAmount: 0, // 本月总产出
        SystemBenefits: 0, // 系统受益
        percentage: 0, // 系统受益
    },
    TotalData: {
        buyRateTotalReturnRate: '0.00%', // 总回报率
        ExtraBetAmount: 0, // 购买小游戏流水
        ExtraWinAmount: 0, // 购买小游戏产出
        extra: '0.00%', // 总额外模式回报率
        TotalReturnRate: '0.00%', // 总回报率
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
const setHistoryList = ref(null)
const loading = ref(false)
onBeforeUnmount(async () => {
    // 在组件销毁之前销毁ECharts图表，释放资源
    /*if (daysChart7.value) {
        daysChart7.value.dispose();
    }*/
    if (DAUChart.value) {
        DAUChart.value.dispose();
    }
    if (totalPressurePointsChart.value) {
        totalPressurePointsChart.value.dispose();
    }
    if (systemBenefitChart.value) {
        systemBenefitChart.value.dispose();
    }
    if (buyTotalBetChart.value) {
        buyTotalBetChart.value.dispose();
    }
    if (responseRateChart.value) {
        responseRateChart.value.dispose();
    }
    if (buyResponseRateChart.value) {
        buyResponseRateChart.value.dispose();
    }
});
onMounted(() => {
});

const defaultOperatorEvent = ref({})
const defaultGameEvent = ref({})
const selectGameListVal = reactive({
    Game: '',
    Type: null

})
const selectGameList = (value, lotteryValue) => {
    selectGameListVal.Game = value.value
    selectGameListVal.Type = lotteryValue.value
}

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
import RuntimeTable from "@/pages/runtime/runtimeTable.vue";
import {AdminSlotsRpc} from "@/api/slots/admin_rpc";

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
const store = useStore();
const HiloData = ref('');
let setHiloData = async () => {
    if (selectGameListVal.Game === 'Hilo') {
        let [datas, err] = await Client.Do(AdminStatsRpc.SetawardPool,
            {AwardPool: Number(HiloData.value)}
        )
        loading.value = false
        if (err) {
            return tip.e(err)
        }
        tip.s(t('设置成功'))
    }
}
let initSearch = async () => {
    echartShow.value = false
    let [data, err] = await Client.Do(AdminStatsRpc.GetGameEarningsMethod,
        selectGameListVal
    )
    if (selectGameListVal.Game === 'Hilo') {
        let [datas, err] = await Client.Do(AdminStatsRpc.GetawardPool,
            {}
        )
        HiloData.value = datas.AwardPool
    }
    if (err) {
        uiData.allTableData = []
        uiData.tableData = []
        uiData.SevenDays = {
            TotalReturnRate: '0.00%', // 7日总回报率
            ExtraBetAmount: 0, // 购买小游戏流水
            ExtraWinAmount: 0, // 购买小游戏产出
            BuyBetAmount: 0, // 额外模式流水
            BuyWinAmount: 0, // 额外模式产出
            extra: '0.00%', // 7日额外模式回报率
            BetAmount: 0, // 7日总流水
            WinAmount: 0, // 7日总产出
            SystemBenefits: 0, // 系统受益
            percentage: 0, // 系统受益
        }
        uiData.MonthData = {
            TotalReturnRate: '0.00%', // 本月总回报率
            ExtraBetAmount: 0, // 购买小游戏流水
            ExtraWinAmount: 0, // 购买小游戏产出
            BuyBetAmount: 0, // 额外模式流水
            BuyWinAmount: 0, // 额外模式产出
            extra: '0.00%', // 本月额外模式回报率
            BetAmount: 0, // 本月总流水
            WinAmount: 0, // 本月总产出
            SystemBenefits: 0, // 系统受益
            percentage: 0, // 系统受益
        }
        uiData.TotalData = {
            TotalReturnRate: '0.00%', // 总回报率
            ExtraBetAmount: 0, // 购买小游戏流水
            ExtraWinAmount: 0, // 购买小游戏产出
            BuyBetAmount: 0, // 额外模式流水
            BuyWinAmount: 0, // 额外模式产出
            extra: '0.00%', // 总额外模式回报率
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
    let chartDAUfun = processTableData(chartLastWeekList, chartNowWeekList, 'DauValue')
    let chartTotalPressurePointsfun = processTableData(chartLastWeekList, chartNowWeekList, 'BetValue')
    let chartSystemBenefitfun = processTableData(chartLastWeekList, chartNowWeekList, 'SystemValue')
    let chartBuyTotalBetfun = processTableData(chartLastWeekList, chartNowWeekList, 'BuyBet')
    let chartResponseRatefun = processTableData(chartLastWeekList, chartNowWeekList, 'RateValue')
    let chartBuyResponseRatefun = processTableData(chartLastWeekList, chartNowWeekList, 'BuyRateValue')
    uiData.chartDAU = dataChartFun(chartDAUfun, t('上周DAU'), t('本周DAU'))
    uiData.chartTotalPressurePoints = dataChartFun(chartTotalPressurePointsfun, t('上周总押分'), t('本周总押分'))
    uiData.chartSystemBenefit = dataChartFun(chartSystemBenefitfun, t('上周系统收益'), t('本周系统收益'))
    uiData.chartBuyTotalBet = dataChartFun(chartBuyTotalBetfun, t('上周购买小游戏总押注'), t('本周购买小游戏总押注'))
    uiData.chartResponseRate = dataChartFunPercent(chartResponseRatefun, t('上周回报率'), t('本周回报率'))
    uiData.chartBuyResponseRate = dataChartFunPercent(chartBuyResponseRatefun, t('上周购买小游戏回报率'), t('本周购买小游戏回报率'))
    await nextTick(() => {
        initChart(DAUChart, 'DAUChart', t('DAU双周分析'), t('上周DAU'), t('本周DAU'), uiData.chartDAU)
        initChart(totalPressurePointsChart, 'totalPressurePointsChart', t('总押分双周分析'), t('上周总押分'), t('本周总押分'), uiData.chartTotalPressurePoints)
        initChart(systemBenefitChart, 'systemBenefitChart', t('系统收益双周分析'), t('上周系统收益'), t('本周系统收益'), uiData.chartSystemBenefit)
        initChart(buyTotalBetChart, 'buyTotalBetChart', t('购买小游戏总押注双周分析'), t('上周购买小游戏总押注'), t('本周购买小游戏总押注'), uiData.chartBuyTotalBet)
        initChart(responseRateChart, 'responseRateChart', t('回报率双周分析'), t('上周回报率'), t('本周回报率'), uiData.chartResponseRate)
        initChart(buyResponseRateChart, 'buyResponseRateChart', t('购买小游戏回报率双周分析'), t('上周购买小游戏回报率'), t('本周购买小游戏回报率'), uiData.chartBuyResponseRate)
    })
    await radioChange()
}
let returnData = (data) => {
    return {
        TotalReturnRate: ut.fmtPercent(data.WinAmount / data.BetAmount || 0), // 7日总回报率
        extra: ut.fmtPercent(data.ExtraWinAmount / data.ExtraBetAmount || 0), // 7日总回报率
        buyRateTotalReturnRate: ut.fmtPercent(data.BuyWinAmount / data.BuyBetAmount || 0), // 7日总回报率
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
    if (!uiData.allTableData || Object.keys(uiData.allTableData).length === 0) {
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

const debounce = (func, delay, immediate = false) => {
    let timer;
    const debounced = function(...args) {
        const context = this;
        const later = () => {
            timer = null;
            if (!immediate) func.apply(context, args);
        };
        const callNow = immediate && !timer;
        clearTimeout(timer);
        timer = setTimeout(later, delay);
        if (callNow) func.apply(context, args);
    };
    debounced.cancel = () => clearTimeout(timer);
    return debounced;
};

const setHistory  = debounce(async () => {

    try{
        let [responseData, err] = await Client.Do(AdminSlotsRpc.setHistoryList, {
            "PageSize":6,
            "PageNumber":1
        })
        if (err) return tip.e(err)
        setHistoryList.value = responseData.List || []
        console.log(responseData);
    } catch (e) {
        tip.e(e)
    } finally {
        console.log('');
    }
}, 200);
</script>
<style lang="scss" scoped>

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

:deep(.el-descriptions__label:not(.is-bordered-label)) {
    max-width: 150px;
    min-width: 80px;
    display: inline-block;
}

</style>
