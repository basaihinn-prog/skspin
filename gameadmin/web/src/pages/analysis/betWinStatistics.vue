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
                                          @select-operatorInfo="operatorListChange" :is-init="true"></operator_container>


                        <el-form-item :label="$t('币种')">
                            <el-input v-model="CurrencyCode" disabled></el-input>
                        </el-form-item>


                        <el-form-item>

                            <el-radio-group v-model="changeShow" @change="radioChange">
                                <el-radio-button :label="$t('打码量')"/>
                                <el-radio-button :label="$t('游戏盈利')"/>
                            </el-radio-group>
                        </el-form-item>
                        <el-form-item>


                            <el-button type="primary" @click="initSearch">{{ $t('搜索') }}</el-button>
                        </el-form-item>
                    </el-space>
                </el-form>


            </el-space>
        </div>

        <div v-loading="loading">
            <el-row class="bottom-chart" :gutter="20">
                <el-col :span="8" class="bottom-chart-col">
                    <div id="NowDay" ref="NowDay-chart" style="width: 100%; height: 600px;" class="NowDay-chart"></div>

                    <el-table :data="uiData.NowDayTableData" style="width: 100%">
                        <el-table-column :label="$t('游戏名称')" prop="name">
                            <template #default="scope">
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}">{{ scope.row.name }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column :label="changeShow===$t('打码量')?$t('打码量（投注金额）'):$t('游戏盈利')">
                            <template #default="scope">
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}" v-if="changeShow===$t('打码量')">{{ scope.row.BetNum }}</span>
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}" v-else>{{ scope.row.WinNum }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column :label="changeShow===$t('打码量')?$t('打码量占比'):$t('游戏盈利占比')" prop="Proportion">
                            <template #default="scope">
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}">{{ scope.row.Proportion }}</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-col>
                <el-col :span="8" class="bottom-chart-col">
                    <div id="SevenDay" ref="SevenDay-chart" style="width: 100%; height: 600px;" class="SevenDay-chart"></div>

                    <el-table :data="uiData.SevenDayTableData" style="width: 100%">
                        <el-table-column :label="$t('游戏名称')" prop="name">
                            <template #default="scope">
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}">{{ scope.row.name }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column :label="changeShow===$t('打码量')?$t('打码量（投注金额）'):$t('游戏盈利')">
                            <template #default="scope">
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}" v-if="changeShow===$t('打码量')">{{ scope.row.BetNum }}</span>
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}" v-else>{{ scope.row.WinNum }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column :label="changeShow===$t('打码量')?$t('打码量占比'):$t('游戏盈利占比')" prop="Proportion">
                            <template #default="scope">
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}">{{ scope.row.Proportion }}</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-col>
                <el-col :span="8" class="bottom-chart-col">
                    <div id="NowMonth" ref="NowMonth-chart" style="width: 100%;height: 600px;" class="NowMonth-chart"></div>

                    <el-table :data="uiData.NowMonthTableData" style="width: 100%">
                        <el-table-column :label="$t('游戏名称')" prop="name">
                            <template #default="scope">
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}">{{ scope.row.name }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column :label="changeShow===$t('打码量')?$t('打码量（投注金额）'):$t('游戏盈利')">
                            <template #default="scope">
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}" v-if="changeShow===$t('打码量')">{{ scope.row.BetNum }}</span>
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}" v-else>{{ scope.row.WinNum }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column :label="changeShow===$t('打码量')?$t('打码量占比'):$t('游戏盈利占比')" prop="Proportion">
                            <template #default="scope">
                                <span :style="scope.row.fontB?{'font-weight':'bolder'}:{'font-weight':'normal'}">{{ scope.row.Proportion }}</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-col>
            </el-row>
        </div>
    </div>
</template>

<script lang="ts" setup>
import {markRaw, nextTick, onBeforeUnmount, onMounted, reactive, ref} from 'vue';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import ut from '@/lib/util'
import {useI18n} from 'vue-i18n';
import * as echarts from 'echarts';
import {AdminGameCenter} from "@/api/gamepb/admin";
import {excel} from "@/lib/excel";
import Analysis_platform from "@/components/analysis_platform.vue";
import Operator_container from "@/components/operator_container.vue";
const NowDayChart = ref(null);
const NowMonthChart = ref(null);
const SevenDayChart = ref(null);
const {t} = useI18n()
const changeShow = ref(t('打码量'))
const echartShow = ref(true)
let uiData = reactive({
    NowDayTableData: [],
    NowMonthTableData: [],
    SevenDayTableData: [],
    gameName: "",
    TotalExpense: 0,
    TotalIncome: 0,
    SevenDay: {
        Bet: 0,
        Win: 0,
        BetAmount: 0,
        GameID: '',
    },
    NowMonth: {
        Bet: 0,
        Win: 0,
        BetAmount: 0,
        GameID: '',
    },
    NowDay: {
        Bet: 0,
        Win: 0,
        BetAmount: 0,
        GameID: '',
    },
    TotalMonthExpense: 0,
    TotalMonthIncome: 0,
    showColumn: {},
    chartNowDay: [],
    chartNowMonth: [],
    chartSevenDay: [],
})
const value = ref([])
const loading = ref(false)
onBeforeUnmount(() => {
    if (NowDayChart.value) {
        NowDayChart.value.dispose();
    }
    if (NowMonthChart.value) {
        NowMonthChart.value.dispose();
    }
    if (SevenDayChart.value) {
        SevenDayChart.value.dispose();
    }
});
let initChart = (refData, id, name, legendData, data) => {
    refData.value = markRaw(echarts.init(document.getElementById(id)))
    const option = {
        title: {
            text: name,
            left:"center",
            top: 10,
            textStyle:{
                fontSize: 14
            }

        },
        tooltip: {
            trigger: 'item',
            position:"inside"
            // formatter: '{a} <br/>{b}: {c} ({d}%)'
        },

        series: {
            type: 'pie',
            radius: '55%',
            top: 150,
            emphasis: {
                itemStyle: {
                    shadowBlur: 10,
                    shadowOffsetX: 0,
                    shadowColor: 'rgba(0, 0, 0, 0.5)'
                }
            },
            height:400,
            data: data
        },
    };
    // 使用配置项绘制图表
    refData.value.setOption(option);
}

let returnData = (data,list) => {
    return data.map((v) => {
        return {
            name:list.find(t=>  v.GameID === t.ID)?.Name ||  t('暂无'),
            value:(changeShow.value === t('打码量'))?(ut.fmtGold(v.Bet)):(ut.fmtGold(Math.abs(v.Bet-v.Win))),
        }
    })
}
let returnTableData = (data,list) => {
    const betAmount = data.reduce((a, b) => {
        return Number(a) + Number(b.Bet);
    }, 0);
    const winAmount = data.reduce((a, b) => {
        return Number(a) + Number(b.Bet)-Number(b.Win);
    }, 0);
    data.forEach((v) => {
        let a = list.find(o=>v.GameID === o.ID)
        if (!a) {
            v.name = t('暂无')
        } else {
            v.name = a.Name
        }
        v.fontB = false
        if (changeShow.value === t('打码量')) {
            v.BetNum = ut.toNumberWithComma(v.Bet).trim()
            v.Proportion = ut.fmtPercentFour(Number(v.Bet)/Number(betAmount ? betAmount : 1),4).trim()
        } else {
            v.WinNum = ut.toNumberWithComma(Number(v.Bet)-Number(v.Win)).trim()
            v.Proportion = ut.fmtPercentFour((Number(v.Bet)-Number(v.Win))/Number(winAmount ? winAmount : 1),4).trim()
        }
    })
    let newArr = JSON.parse(JSON.stringify(data))
    let newArrSort = newArr.sort((a, b) => {
        if (changeShow.value === t('打码量')) {
            return Number(b.Bet) - Number(a.Bet);
        } else {
            return (Number(b.Bet) - Number(b.Win)) - (Number(a.Bet) - Number(a.Win)); // 降序排列
        }

    })
    if (changeShow.value === t('打码量')) {
        newArrSort.push({
            name:t('合计'),
            BetNum: ut.toNumberWithComma(betAmount).trim(),
            Proportion: '100%',
            fontB: true,
        })
    } else {
        newArrSort.push({
            name:t('合计'),
            WinNum: ut.toNumberWithComma(winAmount).trim(),
            Proportion: '100%',
            fontB: true,
        })
    }
    return newArrSort
}

const initflag = ref(true)
const CurrencyCode = ref("")
const OperatorId = ref(0)
const defaultOperatorEvent = ref({})
const operatorListChange = (value) => {
    if (value){

        OperatorId.value = value.Id
        CurrencyCode.value = value.CurrencyKey
    }else{

        OperatorId.value = null
        CurrencyCode.value = ""
    }
    if (initflag.value){

        initSearch()
        initflag.value = false
    }
}

let initSearch = async () => {

    if (OperatorId.value == 0){
        tip.e(t("商户不能为空"))
        return
    }
    echartShow.value = false

    let [data, err] = await Client.Do(AdminGameCenter.GameBetWinDataList, {OperatorId: OperatorId.value})
    if (err) {
        return tip.e(err)
    }
    loading.value = true
    let [datas, errs] = await Client.Do(AdminGameCenter.GameList, {})
    loading.value = false
    if (errs) {
        return tip.e(errs)
    }

    if (err) {
        uiData.NowDayTableData = []
        uiData.SevenDay= {
            Bet: 0,
            Win: 0,
            BetAmount: 0,
            GameID: '',
        },
        uiData.NowMonth= {
            Bet: 0,
            Win: 0,
            BetAmount: 0,
            GameID: '',
        },
        uiData.NowDay= {
            Bet: 0,
            Win: 0,
            BetAmount: 0,
            GameID: '',
        }
        return tip.e(err)
    }
    // gameBetWinDataList.value
    echartShow.value = true
    uiData.SevenDay = returnData(data.SevenDay,datas.List)
    uiData.NowMonth = returnData(data.NowMonth,datas.List)
    uiData.NowDay = returnData(data.NowDay,datas.List)
    uiData.NowDayTableData = returnTableData(data.NowDay,datas.List)
    uiData.NowMonthTableData = returnTableData(data.NowMonth,datas.List)
    uiData.SevenDayTableData = returnTableData(data.SevenDay,datas.List)
    console.log(uiData.SevenDay,'uiData.NowDayTableData');
    let legendDataNowDay = uiData.NowDay.map(t=> {return t?.name || t('暂无')})
    let legendDataNowMonth = uiData.NowMonth.map(t=> {return t?.name || t('暂无')})
    let legendDataSevenDay = uiData.SevenDay.map(t=> {return t?.name || t('暂无')})
    nextTick(() => {
        initChart(NowDayChart, 'NowDay', changeShow.value === t('打码量')?t('打码量占比（当日）'):t('游戏盈利占比（当日）'), legendDataNowDay, uiData.NowDay)
        initChart(NowMonthChart, 'NowMonth', changeShow.value === t('打码量')?t('打码量占比（当月）'):t('游戏盈利占比（当月）'), legendDataNowMonth, uiData.NowMonth)
        initChart(SevenDayChart, 'SevenDay', changeShow.value === t('打码量')?t('打码量占比（7天）'):t('游戏盈利占比（7天）'), legendDataSevenDay, uiData.SevenDay)
    })
}

let radioChange = () => {
    initSearch()
}
let xiazai = (data,num) => {
    let type = [
        changeShow.value === t('打码量')?t('打码量占比（当日）'):t('游戏盈利占比（当日）'),
        changeShow.value === t('打码量')?t('打码量占比（7天）'):t('游戏盈利占比（7天）'),
        changeShow.value === t('打码量')?t('打码量占比（当月）'):t('游戏盈利占比（当月）'),
    ]

    if (changeShow.value === t('打码量')) {
        excel.dump(data, type[num], [
            {key: "name", name: t("游戏名称")},
            {key: "BetNum", name: t("打码量")},
            {key: "Proportion", name: t("打码量占比")},
            {key: "BetAmount", name: t("合计")},
        ])
    } else {
        excel.dump(data, type[num], [
            {key: "name", name: t("游戏名称")},
            {key: "WinNum", name: t("游戏盈利")},
            {key: "Proportion", name: t("游戏盈利占比")},
            {key: "BetAmount", name: t("合计")},
        ])
    }

}
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
            width: 60%;
            transform: translate(30%, -50%);

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
        justify-content: flex-start;
        margin-top: 20px;
        flex-direction: column;
    }
}
:deep(.el-descriptions__body){
    margin-top: 1rem;
    padding: 12px 0 0 12px;
    .el-descriptions__label:not(.is-bordered-label){
        color: rgba(18, 31, 62, 0.74);
    }
    .el-descriptions__content{
        font-weight: bold;
    }
}
</style>
